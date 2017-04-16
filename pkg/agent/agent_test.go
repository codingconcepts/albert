package agent

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/codingconcepts/albert/test"
	nats "github.com/nats-io/go-nats"
)

func TestNewAgentConfigPropertiesAssigned(t *testing.T) {
	c := &Config{
		Application:  application,
		Instructions: instructions,
	}

	a, err := NewAgent(c, &nats.Conn{}, logrus.New())
	test.ErrorNil(t, err)

	test.Equals(t, c.Application, a.Application)
	test.Equals(t, c.Instructions, a.Instructions)
}

func TestNewAgentWithInvalidConfig(t *testing.T) {
	c := &Config{}

	_, err := NewAgent(c, &nats.Conn{}, logrus.New())
	test.ErrorNotNil(t, err)
}
