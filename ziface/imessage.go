package ziface

type IMessage interface {
	GetMsgID() uint32
	SetMsgID(id uint32)

	GetData() []byte
	SetMsgData(data []byte)

	GetLength() uint32
	SetLength(length uint32)
}
