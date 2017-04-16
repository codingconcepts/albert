package agent

import (
	"os"
	"os/exec"

	"github.com/Sirupsen/logrus"
	nats "github.com/nats-io/go-nats"
)

// Agent holds the necessary information to process listen
// for and respond to, instructions from the Orchestrator.
type Agent struct {
	Connection *nats.Conn
	KillInbox  string
	Logger     *logrus.Logger

	Application  string
	Instructions []string
}

// NewAgent returns a pointer to a new instance of an Agent.
func NewAgent(config *Config, conn *nats.Conn, logger *logrus.Logger) (a *Agent, err error) {
	if err = config.Validate(); err != nil {
		return
	}

	a = &Agent{
		Application:  config.Application,
		Instructions: config.Instructions,
		Connection:   conn,
		KillInbox:    nats.NewInbox(),
		Logger:       logger,
	}

	return
}

// Start begins the process of listening for instructions from
// the Orcestrator.
// NOTE:  Needs to be run in a goroutine
func (a *Agent) Start() {
	gatherChan, gatherStop, err := a.chanSubscribe(a.Application)
	if err != nil {
		a.Logger.Fatal(err)
	}
	defer gatherStop()

	killChan, killStop, err := a.chanSubscribe(a.KillInbox)
	if err != nil {
		a.Logger.Fatal(err)
	}
	defer killStop()

	for {
		select {
		case msg := <-gatherChan:
			a.GatherResponse(msg.Reply)
		case <-killChan:
			a.kill()
		}
	}
}

// Stop tears down the Agent.
func (a *Agent) Stop() {
	a.Connection.Close()
}

// GatherResponse responds to a scatter gather request with
// all of the necessary information an Orchestrator will need
// to decide what to kill.
func (a *Agent) GatherResponse(reply string) (err error) {
	a.Logger.WithFields(logrus.Fields{
		"reply": reply,
	}).Info("received scatter gather request")

	return a.Connection.PublishRequest(reply, a.KillInbox, []byte(a.Application))
}

func (a *Agent) chanSubscribe(topic string) (c chan *nats.Msg, stop func(), err error) {
	c = make(chan *nats.Msg)

	var sub *nats.Subscription
	if sub, err = a.Connection.ChanQueueSubscribe(topic, "", c); err != nil {
		return
	}

	stop = func() {
		if err := sub.Unsubscribe(); err != nil {
			a.Logger.Error(err)
		}

		close(c)
	}

	return
}

func (a *Agent) kill() {
	name := a.Instructions[0]

	var args []string
	if len(a.Instructions) > 1 {
		args = a.Instructions[1:]
	}

	a.Logger.WithFields(logrus.Fields{
		"name": name,
		"args": args,
	}).Info("kill")

	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		a.Logger.WithError(err).Error()
	} else {
		a.Logger.Info("success")
	}
}
