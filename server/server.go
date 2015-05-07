package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func Handler(conn net.Conn, msgs chan string) {
	fmt.Println("[server:Handler] connection from[", conn.RemoteAddr().String())

	buf := make([]byte, 1024)
	for {
		length, err := conn.Read(buf)
		if err != nil {
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
		fmt.Println("[server:echoHandler]recive msg [" + msg + "]")

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

	conns := make(map[string]net.Conn)
	msgs := make(chan string, 100)

	go echoHandler(&conns, msgs)

	for {
		fmt.Println("[server]listening ...")

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
	if len(os.Args) != 3 {
		fmt.Println("wrong args like this\nserver port] or [client addr:port]")
		os.Exit(0)
	}

	if os.Args[1] == "server" && len(os.Args) == 3 {
		StartServer(os.Args[2])
	}

	if os.Args[1] == "client" && len(os.Args) == 3 {
		StartClient(os.Args[2])
	}
	fmt.Println("execute")
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
		if  err != nil {
			fmt.Println("client write err info->" + err.Error())
			conn.Close()
			break
		}
	}
}



