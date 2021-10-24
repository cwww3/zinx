package znet

import (
	"fmt"
	"zinx/ziface"
)

type RouterManager struct {
	Apis map[uint32]ziface.IRouter
}

func NewRouterManager() *RouterManager {
	return &RouterManager{Apis: make(map[uint32]ziface.IRouter)}
}

func (m *RouterManager) DoMsgHandler(request ziface.IRequest) {
	router, ok := m.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println(request.GetMsgID(), "router not found")
	}
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

func (m *RouterManager) AddRouter(msgId uint32, router ziface.IRouter) {
	if _, ok := m.Apis[msgId]; ok {
		panic(fmt.Sprintf("%s%d", "repeat add router", msgId))
	}
	m.Apis[msgId] = router
	fmt.Printf("%s%d\n", "add router", msgId)
}
