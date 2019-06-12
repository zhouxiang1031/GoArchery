package gamesvc

import (
	"archerystar/common/constants"
	"archerystar/component/logger"
	"archerystar/message"
)

type GameMsgFuncs map[uint]func(*message.NtolMessage)

//var gameMsgFuncs = &GameMsgFuncs{}

func (f *GameMsgFuncs) registerFun(msgid uint, msgfun func(*message.NtolMessage)) {
	(*f)[msgid] = msgfun
}

func (f *GameMsgFuncs) RegisterAllFuns(svc *GameService) {
	f.registerFun(constants.TestMsg, svc.testMsgProc)
	f.registerFun(constants.CreateRoom, svc.creatRoomProc)
}

func (f *GameMsgFuncs) Func(msg *message.NtolMessage) {
	if fun, ok := (*f)[msg.Msg.ID]; ok {
		fun(msg)
	} else {
		logger.Error(titleGameService, "not found func for msgid:%d", msg.Msg.ID)
	}
}
