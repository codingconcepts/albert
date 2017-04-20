package agent

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/codingconcepts/albert/test"
)

func TestNewConfigFromReader(t *testing.T) {
	reader := strings.NewReader(config)

	c, err := NewConfigFromReader(reader)
	test.ErrorNil(t, err)
	test.Equals(t, application, c.Application)
	test.Equals(t, instructions, c.Instructions)
}

func TestNewConfigFromReaderWithError(t *testing.T) {
	reader := newErrorReader(errSadConfigReader)

	_, err := NewConfigFromReader(reader)
	test.ErrorNotNil(t, err)
}

func TestValidateWithValidConfig(t *testing.T) {
	c := Config{
		Application:  application,
		Instructions: instructions,
	}

	err := c.Validate()
	test.ErrorNil(t, err)
	test.Equals(t, application, c.Application)
	test.Equals(t, instructions, c.Instructions)
}

func TestValidateEmptyConfig(t *testing.T) {
	c := Config{}

	err := c.Validate()
	test.Equals(t, err, ErrMissingApplication)
}

func TestValidateMissingApplication(t *testing.T) {
	c := Config{
		Instructions: instructions,
	}

	err := c.Validate()
	test.Equals(t, err, ErrMissingApplication)
}

func TestValidateMissingInstructions(t *testing.T) {
	c := Config{
		Application: application,
	}

	err := c.Validate()
	test.Equals(t, err, ErrMissingInstructions)
}

func parseConfig(jsonConfig string) (c *Config) {
	json.Unmarshal([]byte(jsonConfig), &c)
	return
}
