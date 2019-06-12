package main

import (
	"fmt"
	"net"

	"archerystar/common/constants"
	"archerystar/protocol/protos"
	"archerystar/testcase/client"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:2333")
	if err != nil {
		fmt.Println("Dial error:", err)
	}

	defer conn.Close()

	req := &protos.GameRoomReq{
		Mode:    1,
		Stage:   3,
		Players: []*protos.RoomPlayer{},
	}

	player1 := &protos.RoomPlayer{
		Id:       1,
		Account:  "A0000000001",
		Name:     "GNU",
		Country:  1,
		Trophies: 36,
		Equips:   []int32{1, 2, 3, 4, 5},
		Head:     3,
		Order:    1,
	}

	player2 := &protos.RoomPlayer{
		Id:       2,
		Account:  "B0000000002",
		Name:     "POSIX",
		Country:  23,
		Trophies: 100,
		Equips:   []int32{4, 3, 8, 10, 6},
		Head:     1,
		Order:    2,
	}

	req.Players = append(req.Players, player1, player2)

	if ret := client.SendRequest(conn, constants.CreateRoom, req); !ret {
		return
	}

	resp := &protos.GameRoomResp{}

	if ret := client.ReadResp(conn, resp); ret {
		return
	}
}
