package response

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func sendJSON(w http.ResponseWriter, code int, payload APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func Success(w http.ResponseWriter, code int, message string, data any) {
	sendJSON(w, code, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, code int, message string, errors any) {
	sendJSON(w, code, APIResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}
