package _interface

type IMsgHandlerManager interface {
	DoMsgHandler(request IRequest)
	AddRouter(msgID uint32, router IRouter)
	StartWorkerPool()
	SendMsgToTaskQueue(request IRequest)
}
