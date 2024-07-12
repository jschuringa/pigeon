package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/brianvoe/gofakeit"
	"github.com/jschuringa/pigeon/pkg/core"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(3)
	go publishALotOfMessages(&wg)
	go publishALotOfMessages(&wg)
	go publishALotOfMessages(&wg)
	wg.Wait()
}

func publishALotOfMessages(wg *sync.WaitGroup) {
	defer wg.Done()
	i := 0
	for i < 1000000 {
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

		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			println("Dial failed:", err.Error())
			os.Exit(1)
		}
		defer conn.Close()

		enc := json.NewEncoder(conn)
		err = enc.Encode(&msg)
		if err != nil {
			fmt.Printf("Err: %v", err)
			os.Exit(1)
		}
		i++
	}
}
