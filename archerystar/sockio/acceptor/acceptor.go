package acceptor

import "net"

type Acceptor interface {
	ListenAndServe()
	Stop()
	GetAddr() string
	GetConnChan() chan net.Conn
}
