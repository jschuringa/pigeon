package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"github.com/jschuringa/pigeon/pkg/core"
)

type Broker struct {
	topics   map[string]*Topic
	messages []string
	queue    chan core.Message
	// sendTo   string
}

func (b *Broker) Receive(ctx context.Context) error {
	b.queue = make(chan core.Message)
	for _, v := range b.topics {
		go v.Listen(ctx)
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-b.queue:
			if !ok {
				return fmt.Errorf("channel closed unexpectedly")
			}
			fmt.Printf("Received message with key: %s", msg.Key)
			if t, ok := b.topics[msg.Key]; ok {
				t.queue <- msg
			}
		}

	}
}

func (b *Broker) Flush() error {
	return fmt.Errorf("not implemented")
}

func (b *Broker) Push(c net.Conn) {
	msg, err := readMessage(c)
	if err != nil {
		fmt.Print("we dropped a message :(\n")
		return
	}
	b.queue <- *msg
}

func (b *Broker) Start(ctx context.Context) error {
	go b.Receive(ctx)
	srv, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		return err
	}
	defer srv.Close()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			fmt.Printf("waiting for accept\n")
			c, err := srv.Accept()
			fmt.Printf("accepted\n")
			if err != nil {
				return err
			}
			go b.Push(c)
		}
	}
}

func New() *Broker {
	return &Broker{
		messages: make([]string, 0),
		topics:   make(map[string]*Topic, 0),
	}
}

func readMessage(c net.Conn) (*core.Message, error) {
	msg := &core.Message{}
	dec := json.NewDecoder(c)
	err := dec.Decode(&msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
