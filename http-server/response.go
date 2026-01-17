package main

// Response represesnts the response from an HTTP serverrequest
type Response struct {
	Status        string // e.g. "200 OK"
	StatusCode    int    // e.g. 200
	Protocol      string // e.g. "HTTP/1.1"
	ProtocolMajor string // e.g. 1
	ProtocolMinor string // e.g. 1

	Header Header

	// ContentLength records the length of the associated content.
	ContentLength int64

	// Request is the request that was sent to obtain
	// this response.
	Request *Request
}
