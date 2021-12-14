package log

import (
	"flag"
	_log "github.com/apex/log"
	termcolor "github.com/jwalton/go-supportscolor"
	"os"
)

var logFile = flag.String("log.file", "", "log file")
var logSyslog = flag.String("log.syslog", "", "syslog tag")
var logLevel = flag.String("log.level", "", "log level")

func init() {
	if termcolor.Stderr().SupportsColor {
		_log.SetHandler(newColored(os.Stderr))
	} else {
		_log.SetHandler(newPlain(os.Stderr))
	}
	_log.SetLevel(_log.InfoLevel)
}

func ParseFlags() {
	if *logLevel != "" {
		l, err := _log.ParseLevel(*logLevel)
		if err != nil {
			_log.WithField("level", *logLevel).Fatal("error parsing log level")
		}

		_log.SetLevel(l)
	}

	if *logFile != "" && *logSyslog != "" {
		_log.Fatal("only one of -log.file or -log.syslog can be enabled")
	}

	if *logFile != "" {
		f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
		if err != nil {
			_log.WithError(err).Fatal("error opening log file")
		}
		_log.SetHandler(newPlain(f))
	}

	if *logSyslog != "" {
		h, err := newSyslog(*logSyslog)
		if err != nil {
			_log.WithError(err).Error("error opening syslog")
		}
		_log.SetHandler(h)
	}
}

type Interface = _log.Interface
type Fields = _log.Fields

var Log = _log.Log
var WithFields = _log.WithFields
var WithField = _log.WithField
var WithError = _log.WithError
var Debug = _log.Debug
var Info = _log.Info
var Warn = _log.Warn
var Error = _log.Error
var Fatal = _log.Fatal
var Debugf = _log.Debugf
var Infof = _log.Infof
var Warnf = _log.Warnf
var Errorf = _log.Errorf
var Fatalf = _log.Fatalf
var Trace = _log.Trace
