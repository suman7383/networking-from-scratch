package websocket

type Frame struct {
	Fin        bool
	Opcode     Opcode
	Masked     bool
	MaskKey    [4]byte
	PayloadLen uint16 // Support 2^16
	Payload    []byte
}
