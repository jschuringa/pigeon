package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
)

type Broker struct {
	messages []string
	queue    chan string
	// sendTo   string
}

type message struct {
	Val string `json:"val"`
}

func (b *Broker) Receive(ctx context.Context) error {
	b.queue = make(chan string)
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-b.queue:
			if !ok {
				return fmt.Errorf("channel closed unexpectedly")
			}
			fmt.Printf("Received message: %s", msg)
			b.messages = append(b.messages, msg)
		}

	}
}

func (b *Broker) Flush() error {
	return fmt.Errorf("not implemented")
}

func (b *Broker) Push(c net.Conn) {
	msg, err := readMessage(c)
	if err != nil {
		fmt.Print("we dropped a message :(")
		return
	}
	b.queue <- msg.Val
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
			fmt.Printf("waiting for accept")
			c, err := srv.Accept()
			fmt.Printf("accepted")
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
	}
}

func readMessage(c net.Conn) (*message, error) {
	// l := make([]byte, 2) // takes the first two bytes: the Length Message.
	// i, _ := c.Read(l)    // read the bytes.

	// lm := binary.BigEndian.Uint16(l[:i]) // convert the bytes into int16.

	// b := make([]byte, lm) // create the byte buffer for the body.
	// e, _ := c.Read(b)     // read the bytes.
	b := make([]byte, 0)
	c.Read(b)

	msg := &message{}
	err := json.Unmarshal(b, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil // returns the body
}
