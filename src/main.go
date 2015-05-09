package main

import (
	"client"
	"fmt"
	"os"
	"server"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("wrong args like this\n[server port] or [client addr:port] ok")
		os.Exit(0)
	}
	if os.Args[1] == "server" && len(os.Args) == 3 {
		server.StartServer(os.Args[2])
	}

	if os.Args[1] == "client" && len(os.Args) == 3 {
		client.StartClient(os.Args[2])
	}
	fmt.Println("exit")
}
