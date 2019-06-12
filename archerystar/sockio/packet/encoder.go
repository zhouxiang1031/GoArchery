package packet

type PacketEncoder interface {
	Encode(typ Type, data []byte) ([]byte, error)
}
