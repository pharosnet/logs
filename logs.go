package logs

var defaultLog = NewLogger(NewStandardSource("default", InfoLevel, NewFlyChannel(NewStandardSink())))

func DefaultLog() Logger {
	return defaultLog
}
