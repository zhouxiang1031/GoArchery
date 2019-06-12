package message

import (
	"archerystar/common/constants"
)

type (
	MsgRoute struct {
		MsgId  uint
		Step   int
		Result bool
		To     constants.ServiceType
	}

	SessionMsg struct {
		Sid    int64 //session unique id
		SUid   string
		ChSend chan NtolMessage // push message queue
	}

	//local service -> net sock
	NtolMessage struct {
		//ctx     context.Context
		//session *session.Session
		SessionMsg *SessionMsg
		MsgRoute   *MsgRoute
		Msg        *Message
	}

	//net sock -> local service
	LtonMessage struct {
		//ctx     context.Context
		//typ message.Type // message type
		//step    int
		//route   string      // message route (push)
		//mid     uint        // response message id (response)
		//payload interface{} // payload
		//err     bool        // if its an error message
		SessionMsg *SessionMsg
		MsgRoute   *MsgRoute
		//session *session.Session
		//Msgs []*Message
		Msg *Message
	}
)
