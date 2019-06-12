package gameroom

type CMD = uint32
type ClassicState = uint32

const (
	_ CMD = iota
	Cmd_JoinRoom
)

const (
	_ ClassicState = iota
	CLASSIC_STATE_WaitJoin
)
