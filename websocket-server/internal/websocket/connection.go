package websocket

import (
	"log/slog"
	"net"
)

// TODO: Write the necessary methods
type WebSocketConn struct {
	conn net.Conn
	r    *FrameReader
	w    *FrameWriter
}

// TODO
func (w *WebSocketConn) Handle() {
	slog.Info("TODO: Handling websocket conn")
}
