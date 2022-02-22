package core

import (
	"fmt"
	"gotoCabbage/interface"
	"gotoCabbage/utils/globalobj"
	"gotoCabbage/utils/log"
	"net"
	"strconv"
)

type Server struct {
	//服务器的名称
	Name string
	//服务器IP的版本
	IPVersion string
	//服务器监听的ip
	IP string
	//服务器监听的端口
	Port int
	//服务器的地址
	Addr string
	//路由
	MsgHandleMgr _interface.IMsgHandlerManager
	//连接管理模块
	ConnMgr _interface.IConnectionManager
	//连接建立之后的hook
	ConnStartHook func(connection _interface.IConnection)
	//连接关闭之前的hook
	ConnStopHook func(connection _interface.IConnection)
}

func (s *Server) SetConnStartHook(f func(connection _interface.IConnection)) {
	s.ConnStartHook = f
}

func (s *Server) SetConnStopHook(f func(connection _interface.IConnection)) {
	s.ConnStopHook = f
}

func (s *Server) CallConnStartHook(connection _interface.IConnection) {
	if s.ConnStartHook != nil {
		fmt.Println("CallConnStartHook")
		s.ConnStartHook(connection)
	}
}

func (s *Server) CallConnStopHook(connection _interface.IConnection) {
	if s.ConnStopHook != nil {
		fmt.Println("CallConnStopHook")
		s.ConnStopHook(connection)
	}
}

func (s *Server) AddRouter(msgID uint32, router _interface.IRouter) {
	s.MsgHandleMgr.AddRouter(msgID, router)
	fmt.Println("Add router success")
}

//启动服务器
func (s *Server) Start() {
	//获取tcp的addr
	s.Addr = s.IP + ":" + strconv.Itoa(s.Port)
	fmt.Println("[Start]Server listener at Address:", s.Addr, "is starting")
	//监听服务器地址
	listener, err := net.Listen(s.IPVersion, s.Addr)
	if err != nil {
		log.Fatal("listen", s.IPVersion, s.Addr, "err:", err)
		return
	} else {
		fmt.Println("[Start]", s.Name, "success now is listening")
	}
	//开启工作池
	s.MsgHandleMgr.StartWorkerPool()
	var cid uint32
	//阻塞等待客户端链接，处理客户端业务
	for {
		cid++
		conn, err := listener.Accept()
		if err != nil {
			log.Error("Accept err:", err)
			continue
		}
		if s.ConnMgr.Count() >= globalobj.GlobalObject.MaxConn {
			fmt.Println("=====>too many connection======<")
			conn.Close()
			continue
		}
		dealConn := NewConnection(s, conn, cid, s.MsgHandleMgr)
		go dealConn.Start()
	}
}

//停止服务器
func (s *Server) Stop() {
	//TODO 回收资源
	s.ConnMgr.ClearConn()
}

//运行服务器
func (s *Server) Serve() {
	//启动服务功能
	go s.Start()
	//TODO，做一些额外工作
	select {}
}

func NewServer() _interface.IServer {
	return &Server{
		Name:         globalobj.GlobalObject.Name,
		IPVersion:    "tcp4",
		IP:           globalobj.GlobalObject.IP,
		Port:         globalobj.GlobalObject.Port,
		MsgHandleMgr: NewMsgHandler(),
		ConnMgr:      NewConnManager(),
	}
}
func (s *Server) GetConnMgr() _interface.IConnectionManager {
	return s.ConnMgr
}
