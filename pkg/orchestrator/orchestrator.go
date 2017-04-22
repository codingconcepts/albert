package orchestrator

import (
	"github.com/Sirupsen/logrus"
	"github.com/codingconcepts/albert/pkg/model"
	"github.com/robfig/cron"
)

// Orchestrator holds the necessary information to process
// Application instances.
type Orchestrator struct {
	// processor provides a basic contract by with the
	// orchestrator is able to request agents for application
	// groups and issue kill orders to selected agents.
	processor Processor

	// applications is a collection of application groups
	// this orchestrator is allowed to control.  For each
	// application group, there will be zero or more agents.
	applications Applications

	logger *logrus.Logger

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
		processor:    processor,
		applications: c.Applications,
		logger:       logger,
	}

	return
}

// Start begins a number of jobs to process each of the applications
// configured in the Orchestrator's config file.
// NOTE:  Needs to be run in a goroutine
func (o *Orchestrator) Start() {
	o.cronRunner = cron.New()
	for _, c := range o.applications {
		if err := o.cronRunner.AddFunc(c.Schedule, func() { o.Process(c) }); err != nil {
			o.logger.Fatal(err)
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
	agents, err := o.processor.Gather(a.Name)
	if err != nil {
		o.logger.WithFields(logrus.Fields{
			"name": a.Name,
		}).WithError(err).Error("error occurred gathering applications")
		return
	}

	// select a number of applications at random to kill
	randomAgents := model.TakeRandom(agents, a.Percentage)

	o.logger.WithFields(logrus.Fields{
		"totalCount": len(agents),
		"killCount":  len(randomAgents),
		"name":       a.Name,
	}).Info("scatter gather responses received")

	for _, topic := range randomAgents {
		o.logger.WithFields(logrus.Fields{
			"topic": topic,
		}).Info("published kill")

		if err := o.processor.IssueKill(topic); err != nil {
			o.logger.WithFields(logrus.Fields{
				"name": a.Name,
			}).WithError(err).Error("error occurred issuing kill command")
		}
	}
}
