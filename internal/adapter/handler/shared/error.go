package shared_handler

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
	if appErr, ok := domain.IsAppError(err); ok {
		return handleAppError(c, appErr)
	}

	log.Printf("Unhandled error: %v", err)
	return c.Status(http.StatusInternalServerError).JSON(ErrorResponse{
		Message: "An unexpected error occurred",
		Type:    string(domain.ErrorTypeInternal),
	})
}

func handleAppError(c fiber.Ctx, err *domain.AppError) error {
	statusCode := errorTypeToStatusCode(err.Type)

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

func errorTypeToStatusCode(errType domain.ErrorType) int {
	switch errType {
	case domain.ErrorTypeBadRequest:
		return http.StatusBadRequest
	case domain.ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case domain.ErrorTypeForbidden:
		return http.StatusForbidden
	case domain.ErrorTypeNotFound:
		return http.StatusNotFound
	case domain.ErrorTypeConflict:
		return http.StatusConflict
	case domain.ErrorTypeValidation:
		return http.StatusUnprocessableEntity
	case domain.ErrorTypeInternal:
		fallthrough
	default:
		return http.StatusInternalServerError
	}
}
