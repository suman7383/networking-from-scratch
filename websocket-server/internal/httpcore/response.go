package httpcore

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	header Header

	status int

	wroteHeader bool

	Protocol string // e.g: HTTP/1.1

	w *bufio.Writer

	body []byte

	contentLength int64
}

// Creates a new Response with protocol set to "HTTP/1.1"
func NewResponse(conn net.Conn) *Response {
	return &Response{
		header:        make(Header),
		body:          make([]byte, 0),
		Protocol:      "HTTP/1.1",
		w:             bufio.NewWriter(conn),
		contentLength: -1,
	}
}

func (r *Response) Write(b []byte) {
	r.body = append(r.body, b...)
}

func (r *Response) WriteHeader(code int) {
	if !r.wroteHeader {
		r.status = code
		r.wroteHeader = true
	}
}

func (r *Response) Header() Header {
	return r.header
}

func (r *Response) Body() []byte {
	return r.body
}

// Sets Date, Content-Type, Connection
func (r *Response) setAutoHeaders(closeConn bool) {
	// Date
	r.Header().Add("Date", time.Now().UTC().Format(time.RFC1123))

	// Connection
	if closeConn {
		r.Header().Set("Connection", "close")
	}

	// Content-Type(defaults to text/plain)
	if v := r.Header().Get("Content-Type"); len(v) == 0 {
		r.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}

}

var CRLF = []byte("\r\n")

// Parse the response and send to wire(conn)
//
// closeConn decides whether to close the underlying connection
// i.e send Connection: close headers
func (r *Response) FinalizeResponse(closeConn bool, setAutoHeaders bool) {
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
	if setAutoHeaders {
		r.setAutoHeaders(closeConn)
	}

	// set contentLength
	//
	// Skip if switching protocol
	if r.status != StatusSwitchingProtocols {
		r.Header().Add("Content-Length", strconv.Itoa(len(r.body)))
	}

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

// TODO: Implement header validation (per RFC 7230)
func (r *Response) validateHeaderField(k string) bool {
	return true
}

// Writes bytes to the wire(conn)
func (r *Response) writeToWire(dataB []byte, dataS string) (n int, err error) {
	if dataB != nil {
		return r.w.Write(dataB)
	} else {
		return r.w.Write([]byte(dataS))
	}
}

func (r *Response) flush() error {
	return r.w.Flush()
}

func (r *Response) SetNotfoundHeader() {
	r.WriteHeader(StatusNotFound)
}

func (r *Response) SetInternalServerErrHeader() {
	r.WriteHeader(StatusInternalServerError)
}

func (r *Response) SetBadRequestHeader() {
	r.WriteHeader(StatusBadRequest)
}
