package main

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	StatusCode int               `json:"status_code"`
	Message    string            `json:"message"`
	Fields     map[string]string `json:"fields,omitempty"`
}

func NewAPIError(message string, statusCode int) *APIError {
	return &APIError{Message: message, StatusCode: statusCode}
}

func (api *APIError) Error() string {
	return api.Message
}
func (api *APIError) respond(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(api.StatusCode)
	json.NewEncoder(w).Encode(api)
}

func badRequest(message string) *APIError {
	return &APIError{Message: message, StatusCode: http.StatusBadRequest}
}

func notFound(message string) *APIError {
	return &APIError{Message: message, StatusCode: http.StatusNotFound}
}

func internalError(message string) *APIError {
	return &APIError{Message: message, StatusCode: http.StatusInternalServerError}
}
func validateError(fields map[string]string) *APIError {
	return &APIError{Message: "Validation failed", StatusCode: http.StatusUnprocessableEntity, Fields: fields}
}

func validationError(fields map[string]string) *APIError {
	return &APIError{
		StatusCode: http.StatusUnprocessableEntity,
		Message:    "validation failed",
		Fields:     fields,
	}
}

func unauthorizedError(message string) *APIError {
	return &APIError{StatusCode: http.StatusUnauthorized, Message: message}
}
