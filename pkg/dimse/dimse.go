// Package dimse implements message types defined in P3.7.
// http://dicom.nema.org/medical/dicom/current/output/pdf/part07.pdf
package dimse

//go:generate ./generate_dimse_messages.py
//go:generate stringer -type StatusCode

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"sort"

	"github.com/kristianvalind/go-netdicom/pkg/pdu"
	dicom "github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/dicomio"
	dicomtag "github.com/suyashkumar/dicom/pkg/tag"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

// Message defines the common interface for all DIMSE message types.
type Message interface {
	fmt.Stringer // Print human-readable description for debugging.
	Encode(*dicomio.Writer) error
	// GetMessageID extracts the message ID field.
	GetMessageID() MessageID
	// CommandField returns the command field value of this message.
	CommandField() int
	// GetStatus returns the the response status value. It is nil for request message
	// types, and non-nil for response message types.
	GetStatus() *Status
	// HasData is true if we expect P_DATA_TF packets after the command packets.
	HasData() bool
}

// Status represents a result of a DIMSE call.  P3.7 C defines list of status
// codes and error payloads.
type Status struct {
	// Status==StatusSuccess on success. A non-zero value on error.
	Status StatusCode

	// Optional error payloads.
	ErrorComment string // Encoded as (0000,0902)
}

// Helper class for extracting values from a list of DicomElement.
type messageDecoder struct {
	elems  []*dicom.Element
	parsed []bool // true if this element was parsed into a message field.
	err    error
}

// Find an element with the given tag. No longer checks requiredness.
func (d *messageDecoder) findElement(tag dicomtag.Tag) (*dicom.Element, error) {
	for i, elem := range d.elems {
		if elem.Tag == tag {
			log.Printf("dimse.findElement: Return %v for %s", elem, tag.String())
			d.parsed[i] = true
			return elem, nil
		}
	}
	return nil, fmt.Errorf("%w: %v", dicom.ErrorElementNotFound, dicomtag.DebugString(tag))
}

// Return the list of elements that did not match any of the prior getXXX calls.
func (d *messageDecoder) unparsedElements() (unparsed []*dicom.Element) {
	for i, parsed := range d.parsed {
		if !parsed {
			unparsed = append(unparsed, d.elems[i])
		}
	}
	return unparsed
}

func (d *messageDecoder) getStatus() (Status, error) {
	s := Status{}
	uStatus, err := d.getUInt16(dicomtag.Status)
	if err != nil {
		return Status{}, err
	}
	s.Status = StatusCode(uStatus)

	errorComment, err := d.getString(dicomtag.ErrorComment)
	if err != nil {
		// ErrorComment is optional
		if !errors.Is(err, dicom.ErrorElementNotFound) {
			return Status{}, err
		}
	}
	s.ErrorComment = errorComment
	return s, nil
}

// Find an element with "tag", and extract a string value from it. Errors are reported in d.err.
func (d *messageDecoder) getString(tag dicomtag.Tag) (string, error) {
	element, err := d.findElement(tag)
	if err != nil {
		return "", err
	}

	val := element.Value.GetValue()
	v, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("element is not string")
	}

	return v, nil
}

// Find an element with "tag", and extract a uint16 from it. Errors are reported in d.err.
func (d *messageDecoder) getUInt16(tag dicomtag.Tag) (uint16, error) {
	element, err := d.findElement(tag)
	if err != nil {
		return 0, err
	}

	val := element.Value.GetValue()
	v, ok := val.(uint16)
	if !ok {
		return 0, fmt.Errorf("element is not uint16")
	}

	return v, nil
}

// Encode the given elements. The elements are sorted in ascending tag order.
func encodeElements(w *dicomio.Writer, elems []*dicom.Element) error {
	sort.Slice(elems, func(i, j int) bool {
		return elems[i].Tag.Compare(elems[j].Tag) < 0
	})

	elementBuffer := bytes.Buffer{}

	elementWriter := dicom.NewWriter(&elementBuffer, dicom.DefaultMissingTransferSyntax())
	for _, elem := range elems {
		err := elementWriter.WriteElement(elem)
		if err != nil {
			return err
		}
	}

	err := w.WriteBytes(elementBuffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// Create a list of elements that represent the dimse status. The list contains
// multiple elements for non-ok status.
func newStatusElements(s Status) ([]*dicom.Element, error) {
	statusElement, err := dicom.NewElement(dicomtag.Status, uint16(s.Status))
	if err != nil {
		return nil, err
	}

	elems := []*dicom.Element{statusElement}
	if s.ErrorComment != "" {
		commentElement, err := dicom.NewElement(dicomtag.ErrorComment, s.ErrorComment)
		if err != nil {
			return nil, err
		}
		elems = append(elems, commentElement)
	}

	return elems, nil
}

// CommandDataSetTypeNull indicates that the DIMSE message has no data payload,
// when set in dicom.TagCommandDataSetType. Any other value indicates the
// existence of a payload.
const CommandDataSetTypeNull uint16 = 0x101

// CommandDataSetTypeNonNull indicates that the DIMSE message has a data
// payload, when set in dicom.TagCommandDataSetType.
const CommandDataSetTypeNonNull uint16 = 1

// Success is an OK status for a call.
var Success = Status{Status: StatusSuccess}

// StatusCode represents a DIMSE service response code, as defined in P3.7
type StatusCode uint16

const (
	StatusSuccess               StatusCode = 0
	StatusCancel                StatusCode = 0xFE00
	StatusSOPClassNotSupported  StatusCode = 0x0112
	StatusInvalidArgumentValue  StatusCode = 0x0115
	StatusInvalidAttributeValue StatusCode = 0x0106
	StatusInvalidObjectInstance StatusCode = 0x0117
	StatusUnrecognizedOperation StatusCode = 0x0211
	StatusNotAuthorized         StatusCode = 0x0124
	StatusPending               StatusCode = 0xff00

	// C-STORE-specific status codes. P3.4 GG4-1
	CStoreOutOfResources              StatusCode = 0xa700
	CStoreCannotUnderstand            StatusCode = 0xc000
	CStoreDataSetDoesNotMatchSOPClass StatusCode = 0xa900

	// C-FIND-specific status codes.
	CFindUnableToProcess StatusCode = 0xc000

	// C-MOVE/C-GET-specific status codes.
	CMoveOutOfResourcesUnableToCalculateNumberOfMatches StatusCode = 0xa701
	CMoveOutOfResourcesUnableToPerformSubOperations     StatusCode = 0xa702
	CMoveMoveDestinationUnknown                         StatusCode = 0xa801
	CMoveDataSetDoesNotMatchSOPClass                    StatusCode = 0xa900

	// Warning codes.
	StatusAttributeValueOutOfRange StatusCode = 0x0116
	StatusAttributeListError       StatusCode = 0x0107
)

// ReadMessage constructs a typed dimse.Message object, given a set of
// dicom.Elements,
func ReadMessage(r dicomio.Reader) (Message, error) {
	// A DIMSE message is a sequence of Elements, encoded in implicit
	// LE.
	//

	p, err := dicom.NewParser(r, r.BytesLeftUntilLimit(), nil)
	if err != nil {
		return nil, err
	}
	var elems []*dicom.Element
	for {
		elem, err := p.Next()
		if err != nil {
			if !errors.Is(err, dicom.ErrorEndOfDICOM) {
				return nil, err
			} else {
				break
			}
		}
		elems = append(elems, elem)
	}

	// Convert elems[] into a golang struct.
	dd := messageDecoder{
		elems:  elems,
		parsed: make([]bool, len(elems)),
		err:    nil,
	}
	commandField, err := dd.getUInt16(dicomtag.CommandField)
	if err != nil {
		return nil, err
	}

	v, err := decodeMessageForType(&dd, commandField)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// EncodeMessage serializes the given message.
func EncodeMessage(w *dicomio.Writer, v Message) error {
	messageBuf := &bytes.Buffer{}

	// DIMSE messages are always encoded Implicit+LE. See P3.7 6.3.1.
	subWriter := dicomio.NewWriter(messageBuf, binary.LittleEndian, true)
	err := v.Encode(&subWriter)
	if err != nil {
		return err
	}
	log.Print(messageBuf.Bytes())

	mBufLen := messageBuf.Bytes()

	elementBuffer := bytes.Buffer{}
	elementWriter := dicom.NewWriter(&elementBuffer, dicom.DefaultMissingTransferSyntax())

	element, err := dicom.NewElement(dicomtag.CommandGroupLength, uint32(len(mBufLen)))
	if err != nil {
		return err
	}

	err = elementWriter.WriteElement(element)
	if err != nil {
		return err
	}

	err = w.WriteBytes(elementBuffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// CommandAssembler is a helper that assembles a DIMSE command message and data
// payload from a sequence of P_DATA_TF PDUs.
type CommandAssembler struct {
	contextID      byte
	commandBytes   []byte
	command        Message
	dataBytes      []byte
	readAllCommand bool

	readAllData bool
}

// AddDataPDU is to be called for each P_DATA_TF PDU received from the
// network. If the fragment is marked as the last one, AddDataPDU returns
// <SOPUID, TransferSyntaxUID, payload, nil>.  If it needs more fragments, it
// returns <"", "", nil, nil>.  On error, it returns a non-nil error.
func (a *CommandAssembler) AddDataPDU(pdu *pdu.PDataTf) (byte, Message, []byte, error) {
	for _, item := range pdu.Items {
		if a.contextID == 0 {
			a.contextID = item.ContextID
		} else if a.contextID != item.ContextID {
			return 0, nil, nil, fmt.Errorf("mixed context: %d %d", a.contextID, item.ContextID)
		}
		if item.Command {
			a.commandBytes = append(a.commandBytes, item.Value...)
			if item.Last {
				if a.readAllCommand {
					return 0, nil, nil, fmt.Errorf("P_DATA_TF: found >1 command chunks with the Last bit set")
				}
				a.readAllCommand = true
			}
		} else {
			a.dataBytes = append(a.dataBytes, item.Value...)
			if item.Last {
				if a.readAllData {
					return 0, nil, nil, fmt.Errorf("P_DATA_TF: found >1 data chunks with the Last bit set")
				}
				a.readAllData = true
			}
		}
	}
	if !a.readAllCommand {
		return 0, nil, nil, nil
	}
	if a.command == nil {
		commandReader := bufio.NewReader(bytes.NewBuffer(a.commandBytes))
		d, err := dicomio.NewReader(commandReader, binary.LittleEndian, int64(len(a.commandBytes)))
		if err != nil {
			return 0, nil, nil, err
		}
		a.command, err = ReadMessage(d)
		if err != nil {
			return 0, nil, nil, err
		}
	}
	if a.command.HasData() && !a.readAllData {
		return 0, nil, nil, nil
	}
	contextID := a.contextID
	command := a.command
	dataBytes := a.dataBytes
	*a = CommandAssembler{}
	return contextID, command, dataBytes, nil
	// TODO(saito) Verify that there's no unread items after the last command&data.
}

type MessageID = uint16
