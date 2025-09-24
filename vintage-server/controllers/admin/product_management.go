package admin

import (
	"net/http"
	"vintage-server/helpers/response_helper"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {

	_, exists := c.Get("currentUser")
	if !exists {
		response_helper.Failed[string](c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	response_helper.Success[any](c, nil, "OK")
}
