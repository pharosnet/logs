package logs

import (
	"path"
	"runtime"
	"strings"
	"time"
)

const (
	skip = 3
	depth = 1
	unknown = "unknown"
)

type Packet struct {
	Gid int64
	PC uintptr
	Occurred time.Time
	Lv Level
	Formatter string
	Elements []interface{}
}

func (pac *Packet) Caller() (home string, fn string, file string, line int) {
	fnPC := runtime.FuncForPC(pac.PC)
	if fnPC == nil {
		home = unknown
		fn = unknown
		file = unknown
	} else {
		filename, lineNo := fnPC.FileLine(pac.PC)
		home, file = fileName(filename)
		line = lineNo
		fn = fnPC.Name()
	}
	return
}

func pack(lv Level, formatter string, args []interface{}) *Packet {
	return &Packet{
		Gid: currentGoruntineId(),
		PC: callers(),
		Occurred: time.Now(),
		Lv :lv,
		Formatter:formatter,
		Elements:args[:],
	}
}

func currentGoruntineId() int64 {
	return 0
}

func callers() uintptr {
	pcs := [depth]uintptr{}
	_ = runtime.Callers(skip, pcs[:])
	return pcs[0]
}

func fileName(src string) (goPath string, file string) {
	file = src
	goHomes := goEnv()
	if goHomes == nil {
		return
	}
	for _, goHome := range goHomes {
		if strings.Contains(file, goHome) {
			goPath = goHome
			file = strings.Replace(file, path.Join(goHome, "src"), "", 1)[1:]
			return
		}
	}
	return
}