package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/jschuringa/pigeon/pkg/core"

	"github.com/gorilla/websocket"
)

type Topic struct {
	Name     string
	Key      string
	queue    chan core.Message
	queues   []*Queue
	queueMtx sync.Mutex
}

func RegisterTopic(b *Broker, name, key string) error {
	t := &Topic{
		Name:  name,
		Key:   key,
		queue: make(chan core.Message),
	}

	b.topics.Store(key, t)
	return nil
}

func (t *Topic) Listen(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-t.queue:
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

func (t *Topic) Subscribe(ctx context.Context, name string, conn *websocket.Conn) error {
	q := NewQueue(name, conn)
	t.queueMtx.Lock()
	t.queues = append(t.queues, q)
	t.queueMtx.Unlock()
	q.Start(ctx)
	return nil
}

func decode(msg core.Message) (*core.BaseModel, error) {
	bm := &core.BaseModel{}
	err := json.Unmarshal(msg.Body, &bm)
	if err != nil {
		return nil, err
	}
	return bm, nil
}
