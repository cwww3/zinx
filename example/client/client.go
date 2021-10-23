package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/utils"
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
			buf := make([]byte, utils.Config.MaxPackageSize)
			_, err = conn.Read(buf)
			if err == io.EOF {
				fmt.Println("read down")
				return
			}
			if err != nil {
				fmt.Println("read err", err)
				break
			}
			fmt.Println("read content", string(buf))
		}
	}()
	for {
		_, err = conn.Write([]byte("Hello World!"))
		if err != nil {
			fmt.Println("write err", err)
		}
		time.Sleep(time.Second * 3)
	}
}
