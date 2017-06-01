package logs

import (
	"os"
)

func New(level Level, out Writer) *logger {
	log := new(logger)
	log.level = level
	log.out = out
	log.hooks = make(map[Level][]Hook)
	return log
}

func NewWithHooks(level Level, out Writer, hooks []Hook) *logger {
	log := New(level, out)
	for _, hook := range hooks {
		log.addHook(hook)
	}
	return log
}

type logger struct {
	level Level
	out Writer
	hooks map[Level][]Hook
}


func (l *logger) addHook(hook Hook) {
	for _, level := range hook.Levels() {
		l.hooks[level] = append(l.hooks[level], hook)
	}
}

func (l *logger) Log(e Element) {
	if !e.Level().LTE(l.level) {
		return
	}
	l.output(e)
}

func (l *logger) Panic(e Element) {
	if !e.Level().LTE(l.level) {
		return
	}
	l.output(e)
	panic(e.String())
}

func (l *logger) Fatal(e Element) {
	if !e.Level().LTE(l.level) {
		return
	}
	l.output(e)
	os.Exit(1)
}


func (l *logger) output(e Element)  {
	l.out.Writer(e)
	l.fire(e)
}

func (l *logger) fire(e Element) {
	if l.hooks == nil || len(l.hooks) == 0 {
		return
	}
	for _, hook := range l.hooks[e.Level()] {
		if hook.IsAsyncFire() {
			go func(l *logger, hook Hook) {
				if err := hook.Fire(e); err != nil {
					hook.OnError(err)
				}
			}(l, hook)
		} else {
			if err := hook.Fire(e); err != nil {
				hook.OnError(err)
			}
		}
	}
}