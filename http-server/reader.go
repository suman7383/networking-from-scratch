package main

import (
	"bufio"
	"net"
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

func (r *Reader) ReadLine() []byte {
	// TODO

	return nil
}

func (r *Reader) ReadN(n int) []byte {
	// TODO

	return nil
}
