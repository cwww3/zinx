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
	fmt.Println("PreHandle")
}

func (b *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Handle")
	_, err := request.GetConnection().GetTCPConnection().Write(request.GetData())
	if err != nil {
		fmt.Println("handle err", err)
	}

}

func (b *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("PostHandle")
}

func main() {
	s := znet.NewServer()
	s.AddRouter(new(PingRouter))
	s.Serve()
}
