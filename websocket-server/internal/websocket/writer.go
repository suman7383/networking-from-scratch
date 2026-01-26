package websocket

import "bufio"

type FrameWriter struct {
	w *bufio.Writer
}

// TODO
func (fw *FrameWriter) WriteFrame(f *Frame) error {
	return nil
}
