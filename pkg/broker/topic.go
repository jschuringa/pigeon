package broker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jschuringa/pigeon/pkg/core"
)

type Topic struct {
	Name  string
	Key   string
	queue chan core.Message
}

func RegisterTopic(b *Broker, name, key string) error {
	t := &Topic{
		Name:  name,
		Key:   key,
		queue: make(chan core.Message),
	}

	b.topics[key] = t
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
		}
	}
}

func decode(msg core.Message) (string, error) {
	bm := &core.BaseModel{}
	err := json.Unmarshal(msg.Content, &bm)
	if err != nil {
		return "", err
	}
	return bm.Val, nil
}
