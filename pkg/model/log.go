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
		// always using the full timestamp formatter because the default
		// number-based timestamp formatter makes distributed logging
		// impossible.
		FullTimestamp: true,

		// always disabling console colours because they don't make for
		// easy reading in terminals that don't support it.
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
