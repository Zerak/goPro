package clientmgr

import(
	"fmt"
)

type Client{
	Id	int
	Name string

}

type ClientMgr interface{
	Client client

	AddClient(id int )
}

func (cm ClientMgr) AddClient(id int) {
	fmt.Printfln("[clientMgr:AddClient] add client id[" , id, "]")
	
}