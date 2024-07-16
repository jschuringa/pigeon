package core

import "encoding/json"

type Message struct {
	Key     string `json:"key"`
	Content []byte `json:"content"`
}

func GetContent[T any](m *Message) (*T, error) {
	var c *T
	err := json.Unmarshal(m.Content, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func NewMessage(key string, content any) (*Message, error) {
	barr, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}

	return &Message{
		Content: barr,
		Key:     key,
	}, nil
}
