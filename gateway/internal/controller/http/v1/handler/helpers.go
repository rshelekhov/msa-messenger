package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

func handleError(w http.ResponseWriter, r *http.Request, log *slog.Logger, err error, statusCode int) {
	errMsg := err.Error()
	log.Error(errMsg)

	render.Status(r, statusCode)
	render.JSON(w, r, ErrorResponse{Error: &errMsg})
}

func handleInternalError(w http.ResponseWriter, r *http.Request, log *slog.Logger, err error) {
	errMsg := err.Error()
	log.Error(errMsg)

	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, ErrorResponse{Error: &errMsg})
}

func handleBadRequestError(w http.ResponseWriter, r *http.Request, log *slog.Logger, err error) {
	errMsg := err.Error()
	log.Error(errMsg)

	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, ErrorResponse{Error: &errMsg})
}

func handleValidationErrors(w http.ResponseWriter, r *http.Request, log *slog.Logger, err error) {
	var validationErrors validator.ValidationErrors

	ok := errors.As(err, &validationErrors)
	if !ok {
		// If this is not a validation error, return a general error
		errMsg := fmt.Errorf("failed to validate request: %w", err).Error()
		log.Error(errMsg)

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{Error: &errMsg})
		return
	}

	errMsg := fmt.Errorf("validation error: %s", processValidationErrors(validationErrors)).Error()
	log.Error(errMsg)

	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, ErrorResponse{Error: &errMsg})
}

func processValidationErrors(errors validator.ValidationErrors) string {
	var errMsgs []string

	for _, err := range errors {
		switch err.Tag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("Field '%s' is required", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("Field '%s' failed on %s validation", err.Field(), err.Tag()))
		}
	}

	return strings.Join(errMsgs, "; ")
}

func isMobileRequest(r *http.Request) bool {
	return r.Header.Get("X-Client-Type") == string(ClientTypeMobile)
}
