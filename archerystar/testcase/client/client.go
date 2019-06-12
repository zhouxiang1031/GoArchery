package client

import (
	"archerystar/message"
	"archerystar/sockio/packet"
	"fmt"
	"net"

	"github.com/golang/protobuf/proto"
)

func SendRequest(conn net.Conn, msgId uint, pb proto.Message) bool {
	msg := &message.Message{}

	msg.Type = message.Request
	msg.ID = msgId
	msg.Compressed = false
	msg.Err = false

	var err error
	msg.Data, err = proto.Marshal(pb)
	if err != nil {
		fmt.Println("Marshaling error: ", err)
		return false
	}

	e, err := message.NewMessagesEncoder(false).Encode(msg)
	if err != nil {
		fmt.Printf("Failed to encode message: %s/n", err.Error())
		return false
	}

	// packet encode
	p, err := packet.NewSockPacketEncoder().Encode(packet.Data, e)
	if err != nil {
		fmt.Printf("Failed to encode packet: %s\n", err.Error())
		return false
	}

	if len(p) > 2048 {
		fmt.Printf("packet too large igore it\n")
		return true
	}

	if _, err := conn.Write(p); err != nil {
		fmt.Printf("Failed to write response: %s\n", err.Error())
		return false
	}
	return true
}

func ReadResp(conn net.Conn, pb proto.Message) bool {
	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("Read message error: %s, session will be closed immediately\n", err.Error())
		return false
	}

	fmt.Printf("recv from server:%v\n", buf[:n])
	//fmt.Printf("Received data on connection")

	// (warning): decoder uses slice for performance, packet data should be copied before next Decode
	p, err := packet.NewSockPacketDecoder().Decode(buf[:n])
	if err != nil {
		fmt.Printf("Failed to decode message: %s\n", err.Error())
		return false
	}

	if len(p) < 1 {
		fmt.Printf("Read no packets, data: %v\n", buf[:n])
		return true
	}

	msg, err := message.Decode(p[0].Data)
	if err != nil {
		fmt.Printf("Read no packets failed,err: %v\n", err.Error())
		return false
	}

	proto.Unmarshal(msg.Data, pb)

	fmt.Printf("test resp:%+v\n", pb)

	return true
}
