package pdu

//go:generate stringer -type AbortReasonType
//go:generate stringer -type PresentationContextResult
//go:generate stringer -type RejectReasonType
//go:generate stringer -type RejectResultType
//go:generate stringer -type SourceType
//go:generate stringer -type Type

// Implements message types defined in P3.8. It sits below the DIMSE layer.
//
// http://dicom.nema.org/medical/dicom/current/output/pdf/part08.pdf
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/suyashkumar/dicom/pkg/dicomio"
)

// PDU is the interface for DUL messages like A-ASSOCIATE-AC, P-DATA-TF.
type PDU interface {
	fmt.Stringer

	// WritePayload encodes the PDU payload. The "payload" here excludes the
	// first 6 bytes that are common to all PDU types - they are encoded in
	// EncodePDU separately.
	WritePayload(*dicomio.Writer) error
}

// Type defines type of the PDU packet.
type Type byte

const (
	TypeAAssociateRq Type = 1 // A_ASSOCIATE_RQ
	TypeAAssociateAc      = 2 // A_ASSOCIATE_AC
	TypeAAssociateRj      = 3 // A_ASSOCIATE_RJ
	TypePDataTf           = 4 // P_DATA_TF
	TypeAReleaseRq        = 5 // A_RELEASE_RQ
	TypeAReleaseRp        = 6 // A_RELEASE_RP
	TypeAAbort            = 7 // A_ABORT
)

// SubItem is the interface for DUL items, such as ApplicationContextItem and
// TransferSyntaxSubItem.
type SubItem interface {
	fmt.Stringer

	// Write serializes the item.
	Write(*dicomio.Writer) error
}

// Possible Type field values for SubItem.
const (
	ItemTypeApplicationContext           = 0x10
	ItemTypePresentationContextRequest   = 0x20
	ItemTypePresentationContextResponse  = 0x21
	ItemTypeAbstractSyntax               = 0x30
	ItemTypeTransferSyntax               = 0x40
	ItemTypeUserInformation              = 0x50
	ItemTypeUserInformationMaximumLength = 0x51
	ItemTypeImplementationClassUID       = 0x52
	ItemTypeAsynchronousOperationsWindow = 0x53
	ItemTypeRoleSelection                = 0x54
	ItemTypeImplementationVersionName    = 0x55
)

func decodeSubItem(r dicomio.Reader) (SubItem, error) {
	itemType, err := r.ReadUInt8()
	if err != nil {
		return nil, err
	}

	r.Skip(1)

	length, err := r.ReadUInt16()
	if err != nil {
		return nil, err
	}

	switch itemType {
	case ItemTypeApplicationContext:
		return decodeApplicationContextItem(r, length)
	case ItemTypeAbstractSyntax:
		return decodeAbstractSyntaxSubItem(r, length)
	case ItemTypeTransferSyntax:
		return decodeTransferSyntaxSubItem(r, length)
	case ItemTypePresentationContextRequest:
		return decodePresentationContextItem(r, itemType, length)
	case ItemTypePresentationContextResponse:
		return decodePresentationContextItem(r, itemType, length)
	case ItemTypeUserInformation:
		return decodeUserInformationItem(r, length)
	case ItemTypeUserInformationMaximumLength:
		return decodeUserInformationMaximumLengthItem(r, length)
	case ItemTypeImplementationClassUID:
		return decodeImplementationClassUIDSubItem(r, length)
	case ItemTypeAsynchronousOperationsWindow:
		return decodeAsynchronousOperationsWindowSubItem(r, length)
	case ItemTypeRoleSelection:
		return decodeRoleSelectionSubItem(r, length)
	case ItemTypeImplementationVersionName:
		return decodeImplementationVersionNameSubItem(r, length)
	default:
		return nil, fmt.Errorf("Unknown item type: 0x%x", itemType)
	}
}

func encodeSubItemHeader(w *dicomio.Writer, itemType byte, length uint16) error {
	err := w.WriteByte(itemType)
	if err != nil {
		return err
	}
	err = w.WriteZeros(1)
	if err != nil {
		return err
	}

	err = w.WriteUInt16(length)
	if err != nil {
		return err
	}

	return nil
}

// P3.8 9.3.2.3
type UserInformationItem struct {
	Items []SubItem // P3.8, Annex D.
}

func (v *UserInformationItem) Write(w *dicomio.Writer) error {
	buf := bytes.NewBuffer(nil)
	itemEncoder := dicomio.NewWriter(buf, binary.BigEndian, false)
	for _, s := range v.Items {
		err := s.Write(&itemEncoder)
		if err != nil {
			return err
		}
	}
	itemBytes := buf.Bytes()
	encodeSubItemHeader(w, ItemTypeUserInformation, uint16(len(itemBytes)))
	w.WriteBytes(itemBytes)

	return nil
}

func decodeUserInformationItem(r dicomio.Reader, length uint16) (*UserInformationItem, error) {
	v := &UserInformationItem{}
	var err error
	r.PushLimit(int64(length))
	defer r.PopLimit()
	for !r.IsLimitExhausted() {
		item, err := decodeSubItem(r)
		if err != nil {
			break
		}
		v.Items = append(v.Items, item)
	}
	return v, err
}

func (v *UserInformationItem) String() string {
	return fmt.Sprintf("UserInformationItem{items: %s}",
		subItemListString(v.Items))
}

// P3.8 D.1
type UserInformationMaximumLengthItem struct {
	MaximumLengthReceived uint32
}

func (v *UserInformationMaximumLengthItem) Write(w *dicomio.Writer) error {
	encodeSubItemHeader(w, ItemTypeUserInformationMaximumLength, 4)
	w.WriteUInt32(v.MaximumLengthReceived)

	return nil
}

func decodeUserInformationMaximumLengthItem(r dicomio.Reader, length uint16) (*UserInformationMaximumLengthItem, error) {
	if length != 4 {
		return nil, fmt.Errorf("UserInformationMaximumLengthItem must be 4 bytes, but found %dB", length)
	}
	maxLengthRecieved, err := r.ReadUInt32()
	if err != nil {
		return nil, err
	}
	return &UserInformationMaximumLengthItem{MaximumLengthReceived: maxLengthRecieved}, nil
}

func (v *UserInformationMaximumLengthItem) String() string {
	return fmt.Sprintf("UserInformationMaximumlengthItem{%d}",
		v.MaximumLengthReceived)
}

// PS3.7 Annex D.3.3.2.1
type ImplementationClassUIDSubItem subItemWithName

func decodeImplementationClassUIDSubItem(r dicomio.Reader, length uint16) (*ImplementationClassUIDSubItem, error) {
	subItemWithName, err := decodeSubItemWithName(r, length)
	if err != nil {
		return nil, err
	}
	return &ImplementationClassUIDSubItem{Name: subItemWithName}, nil
}

func (v *ImplementationClassUIDSubItem) Write(w *dicomio.Writer) error {
	encodeSubItemWithName(w, ItemTypeImplementationClassUID, v.Name)
	return nil
}

func (v *ImplementationClassUIDSubItem) String() string {
	return fmt.Sprintf("ImplementationClassUID{name: \"%s\"}", v.Name)
}

// PS3.7 Annex D.3.3.3.1
type AsynchronousOperationsWindowSubItem struct {
	MaxOpsInvoked   uint16
	MaxOpsPerformed uint16
}

func decodeAsynchronousOperationsWindowSubItem(r dicomio.Reader, length uint16) (*AsynchronousOperationsWindowSubItem, error) {
	maxOpsInvoked, err := r.ReadUInt16()
	if err != nil {
		return nil, err
	}

	maxOpsPerformed, err := r.ReadUInt16()
	if err != nil {
		return nil, err
	}
	return &AsynchronousOperationsWindowSubItem{
		MaxOpsInvoked:   maxOpsInvoked,
		MaxOpsPerformed: maxOpsPerformed,
	}, nil
}

func (v *AsynchronousOperationsWindowSubItem) Write(w *dicomio.Writer) error {
	encodeSubItemHeader(w, ItemTypeAsynchronousOperationsWindow, 2*2)
	w.WriteUInt16(v.MaxOpsInvoked)
	w.WriteUInt16(v.MaxOpsPerformed)
	return nil
}

func (v *AsynchronousOperationsWindowSubItem) String() string {
	return fmt.Sprintf("AsynchronousOpsWindow{invoked: %d performed: %d}",
		v.MaxOpsInvoked, v.MaxOpsPerformed)
}

// PS3.7 Annex D.3.3.4
type RoleSelectionSubItem struct {
	SOPClassUID string
	SCURole     byte
	SCPRole     byte
}

func decodeRoleSelectionSubItem(r dicomio.Reader, length uint16) (*RoleSelectionSubItem, error) {
	uidLen, err := r.ReadUInt16()
	if err != nil {
		return nil, err
	}
	sopClassUID, err := r.ReadString(uint32(uidLen))
	if err != nil {
		return nil, err
	}

	scuRole, err := r.ReadUInt8()
	if err != nil {
		return nil, err
	}

	scpRole, err := r.ReadUInt8()
	if err != nil {
		return nil, err
	}
	return &RoleSelectionSubItem{
		SOPClassUID: sopClassUID,
		SCURole:     scuRole,
		SCPRole:     scpRole,
	}, nil
}

func (v *RoleSelectionSubItem) Write(w *dicomio.Writer) error {
	encodeSubItemHeader(w, ItemTypeRoleSelection, uint16(2+len(v.SOPClassUID)+1*2))
	w.WriteUInt16(uint16(len(v.SOPClassUID)))
	w.WriteString(v.SOPClassUID)
	w.WriteByte(v.SCURole)
	w.WriteByte(v.SCPRole)
	return nil
}

func (v *RoleSelectionSubItem) String() string {
	return fmt.Sprintf("RoleSelection{sopclassuid: %v, scu: %v, scp: %v}", v.SOPClassUID, v.SCURole, v.SCPRole)
}

// PS3.7 Annex D.3.3.2.3
type ImplementationVersionNameSubItem subItemWithName

func decodeImplementationVersionNameSubItem(r dicomio.Reader, length uint16) (*ImplementationVersionNameSubItem, error) {
	subItemWithName, err := decodeSubItemWithName(r, length)
	if err != nil {
		return nil, err
	}
	return &ImplementationVersionNameSubItem{Name: subItemWithName}, nil
}

func (v *ImplementationVersionNameSubItem) Write(w *dicomio.Writer) error {
	encodeSubItemWithName(w, ItemTypeImplementationVersionName, v.Name)
	return nil
}

func (v *ImplementationVersionNameSubItem) String() string {
	return fmt.Sprintf("ImplementationVersionName{name: \"%s\"}", v.Name)
}

// Container for subitems that this package doesnt' support
type SubItemUnsupported struct {
	Type byte
	Data []byte
}

func (item *SubItemUnsupported) Write(w *dicomio.Writer) error {
	encodeSubItemHeader(w, item.Type, uint16(len(item.Data)))
	// TODO: handle unicode properly
	w.WriteBytes(item.Data)
	return nil
}

func (item *SubItemUnsupported) String() string {
	return fmt.Sprintf("SubitemUnsupported{type: 0x%0x data: %dbytes}",
		item.Type, len(item.Data))
}

type subItemWithName struct {
	// Type byte
	Name string
}

func encodeSubItemWithName(w *dicomio.Writer, itemType byte, name string) {
	encodeSubItemHeader(w, itemType, uint16(len(name)))
	// TODO: handle unicode properly
	w.WriteBytes([]byte(name))
}

func decodeSubItemWithName(r dicomio.Reader, length uint16) (string, error) {
	return r.ReadString(uint32(length))
}

type ApplicationContextItem subItemWithName

// The app context for DICOM. The first item in the A-ASSOCIATE-RQ
const DICOMApplicationContextItemName = "1.2.840.10008.3.1.1.1"

func decodeApplicationContextItem(r dicomio.Reader, length uint16) (*ApplicationContextItem, error) {
	subItemWithName, err := decodeSubItemWithName(r, length)
	if err != nil {
		return nil, err
	}
	return &ApplicationContextItem{Name: subItemWithName}, nil
}

func (v *ApplicationContextItem) Write(w *dicomio.Writer) error {
	encodeSubItemWithName(w, ItemTypeApplicationContext, v.Name)
	return nil
}

func (v *ApplicationContextItem) String() string {
	return fmt.Sprintf("ApplicationContext{name: \"%s\"}", v.Name)
}

type AbstractSyntaxSubItem subItemWithName

func decodeAbstractSyntaxSubItem(r dicomio.Reader, length uint16) (*AbstractSyntaxSubItem, error) {
	subItemWithName, err := decodeSubItemWithName(r, length)
	if err != nil {
		return nil, err
	}
	return &AbstractSyntaxSubItem{Name: subItemWithName}, nil
}

func (v *AbstractSyntaxSubItem) Write(w *dicomio.Writer) error {
	encodeSubItemWithName(w, ItemTypeAbstractSyntax, v.Name)
	return nil
}

func (v *AbstractSyntaxSubItem) String() string {
	return fmt.Sprintf("AbstractSyntax{name: \"%s\"}", v.Name)
}

type TransferSyntaxSubItem subItemWithName

func decodeTransferSyntaxSubItem(r dicomio.Reader, length uint16) (*TransferSyntaxSubItem, error) {
	subItemWithName, err := decodeSubItemWithName(r, length)
	if err != nil {
		return nil, err
	}
	return &TransferSyntaxSubItem{Name: subItemWithName}, nil
}

func (v *TransferSyntaxSubItem) Write(w *dicomio.Writer) error {
	encodeSubItemWithName(w, ItemTypeTransferSyntax, v.Name)
}

func (v *TransferSyntaxSubItem) String() string {
	return fmt.Sprintf("TransferSyntax{name: \"%s\"}", v.Name)
}

// Result of abstractsyntax/transfersyntax handshake during A-ACCEPT.  P3.8,
// 90.3.3.2, table 9-18.
type PresentationContextResult byte

const (
	PresentationContextAccepted                                    PresentationContextResult = 0
	PresentationContextUserRejection                               PresentationContextResult = 1
	PresentationContextProviderRejectionNoReason                   PresentationContextResult = 2
	PresentationContextProviderRejectionAbstractSyntaxNotSupported PresentationContextResult = 3
	PresentationContextProviderRejectionTransferSyntaxNotSupported PresentationContextResult = 4
)

// P3.8 9.3.2.2, 9.3.3.2
type PresentationContextItem struct {
	Type      byte // ItemTypePresentationContext*
	ContextID byte
	// 1 byte reserved

	// Result is meaningful iff Type=0x21, zero else.
	Result PresentationContextResult

	// 1 byte reserved
	Items []SubItem // List of {Abstract,Transfer}SyntaxSubItem
}

func decodePresentationContextItem(r dicomio.Reader, itemType byte, length uint16) (*PresentationContextItem, error) {
	v := &PresentationContextItem{Type: itemType}
	r.PushLimit(int64(length))
	defer r.PopLimit()
	var err error
	v.ContextID, err = r.ReadUInt8()
	if err != nil {
		return nil, err
	}
	r.Skip(1)
	result, err := r.ReadUInt8()
	if err != nil {
		return nil, err
	}
	v.Result = PresentationContextResult(result)
	r.Skip(1)
	for !r.IsLimitExhausted() {
		item, err := decodeSubItem(r)
		if err != nil {
			return nil, err
		}
		v.Items = append(v.Items, item)
	}
	if v.ContextID%2 != 1 {
		return nil, fmt.Errorf("PresentationContextItem ID must be odd, but found %x", v.ContextID)
	}
	return v, nil
}

func (v *PresentationContextItem) Write(w *dicomio.Writer) error {
	if v.Type != ItemTypePresentationContextRequest &&
		v.Type != ItemTypePresentationContextResponse {
		panic(*v)
	}

	itemEncoder := dicomio.NewBytesEncoder(binary.BigEndian, dicomio.UnknownVR)
	for _, s := range v.Items {
		s.Write(itemEncoder)
	}
	if err := itemEncoder.Error(); err != nil {
		e.SetError(err)
		return
	}
	itemBytes := itemEncoder.Bytes()
	encodeSubItemHeader(w, v.Type, uint16(4+len(itemBytes)))
	w.WriteByte(v.ContextID)
	w.WriteZeros(3)
	w.WriteBytes(itemBytes)
}

func (v *PresentationContextItem) String() string {
	itemType := "rq"
	if v.Type == ItemTypePresentationContextResponse {
		itemType = "ac"
	}
	return fmt.Sprintf("PresentationContext%s{id: %d result: %d, items:%s}",
		itemType, v.ContextID, v.Result, subItemListString(v.Items))
}

// P3.8 9.3.2.2.1 & 9.3.2.2.2
type PresentationDataValueItem struct {
	// Length: 2 + len(Value)
	ContextID byte

	// P3.8, E.2: the following two fields encode a single byte.
	Command bool // Bit 7 (LSB): 1 means command 0 means data
	Last    bool // Bit 6: 1 means last fragment. 0 means not last fragment.

	// Payload, either command or data
	Value []byte
}

func ReadPresentationDataValueItem(r dicomio.Reader) (PresentationDataValueItem, error) {
	item := PresentationDataValueItem{}
	length, err := r.ReadUInt32()
	if err != nil {
		return PresentationDataValueItem{}, err
	}
	item.ContextID, err = r.ReadUInt8()
	if err != nil {
		return PresentationDataValueItem{}, err
	}
	header, err := r.ReadUInt8()
	if err != nil {
		return PresentationDataValueItem{}, err
	}
	item.Command = (header&1 != 0)
	item.Last = (header&2 != 0)

	valueBytes := make([]byte, length-2)
	_, err = r.Read(valueBytes)
	if err != nil {
		return PresentationDataValueItem{}, err
	}

	item.Value = valueBytes // remove contextID and header
	return item, nil
}

func (v *PresentationDataValueItem) Write(w *dicomio.Writer) error {
	var header byte
	if v.Command {
		header |= 1
	}
	if v.Last {
		header |= 2
	}
	w.WriteUInt32(uint32(2 + len(v.Value)))
	w.WriteByte(v.ContextID)
	w.WriteByte(header)
	w.WriteBytes(v.Value)
}

func (v *PresentationDataValueItem) String() string {
	return fmt.Sprintf("PresentationDataValue{context: %d, cmd:%v last:%v value: %d bytes}", v.ContextID, v.Command, v.Last, len(v.Value))
}

// EncodePDU serializes "pdu" into []byte.
func EncodePDU(pdu PDU) ([]byte, error) {
	var pduType Type
	switch n := pdu.(type) {
	case *AAssociate:
		pduType = n.Type
	case *AAssociateRj:
		pduType = TypeAAssociateRj
	case *PDataTf:
		pduType = TypePDataTf
	case *AReleaseRq:
		pduType = TypeAReleaseRq
	case *AReleaseRp:
		pduType = TypeAReleaseRp
	case *AAbort:
		pduType = TypeAAbort
	default:
		panic(fmt.Sprintf("Unknown PDU %v", pdu))
	}
	e := dicomio.NewBytesEncoder(binary.BigEndian, dicomio.UnknownVR)
	pdu.WritePayload(e)
	if err := e.Error(); err != nil {
		return nil, err
	}
	payload := e.Bytes()
	// Reserve the header bytes. It will be filled in Finish.
	var header [6]byte // First 6 bytes of buf.
	header[0] = byte(pduType)
	header[1] = 0 // Reserved.
	binary.BigEndian.PutUint32(header[2:6], uint32(len(payload)))
	return append(header[:], payload...), nil
}

// EncodePDU reads a "pdu" from a stream. maxPDUSize defines the maximum
// possible PDU size, in bytes, accepted by the caller.
func ReadPDU(in io.Reader, maxPDUSize int) (PDU, error) {
	var pduType Type
	var skip byte
	var length uint32
	err := binary.Read(in, binary.BigEndian, &pduType)
	if err != nil {
		return nil, err
	}
	err = binary.Read(in, binary.BigEndian, &skip)
	if err != nil {
		return nil, err
	}
	err = binary.Read(in, binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}
	if length >= uint32(maxPDUSize)*2 {
		// Avoid using too much memory. *2 is just an arbitrary slack.
		return nil, fmt.Errorf("Invalid length %d; it's much larger than max PDU size of %d", length, maxPDUSize)
	}
	d := dicomio.NewDecoder(
		&io.LimitedReader{R: in, N: int64(length)},
		binary.BigEndian,  // PDU is always big endian
		dicomio.UnknownVR) // irrelevant for PDU parsing
	var pdu PDU
	switch pduType {
	case TypeAAssociateRq:
		fallthrough
	case TypeAAssociateAc:
		pdu, err = decodeAAssociate(d, pduType)
		if err != nil {
			return nil, err
		}
	case TypeAAssociateRj:
		pdu, err = decodeAAssociateRj(d)
		if err != nil {
			return nil, err
		}
	case TypeAAbort:
		pdu, err = decodeAAbort(d)
		if err != nil {
			return nil, err
		}
	case TypePDataTf:
		pdu, err = decodePDataTf(d)
		if err != nil {
			return nil, err
		}
	case TypeAReleaseRq:
		pdu = decodeAReleaseRq(d)
	case TypeAReleaseRp:
		pdu = decodeAReleaseRp(d)
	}
	if pdu == nil {
		err := fmt.Errorf("ReadPDU: unknown message type %d", pduType)
		return nil, err
	}
	if err := d.Finish(); err != nil {
		return nil, err
	}
	return pdu, nil
}

type AReleaseRq struct {
}

func decodeAReleaseRq(r dicomio.Reader) *AReleaseRq {
	pdu := &AReleaseRq{}
	r.Skip(4)
	return pdu
}

func (pdu *AReleaseRq) WritePayload(w *dicomio.Writer) error {
	w.WriteZeros(4)
}

func (pdu *AReleaseRq) String() string {
	return fmt.Sprintf("A_RELEASE_RQ(%v)", *pdu)
}

type AReleaseRp struct {
}

func decodeAReleaseRp(r dicomio.Reader) *AReleaseRp {
	pdu := &AReleaseRp{}
	r.Skip(4)
	return pdu
}

func (pdu *AReleaseRp) WritePayload(w *dicomio.Writer) error {
	w.WriteZeros(4)
}

func (pdu *AReleaseRp) String() string {
	return fmt.Sprintf("A_RELEASE_RP(%v)", *pdu)
}

func subItemListString(items []SubItem) string {
	buf := bytes.Buffer{}
	buf.WriteString("[")
	for i, subitem := range items {
		if i > 0 {
			buf.WriteString("\n")
		}
		buf.WriteString(subitem.String())
	}
	buf.WriteString("]")
	return buf.String()
}

const CurrentProtocolVersion uint16 = 1

// Defines A_ASSOCIATE_{RQ,AC}. P3.8 9.3.2 and 9.3.3
type AAssociate struct {
	Type            Type // One of {TypeA_Associate_RQ,TypeA_Associate_AC}
	ProtocolVersion uint16
	// Reserved uint16
	CalledAETitle  string // For .._AC, the value is copied from A_ASSOCIATE_RQ
	CallingAETitle string // For .._AC, the value is copied from A_ASSOCIATE_RQ
	Items          []SubItem
}

func decodeAAssociate(r dicomio.Reader, pduType Type) (*AAssociate, error) {
	pdu := &AAssociate{}
	pdu.Type = pduType
	var err error
	pdu.ProtocolVersion, err = r.ReadUInt16()
	if err != nil {
		return nil, err
	}
	r.Skip(2) // Reserved
	pdu.CalledAETitle, err = r.ReadString(16)
	if err != nil {
		return nil, err
	}
	pdu.CallingAETitle, err = r.ReadString(16)
	if err != nil {
		return nil, err
	}
	r.Skip(8 * 4)
	for !r.IsLimitExhausted() {
		item, err := decodeSubItem(r)
		if err != nil {
			return nil, err
		}
		pdu.Items = append(pdu.Items, item)
	}
	if pdu.CalledAETitle == "" || pdu.CallingAETitle == "" {
		return nil, fmt.Errorf("A_ASSOCIATE.{Called,Calling}AETitle must not be empty, in %v", pdu.String())
	}
	return pdu, nil
}

func (pdu *AAssociate) WritePayload(w *dicomio.Writer) error {
	if pdu.Type == 0 || pdu.CalledAETitle == "" || pdu.CallingAETitle == "" {
		panic(*pdu)
	}
	w.WriteUInt16(pdu.ProtocolVersion)
	w.WriteZeros(2) // Reserved
	w.WriteString(fillString(pdu.CalledAETitle, 16))
	w.WriteString(fillString(pdu.CallingAETitle, 16))
	w.WriteZeros(8 * 4)
	for _, item := range pdu.Items {
		item.Write(w)
	}
}

func (pdu *AAssociate) String() string {
	name := "AC"
	if pdu.Type == TypeAAssociateRq {
		name = "RQ"
	}
	return fmt.Sprintf("A_ASSOCIATE_%s{version:%v called:'%v' calling:'%v' items:%s}",
		name, pdu.ProtocolVersion,
		pdu.CalledAETitle, pdu.CallingAETitle, subItemListString(pdu.Items))
}

// P3.8 9.3.4
type AAssociateRj struct {
	Result RejectResultType
	Source SourceType
	Reason RejectReasonType
}

// Possible values for AAssociateRj.Result
type RejectResultType byte

const (
	ResultRejectedPermanent RejectResultType = 1
	ResultRejectedTransient RejectResultType = 2
)

// Possible values for AAssociateRj.Reason
type RejectReasonType byte

const (
	RejectReasonNone                               RejectReasonType = 1
	RejectReasonApplicationContextNameNotSupported RejectReasonType = 2
	RejectReasonCallingAETitleNotRecognized        RejectReasonType = 3
	RejectReasonCalledAETitleNotRecognized         RejectReasonType = 7
)

// Possible values for AAssociateRj.Source
type SourceType byte

const (
	SourceULServiceUser                 SourceType = 1
	SourceULServiceProviderACSE         SourceType = 2
	SourceULServiceProviderPresentation SourceType = 3
)

func decodeAAssociateRj(r dicomio.Reader) (*AAssociateRj, error) {
	pdu := &AAssociateRj{}
	r.Skip(1) // reserved
	result, err := r.ReadUInt8()
	if err != nil {
		return nil, err
	}

	source, err := r.ReadUInt8()
	if err != nil {
		return nil, err
	}

	reason, err := r.ReadUInt8()
	if err != nil {
		return nil, err
	}
	pdu.Result = RejectResultType(result)
	pdu.Source = SourceType(source)
	pdu.Reason = RejectReasonType(reason)
	return pdu, nil
}

func (pdu *AAssociateRj) WritePayload(w *dicomio.Writer) error {
	w.WriteZeros(1)
	w.WriteByte(byte(pdu.Result))
	w.WriteByte(byte(pdu.Source))
	w.WriteByte(byte(pdu.Reason))
}

func (pdu *AAssociateRj) String() string {
	return fmt.Sprintf("A_ASSOCIATE_RJ{result: %v, source: %v, reason: %v}", pdu.Result, pdu.Source, pdu.Reason)
}

type AbortReasonType byte

const (
	AbortReasonNotSpecified             AbortReasonType = 0
	AbortReasonUnexpectedPDU            AbortReasonType = 2
	AbortReasonUnrecognizedPDUParameter AbortReasonType = 3
	AbortReasonUnexpectedPDUParameter   AbortReasonType = 4
	AbortReasonInvalidPDUParameterValue AbortReasonType = 5
)

type AAbort struct {
	Source SourceType
	Reason AbortReasonType
}

func decodeAAbort(r dicomio.Reader) (*AAbort, error) {
	pdu := &AAbort{}
	r.Skip(2)
	source, err := r.ReadUInt8()
	if err != nil {
		return nil, err
	}
	pdu.Source = SourceType(source)

	reason, err := r.ReadUInt8()
	if err != nil {
		return nil, err
	}
	pdu.Reason = AbortReasonType(reason)
	return pdu, nil
}

func (pdu *AAbort) WritePayload(w *dicomio.Writer) error {
	err := w.WriteZeros(2)
	if err != nil {
		return err
	}
	err = w.WriteByte(byte(pdu.Source))
	if err != nil {
		return err
	}
	err = w.WriteByte(byte(pdu.Reason))
	if err != nil {
		return err
	}
	return nil
}

func (pdu *AAbort) String() string {
	return fmt.Sprintf("A_ABORT{source:%v reason:%v}", pdu.Source, pdu.Reason)
}

type PDataTf struct {
	Items []PresentationDataValueItem
}

func decodePDataTf(r dicomio.Reader) (*PDataTf, error) {
	pdu := &PDataTf{}
	for !r.IsLimitExhausted() {
		item, err := ReadPresentationDataValueItem(r)
		if err != nil {
			return nil, err
		}
		pdu.Items = append(pdu.Items, item)
	}
	return pdu, nil
}

func (pdu *PDataTf) WritePayload(w *dicomio.Writer) error {
	var err error
	for _, item := range pdu.Items {
		err = item.Write(w)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pdu *PDataTf) String() string {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("P_DATA_TF{items: ["))
	for i, item := range pdu.Items {
		if i > 0 {
			buf.WriteString("\n")
		}
		buf.WriteString(item.String())
	}
	buf.WriteString("]}")
	return buf.String()
}

// fillString pads the string with " " up to the given length.
func fillString(v string, length int) string {
	if len(v) > length {
		return v[:16]
	}
	for len(v) < length {
		v += " "
	}
	return v
}
