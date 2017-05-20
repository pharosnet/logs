package logs

import (
	"time"
	"bytes"
	"fmt"
	"runtime"
	"encoding/json"
)


//
type caller struct {
	file string
	line int
	fn string
	ok bool
}

func (c caller) String() string {
	if c.ok {
		if c.file == "" {

			return fmt.Sprintf("[\x1b[%dm%s:%d\x1b[0m]", 37, c.fn, c.line)
		}
		return fmt.Sprintf("[\x1b[%dm%s:%d\x1b[0m]", 37, c.file, c.line)
	}
	return ""
}

//
type element struct {
	dateTime time.Time
	level Level
	msg string
	c caller
	extra map[interface{}]interface{}
}

func newElement(level Level, format string, v ...interface{}) *element {
	e := new(element)
	e.level = level
	if format == "" {
		e.msg = fmt.Sprint(v...)
	} else {
		e.msg = fmt.Sprintf(format, v...)
	}
	e.dateTime = time.Now()
	e.extra = make(map[interface{}]interface{})
	return e
}

// TODO
func ParseElement(p []byte) (Element, error) {
	return nil, nil
}

func (e *element) Level() Level {
	return e.level
}

func (e *element) Trace() *element {
	pc, _, line, ok := runtime.Caller(1)
	c := caller{}
	if  ok {
		c.line = line
		c.ok = ok
		c.fn = runtime.FuncForPC(pc).Name()
	}
	e.c = c
	return e
}

func (e *element) TraceFile() *element {
	_, f, line, ok := runtime.Caller(1)
	c := caller{}
	if  ok {
		c.file = f
		c.line = line
		c.ok = ok
	}
	e.c = c
	return e
}

func (e *element) WithField(k, v interface{}) *element {
	e.extra[k] = v
	return e
}

// {LEVEL}{TIME}{MSG}{EXTRA}{CALLER}
func (e *element) Bytes() []byte {
	buf := bytes.NewBufferString("")
	// level
	if !e.level.IsNoLevel() {
		buf.WriteString(fmt.Sprintf("[\x1b[%dm%s\x1b[0m] ", e.level.Color(), e.level.String()))
		buf.WriteString(" ")
	}
	// datetime
	buf.WriteString(fmt.Sprintf("[\x1b[%dm%s\x1b[0m]", 37, e.dateTime.String())) // 36
	// msg
	buf.WriteString(" ")
	buf.WriteString(e.msg)
	// extra
	if len(e.extra) > 0 {
		for k, v := range e.extra {
			buf.WriteString(" ")
			buf.WriteString(fmt.Sprintf("\x1b[%dm%s\x1b[0m=", e.level.Color(), k))
			buf.WriteString(fmt.Sprintf("\x1b[%dm%s\x1b[0m=", 37,"="))
			buf.WriteString(fmt.Sprint(v))
		}
	}
	// caller
	if e.c.ok {
		buf.WriteString(" ")
		buf.WriteString(e.c.String())
	}
	return buf.Bytes()
}

func (e element) String() string {
	return string(e.Bytes())
}

func (e element) bytesWithNewLine() string {
	return string(e.Bytes()) + "\n"
}

func (e element) JSON() string {
	ej := new(E)
	ej.DateTime = e.dateTime
	if !e.level.IsNoLevel() {
		ej.Level = e.level.String()
	}
	ej.Msg = e.msg
	if e.c.ok {
		ej.File = e.c.file
		ej.Line = e.c.line
		ej.Func = e.c.fn
	}
	if len(e.extra) > 0 {
		for k, v := range e.extra {
			ej.Extra[fmt.Sprint(k)] = v
		}
	}
	content, err := json.Marshal(ej)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func (e element) Val() E {
	ej := new(E)
	ej.DateTime = e.dateTime
	if !e.level.IsNoLevel() {
		ej.Level = e.level.String()
	}
	ej.Msg = e.msg
	if e.c.ok {
		ej.File = e.c.file
		ej.Line = e.c.line
		ej.Func = e.c.fn
	}
	if len(e.extra) > 0 {
		for k, v := range e.extra {
			ej.Extra[fmt.Sprint(k)] = v
		}
	}
	return ej
}

type Extra map[string]interface{}

type E struct {
	Level 		string            	`json:"Level,omitempty"`
	DateTime 	time.Time      		`json:"Time,omitempty"`
	Msg 		string              	`json:"Message,omitempty"`
	File 		string             	`json:"File,omitempty"`
	Line 		int                	`json:"Line,omitempty"`
	Func 		string             	`json:"Func,omitempty"`
	Extra		 			`json:"Extra,omitempty"`
}
