package publisher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jschuringa/pigeon/internal/connection"
	"github.com/jschuringa/pigeon/internal/core"
)

type Publisher struct {
	// todo: interface TCPPool for testing
	pool      *connection.TCPPool
	sendQueue chan *core.Message
}

func (p *Publisher) Publish(topic string, content any) error {
	body, err := json.Marshal(&content)
	if err != nil {
		return err
	}
	p.sendQueue <- &core.Message{Key: topic, Body: body}
	return nil
}

func NewPublisher(ctx context.Context, pool *connection.TCPPool) *Publisher {
	p := &Publisher{
		pool:      pool,
		sendQueue: make(chan *core.Message),
	}

	// maybe I could build this into fx or something? curious how kafka libraries do it
	// goroutine pool for just the publisher logic?
	go p.start(ctx)
	return p
}

func (p *Publisher) start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-p.sendQueue:
			if !ok {
				// push err to err chan
				panic("this should not happen")
			}
			conn, err := p.pool.Get()
			if err != nil {
				println("Dial failed:", err.Error())
				// push err to err chan
			}

			encMsg, err := json.Marshal(&msg)
			if err != nil {
				fmt.Printf("Marshalling msg failed: %s", err)
				// push err to err chan
			}

			err = conn.Send(encMsg)
			if err != nil {
				fmt.Printf("Err: %v", err)
				// push err to err chan
			}
			p.pool.Close(conn)
		}
	}
}
