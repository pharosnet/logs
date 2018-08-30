package logs

import (
	"os"
	"strings"
	"sync"
)

var gopaths = make([]string, 0, 2)

var _gopathOnce = sync.Once{}

func goEnv() []string {
	_gopathOnce.Do(func() {
		// goroot
		goroot := os.Getenv("GOROOT")
		if len(goroot) > 0 {
			if strings.Contains(goroot, `\`) && strings.Contains(goroot, ":") { // win
				goroot = strings.Replace(goroot, `\`, "/", -1)
			}
			gopaths = append(gopaths, goroot)
		}
		gopath := os.Getenv("GOPATH")
		if len(gopath) == 0 {
			return
		}
		if strings.Contains(gopath, `\`) && strings.Contains(gopath, ":") { // win
			gopath = strings.Replace(gopath, `\`, "/", -1)
			if strings.Contains(gopath, ";") {
				gopaths := strings.Split(gopath, ";")
				for _, item := range gopaths {
					gopaths = append(gopaths, strings.TrimSpace(item))
				}
			} else {
				gopaths = append(gopaths, strings.TrimSpace(gopath))
			}
		} else { // unix
			if strings.Contains(gopath, ":") {
				gopaths := strings.Split(gopath, ":")
				for _, item := range gopaths {
					gopaths = append(gopaths, strings.TrimSpace(item))
				}
			} else {
				gopaths = append(gopaths, strings.TrimSpace(gopath))
			}
		}
	})
	return gopaths
}

