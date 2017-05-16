package logs

import (
	"time"
	"bytes"
	"fmt"
	"runtime"
)



func E() *element {
	return newElement()
}

func newElement() *element {
	e := new(element)
	e.data = newFields()
	return e
}

type element struct {
	dateTime time.Time
	level Level
	msg string
	data *fields
	c *caller
}

func (e *element) now() *element {
	e.dateTime = time.Now()
	return e
}

func (e *element) MsgF(format string, v ...interface{}) *element {
	e.msg = fmt.Sprintf(format, v...)
	return e
}


func (e *element) Msg(v ...interface{}) *element {
	e.msg = fmt.Sprint(v...)
	return e
}

func (e *element) Trace() *element {
	pc, f, line, ok := runtime.Caller(1)
	c := new(caller)
	if  ok {
		c.file = f
		c.line = line
		c.ok = ok
		c.fn = runtime.FuncForPC(pc).Name()
	}
	e.c = c
	return e
}

func (e *element) WithField(k, v interface{}) *element {
	m := e.data.data
	m[k] = v
	e.data.data = m
	return e
}

type fields struct {
	data map[interface{}]interface{}
}

func newFields() *fields {
	f := new(fields)
	f.data = make(map[interface{}]interface{})
	return f
}

func (f *fields) String() string {
	buf := bytes.NewBufferString("")
	for k, v := range f.data {
		buf.WriteString(fmt.Sprintf("\x1b[%dm%s\x1b[0m=", nocolor, k))
		buf.WriteString("=")
		buf.WriteString(fmt.Sprint(v))
		buf.WriteString(" ")
	}
	if buf.Len() > 0 {
		return buf.String()[0: buf.Len() - 1]
	}
	return ""
}

type caller struct {
	file string
	line int
	fn string
	ok bool
}

func (c *caller) String() string {
	if c.ok {
		return "[" + c.fn + ":" + fmt.Sprint(c.line) + "][" + c.file + ":" + fmt.Sprint(c.line) + "]"
	}
	return ""
}