package agent

import (
	"github.com/Sirupsen/logrus"
	nats "github.com/nats-io/go-nats"
)

// Agent holds the necessary information to process listen
// for and respond to, instructions from the Orchestrator.
type Agent struct {
	KillInbox string
	Logger    *logrus.Logger

	Application  string
	Instructions []string

	processor  Processor
	killer     Killer
	gatherChan chan *nats.Msg
	killChan   chan *nats.Msg
	stopSig    chan struct{}
}

// Processor defines the communicative behaviour of an Agent.
type Processor interface {
	GatherSubscribe(topic string) (msgs chan string, stop func() error, err error)
	KillSubscribe(topic string) (msgs chan struct{}, stop func() error, err error)
	GatherResponse(orchInbox string, agentInbox string, application string) (err error)
}

// Killer defines the behaviour of something that can kill
// something else.  Wooly enough?
type Killer interface {
	Kill(instructions []string) (err error)
}

// NewAgent returns a pointer to a new instance of an Agent.
func NewAgent(config *Config, processor Processor, killer Killer, logger *logrus.Logger) (a *Agent, err error) {
	if err = config.Validate(); err != nil {
		return
	}

	a = &Agent{
		Application:  config.Application,
		Instructions: config.Instructions,
		KillInbox:    nats.NewInbox(),
		Logger:       logger,
		processor:    processor,
		killer:       killer,
		stopSig:      make(chan struct{}),
	}

	return
}

// Start begins the process of listening for instructions from
// the Orcestrator.
// NOTE:  this function blocks, launch in a goroutine
func (a *Agent) Start() (err error) {
	gatherChan, gatherStop, err := a.processor.GatherSubscribe(a.Application)
	if err != nil {
		a.Logger.WithError(err).Error("failed to subscribe to scatter gather requests")
		return
	}
	defer func() {
		if err = gatherStop(); err != nil {
			a.Logger.WithError(err).Error("failed to stop scatter gather subscriber")
			return
		}
	}()

	killChan, killStop, err := a.processor.KillSubscribe(a.KillInbox)
	if err != nil {
		a.Logger.WithError(err).Error("failed to subscribe to kill requests")
		return
	}
	defer func() {
		if err = killStop(); err != nil {
			a.Logger.WithError(err).Error("failed to stop kill subscriber")
			return
		}
	}()

	// caller will block here
	a.listenLoop(gatherChan, killChan)
	return
}

// listenLoop selects between the given channels and a.stopSig
// and will block until a message is received on a.stopSig.
func (a *Agent) listenLoop(gatherChan chan string, killChan chan struct{}) {
	for {
		select {
		case reply := <-gatherChan:
			if err := a.processor.GatherResponse(reply, a.KillInbox, a.Application); err != nil {
				a.Logger.WithField("reply", reply).WithError(err).Error("error occurred gathering")
			} else {
				a.Logger.WithField("reply", reply).Debug("responded to scatter gather request")
			}
		case <-killChan:
			if err := a.killer.Kill(a.Instructions); err != nil {
				a.Logger.WithError(err).Error("error occurred killing")
			} else {
				a.Logger.Debug("performed kill")
			}
		case <-a.stopSig:
			a.Logger.Info("received stop signal")
			return
		}
	}
}

// Stop tears down the Agent by sending a message to a.stopSig.
func (a *Agent) Stop() {
	a.stopSig <- struct{}{}
}
