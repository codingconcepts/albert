package agent_test

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/codingconcepts/albert/pkg/agent"
	"github.com/codingconcepts/albert/test"
	nats "github.com/nats-io/go-nats"
)

func TestAgentConfigPropertiesAssigned(t *testing.T) {
	c := &agent.Config{
		Application:  "application",
		Instructions: []string{"one", "two", "three"},
	}

	a, err := agent.NewAgent(c, &nats.Conn{}, logrus.New())
	test.ErrorNil(t, err)

	test.Equals(t, c.Application, a.Application)
	test.Equals(t, c.Instructions, a.Instructions)
}
