package orchestrator

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codingconcepts/marcel/pkg/model"
	nats "github.com/nats-io/go-nats"
	"github.com/robfig/cron"
)

// Orchestrator holds the necessary information to process
// Application instances.
type Orchestrator struct {
	Connection   *nats.Conn
	Applications Applications

	Logger *logrus.Logger

	cronRunner     *cron.Cron
	gatherTimeout  time.Duration
	gatherChanSize int
}

// NewOrchestrator returns a pointer to a new instance of
// an Orchestrator.
func NewOrchestrator(c *Config) (o *Orchestrator) {
	return &Orchestrator{
		Connection:   c.Connection,
		Applications: c.Applications,
		Logger:       c.Logger,

		gatherTimeout:  c.GatherTimeout.Duration,
		gatherChanSize: c.GatherChanSize,
	}
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
	o.Connection.Close()
}

// Process makes a request for Applications and the performs a
// set of kill operations on them.
func (o *Orchestrator) Process(a Application) {
	agents, err := o.ScatterGather(a.Name)
	if err != nil {
		o.Logger.WithFields(logrus.Fields{
			"name": a.Name,
		}).WithError(err).Error("error occurred processing application")
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
		if err := o.IssueKillCommand(topic); err != nil {
			o.Logger.Error(err)
		}
	}
}

// ScatterGather performs a "scatter gather" operation against
// an unknown number of Applications.
// See http://bit.ly/2oEiquY for more information.
func (o *Orchestrator) ScatterGather(application string) (msgs []string, err error) {
	msgs = []string{}
	responses := make(chan *nats.Msg, o.gatherChanSize)
	defer close(responses)

	reply := nats.NewInbox()
	sub, err := o.Connection.ChanQueueSubscribe(reply, "", responses)
	if err != nil {
		return
	}
	defer sub.Unsubscribe()

	if err = o.Connection.PublishRequest(application, reply, nil); err != nil {
		return
	}

	for {
		select {
		case <-time.After(o.gatherTimeout):
			return
		case msg := <-responses:
			msgs = append(msgs, msg.Reply)
		}
	}
}

// IssueKillCommand publishes a kill command for a given Application
// and ApplicationType combination.
func (o *Orchestrator) IssueKillCommand(topic string) (err error) {
	o.Logger.WithFields(logrus.Fields{
		"topic": topic,
	}).Info("published kill")

	return o.Connection.Publish(topic, nil)
}
