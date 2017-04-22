package orchestrator

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/codingconcepts/albert/test"
)

func TestNewOrchestratorConfigPropertiesAssigned(t *testing.T) {
	c := &Config{
		Applications:   applications,
		GatherChanSize: gatherChanSize,
		GatherTimeout:  gatherTimeout,
	}

	p := &mockProcessor{}
	o, err := NewOrchestrator(c, p, logrus.New())
	test.ErrorNil(t, err)

	test.Equals(t, c.Applications, o.applications)
}

func TestNewOrchestratorWithInvalidConfig(t *testing.T) {
	c := &Config{}

	p := &mockProcessor{}
	_, err := NewOrchestrator(c, p, logrus.New())
	test.ErrorNotNil(t, err)
}

func TestProcessHappyPath(t *testing.T) {
	c := &Config{
		Applications:   applications,
		GatherChanSize: gatherChanSize,
		GatherTimeout:  gatherTimeout,
	}

	p := &mockProcessor{}
	o, err := NewOrchestrator(c, p, logger)
	test.ErrorNil(t, err)

	hook.Reset()
	o.Process(testApplication)

	test.AnyLogEntryContainsMessage(t, "scatter gather responses received", hook.Entries)
	test.AnyLogEntryContainsField(t, "totalCount", 2, hook.Entries)
	test.AnyLogEntryContainsField(t, "killCount", 1, hook.Entries)
	test.AnyLogEntryContainsField(t, "name", "notepad", hook.Entries)
	test.AnyLogEntryContainsMessage(t, "published kill", hook.Entries)
	test.AnyLogEntryContainsField(t, "topic", exampleTopic, hook.Entries)
}

func TestProcessFailToGather(t *testing.T) {
	processFailure(t, &mockProcessor{failToGather: true}, errSadGather)

	test.LogEntryContainsField(t, "error", errSadGather.Error(), hook.LastEntry())
}

func TestProcessFailToIssueKill(t *testing.T) {
	processFailure(t, &mockProcessor{failToIssueKill: true}, errSadIssueKill)

	test.LogEntryContainsField(t, "error", errSadIssueKill.Error(), hook.LastEntry())
}

func processFailure(t *testing.T, p Processor, expError error) {
	c := &Config{
		Applications:   applications,
		GatherChanSize: gatherChanSize,
		GatherTimeout:  gatherTimeout,
	}

	o, err := NewOrchestrator(c, p, logger)
	test.ErrorNil(t, err)

	o.Process(testApplication)

	test.LogEntryContainsField(t, "error", expError.Error(), hook.LastEntry())
}
