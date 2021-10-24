package ziface

type IConnManager interface {
	AddConnection(connection IConnection)
	RemoveConnection(connection IConnection)
	GetConnection(connId uint32) (IConnection, error)
	Len() int
	ClearConnection()
}
