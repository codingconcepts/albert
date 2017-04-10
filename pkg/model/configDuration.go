package model

import (
	"fmt"
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

// MarshalJSON marshals a ConfigDuration to JSON.
func (d ConfigDuration) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
}
