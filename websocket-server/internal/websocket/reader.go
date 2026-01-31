package websocket

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"log/slog"
	"net"

	"github.com/suman7383/networking-from-scratch/websocket-server/utils"
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

func (fr *FrameReader) ReadFrame() (*Frame, error) {
	// A Websocket frame looks like this on the wire:
	//
	// | FIN | OPCODE | MASK | PAYLOAD LEN | (EXT LEN) | MASK KEY | PAYLOAD |
	var frame Frame

	err := fr.parseFrameInfo(&frame)
	if err != nil {
		return nil, err
	}

	// MASK KEY
	mKey, err := fr.readMaskKey()
	if err != nil {
		utils.LogErr("could not read mask key", err)

		return nil, err
	}

	frame.MaskKey = [4]byte(mKey)

	// PAYLOAD
	payloadD, err := fr.readPayload(frame.PayloadLen)
	if err != nil {
		utils.LogErr("could not read payload", err)

		return nil, err
	}

	// UNMASK Payload
	fr.unmaskPayload(payloadD, frame.MaskKey[:])
	frame.Payload = payloadD

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

// parseFrameInfo reads from conn and forms these following data:
//
// FIN- whether it is final frame(we only expect FIN 1, fragmentation is unsupported)
//
// OPCODE- Type of frame(continuation, text, binary, close, ping, pong)
//
// MASK- Whether the payload is masked(We always expect this to be 1)
//
// BASE_PAYLOAD_LENGTH- Length of the payload data
//
// EXT_PAYLOAD_LENGTH- If BASE PAYLOAD LEN was > 125, we read EXT PAYLOAD LEN to get actual
// payload length
func (fr *FrameReader) parseFrameInfo(f *Frame) error {
	// Read 2 bytes
	//
	// Parse: FIN(1 bit), OPCODE(4 bits), MASK(1 bit), BASE payload length(7 bits)
	//
	// payload length:
	// 0-125: payload length
	// 126: next 2 bytes = actual length (return error if control frame or actual length < 126)
	// 127: next 8 bytes = actual length (not-supported here so we return error if 127)
	info := make([]byte, 2)

	if _, err := io.ReadFull(fr.r, info); err != nil {
		slog.Error(err.Error())

		return err
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

// Reads the mask key used for unmasking payload data
func (fr *FrameReader) readMaskKey() (key []byte, err error) {
	// Read 4 bytes
	key = make([]byte, 4)

	_, err = io.ReadFull(fr.r, key)
	if err != nil {
		return nil, ErrReadingInfo
	}

	return key, nil
}

// Reads the masked payload data
func (fr *FrameReader) readPayload(payloadLen uint16) (data []byte, err error) {
	data = make([]byte, payloadLen)

	_, err = io.ReadFull(fr.r, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Unmasks the payload data using the mask key
func (fr *FrameReader) unmaskPayload(data, mask []byte) {
	for i := range data {
		data[i] ^= mask[i%4]
	}
}
