package main

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	protoAuth "github.com/iskanye/utilities-payment-proto/auth"
)

type serverAPI struct {
	protoAuth.UnimplementedAuthServer
	auth Auth
}

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	Register(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	protoAuth.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(
	ctx context.Context,
	in *protoAuth.LoginRequest,
) (*protoAuth.LoginResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if in.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}

	token, err := s.auth.Login(ctx, in.GetEmail(), in.GetPassword(), int(in.GetAppId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &protoAuth.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	in *protoAuth.RegisterRequest,
) (*protoAuth.RegisterResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	uid, err := s.auth.Register(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &protoAuth.RegisterResponse{UserId: uid}, nil
}
