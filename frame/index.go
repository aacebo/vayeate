package frame

import (
	"bytes"
	"errors"
	"strconv"
)

type FrameType uint8

const (
	CLOSE  FrameType = 0
	PING   FrameType = 1
	PONG   FrameType = 2
	PUB    FrameType = 3
	SUB    FrameType = 4
	ASSERT FrameType = 5
)

const Delimiter = byte(':')

var InvalidFormatError = errors.New("invalid frame format")

type Frame struct {
	Type FrameType
	Body []byte
}

func New(t FrameType, body []byte) *Frame {
	self := Frame{t, body}
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
	slices := bytes.Split(data, []byte{Delimiter})

	if len(slices) < 1 || len(slices) > 2 {
		return nil, InvalidFormatError
	}

	t, err := strconv.Atoi(string(slices[0]))

	if err != nil {
		return nil, err
	}

	if len(slices) == 2 {
		body = slices[1]
	}

	self := Frame{FrameType(t), body}
	return &self, nil
}

func (self *Frame) Encode() []byte {
	t := []byte(strconv.Itoa(int(self.Type)))

	return append(
		append(t, Delimiter),
		self.Body...,
	)
}

func (self *Frame) GetBody() string {
	return string(self.Body)
}

func (self *Frame) IsClose() bool {
	return self.Type == CLOSE
}

func (self *Frame) IsPing() bool {
	return self.Type == PING
}

func (self *Frame) IsPong() bool {
	return self.Type == PONG
}

func (self *Frame) IsPublish() bool {
	return self.Type == PUB
}

func (self *Frame) IsSubscribe() bool {
	return self.Type == SUB
}

func (self *Frame) IsAssert() bool {
	return self.Type == ASSERT
}
