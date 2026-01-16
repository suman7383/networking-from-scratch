package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Represents a HTTP request received by the server.
type Request struct {
	// It specifies the HTTP methods(GET, POST, PUT, etc).
	Method string

	// URL specifies the URI being requested
	//
	// The URL is parsed from the URI on the Request-Line (See RFC 7230, Section 5.3)
	//
	// TODO: Change this to type *url.URL(create an url struct)
	url string

	Protocol      string // "HTTP/1.1"
	ProtocolMajor int    // 1
	ProtocolMinor int    // 0

	// It contains the request header fields received by
	// the server
	//
	// If a server received a request with header lines,
	//
	// Host: example.com
	// accept-encoding: gzip, deflate
	// fOO: Bar
	// foo: two
	//
	// then
	//
	// Header = map[string][]string{
	// 	"Accept-Encoding": {"gzip, deflate"},
	// 	"Foo": {"Bar", "two"}
	// }
	Header Header

	// Host specifies the host on which the URL is sought. For
	// HTTP/1(per RFC 7230, section 5.4), this is either the value
	// of the "Host" header or the host name givent in the URL itself
	//
	Host string

	// RequestURI is the unmodified request-target of the
	// Request-Line (per RFC 7230, section 3.1.1) as sent by
	// the client to a server.
	RequestURI string

	// Path specifies the URI path for the request("/", "/health" )
	//
	// TODO: Move this to url package
	Path string

	// ctx is the server context.
	ctx context.Context
}

var ErrMalformedRequestLine = errors.New("malformed request line.")
var ErrInvalidRequestMethod = errors.New("method invalid or not supported. Only send GET request")

func badStringError(err, val string) error { return fmt.Errorf("%s %q", err, val) }

func readRequest(r *Reader) (req *Request, err error) {
	req = new(Request)

	// HTTP request-line = method SP request-target SP HTTP-version CRLF
	// Where SP = Single Space
	var reqLine string
	reqLine, err = r.ReadLine()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	var ok bool
	req.Method, req.RequestURI, req.Protocol, ok = parseRequestLine(reqLine)
	if !ok {
		return nil, badStringError("malformed HTTP request", reqLine)
	}

	if len(req.RequestURI) == 0 {
		return nil, ErrMalformedRequestLine
	}

	if !validMethod(req.Method) {
		return nil, ErrInvalidRequestMethod
	}

	// TODO
	// parse url from req.RequestURI
	fmt.Printf("RequestURI is %s", req.RequestURI)

	// parse http version
	if req.ProtocolMajor, req.ProtocolMinor, ok = parseHttpVersion(req.Protocol); !ok {
		return nil, badStringError("malformed HTTP version", req.Protocol)
	}

	// Parse headers
	// header-field   = field-name ":" OWS field-value OWS  (Where OWS = Optional White Space)
	req.Header, err = parseHeaders(r)
	if err != nil {
		return nil, err
	}

	return req, nil
}

var ErrInvalidHeaderField = errors.New("invalid header field")

func parseHeaders(r *Reader) (Header, error) {
	h := make(Header)

	// Sample request(after request line)
	// Host: localhost:8080\r\nUser-Agent: curl/8.0.0\r\nAccept: */*\r\nConnection: close\r\n\r\n
	for {
		line, err := r.ReadLine()
		if err != nil {
			return nil, err
		}

		// Reached end of headers
		if len(line) == 0 {
			return h, nil
		}

		// Check for leading space
		if line[0] == ' ' || line[0] == '\t' {
			return nil, ErrInvalidHeaderField
		}

		// line contains FieldName: Value
		// Check for colon and a OWS before value
		// TODO
		k, v, found := strings.Cut(line, ":")

		if !found {
			return nil, ErrInvalidHeaderField
		}

		// Check for trailing space in key
		if len(k) == 0 || k[len(k)-1] == ' ' || k[len(k)-1] == '\t' || strings.Contains(k, " ") {
			return nil, ErrInvalidHeaderField
		}

		vsr := strings.TrimSpace(v)
		h.Add(k, vsr)
	}
}

func parseHttpVersion(protocol string) (majorProto, minorProto int, ok bool) {
	// HTTP/1.1
	switch protocol {
	case "HTTP/1.1":
		return 1, 1, true
	case "HTTP/1.0":
		return 1, 0, true
	default:
		return 0, 0, false
	}

	// if !strings.HasPrefix(protocol, "HTTP/"){
	// 	return 0, 0, false
	// }
	// if len(protocol) != len("HTTP/X.Y"){
	// 	return 0, 0, false
	// }
	// if protocol[6] != '.'{
	// 	return 0, 0, false
	// }
}

// NewRequest forms and returns *Request using the requst line
func parseRequestLine(requestLine string) (method, requestUri, proto string, ok bool) {
	// requestLine = method SP request-target SP HTTP-version
	// We cut on SP(Single Space)
	method, rest, ok1 := strings.Cut(requestLine, " ")
	requestURI, proto, ok2 := strings.Cut(rest, " ")
	if !ok1 || !ok2 {
		return "", "", "", false
	}

	return method, requestURI, proto, true
}

func validMethod(m string) bool {
	// For now we only check for GET request
	if m != MethodGet {
		return false
	}
	return true
}
