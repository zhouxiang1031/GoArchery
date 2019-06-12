package session

import (
	"fmt"

	"archerystar/common/constants"
	"archerystar/component/logger"
	"archerystar/message"
	"archerystar/route"
	"archerystar/sockio/packet"
)

func (s *Session) readHandle() {
	//buf := make([]byte, MaxBuffSize)
	for {
		select {
		case <-s.chDie:
			return
		default:
			if !s.readFunc() {
				return
			}
		}
	}
}

func (s *Session) readClose() {
	logger.Debug(sessionTitle, "Session read goroutine exit, SessionID=%d, UID=%d", s.GetId(), s.GetUid())
}

func (s *Session) readFunc() bool {
	buf := make([]byte, MaxBuffSize)
	n, err := s.conn.Read(buf)
	if err != nil {
		logger.Error(sessionTitle, "Read message error: %s, session will be closed immediately,", err.Error())
		s.readClose()
		return false
	}

	logger.Debug(sessionTitle, "Received data on connection")

	// (warning): decoder uses slice for performance, packet data should be copied before next Decode
	packets, err := s.decoder.Decode(buf[:n])
	if err != nil {
		logger.Error(sessionTitle, "Failed to decode message: %s", err.Error())
		s.readClose()
		return false
	}

	if len(packets) < 1 {
		logger.Warn(sessionTitle, "Read no packets, data: %v", buf[:n])
		return true
	}

	// process all packet
	for i := range packets {
		fmt.Println(i)
		if err := s.packetHeaderProc(packets[i]); err != nil {
			logger.Error(sessionTitle, "Failed to process packet: %s", err.Error())
			return false
		}
	}

	return true
}

func (s *Session) packetHeaderProc(p *packet.Packet) error {
	switch p.Type {
	case packet.Data:
		if s.GetStatus() < constants.StatusWorking {
			return fmt.Errorf("receive data on socket which is not yet ACK, session will be closed immediately, remote=%s",
				s.RemoteAddr().String())
		}

		msg, err := message.Decode(p.Data)
		if err != nil {
			return err
		}
		msg.NetSeq = GetManagerInstance().NewNetSeq()
		//route  next now
		ntolMsg := &message.NtolMessage{
			//MsgRoute:   message.MsgRoute{},
			MsgRoute:   nil,
			SessionMsg: &message.SessionMsg{Sid: s.id, SUid: s.uid, ChSend: s.chSend},
			Msg:        msg,
		}
		route.Next(s, ntolMsg)

	case packet.Heartbeat:
		// expected
	}

	s.lastHeartBeat()
	return nil
}

func (s *Session) RouteNext(ntolMsg *message.NtolMessage) {
	switch ntolMsg.MsgRoute.To {
	case constants.UnkowSvc:
		logger.Error(sessionTitle, "msg over now!id=%d,step=%d,rst=%v", ntolMsg.MsgRoute.MsgId, ntolMsg.MsgRoute.Step, ntolMsg.MsgRoute.Result)
	case constants.GameSvc:
		s.chNetToGame <- (*ntolMsg)
	default:
		logger.Error(sessionTitle, "msg over now!id=%d,step=%d,rst=%v", ntolMsg.MsgRoute.MsgId, ntolMsg.MsgRoute.Step, ntolMsg.MsgRoute.Result)
	}
}
