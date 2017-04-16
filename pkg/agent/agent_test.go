package agent

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/codingconcepts/albert/test"
	nats "github.com/nats-io/go-nats"
)

func TestAgentConfigPropertiesAssigned(t *testing.T) {
	c := &Config{
		Application:  application,
		Instructions: instructions,
	}

	a, err := NewAgent(c, &nats.Conn{}, logrus.New())
	test.ErrorNil(t, err)

	test.Equals(t, c.Application, a.Application)
	test.Equals(t, c.Instructions, a.Instructions)
}
