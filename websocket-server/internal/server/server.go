package server

import (
	"crypto/tls"
	"errors"
	"log"
	"log/slog"
	"net"
	"os"

	"github.com/suman7383/networking-from-scratch/websocket-server/internal/httpcore"
	"github.com/suman7383/networking-from-scratch/websocket-server/internal/websocket"
	"github.com/suman7383/networking-from-scratch/websocket-server/utils"
)

type Server struct {
	// Addr Specifies the TCP address for the server to listen on,
	// in form "host:port".
	Addr    string
	Handler websocket.HandlerFunc
}

func NewServer(addr string, handler websocket.HandlerFunc) *Server {
	return &Server{
		Addr:    addr,
		Handler: handler,
	}
}

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

	slog.Info("[SERVER] started listening\n")

	// Listen for new connection
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		slog.Info("[SERVER] client connected: ", slog.String("Addr", conn.RemoteAddr().String()))

		go s.serve(conn)
	}
}

func (s *Server) serve(conn net.Conn) {
	defer conn.Close()

	// TODO: Read and parse the incomming request
	reader := httpcore.NewReader(conn)

	req, err := httpcore.ReadRequest(reader)

	if err != nil {
		// Send error response
		utils.WriteErrResponse(conn, httpcore.StatusBadRequest, err.Error())
		return
	}

	// Verify if it is an upgrade(websocket) request
	// If not, respond with a error
	err = s.handleUpgradeControl(req)

	if err != nil {
		// Send error (no upgrade header found)
		utils.WriteErrResponse(conn, httpcore.StatusBadRequest, err.Error())
		return
	}

	// Handshake
	wsc, err := websocket.HandleHandshake(req, conn, s.Handler)

	if err != nil {
		// TODO: Handshake failure. Failure response is already sent
		//
		// maybe log the error here and return ?
		slog.Error(err.Error())

		return
	}

	// Pass the control to WebSocket handler
	wsc.Handle()
}

var ErrMissingConnectionUpgrade = errors.New("Missing Connection upgrade header")
var ErrUnsupportedUpgrade = errors.New("Provided upgrade not support")

// TODO
func (s *Server) handleUpgradeControl(req *httpcore.Request) error {
	if !req.ConnectionUpgrade {
		return ErrMissingConnectionUpgrade
	}

	if req.Upgrade != "websocket" {
		return ErrUnsupportedUpgrade
	}

	slog.Info("Correct upgrade headers found")

	return nil
}
