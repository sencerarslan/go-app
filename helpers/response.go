package helper

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Success bool        `json:"success"`
}

func NewResponse(data interface{}, message string, success bool) Response {
	return Response{
		Data:    data,
		Message: message,
		Success: success,
	}
}

func (r Response) SendJSON(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(r)
}

func SuccessResponse(data interface{}, message string) Response {
	if message == "" {
		message = "Operation successful"
	}
	return NewResponse(data, message, true)
}

func CreatedResponse(data interface{}, message string) Response {
	if message == "" {
		message = "Operation successful"
	}
	return NewResponse(data, message, true)
}

func ErrorResponse(data interface{}, message string) Response {
	if message == "" {
		message = "Operation failed"
	}
	return NewResponse(data, message, false)
}

func UnauthorizedResponse(data interface{}, message string) Response {
	if message == "" {
		message = "Please log in"
	}
	return NewResponse(data, message, false)
}

func NotFoundResponse(data interface{}, message string) Response {
	if message == "" {
		message = "Operation failed"
	}
	return NewResponse(data, message, false)
}

func TooManyRequestsResponse(data interface{}, message string) Response {
	if message == "" {
		message = "Too many requests received"
	}
	return NewResponse(data, message, false)
}
