package znet

type Message struct {
	Id     uint32
	Length uint32
	Data   []byte
}

func NewMessage(msgId uint32, data []byte) *Message {
	return &Message{
		Id:     msgId,
		Length: uint32(len(data)),
		Data:   data,
	}
}

func (m *Message) GetMsgID() uint32 {
	return m.Id
}

func (m *Message) SetMsgID(id uint32) {
	m.Id = id
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsgData(data []byte) {
	m.Data = data
}

func (m *Message) GetLength() uint32 {
	return m.Length
}

func (m *Message) SetLength(length uint32) {
	m.Length = length
}
