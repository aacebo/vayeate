package frame

import (
	"bytes"
	"errors"
	"strconv"
)

type OpCode uint8

const (
	CLOSE  OpCode = 0 // <code>
	PING   OpCode = 1 // <code>
	PONG   OpCode = 2 // <code>
	PUB    OpCode = 3 // <code:queue:body>
	SUB    OpCode = 4 // <code:queue>
	ASSERT OpCode = 5 // <code:body>
)

const (
	START     = byte('<')
	END       = byte('>')
	DELIMITER = byte(':')
)

var InvalidFormatError = errors.New("invalid frame format")
var OpCodeLength = map[OpCode]int{
	CLOSE:  1,
	PING:   1,
	PONG:   1,
	PUB:    3,
	SUB:    2,
	ASSERT: 2,
}

type Frame struct {
	Code OpCode
	Body []byte
}

func New(code OpCode, body []byte) *Frame {
	self := Frame{code, body}
	return &self
}

func NewPing() *Frame {
	self := Frame{PING, []byte{}}
	return &self
}

func NewPong() *Frame {
	self := Frame{PONG, []byte{}}
	return &self
}

func NewClose() *Frame {
	self := Frame{CLOSE, []byte{}}
	return &self
}

func Decode(data []byte) (*Frame, error) {
	body := []byte{}
	slices := bytes.Split(data, []byte{DELIMITER})

	if len(slices) < 1 || len(slices) > 3 {
		return nil, InvalidFormatError
	}

	t, err := strconv.Atoi(string(slices[0]))

	if err != nil {
		return nil, err
	}

	if len(slices) == 2 {
		body = slices[1]
	}

	code := OpCode(t)

	if OpCodeLength[code] != len(slices) {
		return nil, InvalidFormatError
	}

	self := Frame{code, body}
	return &self, nil
}

func (self *Frame) Encode() []byte {
	data := []byte{}
	code := []byte(strconv.Itoa(int(self.Code)))

	data = append(data, START)
	data = append(data, code...)
	data = append(data, END)

	return data
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

func (self *Frame) IsPublish() bool {
	return self.Code == PUB
}

func (self *Frame) IsSubscribe() bool {
	return self.Code == SUB
}

func (self *Frame) IsAssert() bool {
	return self.Code == ASSERT
}
