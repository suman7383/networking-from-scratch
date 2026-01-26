package utils

import (
	"net"

	"github.com/suman7383/networking-from-scratch/websocket-server/internal/http"
)

func WriteErrResponse(conn net.Conn, status int, errMsg string) {
	res := http.NewResponse(conn)

	// write status
	res.WriteHeader(status)

	// write response body
	res.Write([]byte(errMsg))

	res.FinalizeResponse(true, true)
}
