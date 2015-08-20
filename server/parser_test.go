package server

import (
	"bytes"
	"fmt"
	sysmsg "gosyslog/message"
	"io"
	"testing"
	"time"
)

func TestParseValidHeader(t *testing.T) {
	buf := bytes.NewBufferString("<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su - ID47")
	header, err := ParseHeader(buf)

	if err != nil {
		t.Errorf(err.Error())
	}

	if header.Pri != 34 {
		t.Errorf("Priority didn't match")
	}

	if header.Version != 1 {
		t.Errorf("Version didn't match")
	}

	expected, _ := time.Parse(time.RFC3339, "2003-10-11T22:14:15.003Z")
	if header.Timestamp.Unix() != expected.Unix() {
		t.Errorf("Time didn't match")
	}

	if header.Hostname != "mymachine.example.com" {
		t.Errorf("Hostname didn't match")
	}

	if header.Appname != "su" {
		t.Errorf("Appname didn't match")
	}

	if len(header.Procid) != 0 {
		t.Errorf("Procid didn't match")
	}

	if header.Msgid != "ID47" {
		t.Errorf("Msgid didn't match")
	}
}

func TestParseBadHeader(t *testing.T) {
	_, err := parseHeaderString("<340")
	if err != NoPriErr {
		t.Errorf("Must fail with NoPriErr")
	}

	_, err = parseHeaderString("340>")
	if err != PriParseErr {
		t.Errorf("Must fail with PriParseErr")
	}

	_, err = parseHeaderString("<3401>")
	if err != BadPriErr {
		t.Errorf("Must fail with BadPriErr")
	}

	_, err = parseHeaderString("<1>1234")
	if err != BadVerErr {
		t.Errorf("Must fail with BadVerErr")
	}

	_, err = parseHeaderString("<1>1a")
	if err != VerParseErr {
		t.Errorf("Must fail with VerParseErr")
	}

	//PRI VERSION SP TIMESTAMP SP HOSTNAME SP APP-NAME SP PROCID SP MSGID
	_, err = parseHeaderString("<1>1 - - - -")
	if err != io.EOF {
		t.Errorf("Must fail with EOF", err)
	}
}

func TestParseMessage(t *testing.T) {
	headerAndSd := "<34>1 2003-10-11T22:14:15.003Z localhost su 1 ID47 - "
	buf := bytes.NewBufferString(headerAndSd)
	bom := []byte{0xEF, 0xBB, 0xBF}
	buf.Write(bom)

	data := []byte("abc")
	buf.Write(data)

	fmt.Println()

	msg, err := Parse(buf)

	fmt.Println("** ", msg)

	if err != nil {
		t.Errorf("Message parsing should not fail", err)
	}

	if msg == nil {
		t.Errorf("Message parsing failed")
	}

	if !msg.IsUtf8 {
		t.Errorf("Message data must be parsed as UTF8")
	}

	if bytes.Compare(msg.RawMsg, data) != 0 {
		t.Errorf("Raw message data is not identical")
	}

	// non-UTF8
	buf = bytes.NewBufferString(headerAndSd)
	buf.Write(data)

	msg, err = Parse(buf)
	fmt.Println("** ", msg)
	if msg.IsUtf8 {
		t.Errorf("Message data must NOT be parsed as UTF8")
	}

	if bytes.Compare(msg.RawMsg, data) != 0 {
		t.Errorf("Raw message data is not identical")
	}
}

func parseHeaderString(logMsg string) (sysmsg.Header, error) {
	buf := bytes.NewBufferString(logMsg)
	return ParseHeader(buf)
}
