package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/iskanye/utilities-payment-api-gateway/pkg/logger"
	"github.com/iskanye/utilities-payment-auth/internal/app"
	"github.com/iskanye/utilities-payment-auth/internal/config"
)

func main() {
	cfg := config.MustLoad()
	log := setupPrettySlog()
	app := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.Secret, cfg.TokenTTL)

	go func() {
		app.GRPCServer.MustRun()
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.GRPCServer.Stop()
	log.Info("Gracefully stopped")
}

func setupPrettySlog() *slog.Logger {
	opts := logger.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
