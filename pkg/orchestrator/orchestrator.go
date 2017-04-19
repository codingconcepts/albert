package orchestrator

import (
	"github.com/Sirupsen/logrus"
	"github.com/codingconcepts/albert/pkg/model"
	"github.com/robfig/cron"
)

// Orchestrator holds the necessary information to process
// Application instances.
type Orchestrator struct {
	Processor    Processor
	Applications Applications
	Logger       *logrus.Logger

	cronRunner *cron.Cron
}

// Processor defines the communicative behaviour of
// an Orchestrator.
type Processor interface {
	Gather(application string) (resps []string, err error)
	IssueKill(topic string) (err error)
}

// NewOrchestrator returns a pointer to a new instance of
// an Orchestrator.
func NewOrchestrator(c *Config, processor Processor, logger *logrus.Logger) (o *Orchestrator, err error) {
	if err = c.Validate(); err != nil {
		return
	}

	o = &Orchestrator{
		Processor:    processor,
		Applications: c.Applications,
		Logger:       logger,
	}

	return
}

// Start begins a number of jobs to process each of the applications
// configured in the Orchestrator's config file.
// NOTE:  Needs to be run in a goroutine
func (o *Orchestrator) Start() {
	o.cronRunner = cron.New()
	for _, c := range o.Applications {
		if err := o.cronRunner.AddFunc(c.Schedule, func() { o.Process(c) }); err != nil {
			o.Logger.Fatal(err)
		}
	}

	o.cronRunner.Run()
}

// Stop tears down the Orchestrator.
func (o *Orchestrator) Stop() {
	o.cronRunner.Stop()
}

// Process makes a request for Applications and the performs a
// set of kill operations on them.
func (o *Orchestrator) Process(a Application) {
	agents, err := o.Processor.Gather(a.Name)
	if err != nil {
		o.Logger.WithFields(logrus.Fields{
			"name": a.Name,
		}).WithError(err).Error("error occurred gathering applications")
		return
	}

	// select a number of applications at random to kill
	randomAgents := model.TakeRandom(agents, a.Percentage)

	o.Logger.WithFields(logrus.Fields{
		"totalCount": len(agents),
		"killCount":  len(randomAgents),
		"name":       a.Name,
	}).Info("scatter gather responses received")

	for _, topic := range randomAgents {
		o.Logger.WithFields(logrus.Fields{
			"topic": topic,
		}).Info("published kill")

		if err := o.Processor.IssueKill(topic); err != nil {
			o.Logger.WithFields(logrus.Fields{
				"name": a.Name,
			}).WithError(err).Error("error occurred issuing kill command")
		}
	}
}
