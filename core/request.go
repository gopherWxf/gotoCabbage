package core

import (
	"gotoCabbage/interface"
)

type Request struct {
	//已经封装好的连接
	conn _interface.IConnection
	//客户端请求的数据
	msg _interface.IMessage
}

func (r *Request) GetConnection() _interface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetMsgData()
}
func (r *Request) GetMgsID() uint32 {
	return r.msg.GetMsgID()
}
func NewRequest(conn _interface.IConnection, msg _interface.IMessage) _interface.IRequest {
	return &Request{
		conn: conn,
		msg:  msg,
	}
}
