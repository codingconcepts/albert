package orchestrator

import (
	"github.com/Sirupsen/logrus"
	"github.com/codingconcepts/marcel/pkg/model"
	nats "github.com/nats-io/go-nats"
)

// Config holds the configuration variables an orchestrator
// needs to run.
type Config struct {
	NatsHosts      []string             `json:"natsHosts"`
	NatsUser       string               `json:"natsUser"`
	NatsPass       string               `json:"natsPass"`
	GatherTimeout  model.ConfigDuration `json:"gatherTimeout"`
	GatherChanSize int                  `json:"gatherChanSize"`

	Logger     *logrus.Logger
	LogLevel   model.ConfigLogLevel `json:"logLevel"`
	Connection *nats.Conn

	Applications Applications `json:"applications"`
}
