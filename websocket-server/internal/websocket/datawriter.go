package websocket

type DataWriter interface {
	Send(data []byte)
}
