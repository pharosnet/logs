package logs

import (
	"os"
	"sync"
	"sync/atomic"
)

func New(lvl Level, w Writer) *logger {
	log := new(logger)
	log.mu = new(sync.Mutex)
	log.x = uint64(0)
	log.level = lvl
	log.w = w
	log.hooks = make(map[Level][]Hook)
	return log
}

type logger struct {
	mu *sync.Mutex
	x uint64
	level Level
	w Writer
	hooks map[Level][]Hook
}


func (l *logger) Add(hook Hook) {
	for _, level := range hook.Levels() {
		l.hooks[level] = append(l.hooks[level], hook)
	}
}

func (l *logger) Error(e *element)  {
	l.output(ErrorLevel, e)
}

func (l *logger) Warn(e *element)  {
	l.output(WarnLevel, e)
}

func (l *logger) Info(e *element)  {
	l.output(InfoLevel, e)
}

func (l *logger) Debug(e *element)  {
	l.output(DebugLevel, e)
}


func (l *logger) Fatal(v ...interface{})  {
	e := E()
	e.Msg(v...)
	l.output(FatalLevel, e)
	os.Exit(1)
}

func (l *logger) Fatalln(v ...interface{})  {
	e := E()
	e.Msg(v...)
	l.output(FatalLevel, e)
	os.Exit(1)
}

func (l *logger) Fatalf(format string, v ...interface{})  {
}

func (l *logger) Panic(v ...interface{})  {
}

func (l *logger) Panicln(v ...interface{})  {
}

func (l *logger) Panicf(format string, v ...interface{})  {
}

func (l *logger) Print(v ...interface{})  {
}

func (l *logger) Println(v ...interface{})  {
}

func (l *logger) Printf(format string, v ...interface{})  {
}

func (l *logger) output(lvl Level, e *element) {
	retryTimes := 5
	for !atomic.CompareAndSwapUint64(&l.x, l.x, l.x + uint64(1)) {
		retryTimes --
		if retryTimes < 0 {
			l.mu.Lock()
			e.now()
			l.mu.Unlock()
			break
		}
	}
	if e.dateTime.IsZero() {
		e.now()
	}
	// do write
	if l.level.LTE(lvl) {
		_, err := l.w.Write(l.w.Format(e))
		if err != nil && l.w.OnError != nil {
			l.w.OnError(err)
		}
	}
	// hook
	go l.fire(lvl, e)
}

func (l *logger) fire(level Level, e *element) error {
	for _, hook := range l.hooks[level] {
		if err := hook.Fire(e); err != nil {
			return err
		}
	}
	return nil
}