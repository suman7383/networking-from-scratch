# Minimal WebSocket Server (RFC 6455) — Built from Scratch in Go

A **minimal, RFC-compliant WebSocket server implemented from scratch in Go**, without using `net/http`’s WebSocket support or any third-party libraries.

This project focuses on **protocol correctness**, **bit-level framing**, and **real browser compatibility**, rather than feature completeness.

The server has been tested against:
- a custom Go WebSocket client
- modern browsers (Chrome / Firefox)

---

## Motivation

The goal of this project was to deeply understand **how WebSockets actually work** at the protocol level:

- HTTP → WebSocket upgrade
- WebSocket frame format
- Masking rules
- Payload length encoding
- Control frames (PING / PONG / CLOSE)
- Graceful connection shutdown

Instead of relying on existing libraries, everything was implemented directly from **RFC 6455**.

---

## Features Implemented

### HTTP Handshake
- Full HTTP/1.1 request parsing
- Strict CRLF handling
- WebSocket upgrade validation
- Correct `Sec-WebSocket-Accept` computation
- Rejects invalid or malformed upgrade requests

### WebSocket Framing (RFC 6455)
- FIN bit parsing and generation
- Opcode handling:
  - TEXT
  - BINARY
  - PING
  - PONG
  - CLOSE
- Payload length handling:
  - `< 126`
  - `126` (16-bit extended payload length)
- Correct **network byte order (big-endian)** handling
- Strict validation of:
  - RSV bits
  - opcode validity
  - control frame constraints

### Masking
- Enforces **client → server masking**
- Reads and applies masking keys correctly
- Server → client frames are **unmasked**, per spec

### Control Frames
- Proper handling of:
  - PING → PONG
  - CLOSE handshake
- Control frames are never fragmented
- Payload size limits enforced for control frames

### Graceful Close Handshake
- Supports both:
  - client-initiated close
  - server-initiated close
- Proper CLOSE frame exchange before TCP shutdown
- Browser reports clean closure (`1000 Normal Closure`)

### Browser Compatibility
- Successfully tested with real browsers using the JavaScript `WebSocket` API
- Compatible with standard browser behavior (no extensions required)

---

## Explicitly Not Implemented (By Design)

The following features are intentionally **not supported** to keep the server minimal and focused:

- Fragmentation (`FIN = 0`, continuation frames)
- WebSocket extensions (RSV bits must be 0)
- Compression (`permessage-deflate`)
- 64-bit payload lengths (`127` case)
- Subprotocol negotiation

The server **explicitly rejects** unsupported cases instead of silently accepting them.

---

## How to Run

```bash
go run main.go
```

The server listens on:

```
ws://localhost:8080
```

---

## Testing

### Browser Test

Open a regular web page (not `chrome://`) and run in DevTools:

```js
const ws = new WebSocket("ws://localhost:8080");

ws.onopen = () => {
  console.log("CONNECTED");
  ws.send("hello from browser");
};

ws.onmessage = (e) => {
  console.log("MESSAGE FROM SERVER:", e.data);
};

ws.onclose = (e) => {
  console.log("CLOSED:", e.code, e.reason);
};
```

Expected behavior:
- Connection opens successfully
- Messages are exchanged
- `ws.close()` results in a clean close (`1000 Normal`)

---

## Design Principles

- **Fail fast on protocol violations**
- **Validate as early as possible**
- **Never trust client input**
- **Exact byte-level control**
- **No partial or ambiguous frame reads**

All parsing and writing is done explicitly and incrementally.

---

## What This Project Demonstrates

- Understanding of TCP as a byte stream
- Bit-level protocol parsing and construction
- Careful RFC interpretation and enforcement
- Proper error handling and shutdown semantics
- Ability to build network protocols from scratch
- Real-world interoperability with browsers

---

## Disclaimer

This project is intended for **learning and experimentation**.
It is not intended to replace production-grade WebSocket servers.

---

## Author

Built as a learning project to deeply understand WebSocket internals and network protocol implementation in Go.
