package main

import (
	"fmt"
	"net"
)

// TODO
type Server struct {
	// Addr Specifies the TCP address for the server to listen on,
	// in form "host:port".
	Addr string
}

func (s *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", s.Addr)

	if err != nil {
		return err
	}
	return s.Serve(ln)
}

func (s *Server) Serve(ln net.Listener) error {

	// Accept active connections
	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println("error accepting client connection, err", err)
			continue
		}

		go s.serve(conn)
	}
}

func (s *Server) serve(conn net.Conn) {
	defer conn.Close()

	// TODO Parse the incoming request

	// request-line   = method SP request-target SP HTTP-version CRLF

	// header-field   = field-name ":" OWS field-value OWS  (Where OWS = Optional White Space)
}

func ListenAndServe(addr string) error {
	server := &Server{addr}
	return server.ListenAndServe()
}
