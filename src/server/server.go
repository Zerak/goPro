package server

import (
	"fmt"
	"net"
	// "os"

	"github.com/garyburd/redigo/redis"
)

func Handler(conn net.Conn, msgs chan string) {
	fmt.Println("[server:Handler] connection from[", conn.RemoteAddr().String(), "]")

	buf := make([]byte, 1024)
	for {
		length, err := conn.Read(buf)
		if err != nil {
			fmt.Println("[server:Handler] the client[", conn.RemoteAddr().String(), "] Closed")
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
	fmt.Println("[server]StartServer ...")

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

	// init redis db
	client, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("[server]connect to redis Err.", err)
		return
	} else {
		fmt.Println("[server]connect to redis OK")
	}
	defer client.Close()
	//|||||||||||||||||||||||||||||||||||||||||||||||||||||||||
	fmt.Println("[server] test redis")
	// _, err := 
	client.Do("set", "username", "nick")
	if err != nil {
		fmt.Println("[server] set name err")
	} else {
		fmt.Println("[server] set name ok")
	}

	username, err := redis.String(client.Do("get","username"))
	if err != nil {
		fmt.Println("[server] get name err")
	}else{
		fmt.Println("[server] get name ", username)
	}
	//|||||||||||||||||||||||||||||||||||||||||||||||||||||||||

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
