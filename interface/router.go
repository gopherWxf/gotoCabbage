package _interface

type IRouter interface {
	//处理业务之前
	PreHandle(iRequest IRequest)
	//处理业务
	Handle(iRequest IRequest)
	//处理业务之后
	PostHandle(iRequest IRequest)
}
