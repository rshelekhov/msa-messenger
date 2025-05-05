package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/rshelekhov/msa-messenger/gateway/internal/app"
	"github.com/rshelekhov/msa-messenger/gateway/internal/config"
	"github.com/rshelekhov/msa-messenger/gateway/internal/lib/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.AppEnv)

	log = log.With(slog.String("env", cfg.AppEnv))

	log.Info("starting application")
	log.Debug("logger debug mode enabled")

	application, err := app.New(cfg, log)
	if err != nil {
		log.Error("failed to initialize application", slog.String("error", err.Error()))
		os.Exit(1)
	}

	go func() {
		application.HTTPServer.MustRun()
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Info("shutting down...", slog.String("signal", sign.String()))

	if err := application.Stop(); err != nil {
		log.Error("failed to stop application", slog.String("error", err.Error()))
		os.Exit(1)
	}

	log.Info("graceful shutdown completed")
}
