package main

import nats "github.com/nats-io/go-nats"

type natsProcessor struct {
	conn *nats.Conn
}

// newNatsProcessor returns a pointer to a new instance of
// a natsProcessor, passing the connection it'll use.
func newNatsProcessor(conn *nats.Conn) (p *natsProcessor) {
	return &natsProcessor{
		conn: conn,
	}
}

// GatherSubscribe performs a scatter-gather request against
// a NATS cluster and pipes the NATS messages received into a
// simpler channel for later consumption.
func (p *natsProcessor) GatherSubscribe(topic string) (c chan string, stop func() error, err error) {
	c = make(chan string)

	var natsMsgs chan *nats.Msg
	if natsMsgs, stop, err = p.subscribe(topic); err != nil {
		return
	}

	// pipe messages from one channel into another
	go func() {
		for msg := range natsMsgs {
			c <- msg.Reply
		}
	}()

	return
}

// KillSubscribe subscribes for kill requests from the Orchestrator
// and pipes the NATS messages received into a simpler channel for
// later consumption.
func (p *natsProcessor) KillSubscribe(topic string) (c chan struct{}, stop func() error, err error) {
	c = make(chan struct{})

	var natsMsgs chan *nats.Msg
	if natsMsgs, stop, err = p.subscribe(topic); err != nil {
		return
	}

	// pipe messages from one channel into another
	go func() {
		for range natsMsgs {
			c <- struct{}{}
		}
	}()

	return
}

// subscribe creates a NATS subscription and returns the channel created
// and a function that can be used to tear everything down.
func (p *natsProcessor) subscribe(topic string) (msgs chan *nats.Msg, stop func() error, err error) {
	msgs = make(chan *nats.Msg)

	var sub *nats.Subscription
	if sub, err = p.conn.ChanSubscribe(topic, msgs); err != nil {
		return
	}

	// close over the subscription and underlying receive channel, to
	// allow us to tear everything down later on.
	stop = func() (err error) {
		if err = sub.Unsubscribe(); err != nil {
			return
		}

		close(msgs)
		return
	}

	return
}

// GatherResponse allows an Agent to respond to a scatter-gather
// request from the Orchestrator, to allow it to decide which
// Agents for any given application should perform a kill operation.
func (p *natsProcessor) GatherResponse(orchInbox string, agentInbox string, application string) (err error) {
	return p.conn.PublishRequest(orchInbox, agentInbox, []byte(application))
}
