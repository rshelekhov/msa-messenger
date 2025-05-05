package v1

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func registerSwagger(router *chi.Mux) error {
	swagger, err := GetSwagger()
	if err != nil {
		return fmt.Errorf("failed to load swagger spec: %w", err)
	}

	// Clear servers from spec, to use current host
	swagger.Servers = nil

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	return nil
}
