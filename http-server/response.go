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

	body []byte // buffer for response body

	contentLength int64
}

func (r *response) Header() Header {
	// returns the headers
	return r.header
}

func (r *response) WriteHeader(code int) {
	// Write headers if not already written
	if !r.wroteHeader {
		r.wroteHeader = true
		r.status = code
	}
}

func (r *response) Write(data []byte) (n int, err error) {
	return r.write(len(data), data)
}

func (r *response) write(len int, dataB []byte) (n int, err error) {
	// Write header if not written
	if !r.wroteHeader {
		r.WriteHeader(StatusOK)
	}

	if len == 0 {
		return 0, nil
	}

	if dataB != nil {
		r.body = append(r.body, dataB...)

		return n, nil
	}

	return 0, nil
}

// TODO
func (r *response) finalizeResponse() error {
	// write status-line

	// set contentLength

	// write headers

	// write body
}

func (r *response) flush() error {
	return r.w.Flush()
}
