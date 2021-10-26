package znet

import (
	"fmt"
	"github.com/cwww3/zinx/utils"
	"github.com/cwww3/zinx/ziface"
	"net"
)

type Server struct {
	Name          string
	IPVersion     string
	IP            string
	Port          int
	RouterManager ziface.IRouterManager
	ConnManager   ziface.IConnManager
	OnConnStart   func(conn ziface.IConnection)
	OnConnStop    func(conn ziface.IConnection)
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnManager
}

func NewServer() ziface.IServer {
	s := Server{
		Name:          utils.Config.Name,
		IPVersion:     "tcp4",
		IP:            utils.Config.Host,
		Port:          utils.Config.Port,
		RouterManager: NewRouterManager(),
		ConnManager:   NewConnManager(),
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
		if s.ConnManager.Len() >= utils.Config.MaxConnSize {
			fmt.Println("too many conn ")
			conn.Close()
			continue
		}
		connection := NewConnection(s, conn, connID, s.RouterManager)
		connection.Start()

		connID++
	}
}

func (s *Server) Stop() {
	fmt.Println("server stop")
	s.ConnManager.ClearConnection()
}

func (s *Server) Serve() {
	defer s.Stop()
	s.Start()
}

func (s *Server) AddRouter(msgId uint32, routerManager ziface.IRouter) {
	s.RouterManager.AddRouter(msgId, routerManager)
}

func (s *Server) SetOnConnStart(f func(conn ziface.IConnection)) {
	s.OnConnStart = f
}

func (s *Server) SetOnConnStop(f func(conn ziface.IConnection)) {
	s.OnConnStop = f
}

func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		s.OnConnStop(conn)
	}
}
