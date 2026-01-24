package http

const (
	StatusSwitchingProtocols = 101

	StatusOK = 200

	StatusBadRequest = 400
	StatusNotFound   = 404

	StatusInternalServerError = 500
)

func StatusText(code int) string {
	switch code {
	case StatusSwitchingProtocols:
		return "Switching Protocols"
	case StatusOK:
		return "OK"
	case StatusBadRequest:
		return "Bad Request"
	case StatusNotFound:
		return "Not Found"
	case StatusInternalServerError:
		return "Internal Server Error"
	default:
		return ""
	}
}
