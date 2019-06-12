package session

import (
	"context"

	"archerystar/common/constants"
	"archerystar/common/errors"
	"archerystar/component/logger"
	"archerystar/message"
	"archerystar/sockio/packet"
)

func (s *Session) Push(route string, v interface{}) error {
	if s.GetStatus() == constants.StatusClosed {
		return errors.NewError(constants.ErrBrokenPipe, errors.ErrClientClosedRequest)
	}

	switch d := v.(type) {
	case []byte:
		logger.Debug(sessionTitle, "Type=Push, ID=%d, UID=%d, Route=%s, Data=%dbytes",
			s.GetId(), s.GetUid(), route, len(d))
	default:
		logger.Debug(sessionTitle, "Type=Push, ID=%d, UID=%d, Route=%s, Data=%+v",
			s.GetId(), s.GetUid(), route, v)
	}
	//return s.send(pendingMessage{typ: message.Push, route: route, payload: v})
	return s.send(message.NtolMessage{})

}

func (s *Session) ResponseMID(ctx context.Context, mid uint, v interface{}, isError ...bool) error {
	return nil
	/*
		err := false
		if len(isError) > 0 {
			err = isError[0]
		}
		if s.GetStatus() == constants.StatusClosed {
			err := errors.NewError(constants.ErrBrokenPipe, errors.ErrClientClosedRequest)
			return err
		}

		if mid <= 0 {
			err := constants.ErrSessionOnNotify
			return err
		}

		switch d := v.(type) {
		case []byte:
			logger.Debug(sessionTitle, "Type=Response, ID=%d, UID=%d, MID=%d, Data=%dbytes",
				s.GetId(), s.GetUid(), mid, len(d))
		default:
			logger.Debug(sessionTitle, "Type=Response, ID=%d, UID=%d, MID=%d, Data=%+v",
				s.GetId(), s.GetUid(), mid, v)
		}

		return s.send(message.LtonMessage{})
	*/
}

// Close closes the agent, cleans inner state and closes low-level connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (s *Session) Close() error {
	s.closeMutex.Lock()
	defer s.closeMutex.Unlock()
	if s.GetStatus() == constants.StatusClosed {
		return constants.ErrCloseClosedSession
	}
	s.SetStatus(constants.StatusClosed)

	logger.Debug(sessionTitle, "Session closed, ID=%d, UID=%s, IP=%s",
		s.GetId(), s.GetUid(), s.conn.RemoteAddr())

	GetManagerInstance().RemoveSession(s)

	// prevent closing closed channel
	select {
	case <-s.chDie:
		// expect
	default:
		close(s.chDie)
		s.onSessionClosed()
	}

	return s.conn.Close()
}

func (s *Session) Kick(ctx context.Context) error {
	// packet encode
	p, err := s.encoder.Encode(packet.Kick, nil)
	if err != nil {
		return err
	}
	_, err = s.conn.Write(p)
	return err
}
