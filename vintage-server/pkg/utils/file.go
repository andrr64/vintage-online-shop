package utils

import "mime/multipart"

func SizeIsOk(fileHeader *multipart.FileHeader, max_size int64) bool {
	if fileHeader == nil {
		return false
	}
	return fileHeader.Size <= max_size
}

func Megabytes(mb int64) int64 {
	return mb * 1024 * 1024
}
