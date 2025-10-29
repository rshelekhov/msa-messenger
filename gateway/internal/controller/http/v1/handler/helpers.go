package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

func handleError(w http.ResponseWriter, r *http.Request, err error, statusCode int) {
	render.Status(r, statusCode)
	render.JSON(w, r, ErrorResponse{Error: err.Error()})
}

func handleInternalError(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, ErrorResponse{Error: ErrInternalServerError.Error()})
}

func handleBadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, ErrorResponse{Error: err.Error()})
}

func handleValidationErrors(w http.ResponseWriter, r *http.Request, log *slog.Logger, err error) {
	var validationErrors validator.ValidationErrors

	ok := errors.As(err, &validationErrors)
	if !ok {
		// If this is not a validation error, return a general error
		log.Warn("failed to validate request", slog.String("error", err.Error()))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{Error: ErrInvalidRequest.Error()})
		return
	}

	fieldErrors := processValidationErrors(validationErrors)

	log.Warn("failed to validate request", slog.Any("validation_errors", fieldErrors))

	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, ValidationErrorResponse{
		Error:   ErrValidationFailed.Error(),
		Details: fieldErrors,
	})
}

func processValidationErrors(errors validator.ValidationErrors) map[string]string {
	fieldErrors := make(map[string]string)

	for _, err := range errors {
		field := err.Field()

		switch err.Tag() {
		case "required":
			fieldErrors[field] = "field is required"
		case "email":
			fieldErrors[field] = "invalid email format"
		case "min":
			fieldErrors[field] = fmt.Sprintf("must be at least %s characters", err.Param())
		case "max":
			fieldErrors[field] = fmt.Sprintf("must be at most %s characters", err.Param())
		default:
			fieldErrors[field] = "invalid value"
		}
	}

	return fieldErrors
}

func handleMappedError(w http.ResponseWriter, r *http.Request, err error, log *slog.Logger, operation string, errorMappings map[error]int) bool {
	if err == nil {
		return false
	}

	for domainErr, status := range errorMappings {
		if errors.Is(err, domainErr) {
			log.Warn(operation, slog.String("error", err.Error()))
			handleError(w, r, domainErr, status)
			return true
		}
	}

	// No match - let caller handle default case
	return false
}

func isMobileRequest(r *http.Request) bool {
	return r.Header.Get("X-Client-Type") == string(ClientTypeMobile)
}
