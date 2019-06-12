package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"archerystar/component/logger"
	"archerystar/config"
	"archerystar/message"
	"archerystar/service/entrysvc"
	"archerystar/service/gamesvc"
	"archerystar/sockio/packet"
)

type ServerStatus = uint8

const titleGame = "archerystar"

const (
	Stoped ServerStatus = iota
	Running
	Prepare
	maintain
)

type ArcheryStar struct {
	entrySvc entrysvc.EntrySvc
	gameSvc  *gamesvc.GameService

	startAt     time.Time
	dieChan     chan bool
	gameSvcChan chan message.NtolMessage
	release     bool
	status      ServerStatus //0:stop,1:running,2:prepare,3.maintain

}

func NewGame() *ArcheryStar {
	g := &ArcheryStar{
		//entrySvc: entrysvc.NewTcpEntrySvc(),
		//gamesvcSrvs: gamesvc.NewGameSvc(),
		startAt:     time.Now(),
		release:     config.Gameconfig().Server.Release,
		status:      Prepare,
		dieChan:     make(chan bool),
		gameSvcChan: make(chan message.NtolMessage),
	}

	//if config.Gameconfig().Server.NetType == constants.Tcp
	g.entrySvc = entrysvc.NewTcpEntrySvc(
		packet.NewSockPacketDecoder(),
		packet.NewSockPacketEncoder(),
		message.NewMessagesEncoder(config.Gameconfig().Server.IsCompression),
		g.gameSvcChan,
	)

	g.gameSvc = gamesvc.NewGameSvc(g.gameSvcChan)

	return g
}

func main() {
	fmt.Println("archerystar server init now!")
	game := NewGame()

	go game.entrySvc.Run()
	go game.gameSvc.Run()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)

	//waiting for game stop..
	select {
	case <-game.dieChan:
		logger.Warn(titleGame, "game will shutdown in a few seconds!")
	case s := <-sig:
		logger.Warn(titleGame, "got signal: %s  shutting down...", s.String())
		close(game.dieChan)
	}

	logger.Warn(titleGame, "server is stopping...")
	game.entrySvc.Stop()
	game.gameSvc.Stop()
}
