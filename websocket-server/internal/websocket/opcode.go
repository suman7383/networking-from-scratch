package websocket

import "fmt"

type Opcode uint8

const (
	OpContinuation Opcode = 0x0
	OpText         Opcode = 0x1
	OpBinary       Opcode = 0x2
	OpClose        Opcode = 0x8
	OpPing         Opcode = 0x9
	OpPong         Opcode = 0xA
)

func (oc Opcode) IsControlFrame() bool {
	return oc == OpClose || oc == OpPing || oc == OpPong
}

func (oc Opcode) String() string {
	switch oc {
	case 0x0:
		return "CONTINUATION"
	case 0x1:
		return "TEXT"
	case 0x2:
		return "BINARY"
	case 0x8:
		return "CLOSE"
	case 0x9:
		return "PING"
	case 0xA:
		return "PONG"
	default:
		return fmt.Sprintf("UNKNOWN(0x%x)", uint8(oc))
	}
}
