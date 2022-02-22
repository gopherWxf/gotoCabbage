package core

import (
	"bytes"
	"encoding/binary"
	"errors"
	"gotoCabbage/interface"
	"gotoCabbage/utils/globalobj"
)

type DataPack struct {
}

func NewDataPack() _interface.IDataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32 {
	//DataLen 4 + ID 4
	return 8
}

var Dp DataPack

func (d *DataPack) Pack(message _interface.IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})
	var err error
	//将datalen写进databuff中
	err = binary.Write(dataBuff, binary.LittleEndian, message.GetMsgLen())
	if err != nil {
		return nil, err
	}
	//将id写进databuff中
	err = binary.Write(dataBuff, binary.LittleEndian, message.GetMsgID())
	if err != nil {
		return nil, err
	}
	//将data写进databuff中
	err = binary.Write(dataBuff, binary.LittleEndian, message.GetMsgData())
	if err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

func (d *DataPack) Unpack(binaryData []byte) (_interface.IMessage, error) {
	dataBuff := bytes.NewBuffer(binaryData)
	msg := &Message{}
	//读datalen
	binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen)
	//读id
	binary.Read(dataBuff, binary.LittleEndian, &msg.ID)
	//读data
	if globalobj.GlobalObject.MaxPackageSize < msg.GetMsgLen() {
		return nil, errors.New("too large msg data len")
	}
	return msg, nil
}
