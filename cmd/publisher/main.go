package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/jschuringa/pigeon/internal/connection"
	"github.com/jschuringa/pigeon/pkg/core"
	"github.com/jschuringa/pigeon/pkg/publisher"

	"github.com/brianvoe/gofakeit"
)

func main() {

	pool := connection.NewTCPPool("localhost", 9090, 25, 50)
	ctx := context.Background()
	p := publisher.NewPublisher(ctx, pool)
	i := 0
	for i < 1000 {
		time.Sleep(10000 * time.Duration(rand.Intn(10)))
		topic := "topic1"
		if i%rand.Intn(9) == 0 {
			topic = "topic2"
		}

		bm := &core.BaseModel{
			Val: gofakeit.Name(),
		}
		encoded, err := json.Marshal(bm)
		if err != nil {
			fmt.Printf("error: %s", err)
			os.Exit(1)
		}

		if err := p.Publish(topic, encoded); err != nil {
			fmt.Printf("error: %s", err)
		}
	}
}
