package app

import (
	"log/slog"
	"time"

	"github.com/iskanye/utilities-payment-auth/internal/auth/app/grpc"
	"github.com/iskanye/utilities-payment-auth/internal/auth/service/auth"
	"github.com/iskanye/utilities-payment-auth/internal/auth/storage"
)

type App struct {
	GRPCServer *grpc.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := storage.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)
	grpcApp := grpc.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
