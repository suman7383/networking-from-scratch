package websocket

import (
	"bufio"
	"net"
)

type FrameWriter struct {
	w *bufio.Writer
}

func NewFrameWriter(conn net.Conn) *FrameWriter {
	return &FrameWriter{
		w: bufio.NewWriter(conn),
	}
}

// TODO
func (fw *FrameWriter) WriteFrame(f *Frame) error {
	return nil
}
