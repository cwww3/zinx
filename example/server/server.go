package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

type Ping2Router struct {
	znet.BaseRouter
}

func (b *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("ping Handle")
	err := request.GetConnection().SendMsg(request.GetMsgID(), request.GetData())
	if err != nil {
		fmt.Println("handle err", err)
	}
}

func (b *Ping2Router) Handle(request ziface.IRequest) {
	fmt.Println("ping2 Handle")
	err := request.GetConnection().SendMsg(request.GetMsgID(), request.GetData())
	if err != nil {
		fmt.Println("handle err", err)
	}
}

func main() {
	s := znet.NewServer()
	s.AddRouter(1, new(PingRouter))
	s.AddRouter(2, new(Ping2Router))
	s.Serve()
}
