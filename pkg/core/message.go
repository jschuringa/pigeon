package core

type Message struct {
	Key     string `json:"key"`
	Content []byte `json:"content"`
}
