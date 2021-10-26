package znet

import (
	"github.com/cwww3/zinx/ziface"
)

type BaseRouter struct {
}

func (b *BaseRouter) PreHandle(request ziface.IRequest) {
}

func (b *BaseRouter) Handle(request ziface.IRequest) {
}

func (b *BaseRouter) PostHandle(request ziface.IRequest) {
}
