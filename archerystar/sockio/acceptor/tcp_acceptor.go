package acceptor

import (
	"crypto/tls"
	"net"

	"archerystar/common/constants"
	"archerystar/component/logger"
	"archerystar/sockio/session"
)

type TCPAcceptor struct {
	addr string
	//connChan chan net.Conn
	listener net.Listener
	isRun    bool
	certFile string
	keyFile  string
}

const TcpAcceptTitle = "tcp_acceptor"

func NewTCPAcceptor(addr string, certs ...string) *TCPAcceptor {
	keyFile := ""
	certFile := ""
	if len(certs) != 2 && len(certs) != 0 {
		logger.Panic(TcpAcceptTitle, constants.ErrInvalidCert.Error())
	} else if len(certs) == 2 {
		certFile = certs[0]
		keyFile = certs[1]
	}

	return &TCPAcceptor{
		addr:     addr,
		isRun:    false,
		certFile: certFile,
		keyFile:  keyFile,
	}
}

func (a *TCPAcceptor) GetAddr() string {
	if a.listener != nil {
		return a.listener.Addr().String()
	} else {
		return ""
	}
}

func (a *TCPAcceptor) GetConnChan() chan net.Conn {
	//return a.connChan
	return nil
}

// Stop stops the acceptor
func (a *TCPAcceptor) Stop() {
	a.isRun = false
	a.listener.Close()

	session.GetManagerInstance().CloseAllSessions()
}

func (a *TCPAcceptor) useTLSCert() bool {
	return a.certFile != "" && a.keyFile != ""
}

func (a *TCPAcceptor) ListenAndServe() {
	if !a.useTLSCert() {
		a.Listen()
	} else {
		a.ListenWithTLS()
	}

	a.isRun = true
	a.serve()
}

func (a *TCPAcceptor) Listen() {
	var err error
	a.listener, err = net.Listen("tcp", a.addr)
	if err != nil {
		logger.Fatal(TcpAcceptTitle, "Failed to listen: %s", err.Error())
	}
}

func (a *TCPAcceptor) ListenWithTLS() {
	crt, err := tls.LoadX509KeyPair(a.certFile, a.keyFile)
	if err != nil {
		logger.Fatal(TcpAcceptTitle, "Failed to listen: %s", err.Error())
	}

	tlsCfg := &tls.Config{Certificates: []tls.Certificate{crt}}

	a.listener, err = tls.Listen("tcp", a.addr, tlsCfg)
	if err != nil {
		logger.Fatal(TcpAcceptTitle, "Failed to tls listen: %s", err.Error())
	}
}

func (a *TCPAcceptor) serve() {
	//defer a.Stop()
	for a.isRun {
		conn, err := a.listener.Accept()
		if err != nil {
			logger.Error(TcpAcceptTitle, "Failed to accept TCP connection: %s", err.Error())
			continue
		}
		session.GetManagerInstance().CreateSession(conn)
	}
}
