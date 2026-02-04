# HTTP/1.1 Server — Built From Scratch in Go (RFC 7230)

A **minimal, RFC-compliant HTTP/1.1 server implemented from scratch in Go**, without using `net/http` or any third-party HTTP libraries.

This project is part of **Week 3 of my “Networking From Scratch” learning track**, focused on deeply understanding **how HTTP actually works on the wire** — request parsing, response formatting, buffering, connection handling, and protocol correctness.

The server also supports **TLS**, allowing secure HTTPS connections using self-signed certificates.

---

## Motivation

High-level frameworks hide too much.

This project exists to answer questions like:

- What *exactly* does an HTTP request look like on the wire?
- How strict should HTTP request parsing be?
- How does a server correctly format HTTP/1.1 responses?
- When and how is `Content-Length` calculated?
- What happens if CRLF rules are violated?
- How does TLS fit underneath HTTP without changing application logic?

Instead of relying on Go’s `net/http`, **every layer is implemented manually**, guided directly by **RFC 7230**.

---

## Features Implemented

### HTTP/1.1 Request Parsing
- Parses **request-line** (`METHOD SP REQUEST-TARGET SP HTTP-VERSION CRLF`)
- Strict CRLF (`\r\n`) enforcement
- Rejects malformed request lines (wrong token count, missing CRLF)
- Header parsing with:
  - Case-insensitive header names
  - Support for multi-value headers
- Graceful handling of invalid requests with proper HTTP error responses

---

### Response Writing
- Custom `ResponseWriter` implementation
- Buffered writes for response body
- Automatic `Content-Length` calculation
- Proper response serialization order:
  - Status line
  - Headers
  - Blank line
  - Body

---

### Handler Interface (net/http-inspired)
Handlers follow a familiar signature:

```go
func(w ResponseWriter, r *Request)
```

This allows:
- Incremental writes via `w.Write()`
- Deferred response finalization
- Clean separation between parsing and response generation

---

### Connection Handling
- Manual TCP connection handling using `net.Conn`
- Buffered I/O using `bufio.Reader` and `bufio.Writer`
- Correct connection closing semantics
- Basic HTTP/1.1 keep-alive behavior

---

### TLS Support (HTTPS)
- HTTPS support using `crypto/tls`
- Runs HTTP over TLS without modifying HTTP logic
- Uses self-signed certificates for local development
- Demonstrates protocol layering:

```
HTTP → TLS → TCP
```

---

## Example Raw Request

```bash
printf "GET /health HTTP/1.1\r\nHost: localhost:8080\r\n\r\n" | nc localhost 8080
```

---

## Example Response

```http
HTTP/1.1 200 OK
Content-Length: 2
Content-Type: text/plain

OK
```

---

## Running the Server

### HTTP
```bash
go run ./cmd/server --port 8080
```

### HTTPS (TLS)
```bash
go run ./cmd/server --port 8443 --tls
```

> Browsers will show a warning for the self-signed certificate. This is expected for local development.

---

## RFC References

- RFC 7230 — HTTP/1.1 Message Syntax and Routing  
  https://datatracker.ietf.org/doc/html/rfc7230

---

## Key Learnings

- Why HTTP parsing must be **strict and defensive**
- How CRLF impacts request validity
- How response buffering simplifies handler logic
- Why `Content-Length` must be exact
- How browsers react to malformed responses
- How TLS operates *below* application protocols

---

## What This Project Is (and Isn’t)

### ✅ This Project Is
- Protocol-focused
- RFC-driven
- Designed for correctness and learning

### ❌ This Project Is Not
- A production-ready web framework
- Feature-complete (no chunked encoding, no HTTP/2)

---

## Next Steps

- Chunked transfer encoding
- Improved keep-alive handling
- Request body parsing
- Deeper TLS configuration and cipher exploration

---

## Author

Built by **Suman Mukherjee** as part of a deep dive into networking internals using Go.

