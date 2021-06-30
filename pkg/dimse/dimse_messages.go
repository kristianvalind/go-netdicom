
package dimse

// Code generated from generate_dimse_messages.py. DO NOT EDIT.

import (
    "errors"
	"fmt"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/dicomio"
	dicomtag "github.com/suyashkumar/dicom/pkg/tag"
)

        
type CStoreRq struct {
	AffectedSOPClassUID string
	MessageID MessageID
	Priority uint16
	CommandDataSetType uint16
	AffectedSOPInstanceUID string
	MoveOriginatorApplicationEntityTitle string
	MoveOriginatorMessageID MessageID
	Extra []*dicom.Element  // Unparsed elements
}

func (v *CStoreRq) Encode(w *dicomio.Writer) error {
    elems := []*dicom.Element{}
    elem, err := dicom.NewElement(dicomtag.CommandField, uint16(1))
    if err != nil {
        return err
    }
    elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.AffectedSOPClassUID, v.AffectedSOPClassUID)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.MessageID, v.MessageID)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.Priority, v.Priority)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.CommandDataSetType, v.CommandDataSetType)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.AffectedSOPInstanceUID, v.AffectedSOPInstanceUID)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
	if v.MoveOriginatorApplicationEntityTitle != "" {
		elem, err = dicom.NewElement(dicomtag.MoveOriginatorApplicationEntityTitle, v.MoveOriginatorApplicationEntityTitle)
        if err != nil {
            return err
        }
		elems = append(elems, elem)
	}
	if v.MoveOriginatorMessageID != 0 {
		elem, err = dicom.NewElement(dicomtag.MoveOriginatorMessageID, v.MoveOriginatorMessageID)
        if err != nil {
            return err
        }
		elems = append(elems, elem)
	}
    elems = append(elems, v.Extra...)
    return encodeElements(w, elems)
}

func (v *CStoreRq) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v *CStoreRq) CommandField() int {
	return 1
}

func (v *CStoreRq) GetMessageID() MessageID {
	return v.MessageID
}

func (v *CStoreRq) GetStatus() *Status {
	return nil
}

func (v *CStoreRq) String() string {
	return fmt.Sprintf("CStoreRq{AffectedSOPClassUID:%v MessageID:%v Priority:%v CommandDataSetType:%v AffectedSOPInstanceUID:%v MoveOriginatorApplicationEntityTitle:%v MoveOriginatorMessageID:%v}}", v.AffectedSOPClassUID, v.MessageID, v.Priority, v.CommandDataSetType, v.AffectedSOPInstanceUID, v.MoveOriginatorApplicationEntityTitle, v.MoveOriginatorMessageID)
}

func decodeCStoreRq(d *messageDecoder) (*CStoreRq, error) {
	v := &CStoreRq{}
	var err error
	v.AffectedSOPClassUID, err = d.getString(dicomtag.AffectedSOPClassUID)
    if err != nil {
        return nil, err
    }
	v.MessageID, err = d.getUInt16(dicomtag.MessageID)
    if err != nil {
        return nil, err
    }
	v.Priority, err = d.getUInt16(dicomtag.Priority)
    if err != nil {
        return nil, err
    }
	v.CommandDataSetType, err = d.getUInt16(dicomtag.CommandDataSetType)
    if err != nil {
        return nil, err
    }
	v.AffectedSOPInstanceUID, err = d.getString(dicomtag.AffectedSOPInstanceUID)
    if err != nil {
        return nil, err
    }
	v.MoveOriginatorApplicationEntityTitle, err = d.getString(dicomtag.MoveOriginatorApplicationEntityTitle)
    if err != nil {
        if !errors.Is(err, dicom.ErrorElementNotFound) {
            return nil, err
        }
    }
	v.MoveOriginatorMessageID, err = d.getUInt16(dicomtag.MoveOriginatorMessageID)
    if err != nil {
        if !errors.Is(err, dicom.ErrorElementNotFound) {
            return nil, err
        }
    }
	v.Extra = d.unparsedElements()
	return v, nil
}

type CStoreRsp struct {
	AffectedSOPClassUID string
	MessageIDBeingRespondedTo MessageID
	CommandDataSetType uint16
	AffectedSOPInstanceUID string
	Status Status
	Extra []*dicom.Element  // Unparsed elements
}

func (v *CStoreRsp) Encode(w *dicomio.Writer) error {
    elems := []*dicom.Element{}
    elem, err := dicom.NewElement(dicomtag.CommandField, uint16(32769))
    if err != nil {
        return err
    }
    elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.AffectedSOPClassUID, v.AffectedSOPClassUID)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.MessageIDBeingRespondedTo, v.MessageIDBeingRespondedTo)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.CommandDataSetType, v.CommandDataSetType)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.AffectedSOPInstanceUID, v.AffectedSOPInstanceUID)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
	statusElems, err := newStatusElements(v.Status)
    if err != nil {
        return err
    }
	elems = append(elems, statusElems...)
    elems = append(elems, v.Extra...)
    return encodeElements(w, elems)
}

func (v *CStoreRsp) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v *CStoreRsp) CommandField() int {
	return 32769
}

func (v *CStoreRsp) GetMessageID() MessageID {
	return v.MessageIDBeingRespondedTo
}

func (v *CStoreRsp) GetStatus() *Status {
	return &v.Status
}

func (v *CStoreRsp) String() string {
	return fmt.Sprintf("CStoreRsp{AffectedSOPClassUID:%v MessageIDBeingRespondedTo:%v CommandDataSetType:%v AffectedSOPInstanceUID:%v Status:%v}}", v.AffectedSOPClassUID, v.MessageIDBeingRespondedTo, v.CommandDataSetType, v.AffectedSOPInstanceUID, v.Status)
}

func decodeCStoreRsp(d *messageDecoder) (*CStoreRsp, error) {
	v := &CStoreRsp{}
	var err error
	v.AffectedSOPClassUID, err = d.getString(dicomtag.AffectedSOPClassUID)
    if err != nil {
        return nil, err
    }
	v.MessageIDBeingRespondedTo, err = d.getUInt16(dicomtag.MessageIDBeingRespondedTo)
    if err != nil {
        return nil, err
    }
	v.CommandDataSetType, err = d.getUInt16(dicomtag.CommandDataSetType)
    if err != nil {
        return nil, err
    }
	v.AffectedSOPInstanceUID, err = d.getString(dicomtag.AffectedSOPInstanceUID)
    if err != nil {
        return nil, err
    }
	v.Status, err = d.getStatus()
    if err != nil {
        return nil, err
    }
	v.Extra = d.unparsedElements()
	return v, nil
}

type CFindRq struct {
	AffectedSOPClassUID string
	MessageID MessageID
	Priority uint16
	CommandDataSetType uint16
	Extra []*dicom.Element  // Unparsed elements
}

func (v *CFindRq) Encode(w *dicomio.Writer) error {
    elems := []*dicom.Element{}
    elem, err := dicom.NewElement(dicomtag.CommandField, uint16(32))
    if err != nil {
        return err
    }
    elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.AffectedSOPClassUID, v.AffectedSOPClassUID)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.MessageID, v.MessageID)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.Priority, v.Priority)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.CommandDataSetType, v.CommandDataSetType)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elems = append(elems, v.Extra...)
    return encodeElements(w, elems)
}

func (v *CFindRq) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v *CFindRq) CommandField() int {
	return 32
}

func (v *CFindRq) GetMessageID() MessageID {
	return v.MessageID
}

func (v *CFindRq) GetStatus() *Status {
	return nil
}

func (v *CFindRq) String() string {
	return fmt.Sprintf("CFindRq{AffectedSOPClassUID:%v MessageID:%v Priority:%v CommandDataSetType:%v}}", v.AffectedSOPClassUID, v.MessageID, v.Priority, v.CommandDataSetType)
}

func decodeCFindRq(d *messageDecoder) (*CFindRq, error) {
	v := &CFindRq{}
	var err error
	v.AffectedSOPClassUID, err = d.getString(dicomtag.AffectedSOPClassUID)
    if err != nil {
        return nil, err
    }
	v.MessageID, err = d.getUInt16(dicomtag.MessageID)
    if err != nil {
        return nil, err
    }
	v.Priority, err = d.getUInt16(dicomtag.Priority)
    if err != nil {
        return nil, err
    }
	v.CommandDataSetType, err = d.getUInt16(dicomtag.CommandDataSetType)
    if err != nil {
        return nil, err
    }
	v.Extra = d.unparsedElements()
	return v, nil
}

type CFindRsp struct {
	AffectedSOPClassUID string
	MessageIDBeingRespondedTo MessageID
	CommandDataSetType uint16
	Status Status
	Extra []*dicom.Element  // Unparsed elements
}

func (v *CFindRsp) Encode(w *dicomio.Writer) error {
    elems := []*dicom.Element{}
    elem, err := dicom.NewElement(dicomtag.CommandField, uint16(32800))
    if err != nil {
        return err
    }
    elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.AffectedSOPClassUID, v.AffectedSOPClassUID)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.MessageIDBeingRespondedTo, v.MessageIDBeingRespondedTo)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.CommandDataSetType, v.CommandDataSetType)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
	statusElems, err := newStatusElements(v.Status)
    if err != nil {
        return err
    }
	elems = append(elems, statusElems...)
    elems = append(elems, v.Extra...)
    return encodeElements(w, elems)
}

func (v *CFindRsp) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v *CFindRsp) CommandField() int {
	return 32800
}

func (v *CFindRsp) GetMessageID() MessageID {
	return v.MessageIDBeingRespondedTo
}

func (v *CFindRsp) GetStatus() *Status {
	return &v.Status
}

func (v *CFindRsp) String() string {
	return fmt.Sprintf("CFindRsp{AffectedSOPClassUID:%v MessageIDBeingRespondedTo:%v CommandDataSetType:%v Status:%v}}", v.AffectedSOPClassUID, v.MessageIDBeingRespondedTo, v.CommandDataSetType, v.Status)
}

func decodeCFindRsp(d *messageDecoder) (*CFindRsp, error) {
	v := &CFindRsp{}
	var err error
	v.AffectedSOPClassUID, err = d.getString(dicomtag.AffectedSOPClassUID)
    if err != nil {
        return nil, err
    }
	v.MessageIDBeingRespondedTo, err = d.getUInt16(dicomtag.MessageIDBeingRespondedTo)
    if err != nil {
        return nil, err
    }
	v.CommandDataSetType, err = d.getUInt16(dicomtag.CommandDataSetType)
    if err != nil {
        return nil, err
    }
	v.Status, err = d.getStatus()
    if err != nil {
        return nil, err
    }
	v.Extra = d.unparsedElements()
	return v, nil
}

type CGetRq struct {
	AffectedSOPClassUID string
	MessageID MessageID
	Priority uint16
	CommandDataSetType uint16
	Extra []*dicom.Element  // Unparsed elements
}

func (v *CGetRq) Encode(w *dicomio.Writer) error {
    elems := []*dicom.Element{}
    elem, err := dicom.NewElement(dicomtag.CommandField, uint16(16))
    if err != nil {
        return err
    }
    elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.AffectedSOPClassUID, v.AffectedSOPClassUID)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.MessageID, v.MessageID)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.Priority, v.Priority)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.CommandDataSetType, v.CommandDataSetType)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elems = append(elems, v.Extra...)
    return encodeElements(w, elems)
}

func (v *CGetRq) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v *CGetRq) CommandField() int {
	return 16
}

func (v *CGetRq) GetMessageID() MessageID {
	return v.MessageID
}

func (v *CGetRq) GetStatus() *Status {
	return nil
}

func (v *CGetRq) String() string {
	return fmt.Sprintf("CGetRq{AffectedSOPClassUID:%v MessageID:%v Priority:%v CommandDataSetType:%v}}", v.AffectedSOPClassUID, v.MessageID, v.Priority, v.CommandDataSetType)
}

func decodeCGetRq(d *messageDecoder) (*CGetRq, error) {
	v := &CGetRq{}
	var err error
	v.AffectedSOPClassUID, err = d.getString(dicomtag.AffectedSOPClassUID)
    if err != nil {
        return nil, err
    }
	v.MessageID, err = d.getUInt16(dicomtag.MessageID)
    if err != nil {
        return nil, err
    }
	v.Priority, err = d.getUInt16(dicomtag.Priority)
    if err != nil {
        return nil, err
    }
	v.CommandDataSetType, err = d.getUInt16(dicomtag.CommandDataSetType)
    if err != nil {
        return nil, err
    }
	v.Extra = d.unparsedElements()
	return v, nil
}

type CGetRsp struct {
	AffectedSOPClassUID string
	MessageIDBeingRespondedTo MessageID
	CommandDataSetType uint16
	NumberOfRemainingSuboperations uint16
	NumberOfCompletedSuboperations uint16
	NumberOfFailedSuboperations uint16
	NumberOfWarningSuboperations uint16
	Status Status
	Extra []*dicom.Element  // Unparsed elements
}

func (v *CGetRsp) Encode(w *dicomio.Writer) error {
    elems := []*dicom.Element{}
    elem, err := dicom.NewElement(dicomtag.CommandField, uint16(32784))
    if err != nil {
        return err
    }
    elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.AffectedSOPClassUID, v.AffectedSOPClassUID)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.MessageIDBeingRespondedTo, v.MessageIDBeingRespondedTo)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.CommandDataSetType, v.CommandDataSetType)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
	if v.NumberOfRemainingSuboperations != 0 {
		elem, err = dicom.NewElement(dicomtag.NumberOfRemainingSuboperations, v.NumberOfRemainingSuboperations)
        if err != nil {
            return err
        }
		elems = append(elems, elem)
	}
	if v.NumberOfCompletedSuboperations != 0 {
		elem, err = dicom.NewElement(dicomtag.NumberOfCompletedSuboperations, v.NumberOfCompletedSuboperations)
        if err != nil {
            return err
        }
		elems = append(elems, elem)
	}
	if v.NumberOfFailedSuboperations != 0 {
		elem, err = dicom.NewElement(dicomtag.NumberOfFailedSuboperations, v.NumberOfFailedSuboperations)
        if err != nil {
            return err
        }
		elems = append(elems, elem)
	}
	if v.NumberOfWarningSuboperations != 0 {
		elem, err = dicom.NewElement(dicomtag.NumberOfWarningSuboperations, v.NumberOfWarningSuboperations)
        if err != nil {
            return err
        }
		elems = append(elems, elem)
	}
	statusElems, err := newStatusElements(v.Status)
    if err != nil {
        return err
    }
	elems = append(elems, statusElems...)
    elems = append(elems, v.Extra...)
    return encodeElements(w, elems)
}

func (v *CGetRsp) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v *CGetRsp) CommandField() int {
	return 32784
}

func (v *CGetRsp) GetMessageID() MessageID {
	return v.MessageIDBeingRespondedTo
}

func (v *CGetRsp) GetStatus() *Status {
	return &v.Status
}

func (v *CGetRsp) String() string {
	return fmt.Sprintf("CGetRsp{AffectedSOPClassUID:%v MessageIDBeingRespondedTo:%v CommandDataSetType:%v NumberOfRemainingSuboperations:%v NumberOfCompletedSuboperations:%v NumberOfFailedSuboperations:%v NumberOfWarningSuboperations:%v Status:%v}}", v.AffectedSOPClassUID, v.MessageIDBeingRespondedTo, v.CommandDataSetType, v.NumberOfRemainingSuboperations, v.NumberOfCompletedSuboperations, v.NumberOfFailedSuboperations, v.NumberOfWarningSuboperations, v.Status)
}

func decodeCGetRsp(d *messageDecoder) (*CGetRsp, error) {
	v := &CGetRsp{}
	var err error
	v.AffectedSOPClassUID, err = d.getString(dicomtag.AffectedSOPClassUID)
    if err != nil {
        return nil, err
    }
	v.MessageIDBeingRespondedTo, err = d.getUInt16(dicomtag.MessageIDBeingRespondedTo)
    if err != nil {
        return nil, err
    }
	v.CommandDataSetType, err = d.getUInt16(dicomtag.CommandDataSetType)
    if err != nil {
        return nil, err
    }
	v.NumberOfRemainingSuboperations, err = d.getUInt16(dicomtag.NumberOfRemainingSuboperations)
    if err != nil {
        if !errors.Is(err, dicom.ErrorElementNotFound) {
            return nil, err
        }
    }
	v.NumberOfCompletedSuboperations, err = d.getUInt16(dicomtag.NumberOfCompletedSuboperations)
    if err != nil {
        if !errors.Is(err, dicom.ErrorElementNotFound) {
            return nil, err
        }
    }
	v.NumberOfFailedSuboperations, err = d.getUInt16(dicomtag.NumberOfFailedSuboperations)
    if err != nil {
        if !errors.Is(err, dicom.ErrorElementNotFound) {
            return nil, err
        }
    }
	v.NumberOfWarningSuboperations, err = d.getUInt16(dicomtag.NumberOfWarningSuboperations)
    if err != nil {
        if !errors.Is(err, dicom.ErrorElementNotFound) {
            return nil, err
        }
    }
	v.Status, err = d.getStatus()
    if err != nil {
        return nil, err
    }
	v.Extra = d.unparsedElements()
	return v, nil
}

type CMoveRq struct {
	AffectedSOPClassUID string
	MessageID MessageID
	Priority uint16
	MoveDestination string
	CommandDataSetType uint16
	Extra []*dicom.Element  // Unparsed elements
}

func (v *CMoveRq) Encode(w *dicomio.Writer) error {
    elems := []*dicom.Element{}
    elem, err := dicom.NewElement(dicomtag.CommandField, uint16(33))
    if err != nil {
        return err
    }
    elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.AffectedSOPClassUID, v.AffectedSOPClassUID)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.MessageID, v.MessageID)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.Priority, v.Priority)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.MoveDestination, v.MoveDestination)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.CommandDataSetType, v.CommandDataSetType)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elems = append(elems, v.Extra...)
    return encodeElements(w, elems)
}

func (v *CMoveRq) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v *CMoveRq) CommandField() int {
	return 33
}

func (v *CMoveRq) GetMessageID() MessageID {
	return v.MessageID
}

func (v *CMoveRq) GetStatus() *Status {
	return nil
}

func (v *CMoveRq) String() string {
	return fmt.Sprintf("CMoveRq{AffectedSOPClassUID:%v MessageID:%v Priority:%v MoveDestination:%v CommandDataSetType:%v}}", v.AffectedSOPClassUID, v.MessageID, v.Priority, v.MoveDestination, v.CommandDataSetType)
}

func decodeCMoveRq(d *messageDecoder) (*CMoveRq, error) {
	v := &CMoveRq{}
	var err error
	v.AffectedSOPClassUID, err = d.getString(dicomtag.AffectedSOPClassUID)
    if err != nil {
        return nil, err
    }
	v.MessageID, err = d.getUInt16(dicomtag.MessageID)
    if err != nil {
        return nil, err
    }
	v.Priority, err = d.getUInt16(dicomtag.Priority)
    if err != nil {
        return nil, err
    }
	v.MoveDestination, err = d.getString(dicomtag.MoveDestination)
    if err != nil {
        return nil, err
    }
	v.CommandDataSetType, err = d.getUInt16(dicomtag.CommandDataSetType)
    if err != nil {
        return nil, err
    }
	v.Extra = d.unparsedElements()
	return v, nil
}

type CMoveRsp struct {
	AffectedSOPClassUID string
	MessageIDBeingRespondedTo MessageID
	CommandDataSetType uint16
	NumberOfRemainingSuboperations uint16
	NumberOfCompletedSuboperations uint16
	NumberOfFailedSuboperations uint16
	NumberOfWarningSuboperations uint16
	Status Status
	Extra []*dicom.Element  // Unparsed elements
}

func (v *CMoveRsp) Encode(w *dicomio.Writer) error {
    elems := []*dicom.Element{}
    elem, err := dicom.NewElement(dicomtag.CommandField, uint16(32801))
    if err != nil {
        return err
    }
    elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.AffectedSOPClassUID, v.AffectedSOPClassUID)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.MessageIDBeingRespondedTo, v.MessageIDBeingRespondedTo)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.CommandDataSetType, v.CommandDataSetType)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
	if v.NumberOfRemainingSuboperations != 0 {
		elem, err = dicom.NewElement(dicomtag.NumberOfRemainingSuboperations, v.NumberOfRemainingSuboperations)
        if err != nil {
            return err
        }
		elems = append(elems, elem)
	}
	if v.NumberOfCompletedSuboperations != 0 {
		elem, err = dicom.NewElement(dicomtag.NumberOfCompletedSuboperations, v.NumberOfCompletedSuboperations)
        if err != nil {
            return err
        }
		elems = append(elems, elem)
	}
	if v.NumberOfFailedSuboperations != 0 {
		elem, err = dicom.NewElement(dicomtag.NumberOfFailedSuboperations, v.NumberOfFailedSuboperations)
        if err != nil {
            return err
        }
		elems = append(elems, elem)
	}
	if v.NumberOfWarningSuboperations != 0 {
		elem, err = dicom.NewElement(dicomtag.NumberOfWarningSuboperations, v.NumberOfWarningSuboperations)
        if err != nil {
            return err
        }
		elems = append(elems, elem)
	}
	statusElems, err := newStatusElements(v.Status)
    if err != nil {
        return err
    }
	elems = append(elems, statusElems...)
    elems = append(elems, v.Extra...)
    return encodeElements(w, elems)
}

func (v *CMoveRsp) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v *CMoveRsp) CommandField() int {
	return 32801
}

func (v *CMoveRsp) GetMessageID() MessageID {
	return v.MessageIDBeingRespondedTo
}

func (v *CMoveRsp) GetStatus() *Status {
	return &v.Status
}

func (v *CMoveRsp) String() string {
	return fmt.Sprintf("CMoveRsp{AffectedSOPClassUID:%v MessageIDBeingRespondedTo:%v CommandDataSetType:%v NumberOfRemainingSuboperations:%v NumberOfCompletedSuboperations:%v NumberOfFailedSuboperations:%v NumberOfWarningSuboperations:%v Status:%v}}", v.AffectedSOPClassUID, v.MessageIDBeingRespondedTo, v.CommandDataSetType, v.NumberOfRemainingSuboperations, v.NumberOfCompletedSuboperations, v.NumberOfFailedSuboperations, v.NumberOfWarningSuboperations, v.Status)
}

func decodeCMoveRsp(d *messageDecoder) (*CMoveRsp, error) {
	v := &CMoveRsp{}
	var err error
	v.AffectedSOPClassUID, err = d.getString(dicomtag.AffectedSOPClassUID)
    if err != nil {
        return nil, err
    }
	v.MessageIDBeingRespondedTo, err = d.getUInt16(dicomtag.MessageIDBeingRespondedTo)
    if err != nil {
        return nil, err
    }
	v.CommandDataSetType, err = d.getUInt16(dicomtag.CommandDataSetType)
    if err != nil {
        return nil, err
    }
	v.NumberOfRemainingSuboperations, err = d.getUInt16(dicomtag.NumberOfRemainingSuboperations)
    if err != nil {
        if !errors.Is(err, dicom.ErrorElementNotFound) {
            return nil, err
        }
    }
	v.NumberOfCompletedSuboperations, err = d.getUInt16(dicomtag.NumberOfCompletedSuboperations)
    if err != nil {
        if !errors.Is(err, dicom.ErrorElementNotFound) {
            return nil, err
        }
    }
	v.NumberOfFailedSuboperations, err = d.getUInt16(dicomtag.NumberOfFailedSuboperations)
    if err != nil {
        if !errors.Is(err, dicom.ErrorElementNotFound) {
            return nil, err
        }
    }
	v.NumberOfWarningSuboperations, err = d.getUInt16(dicomtag.NumberOfWarningSuboperations)
    if err != nil {
        if !errors.Is(err, dicom.ErrorElementNotFound) {
            return nil, err
        }
    }
	v.Status, err = d.getStatus()
    if err != nil {
        return nil, err
    }
	v.Extra = d.unparsedElements()
	return v, nil
}

type CEchoRq struct {
	MessageID MessageID
	CommandDataSetType uint16
	Extra []*dicom.Element  // Unparsed elements
}

func (v *CEchoRq) Encode(w *dicomio.Writer) error {
    elems := []*dicom.Element{}
    elem, err := dicom.NewElement(dicomtag.CommandField, uint16(48))
    if err != nil {
        return err
    }
    elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.MessageID, v.MessageID)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.CommandDataSetType, v.CommandDataSetType)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elems = append(elems, v.Extra...)
    return encodeElements(w, elems)
}

func (v *CEchoRq) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v *CEchoRq) CommandField() int {
	return 48
}

func (v *CEchoRq) GetMessageID() MessageID {
	return v.MessageID
}

func (v *CEchoRq) GetStatus() *Status {
	return nil
}

func (v *CEchoRq) String() string {
	return fmt.Sprintf("CEchoRq{MessageID:%v CommandDataSetType:%v}}", v.MessageID, v.CommandDataSetType)
}

func decodeCEchoRq(d *messageDecoder) (*CEchoRq, error) {
	v := &CEchoRq{}
	var err error
	v.MessageID, err = d.getUInt16(dicomtag.MessageID)
    if err != nil {
        return nil, err
    }
	v.CommandDataSetType, err = d.getUInt16(dicomtag.CommandDataSetType)
    if err != nil {
        return nil, err
    }
	v.Extra = d.unparsedElements()
	return v, nil
}

type CEchoRsp struct {
	MessageIDBeingRespondedTo MessageID
	CommandDataSetType uint16
	Status Status
	Extra []*dicom.Element  // Unparsed elements
}

func (v *CEchoRsp) Encode(w *dicomio.Writer) error {
    elems := []*dicom.Element{}
    elem, err := dicom.NewElement(dicomtag.CommandField, uint16(32816))
    if err != nil {
        return err
    }
    elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.MessageIDBeingRespondedTo, v.MessageIDBeingRespondedTo)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
    elem, err = dicom.NewElement(dicomtag.CommandDataSetType, v.CommandDataSetType)
    if err != nil {
        return err
    }
	elems = append(elems, elem)
	statusElems, err := newStatusElements(v.Status)
    if err != nil {
        return err
    }
	elems = append(elems, statusElems...)
    elems = append(elems, v.Extra...)
    return encodeElements(w, elems)
}

func (v *CEchoRsp) HasData() bool {
	return v.CommandDataSetType != CommandDataSetTypeNull
}

func (v *CEchoRsp) CommandField() int {
	return 32816
}

func (v *CEchoRsp) GetMessageID() MessageID {
	return v.MessageIDBeingRespondedTo
}

func (v *CEchoRsp) GetStatus() *Status {
	return &v.Status
}

func (v *CEchoRsp) String() string {
	return fmt.Sprintf("CEchoRsp{MessageIDBeingRespondedTo:%v CommandDataSetType:%v Status:%v}}", v.MessageIDBeingRespondedTo, v.CommandDataSetType, v.Status)
}

func decodeCEchoRsp(d *messageDecoder) (*CEchoRsp, error) {
	v := &CEchoRsp{}
	var err error
	v.MessageIDBeingRespondedTo, err = d.getUInt16(dicomtag.MessageIDBeingRespondedTo)
    if err != nil {
        return nil, err
    }
	v.CommandDataSetType, err = d.getUInt16(dicomtag.CommandDataSetType)
    if err != nil {
        return nil, err
    }
	v.Status, err = d.getStatus()
    if err != nil {
        return nil, err
    }
	v.Extra = d.unparsedElements()
	return v, nil
}

const CommandFieldCStoreRq = 1
const CommandFieldCStoreRsp = 32769
const CommandFieldCFindRq = 32
const CommandFieldCFindRsp = 32800
const CommandFieldCGetRq = 16
const CommandFieldCGetRsp = 32784
const CommandFieldCMoveRq = 33
const CommandFieldCMoveRsp = 32801
const CommandFieldCEchoRq = 48
const CommandFieldCEchoRsp = 32816
func decodeMessageForType(d* messageDecoder, commandField uint16) (Message, error) {
	switch commandField {
	case 0x1:
		return decodeCStoreRq(d)
	case 0x8001:
		return decodeCStoreRsp(d)
	case 0x20:
		return decodeCFindRq(d)
	case 0x8020:
		return decodeCFindRsp(d)
	case 0x10:
		return decodeCGetRq(d)
	case 0x8010:
		return decodeCGetRsp(d)
	case 0x21:
		return decodeCMoveRq(d)
	case 0x8021:
		return decodeCMoveRsp(d)
	case 0x30:
		return decodeCEchoRq(d)
	case 0x8030:
		return decodeCEchoRsp(d)
	default:
	    return nil, fmt.Errorf("Unknown DIMSE command 0x%x", commandField)
	}
}
