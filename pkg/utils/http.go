package utils

import (
	"encoding/json"
	"net/http"
)

// WriteJSON is a utility method that encodes the response as JSON
// sets appropriate "Content-Type" response header and writes full response
func WriteJSON(w http.ResponseWriter, data interface{}, statusCode int) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	return err // forward error (usually it's nil)
}

// MustWriteJSON that calls WriteJSON and checks error
// Note, error (if any) just printed to std log!
// Don't be confused with "Must" prefix, it doesn't panic.
func MustWriteJSON(w http.ResponseWriter, response interface{}, statusCode int) {
	IgnoreError(WriteJSON(w, response, statusCode))
}
