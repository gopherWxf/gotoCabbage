package _interface

type IDataPack interface {
	GetHeadLen() uint32
	//封包
	Pack(message IMessage) ([]byte, error)
	//拆包
	Unpack([]byte) (IMessage, error)
}
