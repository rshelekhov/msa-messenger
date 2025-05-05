package v1

import (
	"fmt"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/rshelekhov/jwtauth"
	"github.com/rshelekhov/msa-messenger/gateway/internal/config"
	"github.com/rshelekhov/msa-messenger/gateway/internal/controller/http/v1/handler"
	mwlogger "github.com/rshelekhov/msa-messenger/gateway/internal/lib/middleware/logger"
)

type appRouter struct {
	mux    *chi.Mux
	jwtMgr jwtauth.Middleware
	h      *handler.Handler
}

func NewRouter(
	cfg *config.CORS,
	log *slog.Logger,
	jwtMgr jwtauth.Middleware,
	h *handler.Handler,
) (*chi.Mux, error) {
	ar := &appRouter{
		mux:    chi.NewRouter(),
		jwtMgr: jwtMgr,
		h:      h,
	}

	// Middleware
	ar.mux.Use(middleware.StripSlashes)
	ar.mux.Use(middleware.RequestID)
	ar.mux.Use(mwlogger.New(log))
	ar.mux.Use(middleware.Recoverer)
	ar.mux.Use(render.SetContentType(render.ContentTypeJSON))

	// CORS
	ar.mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   cfg.AllowedMethods,
		AllowedHeaders:   cfg.AllowedHeaders,
		ExposedHeaders:   cfg.ExposedHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           cfg.MaxAge,
		Debug:            cfg.Debug,
	}))

	// Swagger
	if err := registerSwagger(ar.mux); err != nil {
		return nil, fmt.Errorf("failed to register swagger: %w", err)
	}

	// Routes
	return ar.initRoutes()
}
