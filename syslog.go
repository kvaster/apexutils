package apexutils

import (
	"fmt"
	"github.com/apex/log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	LOG_EMERG int = iota
	LOG_ALERT
	LOG_CRIT
	LOG_ERR
	LOG_WARNING
	LOG_NOTICE
	LOG_INFO
	LOG_DEBUG
)

const (
	LOG_KERN int = iota << 3
	LOG_USER
	LOG_MAIL
	LOG_DAEMON
	LOG_AUTH
	LOG_SYSLOG
	LOG_LPR
	LOG_NEWS
	LOG_UUCP
	LOG_CRON
	LOG_AUTHPRIV
	LOG_FTP
	_ // unused
	_ // unused
	_ // unused
	_ // unused
	LOG_LOCAL0
	LOG_LOCAL1
	LOG_LOCAL2
	LOG_LOCAL3
	LOG_LOCAL4
	LOG_LOCAL5
	LOG_LOCAL6
	LOG_LOCAL7
)

var syslogLevels = [...]int{
	log.DebugLevel: LOG_DEBUG,
	log.InfoLevel:  LOG_NOTICE,
	log.WarnLevel:  LOG_WARNING,
	log.ErrorLevel: LOG_ERR,
	log.FatalLevel: LOG_ERR,
}

type Syslog struct {
	c   net.Conn
	tag string
}

func newSyslog(tag string) (*Syslog, error) {
	c, err := net.Dial("unixgram", "/dev/log")
	if err != nil {
		return nil, err
	}
	return &Syslog{c: c, tag: tag}, nil
}

func (s *Syslog) HandleLog(e *log.Entry) error {
	b := &strings.Builder{}

	pr := LOG_DAEMON | syslogLevels[e.Level]
	timestamp := time.Now().Format(time.Stamp)

	_, _ = fmt.Fprintf(b, "<%d>%s %s[%d]: [%s] %s", pr, timestamp, s.tag, os.Getpid(), levelNames[e.Level], e.Message)

	for k, v := range e.Fields {
		_, _ = fmt.Fprintf(b, "; %s=%v", k, v)
	}

	m := b.String()

	_, err := s.c.Write([]byte(m))

	return err
}
