package client

func NewPingAckMessage() *Message {
	return &Message{
		Code:    PING_ACK,
		Payload: []byte{},
	}
}
