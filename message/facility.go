package message

type Facility struct {
	code    int
	Keyword string
	Descr   string
}

var Kern = Facility{code: 0, Keyword: "kern", Descr: "kernel messages"}
var User = Facility{code: 1, Keyword: "user", Descr: "user-level messages"}
var Mail = Facility{code: 2, Keyword: "mail", Descr: "mail system"}
var Daemon = Facility{code: 3, Keyword: "daemon", Descr: "system daemons"}
var Auth = Facility{code: 4, Keyword: "auth", Descr: "security/authorization messages"}
var Syslog = Facility{code: 5, Keyword: "syslog", Descr: "messages generated internally by syslogd"}
var Lpr = Facility{code: 6, Keyword: "lpr", Descr: "line printer subsystem"}
var News = Facility{code: 7, Keyword: "news", Descr: "network news subsystem"}
var Uucp = Facility{code: 8, Keyword: "uucp", Descr: "UUCP subsystem"}
var Clock = Facility{code: 9, Keyword: "clock", Descr: "clock daemon"}
var Authpriv = Facility{code: 10, Keyword: "authpriv", Descr: "security/authorization messages"}
var Ftp = Facility{code: 11, Keyword: "ftp", Descr: "FTP daemon"}
var Ntp = Facility{code: 12, Keyword: "ntp", Descr: "NTP subsystem"}
var Audit = Facility{code: 13, Keyword: "audit", Descr: "log audit"}
var AlertFacility = Facility{code: 14, Keyword: "alert", Descr: "log alert"}
var Cron = Facility{code: 15, Keyword: "cron", Descr: "scheduling daemon"}
var Local0 = Facility{code: 16, Keyword: "local0", Descr: "local use 0 (local0)"}
var Local1 = Facility{code: 17, Keyword: "local1", Descr: "local use 1 (local1)"}
var Local2 = Facility{code: 18, Keyword: "local2", Descr: "local use 2 (local2)"}
var Local3 = Facility{code: 19, Keyword: "local3", Descr: "local use 3 (local3)"}
var Local4 = Facility{code: 20, Keyword: "local4", Descr: "local use 4 (local4)"}
var Local5 = Facility{code: 21, Keyword: "local5", Descr: "local use 5 (local5)"}
var Local6 = Facility{code: 22, Keyword: "local6", Descr: "local use 6 (local6)"}
var Local7 = Facility{code: 23, Keyword: "local7", Descr: "local use 7 (local7)"}
