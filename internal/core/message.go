package core

import "encoding/json"

type Message struct {
	Key string `json:"key"`
	// eventually msg id, other message metadata

	// this should be encrypted (eventually)
	Body []byte `json:"body"`
}

func Content[T any](m *Message) (*T, error) {
	var c *T
	err := json.Unmarshal(m.Body, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func NewMessage(key string, body any) (*Message, error) {
	barr, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return &Message{
		Body: barr,
		Key:  key,
	}, nil
}
