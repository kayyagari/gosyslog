package server

import (
	"bytes"
	"errors"
	"fmt"
	sysmsg "gosyslog/message"
	"io"
	"strconv"
	"strings"
	"time"
)

var VerParseErr = errors.New("Failed to parse the version")
var PriParseErr = errors.New("Failed to parse the priority")
var BadPriErr = errors.New("Invalid header, priority value cannot contain more than 3 digits")
var NoPriErr = errors.New("Invalid header, does not contain <pri>")
var NoVerErr = errors.New("Invalid header, does not contain version number")
var TimeMillisErr = errors.New("Invalid timestamp format, misplaced milliseconds")
var TimeNanoErr = errors.New("Invalid timestamp format, TIME-SECFRAC is in nanoseconds")
var BadVerErr = errors.New("Invalid header, version value cannot contain more than 2 digits")
var BadSDErr = errors.New("Invalid SData")

func ParseString(logMsg string) (*sysmsg.Message, error) {
	buf := bytes.NewBufferString(logMsg)
	return Parse(buf)
}

func Parse(buf *bytes.Buffer) (*sysmsg.Message, error) {
	header, err := ParseHeader(buf)
	if err != nil {
		return nil, err
	}

	sd, err := ParseSData(buf)
	if err != nil {
		return nil, err
	}

	var rawMsg []byte
	msg := &sysmsg.Message{header, sd, rawMsg}
	return msg, nil
}

func ParseSData(buf *bytes.Buffer) (*sysmsg.StrctData, error) {
	var sData *sysmsg.StrctData
	startChar, _, err := buf.ReadRune()

	if err != nil {
		return sData, err
	}
	// [ = 91
	// ] = 93
	// \ = 92

	// if it doesn't start with '[' return error
	if startChar != 91 {
		return sData, BadSDErr
	}

	buf.UnreadRune()

	// now, there exists SData, parse it
	readSData := true
	endSData := false
	//[exampleSDID iut="3" eventSource="Application" eventID="1011"][examplePriority@32473 class="high"]
	var element sysmsg.SDElement
	var readParam bool

loop:
	for {
		startChar, _, err := buf.ReadRune()

		if err == io.EOF && endSData {
			return sData, nil
		}else if err != nil {
			return sData, err
		}

		switch {
		case startChar == '[' || readSData:
			//fmt.Println("begin SData")
			readSData = false
			endSData = false
			id, err := parseSDName(buf, ' ')
			//fmt.Println("Parsed ID ", id)
			if err != nil {
				return sData, err
			}
			element = sysmsg.SDElement{}
			element.Id = id
			element.Params = make([]sysmsg.SDParam, 0, 2)
			readParam = true
			break

		case startChar == ' ' || readParam:
			//fmt.Println("begin SData SP")

			// step back one char
			if readParam {
				buf.UnreadRune()
			}

			readParam = false
			if endSData {
				break loop
			}

			name, err := parseSDName(buf, '=')
			if err != nil {
				return sData, err
			}

			val, err := parseSDParamVal(buf)
			if err != nil {
				return sData, err
			}

			param := sysmsg.SDParam{name, val}
			element.Params = append(element.Params, param)

			break

		case startChar == ']':
			//fmt.Println("end SData")
			endSData = true
			if sData == nil {
				sData = &sysmsg.StrctData{}
				sData.Elements = make([]sysmsg.SDElement, 0, 2)
			}

			sData.Elements = append(sData.Elements, element)
			continue loop
		}
	}

	fmt.Println("completed parsing SData", sData)

	return sData, nil
}

func parseSDParamVal(buf *bytes.Buffer) (string, error) {
	slice := make([]rune, 0, 128)
	firstQuote := true

	for {
		startChar, _, err := buf.ReadRune()
		if err != nil {
			return "", err
		}

		switch {
		case startChar == '\\':
			// check the next char it must be an escape char
			nextChar, _, err := buf.ReadRune()
			if err != nil {
				return "", err
			}

			if nextChar != '\\' && nextChar != '"' && nextChar != ']' {
				slice = append(slice, startChar)
			}

			slice = append(slice, nextChar)
			break

		case startChar == '"':
			if !firstQuote {
				return string(slice), nil
			}

			if firstQuote {
				firstQuote = false
			}
			break

		default:
			slice = append(slice, startChar)
			break
		}
	}

}

//SD-NAME = 1*32PRINTUSASCII ; except '=', SP, ']', %d34 (")
func parseSDName(buf *bytes.Buffer, delim rune) (string, error) {
	var count int
	slice := make([]rune, 0, 32)
	for {
		char, _, err := buf.ReadRune()
		count = count + 1
		if err != nil {
			return "", err
		}

		if count > 32 {
			return "", BadSDErr
		}

		switch char {
		case ']', '"':
			return "", BadSDErr

		case delim:
			return string(slice), nil

		default:
			slice = append(slice, char)
		}
	}
}

func ParseHeader(buf *bytes.Buffer) (sysmsg.Header, error) {
	prefix, err := getToken(buf)

	var h sysmsg.Header

	pos := strings.Index(prefix, "<")
	if pos != 0 {
		return h, PriParseErr
	}

	pos = strings.Index(prefix, ">")

	var pri int
	if pos > 1 {
		p := prefix[1:pos]
		if len(p) > 3 {
			return h, BadPriErr
		}

		pri, err = strconv.Atoi(p)
		if err != nil {
			return h, PriParseErr
		}
	} else {
		return h, NoPriErr
	}

	l := len(prefix)
	if (l - 1) < (pos + 1) {
		return h, NoVerErr
	}

	v := prefix[pos+1 : len(prefix)]
	if len(v) > 2 {
		return h, BadVerErr
	}
	ver, err := strconv.Atoi(v)
	if err != nil {
		return h, VerParseErr
	}

	timestamp, err := parseTime(buf)
	if err != nil {
		return h, err
	}

	hostName, err := getToken(buf)
	hostLen := len(hostName)
	if hostLen > 255 {
		hostName = hostName[0:255]
	}
	if err != nil {
		return h, err
	}

	appName, err := getToken(buf)
	appLen := len(appName)
	if appLen > 48 {
		appName = appName[0:48]
	}
	if err != nil {
		return h, err
	}

	procId, err := getToken(buf)
	procLen := len(procId)
	if procLen > 128 {
		procId = procId[0:128]
	}
	if err != nil {
		return h, err
	}

	msgId, err := getToken(buf)
	msgLen := len(msgId)
	if msgLen > 32 {
		msgId = msgId[0:32]
	}
	// message can be empty after the header

	return sysmsg.Header{pri, ver, timestamp, hostName, appName, procId, msgId}, nil
}

func parseTime(buf *bytes.Buffer) (time.Time, error) {
	var t time.Time
	timestamp, _ := getToken(buf)
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
			return t, TimeMillisErr
		}

		millis := timestamp[dotPos+1 : zPos]

		if len(millis) > 6 {
			return t, TimeNanoErr
		}
	}

	t, err := time.Parse(time.RFC3339, timestamp)
	return t, err
}

func getToken(buf *bytes.Buffer) (string, error) {
	token, err := buf.ReadString(' ')

	//if err == io.EOF {
	//	panic("EOF while trying to read token")
	//}

	token = strings.TrimSpace(token)

	// check if the value is nil
	if strings.EqualFold(token, "-") {
		return "", err
	}

	//fmt.Println("parsed token ", token)
	return token, err
}
