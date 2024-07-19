package broker

import (
	"context"
	"fmt"

	"github.com/jschuringa/pigeon/internal/core"

	"github.com/gorilla/websocket"
)

// Queue allows us to build a queue for each
// subscriber that connects
type Queue struct {
	name  string
	conn  *websocket.Conn
	queue chan *core.BaseModel
}

func NewQueue(name string, conn *websocket.Conn) *Queue {
	return &Queue{
		name:  name,
		conn:  conn,
		queue: make(chan *core.BaseModel),
	}
}

func (q *Queue) Push(msg *core.BaseModel) {
	q.queue <- msg
}

func (q *Queue) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-q.queue:
			if !ok {
				// this would mean a dropped message, so need some retry logic
				return fmt.Errorf("Channel for queue %s closed unexpectedly", q.name)
			}
			err := q.conn.WriteJSON(msg)
			if err != nil {
				return err
			}
		}
	}
}
