package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	Router    ziface.IRouter
}

func NewServer(name string) ziface.IServer {
	s := Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
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
	var connID uint32 = 1
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("accept err", err)
			continue
		}
		connection := NewConnection(conn, connID, s.Router)
		connID++
		go connection.Start()
	}
}

func (s *Server) Stop() {

}

func (s *Server) Serve() {
	defer s.Stop()
	s.Start()
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
}
