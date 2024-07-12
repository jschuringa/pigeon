package main

import (
	"context"
	"os"

	"github.com/jschuringa/pigeon/pkg/broker"
)

func main() {
	b := broker.New()
	broker.RegisterTopic(b, "Topic 1", "topic1")
	broker.RegisterTopic(b, "Topic 2", "topic2")
	err := b.Start(context.Background())
	if err != nil {
		os.Exit(1)
	}
}
