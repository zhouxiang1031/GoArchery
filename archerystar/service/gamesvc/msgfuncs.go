package gamesvc

import (
	"fmt"

	"archerystar/component/logger"
	"archerystar/message"
	gmpb "archerystar/protocol/protos"
	"archerystar/route"
	"archerystar/service/gamesvc/gameroom"

	"github.com/golang/protobuf/proto"
)

func (s *GameService) testMsgProc(msg *message.NtolMessage) {
	logger.Debug(titleGameService, "begin testMsgProc")

	testReq := &gmpb.TestMsgReq{}

	proto.Unmarshal(msg.Msg.Data, testReq)

	fmt.Printf("test req:%+v", testReq)

	//resp
	s.resetMsgForReso(msg)

	testResp := &gmpb.TestMsgResp{
		User: "server",
		Desc: "hello, I'm server!",
		Flag: 9527,
	}
	var ok error
	if msg.Msg.Data, ok = proto.Marshal(testResp); ok != nil {
		fmt.Println("Marshal proto failed:", ok.Error())
		return
	}

	route.Next(s, msg)
}

func (s *GameService) creatRoomProc(msg *message.NtolMessage) {
	logger.Debug(titleGameService, "begin creatRoomProc")

	roomReq := &gmpb.GameRoomReq{}

	proto.Unmarshal(msg.Msg.Data, roomReq)

	logger.Debug(titleGameService, "test req:%+v", roomReq)

	//create room
	room := gameroom.GetRoomsManager().CreateRomm(roomReq)

	logic := gameroom.NewRoomLogic(room)
	go logic.Handle()

	s.resetMsgForReso(msg)
	var ok error
	resp := &gmpb.GameRoomResp{
		Id:      room.Id,
		Mode:    room.Mode,
		Stage:   room.Stage,
		LevId:   3,
		Players: roomReq.Players[:len(roomReq.Players)],
	}

	if msg.Msg.Data, ok = proto.Marshal(resp); ok != nil {
		logger.Error(titleGameService, "Marshal proto failed:", ok.Error())
		return
	}

	route.Next(s, msg)

}

func (s *GameService) joinRoomProc(msg *message.NtolMessage) {
	logger.Debug(titleGameService, "begin creatRoomProc")

	req := &gmpb.JoinRoomReq{}

	proto.Unmarshal(msg.Msg.Data, req)

	logger.Debug(titleGameService, "test req:%+v", req)

	//create room
	room, ok := gameroom.GetRoomsManager().FindRoom(req.RoomId)
	if !ok {
		logger.Error(titleGameService, "not found room for id %d", req.RoomId)
		resp := &gmpb.JoinRoomResp{}
		resp.IsReady = false

		s.resetMsgForReso(msg)

		var err error
		if msg.Msg.Data, err = proto.Marshal(resp); err != nil {
			logger.Error(titleGameService, "Marshal proto failed:", err.Error())
			return
		}
		route.Next(s, msg)
		return
	}

	room.Session[req.PlayerId-1] = msg.SessionMsg
	cmd := gameroom.StateCmdMsg{
		Cmd:          gameroom.Cmd_JoinRoom,
		PlayerStatus: nil,
	}

	room.ChCmd <- cmd
}

func (s *GameService) resetMsgForReso(msg *message.NtolMessage) {
	msg.Msg.Type = message.Response
	msg.Msg.Err = false
	msg.Msg.Data = msg.Msg.Data[:0]
}
