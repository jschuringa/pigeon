package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jschuringa/pigeon/internal/core"
	"github.com/jschuringa/pigeon/pkg/subscriber"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	cfg := &subscriber.Config{
		Host: "localhost",
		Port: 8080,
	}

	topic1Sub := subscriber.NewSubscriber(cfg, "Topic 1")
	topic2Sub := subscriber.NewSubscriber(cfg, "Topic 2")
	go func() {
		err := topic1Sub.Subscribe(ctx, "topic1", printMessageTopicOne)
		if err != nil {
			// err chan + retry?
			return
		}
	}()

	go func() {
		err := topic2Sub.Subscribe(ctx, "topic2", printMessageTopicTwo)
		if err != nil {
			return
		}
	}()

	<-interrupt
	cancel()
}

func printMessageTopicOne(data []byte) error {
	var bm *core.BaseModel
	err := json.Unmarshal(data, &bm)
	if err != nil {
		return err
	}
	log.Printf("Received on t1: %s", bm.Val)
	return nil
}

func printMessageTopicTwo(data []byte) error {
	var bm *core.BaseModel
	err := json.Unmarshal(data, &bm)
	if err != nil {
		return err
	}
	log.Printf("Received on t2: %s", bm.Val)
	return nil
}
