package gamesvc

import (
	"archerystar/common/constants"
	"archerystar/component/logger"
	"archerystar/message"
	"sync"
)

var once = sync.Once{}

const titleGameService = "gameservice"

type GameService struct {
	ChSrv chan message.NtolMessage

	gameMsgFuncs *GameMsgFuncs
}

func NewGameSvc(chSrv chan message.NtolMessage) *GameService {
	gameSvc := &GameService{
		ChSrv:        chSrv,
		gameMsgFuncs: &GameMsgFuncs{},
	}

	gameSvc.gameMsgFuncs.RegisterAllFuns(gameSvc)

	return gameSvc
}

func (s *GameService) Run() {
	s.msgHandle()
}

func (s *GameService) Stop() {
	close(s.ChSrv)
}

func (s *GameService) msgHandle() {
	for {
		select {
		case msg, ok := <-s.ChSrv:
			if ok {
				s.gameMsgFuncs.Func(&msg)
			} else {
				logger.Error(titleGameService, constants.ErrGameServiceChannelClosed.Error())
				return
			}
		}
	}
}

func (s *GameService) RouteNext(ntolMsg *message.NtolMessage) {
	switch ntolMsg.MsgRoute.To {
	case constants.UnkowSvc:
		logger.Debug(titleGameService, "msg over now!id=%d,step=%d,rst=%t\n", ntolMsg.MsgRoute.MsgId, ntolMsg.MsgRoute.Step, ntolMsg.MsgRoute.Result)
	case constants.TCPEntrySvc:
		ntolMsg.SessionMsg.ChSend <- (*ntolMsg)
	default:
		logger.Debug(titleGameService, "msg over now!id=%d,step=%d,rst=%t\n", ntolMsg.MsgRoute.MsgId, ntolMsg.MsgRoute.Step, ntolMsg.MsgRoute.Result)
	}
}
