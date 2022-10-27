package mqttc

import (
	"encoding/json"
	"errors"
	"time"
)

type Message[T any] struct {
	Command string `json:"command"`
	Data    T      `json:"data"`
}

func NewMessage[T any](command string, data T) *Message[T] {
	return &Message[T]{
		Command: command,
		Data:    data,
	}
}

type BootstrapMessage struct {
	ClientId  string `json:"client_id"`
	Timestamp int64  `json:"timestamp"`
}

func newBootstrapMessage(clientId string) *Message[BootstrapMessage] {
	return NewMessage[BootstrapMessage]("bootstrap", BootstrapMessage{
		ClientId:  clientId,
		Timestamp: time.Now().Unix(),
	})
}

type InitializeMessage struct {
	ServerId  string `json:"server_id"`
	Timestamp int64  `json:"timestamp"`
}

func (m *Message[T]) Encode() ([]byte, error) {
	if m.Command == "" {
		return nil, errors.New("command is empty")
	}
	return json.MarshalIndent(m, "", "  ")
}

func DecodeMessage[T any](data []byte) (*Message[T], error) {
	var m Message[T]
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
