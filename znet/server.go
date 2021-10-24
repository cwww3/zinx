package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	Name          string
	IPVersion     string
	IP            string
	Port          int
	RouterManager ziface.IRouterManager
}

func NewServer() ziface.IServer {
	s := Server{
		Name:          utils.Config.Name,
		IPVersion:     "tcp4",
		IP:            utils.Config.Host,
		Port:          utils.Config.Port,
		RouterManager: NewRouterManager(),
	}
	return &s
}

func (s *Server) Start() {
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%v:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("resolve err", err)
		return
	}

	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("listen err", err)
		return
	}
	fmt.Printf("%s服务已启动正在监听%s:%d\n", s.Name, s.IP, s.Port)
	s.RouterManager.StartWorkers()
	var connID uint32 = 1
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("accept err", err)
			continue
		}
		NewConnection(conn, connID, s.RouterManager).Start()
		connID++
	}
}

func (s *Server) Stop() {

}

func (s *Server) Serve() {
	defer s.Stop()
	s.Start()
}

func (s *Server) AddRouter(msgId uint32, routerManager ziface.IRouter) {
	s.RouterManager.AddRouter(msgId, routerManager)
}
