package main

import (
	"archerystar/common/constants"
	"archerystar/protocol/protos"
	"archerystar/testcase/client"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:2333")
	if err != nil {
		fmt.Println("Dial error:", err)
	}

	defer conn.Close()

	user := &protos.TestMsgReq{}
	user.User = "fanic"
	user.Desc = "hello server"

	if ret := client.SendRequest(conn, constants.TestMsg, user); !ret {
		return
	}

	testResp := &protos.TestMsgResp{}

	if ret := client.ReadResp(conn, testResp); ret {
		return
	}
}
