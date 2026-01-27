package websocket

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"net"

	"github.com/suman7383/networking-from-scratch/websocket-server/internal/http"
	"github.com/suman7383/networking-from-scratch/websocket-server/utils"
)

var ErrBadHandshake = errors.New("Bad handshake")
var ErrClientVersion = errors.New("Unsupported version")

var guid = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

var swsvk = "Sec-WebSocket-Version"
var swsak = "Sec-WebSocket-Accept"
var swsk = "Sec-WebSocket-Key"

func HandleHandshake(req *http.Request, conn net.Conn) (WebsocketConn *WebSocketConn, err error) {
	// TODO: Validate WebSocket-only headers\
	key, err := validateHeaders(req)
	if err != nil {
		// Send HTTP error response
		utils.WriteErrResponse(conn, http.StatusBadRequest, err.Error())

		return nil, err
	}

	// Compute Sec-WebSocket-Accept
	swsa := computeWebsocketAccept(key)

	// Write 101 Switching Protocols response
	sendSwitchingProtoResponse(swsa, conn)

	// Take ownership of the connection and create WebsocketConn
	wsc := &WebSocketConn{
		conn: conn,
		r:    NewFrameReader(conn),
		w:    NewFrameWriter(conn),
	}

	return wsc, nil
}

// Sends 101 Switching Protocols response
//
// Adds Sec-WebSocket-Accept: <value> header
func sendSwitchingProtoResponse(swsa string, conn net.Conn) {
	res := http.NewResponse(conn)

	// Set Sec-WebSocket-Accept and Sec-WebSocket-Version: 13 header
	res.Header()[swsak] = []string{swsa}
	res.Header().Set("Connection", "Upgrade")
	res.Header().Set("Upgrade", "websocket")

	// Set 101 status
	res.WriteHeader(http.StatusSwitchingProtocols)

	res.FinalizeResponse(false, false)
}

func computeWebsocketAccept(key string) string {
	// Concatenate the key and guid
	t := key + guid

	h := sha1.New()
	h.Write([]byte(t))
	sha1Bytes := h.Sum(nil)

	return base64.StdEncoding.EncodeToString(sha1Bytes)
}

func validateHeaders(req *http.Request) (key string, err error) {
	// Sec-WebSocket-Key, Sec-WebSocket-Version
	//
	// We validate only the above two as for now we don't care about others
	key = req.Header.Get(swsk)

	// Checks if Version is strictly 13(and not list of versions)
	ver13 := func() bool {

		if v := req.Header.Values(swsvk); len(v) == 1 {
			return v[0] == "13"
		} else {
			return false
		}
	}()

	if len(key) == 0 {
		return key, ErrBadHandshake
	}

	if !ver13 {
		return key, ErrClientVersion
	}

	return key, nil
}

var CRLF = []byte("\r\n")
