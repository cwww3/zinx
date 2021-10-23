package znet

import (
	"zinx/ziface"
)

type Request struct {
	connection ziface.IConnection
	data       []byte
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.connection
}

func (r *Request) GetData() []byte {
	return r.data
}
