package dimse_test

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/kristianvalind/go-netdicom/pkg/dimse"
	"github.com/suyashkumar/dicom/pkg/dicomio"
)

func testDIMSE(t *testing.T, v dimse.Message) {
	b := &bytes.Buffer{}
	w := dicomio.NewWriter(b, binary.LittleEndian, true)

	w.SetTransferSyntax(binary.LittleEndian, true)

	err := dimse.EncodeMessage(&w, v)
	if err != nil {
		t.Fatal(err)
	}

	b2 := bufio.NewReader(bytes.NewBuffer(b.Bytes()))

	d, err := dicomio.NewReader(b2, binary.LittleEndian, int64(len(b.Bytes())))
	if err != nil {
		t.Fatal(err)
	}
	v2, err := dimse.ReadMessage(d)
	if err != nil {
		t.Fatal(err)
	}
	if v.String() != v2.String() {
		t.Errorf("%v <-> %v", v, v2)
	}
}

func TestCStoreRq(t *testing.T) {
	testDIMSE(t, &dimse.CStoreRq{
		"1.2.3",
		0x1234,
		0x2345,
		1,
		"3.4.5",
		"foohah",
		0x3456, nil})
}

func TestCStoreRsp(t *testing.T) {
	testDIMSE(t, &dimse.CStoreRsp{
		"1.2.3",
		0x1234,
		dimse.CommandDataSetTypeNull,
		"3.4.5",
		dimse.Status{Status: dimse.StatusCode(0x3456)},
		nil})
}

func TestCEchoRq(t *testing.T) {
	testDIMSE(t, &dimse.CEchoRq{0x1234, 1, nil})
}

func TestCEchoRsp(t *testing.T) {
	testDIMSE(t, &dimse.CEchoRsp{0x1234, 1,
		dimse.Status{Status: dimse.StatusCode(0x2345)},
		nil})
}
