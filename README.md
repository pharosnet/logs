
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
	"time"
	"context"
)

func main() {
    logs.DefaultLog().Infof("some info %v", time.Now())
	logs.DefaultLog().Debugf("some debug %v", time.Now())
	logs.DefaultLog().Warnf("some warn %v", time.Now())
	logs.DefaultLog().Errorf("some error %v", time.Now())
	logs.DefaultLog().Close(context.Background())
}

```

The customize way to use Logs is:
```go
package main

import (
	"github.com/pharosnet/logs"
	"os"
	"time"
	"context"
)

func main() {
    logger := logs.NewLogger(logs.NewStandardSource("std", logs.DebugLevel, logs.NewFlyChannel(logs.NewStandardSink())))
	logger.Infof("some info %v", time.Now())
	logger.Debugf("some debug %v", time.Now())
	logger.Warnf("some warn %v", time.Now())
	logger.Errorf("some error %v", time.Now())
	logger.Close(context.Background())    
}

```

The way to use json-output.

```go
package main

import (
	"context"
    "time"
    "os"
    "github.com/pharosnet/logs"
)

func main() {
    logger := logs.NewLogger(logs.NewStandardSource("std", logs.DebugLevel, logs.NewFlyChannel(logs.NewJsonSink())))
    logger.Infof("some info %v", time.Now())
    logger.Debugf("some debug %v", time.Now())
    logger.Warnf("some warn %v", time.Now())
    logger.Errorf("some error %v", time.Now())
    logger.Close(context.Background())
}

```

#### Thread safety

By default Logs is protected by memory barrier for concurrent writes, this memory barrier is invoked when calling send() and recv() of channel.


