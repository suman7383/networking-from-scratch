# Minimal WebSocket Server (RFC 6455) — Built from Scratch in Go (with TLS)

A **minimal, RFC-compliant WebSocket server implemented from scratch in Go**, without using `net/http`’s WebSocket support or any third‑party libraries.

This project focuses on **protocol correctness**, **bit‑level framing**, and **real browser compatibility**, and now includes **TLS support (`wss://`)** to mirror real‑world production setups.

The server has been tested against:
- a custom Go WebSocket client
- modern browsers (Chrome / Firefox)
- both plain WebSocket (`ws://`) and secure WebSocket (`wss://`)

---

## Motivation

The goal of this project was to deeply understand **how WebSockets actually work** at the protocol level, and how they fit into the real network stack:

- TCP as a byte stream
- TLS as a transport‑layer security wrapper
- HTTP → WebSocket upgrade
- WebSocket frame format
- Masking rules
- Payload length encoding
- Control frames (PING / PONG / CLOSE)
- Graceful connection shutdown

Instead of relying on existing libraries, everything was implemented directly from **RFC 6455**, with TLS handled explicitly at the listener level.

---

## Features Implemented

### TLS (`wss://`)
- TLS enabled using Go’s `crypto/tls`
- Self‑signed certificate for local development
- Correct certificate handling for `localhost`
- Browser‑compatible secure WebSocket connections
- TLS cleanly layered below HTTP and WebSocket logic

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
  - `126` (16‑bit extended payload length)
- Correct **network byte order (big‑endian)** handling
- Strict validation of:
  - RSV bits
  - opcode validity
  - control‑frame constraints

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
  - client‑initiated close
  - server‑initiated close
- Proper CLOSE frame exchange before TCP shutdown
- Browser reports clean closure (`1000 Normal Closure`)

### Browser Compatibility
- Successfully tested with real browsers using the JavaScript `WebSocket` API
- Compatible with standard browser behavior (no extensions required)
- Works correctly over **TLS (`wss://`)**

---

## Explicitly Not Implemented (By Design)

The following features are intentionally **not supported** to keep the server minimal and focused:

- Fragmentation (`FIN = 0`, continuation frames)
- WebSocket extensions (RSV bits must be 0)
- Compression (`permessage-deflate`)
- 64‑bit payload lengths (`127` case)
- Subprotocol negotiation

The server **explicitly rejects** unsupported cases instead of silently accepting them.

---

## How to Run

### 1. Generate a local TLS certificate

```bash
openssl req -x509 -newkey rsa:2048 \
  -keyout certs/server.key \
  -out certs/server.crt \
  -days 365 \
  -nodes \
  -subj "/CN=localhost" \
  -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"
```

> The browser will warn about the certificate being self‑signed.  
> This is expected for local development.

---

### 2. Start the server

Run from the **project root**:

```bash
go run internal/server/server.go
```

The server listens on:

```
wss://localhost:8443
```

---

## Testing

### Browser Test (TLS)

1. First, visit in the browser:
   ```
   https://localhost:8443
   ```
   Accept the certificate warning.

2. Then open DevTools and run:

```js
const ws = new WebSocket("wss://localhost:8443");

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
- Secure connection opens successfully
- Messages are exchanged
- `ws.close()` results in a clean close (`1000 Normal`)

---

## Design Principles

- **Fail fast on protocol violations**
- **Validate as early as possible**
- **Never trust client input**
- **Exact byte‑level control**
- **Clear separation of layers (TCP → TLS → HTTP → WebSocket)**

All parsing and writing is done explicitly and incrementally.

---

## What This Project Demonstrates

- Understanding of TCP as a byte stream
- Correct layering of TLS over TCP
- Bit‑level protocol parsing and construction
- Careful RFC interpretation and enforcement
- Proper error handling and shutdown semantics
- Ability to build network protocols from scratch
- Real‑world interoperability with browsers over `wss://`

---

## Disclaimer

This project is intended for **learning and experimentation**.
It is not intended to replace production‑grade WebSocket servers.

---

## Author

Built as a learning project to deeply understand WebSocket internals, TLS integration, and network protocol implementation in Go.
