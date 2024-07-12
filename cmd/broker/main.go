package main

import (
	"context"
	"os"

	"github.com/jschuringa/pigeon/pkg/broker"
)

func main() {
	b := broker.New()
	err := b.Start(context.Background())
	if err != nil {
		os.Exit(1)
	}
}
