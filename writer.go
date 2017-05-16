package logs

import (
	"io"
	"os"
)

type Writer interface {
	Formatter
	io.Writer
	OnError(error)
}

func NewStdOut(onError func(error)) *stdWriter {
	w := new(stdWriter)
	w.f = new(TextFormatter)
	w.fn = onError()
	return
}

type stdWriter struct {
	f *TextFormatter
	fn func(error)
}

func (w *stdWriter) OnError(err error)  {
	if w.fn != nil {
		w.fn(err)
	}
}

func (w *stdWriter) Write(p []byte) (int, error) {
	return os.Stdout.Write(append(p, byte('\n')))
}
