package main

import "strings"

// It represents the key-value pairs in an HTTP header
type Header map[string][]string

// Adds key, value pait to header
//
// The key is case-insensitive
func (h Header) Add(key, value string) {

	// Convert key to lower case
	kl := strings.ToLower(key)
	if _, ok := h[kl]; !ok {
		h[kl] = make([]string, 0)
	}

	// Append the value at the given key
	h[kl] = append(h[kl], value)
}

// Values returns all the values associated with the given key.
//
// It is case insensitive
func (h Header) Values(key string) []string {
	kl := strings.ToLower(key)

	if h == nil {
		return nil
	}

	return h[kl]
}

// Del deletes the values associated with the key
func (h Header) Del(key string) {
	delete(h, strings.ToLower(key))
}
