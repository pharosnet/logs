package main

import (
	"context"
	"fmt"
	"github.com/pharosnet/logs"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(8)
	logs.DefaultLog().Infof("some info %v", time.Now())
	logs.DefaultLog().Debugf("some debug %v", time.Now())
	logs.DefaultLog().Warnf("some warn %v", time.Now())
	logs.DefaultLog().Errorf("some error %v", time.Now())
	fmt.Println("close", logs.DefaultLog().Close(context.Background()))
}
