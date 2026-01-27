package websocket

import (
	"bufio"
	"errors"
	"io"
	"log/slog"
	"net"
)

type FrameReader struct {
	r *bufio.Reader
}

func NewFrameReader(conn net.Conn) *FrameReader {
	return &FrameReader{
		r: bufio.NewReader(conn),
	}
}

var ErrReadingInfo = errors.New("could not read frame")

// TODO
func (fr *FrameReader) ReadFrame() (*Frame, error) {
	// A Websocket frame looks like this on the wire:
	//
	// | FIN | OPCODE | MASK | PAYLOAD LEN | (EXT LEN) | MASK KEY | PAYLOAD |
	var frame Frame

	err := fr.parseFrameInfo(&frame)
	if err != nil {
		return nil, err
	}

	// TODO
	// MASK | PAYLOAD LEN

	// TODO
	// EXT LEN

	// TODO
	// MASK KEY

	// TODO
	// PAYLOAD

	return &frame, nil
}

const fin_mask = (1 << 7)            // 7th bit
const rsv_mask = ((1 << 3) - 1) << 4 // 4th, 5th, 6th bits (We do not care about RSV bits now)
const opcode_mask = (1 << 4) - 1     // 0 to 3rd bits set

var ErrUnsupportedFragmentation = errors.New("fragmentation not supported")
var ErrExtensionNotSupported = errors.New("extension not supported")

func (fr *FrameReader) parseFrameInfo(f *Frame) error {
	// Read 2 bytes
	//
	// Parse: FIN(1 bit), OPCODE(4 bits), MASK(1 bit), BASE payload length(7 bits)
	//
	// payload length:
	// 0-125: payload length
	// 126: next 2 bytes = actual length
	// 127: next 8 bytes = actual length
	info := make([]byte, 2)

	if _, err := io.ReadFull(fr.r, info); err != nil {
		slog.Error(err.Error())

		return ErrReadingInfo
	}

	// FIN
	fin := info[0] & fin_mask

	if fin == 0 {
		return ErrUnsupportedFragmentation
	}

	f.Fin = true

	// RSV
	//
	// We ignore if rsv is 0 but throw error if > 0
	rsv := info[0] & rsv_mask

	if rsv > 0 {
		return ErrExtensionNotSupported
	}

	opcode := info[0] & opcode_mask
	f.Opcode = Opcode(opcode)

	return nil
}
