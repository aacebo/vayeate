package client

import (
	"bufio"
	"encoding/binary"
)

type Code uint8

const (
	CONNECT         Code = 1
	CONNECT_ACK     Code = 2
	PUBLISH         Code = 3
	PUBLISH_ACK     Code = 4
	PUBLISH_REC     Code = 5
	PUBLISH_REL     Code = 6
	PUBLISH_COMP    Code = 7
	SUBSCRIBE       Code = 8
	SUBSCRIBE_ACK   Code = 9
	UNSUBSCRIBE     Code = 10
	UNSUBSCRIBE_ACK Code = 11
	PING            Code = 12
	PING_ACK        Code = 13
	DISCONNECT      Code = 14
)

type Message struct {
	Code    Code
	Payload []byte
}

func ReadMessage(reader *bufio.Reader) (*Message, error) {
	var code Code
	err := binary.Read(reader, binary.BigEndian, &code)

	if err != nil {
		return nil, err
	}

	var length uint32
	err = binary.Read(reader, binary.BigEndian, &length)

	if err != nil {
		return nil, err
	}

	payload := make([]byte, length)
	_, err = reader.Read(payload)

	if err != nil {
		return nil, err
	}

	return &Message{code, payload}, nil
}

func (self *Message) Serialize() []byte {
	b := []byte{byte(self.Code)}
	length := make([]byte, 4)
	binary.BigEndian.PutUint32(length, uint32(len(self.Payload)))
	b = append(b, length...)
	b = append(b, self.Payload...)
	return b
}
