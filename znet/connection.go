package znet

import (
	"fmt"
	"io"
	"net"
	"zinx/ziface"
)

type Connection struct {
	Conn      *net.TCPConn
	ConnID    uint32
	isClosed  bool
	handleAPI ziface.HandleFunc
	ExitCh    chan bool // 通知连接关闭
}

func (c *Connection) Start() {
	fmt.Println("conn start id", c.ConnID)
	c.StartHandle()
}

func (c *Connection) Stop() {
	if c.isClosed == true {
		fmt.Println("conn", c.ConnID, "has already closed")
		return
	}
	c.isClosed = true
	c.Conn.Close()
	close(c.ExitCh)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	panic("implement me")
}

func NewConnection(conn *net.TCPConn, connID uint32, callbackAPI ziface.HandleFunc) ziface.IConnection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		isClosed:  false,
		handleAPI: callbackAPI,
		ExitCh:    make(chan bool, 1),
	}
	return c
}

func (c *Connection) StartHandle() {
	defer fmt.Println(c.ConnID, "stop")
	defer c.Stop()
	conn := c.Conn
	for {
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err == io.EOF {
			fmt.Println(c.ConnID, "read done")
			break
		}
		if err != nil {
			fmt.Println(c.ConnID, "read err", err)
			break
		}
		fmt.Println(c.ConnID, "read content", string(buf))
		err = c.handleAPI(conn, buf[:], cnt)
		if err != nil {
			fmt.Println(c.ConnID, "write err", err)
			break
		}
	}
}
