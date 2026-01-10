package main

import (
	"fmt"
	"net"
)

// TODO
type HttpServer struct {
	ln net.Listener
}

func NewHttpServer(port string) (*HttpServer, error) {
	ln, err := net.Listen("tcp", port)

	if err != nil {
		fmt.Println("Error creating http server, err:", err)
		return nil, err
	}

	return &HttpServer{
		ln: ln,
	}, nil
}

func (h *HttpServer) Start() error {
	// Accept active connections
	for {
		conn, err := h.ln.Accept()

		if err != nil {
			fmt.Println("error accepting client connection, err", err)
			continue
		}

		go h.handleConnection(conn)
	}
}

func (h *HttpServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	// TODO Parse the incoming request

	// request-line   = method SP request-target SP HTTP-version CRLF

	// header-field   = field-name ":" OWS field-value OWS  (Where OWS = Optional White Space)
}

// TODO
func main() {

}
