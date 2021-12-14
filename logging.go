package apexutils

import (
	"flag"
	"github.com/apex/log"
	termcolor "github.com/jwalton/go-supportscolor"
	"os"
)

var logFile = flag.String("log.file", "", "log file")
var logSyslog = flag.String("log.syslog", "", "syslog tag")
var logLevel = flag.String("log.level", "", "log level")

func init() {
	if termcolor.Stderr().SupportsColor {
		log.SetHandler(newColored(os.Stderr))
	} else {
		log.SetHandler(newPlain(os.Stderr))
	}
	log.SetLevel(log.InfoLevel)
}

func ParseFlags() {
	if *logLevel != "" {
		l, err := log.ParseLevel(*logLevel)
		if err != nil {
			log.WithField("level", *logLevel).Fatal("error parsing log level")
		}

		log.SetLevel(l)
	}

	if *logFile != "" && *logSyslog != "" {
		log.Fatal("only one of -log.file or -log.syslog can be enabled")
	}

	if *logFile != "" {
		f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
		if err != nil {
			log.WithError(err).Fatal("error opening log file")
		}
		log.SetHandler(newPlain(f))
	}

	if *logSyslog != "" {
		h, err := newSyslog(*logSyslog)
		if err != nil {
			log.WithError(err).Error("error opening syslog")
		}
		log.SetHandler(h)
	}
}
