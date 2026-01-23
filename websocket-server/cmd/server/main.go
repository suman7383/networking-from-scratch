package main

import "net"

type Server struct {
	// Addr Specifies the TCP address for the server to listen on,
	// in form "host:port".
	Addr string
}

func NewServer(addr string) *Server {
	return &Server{
		Addr: addr,
	}
}

func (s *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	// Listen for new connection
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		go s.serve(conn)
	}
}

func (s *Server) serve(conn net.Conn) {
	defer conn.Close()

	// TODO: Read and parse the incomming request
	//
	// Verify if it is an upgrade(websocket) request
	//
	// If not, respond with a error

	// TODO: If handshake success, handle the conn to websocket handler
}
