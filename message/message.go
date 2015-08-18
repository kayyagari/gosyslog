package message

import (
	"time"
)

/*
 * HEADER = PRI VERSION SP TIMESTAMP SP HOSTNAME SP APP-NAME SP PROCID SP MSGID
 */
type Header struct {
	Pri       int
	Version   int
	Timestamp time.Time
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
