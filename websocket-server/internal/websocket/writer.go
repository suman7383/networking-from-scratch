package websocket

import (
	"bufio"
	"encoding/binary"
	"net"

	"github.com/suman7383/networking-from-scratch/websocket-server/utils"
)

type FrameWriter struct {
	w       *bufio.Writer
	closeCh chan struct{}
}

func NewFrameWriter(conn net.Conn) *FrameWriter {
	return &FrameWriter{
		w:       bufio.NewWriter(conn),
		closeCh: make(chan struct{}),
	}
}

func (fw *FrameWriter) WriteFrame(f *Frame) error {
	// Write 2 bytes
	//
	// FIN(1 bit): 1, RSV(3 bit): 0
	// OPCODE(4 bits), MASK(1 bit): 0
	// BASE PAYLOAD(7 bits)
	// payload length:
	// 0-125: payload length
	//
	// payload length
	// 0-125: payload length
	// 126: next 2 bytes = actual length
	// 127: next 8 bytes = actual length (not-supported here)
	info := make([]byte, 2)

	info[0] = (1 << 7)        // FIN = 1 (bit 7)
	info[0] |= byte(f.Opcode) // OPCODE in bits 0-3

	// payload len
	//
	// 0 - 125: payload length
	if f.PayloadLen < 126 {
		info[1] = byte(f.PayloadLen)
	} else {
		// 126: next 2 bytes = actual length
		info[1] = 126
	}

	// Write the first 2 bytes
	err := fw.write(info)
	if err != nil {
		utils.LogErr("could not write info bytes to conn", err)

		return err
	}

	// EXT payload len
	if f.PayloadLen >= 126 {
		extLen := make([]byte, 2)
		binary.BigEndian.PutUint16(extLen, f.PayloadLen)

		err := fw.write(extLen)
		if err != nil {
			utils.LogErr("could not write EXT PAYLOAD Len to conn", err)

			return err
		}
	}

	// Payload
	err = fw.write(f.Payload)
	if err != nil {
		utils.LogErr("could not write PAYLOAD to conn", err)

		return err
	}

	// Flush the data
	fw.flush()

	return nil
}

func (fw *FrameWriter) write(b []byte) error {
	if _, err := fw.w.Write(b); err != nil {
		return err
	}

	return nil
}

func (fw *FrameWriter) Send(data []byte, dt DataType) {
	// Write data Frame
	var op Opcode

	if dt == DataTypeText {
		op = OpText
	} else {
		op = OpBinary
	}

	fr := &Frame{
		Fin:        true,
		Opcode:     op,
		Masked:     false,
		PayloadLen: uint16(len(data)),
		Payload:    data,
	}

	err := fw.WriteFrame(fr)
	if err != nil {
		// Inform about an error to initiate close
		fw.closeCh <- struct{}{}
	}
}

func (fw *FrameWriter) flush() {
	if err := fw.w.Flush(); err != nil {
		utils.LogErr("could not flush data", err)
	}
}
