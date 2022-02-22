package main

import (
	"fmt"
	"gotoCabbage/core"
	"io"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp4", "0.0.0.0:6872")
	if err != nil {
		fmt.Println("client conn err", err)
		return
	}
	fmt.Println("client start")
	cnt := 1
	for {
		mid := 0
		cnt = -cnt
		if cnt == 1 {
			mid = 1
		} else {
			mid = 0
		}
		binaryMsgData, err := core.Dp.Pack(core.NewMessage(uint32(mid), []byte("hello wxf v1.0")))
		if err != nil {
			fmt.Println(err)
		}
		conn.Write(binaryMsgData)
		headData := make([]byte, core.Dp.GetHeadLen())
		if _, err := io.ReadFull(conn, headData); err != nil {
			fmt.Println("read msg head err", err)
			break
		}
		msg, err := core.Dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack err", err)
			break
		}
		if msg.GetMsgLen() > 0 {
			data := make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, data); err != nil {
				fmt.Println("read data err", err)
				break
			}
			msg.SetMsgData(data)
		}
		fmt.Println("recv msgID=", msg.GetMsgID(), "recv msgLen=", msg.GetMsgLen(), "recv msgData=", string(msg.GetMsgData()))
		time.Sleep(5 * time.Second)
	}
}
