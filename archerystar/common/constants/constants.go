package constants

const (
	_ int32 = iota
	// StatusStart status
	StatusStart
	// StatusHandshake status
	StatusHandshake
	// StatusWorking status
	StatusWorking
	// StatusClosed status
	StatusClosed
)

type propagateKey struct{}

// PropagateCtxKey is the context key where the content that will be
var PropagateCtxKey = propagateKey{}

// propagated through rpc calls is set

const (
	Tcp = iota
	Websocket
	Udp
)

type ServiceType = int

const (
	UnkowSvc ServiceType = iota
	TCPEntrySvc
	GameSvc
)
