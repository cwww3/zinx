package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/ziface"
)

type Connection struct {
	Conn     *net.TCPConn
	ConnID   uint32
	isClosed bool
	ExitCh   chan bool // 通知连接关闭
	Router   ziface.IRouter
}

func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) ziface.IConnection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		Router:   router,
		ExitCh:   make(chan bool, 1),
	}
	return c
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

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection is already closed")
	}
	dp := NewDataPack()
	msg := NewMessage(msgId, data)
	pack, err := dp.Pack(msg)
	if err != nil {
		return err
	}
	_, err = c.GetTCPConnection().Write(pack)
	return err
}

func (c *Connection) StartHandle() {
	defer fmt.Println(c.ConnID, "stop")
	defer c.Stop()
	for {
		dp := NewDataPack()
		buf := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.GetTCPConnection(), buf)
		if err != nil {
			fmt.Println(c.ConnID, "read head err", err)
			break
		}
		msg, err := dp.Unpack(buf)
		if err != nil {
			fmt.Println("unpack err", err)
			break
		}
		if msg.GetLength() > 0 {
			data := make([]byte, msg.GetLength())
			_, err = io.ReadFull(c.GetTCPConnection(), data)
			if err != nil {
				fmt.Println(c.ConnID, "read content err", err)
				break
			}
			msg.SetMsgData(data)
			fmt.Println(c.ConnID, "read content", string(data), "msgId", msg.GetMsgID())
		}

		req := &Request{
			connection: c,
			msg:        msg,
		}
		c.Router.PreHandle(req)
		c.Router.Handle(req)
		c.Router.PostHandle(req)
	}
}
