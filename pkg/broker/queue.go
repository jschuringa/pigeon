package broker

import (
	"context"
	"fmt"

	"github.com/jschuringa/pigeon/pkg/core"

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

func (q *Queue) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-q.queue:
			if !ok {
				fmt.Printf("Channel for queue %s closed unexpectedly", q.name)
			}
			q.conn.WriteJSON(msg)
		}
	}
}
