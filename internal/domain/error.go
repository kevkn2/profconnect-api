package domain

import "fmt"

// ErrorType defines the category of an error
type ErrorType string

const (
	// 4xx Client Errors
	ErrorTypeBadRequest   ErrorType = "bad_request"      // 400
	ErrorTypeUnauthorized ErrorType = "unauthorized"     // 401
	ErrorTypeForbidden    ErrorType = "forbidden"        // 403
	ErrorTypeNotFound     ErrorType = "not_found"        // 404
	ErrorTypeConflict     ErrorType = "conflict"         // 409
	ErrorTypeValidation   ErrorType = "validation_error" // 422

	// 5xx Server Errors
	ErrorTypeInternal ErrorType = "internal_error" // 500
)

// AppError represents a domain-aware application error
type AppError struct {
	Type    ErrorType
	Message string
	Err     error // Original error for logging
}

// NewAppError creates a new application error
func NewAppError(errorType ErrorType, message string) *AppError {
	return &AppError{
		Type:    errorType,
		Message: message,
		Err:     nil,
	}
}

// NewAppErrorWithErr creates a new application error with the underlying error
func NewAppErrorWithErr(errorType ErrorType, message string, err error) *AppError {
	return &AppError{
		Type:    errorType,
		Message: message,
		Err:     err,
	}
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}

// Common error constructors for convenience
func BadRequest(message string) *AppError {
	return NewAppError(ErrorTypeBadRequest, message)
}

func BadRequestErr(message string, err error) *AppError {
	return NewAppErrorWithErr(ErrorTypeBadRequest, message, err)
}

func Unauthorized(message string) *AppError {
	return NewAppError(ErrorTypeUnauthorized, message)
}

func UnauthorizedErr(message string, err error) *AppError {
	return NewAppErrorWithErr(ErrorTypeUnauthorized, message, err)
}

func Forbidden(message string) *AppError {
	return NewAppError(ErrorTypeForbidden, message)
}

func NotFound(message string) *AppError {
	return NewAppError(ErrorTypeNotFound, message)
}

func Conflict(message string) *AppError {
	return NewAppError(ErrorTypeConflict, message)
}

func ValidationError(message string) *AppError {
	return NewAppError(ErrorTypeValidation, message)
}

func Internal(message string) *AppError {
	return NewAppError(ErrorTypeInternal, message)
}

func InternalErr(message string, err error) *AppError {
	return NewAppErrorWithErr(ErrorTypeInternal, message, err)
}
