package agent

import (
	"testing"

	"github.com/codingconcepts/albert/test"
)

func TestNewAgentConfigPropertiesAssigned(t *testing.T) {
	c := &Config{
		Application:  application,
		Instructions: instructions,
	}

	p := &mockProcessor{}
	k := &mockKiller{}
	a, err := NewAgent(c, p, k, logger)
	test.ErrorNil(t, err)

	test.Equals(t, c.Application, a.Application)
	test.Equals(t, c.Instructions, a.Instructions)
}

func TestNewAgentWithInvalidConfig(t *testing.T) {
	c := &Config{}

	p := &mockProcessor{}
	k := &mockKiller{}
	_, err := NewAgent(c, p, k, logger)
	test.ErrorNotNil(t, err)
}

func TestStartStop(t *testing.T) {
	a, _, _ := createTestAgent(t)

	go func() {
		a.Stop()
	}()

	a.Start()

	test.AnyLogEntryContainsMessage(t, "received stop signal", hook.Entries)
}

func TestGatherResponse(t *testing.T) {
	a, p, _ := createTestAgent(t)

	go func() {
		p.gatherChan <- "test message"
		a.Stop()
	}()

	a.Start()

	test.AnyLogEntryContainsField(t, "reply", "test message", hook.Entries)
	test.AnyLogEntryContainsMessage(t, "responded to scatter gather request", hook.Entries)
}

func TestKill(t *testing.T) {
	a, p, _ := createTestAgent(t)

	go func() {
		p.killChan <- struct{}{}
		a.Stop()
	}()

	a.Start()

	test.AnyLogEntryContainsMessage(t, "performed kill", hook.Entries)
}

func TestGatherResponseFailure(t *testing.T) {
	a, p, _ := createTestAgent(t)
	p.failToGather = true

	go func() {
		p.gatherChan <- "test message"
		a.Stop()
	}()

	a.Start()

	test.AnyLogEntryContainsField(t, "reply", "test message", hook.Entries)
	test.AnyLogEntryContainsMessage(t, "error occurred gathering", hook.Entries)
}

func TestKillFailure(t *testing.T) {
	a, p, k := createTestAgent(t)
	k.failToKill = true

	go func() {
		p.killChan <- struct{}{}
		a.Stop()
	}()

	a.Start()

	test.AnyLogEntryContainsMessage(t, "error occurred killing", hook.Entries)
}

func TestFailToGatherSubscribe(t *testing.T) {
	a, p, _ := createTestAgent(t)
	p.failToGatherSubcribe = true

	a.Start()

	// no test necessary here because if it doesn't fail,
	// we'll deadlock and that'll fail the test
}

func TestFailToKillSubscribe(t *testing.T) {
	a, p, _ := createTestAgent(t)
	p.failToKillSubcribe = true

	a.Start()

	// no test necessary here because if it doesn't fail,
	// we'll deadlock and that'll fail the test
}

func TestFailToStopGatherSubscription(t *testing.T) {
	a, p, _ := createTestAgent(t)
	p.failToStopGatherChan = true

	go func() {
		a.Stop()
	}()

	a.Start()

	test.AnyLogEntryContainsMessage(t, "failed to stop scatter gather subscriber", hook.Entries)
}

func TestFailToStopKillSubscription(t *testing.T) {
	a, p, _ := createTestAgent(t)
	p.failToStopKillChan = true

	go func() {
		a.Stop()
	}()

	a.Start()

	test.AnyLogEntryContainsMessage(t, "failed to stop kill subscriber", hook.Entries)
}
