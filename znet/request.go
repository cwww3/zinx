package znet

import (
	"zinx/ziface"
)

type Request struct {
	connection ziface.IConnection
	msg        ziface.IMessage
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.connection
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}
