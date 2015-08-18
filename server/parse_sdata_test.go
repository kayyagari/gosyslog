package server

import (
	"bytes"
	_ "fmt"
	sysmsg "gosyslog/message"
	"io"
	"testing"
)

func TestParseValidSData(t *testing.T) {
	buf := bytes.NewBufferString("[exampleSDID iut=\"7\" eventSource=\"Application\" eventID=\"101\"][examplePriority class=\"high\"]")
	sdata, err := ParseSData(buf)
	verifyBasicState(t, sdata, err, 2)

	sdElem := sdata.Elements[0]
	if sdElem.Id != "exampleSDID" {
		t.Errorf("SDID of first element is not matching - %v", sdElem.Id)
	}

	sdParams := sdElem.Params
	if len(sdParams) != 3 {
		t.Errorf("Incorrect number of SDParams")
	}

	verifySDParam(t, &sdParams[0], "iut", "7")
	verifySDParam(t, &sdParams[1], "eventSource", "Application")
	verifySDParam(t, &sdParams[2], "eventID", "101")

	// the second SDElement
	sdElem = sdata.Elements[1]
	if sdElem.Id != "examplePriority" {
		t.Errorf("SDID of second element is not matching - %v", sdElem.Id)
	}

	sdParams = sdElem.Params
	if len(sdParams) != 1 {
		t.Errorf("Incorrect number of SDParams")
	}

	verifySDParam(t, &sdParams[0], "class", "high")
}

func TestParseInvalidSData(t *testing.T) {
	buf := bytes.NewBufferString("[quoteInValue k=\"v\\\"\"]")
	sdata, err := ParseSData(buf)
	verifyBasicState(t, sdata, err, 1)
	sdElem := sdata.Elements[0]
	sdParams := sdElem.Params
	verifySDParam(t, &sdParams[0], "k", "v\"")

	buf = bytes.NewBufferString("[valueWithBadEscape k=\"v\\|\"]")
	sdata, err = ParseSData(buf)
	verifyBasicState(t, sdata, err, 1)
	sdElem = sdata.Elements[0]
	sdParams = sdElem.Params
	verifySDParam(t, &sdParams[0], "k", "v\\|")

	buf = bytes.NewBufferString("[spaceAtEndErr k=\"v\" ]")
	checkBadSData(t, buf)
	buf = bytes.NewBufferString("[wrongCharIn]Id k=\"v\"]")
	checkBadSData(t, buf)

	// empty SDdata
	buf = bytes.NewBufferString("[]")
	checkBadSData(t, buf)

	// incomplete SData
	buf = bytes.NewBufferString("[")
	_, err = ParseSData(buf)
	if err != io.EOF {
		t.Errorf("Must fail with EOF")
	}

}

func checkBadSData(t *testing.T, buf *bytes.Buffer) {
	_, err := ParseSData(buf)
	if err != BadSDErr {
		t.Errorf("Must fail with BadSDErr")
	}
}

func verifySDParam(t *testing.T, sdParam *sysmsg.SDParam, name, value string) {
	if sdParam.Name != name {
		t.Errorf("SDParam Name didn't match - %v", name)
	}

	if sdParam.Value != value {
		t.Errorf("SDParam Value didn't match - %v", value)
	}
}

func verifyBasicState(t *testing.T, sdata *sysmsg.StrctData, err error, count int) {
	if err != nil {
		t.Errorf(err.Error())
	}

	if sdata == nil {
		t.Errorf("Failed to parse SData")
		t.FailNow()
	}

	if len(sdata.Elements) != count {
		t.Errorf("Incorrect number of SDElements")
		t.FailNow()
	}
}
