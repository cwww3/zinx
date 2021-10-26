package znet

import (
	"errors"
	"fmt"
	"github.com/cwww3/zinx/ziface"
	"sync"
)

type ConnManager struct {
	sync.Mutex
	connections map[uint32]ziface.IConnection
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		Mutex:       sync.Mutex{},
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (c *ConnManager) AddConnection(connection ziface.IConnection) {
	c.Lock()
	defer c.Unlock()
	c.connections[connection.GetConnID()] = connection
}

func (c *ConnManager) RemoveConnection(connection ziface.IConnection) {
	c.Lock()
	defer c.Unlock()
	delete(c.connections, connection.GetConnID())
}

func (c *ConnManager) GetConnection(connId uint32) (ziface.IConnection, error) {
	c.Lock()
	defer c.Unlock()
	if conn, ok := c.connections[connId]; !ok {
		return nil, errors.New(fmt.Sprintf("未找到连接 id=%d\n", connId))
	} else {
		return conn, nil
	}
}

func (c *ConnManager) Len() int {
	return len(c.connections)
}

func (c *ConnManager) ClearConnection() {
	c.Lock()
	defer c.Unlock()
	for _, conn := range c.connections {
		conn.Stop()
		delete(c.connections, conn.GetConnID())
	}
}
