package message

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

const NIL_VAL_STR = "-"

const SP_VAL_BYTE = ' '

const SYSLOG_TIME_FORMAT = "2006-01-02T15:04:05.999999Z07:00"

/*
 * HEADER = PRI VERSION SP TIMESTAMP SP HOSTNAME SP APP-NAME SP PROCID SP MSGID
 */
type Header struct {
	Pri       int
	Version   int
	Timestamp *time.Time
	Hostname  string
	Appname   string
	Procid    string
	Msgid     string
}

/*
 * SD-PARAM = PARAM-NAME "=" %d34 PARAM-VALUE %d34
 */
type SDParam struct {
	Name  string
	Value string
}

/*
 * SD-ELEMENT = "[" SD-ID *(SP SD-PARAM) "]"
 */
type SDElement struct {
	Id     string
	Params []SDParam
}

/*
 * STRUCTURED-DATA = NILVALUE / 1*SD-ELEMENT
 */
type StrctData struct {
	Elements []SDElement
}

/*
 * SYSLOG-MSG = HEADER SP STRUCTURED-DATA [SP MSG]
 */
type Message struct {
	Header Header
	SData  *StrctData
	RawMsg []byte
	IsUtf8 bool
}

func (sd *StrctData) String() string {
	str := NIL_VAL_STR
	if sd.Elements != nil {
		str = ""
		for _, e := range sd.Elements {
			str = str + e.String()
		}
	}

	return str
}

func (sde *SDElement) String() string {
	prefix := fmt.Sprintf("[%s", sde.Id)
	if sde.Params != nil {
		for _, param := range sde.Params {
			prefix = prefix + " " + param.String()
		}
	}

	return prefix + "]"
}

func (sdparam *SDParam) String() string {
	return fmt.Sprintf("%s=\"%s\"", sdparam.Name, sdparam.Value)
}

func (h *Header) String() string {
	t := NIL_VAL_STR
	if h.Timestamp != nil {
		t = h.Timestamp.Format(SYSLOG_TIME_FORMAT)
		fmt.Println("time ", t)
	}

	host := NIL_VAL_STR
	if len(h.Hostname) != 0 {
		host = h.Hostname
	}

	app := NIL_VAL_STR
	if len(h.Appname) != 0 {
		app = h.Appname
	}

	pid := NIL_VAL_STR
	if len(h.Procid) != 0 {
		pid = h.Procid
	}

	mid := NIL_VAL_STR
	if len(h.Msgid) != 0 {
		mid = h.Msgid
	}

	return fmt.Sprintf("<%d>%d %s %s %s %s %s", h.Pri, h.Version, t, host, app, pid, mid)

}

func (m *Message) String() string {
	header := m.Header.String()

	sd := " NIL-SD "
	if m.SData != nil {
		sd = fmt.Sprintf(" %+v ", m.SData)
	}

	utf := "(binary) "
	data := ""
	if m.IsUtf8 {
		utf = "(UTF-8) "
	}
	if m.RawMsg != nil {
		if m.IsUtf8 {
			data = string(m.RawMsg)
		} else if !m.IsUtf8 {
			data = "(" + strconv.Itoa(len(m.RawMsg)) + " bytes)"
		}
	} else if m.RawMsg == nil {
		data = "NIL-DATA"
	}

	return header + sd + utf + data
}

func (m *Message) Bytes() []byte {
	buf := bytes.NewBufferString(m.Header.String())
	sd := NIL_VAL_STR

	if m.SData != nil {
		sd = m.SData.String()
	}

	buf.WriteByte(SP_VAL_BYTE)
	buf.WriteString(sd)

	if m.RawMsg != nil {
		buf.WriteByte(SP_VAL_BYTE)
		if m.IsUtf8 {
			BOM := []byte{0xEF, 0xBB, 0xBF}
			buf.Write(BOM)
		}
		buf.Write(m.RawMsg)
	}

	return buf.Bytes()
}
