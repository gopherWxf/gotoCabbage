package main

import (
	"fmt"
	"gotoCabbage/core"
	"gotoCabbage/interface"
	"gotoCabbage/utils/log"
)

func main() {
	fmt.Print(`
                __         _________       ___.  ___.                          
   ____   _____/  |_  ____ \_   ___ \_____ \_ |__\_ |__ _____     ____   ____  
  / ___\ /  _ \   __\/  _ \/    \  \/\__  \ | __ \| __ \\__  \   / ___\_/ __ \ 
 / /_/  >  <_> )  | (  <_> )     \____/ __ \| \_\ \ \_\ \/ __ \_/ /_/  >  ___/ 
 \___  / \____/|__|  \____/ \______  (____  /___  /___  (____  /\___  / \___  >
/_____/                            \/     \/    \/    \/     \//_____/      \/ 
`)
	log.Setup(&log.Settings{
		Path:       "logs",
		Name:       "filename",
		Ext:        "log",
		TimeFormat: "2006-01-02",
		MaxBackups: 20, //最多MaxBackups个文件
		MaxSize:    10, //文件最大MaxSize(k)
		Cnt:        1,  //从filename_cnt开始，后期做config配置
	})

	server := core.NewServer()
	server.AddRouter(0, &PingRouter{})
	server.AddRouter(1, &HelloRouter{})
	server.SetConnStartHook(DoConnBegin)
	server.SetConnStopHook(DoConnStop)
	server.Serve()
}

//自定义路由
type PingRouter struct {
	core.BaseRouter
}

func (p *PingRouter) Handle(iRequest _interface.IRequest) {
	fmt.Println("Call Router Handle")
	//先读取客户端的数据
	fmt.Println("msgID=", iRequest.GetMgsID(), "msg=", string(iRequest.GetData()))
	err := iRequest.GetConnection().SendMsg(200, []byte("ping ping ping"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloRouter struct {
	core.BaseRouter
}

func (h *HelloRouter) Handle(iRequest _interface.IRequest) {
	fmt.Println("Call Router Handle")
	//先读取客户端的数据
	fmt.Println("msgID=", iRequest.GetMgsID(), "msg=", string(iRequest.GetData()))
	err := iRequest.GetConnection().SendMsg(404, []byte("err err err"))
	if err != nil {
		fmt.Println(err)
	}
}

//连接链接之后的hook
func DoConnBegin(connection _interface.IConnection) {
	fmt.Println("====>Do Conn Begin<====")
	connection.SendMsg(200, []byte("Do Conn Begin"))
	fmt.Println("set conn Name")
	connection.SetProperty("name", "wxf")
	connection.SetProperty("age", "20")

}

//连接销毁之前的hook
func DoConnStop(connection _interface.IConnection) {
	fmt.Println("====>Do Conn Stop<====")
	if name, err := connection.GetProperty("name"); err == nil {
		fmt.Println(name)
	}
	if age, err := connection.GetProperty("age"); err == nil {
		fmt.Println(age)
	}
}
