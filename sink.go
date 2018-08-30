package logs

import (
	"context"
	"fmt"
	"os"
)

type Sink interface {
	FlowInto(pac *Packet) error
	Close(ctx context.Context) error
}

func NewStandardSink() Sink {
	return &standardSink{}
}

type standardSink struct {}

func (s *standardSink) FlowInto(pac *Packet) error {
	buf := getBuffer()
	defer putBuffer(buf)
	buf.WriteString(fmt.Sprintf("[ \x1b[%dm%s\x1b[0m ]", pac.Lv.Color(), pac.Lv.String()))
	buf.WriteString(fmt.Sprintf("[ \x1b[%dm%s\x1b[0m ]", 37, pac.Occurred.Format("2006-01-02 15:04:05.999999999 -0700 MST"))) // 36
	buf.WriteString(fmt.Sprintf("[ %d ]", pac.Gid))
	_, fn, file, line := pac.Caller()
	buf.WriteString(fmt.Sprintf("[ %s ][ %s:%d ][ ", fn, file, line))
	buf.WriteString(fmt.Sprintf(pac.Formatter, pac.Elements...))
	buf.WriteByte(' ')
	buf.WriteByte(']')

	buf.WriteByte('\n')
	var wErr error = nil
	_, wErr = buf.WriteTo(os.Stdout)
	return wErr
}

func (s *standardSink) Close(ctx context.Context) error {
	return nil
}


func NewJsonSink() Sink {
	return &jsonSink{}
}

type jsonSink struct {}

func (s *jsonSink) FlowInto(pac *Packet) error {
	buf := getBuffer()
	defer putBuffer(buf)
	_, fn, file, line := pac.Caller()
	buf.WriteString(
		fmt.Sprintf(`{"level":"%s","fn":"%s","file":"%s","line":%d,"msg":"%s"}`,
			pac.Lv.String(), fn, file, line, fmt.Sprintf(pac.Formatter, pac.Elements...)))
	buf.WriteByte('\n')
	_, wErr := buf.WriteTo(os.Stdout)
	return wErr
}


func (s *jsonSink) Close(ctx context.Context) error {
	return nil
}