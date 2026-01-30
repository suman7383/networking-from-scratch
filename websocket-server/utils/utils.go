package utils

import (
	"log/slog"
	"net"

	"github.com/suman7383/networking-from-scratch/websocket-server/internal/httpcore"
)

func WriteErrResponse(conn net.Conn, status int, errMsg string) {
	res := httpcore.NewResponse(conn)

	// write status
	res.WriteHeader(status)

	// write response body
	res.Write([]byte(errMsg))

	res.FinalizeResponse(true, true)
}

func LogErr(msg string, err error) {
	slog.Error(msg, slog.String("err", err.Error()))
}
