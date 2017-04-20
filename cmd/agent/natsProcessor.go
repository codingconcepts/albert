package main

import nats "github.com/nats-io/go-nats"

type natsProcessor struct {
	conn *nats.Conn
}

func newNatsProcessor(conn *nats.Conn) (p *natsProcessor) {
	return &natsProcessor{
		conn: conn,
	}
}

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

func (p *natsProcessor) subscribe(topic string) (msgs chan *nats.Msg, stop func() error, err error) {
	msgs = make(chan *nats.Msg)

	var sub *nats.Subscription
	if sub, err = p.conn.ChanSubscribe(topic, msgs); err != nil {
		return
	}

	stop = func() (err error) {
		if err = sub.Unsubscribe(); err != nil {
			return
		}

		close(msgs)
		return
	}

	return
}

func (p *natsProcessor) GatherResponse(orchInbox string, agentInbox string, application string) (err error) {
	return p.conn.PublishRequest(orchInbox, agentInbox, []byte(application))
}
