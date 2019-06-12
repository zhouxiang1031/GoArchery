package session

import (
	"time"

	"archerystar/component/logger"
	"archerystar/config"
	"archerystar/sockio/packet"
)

func (s *Session) heartbeat() {
	ticker := time.NewTicker(time.Duration(config.Gameconfig().Server.Heartbeat.Interval) * time.Second)

	stop := func() {
		ticker.Stop()
		s.Close()
	}

	for {
		select {
		case <-ticker.C:
			deadline := time.Now().Add(time.Second * (-s.heartbeatTimeout)).Unix()
			if s.lastHeartTime < deadline {
				logger.Debug(sessionTitle, "Session heartbeat timeout, LastTime=%v, Deadline=%v", s.lastHeartTime, deadline)
				stop()
				return
			}

			hbd := []byte{byte(packet.Heartbeat)}
			if _, err := s.conn.Write(hbd); err != nil {
				stop()
				return
			}
		case <-s.chDie:
			return
		}
	}
}

func (s *Session) lastHeartBeat() {
	s.lastHeartTime = time.Now().Unix()
}
