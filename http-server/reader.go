package main

import (
	"bufio"
	"net"
	"strings"
)

// Reader is responsible for reading from the connection
type Reader struct {
	// reader is a buffered reader for the connection
	reader *bufio.Reader
}

func NewReader(conn net.Conn) *Reader {
	return &Reader{
		reader: bufio.NewReader(conn),
	}
}

// ReadLine reads until delimiter '\n' from the reader and
// checks for valid CRLF('\r\n').
// It returns the line(string) without CRLF('\r\n').
//
// Note: It returns string because HTTP request line & headers
// are ASCII text
func (r *Reader) ReadLine() (string, error) {
	line, err := r.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	// Check for valid CRLF('\r\n')
	if !strings.HasSuffix(line, "\r\n") {
		return "", ErrMalformedRequestLine
	}

	return strings.TrimSuffix(line, "\r\n"), nil
}

// ReadN reads n bytes from reader.
// The bytes are taken from at most one Read on the underlying Reader,
// hence bytes read may be less than n.
func (r *Reader) ReadN(n int) (bytesRead int, buf []byte, err error) {
	buf = make([]byte, n)

	bytesRead, err = r.reader.Read(buf)

	return
}
