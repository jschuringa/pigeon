package main

import (
	"encoding/json"
	"net"
	"os"

	"github.com/brianvoe/gofakeit"
)

type Message struct {
	Val string `json:"val"`
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	i := 0
	for i < 1000 {
		msg := &Message{
			Val: gofakeit.Name(),
		}

		data, err := json.Marshal(msg)
		if err != nil {
			os.Exit(1)
		}
		_, err = conn.Write(data)
		if err != nil {
			os.Exit(1)
		}
		i++
	}
}
