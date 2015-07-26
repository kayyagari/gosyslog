package server

import (
	//"strconv"
	"bytes"
	sysmsg "gosyslog/message"
	"strings"
)

func Parse(buf Buffer) (sysmsg.Message, error) {
	msg := sysmsg.Message{}

	header, err = praseHeader(buf)

	return msg, nil
}

func parseHeader(buf Buffer) (sysmsg.Header, error) {
	prefix, _ := buf.ReadString(' ')
	pos := strings.Index(prefix, ">")

	var pri string
	if pos > 0 {
		pri = prefix[1:pos]
	}
	if pos <= 0 {
		panic("Invalid header, does not contain <pri>")
	}
}

func getToken(buf Buffer) (string, error) {
	for {
		buf.ReadBytes(' ')
	}
}
