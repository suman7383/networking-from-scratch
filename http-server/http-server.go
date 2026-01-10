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

// TODO
func main() {

}
