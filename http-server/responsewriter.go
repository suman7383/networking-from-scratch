package main

type ResponseWriter interface {
	// It returns the header map that will be sent by [ResponseWriter.WriteHeader].
	Header() Header

	// Write writes the data to the connection as part of an HTTP reply
	Write([]byte) (int, error)

	// WriteHeader sends an HTTP response header with the provided status code
	//
	// If WriteHeader is not called explicitly, the first call to Writewill trigger
	// an implicit WriteHeader(StatusOK)
	WriteHeader(statusCode int)
}
