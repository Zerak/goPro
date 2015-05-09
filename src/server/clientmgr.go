// package clientmgr
package server

import (
	"fmt"
	// "errors"
)

type Client struct {
	Id   int
	Name string
}

type ClientMgr struct {
	// Client
	// func AddClient(id int ){   }
}

func (cm ClientMgr) AddClient(id int) {
	fmt.Println("[clientMgr:AddClient] add client id[", id, "]")
}
