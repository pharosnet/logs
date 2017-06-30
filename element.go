package logs

import (
	"time"
	"bytes"
	"fmt"
	"runtime"
	"encoding/json"
	"strings"
	"strconv"
	"os"
)

var gopath string

func init()  {
	found := false
	gopath, found = os.LookupEnv("GOPATH")
	if found {
		gopath = gopath + "/src"
	}
}

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
			return fmt.Sprintf("[\x1b[%dm%s:%d\x1b[0m]", 36, c.fn, c.line)
		}
		return fmt.Sprintf("[\x1b[%dm%s:%s:%d\x1b[0m]", 36, c.file, c.fn, c.line)
	}
	return ""
}

//
type element struct {
	dateTime time.Time
	level Level
	msg string
	c caller
	extra Extra
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
	e.extra = Extra{}
	return e
}

// TODO
// [{LEVEL}][{TIME}][{MSG}][{EXTRA}][{CALLER}]
func ParseElement(content []byte) (Element, error) {
	level := ""
	datetime := ""
	msg := ""
	extra := ""
	call := ""
	for i := 0 ; i < len(content) ; i ++ {
		b := content[i]
		if b == 91 {
			if level == "" {
				for {
					i ++
					b = content[i]
					if b == 93 {
						break
					}
					level = level + string(b)
				}
				continue
			}
			if datetime == "" {
				for {
					i ++
					b = content[i]
					if b == 93 {
						break
					}
					datetime = datetime + string(b)
				}
				continue
			}
			if msg == "" {
				for {
					i ++
					b = content[i]
					if b == 93  {
						j := i
						if content[j + 1] == 91 && content[j + 2] == 123 {
							break
						}

					}
					msg = msg + string(b)
				}
				continue
			}
			if extra == "" {
				for {
					i ++
					b = content[i]
					if b == 93  {
						j := i
						if content[j - 1] == 125 {
							break
						}
					}
					extra = extra + string(b)
				}
				continue
			}
			if call == "" {
				for {
					i ++
					b = content[i]
					if b == 93  {
						break
					}
					call = call + string(b)
				}
				continue
			}
		}
	}
	e := new(element)
	lvl, lvlErr := ParseLevel(level)
	if lvlErr != nil {
		return nil, lvlErr
	}
	e.level = lvl
	date, dateErr := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", datetime)
	if dateErr != nil {
		return nil, dateErr
	}
	e.dateTime = date
	e.msg = msg
	if extra != "" { // {k=v}
		e.extra = Extra{}
		extra = extra[1:len(extra)-1]
		items := strings.Split(extra, "}{")
		for _, item := range items {
			kv := strings.Split(item, "=")
			e.extra[kv[0]] = kv[1]
		}
	}
	if call != "" { // /Users/doaman/workspace/projects/pharosnet/src/github.com/pharosnet/test/logs_std.go:func_name:16
		e.c = caller{}
		items := strings.Split(call, ":")
		if len(items) == 3 {
			e.c.file = items[0]
			e.c.fn = items[1]
			line, err := strconv.Atoi(items[2])
			if err != nil {
				return nil, err
			}
			e.c.line = line
		} else if len(items) == 2 {
			e.c.fn = items[0]
			line, err := strconv.Atoi(items[1])
			if err != nil {
				return nil, err
			}
			e.c.line = line
		}
		e.c.ok = true
	}
	return e, nil
}

func (e *element) Level() Level {
	return e.level
}

func (e *element) Trace() Element {
	pc, _, line, ok := runtime.Caller(1)
	c := caller{}
	if  ok {
		c.line = line
		c.fn = runtime.FuncForPC(pc).Name()
		c.ok = ok
	}
	e.c = c
	return e
}

func (e *element) TraceFile() Element {
	pc, f, line, ok := runtime.Caller(1)
	c := caller{}
	if  ok {
		c.file = f
		c.line = line
		c.fn = runtime.FuncForPC(pc).Name()
		c.ok = ok
	}
	e.c = c
	return e
}

func (e *element) TraceFileWithoutGoPath() Element {
	pc, f, line, ok := runtime.Caller(1)
	c := caller{}
	if  ok {
		c.file = strings.Replace(f, gopath, "", 1)
		c.line = line
		c.fn = runtime.FuncForPC(pc).Name()
		c.ok = ok
	}
	e.c = c
	return e
}

func (e *element) Extra(fields ...F) Element {
	for _, field := range fields {
		e.extra[field.K] = field.V
	}
	return e
}

// [{LEVEL}][{TIME}][{MSG}][{EXTRA}][{CALLER}]
func (e *element) Bytes() []byte {
	buf := bytes.NewBufferString("")
	// level
	if !e.level.IsNoLevel() {
		buf.WriteString(fmt.Sprintf("[\x1b[%dm%s\x1b[0m]", e.level.Color(), e.level.String()))
	}
	// datetime
	buf.WriteString(fmt.Sprintf("[\x1b[%dm%s\x1b[0m]", 37, e.dateTime.String())) // 36
	// msg
	buf.WriteString("[")
	buf.WriteString(e.msg)
	buf.WriteString("]")
	// extra
	if len(e.extra) > 0 {
		buf.WriteString("[")
		for k, v := range e.extra {
			buf.WriteString("{")
			buf.WriteString(fmt.Sprintf("\x1b[%dm%s\x1b[0m", e.level.Color(), k))
			buf.WriteString(fmt.Sprintf("\x1b[%dm%s\x1b[0m", 37, "="))
			buf.WriteString(fmt.Sprint(v))
			buf.WriteString("}")
		}
		buf.WriteString("]")
	}
	// caller
	if e.c.ok {
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
		ej.Extra = Extra{}
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
	ej := E{}
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

func (e *element) SetLevel(level Level) Element {
	e.level = level
	return e
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
