package server

import (
	"bytes"
	"errors"
	sysmsg "gosyslog/message"
	"strconv"
	"strings"
	"time"
)

func Parse(buf *bytes.Buffer) (sysmsg.Message, error) {
	msg := sysmsg.Message{}

	header, err := parseHeader(buf)

	return msg, nil
}

func parseHeader(buf *bytes.Buffer) (*sysmsg.Header, error) {
	prefix := getToken(buf)
	pos := strings.Index(prefix, ">")

	var pri int
	if pos > 1 {
		p := prefix[1:pos]
		pri, err := strconv.Atoi(p)
		if err != nil {
			return nil, errors.New("Failed to parse the priority " + p)
		}
	} else {
		return nil, errors.New("Invalid header, does not contain <pri>")
	}

	l := len(prefix)
	if (l - 1) < (pos + 1) {
		return nil, errors.New("Invalid header, does not contain version number")
	}

	v := prefix[pos+1 : len(prefix)]
	ver, err := strconv.Atoi(v)
	if err != nil {
		return nil, errors.New("Failed to parse the version " + v)
	}

	timestamp, err := parseTime(buf)
	if err != nil {
		return nil, err
	}

	hostName := getToken(buf)
	appName := getToken(buf)
	procId := getToken(buf)
	msgId := getToken(buf)

	return nil, nil
}

func parseTime(buf *bytes.Buffer) (time.Time, error) {
	var t time.Time
	timestamp := getToken(buf)
	if len(timestamp) == 0 {
		return t, nil
	}

	dotPos := strings.Index(timestamp, ".")
	if dotPos > 0 {
		zPos := strings.Index(timestamp, "Z")
		if zPos == -1 {
			zPos := strings.LastIndex(timestamp, "-")
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

func getToken(buf *bytes.Buffer) (string) {
	token, err := buf.ReadString(' ')
	token = strings.TrimSpace(token)

	// check if the value is nil
	if strings.EqualFold(token, "-") {
		return ""
	}

	return token
}
