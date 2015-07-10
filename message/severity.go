package message

type Severity struct {
	code int
	Keyword string
	Descr string
}

var Emergency = Severity{code: 0, Keyword: "emerg", Descr: "system is unusable"}
var Alert = Severity{code: 1, Keyword: "alert", Descr: "action must be taken immediately"}
var Critical = Severity{code: 2, Keyword: "crit", Descr: "critical conditions"}
var Error = Severity{code: 3, Keyword: "err", Descr: "error conditions"}
var Warning = Severity{code: 4, Keyword: "warning", Descr: "warning conditions"}
var Notice = Severity{code: 5, Keyword: "notice", Descr: "normal but significant condition"}
var Informational = Severity{code: 6, Keyword: "info", Descr: "informational messages"}
var Debugging = Severity{code: 7, Keyword: "debug", Descr: "debug-level messages"}
