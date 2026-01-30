package httpcore

import (
	"net/textproto"
	"strings"
)

// It represents the key-value pairs in an HTTP header
type Header map[string][]string

// Adds key, value pair to header
//
// The key is case-insensitive
func (h Header) Add(key, value string) {

	// Convert key to canonicalMIMEHeaderKey
	kl := textproto.CanonicalMIMEHeaderKey(key)
	if _, ok := h[kl]; !ok {
		h[kl] = make([]string, 0)
	}

	// Append the value at the given key
	h[kl] = append(h[kl], value)
}

// Sets key to the new value
func (h Header) Set(key, value string) {

	kl := textproto.CanonicalMIMEHeaderKey(key)

	// Delete any previous entry
	delete(h, kl)

	h[kl] = []string{value}
}

// Get returns the first value associated with the given key
//
// If no value for a key, it returns empty string
func (h Header) Get(key string) string {
	kl := textproto.CanonicalMIMEHeaderKey(key)

	if _, ok := h[kl]; !ok {
		return ""
	}

	return h[kl][0]
}

// Values returns all the values associated with the given key.
//
// It is case insensitive
func (h Header) Values(key string) []string {
	kl := textproto.CanonicalMIMEHeaderKey(key)

	if h == nil {
		return nil
	}

	return h[kl]
}

// Del deletes the values associated with the key
func (h Header) Del(key string) {
	delete(h, strings.ToLower(key))
}
