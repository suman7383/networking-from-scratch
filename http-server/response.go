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

	w *bufio.Writer

	contentLength int64
}
