package logs

import "time"

type Element interface {
	MsgF(string, ...interface{}) Element
	Msg(...interface{}) Element
	Trace() Element
	Now() Element
	WithField(interface{}, interface{}) Element
}

type Logs interface {
	Error(Element)
	Warn(Element)
	Info(Element)
	Debug(Element)
}