package response

import (
	"net/http"
	"vintage-server/pkg/utils"

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

// Success mengirimkan response sukses dengan data.
func Success[T any](c *gin.Context, statusCode int, data T) {
	c.JSON(statusCode, APIResponse[T]{Data: data})
}

// SuccessWithoutData mengirimkan response sukses tanpa data.
func SuccessWithoutData(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse[any]{Detail: &message})
}

// Error generic untuk custom error.
func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse[any]{Detail: &message})
}

// =================================================================================
// SHORTHAND UNTUK ERROR UMUM
// =================================================================================

func ErrorBadRequest(c *gin.Context, message ...string) {
	msg := "Invalid request"
	if len(message) > 0 {
		msg = message[0]
	}
	c.JSON(http.StatusBadRequest, APIResponse[any]{Detail: utils.Ptr(msg)})
}

func ErrorUnauthorized(c *gin.Context, message ...string) {
	msg := "Unauthorized"
	if len(message) > 0 {
		msg = message[0]
	}
	c.JSON(http.StatusUnauthorized, APIResponse[any]{Detail: utils.Ptr(msg)})
}

func ErrorForbidden(c *gin.Context, message ...string) {
	msg := "Forbidden"
	if len(message) > 0 {
		msg = message[0]
	}
	c.JSON(http.StatusForbidden, APIResponse[any]{Detail: utils.Ptr(msg)})
}

func ErrorInternalServer(c *gin.Context, message ...string) {
	msg := "Internal server error"
	if len(message) > 0 {
		msg = message[0]
	}
	c.JSON(http.StatusInternalServerError, APIResponse[any]{Detail: utils.Ptr(msg)})
}
