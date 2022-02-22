package core

import (
	"fmt"
	"gotoCabbage/interface"
	"gotoCabbage/utils/globalobj"
	"strconv"
)

type MsgHandlerManager struct {
	Apis map[uint32]_interface.IRouter
	//负责worker取任务的消息队列
	TaskQueue chan _interface.IRequest
	//worker数量
	WorkerPoolSize uint32
	//消息队列的长度
	TaskQueueChanSize uint32
}

func NewMsgHandler() *MsgHandlerManager {
	return &MsgHandlerManager{
		Apis:              map[uint32]_interface.IRouter{},
		TaskQueue:         make(chan _interface.IRequest, globalobj.GlobalObject.WorkerPoolSize),
		WorkerPoolSize:    globalobj.GlobalObject.WorkerPoolSize,
		TaskQueueChanSize: globalobj.GlobalObject.TaskQueueChanSize,
	}
}
func (m *MsgHandlerManager) DoMsgHandler(request _interface.IRequest) {
	handler, ok := m.Apis[request.GetMgsID()]
	if !ok {
		fmt.Println("api msgID", request.GetMgsID(), "is not found")
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (m *MsgHandlerManager) AddRouter(msgID uint32, router _interface.IRouter) {
	if _, ok := m.Apis[msgID]; ok {
		panic("repeat api,msgID=" + strconv.Itoa(int(msgID)))
	}
	m.Apis[msgID] = router
	fmt.Println("Add Api msgID", msgID, "succ")
}

//启动worker工作池
func (m *MsgHandlerManager) StartWorkerPool() {
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		go m.startOneWorker(i)
	}
}

//启动worker流程
func (m *MsgHandlerManager) startOneWorker(workerID int) {
	fmt.Println("workID=", workerID, "is started")
	for request := range m.TaskQueue {
		fmt.Println("workID=", workerID, "deal with", request.GetConnection().GetConnID())
		m.DoMsgHandler(request)
	}
}
func (m *MsgHandlerManager) SendMsgToTaskQueue(request _interface.IRequest) {
	m.TaskQueue <- request
}
