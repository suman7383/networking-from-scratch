package websocket

type Frame struct {
	Fin     bool
	Opcode  Opcode
	Masked  bool
	MaskKey [4]byte
	Payload []byte
}
