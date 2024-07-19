package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/jschuringa/pigeon/internal/connection"
	"github.com/jschuringa/pigeon/internal/core"
	"github.com/jschuringa/pigeon/pkg/publisher"

	"github.com/brianvoe/gofakeit"
)

func main() {
	pool := connection.NewTCPPool("localhost", 9090, 25, 50)
	ctx := context.Background()
	p := publisher.NewPublisher(ctx, pool)
	i := 1
	for i <= 1000 {
		time.Sleep(10000 * time.Duration(rand.Intn(10)))
		topic := "topic1"
		if i%2 == 0 {
			topic = "topic2"
		}

		bm := &core.BaseModel{
			Val: gofakeit.Name(),
		}

		if err := p.Publish(topic, bm); err != nil {
			fmt.Printf("error: %s", err)
		}
		i++
	}
}
