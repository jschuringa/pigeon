package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/jschuringa/pigeon/internal/core"
	"github.com/jschuringa/pigeon/internal/listener"

	"github.com/gorilla/websocket"
)

type Broker struct {
	// we should only write each topic once to this
	routes  sync.Map
	queue   chan core.Message
	errChan chan error
}

// Right now, running as a websockets server.
// Raw TCP proved very cumbersome
// Long term, move off websockets to a better implemented TPC and IPC solution
func (b *Broker) Server(ctx context.Context, host string, port int) {
	srv := &http.Server{Addr: fmt.Sprintf("%s:%d", host, port)}
	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			// just return for now until sending errors figured out
			fmt.Printf("failed to start websocket")
			return
		}

		// expect subscriber id, and last record id from connection
		// but for now, ask for topic, get all messages that show up

		defer conn.Close()
		b.handleSubscribe(ctx, conn)
	})

	// http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
	// 	upgrader := websocket.Upgrader{}
	// 	conn, err := upgrader.Upgrade(w, r, nil)
	// 	if err != nil {
	// 		fmt.Printf("failed to start websocket")
	// 	}

	// 	// expect subscriber id, and last record id from connection
	// 	// but for now, ask for topic, get all messages that show up

	// 	defer conn.Close()
	// 	b.handlePublish(ctx, conn)
	// })

	go func(errChan chan error) {
		if err := srv.ListenAndServe(); err != nil {
			errChan <- err
		}
	}(b.errChan)
	<-ctx.Done()
	srv.Close()
}

func (b *Broker) handleSubscribe(ctx context.Context, conn *websocket.Conn) error {
	// first we need to read the topic from the request
	var req *core.SubscriberRequest
	err := conn.ReadJSON(&req)
	if err != nil {
		return err
	}

	// then ensure the topic exists by getting the registered router
	t, err := b.getRouter(req.Key)
	if err != nil {
		return err
	}

	// then start subscribing to the topic
	return t.Subscribe(ctx, req.Name, conn)
}

// todo: if need to refactor TCP, can use websockets temporarily to communicate
// func (b *Broker) handlePublish(ctx context.Context, conn *websocket.Conn) error {
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return nil
// 		default:
// 			var msg *core.Message
// 			err := conn.ReadJSON(&msg)
// 			if err != nil {
// 				return err
// 			}
// 		case msg, ok := <-q.queue:
// 			if !ok {
// 				// this would mean a dropped message, so need some retry logic
// 				return fmt.Errorf("Channel for queue %s closed unexpectedly", q.name)
// 			}
// 			err := q.conn.WriteJSON(msg)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// }

func (b *Broker) RegisterRouter(name, key string) {
	rtr := NewRouter(name, key)
	b.routes.Store(key, rtr)
}

func (b *Broker) getRouter(key string) (*Router, error) {
	v, ok := b.routes.Load(key)
	if !ok {
		return nil, fmt.Errorf("topic does not exist")
	}

	t, ok := v.(*Router)
	if !ok {
		// better error message - probably could add a wrapper around the sync map (Topics?)
		// so that we're not worried about the type being wrong ever
		return nil, fmt.Errorf("value was not topic")
	}

	return t, nil
}

func (b *Broker) Receive(ctx context.Context) error {
	b.queue = make(chan core.Message)
	b.routes.Range(func(_, v any) bool {
		// currently can only register topics at startup
		t, ok := v.(*Router)
		if !ok {
			b.errChan <- fmt.Errorf("couldn't range over topic")
		}
		go func(errChan chan error) {
			if err := t.Listen(ctx); err != nil {
				errChan <- err
			}
		}(b.errChan)
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
			fmt.Printf("Received message with key: %s\n", msg.Key)
			rtr, err := b.getRouter(msg.Key)
			if err != nil {
				return err
			}
			rtr.inbound <- msg
		}
	}
}

func (b *Broker) Handle(c net.Conn) {
	defer c.Close()
	msg, err := readMessage(c)
	if err != nil {
		b.errChan <- err
		return
	}
	b.queue <- *msg
}

func (b *Broker) Start(ctx context.Context) error {
	go func(errChan chan error) {
		if err := b.Receive(ctx); err != nil {
			errChan <- err
		}
	}(b.errChan)
	go func(errChan chan error) {
		b.Server(ctx, "localhost", 8080)
	}(b.errChan)

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
		case err, ok := <-b.errChan:
			if !ok {
				return fmt.Errorf("broker err chan closed unexpectedly")
			}
			return err
		default:
			c, err := srv.Accept()
			if err != nil {
				return err
			}
			go b.Handle(c)
		}
	}
}

func New() *Broker {
	return &Broker{
		routes: sync.Map{},
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
