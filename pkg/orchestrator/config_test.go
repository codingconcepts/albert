package orchestrator

import (
	"strings"
	"testing"

	"github.com/codingconcepts/albert/test"
)

func TestNewConfigFromReader(t *testing.T) {
	reader := strings.NewReader(config)

	c, err := NewConfigFromReader(reader)
	test.ErrorNil(t, err)
	test.Equals(t, applications, c.Applications)
	test.Equals(t, gatherChanSize, c.GatherChanSize)
	test.Equals(t, gatherTimeout, c.GatherTimeout)
}

func TestValidateWithValidConfig(t *testing.T) {
	c := Config{
		Applications:   applications,
		GatherChanSize: gatherChanSize,
		GatherTimeout:  gatherTimeout,
	}

	err := c.Validate()
	test.ErrorNil(t, err)
	test.Equals(t, applications, c.Applications)
	test.Equals(t, gatherChanSize, c.GatherChanSize)
	test.Equals(t, gatherTimeout, c.GatherTimeout)
}

func TestValidateEmptyConfig(t *testing.T) {
	c := Config{}

	err := c.Validate()
	test.Equals(t, err, ErrMissingApplications)
}

func TestValidateMissingApplications(t *testing.T) {
	c := Config{
		GatherChanSize: gatherChanSize,
		GatherTimeout:  gatherTimeout,
	}

	err := c.Validate()
	test.Equals(t, err, ErrMissingApplications)
}

func TestValidateMissingGatherTimeout(t *testing.T) {
	c := Config{
		Applications:   applications,
		GatherChanSize: gatherChanSize,
	}

	err := c.Validate()
	test.Equals(t, err, ErrMissingGatherTimeout)
}

func TestValidateMissingGatherChanSize(t *testing.T) {
	c := Config{
		Applications:  applications,
		GatherTimeout: gatherTimeout,
	}

	err := c.Validate()
	test.Equals(t, err, ErrInvalidGatherChanSize)
}
