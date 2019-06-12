package entrysvc

import (
	"archerystar/config"
	"archerystar/message"
	"archerystar/sockio/acceptor"
	"archerystar/sockio/packet"
	"archerystar/sockio/session"
	"sync"
)

var once = sync.Once{}
var accInstance acceptor.Acceptor

type TcpEntryService struct {
	acceptor acceptor.Acceptor
}

func NewTcpEntrySvc(
	pDecoder packet.PacketDecoder,
	pEncoder packet.PacketEncoder,
	msgEncoder message.Encoder,
	chNet2Game chan message.NtolMessage,
) *TcpEntryService {
	tcpsvc := &TcpEntryService{}
	tcpsvc.acceptor = acceptor.NewTCPAcceptor(config.Gameconfig().Server.Addr, config.Gameconfig().Server.Certfile, config.Gameconfig().Server.Keyfile)

	session.GetManagerInstance().Init(pDecoder, pEncoder, msgEncoder, chNet2Game)

	return tcpsvc
}

func (s *TcpEntryService) Run() {
	s.acceptor.ListenAndServe()
}

func (s *TcpEntryService) Stop() {
	s.acceptor.Stop()
}
