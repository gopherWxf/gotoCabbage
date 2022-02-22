package core

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack_Pack(t *testing.T) {
	listener, _ := net.Listen("tcp4", "127.0.0.1:6872")
	go func() {
		for {
			conn, _ := listener.Accept()
			go func(conn2 net.Conn) {
				dp := NewDataPack()
				for {
					headData := make([]byte, dp.GetHeadLen())
					io.ReadFull(conn, headData)
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println(err)
					}
					if msgHead.GetMsgLen() > 0 {
						data := make([]byte, msgHead.GetMsgLen())
						io.ReadFull(conn, data)
						msgHead.SetMsgData(data)
					}
					fmt.Println(msgHead.GetMsgID(), msgHead.GetMsgLen(), string(msgHead.GetMsgData()))
				}
			}(conn)
		}
	}()
	conn, _ := net.Dial("tcp4", "127.0.0.1:6872")
	dp := NewDataPack()
	msg1 := &Message{
		ID:      1,
		DataLen: 5,
		Data:    []byte("wxfaa"),
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println(err)
	}
	msg2 := &Message{
		ID:      1,
		DataLen: 8,
		Data:    []byte("qwerasdf"),
	}
	sendData2, err := dp.Pack(msg2)
	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)
	select {}
}
