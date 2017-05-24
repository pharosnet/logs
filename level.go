package logs

import (
	"strings"
	"fmt"
)

const (
	NoLevel Level = iota
	PanicLevel
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

func (level Level) IsNoLevel() bool {
	if level == NoLevel {
		return true
	}
	return false
}

func (level Level) Color() int {
	var levelColor int
	switch level {
	case DebugLevel:
		levelColor = 36
	case InfoLevel:
		levelColor = 34
	case WarnLevel:
		levelColor = 33
	case ErrorLevel:
		levelColor = 31
	case FatalLevel:
		levelColor = 45
	case PanicLevel:
		levelColor = 35
	default:
		levelColor = 0
	}

	return levelColor
}

func (level Level) String() string {
	switch level {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	case PanicLevel:
		return "PANIC"
	}
	return "UNKNOWN"
}

func ParseLevel(lvl string) (Level, error) {
	switch strings.ToUpper(lvl) {
	case "PANIC":
		return PanicLevel, nil
	case "FATAL":
		return FatalLevel, nil
	case "ERROR":
		return ErrorLevel, nil
	case "WARN":
		return WarnLevel, nil
	case "INFO":
		return InfoLevel, nil
	case "DEBUG":
		return DebugLevel, nil
	}
	return NoLevel, fmt.Errorf("not a valid Level: %q", lvl)
}



