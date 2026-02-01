package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// =====================
// Entry point
// =====================

func main() {
	testTextThenClose()
}

// =====================
// Handshake
// =====================

func dialAndHandshake() (*bufio.Reader, *bufio.Writer, net.Conn) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	req := "" +
		"GET / HTTP/1.1\r\n" +
		"Host: localhost:8080\r\n" +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==\r\n" +
		"Sec-WebSocket-Version: 13\r\n" +
		"\r\n"

	w.WriteString(req)
	w.Flush()

	for {
		line, _ := r.ReadString('\n')
		if line == "\r\n" {
			break
		}
	}

	return r, w, conn
}

// =====================
// Frame reader (server → client)
// =====================

func readFrame(r *bufio.Reader) (opcode byte, payload []byte, err error) {
	b1, err := r.ReadByte()
	if err != nil {
		return
	}

	b2, err := r.ReadByte()
	if err != nil {
		return
	}

	opcode = b1 & 0x0F
	masked := (b2 & 0x80) != 0
	payloadLen := int(b2 & 0x7F)

	if masked {
		return 0, nil, fmt.Errorf("protocol error: server frame masked")
	}

	if payloadLen == 126 {
		ext := make([]byte, 2)
		io.ReadFull(r, ext)
		payloadLen = int(binary.BigEndian.Uint16(ext))
	}

	payload = make([]byte, payloadLen)
	if payloadLen > 0 {
		io.ReadFull(r, payload)
	}

	return
}

// =====================
// Frame writer (client → server)
// =====================

func writeCloseFrame(w *bufio.Writer) {
	// FIN=1, OPCODE=8
	w.WriteByte(0x88)

	// MASK=1, LEN=0
	w.WriteByte(0x80)

	// Mask key (required even for len=0)
	w.Write([]byte{0x01, 0x02, 0x03, 0x04})
	w.Flush()
}

func writeTextFrame(w *bufio.Writer, text string) {
	payload := []byte(text)
	mask := []byte{0x01, 0x02, 0x03, 0x04}

	w.WriteByte(0x81)                      // FIN=1, TEXT
	w.WriteByte(0x80 | byte(len(payload))) // MASK=1, LEN

	w.Write(mask)

	for i := range payload {
		payload[i] ^= mask[i%4]
	}
	w.Write(payload)
	w.Flush()
}

// =====================
// Test: TEXT → CLOSE handshake
// =====================

func testTextThenClose() {
	fmt.Println("TEST: TEXT frame then graceful CLOSE")

	r, w, conn := dialAndHandshake()
	defer conn.Close()

	// 1️⃣ Send TEXT
	writeTextFrame(w, "hi")

	// 2️⃣ Read server response
	opcode, payload, err := readFrame(r)
	if err != nil {
		fmt.Println("read error:", err)
		return
	}

	fmt.Printf("SERVER OPCODE: 0x%x PAYLOAD: %q\n", opcode, payload)

	// 3️⃣ Expect CLOSE or normal response
	if opcode == 0x8 {
		fmt.Println("SERVER sent CLOSE → replying with CLOSE")
		writeCloseFrame(w)
		return
	}

	// 4️⃣ Client initiates CLOSE
	fmt.Println("CLIENT initiating CLOSE")
	writeCloseFrame(w)

	// 5️⃣ Expect server CLOSE
	opcode, _, err = readFrame(r)
	if err != nil {
		fmt.Println("read error:", err)
		return
	}

	if opcode == 0x8 {
		fmt.Println("SERVER replied with CLOSE → closing TCP")
	}
}
