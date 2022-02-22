package core

import (
	"errors"
	"fmt"
	"gotoCabbage/interface"
	"gotoCabbage/utils/log"
	"io"
	"net"
	"sync"
)

type Connection struct {
	//当前连接的socket
	SocketConn net.Conn
	//连接的id
	ConnID uint32
	//当前连接的状态
	isClosed bool
	//告知当前连接已经停止的channel
	ExitChan chan bool
	//路由
	MsgHandleManage _interface.IMsgHandlerManager
	//读写协程通信chan
	msgChan chan []byte
	//当前Conn隶属于那个Server
	parent _interface.IServer
	//链接属性的集合
	property      map[string]interface{}
	propertyMutex sync.RWMutex
}

func (c *Connection) SetProperty(key string, val interface{}) {
	c.propertyMutex.Lock()
	defer c.propertyMutex.Unlock()
	c.property[key] = val
}

func (c *Connection) GetProperty(key string) (val interface{}, err error) {
	c.propertyMutex.RLock()
	defer c.propertyMutex.RUnlock()
	val, ok := c.property[key]
	if ok {
		return val, nil
	}
	return nil, errors.New("not found property")
}

func (c *Connection) RemobeProperty(key string) {
	c.propertyMutex.Lock()
	defer c.propertyMutex.Unlock()
	delete(c.property, key)
}

//读业务
func (c *Connection) StartRead() {
	fmt.Println("[Read Goroutine] is running")
	defer c.Stop()
	for {
		headData := make([]byte, Dp.GetHeadLen())
		if _, err := io.ReadFull(c.SocketConn, headData); err != nil {
			fmt.Println("Read Goroutine: read msg err:", err)
			break
		}
		msg, err := Dp.Unpack(headData)
		if err != nil {
			fmt.Println("Read Goroutine: unpack err:", err)
			break
		}
		if msg.GetMsgLen() > 0 {
			data := make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.SocketConn, data); err != nil {
				log.Error("Read Goroutine: readFull data err:", err)
				break
			}
			msg.SetMsgData(data)
		}
		req := NewRequest(c, msg)
		//交给工作池来处理
		c.MsgHandleManage.SendMsgToTaskQueue(req)
	}
}

//写业务
func (c *Connection) StartWrite() {
	fmt.Println("[Write Goroutine] is running")
	defer fmt.Println(c.SocketConn.RemoteAddr().String(), "[write goroutine] is exit")
	for {
		select {
		case binaryData := <-c.msgChan:
			if _, err := c.SocketConn.Write(binaryData); err != nil {
				log.Error("Write Goroutine: write msg err:", err)
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn ID:", c.ConnID, " Start")
	//启动读业务
	go c.StartRead()
	//启动写业务
	go c.StartWrite()
	//执行hook
	c.parent.CallConnStartHook(c)
}

func (c *Connection) Stop() {
	fmt.Println("Conn ID:", c.ConnID, " Stop Remote Addr=", c.GetRemoteAddr())
	if c.isClosed {
		return
	}
	//关闭连接
	c.isClosed = true
	c.ExitChan <- true
	//连接销毁前执行的hook
	c.parent.CallConnStopHook(c)
	c.SocketConn.Close()
	//回收资源
	c.parent.GetConnMgr().Remove(c)
	close(c.ExitChan)
	close(c.msgChan)
}

func (c *Connection) GetSocketConn() net.Conn {
	return c.SocketConn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) GetRemoteAddr() string {
	return c.SocketConn.RemoteAddr().String()
}

//对发给客户端的数据进行封包再发送
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection is closed")
	}
	//进行封包
	binaryMsgData, err := Dp.Pack(NewMessage(msgID, data))
	if err != nil {
		fmt.Println("pack err", err)
		return errors.New("pack msg err")
	}
	c.msgChan <- binaryMsgData
	return nil
}
func NewConnection(server _interface.IServer, conn net.Conn, connID uint32, msgHandler _interface.IMsgHandlerManager) _interface.IConnection {
	connection := &Connection{
		parent:          server,
		SocketConn:      conn,
		ConnID:          connID,
		isClosed:        false,
		MsgHandleManage: msgHandler,
		ExitChan:        make(chan bool, 1),
		msgChan:         make(chan []byte),
		property:        make(map[string]interface{}),
	}
	connection.parent.GetConnMgr().Add(connection)
	return connection
}
