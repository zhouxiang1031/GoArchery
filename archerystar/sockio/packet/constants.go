package packet

import "errors"

// Type represents the network packet's type such as: handshake and so on.
type Type byte

const (
	_ Type = iota
	// Handshake represents a handshake: request(client) <====> handshake response(server)
	Handshake = 0x01

	// HandshakeAck represents a handshake ack from client to server
	HandshakeAck = 0x02

	// Heartbeat represents a heartbeat
	Heartbeat = 0x03

	// Data represents a common data packet
	Data = 0x04

	// Kick represents a kick off packet
	Kick = 0x05 // disconnect message from server
)

const (
	HeadLength    = 3
	MaxPacketSize = 1 << 16 //65K
)

// ErrWrongPomeloPacketType represents a wrong packet type.
var ErrWrongSockPacketType = errors.New("wrong packet type")
var ErrPacketSizeExcced = errors.New("codec: packet size exceed")
