package agent

import (
	"github.com/Sirupsen/logrus"
	"github.com/codingconcepts/albert/pkg/model"
	nats "github.com/nats-io/go-nats"
)

// Agent holds the necessary information to process listen
// for and respond to, instructions from the Orchestrator.
type Agent struct {
	// application is the name of the name of the application
	// group that this agent serves.  This must match with the
	// name of the application group expected by the orchestrator.
	application string

	// instructions is a slice of strings used to execute commands.
	// If for example the cmdKiller is used to carry out kills,
	// Instructions[0] will be the name of the command line
	// application, while Instructions [1:] represent the optional
	// arguments to pass to it.
	instructions []string

	logger *logrus.Logger

	processor  Processor
	killer     Killer
	inbox      string
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

	inbox, err := model.InboxName(model.InboxMaxLength)
	if err != nil {
		return
	}

	a = &Agent{
		application:  config.Application,
		instructions: config.Instructions,
		logger:       logger,
		inbox:        inbox,
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
	gatherChan, gatherStop, err := a.processor.GatherSubscribe(a.application)
	if err != nil {
		a.logger.WithError(err).Error("failed to subscribe to scatter gather requests")
		return
	}
	defer func() {
		if err = gatherStop(); err != nil {
			a.logger.WithError(err).Error("failed to stop scatter gather subscriber")
			return
		}
	}()

	killChan, killStop, err := a.processor.KillSubscribe(a.inbox)
	if err != nil {
		a.logger.WithError(err).Error("failed to subscribe to kill requests")
		return
	}
	defer func() {
		if err = killStop(); err != nil {
			a.logger.WithError(err).Error("failed to stop kill subscriber")
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
			if err := a.processor.GatherResponse(reply, a.inbox, a.application); err != nil {
				a.logger.WithField("reply", reply).WithError(err).Error("error occurred gathering")
			} else {
				a.logger.WithField("reply", reply).Debug("responded to scatter gather request")
			}
		case <-killChan:
			if err := a.killer.Kill(a.instructions); err != nil {
				a.logger.WithError(err).Error("error occurred killing")
			} else {
				a.logger.Debug("performed kill")
			}
		case <-a.stopSig:
			a.logger.Info("received stop signal")
			return
		}
	}
}

// Stop tears down the Agent by sending a message to a.stopSig.
func (a *Agent) Stop() {
	a.stopSig <- struct{}{}
}
