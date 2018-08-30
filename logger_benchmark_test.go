package logs

import (
	"context"
	"runtime"
	"testing"
)

func BenchmarkNewLogger(b *testing.B) {
	runtime.GOMAXPROCS(3)
	defer runtime.GOMAXPROCS(1)
	logger := NewLogger(NewStandardSource("std", DebugLevel, NewFlyChannel(NewStandardSink())))
	for i := 0 ; i < b.N ; i ++ {
		logger.Infof("%v", i)
	}
	logger.Close(context.Background())
}
