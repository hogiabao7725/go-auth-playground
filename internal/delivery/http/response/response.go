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

// --- Helpers ---

// ValidationError - 422: Using for validation errors, with details of which fields failed validation
func ValidationError(w http.ResponseWriter, fields any) {
	sendJSON(w, http.StatusUnprocessableEntity, APIResponse{
		Success: false,
		Message: "validation error",
		Errors:  fields,
	})
}

// BadRequest - 400: Using for bad request errors, with details of the issue
func BadRequest(w http.ResponseWriter, msg string, details any) {
	if msg == "" {
		msg = "invalid request"
	}
	sendJSON(w, http.StatusBadRequest, APIResponse{
		Success: false,
		Message: msg,
		Errors:  details,
	})
}

// --- Success Responses ---

func Created(w http.ResponseWriter, data any, msg string) {
	sendJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Message: msg,
		Data:    data,
	})
}

func OK(w http.ResponseWriter, data any, msg string) {
	sendJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: msg,
		Data:    data,
	})
}

func InternalServerError(w http.ResponseWriter, msg string) {
	if msg == "" {
		msg = "internal server error"
	}
	sendJSON(w, http.StatusInternalServerError, APIResponse{
		Success: false,
		Message: msg,
	})
}

// Error - Function used for special error cases (401, 403, 409...)
func Error(w http.ResponseWriter, code int, msg string, details any) {
	sendJSON(w, code, APIResponse{
		Success: false,
		Message: msg,
		Errors:  details,
	})
}
