package app

import (
	"log/slog"
	"time"

	"github.com/iskanye/utilities-payment-auth/internal/app/grpc"
	"github.com/iskanye/utilities-payment-auth/internal/lib/jwt"
	"github.com/iskanye/utilities-payment-auth/internal/service/auth"
	"github.com/iskanye/utilities-payment-auth/internal/storage"
)

type App struct {
	GRPCServer *grpc.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	secret string,
	tokenTTL time.Duration,
) *App {
	storage, err := storage.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, jwt.Validate, secret, tokenTTL)
	grpcApp := grpc.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
