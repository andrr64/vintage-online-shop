package response

import (
	"github.com/gin-gonic/gin"
)

// APIResponse adalah struct generic untuk semua response JSON dari API kita.
type APIResponse[T any] struct {
	Data   T       `json:"data,omitempty"`
	Detail *string `json:"detail,omitempty"`
}

// =================================================================================
// HELPER FUNCTIONS - Untuk membuat response di handler jadi lebih bersih
// =================================================================================

// Success mengirimkan response 200 OK dengan data.
func Success[T any](c *gin.Context, statusCode int, data T) {
	c.JSON(statusCode, APIResponse[T]{Data: data})
}

// Error mengirimkan response error dengan pesan detail.
func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse[any]{Detail: &message})
}

func SuccessWithoutData(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse[any]{Detail: &message})
}