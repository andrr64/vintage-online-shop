// File: pkg/utils/pointer.go
package utils

func Ptr[T any](v T) *T {
	return &v
}
