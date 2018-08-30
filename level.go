package logs

const (
	ErrorLevel = iota
	WarnLevel
	InfoLevel
	DebugLevel
)

type Level uint32

func (l Level) LTE(lv Level) bool {
	if l <= lv {
		return true
	}
	return false
}

func (l Level) LT(lv Level) bool {
	if l < lv {
		return true
	}
	return false
}

func (l Level) Color() int {
	var levelColor int
	switch l {
	case DebugLevel:
		levelColor = 36
	case InfoLevel:
		levelColor = 34
	case WarnLevel:
		levelColor = 33
	case ErrorLevel:
		levelColor = 31
	default:
		levelColor = 0
	}
	return levelColor
}

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBU"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERRO"
	}
	return "    "
}



