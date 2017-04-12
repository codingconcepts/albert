package agent

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"github.com/codingconcepts/albert/pkg/model"
)

// Config holds the configuration variables an agent
// needs to run.
type Config struct {
	NatsHosts []string `json:"natsHosts"`
	NatsUser  string   `json:"natsUser"`
	NatsPass  string   `json:"natsPass"`

	LogLevel model.ConfigLogLevel `json:"logLevel"`

	Application        string                `json:"application"`
	ApplicationType    model.ApplicationType `json:"applicationType"`
	Identifier         string                `json:"identifier"`
	CustomInstructions []string              `json:"customInstructions"`
}

// NewConfigFromFile loads Agent configuration from a
// given file path and returns any errors encountered.
func NewConfigFromFile(path string) (c *Config, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}

	buffer := new(bytes.Buffer)
	if _, err = io.Copy(buffer, file); err != nil {
		return
	}

	c = new(Config)
	err = json.Unmarshal(buffer.Bytes(), c)

	return
}

func (c *Config) Validate() (err error) {
	if c.Application == "" {
		return ErrMissingApplication
	}

	if c.ApplicationType == model.UnknownApplicationType {
		return ErrMissingApplicationType
	}

	if c.ApplicationType == model.CustomApplicationType {
		if len(c.CustomInstructions) == 0 {
			return ErrMissingCustomInstructions
		}
	} else {
		if c.Identifier == "" {
			return ErrMissingIdentifier
		}
	}

	return
}
