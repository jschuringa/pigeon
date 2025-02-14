package subscriber

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jschuringa/pigeon/internal/core"

	"github.com/gorilla/websocket"
)

type Subscriber struct {
	// goroutine manager
	host string
	port int
	name string
}

// this feels so unnecessary right now lol
type Config struct {
	Host string
	Port int
}

func NewSubscriber(cfg *Config, name string) *Subscriber {
	return &Subscriber{
		host: cfg.Host,
		port: cfg.Port,
		name: name,
	}
}

// I think this should probably be made private and called on NewSubscriber. Put the handler and the topic in the new definition
// I could also use consumer groups to start all consumers at once with a cg.Start() and a private method?
// I dunno yet
func (s *Subscriber) Subscribe(ctx context.Context, topic string, handler func([]byte) error) error {
	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%d", s.host, s.port), Path: "/subscribe"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	defer c.Close()

	// Send message to connect to topics queue
	msg := &core.SubscriberRequest{
		Key:  topic,
		Name: s.name,
	}
	err = c.WriteJSON(&msg)
	if err != nil {
		return err
	}

	// read ack registration here eventually

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, barr, err := c.ReadMessage()
				if err != nil {
					return
				}
				handler(barr)
			}
		}
	}()
	<-ctx.Done()
	return nil
}
