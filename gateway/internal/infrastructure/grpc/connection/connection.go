package connection

import (
	"errors"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	GRPCConnections struct {
		SSOConn        *grpc.ClientConn
		ChatConn       *grpc.ClientConn
		SubscriberConn *grpc.ClientConn
	}

	Addresses struct {
		SSO        string
		Chat       string
		Subscriber string
	}
)

func NewGRPCConnections(addresses Addresses, clientID string, authInterceptor grpc.UnaryClientInterceptor) (*GRPCConnections, error) {
	const op = "grpc.connection.NewGRPCConnections"

	ssoConn, err := grpc.NewClient(
		addresses.SSO,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(authInterceptor),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create sso connection: %w", op, err)
	}

	chatConn, err := grpc.NewClient(
		addresses.Chat,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(authInterceptor),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create chat connection: %w", op, err)
	}

	subscriberConn, err := grpc.NewClient(
		addresses.Subscriber,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(authInterceptor),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create subscriber connection: %w", op, err)
	}

	return &GRPCConnections{
		SSOConn:        ssoConn,
		ChatConn:       chatConn,
		SubscriberConn: subscriberConn,
	}, nil
}

func (c *GRPCConnections) Close() error {
	var errs []error

	if c.SSOConn != nil {
		if err := c.SSOConn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("sso connection: %w", err))
		}
	}

	if c.ChatConn != nil {
		if err := c.ChatConn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("chat connection: %w", err))
		}
	}

	if c.SubscriberConn != nil {
		if err := c.SubscriberConn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("subscriber connection: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("grpc.connection.GRPCConnections.Close: failed to close connections: %w", errors.Join(errs...))
	}

	return nil
}
