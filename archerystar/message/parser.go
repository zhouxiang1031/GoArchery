package message

import (
	"encoding/binary"

	"archerystar/common/compression"
)

// Encoder interface
type Encoder interface {
	IsCompressionEnabled() bool
	Encode(message *Message) ([]byte, error)
}

// MessagesEncoder implements MessageEncoder interface
type MessagesEncoder struct {
	DataCompression bool
}

// NewMessagesEncoder returns a new message encoder
func NewMessagesEncoder(dataCompression bool) *MessagesEncoder {
	me := &MessagesEncoder{dataCompression}
	return me
}

// IsCompressionEnabled returns wether the compression is enabled or not
func (me *MessagesEncoder) IsCompressionEnabled() bool {
	return me.DataCompression
}

// Encode marshals message to binary format. Different message types is corresponding to
// different message header, message types is identified by 2-4 bit of flag field. The
// relationship between message types and message header is presented as follows:
// ------------------------------------------
// |   type   |  flag  |       other        |
// |----------|--------|--------------------|
// | request  |----000-|<message id>|<route>|
// | notify   |----001-|<route>             |
// | response |----010-|<message id>        |
// | push     |----011-|<route>             |
// ------------------------------------------
// The figure above indicates that the bit does not affect the type of message.
// See ref: https://github.com/topfreegames/pitaya/blob/master/docs/communication_protocol.md
func (me *MessagesEncoder) Encode(message *Message) ([]byte, error) {
	if invalidType(message.Type) {
		return nil, ErrWrongMessageType
	}

	buf := make([]byte, 0)
	flag := byte(message.Type) << 1

	if message.Err {
		flag |= errorMask
	}

	buf = append(buf, flag)

	if message.Type == Request || message.Type == Response {

		mid := make([]byte, 2)
		binary.BigEndian.PutUint16(mid, uint16(message.ID))
		buf = append(buf, mid...)
	}

	if me.DataCompression {
		d, err := compression.DeflateData(message.Data)
		if err != nil {
			return nil, err
		}

		if len(d) < len(message.Data) {
			message.Data = d
			buf[0] |= gzipMask
		}
	}

	buf = append(buf, message.Data...)
	return buf, nil
}

// Decode decodes the message
func (me *MessagesEncoder) Decode(data []byte) (*Message, error) {
	return Decode(data)
}

// Decode unmarshal the bytes slice to a message
func Decode(data []byte) (*Message, error) {
	if len(data) < msgHeadLength {
		return nil, ErrInvalidMessage
	}
	m := New()
	flag := data[0]
	offset := 1
	m.Type = Type((flag >> 1) & msgTypeMask)

	if invalidType(m.Type) {
		return nil, ErrWrongMessageType
	}

	m.Err = flag&errorMask == errorMask

	if m.Type == Request || m.Type == Response {

		m.ID = uint(binary.BigEndian.Uint16(data[offset : offset+2]))
		offset += 2
	}

	m.Data = data[offset:]
	var err error
	if flag&gzipMask == gzipMask {
		m.Data, err = compression.InflateData(m.Data)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}
