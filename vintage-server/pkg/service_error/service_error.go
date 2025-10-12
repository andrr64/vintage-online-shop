package serviceerror

import "fmt"

type Code string

const (
	ErrInternal     Code = "internal_error"
	ErrBadRequest   Code = "bad_request"
	ErrConflict     Code = "conflict"
	ErrNotFound     Code = "not_found"
	ErrUnauthorized Code = "unauthorized" // ditambahkan
)

// ServiceError adalah error kustom untuk service
type ServiceError struct {
	Code    Code
	Message string
	Err     error // optional, untuk wrapping original error
}

func (e *ServiceError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s | %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func New(code Code, message string, err error) *ServiceError {
	return &ServiceError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (e *ServiceError) IsInternal() bool {
	return e.Code == ErrInternal
}

func (e *ServiceError) IsBadRequest() bool {
	return e.Code == ErrBadRequest
}

func (e *ServiceError) IsConflict() bool {
	return e.Code == ErrConflict
}

func (e *ServiceError) IsNotFound() bool {
	return e.Code == ErrNotFound
}

func (e *ServiceError) IsUnauthorized() bool {
	return e.Code == ErrUnauthorized
}
