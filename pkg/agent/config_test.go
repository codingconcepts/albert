package agent_test

import (
	"encoding/json"
	"testing"

	"github.com/codingconcepts/albert/pkg/agent"
	"github.com/codingconcepts/albert/test"
)

var (
	application  = "notepad"
	instructions = []string{"taskkill", "/f", "/t", "/im", "PROCNAME.exe"}
)

func TestValidateWithValidConfig(t *testing.T) {
	c := agent.Config{
		Application:  application,
		Instructions: instructions,
	}

	err := c.Validate()
	test.ErrorNil(t, err)
	test.Equals(t, application, c.Application)
	test.Equals(t, instructions, c.Instructions)
}

func TestValidateEmptyConfig(t *testing.T) {
	c := agent.Config{}

	err := c.Validate()
	test.Equals(t, err, agent.ErrMissingApplication)
}

func TestValidateMissingApplication(t *testing.T) {
	c := agent.Config{
		Instructions: instructions,
	}

	err := c.Validate()
	test.Equals(t, err, agent.ErrMissingApplication)
}

func TestValidateMissingInstructions(t *testing.T) {
	c := agent.Config{
		Application: application,
	}

	err := c.Validate()
	test.Equals(t, err, agent.ErrMissingInstructions)
}

func parseConfig(jsonConfig string) (c *agent.Config) {
	json.Unmarshal([]byte(jsonConfig), &c)
	return
}
