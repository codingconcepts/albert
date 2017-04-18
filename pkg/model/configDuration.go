package model

import (
	"strings"
	"time"
)

// ConfigDuration allows for the configuration of durations
// in the form of "5m30s", as opposed to the default Unix
// epoch timestamp.
type ConfigDuration struct {
	time.Duration
}

// UnmarshalJSON unmarshals a ConfigDuration from JSON.
func (d *ConfigDuration) UnmarshalJSON(b []byte) (err error) {
	d.Duration, err = time.ParseDuration(strings.Trim(string(b), `"`))
	return
}
