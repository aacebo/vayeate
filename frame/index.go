package frame

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

type OpCode uint8

const (
	CLOSE    OpCode = 0 // <code::>
	PING     OpCode = 1 // <code::>
	PONG     OpCode = 2 // <code::>
	ASSERT   OpCode = 3 // <code:subject:>
	PRODUCE  OpCode = 4 // <code:subject:body>
	CONSUME  OpCode = 5 // <code:subject:>
	ACK      OpCode = 6 // <code:subject:>
	DELEGATE OpCode = 7 // <code:subject:body>
)

const (
	START     = byte('<') // frame start byte
	END       = byte('>') // frame end byte
	DELIMITER = byte(':') // frame slice delimiter
)

var InvalidFormatError = errors.New("invalid frame format")

// A Frame comprised of many packets, utilizing a custom Netstring format
// ex. <code:subject:body>
// https://cr.yp.to/proto/netstrings.txt
type Frame struct {
	Code    OpCode
	Subject []byte
	Body    []byte
}

func New(code OpCode, subject []byte, body []byte) *Frame {
	return &Frame{code, subject, body}
}

func NewClose() *Frame {
	return &Frame{CLOSE, []byte{}, []byte{}}
}

func NewPing() *Frame {
	return &Frame{PING, []byte{}, []byte{}}
}

func NewPong() *Frame {
	return &Frame{PONG, []byte{}, []byte{}}
}

func NewAssert(subject []byte) *Frame {
	return &Frame{ASSERT, subject, []byte{}}
}

func NewProduce(subject []byte, body []byte) *Frame {
	return &Frame{PRODUCE, subject, body}
}

func NewConsume(subject []byte) *Frame {
	return &Frame{CONSUME, subject, []byte{}}
}

func NewAck(subject []byte) *Frame {
	return &Frame{ACK, subject, []byte{}}
}

func NewDelegate(subject []byte, body []byte) *Frame {
	return &Frame{DELEGATE, subject, body}
}

func Decode(reader *bufio.Reader) (*Frame, error) {
	subject := []byte{}
	body := []byte{}
	b, err := reader.ReadByte()

	if err != nil {
		return nil, err
	}

	if b != START {
		return nil, InvalidFormatError
	}

	// read opcode
	b, err = reader.ReadByte()

	if err != nil {
		return nil, err
	}

	t, err := strconv.Atoi(string(b))

	if err != nil {
		return nil, err
	}

	code := OpCode(t)

	b, err = reader.ReadByte()

	if err != nil {
		return nil, err
	}

	if b != DELIMITER {
		return nil, InvalidFormatError
	}

	// read subject
	for {
		b, err := reader.ReadByte()

		if err == io.EOF {
			return nil, InvalidFormatError
		}

		if err != nil {
			return nil, err
		}

		if b == DELIMITER {
			break
		}

		subject = append(subject, b)
	}

	// read body
	for {
		b, err := reader.ReadByte()

		if err == io.EOF {
			return nil, InvalidFormatError
		}

		if err != nil {
			return nil, err
		}

		if b == END {
			break
		}

		body = append(body, b)
	}

	return &Frame{code, subject, body}, nil
}

func (self *Frame) Encode() []byte {
	data := []byte{}
	code := []byte(strconv.Itoa(int(self.Code)))

	data = append(data, START)
	data = append(data, code...)
	data = append(data, DELIMITER)
	data = append(data, self.Subject...)
	data = append(data, DELIMITER)
	data = append(data, self.Body...)
	data = append(data, END)

	return data
}

func (self *Frame) GetSubject() string {
	return string(self.Subject)
}

func (self *Frame) GetBody() string {
	return string(self.Body)
}

func (self *Frame) IsClose() bool {
	return self.Code == CLOSE
}

func (self *Frame) IsPing() bool {
	return self.Code == PING
}

func (self *Frame) IsPong() bool {
	return self.Code == PONG
}

func (self *Frame) IsAssert() bool {
	return self.Code == ASSERT
}

func (self *Frame) IsProduce() bool {
	return self.Code == PRODUCE
}

func (self *Frame) IsConsume() bool {
	return self.Code == CONSUME
}

func (self *Frame) IsAck() bool {
	return self.Code == ACK
}

func (self *Frame) IsDelegate() bool {
	return self.Code == DELEGATE
}
