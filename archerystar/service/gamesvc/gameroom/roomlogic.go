package gameroom

import (
	"archerystar/component/logger"
	"archerystar/config"
	"time"
)

type RommLogic struct {
	Room         *GameRoom
	PlayerStatus []*PlayerStatus
	ticker       *time.Ticker
	frames       uint64
	State        ClassicState
}

func NewRoomLogic(room *GameRoom) *RommLogic {
	logic := &RommLogic{
		Room:         room,
		PlayerStatus: []*PlayerStatus{},
		ticker:       time.NewTicker(time.Millisecond * time.Duration(config.Gameconfig().Server.Room.Update)),
		frames:       0,
	}

	for i := 0; i < len(room.Players); i++ {
		logic.PlayerStatus = append(logic.PlayerStatus, &PlayerStatus{
			Rotation: Rotation{
				X: 0,
				Y: 0,
				Z: 0,
			},
			TotalScore: 0,
			Scores:     make([]float32, 0),
			IsBurst:    false,
			IsReady:    false,
		})
	}

	return logic
}

func (rl *RommLogic) Handle() {
	for {
		select {
		case data, ok := <-rl.Room.ChCmd:
			if ok {
				rl.CmdInput(&data)
			} else {
				logger.Debug(titleGameRoom, "room %d exit now!", rl.Room.Id)
				return
			}
		case <-rl.ticker.C:
			rl.CheckFrameAdd()
			rl.Update()
			rl.Multicast()
		}
	}
}

func (rl *RommLogic) CheckFrameAdd() {
	if rl.State != CLASSIC_STATE_WaitJoin {
		rl.frames++
	}
}

func (rl *RommLogic) CmdInput(data *StateCmdMsg) {

}

func (rl *RommLogic) Update() {

}

func (rl *RommLogic) Multicast() {

}
