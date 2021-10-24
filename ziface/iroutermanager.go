package ziface

type IRouterManager interface {
	DoMsgHandler(request IRequest)
	AddRouter(msgId uint32, router IRouter)
}
