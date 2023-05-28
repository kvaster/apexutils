package apexutils

import (
	"flag"
	"github.com/apex/log"
	"github.com/apex/log/handlers/json"
	termcolor "github.com/jwalton/go-supportscolor"
	"os"
)

var logFile = flag.String("log.file", "", "log file")
var logSyslog = flag.String("log.syslog", "", "syslog tag")
var logStdout = flag.Bool("log.stdout", false, "log to stdout")
var logStderr = flag.Bool("log.stderr", false, "log to stderr")
var logLevel = flag.String("log.level", "", "log level")
var logJson = flag.Bool("log.json", false, "use json format, valid only for stdout/stderr/file")

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

		defer log.SetLevel(l)
	}

	initialized := false

	if *logFile != "" {
		checkInitialized(&initialized)

		f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
		if err != nil {
			log.WithError(err).Fatal("error opening log file")
		}
		defer log.SetHandler(newPlain(f))
	}

	if *logSyslog != "" {
		checkInitialized(&initialized)

		h, err := newSyslog(*logSyslog)
		if err != nil {
			log.WithError(err).Error("error opening syslog")
		}
		defer log.SetHandler(h)
	}

	if *logStdout {
		checkInitialized(&initialized)

		if *logJson {
			defer log.SetHandler(json.New(os.Stdout))
		} else if termcolor.Stdout().SupportsColor {
			defer log.SetHandler(newColored(os.Stdout))
		} else {
			defer log.SetHandler(newPlain(os.Stdout))
		}
	}

	if *logStderr || !initialized {
		checkInitialized(&initialized)

		if *logJson {
			defer log.SetHandler(json.New(os.Stderr))
		} else if termcolor.Stderr().SupportsColor {
			defer log.SetHandler(newColored(os.Stderr))
		} else {
			defer log.SetHandler(newPlain(os.Stderr))
		}
	}
}

func checkInitialized(initialized *bool) {
	if *initialized {
		log.Fatal("only one of -log.file, -log.syslog, -log.stdout or -log.stderr can be enabled")
	}
	*initialized = true
}
