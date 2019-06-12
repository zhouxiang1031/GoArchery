package session

import (
	"net"
	"sync"
	"sync/atomic"

	"archerystar/component/logger"
	"archerystar/message"
	"archerystar/sockio/packet"
)

var once sync.Once

type SessionManger struct {
	//session info
	gSessionId    int64
	sessionsByID  sync.Map
	sessionsByUID sync.Map
	SessionCount  int64

	//message info
	packetDecoder  packet.PacketDecoder
	packetEncoder  packet.PacketEncoder
	messageEncoder message.Encoder

	//channel
	chNetToGame chan message.NtolMessage

	gNetSeq uint64
}

//var SessionManager = newSessionManger()

var sessionManager *SessionManger

func GetManagerInstance() *SessionManger {
	once.Do(func() {
		sessionManager = &SessionManger{
			gSessionId:   0,
			SessionCount: 0,
			gNetSeq:      0,
		}
	})

	return sessionManager
}

func (mng *SessionManger) Init(
	pDecoder packet.PacketDecoder,
	pEncoder packet.PacketEncoder,
	msgEncoder message.Encoder,
	chNet2Game chan message.NtolMessage,
) {
	mng.packetDecoder = pDecoder
	mng.packetEncoder = pEncoder
	mng.messageEncoder = msgEncoder
	mng.chNetToGame = chNet2Game
}

func (mng *SessionManger) NewSessionID() int64 {
	return atomic.AddInt64(&mng.gSessionId, 1)
}

func (mng *SessionManger) CreateSession(
	conn net.Conn,
) {
	s := New(mng.NewSessionID(), conn, mng.packetDecoder, mng.packetEncoder, mng.messageEncoder, mng.chNetToGame)

	mng.AddSession(s)
	s.Handle()
}

func (mng *SessionManger) CloseAllSessions() {
	logger.Debug(sessionTitle, "closing all sessions, %d sessions", mng.SessionCount)
	mng.sessionsByID.Range(func(_, value interface{}) bool {
		s := value.(*Session)
		s.Close()
		return true
	})
	logger.Debug(sessionTitle, "finished closing sessions")
}

func (mng *SessionManger) AddSession(s *Session) {
	mng.sessionsByID.Store(s.id, s)
	atomic.AddInt64(&mng.SessionCount, 1)
	if s.GetUid() != "" {
		mng.sessionsByID.Store(s.uid, s)
	}
}

func (mng *SessionManger) RemoveSession(s *Session) {
	atomic.AddInt64(&mng.SessionCount, -1)
	mng.sessionsByID.Delete(s.id)
	if s.GetUid() != "" {
		mng.sessionsByUID.Delete(s.uid)
	}

}

func (mng *SessionManger) NewNetSeq() uint64 {
	return atomic.AddUint64(&mng.gNetSeq, 1)
}
