package response_helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"vintage-server/dto"
)

// Success response generic
func Success[T any](c *gin.Context, data *T, message string) {
	c.JSON(http.StatusOK, dto.CommonResponse[T]{
		Message: message,
		Success: true,
		Data:    data,
	})
}

// Failed response generic
func Failed[T any](c *gin.Context, httpCode int, message string, data *T) {
	c.JSON(httpCode, dto.CommonResponse[T]{
		Message: message,
		Success: false,
		Data:    data,
	})
}
