package helper

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Ambil param sebagai int32
func GetParamInt32(c *gin.Context, param string) (int32, error) {
	val := c.Param(param)
	if val == "" {
		return 0, fmt.Errorf("param %s is empty", param)
	}

	v, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return 0, err
	}

	return int32(v), nil
}

func GetParamInt(c *gin.Context, param string) (int, error) {
	val := c.Param(param)
	if val == "" {
		return 0, fmt.Errorf("param %s is empty", param)
	}

	v, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("invalid param %s: %v", param, err)
	}

	return v, nil
}

// Ambil param sebagai string
func GetParamString(c *gin.Context, param string) (string, error) {
	val := c.Param(param)
	if val == "" {
		return "", fmt.Errorf("param %s is empty", param)
	}
	return val, nil
}

// Ambil param sebagai int64
func GetParamInt64(c *gin.Context, param string) (int64, error) {
	val := c.Param(param)
	if val == "" {
		return 0, fmt.Errorf("param %s is empty", param)
	}

	v, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}

	return v, nil
}

func GetParamUUID(c *gin.Context, param string) (uuid.UUID, error) {
	val := c.Param(param)
	if val == "" {
		return uuid.Nil, fmt.Errorf("param %s is empty", param)
	}

	id, err := uuid.Parse(val)
	if err != nil {
		return uuid.Nil, fmt.Errorf("param %s is not a valid UUID", param)
	}

	return id, nil
}

// GetQueryInt mengambil query param sebagai int, jika kosong atau invalid pakai default
func GetQueryInt(c *gin.Context, key string, defaultVal int) int {
	valStr := c.Query(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}

// GetQueryInt64 mengambil query param sebagai int64, jika kosong atau invalid pakai default
func GetQueryInt64(c *gin.Context, key string, defaultVal int64) int64 {
	valStr := c.Query(key)
	if valStr == "" {
		return defaultVal
	}

	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return defaultVal
	}

	return val
}
