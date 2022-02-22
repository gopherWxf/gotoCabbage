package _interface

type IConnectionManager interface {
	Add(conn IConnection)
	Remove(conn IConnection)
	Get(connID uint32) (IConnection, error)
	Count() int
	ClearConn()
}
