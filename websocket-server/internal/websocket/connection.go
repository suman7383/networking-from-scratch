package websocket

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/suman7383/networking-from-scratch/websocket-server/utils"
)

// TODO: Write the necessary methods
type WebSocketConn struct {
	conn   net.Conn
	r      *FrameReader
	w      *FrameWriter
	closed bool // Indicates whether closed Frame was recieved.
}

// TODO
func (w *WebSocketConn) Handle() {
	slog.Info("TODO: Handling websocket conn")
	defer w.conn.Close()

	// Read for incoming frames
	for {
		fr, err := w.r.ReadFrame()
		if err != nil {
			// TODO
			//
			// Send Close control frame with error status
			fmt.Fprintf(w.w.w, "Error, err: %s\r\n", err.Error())

			return
		}

		fmt.Println("Frame parsed successfully")
		printFrame(fr)

		// TODO
		//
		// Decide what to do with the frame
		// REMOVE THE BELOW CODE LATER
		n, err := fmt.Fprintf(w.w.w, "OK\r\n")
		fmt.Printf("Sent bytes %d", n)
		if err != nil {
			utils.LogErr("error sending OK", err)

			return
		}
	}
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
