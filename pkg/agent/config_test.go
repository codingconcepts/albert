package agent_test

import (
	"encoding/json"
	"testing"

	"github.com/codingconcepts/albert/pkg/agent"
	"github.com/codingconcepts/albert/pkg/model"
	"github.com/codingconcepts/albert/test"
)

var (
	application        = "application"
	applicationType    = model.ProcessApplicationType
	identifier         = "identifier"
	customInstructions = []string{"one", "two", "three"}
)

func TestValidateMissingApplication(t *testing.T) {
	c := parseConfig(`
	{
		"applicationType": "process",
    	"identifier": "notepad.exe"
	}`)

	err := c.Validate()
	test.Equals(t, err, agent.ErrMissingApplication)
}

func TestValidateMissingApplicationType(t *testing.T) {
	c := parseConfig(`
	{
		"application": "notepad",
    	"identifier": "notepad.exe"
	}`)

	err := c.Validate()
	test.Equals(t, err, agent.ErrMissingApplicationType)
}

func TestValidateMissingIdentifier(t *testing.T) {
	c := parseConfig(`
	{
		"application": "notepad",
		"applicationType": "process"
	}`)

	err := c.Validate()
	test.Equals(t, err, agent.ErrMissingIdentifier)
}

func TestValidateMissingCustomInstructions(t *testing.T) {
	c := parseConfig(`
	{
		"application": "notepad",
		"applicationType": "custom"
	}`)

	err := c.Validate()
	test.Equals(t, err, agent.ErrMissingCustomInstructions)
}

func parseConfig(jsonConfig string) (c *agent.Config) {
	json.Unmarshal([]byte(jsonConfig), &c)
	return
}
