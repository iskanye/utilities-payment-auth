package auth

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
	) (token string, err error)
	Register(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
	Validate(
		ctx context.Context,
		token string,
	) (isValid bool, err error)
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

	token, err := s.auth.Login(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
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
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &protoAuth.RegisterResponse{UserId: uid}, nil
}

func (s *serverAPI) Validate(
	ctx context.Context,
	in *protoAuth.ValidateRequest,
) (*protoAuth.ValidateResponse, error) {
	if in.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	isValid, err := s.auth.Validate(ctx, in.GetToken())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &protoAuth.ValidateResponse{IsValid: isValid}, nil
}
