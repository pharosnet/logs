package logs

import (
	"strings"
	"fmt"
)

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel

)

type Level uint32

func (level Level) LTE(lvl Level) bool {
	if level <= lvl {
		return true
	}
	return false
}

func (level Level) Color() int {
	var levelColor int
	switch level {
	case DebugLevel:
		levelColor = gray
	case InfoLevel:
		levelColor = blue
	case WarnLevel:
		levelColor = yellow
	case ErrorLevel:
		levelColor = red
	case FatalLevel:
		levelColor = red
	case PanicLevel:
		levelColor = red
	default:
		levelColor = nocolor
	}

	return levelColor
}

func (level Level) String() string {
	switch level {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warning"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	}

	return "unknown"
}

func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	}
	var l Level
	return l, fmt.Errorf("not a valid logrus Level: %q", lvl)
}



