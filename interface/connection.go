package _interface

import "net"

type IConnection interface {
	//启动连接
	Start()
	//停止连接
	Stop()
	//获取当前连接绑定的socket conn
	GetSocketConn() net.Conn
	//获取当前连接的id
	GetConnID() uint32
	//获取远程client的地址
	GetRemoteAddr() string
	//发送数据
	SendMsg(msgID uint32, data []byte) error
	//设置属性
	SetProperty(key string, val interface{})
	//获取属性
	GetProperty(key string) (val interface{}, err error)
	//移除属性
	RemobeProperty(key string)
}

//处理业务的方法
type HandleFunc func(net.Conn, []byte, int) error
