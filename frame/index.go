package frame

type Frame struct {
	Type []byte
	Data []byte
}

func New(t string, data []byte) *Frame {
	self := Frame{[]byte(t), data}
	return &self
}

func Decode(data []byte) *Frame {
	self := Frame{data[0:4],data[4:]}
	return &self
}

func (self *Frame) Encode() []byte {
	return append(
		self.Type,
		self.Data...
	)
}

func (self *Frame) GetType() string {
	return string(self.Type)
}

func (self *Frame) GetData() string {
	return string(self.Data)
}
