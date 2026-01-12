package main

import "context"

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

	// ctx is the server context.
	ctx context.Context
}
