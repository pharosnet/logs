package logs

import (
	"io"
	"os"
)

type ErrorHandle func(int, error)

type Writer interface {
	io.Writer
	OnError(int, error)
}

// SimpleWriter <- file
func NewSimpleWriter(out io.Writer, newline bool, errorHandle ErrorHandle) *simpleWriter {
	w := new(simpleWriter)
	w.onError = errorHandle
	w.out = out
	w.newline = newline
	return
}

type simpleWriter struct {
	newline bool
	out io.Writer
	onError func(int, error)
}

func (w *simpleWriter) OnError(n int, err error)  {
	if w.onError != nil {
		w.onError(n, err)
	}
}

func (w *simpleWriter) Write(p []byte) (int, error) {
	if w.newline {
		return w.out.Write(append(p, byte('\n')))
	}
	return w.out.Write(p)
}

// StdoutWriter
func NewStdoutWriter(errorHandle ErrorHandle) *stdoutWriter {
	w := new(stdoutWriter)
	w.onError = errorHandle
	w.out = os.Stdout
	return
}

type stdoutWriter struct {
	out io.Writer
	onError func(int, error)
}

func (w *stdoutWriter) OnError(n int, err error)  {
	if w.onError != nil {
		w.onError(n, err)
	}
}

func (w *stdoutWriter) Write(p []byte) (int, error) {
	return w.out.Write(append(p, byte('\n')))
}

//
type LogWrappedWriter struct {
	out io.Writer
	decoder Decode
}

func (w *LogWrappedWriter) Write(p []byte) (int, error) {
	_, _, msg, newLine, err := w.decoder.Decode(p)
	if err != nil {
		return 0, err
	}
	e, parseErr := parseElement(msg)
	if parseErr != nil {
		return 0, parseErr
	}
	if newLine {
		return w.out.Write(e.bytesWithNewLine())
	}
	return w.out.Write(e.Bytes())
}