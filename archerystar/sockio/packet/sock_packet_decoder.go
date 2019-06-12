package packet

import (
	"bytes"
	"encoding/binary"
)

// SockPacketDecoder reads and decodes network data slice following pomelo's protocol
type SockPacketDecoder struct{}

// NewSockPacketDecoder returns a new decoder that used for decode network bytes slice.
func NewSockPacketDecoder() *SockPacketDecoder {
	return &SockPacketDecoder{}
}

func (c *SockPacketDecoder) forward(buf *bytes.Buffer) (int, Type, error) {
	header := buf.Next(HeadLength)
	typ := header[0]
	if typ < Handshake || typ > Kick {
		return 0, 0x00, ErrWrongSockPacketType
	}

	//size := bytesToInt(header[1:])
	size := int(binary.BigEndian.Uint16(header[1:]))

	// packet length limitation
	if size > MaxPacketSize {
		return 0, 0x00, ErrPacketSizeExcced
	}
	return size, Type(typ), nil
}

// Decode decode the network bytes slice to packet.Packet(s)
func (c *SockPacketDecoder) Decode(data []byte) ([]*Packet, error) {
	buf := bytes.NewBuffer(nil)
	buf.Write(data)

	var (
		packets []*Packet
		err     error
	)
	// check length
	if buf.Len() < HeadLength {
		return nil, nil
	}

	// first time
	size, typ, err := c.forward(buf)
	if err != nil {
		return nil, err
	}

	for size <= buf.Len() {
		p := &Packet{Type: typ, Length: size, Data: buf.Next(size)}
		packets = append(packets, p)

		// more packet
		if buf.Len() < HeadLength {
			break
		}

		size, typ, err = c.forward(buf)
		if err != nil {
			return nil, err
		}
	}

	return packets, nil
}

// Decode packet data length byte to int(Big end)
func bytesToInt(b []byte) int {
	result := 0
	for _, v := range b {
		result = result<<8 + int(v)
	}
	return result
}
