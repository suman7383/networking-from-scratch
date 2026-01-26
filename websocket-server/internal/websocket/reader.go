package websocket

import "bufio"

type FrameReader struct {
	r *bufio.Reader
}

// TODO
func (fr *FrameReader) ReadFrame() (*Frame, error) {
	return nil, nil
}
