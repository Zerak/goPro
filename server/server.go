package main

import (
	"fmt"
	"net"
	"os"

	"clientmgr"
)

func Handler(conn net.Conn, msgs chan string) {
	fmt.Println("[server:Handler] connection from[", conn.RemoteAddr().String() , "]")

	buf := make([]byte, 1024)
	for {
		length, err := conn.Read(buf)
		if err != nil {
			fmt.Println("[server:Handler] the client[", conn.RemoteAddr().String() , "] Closed")
			conn.Close()
			break
		}
		if length > 0 {
			buf[length] = 0
		}

		recvStr := string(buf[0:length])
		msgs <- recvStr
	}
}

func echoHandler(conns *map[string]net.Conn, msgs chan string) {
	for {
		msg := <-msgs
		fmt.Println("[server:echoHandler] msg [" + msg + "]")

		for key, val := range *conns {
			//fmt.Println("[server:echoHandler]connection from ->", key)
			_, err := val.Write([]byte(msg))
			if err != nil {
				fmt.Println("[server:echoHandler]recive err ", err.Error())
				delete(*conns, key)
			}
		}
	}
}

func StartServer(port string) {
	service := ":" + port
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		fmt.Println("[server]ResolveTCPAddr err")
		return
	}

	listen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("[server]ListenTCP err")
		return
	}
	fmt.Println("[server]ListenTCP ok")

	conns := make(map[string]net.Conn)
	msgs := make(chan string, 100)

	go echoHandler(&conns, msgs)

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("[server]Accept err")
			continue
		}
		conns[conn.RemoteAddr().String()] = conn

		go Handler(conn, msgs)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("wrong args like this\n [port]")
		os.Exit(0)
	}

	if len(os.Args) == 2 {
		StartServer(os.Args[1])
	} else {
		fmt.Println("wrong args like this\n [port]")
	}

	fmt.Println("exit")
}
