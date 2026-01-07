package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

type ChatServer struct {
	clients []net.Conn
	ln      net.Listener
}

// Creates a new Chat server
func NewChatServer(port string) (*ChatServer, error) {
	ln, err := net.Listen("tcp", port)

	if err != nil {
		fmt.Println("Error starting the server")
		return nil, err
	}

	return &ChatServer{
		clients: make([]net.Conn, 0),
		ln:      ln,
	}, nil
}

func (c *ChatServer) Start() error {
	fmt.Printf("[SERVER] started accepting connections on: %s\n", c.ln.Addr())

	for {
		conn, err := c.ln.Accept()

		if err != nil {
			fmt.Printf("[ERROR] could not accept client connection from %s\n", conn.RemoteAddr())
			continue
		}

		fmt.Printf("[SERVER] new client connected: %s\n", conn.RemoteAddr())

		// Push this client to clients slice
		c.clients = append(c.clients, conn)

		// Handle the connection
		go c.handleConnection(conn)
	}
}

// Handle connection life cycle
func (c *ChatServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	// Read from connection
	for {
		bytes, err := reader.ReadBytes(byte('\n'))

		if err != nil {
			if err != io.EOF {
				fmt.Printf("[ERROR] reading from connection %s, err: %s\n", conn.RemoteAddr(), err)
			} else {
				fmt.Printf("[ERROR] client connection closed %s\n", conn.RemoteAddr())
			}

			return
		}

		// Broadcast the message to all
		msg := fmt.Sprintf("[CLIENT] %s", bytes)
		c.broadcastExceptSelf(conn, []byte(msg))
	}
}

func (c *ChatServer) broadcastExceptSelf(conn net.Conn, msg []byte) {
	// iterate over clients slice
	for _, client := range c.clients {
		// Skip the sender
		if client == conn {
			continue
		}

		_, err := client.Write(msg)

		if err != nil {
			fmt.Printf("[ERROR] writing to client, %s\n", err)
			continue
		}
	}
}

func run(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("Expected 2 arguments got %d\n", len(args))
	}
	// Format the port properly(:8080)
	port := fmt.Sprintf(":%s", args[1])

	server, err := NewChatServer(port)

	if err != nil {
		os.Exit(1)
	}

	// Start accepting connections
	return server.Start()
}

func main() {
	if err := run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	run(os.Args)
}
