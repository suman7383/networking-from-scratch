package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func main() {
	// Choose ONE test at a time
	// Comment / uncomment what you want to test

	testTextPayloadSmall() // payload <= 125 (success)
	// testControlFramePing() // control frame (PING)
	// testPayloadLen126() // payload len == 126 (success)
	// testFragmentedFrame() // FIN = 0 (should fail)
	// testMaskBitZero() // MASK = 0 (should fail)
	// testRSVBitSet() // RSV bit set (should fail)
}

func dialAndHandshake() (*bufio.Reader, *bufio.Writer, func()) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}

	cleanup := func() {
		conn.Close()
	}

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	handshake := "" +
		"GET / HTTP/1.1\r\n" +
		"Host: localhost:8080\r\n" +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==\r\n" +
		"Sec-WebSocket-Version: 13\r\n" +
		"\r\n"

	if _, err := w.WriteString(handshake); err != nil {
		panic(err)
	}
	w.Flush()

	// Read handshake response headers
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			panic(err)
		}
		if line == "\r\n" {
			break
		}
	}

	return r, w, cleanup
}

func readServerResponse(r *bufio.Reader) {
	for {

		buf := make([]byte, 4)
		n, err := io.ReadFull(r, buf)
		if err == nil {
			fmt.Println("Server response:")
			fmt.Println(string(buf[:n]))
		}
	}
}

// =====================
// Success: payload len <= 125
// =====================

func testTextPayloadSmall() {
	fmt.Println("TEST: Text frame, payload <= 125")

	r, w, cleanup := dialAndHandshake()
	defer cleanup()

	// "hi"
	frame := []byte{
		0x81,                   // FIN=1, TEXT
		0x82,                   // MASK=1, len=2
		0x01, 0x02, 0x03, 0x04, // mask
		0x69, 0x6b, // masked "hi"
	}

	w.Write(frame)
	w.Flush()

	readServerResponse(r)
}

// =====================
// Control frame (PING)
// =====================

func testControlFramePing() {
	fmt.Println("TEST: Control frame (PING)")

	r, w, cleanup := dialAndHandshake()
	defer cleanup()

	frame := []byte{
		0x89,                   // FIN=1, PING
		0x80,                   // MASK=1, len=0
		0x01, 0x02, 0x03, 0x04, // mask (required even if len=0)
	}

	w.Write(frame)
	w.Flush()

	readServerResponse(r)
}

// =====================
// Success: payload len == 126
// =====================

func testPayloadLen126() {
	fmt.Println("TEST: Text frame, payload len == 126")

	r, w, cleanup := dialAndHandshake()
	defer cleanup()

	payloadLen := 126
	payload := make([]byte, payloadLen)
	mask := []byte{0x01, 0x02, 0x03, 0x04}

	for i := 0; i < payloadLen; i++ {
		payload[i] = byte('a') ^ mask[i%4]
	}

	frame := []byte{
		0x81,       // FIN=1, TEXT
		0xFE,       // MASK=1, len=126
		0x00, 0x7E, // extended payload length = 126
	}
	frame = append(frame, mask...)
	frame = append(frame, payload...)

	w.Write(frame)
	w.Flush()

	readServerResponse(r)
}

// =====================
// FAIL: FIN = 0 (fragmented frame)
// =====================

func testFragmentedFrame() {
	fmt.Println("TEST: FIN = 0 (fragmentation)")

	r, w, cleanup := dialAndHandshake()
	defer cleanup()

	frame := []byte{
		0x01,                   // FIN=0, TEXT
		0x82,                   // MASK=1, len=2
		0x01, 0x02, 0x03, 0x04, // mask
		0x69, 0x6b,
	}

	w.Write(frame)
	w.Flush()

	readServerResponse(r)
}

// =====================
// FAIL: MASK = 0
// =====================

func testMaskBitZero() {
	fmt.Println("TEST: MASK = 0 (client violation)")

	r, w, cleanup := dialAndHandshake()
	defer cleanup()

	frame := []byte{
		0x81, // FIN=1, TEXT
		0x02, // MASK=0, len=2 ❌
		'h', 'i',
	}

	w.Write(frame)
	w.Flush()

	readServerResponse(r)
}

// =====================
// FAIL: RSV bit set
// =====================

func testRSVBitSet() {
	fmt.Println("TEST: RSV bit set")

	r, w, cleanup := dialAndHandshake()
	defer cleanup()

	frame := []byte{
		0xC1,                   // FIN=1, RSV1=1, TEXT ❌
		0x82,                   // MASK=1, len=2
		0x01, 0x02, 0x03, 0x04, // mask
		0x69, 0x6b,
	}

	w.Write(frame)
	w.Flush()

	readServerResponse(r)
}
