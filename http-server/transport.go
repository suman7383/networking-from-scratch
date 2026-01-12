package main

import (
	"fmt"
	"net"
)

// TODO
type Tcp struct {
	ln net.Listener
}

func NewTcpServer(port string) (*Tcp, error) {
	ln, err := net.Listen("tcp", port)

	if err != nil {
		fmt.Println("Error creating http server, err:", err)
		return nil, err
	}

	return &Tcp{
		ln: ln,
	}, nil
}

func (h *Tcp) ListenAndServe() error {
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

func (h *Tcp) handleConnection(conn net.Conn) {
	defer conn.Close()

	// TODO Parse the incoming request

	// request-line   = method SP request-target SP HTTP-version CRLF

	// header-field   = field-name ":" OWS field-value OWS  (Where OWS = Optional White Space)
}

// TODO
func main() {

}
