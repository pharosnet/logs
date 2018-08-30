package logs

import (
	"context"
	"errors"
	"fmt"
)

type Logger interface {
	Errorf(formatter string, args ...interface{})
	Error(args ...interface{})
	Warnf(formatter string, args ...interface{})
	Warn(args ...interface{})
	Infof(formatter string, args ...interface{})
	Info(args ...interface{})
	Debugf(formatter string, args ...interface{})
	Debug(args ...interface{})
	Close(ctx context.Context) error
}

type LS struct {
	L Level
	S Source
}


func NewLogger(sources ...Source) Logger {
	if len(sources) == 0 {
		return nil
	}
	return &standardLogger{sources: sources}
}

type standardLogger struct {
	sources []Source
}

func (l *standardLogger) Errorf(formatter string, args ...interface{}) {
	if formatter == "" || len(args) == 0 {
		return
	}
	for _, s := range l.sources {
		if s.Level().LTE(ErrorLevel) {
			s.Put(formatter, args)
		}
	}
}

func (l *standardLogger) Error(args ...interface{}) {
	if len(args) == 0 {
		return
	}
	for _, s := range l.sources {
		if s.Level().LTE(ErrorLevel) {
			s.Put(argsToFmt(args), args)
		}
	}
}

func (l *standardLogger) Warnf(formatter string, args ...interface{}) {
	if formatter == "" || len(args) == 0 {
		return
	}
	for _, s := range l.sources {
		if s.Level().LTE(WarnLevel) {
			s.Put(formatter, args)
		}
	}
}

func (l *standardLogger) Warn(args ...interface{}) {
	if len(args) == 0 {
		return
	}
	for _, s := range l.sources {
		if s.Level().LTE(WarnLevel) {
			s.Put(argsToFmt(args), args)
		}
	}
}

func (l *standardLogger) Infof(formatter string, args ...interface{}) {
	if formatter == "" || len(args) == 0 {
		return
	}
	for _, s := range l.sources {
		if s.Level().LTE(InfoLevel) {
			s.Put(formatter, args)
		}
	}
}

func (l *standardLogger) Info(args ...interface{}) {
	if len(args) == 0 {
		return
	}
	for _, s := range l.sources {
		if s.Level().LTE(InfoLevel) {
			s.Put(argsToFmt(args), args)
		}
	}
}

func (l *standardLogger) Debugf(formatter string, args ...interface{}) {
	if formatter == "" || len(args) == 0 {
		return
	}
	for _, s := range l.sources {
		if s.Level().LTE(DebugLevel) {
			s.Put(formatter, args)
		}
	}
}

func (l *standardLogger) Debug(args ...interface{}) {
	if len(args) == 0 {
		return
	}
	for _, s := range l.sources {
		if s.Level().LTE(DebugLevel) {
			s.Put(argsToFmt(args), args)
		}
	}
}

func (l *standardLogger) Close(ctx context.Context) error {
	buf := defaultPool.Get()
	defer defaultPool.Put(buf)
	for _, s := range l.sources {
		closeSourceErr := s.Close(ctx)
		if closeSourceErr != nil {
			buf.WriteString(fmt.Sprintf("%s: %v", s.Name(), closeSourceErr))
		}
	}
	if buf.Len() == 0 {
		return nil
	}
	return errors.New(buf.String())
}

func argsToFmt(args []interface{}) string {
	buf := defaultPool.Get()
	defer defaultPool.Put(buf)
	for i := 0 ; i < len(args) ; i ++ {
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString("%v")
	}
	return buf.String()
}