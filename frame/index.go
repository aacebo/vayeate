package frame

import (
	"bytes"
	"errors"
	"strconv"
)

type FrameType uint8

const (
	CLOSE FrameType = 0
	PING  FrameType = 1
	PONG  FrameType = 2
	PUB   FrameType = 3
	SUB   FrameType = 4
)

const Delimiter = byte(':')

var InvalidFormatError = errors.New("invalid frame format")

type Frame struct {
	Type    FrameType
	Payload []byte
}

func New(t FrameType, payload []byte) *Frame {
	self := Frame{t, payload}
	return &self
}

func Ping() *Frame {
	self := Frame{PING, []byte{}}
	return &self
}

func Pong() *Frame {
	self := Frame{PONG, []byte{}}
	return &self
}

func Close() *Frame {
	self := Frame{CLOSE, []byte{}}
	return &self
}

func Decode(data []byte) (*Frame, error) {
	slices := bytes.Split(data, []byte{Delimiter})

	if len(slices) != 3 {
		return nil, InvalidFormatError
	}

	t, err := strconv.Atoi(string(slices[1]))

	if err != nil {
		return nil, err
	}

	self := Frame{FrameType(t), slices[2]}
	return &self, nil
}

func (self *Frame) Encode() []byte {
	len := []byte(strconv.Itoa(len(self.Payload)))
	t := []byte(strconv.Itoa(int(self.Type)))

	return append(
		append(append(len, Delimiter), append(t, Delimiter)...),
		self.Payload...,
	)
}

func (self *Frame) GetPayload() string {
	return string(self.Payload)
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
