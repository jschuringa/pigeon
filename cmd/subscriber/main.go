package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/jschuringa/pigeon/pkg/core"

	"github.com/gorilla/websocket"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/subscribe"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	msg := &core.SubscriberRequest{
		Key:  "topic1",
		Name: "test-subscriber",
	}
	c.WriteJSON(&msg)

	go func() {
		defer close(done)
		for {
			bm := &core.BaseModel{}
			err := c.ReadJSON(&bm)
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", bm.Val)
		}
	}()
	select {}
}
