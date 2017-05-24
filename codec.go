package logs

import (
	"time"
	"log"
	"bytes"
)

type Decode interface {
	Decode([]byte) (string, time.Time, []byte, bool, error)
}

// log.logger
type StdLogDecoder struct {
	HasPrefix bool
	Flags int
}

func (d *StdLogDecoder) Decode(p []byte) (prefix string, datetime time.Time, file string, msg []byte, newLine bool, err error) {
	if p == nil || len(p) == 0 {
		return
	}
	if p[len(p)-1] == '\n' {
		p = p[0:len(p)-1]
		newLine = true
	}
	buf := bytes.NewBuffer(p)
	// prefix
	if d.HasPrefix {
		token, readErr := buf.ReadString(' ')
		if readErr != nil {
			err = readErr
			return
		}
		prefix = token[0:len(token) - 1]
	}
	if d.Flags & (log.Ldate|log.Ltime|log.Lmicroseconds) != 0 {
		layout := ""
		logTime := ""
		if d.Flags & log.Ldate != 0 {
			layout = "2006/01/02"
			token, readErr := buf.ReadString(' ')
			if readErr != nil {
				err = readErr
				return
			}
			logTime = token[0:len(token) - 1]
		}
		if d.Flags & log.Ltime != 0 {
			if layout == "" {
				layout = "15:04:05"
			} else {
				layout = layout + " 15:04:05"
			}
		}
		if d.Flags & log.Lmicroseconds != 0 {
			if layout == "" {
				layout = "999999999"
			} else {
				layout = layout + ".999999999"
			}
		}
		token, readErr := buf.ReadString(' ')
		if readErr != nil {
			err = readErr
			return
		}
		logTime = logTime + " " + token[0:len(token) - 1]
		if d.Flags & log.LUTC != 0 {
			datetime, err = time.ParseInLocation(layout, logTime, time.UTC)
		} else {
			datetime, err = time.Parse(layout, logTime)
		}
		if err != nil {
			return
		}
	}
	if d.Flags & (log.Lshortfile|log.Llongfile) != 0 {
		token, readErr := buf.ReadString(' ')
		if readErr != nil {
			err = readErr
			return
		}
		file = token[0:len(token) - 1]
	}
	msg = buf.Bytes()
	return
}