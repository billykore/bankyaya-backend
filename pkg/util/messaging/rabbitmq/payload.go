package rabbitmq

import "encoding/json"

// MessagePayload is message broker payload.
type MessagePayload[T any] struct {
	// Origin is where the payload come from.
	Origin string `json:"origin"`
	// Data to publish or consume from the message broker.
	Data T `json:"data"`
}

// MarshalBinary return JSON encoding of MessagePayload.
func (p *MessagePayload[T]) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

// UnmarshalBinary parses encoded data into MessagePayload.
func (p *MessagePayload[T]) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}
