package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
)

// TODO
type Server struct {
	// Addr Specifies the TCP address for the server to listen on,
	// in form "host:port".
	Addr string

	router *Router
}

var CRLF = []byte("\r\n")

func (s *Server) ListenAndServe() error {
	cwd, _ := os.Getwd()
	log.Println("Working directory:", cwd)

	// load the cert
	cert, err := tls.LoadX509KeyPair("certs/server.crt", "certs/server.key")
	if err != nil {
		log.Fatal(err)
	}

	// Create the Tls config
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	ln, err := tls.Listen("tcp", s.Addr, tlsConfig)
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
	res, err := readRequest(reader, conn)

	if err != nil {

		if res == nil {
			// Drop the connection
			return
		}

		switch {
		case errors.Is(err, ErrMalformedRequestLine) || errors.Is(err, ErrInvalidRequestMethod):
			// TODO
			// Send 400 Error response
			slog.Error(err.Error())
			res.SetBadRequestHeader()
			res.finalizeResponse()

		default:
			// TODO
			// Send 500 Error response
			slog.Error(err.Error())
			res.SetInternalServerErrHeader()
			res.finalizeResponse()
		}
		return
	}

	// Testing (Remove this later)
	// fmt.Printf("%v", req)
	// _, err = io.WriteString(conn, "HTTP/1.1 200 OK\r\nContent-Length: 2\r\nContent-Type: text/plain\r\nConnection: close\r\n\r\nOK")

	// if err != nil {
	// 	slog.Error(err.Error())
	// 	return
	// }

	// Route the request according to target-path
	//
	// Change this to res.req.Path later
	h, err := s.route(res.req.RequestURI)

	if err != nil {
		// No route registered for this path
		// Send 404 error
		res.SetNotfoundHeader()
	} else {
		h.ServerHTTP(res, res.req)
	}

	res.finalizeResponse()
}

type HandlerFunc func(ResponseWriter, *Request)

// ServerHTTP calls f(w, r)
func (f HandlerFunc) ServerHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}

var ErrRouteNotFound = errors.New("route not found")

func (s *Server) route(pattern string) (HandlerFunc, error) {
	if h, ok := s.router.get(pattern); !ok {
		return nil, ErrRouteNotFound
	} else {
		return h, nil
	}

}

func ListenAndServe(addr string) (*Server, error) {
	server := &Server{Addr: addr, router: NewRouter()}
	return server, server.ListenAndServe()
}
