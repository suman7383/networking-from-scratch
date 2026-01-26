package websocket

import (
	"log/slog"
	"net"
)

// TODO: Write the necessary methods
type WebSocketConn struct {
	conn net.Conn
}

// TODO
func (w *WebSocketConn) Handle() {
	slog.Info("TODO: Handling websocket conn")
}
