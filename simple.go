package apexutils

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/apex/log"
)

// levelColors.
const (
	none   = 0
	red    = 31
	green  = 32
	yellow = 33
	blue   = 34
	gray   = 37
)

// Colors mapping.
var levelColors = [...]int{
	log.DebugLevel: gray,
	log.InfoLevel:  blue,
	log.WarnLevel:  yellow,
	log.ErrorLevel: red,
	log.FatalLevel: red,
}

// Strings mapping.
var levelNames = [...]string{
	log.DebugLevel: "debug",
	log.InfoLevel:  "info",
	log.WarnLevel:  "warn",
	log.ErrorLevel: "error",
	log.FatalLevel: "fatal",
}

// Simple implementation.
type Simple struct {
	mu              sync.Mutex
	Writer          io.Writer
	TimestampFormat string
	colored         bool
}

func newColored(w io.Writer) *Simple {
	return &Simple{
		Writer:          w,
		TimestampFormat: "2006-01-02 15:04:05",
		colored:         true,
	}
}

func newPlain(w io.Writer) *Simple {
	return &Simple{
		Writer:          w,
		TimestampFormat: "2006-01-02 15:04:05",
		colored:         false,
	}
}

// HandleLog implements log.Simple.
func (s *Simple) HandleLog(e *log.Entry) error {
	color := levelColors[e.Level]
	level := levelNames[e.Level]
	names := e.Fields.Names()

	s.mu.Lock()
	defer s.mu.Unlock()
	ts := time.Now().Format(s.TimestampFormat)

	if s.colored {
		_, _ = fmt.Fprintf(s.Writer, "%s [\033[%dm%6s\033[0m] %s", ts, color, level, e.Message)

		for _, name := range names {
			_, _ = fmt.Fprintf(s.Writer, "; \033[%dm%s\033[0m=%v", color, name, e.Fields.Get(name))
		}
	} else {
		_, _ = fmt.Fprintf(s.Writer, "%s [%6s] %s", ts, level, e.Message)

		for _, name := range names {
			_, _ = fmt.Fprintf(s.Writer, "; %s=%v", name, e.Fields.Get(name))
		}
	}

	_, _ = fmt.Fprintln(s.Writer)

	return nil
}
