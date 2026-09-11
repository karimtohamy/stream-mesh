package payload

import (
	"github.com/gin-gonic/gin"
)

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Meta struct {
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
	NextPage   int `json:"next_page,omitempty"`
	PrevPage   int `json:"prev_page,omitempty"`
	FirstPage  int `json:"first_page,omitempty"`
	LastPage   int `json:"last_page,omitempty"`
}

type Response struct {
	Success bool       `json:"success"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
	Meta    *Meta      `json:"meta,omitempty"`
}

func OK(c *gin.Context, statusCode int, data any) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
	})
}

func Fail(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}
