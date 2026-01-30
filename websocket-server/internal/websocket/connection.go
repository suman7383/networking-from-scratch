package websocket

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/suman7383/networking-from-scratch/websocket-server/utils"
)

// TODO: Write the necessary methods
type WebSocketConn struct {
	conn         net.Conn
	r            *FrameReader
	w            *FrameWriter
	closeSent    bool // Whether close Frame was sent.
	closeReceive bool // Whether close frame received
	closeCh      chan struct{}
	closed       bool // Whether the TCP conn is closed
}

// TODO
func (w *WebSocketConn) Handle() {
	defer func() {
		if !w.closed {
			w.closeTCPConn()
		}
	}()

	// Read for incoming frames
	for {
		fr, err := w.r.ReadFrame()
		if err != nil {
			// Send Close control frame with error status
			utils.LogErr("reading frame error", err)
			w.initiateClose(CloseProtocolErr, CloseProtocolErr.String())
			continue
		}

		// FOR DEBUGGING
		printFrame(fr)

		switch fr.Opcode {
		case OpPing:
			// Send PONG FRAME

			frW := fr.Clone()

			// Set Opcode to PONG
			frW.Opcode = OpPong

			// Set Masked to false(server to client payload is not masked)
			frW.Masked = false

			w.w.WriteFrame(frW)
		case OpClose:
			// TODO
			//
			// Send CLOSE FRAME
			//
			// If CLOSE frame is not already sent by server
			// we send CLOSE FRAME
			if !w.closeSent {
				w.sendCloseFrame()
			} else {
				w.closeReceived()
			}

			return
		case OpContinuation:
			// Send CLOSE FRAME
			w.initiateClose(CloseProtocolErr, "Fragmentation unsupported")
		case OpText, OpBinary:
			// TODO
			//
			// Handle this data to user(application layer) to handle
		default:
			// CONTROL SHOULD NEVER REACH HERE
			// Send close frame
			w.initiateClose(ClosePolicyViolation, ClosePolicyViolation.String())
		}

	}
}

// Marks closeReceived to true and sends signal to "close" channel
func (w *WebSocketConn) closeReceived() {
	w.closeReceive = true
	close(w.closeCh)
}

const DEFAULT_CLOSE_TIMEOUT = 5 * time.Second

func (w *WebSocketConn) initiateClose(code CloseStatus, reason string) {
	if w.closeSent {
		return
	}

	w.sendErrCloseFrame(code, reason)
	w.closeSent = true

	// Wait for client close frame or timeout
	go func() {
		select {
		case <-w.closeCh:
			// Safe to close the conn
			w.closeTCPConn()
		case <-time.After(DEFAULT_CLOSE_TIMEOUT):
			w.closeTCPConn()
		}
	}()
}

func (w *WebSocketConn) closeTCPConn() {
	w.conn.Close()
	w.closed = true
}

var ErrConnectionClosing = errors.New("Writes closed, connection closing")

func (w *WebSocketConn) writeFrame(f *Frame) error {
	if w.closeSent {
		return ErrConnectionClosing
	}

	return w.w.WriteFrame(f)
}

func (w *WebSocketConn) sendErrCloseFrame(code CloseStatus, reason string) {
	if w.closeSent {
		return
	}

	fr := CloseFrame(code, []byte(reason))

	w.w.WriteFrame(fr)
}

func (w *WebSocketConn) sendCloseFrame() {
	frW := CloseFrame(CloseNormal, []byte(CloseNormal.String()))

	w.writeFrame(frW)
	w.closeSent = true
}

type CloseStatus uint16

const (
	CloseNormal          CloseStatus = 1000
	CloseGoingAway       CloseStatus = 1001
	CloseProtocolErr     CloseStatus = 1002
	CloseUnsupportedData CloseStatus = 1003
	CloseInvalidUTF8     CloseStatus = 1007
	ClosePolicyViolation CloseStatus = 1008 // General error if we don't know a status code
	CloseMessageTooBig   CloseStatus = 1009
	CloseInternalError   CloseStatus = 1011
)

func (cs CloseStatus) String() string {
	switch cs {
	case CloseNormal:
		return "Normal"
	case CloseGoingAway:
		return "Going Away"
	case CloseProtocolErr:
		return "Protocol Error"
	case CloseUnsupportedData:
		return "Unsupported Data"
	case CloseInvalidUTF8:
		return "Invalid UTF-8"
	case ClosePolicyViolation:
		return "Policy Violation"
	case CloseMessageTooBig:
		return "Message Too Big"
	default:
		return "Internal Error"
	}
}

func CloseFrame(statusCode CloseStatus, reason []byte) *Frame {
	fr := &Frame{
		Fin:    true,
		Opcode: OpClose,
		Masked: false,
	}

	if len(reason) > 0 && len(reason) < 125 {
		// 2 bytes statusCode rest reason
		payload := make([]byte, 2+len(reason))
		binary.BigEndian.PutUint16(payload[:2], uint16(statusCode))
		copy(payload[2:], reason)

		fr.PayloadLen = uint16(len(payload))
		fr.Payload = payload
	} else {
		// 2 bytes statusCode only
		payload := make([]byte, 2)
		binary.BigEndian.PutUint16(payload, uint16(statusCode))

		fr.PayloadLen = uint16(len(payload))
		fr.Payload = payload
	}

	return fr
}

// Prints the frame for debugging
func printFrame(f *Frame) {
	fmt.Println("====== WebSocket Frame ======")

	fmt.Printf("FIN        : %v\n", f.Fin)
	fmt.Printf("Opcode     : %v (%s)\n", f.Opcode, f.Opcode.String())
	fmt.Printf("PayloadLen : %d\n", f.PayloadLen)
	fmt.Printf("Masked     : %v\n", f.Masked)

	if f.Masked {
		fmt.Printf("Mask Key   : %02x %02x %02x %02x\n",
			f.MaskKey[0],
			f.MaskKey[1],
			f.MaskKey[2],
			f.MaskKey[3],
		)
	}

	if len(f.Payload) == 0 {
		fmt.Println("Payload    : <empty>")
	} else {
		fmt.Printf("Payload    : %q\n", f.Payload)
		fmt.Printf("PayloadHex :")
		for _, b := range f.Payload {
			fmt.Printf(" %02x", b)
		}
		fmt.Println()
	}

	fmt.Println("=============================")
}
