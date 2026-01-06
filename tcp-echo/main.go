package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("[SERVER] error starting TCP server, err:", err)
		os.Exit(1)
	}

	fmt.Printf("[SERVER] listening on %s\n", ln.Addr())

	// Continuously listen for new connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("error connecting to client, err:", err)
			continue
		}

		fmt.Printf("[SERVER] client connected %s\n", conn.LocalAddr())

		// handle the connection
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		bytes, err := reader.ReadBytes(byte('\n'))
		if err != nil {
			if err != io.EOF {
				fmt.Println("failed to read data, err:", err)
			}
			fmt.Printf("[SERVER] client closed connection: %s\n", conn.LocalAddr())
			return
		}
		fmt.Printf("request: %s", bytes)
		bytes = append(bytes, "sent by server\n"...)

		// send the message back to the client
		conn.Write(bytes)
	}
}
