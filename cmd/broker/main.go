package main

import (
	"context"
	"os"

	"github.com/jschuringa/pigeon/pkg/broker"
)

func main() {
	b := broker.New(&broker.Config{
		WSHost:  "localhost",
		WSPort:  8080,
		TCPHost: "localhost",
		TCPPort: 9090,
	})
	b.RegisterRouter("Topic 1", "topic1")
	b.RegisterRouter("Topic 2", "topic2")
	err := b.Start(context.Background())
	if err != nil {
		os.Exit(1)
	}
}
