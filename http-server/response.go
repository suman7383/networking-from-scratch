package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"time"
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

	wantKeepAlive bool // Used for Connection header
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
	return r.write(len(data), data, "")
}

func (r *response) WriteString(data string) (n int, err error) {
	return r.write(len(data), nil, data)
}

func (r *response) write(len int, dataB []byte, dataS string) (n int, err error) {
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
	} else {
		r.body = append(r.body, []byte(dataS)...)

		return n, nil
	}
}

// Parse the response and send to wire(conn)
func (r *response) finalizeResponse() {
	// write status-line
	if !r.wroteHeader {
		r.WriteHeader(StatusOK)
	}

	// example status line : "HTTP/1.1 <status-code> <reason-phrase>\r\n"
	var sb strings.Builder

	sb.WriteString("HTTP/1.1 ")
	sb.WriteString(strconv.Itoa(r.status) + " ")
	sb.WriteString(StatusText(r.status))
	sb.Write(CRLF)

	sl := sb.String()

	// Write status line to wire
	r.writeToWire(nil, sl)

	// Set Auto headers
	r.setAutoHeaders()

	// set contentLength
	r.Header().Add("Content-Length", strconv.Itoa(len(r.body)))

	// Reset the string builder for headers
	sb.Reset()

	// Parse Headers
	for k, v := range r.Header() {

		if !r.validateHeaderField(k) {
			panic(fmt.Sprintf("Invalid response header field name, %s", k))
		}

		sb.WriteString(k)
		sb.WriteString(": ")

		// comma separeted values
		for i, vv := range v {
			// No comma for last element
			if i == len(v)-1 {
				sb.WriteString(vv)
			} else {
				sb.WriteString(vv + ", ")
			}
		}

		sb.Write(CRLF)
	}

	// write headers
	r.writeToWire(nil, sb.String())

	r.writeToWire(CRLF, "")

	// write body
	r.writeToWire(r.body, "")

	err := r.flush()

	if err != nil {
		slog.Error(err.Error())
	}
}

// Writes bytes to the wire(conn)
func (r *response) writeToWire(dataB []byte, dataS string) (n int, err error) {
	if dataB != nil {
		return r.w.Write(dataB)
	} else {
		return r.w.Write([]byte(dataS))
	}
}

// TODO: Implement header validation (per RFC 7230)
func (r *response) validateHeaderField(k string) bool {
	return true
}

// Sets Date, Content-Type, Connection
func (r *response) setAutoHeaders() {
	// Date
	r.Header().Add("Date", time.Now().UTC().Format(time.RFC1123))

	// Connection
	// For now we close on each req
	r.Header().Set("Connection", "close")

	// Content-Type(defaults to text/plain)
	if v := r.Header().Get("Content-Type"); len(v) == 0 {
		r.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}

}

func (r *response) flush() error {
	return r.w.Flush()
}

func (r *response) SetNotfoundHeader() {
	r.WriteHeader(StatusNotFound)
}

func (r *response) SetInternalServerErrHeader() {
	r.WriteHeader(StatusInternalServerError)
}

func (r *response) SetBadRequestHeader() {
	r.WriteHeader(StatusBadRequest)
}
