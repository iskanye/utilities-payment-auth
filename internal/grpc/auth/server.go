package auth

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	protoAuth "github.com/iskanye/utilities-payment-proto/auth"
	"github.com/iskanye/utilities-payment/pkg/models"
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
	) (string, error)
	Register(
		ctx context.Context,
		email string,
		password string,
	) (int64, error)
	GetUsers(
		ctx context.Context,
	) ([]models.User, error)
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

	return &protoAuth.LoginResponse{
		Token: token,
	}, nil
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

func (s *serverAPI) Users(
	in *protoAuth.UsersRequest,
	stream grpc.ServerStreamingServer[protoAuth.User],
) error {
	users, err := s.auth.GetUsers(stream.Context())
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for _, i := range users {
		id := i.ID
		email := i.Email
		user := &protoAuth.User{
			Id:    id,
			Email: email,
		}
		if err = stream.SendMsg(user); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}

	return nil
}
