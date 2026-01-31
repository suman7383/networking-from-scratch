package websocket

type DataType string

const (
	DataTypeText   DataType = "text"
	DataTypeBinary DataType = "binary"
)

type DataWriter interface {
	// Data to send and the type of data(text/binary)
	Send(data []byte, dt DataType)
}
