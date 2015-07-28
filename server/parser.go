package server

import (
	"bytes"
	"errors"
	sysmsg "gosyslog/message"
	"io"
	"strconv"
	"strings"
	"time"
)

func Parse(buf *bytes.Buffer) (sysmsg.Message, error) {
	header, err := parseHeader(buf)
	if err != nil {
		var s sysmsg.Message
		return s, err
	}
	var sd sysmsg.StrctData
	var rawMsg []byte
	msg := sysmsg.Message{header, sd, rawMsg}
	return msg, nil
}

func parseHeader(buf *bytes.Buffer) (sysmsg.Header, error) {
	prefix := getToken(buf)
	pos := strings.Index(prefix, ">")

	var h sysmsg.Header

	var pri int
	if pos > 1 {
		p := prefix[1:pos]
		// FIXME workaround for the error :- server/parser.go:35: pri declared and not used
		err := errors.New("")
		pri, err = strconv.Atoi(p)
		if err != nil {
			return h, errors.New("Failed to parse the priority " + p)
		}
	} else {
		return h, errors.New("Invalid header, does not contain <pri>")
	}

	l := len(prefix)
	if (l - 1) < (pos + 1) {
		return h, errors.New("Invalid header, does not contain version number")
	}

	v := prefix[pos+1 : len(prefix)]
	ver, err := strconv.Atoi(v)
	if err != nil {
		return h, errors.New("Failed to parse the version " + v)
	}

	timestamp, err := parseTime(buf)
	if err != nil {
		return h, err
	}

	hostName := getToken(buf)
	appName := getToken(buf)
	procId := getToken(buf)
	msgId := getToken(buf)

	return sysmsg.Header{pri, ver, timestamp, hostName, procId, appName, msgId}, nil
}

func parseTime(buf *bytes.Buffer) (time.Time, error) {
	var t time.Time
	timestamp := getToken(buf)
	if len(timestamp) == 0 {
		return t, nil
	}

	dotPos := strings.Index(timestamp, ".")
	if dotPos > 0 {
		var zPos int = strings.Index(timestamp, "Z")
		if zPos == -1 {
			zPos = strings.LastIndex(timestamp, "-")
		}

		if dotPos > zPos {
			return t, errors.New("Invalid timestamp format, misplaced milliseconds " + timestamp)
		}

		millis := timestamp[dotPos+1 : zPos]

		if len(millis) > 6 {
			return t, errors.New("Invalid timestamp format, TIME-SECFRAC is in nanoseconds " + timestamp)
		}
	}

	t, err := time.Parse(time.RFC3339, timestamp)
	return t, err
}

func getToken(buf *bytes.Buffer) string {
	token, err := buf.ReadString(' ')
	if err == io.EOF {
		panic("EOF while trying to read token")
	}
	token = strings.TrimSpace(token)

	// check if the value is nil
	if strings.EqualFold(token, "-") {
		return ""
	}

	return token
}
