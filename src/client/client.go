package client

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func StartClient(tcpaddr string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", tcpaddr)
	if err != nil {
		fmt.Println("ResolveTCPAddr err")
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("DiaTCP err")
		return
	}
	go chatSend(conn)

	buf := make([]byte, 1024)
	for {
		length, err := conn.Read(buf)
		if err != nil {
			conn.Close()
			fmt.Println("servive may down. exit")
			os.Exit(0)
		}

		fmt.Println("recv from server[" + string(buf[0:length]) + "]")
	}
}

func chatSend(conn net.Conn) {
	var input string
	uname := conn.LocalAddr().String()
	for {
		fmt.Scanln(&input)
		if input == "/quit" {
			fmt.Println("bye...")
			conn.Close()
			os.Exit(0)
		}

		length, err := conn.Write([]byte(uname + " Say :::" + input))
		fmt.Println("send len[" + strconv.Itoa(length) + "]")
		if err != nil {
			fmt.Println("client write err info->" + err.Error())
			conn.Close()
			break
		}
	}
}
