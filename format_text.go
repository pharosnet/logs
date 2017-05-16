package logs

import (
	"bytes"
	"fmt"
)

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 34
	gray    = 37
)

type TextFormatter struct {}

func (f *TextFormatter) Format(e *element) []byte {
	buf := bytes.NewBufferString("")
	buf.WriteString(fmt.Sprintf("\x1b[%dm%s\x1b[0m[%s] %-44s ", e.level.Color(), e.level.String(), e.dateTime.String(), e.msg))
	for k, v := range e.data.data {
		buf.WriteString(fmt.Sprintf(" \x1b[%dm%s\x1b[0m=", green, k))
		buf.WriteString("=")
		buf.WriteString(fmt.Sprint(v))
		buf.WriteString(" ")
	}
	// call
	if e.c.ok {
		buf.WriteString(e.c.String())
	}

	return buf.Bytes()
}

/*
if f.DisableTimestamp {
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m %-44s ", levelColor, levelText, entry.Message)
	} else if !f.FullTimestamp {
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%04d] %-44s ", levelColor, levelText, int(entry.Time.Sub(baseTimestamp)/time.Second), entry.Message)
	} else {
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%s] %-44s ", levelColor, levelText, entry.Time.Format(timestampFormat), entry.Message)
	}
	for _, k := range keys {
		v := entry.Data[k]
		fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=", levelColor, k)
		f.appendValue(b, v)
	}

*/
