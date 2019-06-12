package gameroom

import (
	"archerystar/message"
)

type PlayerInfo struct {
	Id       int32
	Account  string
	Name     string
	Country  int32
	Trophies int32
	Equips   []int32
	Head     int32
	Order    int32
}

type Rotation struct {
	X float32
	Y float32
	Z float32
}

type Positin struct {
	X float32
	Y float32
	Z float32
}

type TragetStatus struct {
	TargetType int32
	Index      int32
	Pos        Positin
}

type PlayerStatus struct {
	Rotation   Rotation
	TotalScore int32
	Scores     []float32
	IsBurst    bool
	IsReady    bool
}

type StateCmdMsg struct {
	Cmd          CMD
	PlayerStatus *PlayerStatus
}

type GameRoom struct {
	Id       uint64
	Mode     int32
	Stage    int32
	StartAt  int64
	Status   int16
	CurSet   int32
	CurRound int32
	ChCmd    chan StateCmdMsg
	Players  []*PlayerInfo
	Session  []*message.SessionMsg
}

type SessionRoomMap map[int64]int64
type PlayerRoomMap map[string]int64
type GameRoomMap map[uint64]*GameRoom

const titleGameRoom = "gameroom"
