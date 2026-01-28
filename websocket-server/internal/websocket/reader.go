package websocket

import (
	"bufio"
	"encoding/binary"
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
const maskP_mask = (1 << 7)          // 7th bit
const payloadLen_mask = (1 << 7) - 1 // 0 to 6th bits set

var ErrUnsupportedFragmentation = errors.New("fragmentation not supported")
var ErrExtensionNotSupported = errors.New("extension not supported")
var ErrProtocol = errors.New("protocol error")
var ErrPayloadTooLarge = errors.New("payload is too large")

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

	// OPCODE
	opcode := info[0] & opcode_mask
	f.Opcode = Opcode(opcode)

	// MASK
	maskP := info[1] & maskP_mask
	if maskP == 0 {
		return ErrProtocol
	}

	f.Masked = true

	// PAYLOAD Len
	// Check for len 127, 126 and <125
	//
	// 127 not supported for now
	switch payloadLen := info[1] & payloadLen_mask; payloadLen {
	case 127:
		return ErrPayloadTooLarge
	case 126:
		// Throw if Control frame
		if f.Opcode.IsControlFrame() {
			return ErrProtocol
		}

		// Read next 2 bytes to get actual length
		if epl, err := fr.readExtPayloadLen16(); err != nil {
			return err
		} else {
			// Reject if actual length is < 126
			// And
			// Reject if actual length > 126 and control frame
			if epl < 126 {
				return ErrProtocol
			}

			f.PayloadLen = epl
		}
	default:
		// payloadLen is <=125
		f.PayloadLen = uint16(payloadLen)
	}

	return nil
}

// Reads next 2 bytes(16 bit)
func (fr *FrameReader) readExtPayloadLen16() (len uint16, err error) {
	extPayloadLen := make([]byte, 2)

	if _, err := io.ReadFull(fr.r, extPayloadLen); err != nil {
		slog.Error("could not read extended payload length", slog.String("err", err.Error()))

		return 0, ErrReadingInfo
	}

	return binary.BigEndian.Uint16(extPayloadLen), nil
}
