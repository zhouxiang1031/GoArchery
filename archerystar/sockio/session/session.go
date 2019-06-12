package session

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"archerystar/common/constants"
	"archerystar/common/errors"
	"archerystar/component/logger"
	"archerystar/config"
	"archerystar/message"
	"archerystar/sockio/packet"
)

type (
	Session struct {
		//base of seeion
		sync.RWMutex
		id    int64  //session unique id in game
		uid   string //bind user id
		state int32  //current session state

		//net
		conn             net.Conn //connection fd
		OnCloseCallbacks []func() //onClose callbacks
		closeMutex       sync.Mutex

		//channel
		appDieChan  chan bool                // app die channel
		chDie       chan struct{}            // wait for close
		chSend      chan message.NtolMessage // push message queue
		chNetToGame chan message.NtolMessage //send message to game service

		//parse for message
		decoder        packet.PacketDecoder // binary decoder
		encoder        packet.PacketEncoder // binary encoder
		messageEncoder message.Encoder
		//serializer         serialize.Serializer // message serializer

		//hearbeat
		heartbeatTimeout time.Duration
		lastHeartTime    int64 //last heartbeat time

		//data
		data map[string]interface{} // session data store
		//handshakeData     *HandshakeData         // handshake data received by the client
		encodedData []byte // session data encoded as a byte array

	}

	pendingMessage struct {
		ctx context.Context
		typ message.Type // message type
		//route   string       // message route (push)
		mid     uint        // response message id (response)
		payload interface{} // payload
		err     bool        // if its an error message
		netSeq  uint64
	}
)

const sessionTitle = "session"

func New(
	sid int64,
	conn net.Conn,
	packetDecoder packet.PacketDecoder,
	packetEncoder packet.PacketEncoder,
	messageEncoder message.Encoder,
	chNetToGame chan message.NtolMessage,
) *Session {
	// initialize heartbeat and handshake data on first user connection
	/*
		once.Do(func() {
			//InitHeartbeatAndHandshake(heartbeatTime, packetEncoder, messageEncoder.IsCompressionEnabled(), serializer.GetName())
		})
	*/

	s := &Session{
		id:               sid,
		chDie:            make(chan struct{}),
		chSend:           make(chan message.NtolMessage, config.Gameconfig().Server.BuffLimit.NtogMax),
		chNetToGame:      chNetToGame,
		conn:             conn,
		decoder:          packetDecoder,
		encoder:          packetEncoder,
		heartbeatTimeout: time.Duration(config.Gameconfig().Server.Heartbeat.Ttl * config.Gameconfig().Server.Heartbeat.Interval),
		lastHeartTime:    time.Now().Unix(),
		state:            constants.StatusWorking,
		messageEncoder:   messageEncoder,

		data:             make(map[string]interface{}),
		OnCloseCallbacks: []func(){},
	}

	return s
}

/*
func InitHeartbeatAndHandshake(heartbeatTimeout time.Duration, packetEncoder codec.PacketEncoder, dataCompression bool, serializerName string) {
	hData := map[string]interface{}{
		"code": 200,
		"sys": map[string]interface{}{
			"heartbeat": heartbeatTimeout.Seconds(),
			//"dict":       message.GetDictionary(),
			//"serializer": serializerName,
		},
	}
	data, err := gojson.Marshal(hData)
	if err != nil {
		panic(err)
	}

	if dataCompression {
		compressedData, err := compression.DeflateData(data)
		if err != nil {
			panic(err)
		}

		if len(compressedData) < len(data) {
			data = compressedData
		}
	}

	hrd, err = packetEncoder.Encode(packet.Handshake, data)
	if err != nil {
		panic(err)
	}

	hbd, err = packetEncoder.Encode(packet.Heartbeat, nil)
	if err != nil {
		panic(err)
	}
}
*/

func (s *Session) GetStatus() int32 {
	return atomic.LoadInt32(&s.state)
}

func (s *Session) GetId() int64 {
	return s.id
}

func (s *Session) GetUid() string {
	return s.uid
}

// GetData gets the data
func (s *Session) GetData() map[string]interface{} {
	s.RLock()
	defer s.RUnlock()

	return s.data
}

func (s *Session) send(m message.NtolMessage) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.NewError(constants.ErrBrokenPipe, errors.ErrClientClosedRequest)
		}
	}()

	s.chSend <- m
	return
}

func (s *Session) SetStatus(state int32) {
	atomic.StoreInt32(&s.state, state)
}

func (s *Session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *Session) onSessionClosed() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(sessionTitle, "onSessionClosed: %v", err)
		}
	}()

	for _, fn1 := range s.OnCloseCallbacks {
		fn1()
	}

	/*
		for _, fn2 := range session.SessionCloseCallbacks {
			fn2(s)
		}
	*/
}

func (s *Session) OnClose(c func()) error {
	/*
		if !s.IsFrontend {
			return constants.ErrOnCloseBackend
		}
	*/

	s.OnCloseCallbacks = append(s.OnCloseCallbacks, c)
	return nil
}

func (s *Session) String() string {
	return fmt.Sprintf("Remote=%s, LastTime=%d", s.conn.RemoteAddr().String(), s.lastHeartTime)
}

func (s *Session) Handle() {
	logger.Debug(sessionTitle, "New session established: %s", s.String())

	go s.readHandle()
	go s.writeHandle()
	go s.heartbeat()
}
