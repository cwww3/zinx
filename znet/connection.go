package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	Conn          *net.TCPConn
	ConnID        uint32
	isClosed      bool
	ExitCh        chan bool // 通知连接关闭
	MsgCh         chan []byte
	RouterManager ziface.IRouterManager
}

func NewConnection(conn *net.TCPConn, connID uint32, routerManager ziface.IRouterManager) ziface.IConnection {
	c := &Connection{
		Conn:          conn,
		ConnID:        connID,
		isClosed:      false,
		RouterManager: routerManager,
		MsgCh:         make(chan []byte, 3),
		ExitCh:        make(chan bool),
	}
	return c
}

func (c *Connection) Start() {
	fmt.Println("conn start id", c.ConnID)
	go c.StartRead()
	go c.StartWrite()
}

func (c *Connection) Stop() {
	if c.isClosed == true {
		fmt.Println("conn", c.ConnID, "has already closed")
		return
	}
	c.isClosed = true
	c.Conn.Close()
	close(c.ExitCh)
	close(c.MsgCh)
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
	c.MsgCh <- pack
	return err
}

func (c *Connection) StartRead() {
	defer fmt.Println(c.ConnID, "read stop")
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
		if utils.Config.WorkerSize > 0 {
			c.RouterManager.SendRequest(req)
		} else {
			go c.RouterManager.DoMsgHandler(req)
		}
	}
}

func (c *Connection) StartWrite() {
	defer fmt.Println(c.ConnID, "write stop")
	for {
		select {
		case data := <-c.MsgCh:
			_, err := c.GetTCPConnection().Write(data)
			if err != nil {
				fmt.Println(c.ConnID, "write err", err)
				return
			}
		case <-c.ExitCh:
			return
		}
	}
}
