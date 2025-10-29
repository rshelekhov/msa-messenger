package app

import (
	"fmt"
	"log/slog"

	"github.com/rshelekhov/msa-messenger/gateway/internal/app/http"
	"github.com/rshelekhov/msa-messenger/gateway/internal/config"
	v1 "github.com/rshelekhov/msa-messenger/gateway/internal/controller/http/v1"
	"github.com/rshelekhov/msa-messenger/gateway/internal/controller/http/v1/handler"
	"github.com/rshelekhov/msa-messenger/gateway/internal/domain/usecase"
	"github.com/rshelekhov/msa-messenger/gateway/internal/infrastructure/grpc/client"
	"github.com/rshelekhov/msa-messenger/gateway/internal/infrastructure/grpc/connection"
	"github.com/rshelekhov/msa-messenger/gateway/internal/lib/validator"
	"github.com/rshelekhov/sso/pkg/jwtauth"
)

type App struct {
	HTTPServer *http.App
}

func New(cfg *config.Config, log *slog.Logger) (*App, error) {
	validate, err := validator.New(&cfg.Validator)
	if err != nil {
		return nil, fmt.Errorf("failed to init validator: %w", err)
	}

	jwksProvider, err := jwtauth.NewRemoteJWKSProvider(cfg.JWT.JWKSEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to init jwks provider: %w", err)
	}

	jwtManager := jwtauth.NewManager(jwksProvider)

	grpcConnections, err := connection.NewGRPCConnections(
		connection.Addresses{
			SSO:        cfg.GRPCServices.SSOService.Address,
			Chat:       cfg.GRPCServices.ChatService.Address,
			Subscriber: cfg.GRPCServices.SubscriberService.Address,
		},
		cfg.App.ID,
		jwtManager.AuthUnaryClientInterceptor(cfg.App.ID),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init grpc connections: %w", err)
	}

	authClient := client.NewAuthClient(log, grpcConnections.SSOConn, cfg.GRPCServices.SSOService)
	userClient := client.NewUserClient(log, grpcConnections.SSOConn)
	subscriberClient := client.NewSubscriberClient(log, grpcConnections.SubscriberConn)
	chatClient := client.NewChatClient(log, grpcConnections.ChatConn)

	authUsecase := usecase.NewAuthUsecase(authClient)
	userUsecase := usecase.NewUserUsecase(userClient)
	subscriberUsecase := usecase.NewSubscriberUsecase(subscriberClient)
	chatUsecase := usecase.NewChatUsecase(chatClient)

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
