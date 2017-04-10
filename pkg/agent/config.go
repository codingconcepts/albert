package agent

import (
	"github.com/Sirupsen/logrus"
	"github.com/codingconcepts/albert/pkg/model"
	nats "github.com/nats-io/go-nats"
)

// Config holds the configuration variables an agent
// needs to run.
type Config struct {
	NatsHosts []string `json:"natsHosts"`
	NatsUser  string   `json:"natsUser"`
	NatsPass  string   `json:"natsPass"`

	Logger   *logrus.Logger
	LogLevel model.ConfigLogLevel `json:"logLevel"`

	Connection *nats.Conn

	Application     string                `json:"application"`
	ApplicationType model.ApplicationType `json:"applicationType"`
	Identifier      string                `json:"identifier"`
}
