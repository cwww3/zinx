package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/cwww3/zinx/utils"
	"github.com/cwww3/zinx/ziface"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return new(DataPack)
}

func (d *DataPack) GetHeadLen() uint32 {
	return 8
}

func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	var err error
	buf := bytes.NewBuffer([]byte{})
	if err = binary.Write(buf, binary.LittleEndian, msg.GetLength()); err != nil {
		return nil, err
	}
	if err = binary.Write(buf, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}
	if err = binary.Write(buf, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

func (d *DataPack) Unpack(data []byte) (ziface.IMessage, error) {
	reader := bytes.NewReader(data)
	msg := &Message{}
	var err error
	if err = binary.Read(reader, binary.LittleEndian, &msg.Length); err != nil {
		return nil, err
	}
	if err = binary.Read(reader, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	if utils.Config.MaxPackageSize > 0 && utils.Config.MaxPackageSize < msg.Length {
		err = errors.New("pack too large")
		return nil, err
	}
	return msg, err
}
