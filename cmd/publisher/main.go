package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/jschuringa/pigeon/internal/connection"
	"github.com/jschuringa/pigeon/pkg/core"

	"github.com/brianvoe/gofakeit"
)

func main() {

	pool := connection.NewTCPPool("localhost", 8080, 25, 50)

	var wg sync.WaitGroup
	wg.Add(1)
	go publishALotOfMessages(&wg, pool)
	// go publishALotOfMessages(&wg, pool)
	// go publishALotOfMessages(&wg, pool)
	wg.Wait()
	os.Exit(0)
}

func publishALotOfMessages(wg *sync.WaitGroup, pool *connection.TCPPool) {
	defer wg.Done()
	i := 0
	for i < 10000 {
		time.Sleep(1000 * time.Duration(rand.Int()))
		bm := &core.BaseModel{
			Val: gofakeit.Name(),
		}

		key := "topic1"
		if i%2 == 0 {
			key = "topic2"
		}

		encoded, err := json.Marshal(bm)
		if err != nil {
			fmt.Printf("error: %v", err)
			os.Exit(1)
		}

		msg := &core.Message{
			Key:     key,
			Content: encoded,
		}

		conn, err := pool.Get()
		if err != nil {
			println("Dial failed:", err.Error())
			os.Exit(1)
		}

		encMsg, err := json.Marshal(&msg)
		if err != nil {
			fmt.Printf("Marshalling msg failed: %s", err)
			os.Exit(1)
		}

		err = conn.Send(encMsg)
		if err != nil {
			fmt.Printf("Err: %v", err)
			os.Exit(1)
		}
		i++
		pool.Put(conn)
	}
}
