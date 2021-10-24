package znet

import (
	"fmt"
	"zinx/utils"
	"zinx/ziface"
)

var cnt int

type RouterManager struct {
	Apis         map[uint32]ziface.IRouter
	RequestQueue []chan ziface.IRequest
}

func NewRouterManager() *RouterManager {
	return &RouterManager{
		Apis: make(map[uint32]ziface.IRouter),
	}
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

func (m *RouterManager) StartWorkers() {
	m.RequestQueue = make([]chan ziface.IRequest, utils.Config.WorkerSize)
	for i := 0; i < utils.Config.WorkerSize; i++ {
		m.RequestQueue[i] = make(chan ziface.IRequest, utils.Config.RequestPoolSize)
		go func(requestCh chan ziface.IRequest) {
			for {
				req := <-requestCh
				m.DoMsgHandler(req)
			}
		}(m.RequestQueue[i])
	}
}

func (m *RouterManager) SendRequest(request ziface.IRequest) {
	cnt %= utils.Config.WorkerSize
	m.RequestQueue[cnt] <- request
	cnt++
}
