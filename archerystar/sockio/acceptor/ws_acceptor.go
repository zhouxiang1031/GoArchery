package acceptor

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"archerystar/common/constants"
	"archerystar/component/logger"
)

type WSAcceptor struct {
	addr     string
	connChan chan net.Conn
	listener net.Listener
	certFile string
	keyFile  string
}

const WsAcceptTitle = "tcp_acceptor"

func NewWSAcceptor(addr string, certs ...string) *WSAcceptor {
	keyFile := ""
	certFile := ""
	if len(certs) != 2 && len(certs) != 0 {
		logger.Panic(WsAcceptTitle, constants.ErrInvalidCert.Error())
	} else if len(certs) == 2 {
		certFile = certs[0]
		keyFile = certs[1]
	}

	w := &WSAcceptor{
		addr:     addr,
		connChan: make(chan net.Conn),
		certFile: certFile,
		keyFile:  keyFile,
	}
	return w
}

func (w *WSAcceptor) GetAddr() string {
	if w.listener != nil {
		return w.listener.Addr().String()
	}
	return ""
}

func (w *WSAcceptor) GetConnChan() chan net.Conn {
	return w.connChan
}

type connHandler struct {
	upgrader *websocket.Upgrader
	connChan chan net.Conn
}

func (h *connHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(rw, r, nil)
	if err != nil {
		logger.Error(WsAcceptTitle, "Upgrade failure, URI=%s, Error=%s", r.RequestURI, err.Error())
		return
	}

	c, err := newWSConn(conn)
	if err != nil {
		logger.Error(WsAcceptTitle, "Failed to create new ws connection: %s", err.Error())
		return
	}
	h.connChan <- c
}

func (w *WSAcceptor) useTLSCert() bool {
	return w.certFile != "" && w.keyFile != ""
}

func (w *WSAcceptor) ListenAndServe() {
	if w.useTLSCert() {
		w.ListenAndServeTLS(w.certFile, w.keyFile)
		return
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	listener, err := net.Listen("tcp", w.addr)
	if err != nil {
		logger.Fatal(WsAcceptTitle, "Failed to listen: %s", err.Error())
	}
	w.listener = listener

	w.serve(&upgrader)
}

func (w *WSAcceptor) ListenAndServeTLS(cert, key string) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	crt, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		logger.Fatal(WsAcceptTitle, "Failed to load x509: %s", err.Error())
	}

	tlsCfg := &tls.Config{Certificates: []tls.Certificate{crt}}
	listener, err := tls.Listen("tcp", w.addr, tlsCfg)
	if err != nil {
		logger.Fatal(WsAcceptTitle, "Failed to listen: %s", err.Error())
	}
	w.listener = listener
	w.serve(&upgrader)
}

func (w *WSAcceptor) serve(upgrader *websocket.Upgrader) {
	defer w.Stop()

	http.Serve(w.listener, &connHandler{
		upgrader: upgrader,
		connChan: w.connChan,
	})
}

func (w *WSAcceptor) Stop() {
	err := w.listener.Close()
	if err != nil {
		logger.Error(WsAcceptTitle, "Failed to stop: %s", err.Error())
	}
}

// interface base on *websocket.Conn
type wsConn struct {
	conn   *websocket.Conn
	typ    int // message type
	reader io.Reader
}

func newWSConn(conn *websocket.Conn) (*wsConn, error) {
	c := &wsConn{conn: conn}

	t, r, err := conn.NextReader()
	if err != nil {
		return nil, err
	}

	c.typ = t
	c.reader = r

	return c, nil
}

func (c *wsConn) Read(b []byte) (int, error) {
	n, err := c.reader.Read(b)
	if err != nil && err != io.EOF {
		return n, err
	} else if err == io.EOF {
		_, r, err := c.conn.NextReader()
		if err != nil {
			return 0, err
		}
		c.reader = r
	}

	return n, nil
}

func (c *wsConn) Write(b []byte) (int, error) {
	err := c.conn.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		return 0, err
	}

	return len(b), nil
}

func (c *wsConn) Close() error {
	return c.conn.Close()
}

// LocalAddr returns the local network address.
func (c *wsConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *wsConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail with a timeout (see type Error) instead of
// blocking. The deadline applies to all future and pending
// I/O, not just the immediately following call to Read or
// Write. After a deadline has been exceeded, the connection
// can be refreshed by setting a deadline in the future.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
func (c *wsConn) SetDeadline(t time.Time) error {
	if err := c.SetReadDeadline(t); err != nil {
		return err
	}

	return c.SetWriteDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (c *wsConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (c *wsConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
