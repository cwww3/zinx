package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (b *PingRouter) PreHandle(request ziface.IRequest) {
}

func (b *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Handle")
	err := request.GetConnection().SendMsg(1, []byte("Hello World!"))
	if err != nil {
		fmt.Println("handle err", err)
	}
}

func (b *PingRouter) PostHandle(request ziface.IRequest) {
}

func main() {
	s := znet.NewServer()
	s.AddRouter(new(PingRouter))
	s.Serve()
}
