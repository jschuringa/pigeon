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
	name     string
	conn     *websocket.Conn
	outbound chan *core.Message
}

func NewQueue(name string, conn *websocket.Conn) *Queue {
	return &Queue{
		name:     name,
		conn:     conn,
		outbound: make(chan *core.Message),
	}
}

func (q *Queue) Push(msg *core.Message) {
	q.outbound <- msg
}

func (q *Queue) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-q.outbound:
			if !ok {
				// this would mean a dropped message, so need some retry logic
				return fmt.Errorf("channel for queue %s closed unexpectedly", q.name)
			}
			err := q.conn.WriteJSON(msg)
			if err != nil {
				return err
			}
		}
	}
}
