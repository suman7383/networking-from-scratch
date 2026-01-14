package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
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

	fmt.Println("[SERVER] started listening")
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

		slog.Info(fmt.Sprintf("client connected: %s\n", conn.RemoteAddr()))

		go s.serve(conn)
	}
}

func (s *Server) serve(conn net.Conn) {
	defer conn.Close()

	// Wrap the conn with Reader
	reader := NewReader(conn)

	// Parse the HTTP request line
	req, err := readRequest(reader)

	if err != nil {
		switch {
		case errors.Is(err, ErrMalformedRequestLine) || errors.Is(err, ErrInvalidRequestMethod):
			// TODO
			// Send 400 Error response
			slog.Error(err.Error())
			_, errW := io.WriteString(conn, fmt.Sprintf("HTTP/1.1 400 Bad Request\r\n\r\n%s\n", err.Error()))
			if errW != nil {
				slog.Error(err.Error())
				return
			}

		default:
			// TODO
			// Send 500 Error response
			slog.Error(err.Error())
			_, errW := io.WriteString(conn, fmt.Sprintf("HTTP/1.1 500 Internal Server Error\r\n\r\n%s\n", err.Error()))
			if errW != nil {
				slog.Error(err.Error())
				return
			}
		}
		return
	}

	// Testing (Remove this later)
	fmt.Printf("%v", req)
	_, err = io.WriteString(conn, "HTTP/1.1 200 OK\r\nContent-Length: 2\r\nContent-Type: text/plain\r\nConnection: close\r\n\r\nOK")

	if err != nil {
		slog.Error(err.Error())
		return
	}
	// header-field   = field-name ":" OWS field-value OWS  (Where OWS = Optional White Space)
}

func ListenAndServe(addr string) error {
	server := &Server{addr}
	return server.ListenAndServe()
}
