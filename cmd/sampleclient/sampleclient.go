// A sample program for issuing C-STORE or C-FIND to a remote server.
package main

import (
	"flag"
	"log"

	"github.com/kristianvalind/go-netdicom"
	"github.com/kristianvalind/go-netdicom/pkg/dimse"
	"github.com/kristianvalind/go-netdicom/pkg/sopclass"
	"github.com/suyashkumar/dicom"
	dicomtag "github.com/suyashkumar/dicom/pkg/tag"
)

var (
	serverFlag        = flag.String("server", "localhost:10000", "host:port of the remote application entity")
	storeFlag         = flag.String("store", "", "If set, issue C-STORE to copy this file to the remote server")
	aeTitleFlag       = flag.String("ae-title", "testclient", "AE title of the client")
	remoteAETitleFlag = flag.String("remote-ae-title", "testserver", "AE title of the server")
	findFlag          = flag.Bool("find", false, "Issue a C-FIND.")
	getFlag           = flag.Bool("get", false, "Issue a C-GET.")
	seriesFlag        = flag.String("series", "", "Study series UID to retrieve in C-{FIND,GET}.")
	studyFlag         = flag.String("study", "", "Study instance UID to retrieve in C-{FIND,GET}.")
)

func mustNewElement(t dicomtag.Tag, v interface{}) *dicom.Element {
	elem, err := dicom.NewElement(t, v)
	if err != nil {
		log.Panic(err)
	}

	return elem
}

func newServiceUser(sopClasses []string) *netdicom.ServiceUser {
	su, err := netdicom.NewServiceUser(netdicom.ServiceUserParams{
		CalledAETitle:  *remoteAETitleFlag,
		CallingAETitle: *aeTitleFlag,
		SOPClasses:     sopClasses})
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Connecting to %s", *serverFlag)
	su.Connect(*serverFlag)
	return su
}

func cStore(inPath string) {
	su := newServiceUser(sopclass.StorageClasses)
	defer su.Release()
	dataset, err := dicom.ParseFile(inPath, nil)
	if err != nil {
		log.Panicf("%s: %v", inPath, err)
	}
	err = su.CStore(&dataset)
	if err != nil {
		log.Panicf("%s: cstore failed: %v", inPath, err)
	}
	log.Printf("C-STORE finished successfully")
}

func generateCFindElements() (netdicom.QRLevel, []*dicom.Element) {
	if *seriesFlag != "" {
		return netdicom.QRLevelSeries, []*dicom.Element{mustNewElement(dicomtag.SeriesInstanceUID, *seriesFlag)}
	}
	if *studyFlag != "" {
		return netdicom.QRLevelStudy, []*dicom.Element{mustNewElement(dicomtag.StudyInstanceUID, *studyFlag)}
	}
	args := []*dicom.Element{
		mustNewElement(dicomtag.SpecificCharacterSet, []string{"ISO_IR 100"}),
		mustNewElement(dicomtag.AccessionNumber, []string{}),
		mustNewElement(dicomtag.ReferringPhysicianName, []string{}),
		mustNewElement(dicomtag.PatientName, []string{}),
		mustNewElement(dicomtag.PatientID, []string{}),
		mustNewElement(dicomtag.PatientBirthDate, []string{}),
		mustNewElement(dicomtag.PatientSex, []string{}),
		mustNewElement(dicomtag.StudyInstanceUID, []string{}),
		mustNewElement(dicomtag.RequestedProcedureDescription, []string{}),
		mustNewElement(dicomtag.ScheduledProcedureStepSequence, [][]*dicom.Element{
			{
				mustNewElement(dicomtag.Modality, []string{}),
				mustNewElement(dicomtag.ScheduledProcedureStepStartDate, []string{}),
				mustNewElement(dicomtag.ScheduledProcedureStepStartTime, []string{}),
				mustNewElement(dicomtag.ScheduledPerformingPhysicianName, []string{}),
				mustNewElement(dicomtag.ScheduledProcedureStepStatus, []string{}),
			}},
		),
	}
	return netdicom.QRLevelPatient, args
}

func cGet() {
	su := newServiceUser(sopclass.QRGetClasses)
	defer su.Release()
	qrLevel, args := generateCFindElements()
	n := 0
	err := su.CGet(qrLevel, args,
		func(transferSyntaxUID, sopClassUID, sopInstanceUID string, data []byte) dimse.Status {
			log.Printf("%d: C-GET data; transfersyntax=%v, sopclass=%v, sopinstance=%v data %dB",
				n, transferSyntaxUID, sopClassUID, sopInstanceUID, len(data))
			n++
			return dimse.Success
		})
	log.Printf("C-GET finished: %v", err)
}

func cFind() {
	su := newServiceUser(sopclass.QRFindClasses)
	defer su.Release()
	qrLevel, args := generateCFindElements()
	for result := range su.CFind(qrLevel, args) {
		if result.Err != nil {
			log.Printf("C-FIND error: %v", result.Err)
			continue
		}
		log.Printf("Got response with %d elems", len(result.Elements))
		for _, elem := range result.Elements {
			log.Printf("Got elem: %v", elem.String())
		}
	}
}

func main() {
	flag.Parse()
	if *storeFlag != "" {
		cStore(*storeFlag)
	} else if *findFlag {
		cFind()
	} else if *getFlag {
		cGet()
	} else {
		log.Panic("Either -store, -get, or -find must be set")
	}
}
