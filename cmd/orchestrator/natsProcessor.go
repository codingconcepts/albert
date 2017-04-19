package main

import (
	"time"

	"github.com/codingconcepts/albert/pkg/orchestrator"
	nats "github.com/nats-io/go-nats"
)

type natsProcessor struct {
	conn           *nats.Conn
	gatherChanSize int
	gatherTimeout  time.Duration
}

// NewNatsProcessor returns the pointer to a new instance of a
// NatsProcessor.
func newNatsProcessor(conn *nats.Conn, config *orchestrator.Config) (p *natsProcessor) {
	return &natsProcessor{
		conn:           conn,
		gatherChanSize: config.GatherChanSize,
		gatherTimeout:  config.GatherTimeout.Duration,
	}
}

// Gather performs a "scatter gather" operation against
// an unknown number of Applications.
// See http://bit.ly/2oEiquY for more information.
func (p *natsProcessor) Gather(application string) (msgs []string, err error) {
	msgs = []string{}
	responses := make(chan *nats.Msg, p.gatherChanSize)
	defer close(responses)

	reply := nats.NewInbox()
	sub, err := p.conn.ChanQueueSubscribe(reply, "", responses)
	if err != nil {
		return
	}
	defer sub.Unsubscribe()

	if err = p.conn.PublishRequest(application, reply, nil); err != nil {
		return
	}

	for {
		select {
		case <-time.After(p.gatherTimeout):
			return
		case msg := <-responses:
			msgs = append(msgs, msg.Reply)
		}
	}
}

// IssueKill publishes a kill command for a given Application
// and ApplicationType combination.
func (p *natsProcessor) IssueKill(topic string) (err error) {
	return p.conn.Publish(topic, nil)
}
