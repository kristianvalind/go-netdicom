package netdicom

import (
	"bytes"
	"fmt"
	"log"

	"github.com/kristianvalind/go-netdicom/pkg/dimse"
	"github.com/suyashkumar/dicom"
	dicomtag "github.com/suyashkumar/dicom/pkg/tag"
	dicomuid "github.com/suyashkumar/dicom/pkg/uid"
)

// Helper function used by C-{STORE,GET,MOVE} to send a dataset using C-STORE
// over an already-established association.
func runCStoreOnAssociation(upcallCh chan upcallEvent, downcallCh chan stateEvent,
	cm *contextManager,
	messageID dimse.MessageID,
	ds *dicom.Dataset) error {
	var getElement = func(tag dicomtag.Tag) (string, error) {
		elem, err := ds.FindElementByTag(tag)
		if err != nil {
			return "", fmt.Errorf("dicom.cstore: data lacks %s: %v", tag.String(), err)
		}

		elemStrings, ok := elem.Value.GetValue().([]string)
		if !ok {
			return "", fmt.Errorf("could not get string")
		}
		return elemStrings[0], nil
	}
	sopInstanceUID, err := getElement(dicomtag.MediaStorageSOPInstanceUID)
	if err != nil {
		return fmt.Errorf("dicom.cstore: data lacks SOPInstanceUID: %v", err)
	}
	sopClassUID, err := getElement(dicomtag.MediaStorageSOPClassUID)
	if err != nil {
		return fmt.Errorf("dicom.cstore: data lacks MediaStorageSOPClassUID: %v", err)
	}
	log.Printf("dicom.cstore(%s): DICOM abstractsyntax: %s, sopinstance: %s", cm.label, dicomuid.UIDString(sopClassUID), sopInstanceUID)
	context, err := cm.lookupByAbstractSyntaxUID(sopClassUID)
	if err != nil {
		log.Printf("dicom.cstore(%s): sop class %v not found in context %v", cm.label, sopClassUID, err)
		return err
	}
	log.Printf("dicom.cstore(%s): using transfersyntax %s to send sop class %s, instance %s",
		cm.label,
		dicomuid.UIDString(context.transferSyntaxUID),
		dicomuid.UIDString(sopClassUID),
		sopInstanceUID)
	bo, implicit, err := dicomuid.ParseTransferSyntaxUID(context.transferSyntaxUID)
	if err != nil {
		return fmt.Errorf("could not parse transfer syntax uid: %w", err)
	}
	b := bytes.Buffer{}
	bodyEncoder := dicom.NewWriter(&b)
	bodyEncoder.SetTransferSyntax(bo, implicit)
	for _, elem := range ds.Elements {
		if elem.Tag.Group == dicomtag.MetadataGroup {
			continue
		}
		err := bodyEncoder.WriteElement(elem)
		if err != nil {
			return fmt.Errorf("could not write metadata element: %w", err)
		}
	}
	downcallCh <- stateEvent{
		event: evt09,
		dimsePayload: &stateEventDIMSEPayload{
			abstractSyntaxName: sopClassUID,
			command: &dimse.CStoreRq{
				AffectedSOPClassUID:    sopClassUID,
				MessageID:              messageID,
				CommandDataSetType:     dimse.CommandDataSetTypeNonNull,
				AffectedSOPInstanceUID: sopInstanceUID,
			},
			data: b.Bytes(),
		},
	}
	for {
		log.Printf("dicom.cstore(%s): Start reading resp w/ messageID:%v", cm.label, messageID)
		event, ok := <-upcallCh
		if !ok {
			return fmt.Errorf("dicom.cstore(%s): Connection closed while waiting for C-STORE response", cm.label)
		}
		log.Printf("dicom.cstore(%s): resp event: %v", cm.label, event.command)
		doassert(event.eventType == upcallEventData)
		doassert(event.command != nil)
		resp, ok := event.command.(*dimse.CStoreRsp)
		doassert(ok) // TODO(saito)
		if resp.Status.Status != 0 {
			return fmt.Errorf("dicom.cstore(%s): failed: %v", cm.label, resp.String())
		}
		return nil
	}
}
