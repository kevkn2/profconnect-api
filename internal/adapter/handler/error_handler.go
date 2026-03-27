package handler

import (
	"log"
	"net/http"
	"profconnect-api/internal/domain"

	"github.com/gofiber/fiber/v3"
)

// ErrorResponse is the standard error response format
type ErrorResponse struct {
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
}

// HandleError converts domain errors to HTTP responses
func HandleError(c fiber.Ctx, err error) error {
	// Check if it's an AppError
	if appErr, ok := domain.IsAppError(err); ok {
		return handleAppError(c, appErr)
	}

	// Handle generic errors as internal server errors
	log.Printf("Unhandled error: %v", err)
	return c.Status(http.StatusInternalServerError).JSON(ErrorResponse{
		Message: "An unexpected error occurred",
		Type:    string(domain.ErrorTypeInternal),
	})
}

// handleAppError maps AppError to HTTP status codes
func handleAppError(c fiber.Ctx, err *domain.AppError) error {
	statusCode := errorTypeToStatusCode(err.Type)

	// Log errors with context
	if err.Err != nil {
		log.Printf("[%s] %v (caused by: %v)", err.Type, err.Message, err.Err)
	} else {
		log.Printf("[%s] %v", err.Type, err.Message)
	}

	return c.Status(statusCode).JSON(ErrorResponse{
		Message: err.Message,
		Type:    string(err.Type),
	})
}

// errorTypeToStatusCode maps domain error types to HTTP status codes
func errorTypeToStatusCode(errType domain.ErrorType) int {
	switch errType {
	case domain.ErrorTypeBadRequest:
		return http.StatusBadRequest // 400
	case domain.ErrorTypeUnauthorized:
		return http.StatusUnauthorized // 401
	case domain.ErrorTypeForbidden:
		return http.StatusForbidden // 403
	case domain.ErrorTypeNotFound:
		return http.StatusNotFound // 404
	case domain.ErrorTypeConflict:
		return http.StatusConflict // 409
	case domain.ErrorTypeValidation:
		return http.StatusUnprocessableEntity // 422
	case domain.ErrorTypeInternal:
		fallthrough
	default:
		return http.StatusInternalServerError // 500
	}
}
