// Package handlers contains HTTP request handlers for the TCG API
package handlers

// ErrorResponse
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
