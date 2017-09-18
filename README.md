
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
	"os"
)

func main() {
	loggers := logs.New(logs.DebugLevel, logs.NewByteBufferWriter(os.Stdout, nil))
	loggers.Log(logs.Debugf("msg level : %s", "debug").Extra(logs.F{"k1", "v1"}, logs.F{"k2", 2}).CallFile())
	loggers.Log(logs.Infof("msg level : %s", "info").Extra(logs.F{"k1", "v1"}, logs.F{"k2", 2}).CallFunc())
	loggers.Log(logs.Warnf("msg level : %v", "warn").Extra(logs.F{"k1", "v1"}, logs.F{"k2", 2}))
	loggers.Log(logs.Errorf("msg level : %v", "error").Extra(logs.F{"k1", "v1"}, logs.F{"k2", 2}).CallFileWithGoPath())
	loggers.Panic(logs.Errorf("msg level : %v", "panic, it will call panic(logs.Element) and swap level with PanicLevel.").Extra(logs.F{"k1", "v1"}, logs.F{"k2", 2}))
	loggers.Fatal(logs.Errorf("msg level : %v", "fatal, it will call os.Exit(1) and swap level with FatalLevel.").Extra(logs.F{"k1", "v1"}, logs.F{"k2", 2}))
}

```

The simplest way to use Logs is simply the package-level exported logger in async model:
```go
package main

import (
	"github.com/pharosnet/logs"
	"os"
)

func main() {
    w := logs.NewAsyncByteBufferWriter(os.Stdout, nil, 64)
    loggers := logs.New(logs.DebugLevel, w)
    defer w.Flush()
    loggers.Log(logs.Debugf("msg level : %s", "debug").Extra(logs.F{"k1", "v1"}, logs.F{"k2", 2}).CallFile())
    loggers.Log(logs.Infof("msg level : %s", "info").Extra(logs.F{"k1", "v1"}, logs.F{"k2", 2}).CallFunc())
    loggers.Log(logs.Warnf("msg level : %v", "warn").Extra(logs.F{"k1", "v1"}, logs.F{"k2", 2}))
    loggers.Log(logs.Errorf("msg level : %v", "error").Extra(logs.F{"k1", "v1"}, logs.F{"k2", 2}).CallFileWithGoPath())
    loggers.Panic(logs.Errorf("msg level : %v", "panic, it will call panic(logs.Element) and swap level with PanicLevel.").Extra(logs.F{"k1", "v1"}, logs.F{"k2", 2}))
    loggers.Fatal(logs.Errorf("msg level : %v", "fatal, it will call os.Exit(1) and swap level with FatalLevel.").Extra(logs.F{"k1", "v1"}, logs.F{"k2", 2}))
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
    logger.Println(logs.Infof("msg %s", "some message").Extra(logs.F{"k1", "v1"}, logs.F{"k2", 2}).CallFile())
}

```

#### Thread safety

By default Logs is protected by mutex for concurrent writes, this mutex is invoked when calling hooks and writing logs.
If you are sure such locking is not needed, like no fatal, you can call logger.SetNoLock() to disable the locking.

