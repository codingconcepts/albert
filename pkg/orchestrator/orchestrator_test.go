package orchestrator

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/codingconcepts/albert/test"
	nats "github.com/nats-io/go-nats"
)

func TestNewOrchestratorConfigPropertiesAssigned(t *testing.T) {
	c := &Config{
		Applications:   applications,
		GatherChanSize: gatherChanSize,
		GatherTimeout:  gatherTimeout,
	}

	o, err := NewOrchestrator(c, &nats.Conn{}, logrus.New())
	test.ErrorNil(t, err)

	test.Equals(t, c.Applications, o.Applications)
	test.Equals(t, c.GatherChanSize, o.gatherChanSize)
	test.Equals(t, c.GatherTimeout.Duration, o.gatherTimeout)
}

func TestNewOrchestratorWithInvalidConfig(t *testing.T) {
	c := &Config{}

	_, err := NewOrchestrator(c, &nats.Conn{}, logrus.New())
	test.ErrorNotNil(t, err)
}
