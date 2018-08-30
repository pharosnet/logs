package main

import (
	"context"
	"github.com/pharosnet/logs"
	"time"
)

func main() {
	logger := logs.NewLogger(logs.NewStandardSource("std", logs.DebugLevel, logs.NewFlyChannel(logs.NewStandardSink())))
	logger.Infof("some info %v", time.Now())
	logger.Debugf("some debug %v", time.Now())
	logger.Warnf("some warn %v", time.Now())
	logger.Errorf("some error %v", time.Now())
	logger.Close(context.Background())
}
