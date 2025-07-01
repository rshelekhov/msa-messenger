package app

import (
	"fmt"
	"log/slog"

	"github.com/rshelekhov/jwtauth"
	"github.com/rshelekhov/msa-messenger/gateway/internal/app/http"
	"github.com/rshelekhov/msa-messenger/gateway/internal/config"
	v1 "github.com/rshelekhov/msa-messenger/gateway/internal/controller/http/v1"
	"github.com/rshelekhov/msa-messenger/gateway/internal/controller/http/v1/handler"
	"github.com/rshelekhov/msa-messenger/gateway/internal/domain/usecase/auth"
	"github.com/rshelekhov/msa-messenger/gateway/internal/domain/usecase/chat"
	"github.com/rshelekhov/msa-messenger/gateway/internal/domain/usecase/subscriber"
	"github.com/rshelekhov/msa-messenger/gateway/internal/domain/usecase/user"
	"github.com/rshelekhov/msa-messenger/gateway/internal/lib/validator"
)

type App struct {
	HTTPServer *http.App
}

func New(cfg *config.Config, log *slog.Logger) (*App, error) {
	validate, err := validator.New(&cfg.Validator)
	if err != nil {
		return nil, fmt.Errorf("failed to init validator: %w", err)
	}

	jwksProvider := jwtauth.NewRemoteJWKSProvider(cfg.JWT.JWKSEndpoint)
	jwtManager, err := jwtauth.NewManager(jwksProvider, jwtauth.WithAppID(cfg.App.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to init jwt manager: %w", err)
	}

	authUsecase := auth.NewUsecase(log)
	userUsecase := user.NewUsecase(log)
	subscriberUsecase := subscriber.NewUsecase(log)
	chatUsecase := chat.NewUsecase(log)

	handler := handler.New(log, validate, jwtManager, authUsecase, userUsecase, subscriberUsecase, chatUsecase)

	router, err := v1.NewRouter(&cfg.CORS, log, jwtManager, handler)
	if err != nil {
		return nil, fmt.Errorf("failed to init router: %w", err)
	}

	httpServer := http.New(&cfg.HTTPServer, log, router)

	return &App{
		HTTPServer: httpServer,
	}, nil
}

func (a *App) Stop() error {
	const method = "app.Stop"

	// Shutdown HTTP server
	if err := a.HTTPServer.Stop(); err != nil {
		return fmt.Errorf("%s: failed to stop http server: %w", method, err)
	}

	return nil
}
