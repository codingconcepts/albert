package agent

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/codingconcepts/albert/pkg/model"
)

// Config holds the configuration variables an agent
// needs to run.
type Config struct {
	NatsHosts []string `json:"natsHosts"`
	NatsUser  string   `json:"natsUser"`
	NatsPass  string   `json:"natsPass"`

	LogLevel model.ConfigLogLevel `json:"logLevel"`

	Application  string   `json:"application"`
	Instructions []string `json:"instructions"`
}

// NewConfigFromReader loads Agent configuration from a
// given reader and returns any errors encountered.
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
	if c.Application == "" {
		return ErrMissingApplication
	}

	if len(c.Instructions) == 0 {
		return ErrMissingInstructions
	}

	return
}
