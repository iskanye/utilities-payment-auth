package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/iskanye/utilities-payment-auth/internal/app"
	"github.com/iskanye/utilities-payment-auth/internal/config"
	pkgConfig "github.com/iskanye/utilities-payment-utils/pkg/config"
	"github.com/iskanye/utilities-payment-utils/pkg/logger"
)

func main() {
	cfg := pkgConfig.MustLoad[config.Config]()
	cfg.LoadSecret()

	log := logger.SetupPrettySlog()
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
