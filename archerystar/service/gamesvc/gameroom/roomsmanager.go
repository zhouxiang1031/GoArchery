package gameroom

import (
	"archerystar/config"
	"archerystar/message"
	"archerystar/protocol/protos"
	"sync"
	"sync/atomic"
	"time"
)

var once = sync.Once{}

type RoomsManager struct {
	GRoomId        uint64
	GameRoomMap    GameRoomMap
	SessionRoomMap SessionRoomMap
}

var roomsManager *RoomsManager

func GetRoomsManager() *RoomsManager {
	once.Do(func() {
		roomsManager = &RoomsManager{
			GameRoomMap:    GameRoomMap{},
			SessionRoomMap: SessionRoomMap{},
		}
	})

	return roomsManager
}

func (m *RoomsManager) CreateRomm(roomData *protos.GameRoomReq) *GameRoom {

	room := &GameRoom{
		Id:      m.NewRoomID(),
		Mode:    roomData.Mode,
		Stage:   roomData.Stage,
		StartAt: time.Now().UnixNano() / 1e6,
		Status:  0,
		Session: []*message.SessionMsg{},
		ChCmd:   make(chan StateCmdMsg, config.Gameconfig().Server.Room.Buff),
	}

	for i := 0; i < len(roomData.Players); i++ {
		player := &PlayerInfo{
			Id:       roomData.Players[i].Id,
			Account:  roomData.Players[i].Account,
			Name:     roomData.Players[i].Name,
			Country:  roomData.Players[i].Country,
			Trophies: roomData.Players[i].Trophies,
			Equips:   roomData.Players[i].Equips[:len(roomData.Players[i].Equips)],
			Head:     roomData.Players[i].Head,
			Order:    roomData.Players[i].Order,
		}

		room.Players = append(room.Players, player)
	}

	m.GameRoomMap[room.Id] = room

	return room
}

func (m *RoomsManager) NewRoomID() uint64 {
	return atomic.AddUint64(&m.GRoomId, 1)
}

func (m *RoomsManager) RemoveRoom(id uint64) {
	if room, ok := m.GameRoomMap[id]; ok {
		for i := 0; i < len(room.Session); i++ {
			delete(m.SessionRoomMap, room.Session[i].Sid)
		}

		delete(m.GameRoomMap, id)
	}
}

func (m *RoomsManager) FindRoom(id uint64) (room *GameRoom, ok bool) {
	room, ok = m.GameRoomMap[id]
	return
}
