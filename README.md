
#### Description
 
* Logs is a structured logger for Go, and it itself is completely stable.
* Nicely color-coded in development (when a TTY is attached, otherwise just
  plain text).
* Logs can be used with the standard library logger.

#### Example

The simplest way to use Logs is simply the package-level exported logger:
```go
package main

import (
	"github.com/pharosnet/logs"
)

func main() {
	loggers := logs.New(logs.NewStdoutWriter(nil))
	loggers.Log(logs.Debugf("msg level : %s", "debug").WithField("k1","v1").WithField("k2", 2).Trace())
	loggers.Log(logs.Infof("msg level : %s", "info").WithField("k1","v1").WithField("k2", 2))
	loggers.Log(logs.Warnf("msg level : %v", "warn").WithField("k1","v1").WithField("k2", 2))
	loggers.Log(logs.Errorf("msg level : %v", "error").WithField("k1","v1").WithField("k2", 2).TraceFile())
	loggers.Panic(logs.Errorf("msg level : %v", "panic, it will call panic(logs.Element) and swap level with PanicLevel.").WithField("k1","v1").WithField("k2", 2).TraceFile())
	loggers.Fatal(logs.Errorf("msg level : %v", "fatal, it will call os.Exit(1) and swap level with FatalLevel.").WithField("k1","v1").WithField("k2", 2).TraceFile())
}

```

The way to use Logs with the standard library logger.

```go
package main

import (
    "log"
    "os"
	"github.com/pharosnet/logs"
)

func main() {
    logger := log.New(os.Stdout, "", 0) // prefix must be empty, and flag must be zero. in future, prefix and flag can be used.
	logger.Println(logs.Infof("msg %s", "some message").WithField("k1","v1").WithField("k2", 2).TraceFile())
}

```

#### Thread safety

By default Logs is protected by mutex for concurrent writes, this mutex is invoked when calling hooks and writing logs.
If you are sure such locking is not needed, like no fatal, you can call logger.SetNoLock() to disable the locking.

