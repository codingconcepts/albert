package orchestrator

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

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

// NewConfigFromFile loads Orchestrator configuration from a
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
