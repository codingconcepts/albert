package agent

import (
	"errors"
	"testing"

	"github.com/Sirupsen/logrus"
	logrustest "github.com/Sirupsen/logrus/hooks/test"
	"github.com/codingconcepts/albert/pkg/model"
	"github.com/codingconcepts/albert/test"
)

var (
	logger, hook = logrustest.NewNullLogger()

	application  = "notepad"
	instructions = []string{"taskkill", "/f", "/t", "/im", "notepad.exe"}
	config       = `
	{
		"application": "notepad",
		"instructions": [ "taskkill", "/f", "/t", "/im", "notepad.exe" ]	
	}`

	errSadGather          = errors.New("sadGatherError")
	errSadGatherSubscribe = errors.New("sadGatherSubscribe")
	errSadKillSubscribe   = errors.New("sadKillSubscribe")
	errSadGatherChanStop  = errors.New("sadGatherChanStop")
	errSadKillChanStop    = errors.New("sadKillChanStop")
	errSadKill            = errors.New("errSadKill")
	errSadConfigReader    = errors.New("errSadConfigReder")
)

func createTestAgent(tb testing.TB) (a *Agent, p *mockProcessor, k *mockKiller) {
	c := &Config{
		Application:  application,
		Instructions: instructions,
		LogLevel: model.ConfigLogLevel{
			Level: logrus.DebugLevel,
		},
	}

	logger.Level = c.LogLevel.Level

	p = newMockProcessor()
	k = &mockKiller{}

	var err error
	if a, err = NewAgent(c, p, k, logger); err != nil {
		test.ErrorNil(tb, err)
	}

	return
}

type mockProcessor struct {
	failToGatherSubcribe bool
	failToKillSubcribe   bool

	failToGather bool

	gatherChan           chan string
	gatherChanStop       func() error
	failToStopGatherChan bool

	killChan           chan struct{}
	killChanStop       func() error
	failToStopKillChan bool
}

func newMockProcessor() (p *mockProcessor) {
	p = &mockProcessor{
		gatherChan: make(chan string),
		killChan:   make(chan struct{}),
	}

	p.gatherChanStop = func() (err error) {
		if p.failToStopGatherChan {
			return errSadGatherChanStop
		}
		close(p.gatherChan)
		return
	}

	p.killChanStop = func() (err error) {
		if p.failToStopKillChan {
			return errSadKillChanStop
		}
		close(p.killChan)
		return
	}

	return
}

func (p *mockProcessor) GatherSubscribe(topic string) (msgs chan string, stop func() error, err error) {
	if p.failToGatherSubcribe {
		err = errSadGatherSubscribe
		return
	}

	msgs = p.gatherChan
	stop = p.gatherChanStop
	return
}

func (p *mockProcessor) KillSubscribe(topic string) (msgs chan struct{}, stop func() error, err error) {
	if p.failToKillSubcribe {
		err = errSadKillSubscribe
		return
	}

	msgs = p.killChan
	stop = p.killChanStop
	return
}

func (p *mockProcessor) GatherResponse(orchInbox string, agentInbox string, application string) (err error) {
	if p.failToGather {
		err = errSadGather
		return
	}

	return
}

type mockKiller struct {
	failToKill bool
}

func (k *mockKiller) Kill(instructions []string) (err error) {
	if k.failToKill {
		err = errSadKill
	}

	return
}

type errorReader struct {
	err error
}

func newErrorReader(err error) (r *errorReader) {
	return &errorReader{
		err: err,
	}
}

func (r *errorReader) Read(data []byte) (n int, err error) {
	return 0, r.err
}
