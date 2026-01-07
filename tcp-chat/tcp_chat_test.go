package main

import (
	"bufio"
	"io"
	"net"
	"testing"
	"time"
)

func TestChatServer(t *testing.T) {
	send := "Test"
	expect := "[CLIENT] Test\n"

	args := []string{"cmd", "8080"}

	// Start the server
	go func() {
		run(args)
	}()

	// Wait for the server to start
	time.Sleep(2 * time.Second)

	c1, err := net.Dial("tcp", ":8080")
	if err != nil {
		t.Fatalf("Error connecting to client: %s", err)
	}

	c2, err := net.Dial("tcp", ":8080")
	if err != nil {
		t.Fatalf("Error connecting to client: %s", err)
	}

	// Send a message
	c1.Write([]byte(send + "\n"))

	// Expect c2 to receive the message
	reader := bufio.NewReader(c2)

	bytes, err := reader.ReadBytes(byte('\n'))

	if err != nil && err != io.EOF {
		t.Fatalf("Error reading from connection: %s", err)
	}

	msgR := string(bytes)

	if msgR != expect {
		t.Fatalf("Expected %s, got %s", expect, msgR)
	}
}
