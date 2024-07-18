package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/jschuringa/pigeon/internal/listener"
	"github.com/jschuringa/pigeon/pkg/core"

	"github.com/gorilla/websocket"
)

type Broker struct {
	// we should only write each topic once to this
	topics   sync.Map
	messages []string
	queue    chan core.Message
}

func (b *Broker) SubscriberServer(ctx context.Context) {
	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Printf("failed to start websocket")
		}

		// expect subscriber id, and last record id from connection
		// but for now, ask for topic, get all messages stored from that topic

		defer conn.Close()

		// first we need to read the topic from the request
		var req *core.SubscriberRequest
		err = conn.ReadJSON(&req)
		if err != nil {
			fmt.Printf("failed to read subscriber request: %s", err)
			return
		}

		// can probably add a helper function that wraps this all so don't need to constantly do it
		t, err := b.getTopic(req.Key)
		if err != nil {
			fmt.Printf("failed to get topic: %s", err)
			return
		}

		t.Subscribe(ctx, req.Name, conn)
	})

	http.ListenAndServe("localhost:8080", nil)
}

func (b *Broker) getTopic(key string) (*Topic, error) {
	v, ok := b.topics.Load(key)
	if !ok {
		return nil, fmt.Errorf("topic does not exist")
	}

	t, ok := v.(*Topic)
	if !ok {
		// better error message - probably could add a wrapper around the sync map (Topics?)
		// so that we're not worried about the type being wrong ever
		return nil, fmt.Errorf("value was not topic")
	}

	return t, nil
}

func (b *Broker) Receive(ctx context.Context) error {
	b.queue = make(chan core.Message)
	b.topics.Range(func(_, v any) bool {
		// currently only subscribers registered through startup would work
		// can expose api eventually to add during runtime
		t, ok := v.(*Topic)
		if !ok {
			fmt.Printf("range over topic failed")
		}
		go t.Listen(ctx)
		return true
	})
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-b.queue:
			if !ok {
				return fmt.Errorf("channel closed unexpectedly")
			}
			fmt.Printf("Received message with key: %s\n\n", msg.Key)
			if v, ok := b.topics.Load(msg.Key); ok {
				t, ok := v.(*Topic)
				if !ok {
					fmt.Printf("failed to load topic")
				}
				t.queue <- msg
			}
		}
	}
}

func (b *Broker) Handle(c net.Conn) {
	defer c.Close()
	msg, err := readMessage(c)
	if err != nil {
		fmt.Print("we dropped a message :(\n")
		return
	}
	b.queue <- *msg
}

func (b *Broker) Start(ctx context.Context) error {
	go b.Receive(ctx)
	go b.SubscriberServer(ctx)
	l := listener.NewListener("localhost", 9090)
	srv, err := l.Start(ctx)
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
			go b.Handle(c)
		}
	}
}

func New() *Broker {
	return &Broker{
		messages: make([]string, 0),
		topics:   sync.Map{},
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
