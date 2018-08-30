package logs

import "context"

type Source interface {
	Name() string
	Level() Level
	Put(lv Level, formatter string, args []interface{}) error
	Close(ctx context.Context) error
}

func NewStandardSource(name string, lv Level, ch Channel) Source {
	return &standardSource{name:name, lv:lv, ch:ch}
}

type standardSource struct {
	name string
	lv Level
	ch Channel
}


func (s *standardSource) Name() string {
	return s.name
}

func (s *standardSource) Level() Level {
	return s.lv
}

func (s *standardSource) Put(lv Level, formatter string, args []interface{}) error {
		return s.ch.Send(pack(lv, formatter, args))
}

func (s *standardSource) Close(ctx context.Context) error {
	return s.ch.Close(ctx)
}
