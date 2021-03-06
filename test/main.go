package main

import (
	"flag"
	"github.com/apex/log"
	"github.com/kvaster/apexutils"
	"time"
)

// simple test
func main() {
	flag.Parse()
	apexutils.ParseFlags()

	for i := 0; i < 5; i++ {
		someLog()
		time.Sleep(time.Second * 3)
	}
}

func someLog() {
	for i := 0; i < 10; i++ {
		_i := i
		go func() {
			for j := 0; j < 20; j++ {
				log.WithField("v", _i*20+j).Info("info")
			}
		}()
	}
}
