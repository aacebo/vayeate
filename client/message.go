package client

import (
	"bufio"
	"encoding/binary"
	"io"
)

type Code uint8

const (
	ERROR           Code = 0
	CONNECT         Code = 1
	CONNECT_ACK     Code = 2
	PUBLISH         Code = 3
	PUBLISH_ACK     Code = 4
	CONSUME         Code = 5
	CONSUME_ACK     Code = 6
	SUBSCRIBE       Code = 7
	SUBSCRIBE_ACK   Code = 8
	UNSUBSCRIBE     Code = 9
	UNSUBSCRIBE_ACK Code = 10
	PING            Code = 11
	PING_ACK        Code = 12
)

type Message struct {
	Code    Code
	SentAt  int64
	Payload []byte
}

func ReadMessage(reader *bufio.Reader) (*Message, error) {
	var code Code
	err := binary.Read(reader, binary.BigEndian, &code)

	if err != nil {
		return nil, err
	}

	var sentAt uint64
	err = binary.Read(reader, binary.BigEndian, &sentAt)

	if err != nil {
		return nil, err
	}

	var length uint32
	err = binary.Read(reader, binary.BigEndian, &length)

	if err != nil {
		return nil, err
	}

	payload := make([]byte, length)
	_, err = io.ReadFull(reader, payload)

	if err != nil {
		return nil, err
	}

	return &Message{code, int64(sentAt), payload}, nil
}

func (self *Message) Serialize() []byte {
	b := []byte{byte(self.Code)}

	sentAt := make([]byte, 8)
	binary.BigEndian.PutUint64(sentAt, uint64(self.SentAt))

	length := make([]byte, 4)
	binary.BigEndian.PutUint32(length, uint32(len(self.Payload)))

	b = append(b, sentAt...)
	b = append(b, length...)
	b = append(b, self.Payload...)

	return b
}
