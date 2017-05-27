package logs

import (
	"runtime"
	"strings"
	"os"
)

func Trace() (string, int) {
	pc, _, line, ok := runtime.Caller(1)
	if !ok {
		return "?", 0
	}
	return runtime.FuncForPC(pc).Name(), line
}

func TraceFile() (string, string, int) {
	pc, f, line, ok := runtime.Caller(1)
	if !ok {
		return "?", "?", 0
	}
	return f, runtime.FuncForPC(pc).Name(), line
}

func TraceFileWithoutGoPath() (string, string, int) {
	pc, f, line, ok := runtime.Caller(1)
	if !ok {
		return "?", "?", 0
	}
	gopath, found := os.LookupEnv("GOPATH")
	if found {
		gopath = gopath + "/src"
		f = strings.Replace(f, gopath, "", 1)
	}
	return f, runtime.FuncForPC(pc).Name(), line
}

