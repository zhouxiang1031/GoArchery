package constants

import "errors"

var (
	ErrInvalidCert              = errors.New("certificates must be exactly two")
	ErrBrokenPipe               = errors.New("broken low-level pipe")
	ErrSessionOnNotify          = errors.New("current session working on notify mode")
	ErrCloseClosedSession       = errors.New("the session has been closed")
	ErrGameServiceChannelClosed = errors.New("the game service channel has been closed")
)
