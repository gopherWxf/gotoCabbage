package globalobj

import (
	"encoding/json"
	server2 "gotoCabbage/interface"
	"io/ioutil"
)

type GlobalObj struct {
	Server            server2.IServer
	IP                string
	Port              int
	Name              string
	Version           string
	MaxConn           int
	MaxPackageSize    uint32
	WorkerPoolSize    uint32
	TaskQueueChanSize uint32
}

var GlobalObject *GlobalObj

func init() {
	//默认
	GlobalObject = &GlobalObj{
		IP:                "0.0.0.0",
		Port:              6872,
		Name:              "Cabbage",
		Version:           "v1.0",
		MaxConn:           1000,
		MaxPackageSize:    4096,
		WorkerPoolSize:    10,
		TaskQueueChanSize: 4096,
	}
	//加载
	GlobalObject.Reload()
}
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/Cabbage.json")
	if err != nil {
		panic(err)
	}
	//解析
	err = json.Unmarshal(data, GlobalObject)
	if err != nil {
		panic(err)
	}
}
