package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/jschuringa/pigeon/internal/core"

	"github.com/gorilla/websocket"
)

type Router struct {
	Name     string
	Key      string
	inbound  chan core.Message
	queues   []*Queue
	queueMtx sync.Mutex
}

func NewRouter(name, key string) *Router {
	return &Router{
		Name:    name,
		Key:     key,
		inbound: make(chan core.Message),
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
			res, err := decode(msg)
			if err != nil {
				return err
			}
			fmt.Printf("Message received on topic %s: %s\n", t.Name, res)
			for _, q := range t.queues {
				q.Push(res)
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

func decode(msg core.Message) (*core.BaseModel, error) {
	bm := &core.BaseModel{}
	err := json.Unmarshal(msg.Body, &bm)
	if err != nil {
		return nil, err
	}
	return bm, nil
}
