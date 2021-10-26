package main

import (
	"fmt"
	"github.com/cwww3/zinx/znet"
	"io"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp4", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("dial err", err)
		return
	}
	defer conn.Close()
	go func() {
		for {
			var err error
			dp := znet.NewDataPack()
			buf := make([]byte, dp.GetHeadLen())
			_, err = io.ReadFull(conn, buf)
			if err == io.EOF {
				fmt.Println("read done")
				break
			}
			if err != nil {
				fmt.Println("read err", err)
				break
			}
			msg, err := dp.Unpack(buf)
			if err != nil {
				fmt.Println("unpack err", err)
				break
			}
			content := make([]byte, msg.GetLength())
			_, err = io.ReadFull(conn, content)
			if err != nil {
				fmt.Println("read content err", err)
				break
			}
			msg.(*znet.Message).Data = content
			fmt.Println("read content", string(content))
			fmt.Println("read msg", msg)
		}
	}()
	cnt := 0
	for {
		dp := znet.NewDataPack()
		var msgId uint32
		var content string
		if cnt%2 == 0 {
			msgId = 1
			content = "ping1 Hello World!"
		} else {
			msgId = 2
			content = "ping2 Hello World!"
		}
		data, err := dp.Pack(znet.NewMessage(msgId, []byte(content)))
		if err != nil {
			fmt.Println("pack err", err)
			break
		}
		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("write err", err)
			break
		}
		cnt++
		time.Sleep(time.Second * 3)
	}
}
