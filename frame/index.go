package frame

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

type OpCode uint8

const (
	CLOSE   OpCode = 0 // <code::>
	PING    OpCode = 1 // <code::>
	PONG    OpCode = 2 // <code::>
	ASSERT  OpCode = 3 // <code:subject:>
	PRODUCE OpCode = 4 // <code:subject:body>
	CONSUME OpCode = 5 // <code:subject:>
	ACK     OpCode = 6 // <code:subject:>
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
	self := Frame{code, subject, body}
	return &self
}

func NewClose() *Frame {
	self := Frame{CLOSE, []byte{}, []byte{}}
	return &self
}

func NewPing() *Frame {
	self := Frame{PING, []byte{}, []byte{}}
	return &self
}

func NewPong() *Frame {
	self := Frame{PONG, []byte{}, []byte{}}
	return &self
}

func NewAssert(subject []byte) *Frame {
	self := Frame{ASSERT, subject, []byte{}}
	return &self
}

func NewProduce(subject []byte, body []byte) *Frame {
	self := Frame{PRODUCE, subject, body}
	return &self
}

func NewConsume(subject []byte) *Frame {
	self := Frame{CONSUME, subject, []byte{}}
	return &self
}

func NewAck(subject []byte) *Frame {
	self := Frame{ACK, subject, []byte{}}
	return &self
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

	self := Frame{code, subject, body}
	return &self, nil
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
