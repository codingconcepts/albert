package model

import (
	"strings"

	"github.com/Sirupsen/logrus"
)

// ConfigLogLevel allows for the configuration of logrus
// log levels in the form of "debug", as opposed to the
// default numeric representation
type ConfigLogLevel struct {
	logrus.Level
}

// UnmarshalJSON unmarshals a ConfigLogLevel from JSON.
func (d *ConfigLogLevel) UnmarshalJSON(b []byte) (err error) {
	// trim off the quotes to get at the log level (the
	// JSON object will appear as "debug")
	raw := strings.Trim(string(b), `"`)
	d.Level, err = logrus.ParseLevel(raw)

	return
}
