# Networking From Scratch (Go)

Low-level networking servers implemented from scratch in Go to understand TCP, HTTP, and WebSockets without frameworks.

---

## 📌 Overview

This repository is a hands-on exploration of how network servers work **under the hood**.  
Instead of relying on high-level abstractions, each project builds directly on top of Go’s `net` package to understand:

- How TCP connections behave
- How protocols frame and parse data
- Why HTTP and WebSockets are designed the way they are
- How concurrency and backpressure affect servers

Each directory contains a **standalone binary** focused on one concept.

---

## 🎯 Goals

- Understand TCP as a **byte stream**, not messages
- Learn protocol **framing and parsing**
- Implement HTTP without `net/http`
- Understand WebSocket handshakes and frames
- Build intuition for concurrency and connection handling

This is a **learning-first repository**, not a production-ready framework.

---

## 📂 Projects

### 1️⃣ TCP Echo Server
**Location:** `tcp-echo/`

A minimal TCP server that echoes back whatever the client sends.

**Concepts learned**
- TCP connection lifecycle
- Blocking I/O
- Partial reads and writes
- Goroutine-per-connection model

---

### 2️⃣ TCP Chat Server
**Location:** `tcp-chat/`

A multi-client TCP chat server where messages from one client are broadcast to all others.

**Concepts learned**
- Managing multiple concurrent connections
- Shared state and coordination
- Fan-out message broadcasting
- Handling slow or disconnected clients

---

### 3️⃣ HTTP Server (No Framework)
**Location:** `http-server/`

A minimal HTTP/1.1 server implemented directly over TCP, without using `net/http`.

**Features**
- Parses request line and headers
- Supports basic routing
- Returns proper HTTP responses
- Handles `Content-Length`

**Concepts learned**
- HTTP is just text over TCP
- Request/response framing
- Header parsing
- Persistent connections (optional extension)

---

### 4️⃣ Minimal WebSocket Server
**Location:** `websocket/`

A basic WebSocket server implementing the HTTP upgrade handshake and frame parsing.

**Concepts learned**
- HTTP → WebSocket upgrade
- SHA1 + Base64 handshake
- WebSocket frame format
- Masking rules
- Binary protocol parsing

---

## 🧠 Core Concepts Covered

- TCP streams vs messages
- Framing strategies (delimiter-based, length-based)
- Protocol parsing
- Connection lifecycle management
- Concurrency with goroutines and channels
- Backpressure and slow clients

---

## 📚 Learning Resources

These resources were used while building the projects:

- Beej’s Guide to Networking — https://beej.us/guide/bgnet/
- Go `net` package documentation
- HTTP/1.1 RFC (RFC 7230)
- HTTP/2 Explained — https://http2.github.io/
- High Performance Browser Networking — https://hpbn.co/

---

## 🧪 How to Run

Each project is independent.

```bash
cd tcp-echo
go run .
