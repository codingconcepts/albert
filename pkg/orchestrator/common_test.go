package orchestrator

import (
	"errors"
	"time"

	logrustest "github.com/Sirupsen/logrus/hooks/test"

	"github.com/codingconcepts/albert/pkg/model"
)

var (
	logger, hook = logrustest.NewNullLogger()

	testApplication = Application{
		Name:       "notepad",
		Schedule:   "1 2 3 4 5 6",
		Percentage: 0.75,
	}

	gatherTimeout  = model.ConfigDuration{Duration: time.Second}
	gatherChanSize = 10
	applications   = Applications{
		testApplication,
	}
	config = `
	{
		"gatherTimeout": "1s",
		"gatherChanSize": 10,
		
		"applications": [
			{
				"name": "notepad",
				"schedule": "1 2 3 4 5 6",
				"percentage": 0.75
			}
		]
	}`

	// these are the same to make testing deterministic
	exampleTopic     = "example_topic"
	happyGatherResps = []string{exampleTopic, exampleTopic}
	errSadGather     = errors.New("sadGatherError")
	errSadIssueKill  = errors.New("sadIssueKill")
)

type mockProcessor struct {
	failToGather    bool
	failToIssueKill bool
}

func (p *mockProcessor) Gather(application string) (resps []string, err error) {
	if p.failToGather {
		err = errSadGather
	} else {
		resps = happyGatherResps
	}

	return
}

func (p *mockProcessor) IssueKill(topic string) (err error) {
	if p.failToIssueKill {
		err = errSadIssueKill
	}

	return
}
