package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jschuringa/pigeon/internal/core"

	"github.com/gorilla/websocket"
)

type Router struct {
	Name     string
	Key      string
	inbound  chan *core.Message
	queues   []*Queue
	queueMtx sync.Mutex
}

func NewRouter(name, key string) *Router {
	return &Router{
		Name:    name,
		Key:     key,
		inbound: make(chan *core.Message),
	}
}

func (t *Router) Listen(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-t.inbound:
			if !ok {
				return fmt.Errorf("topic %s channel closed unexpectedly", t.Name)
			}
			log.Printf("Message received on topic %s", t.Name)
			t.writeToFile(msg)
			for _, q := range t.queues {
				q.Push(msg)
			}
		}
	}
}

func (t *Router) Subscribe(ctx context.Context, name string, conn *websocket.Conn) error {
	q := NewQueue(name, conn)
	t.queueMtx.Lock()
	t.queues = append(t.queues, q)
	t.queueMtx.Unlock()
	return q.Start(ctx)
}

func (t *Router) writeToFile(msg *core.Message) error {
	barr, err := json.Marshal(&msg)
	if err != nil {
		return err
	}
	err = os.MkdirAll(fmt.Sprintf("./msgs/%s", t.Key), os.ModePerm)
	if err != nil {
		return err
	}
	err = os.WriteFile(fmt.Sprintf("./msgs/%s/%d.json", t.Key, time.Now().Unix()), barr, 0644)
	if err != nil {
		return err
	}
	return nil
}
