# Utils for apex log

Usage in go:

```go
package main

import "github.com/kvaster/apexutils/log"

func main() {
    flag.Parse()
    log.ParseFlags()
	
	log.Info("some log")
}
```

Available cli flags:

`-log.level <level>` - logging level (debug, info, warn, error), default is info

`-log.file <file>` - log to file

`-log.syslog <tag>` - log to syslog (`unixgram:/dev/log`)
