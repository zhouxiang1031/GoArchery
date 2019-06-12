package packet

type PacketDecoder interface {
	Decode(data []byte) ([]*Packet, error)
}
