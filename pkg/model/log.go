package model

import (
	"io"

	log "github.com/Sirupsen/logrus"
)

// NewLogger returns a pointer to a new logrus logger with all
// sensible defaults configured.
func NewLogger(out io.Writer, level log.Level) (logger *log.Logger) {
	logger = log.New()
	logger.Level = level
	logger.Out = out
	logger.Formatter = UTCFormatter{&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: true,
	}}

	return
}

// UTCFormatter formats logrus timestamps in UTC format
// rather than the default local format.
type UTCFormatter struct {
	log.Formatter
}

// Format takes a local timestamp and converts it to UTC.
func (u UTCFormatter) Format(e *log.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return u.Formatter.Format(e)
}
