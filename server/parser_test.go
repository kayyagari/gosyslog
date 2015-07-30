package server

import (
	"bytes"
	"fmt"
	"testing"
	"time"
	//sysmsg "gosyslog/message"
)

func TestParseValidHeader(t *testing.T) {
	buf := bytes.NewBufferString("<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su - ID47")
	msg, err := Parse(buf)

	header := msg.Header
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

	fmt.Println("Appname ", header.Appname)
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
