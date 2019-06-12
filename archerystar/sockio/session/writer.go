package session

import (
	"archerystar/component/logger"
	"archerystar/message"
	"archerystar/sockio/packet"
)

func (s *Session) writeHandle() {
	for {
		select {
		case msgs, ok := <-s.chSend:
			if !ok || !s.writeFunc(msgs) {
				s.closeWrite()
				return
			}
		case <-s.chDie:
			return
		}
	}
}

func (s *Session) writeFunc(ltonmsg message.NtolMessage) bool {
	e, err := s.messageEncoder.Encode(ltonmsg.Msg)
	if err != nil {
		logger.Error(sessionTitle, "Failed to encode message: %s", err.Error())
		return false
	}

	// packet encode
	p, err := s.encoder.Encode(packet.Data, e)
	if err != nil {
		logger.Error(sessionTitle, "Failed to encode packet: %s", err.Error())
		return false
	}

	if len(p) > MaxBuffSize {
		logger.Error(sessionTitle, "packet too large igore it")
		return true
	}

	if _, err := s.conn.Write(p); err != nil {
		logger.Error("Failed to write response: %s", err.Error())
		return false
	}
	return true
}

/*
func (s *Session) writeBatchFunc(ltonmsgs message.LtonMessage) bool {
	buf := make([]byte, 0)

	for _, msg := range ltonmsgs.Msgs {
		// construct message and encode
		e, err := s.messageEncoder.Encode(msg)
		if err != nil {
			logger.Error(sessionTitle, "Failed to encode message: %s", err.Error())
			return false
		}

		// packet encode
		p, err := s.encoder.Encode(packet.Data, e)
		if err != nil {
			logger.Error(sessionTitle, "Failed to encode packet: %s", err.Error())
			return false
		}

		if len(p) > MaxBuffSize {
			logger.Error(sessionTitle, "packet too large igore it")
			continue
		}

		if len(buf)+len(p) <= MaxBuffSize {
			buf = append(buf, p[:]...)
		} else {
			logger.Debug(sessionTitle, "packet to large,split it!")
			// close session if low-level Conn broken
			if _, err := s.conn.Write(buf); err != nil {
				logger.Error("Failed to write response: %s", err.Error())
				return false
			}
			buf = p[:]
		}
	}

	if len(buf) > 0 {
		// close session if low-level Conn broken
		if _, err := s.conn.Write(buf); err != nil {
			logger.Error("Failed to write response: %s", err.Error())
			return false
		}
	}

	return true
}
*/

func (s *Session) closeWrite() {
	logger.Debug(sessionTitle, "Session write goroutine exit, SessionID=%d, UID=%d", s.GetId(), s.GetUid())
	close(s.chSend)
	s.Close()
}
