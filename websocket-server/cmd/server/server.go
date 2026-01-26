package server

import (
	"errors"
	"log/slog"
	"net"

	"github.com/suman7383/networking-from-scratch/websocket-server/internal/http"
	"github.com/suman7383/networking-from-scratch/websocket-server/internal/websocket"
	"github.com/suman7383/networking-from-scratch/websocket-server/utils"
)

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
	reader := http.NewReader(conn)

	req, err := http.ReadRequest(reader)

	if err != nil {
		// Send error response
		utils.WriteErrResponse(conn, http.StatusBadRequest, err.Error())
		return
	}

	// Verify if it is an upgrade(websocket) request
	// If not, respond with a error
	err = s.handleUpgradeControl(req)

	if err != nil {
		// Send error (no upgrade header found)
		utils.WriteErrResponse(conn, http.StatusBadRequest, err.Error())
		return
	}

	// Handshake
	wsc, err := websocket.HandleHandshake(req, conn)

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
func (s *Server) handleUpgradeControl(req *http.Request) error {
	if !req.ConnectionUpgrade {
		return ErrMissingConnectionUpgrade
	}

	if req.Upgrade != "websocket" {
		return ErrUnsupportedUpgrade
	}

	slog.Info("Correct upgrade headers found")

	return nil
}
