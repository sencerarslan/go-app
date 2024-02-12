package utils

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  int         `json:"status"`
}

func (r *Response) Send(res *gin.Context) {
	res.JSON(r.Status, gin.H{
		"data":    r.Data,
		"message": r.Message,
	})
}

type APIError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func (e *APIError) Send(res *gin.Context) {
	res.JSON(e.StatusCode, gin.H{"message": e.Message})
}

func NewAPIError(message string, statusCode int) *APIError {
	return &APIError{
		Message:    message,
		StatusCode: statusCode,
	}
}
