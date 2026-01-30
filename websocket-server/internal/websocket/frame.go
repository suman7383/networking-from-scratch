package websocket

type Frame struct {
	Fin        bool
	Opcode     Opcode
	Masked     bool
	MaskKey    [4]byte
	PayloadLen uint16 // Support 2^16
	Payload    []byte
}

// Does shallow clonning of the FRAME
func (f *Frame) Clone() *Frame {
	frW := new(Frame)
	// Copy the FRAME
	*frW = *f
	return frW
}
