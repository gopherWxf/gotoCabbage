package _interface

type IRequest interface {
	//得到当前连接
	GetConnection() IConnection
	//得到请求的消息
	GetData() []byte

	GetMgsID() uint32
}
