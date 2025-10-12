package common

type Pagination[T any] struct {
	Page       int   `json:"page"`
	Size       int   `json:"size"`
	TotalItems int   `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	Items      []T   `json:"items"`
}
