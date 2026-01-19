package main

import (
	"bufio"
	"net"
)

// Response represesnts the server side of an HTTP response
type response struct {
	conn net.Conn

	req *Request

	header Header

	wroteHeader bool // Tells that a non 1xx header has been written

	status int // status code passed to WriteHeader

	w *bufio.Writer

	contentLength int64
}

func (r *response) Header() Header {
	return r.header
}

func (r *response) WriteHeader(code int) {
	r.wroteHeader = true
	r.status = code
}

func (r *response) Write(data string) (int, error) {
	// Write header if not written
	if !r.wroteHeader {
		r.WriteHeader(StatusOK)
	}

	if len(data) == 0 {
		return 0, nil
	}

	return r.w.WriteString(data)
}

func (r *response) Flush() error {
	return r.w.Flush()
}
