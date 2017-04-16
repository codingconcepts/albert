package orchestrator

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/codingconcepts/albert/pkg/model"
)

// Config holds the configuration variables an orchestrator
// needs to run.
type Config struct {
	NatsHosts      []string             `json:"natsHosts"`
	NatsUser       string               `json:"natsUser"`
	NatsPass       string               `json:"natsPass"`
	GatherTimeout  model.ConfigDuration `json:"gatherTimeout"`
	GatherChanSize int                  `json:"gatherChanSize"`

	LogLevel model.ConfigLogLevel `json:"logLevel"`

	Applications Applications `json:"applications"`
}

// NewConfigFromReader loads Orchestrator configuration from a
// given file path and returns any errors encountered.
func NewConfigFromReader(reader io.Reader) (c *Config, err error) {
	buffer := new(bytes.Buffer)
	if _, err = io.Copy(buffer, reader); err != nil {
		return
	}

	c = new(Config)
	err = json.Unmarshal(buffer.Bytes(), c)

	return
}

// Validate performs basic config-time validation.
func (c *Config) Validate() (err error) {
	if len(c.Applications) == 0 {
		return ErrMissingApplications
	}

	if c.GatherTimeout.Duration == time.Duration(0) {
		return ErrMissingGatherTimeout
	}

	if c.GatherChanSize == 0 {
		return ErrInvalidGatherChanSize
	}

	return
}
