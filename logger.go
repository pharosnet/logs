package logs

import (
	"sync"
	"bytes"
	"os"
)

func New(out Writer) *logger {
	log := new(logger)
	log.mu = new(sync.Mutex)
	log.x = uint64(0)
	log.out = out
	log.hooks = make(map[Level][]Hook)
	return log
}

type logger struct {
	mu *sync.Mutex
	x uint64
	out Writer
	hooks map[Level][]Hook
}


func (l *logger) Add(hook Hook) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, level := range hook.Levels() {
		l.hooks[level] = append(l.hooks[level], hook)
	}
}

func (l *logger) Log(e Element) {
	n, err := l.output(e)
	if err != nil {
		l.out.OnError(n, err)
	}
}

func (l *logger) Panic(e Element) {
	e.SetLevel(PanicLevel)
	n, err := l.output(e)
	if err != nil {
		l.out.OnError(n, err)
	}
	panic(e.String())
}

func (l *logger) Fatal(e Element) {
	e.SetLevel(FatalLevel)
	n, err := l.output(e)
	if err != nil {
		l.out.OnError(n, err)
	}
	os.Exit(1)
}


func (l *logger) output(e Element) (int64, error) {
	buf := getBuffer()
	defer putBuffer(buf)
	buf.ReadFrom(bytes.NewReader(e.Bytes()))
	n, err := buf.WriteTo(l.out)
	if err != nil {
		return n, err
	}
	// hook
	go l.fire(e)
	return n, nil
}

func (l *logger) fire(e Element) {
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