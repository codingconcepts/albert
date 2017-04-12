package agent_test

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/codingconcepts/albert/pkg/agent"
	"github.com/codingconcepts/albert/pkg/model"
	"github.com/codingconcepts/albert/test"
	nats "github.com/nats-io/go-nats"
)

func TestAgentConfigPropertiesAssigned(t *testing.T) {
	c := &agent.Config{
		Application:        "application",
		ApplicationType:    model.CustomApplicationType,
		CustomInstructions: []string{"one", "two", "three"},
		Identifier:         "identifier",
	}

	a, err := agent.NewAgent(c, &nats.Conn{}, logrus.New())
	test.ErrorNil(t, err)

	test.Equals(t, c.Application, a.Application)
	test.Equals(t, c.ApplicationType, a.ApplicationType)
	test.Equals(t, c.Identifier, a.Identifier)
	test.Equals(t, c.CustomInstructions, a.CustomInstructions)
}
