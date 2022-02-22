package _interface

//服务器接口
type IServer interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//运行服务器
	Serve()
	//添加路由方法
	AddRouter(msgID uint32, router IRouter)

	GetConnMgr() IConnectionManager

	SetConnStartHook(func(connection IConnection))
	SetConnStopHook(func(connection IConnection))
	CallConnStartHook(connection IConnection)
	CallConnStopHook(connection IConnection)
}
