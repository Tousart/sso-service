package auth

import (
	"context"
	"fmt"
	"log"

	ssov1 "github.com/tousart/protos/gen/go/sso"
	"github.com/tousart/sso/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth usecase.AuthService
}

func CreateServerAPI(auth usecase.AuthService) *serverAPI {
	return &serverAPI{auth: auth}
}

func Register(gRPCServer *grpc.Server, api *serverAPI) {
	ssov1.RegisterAuthServer(gRPCServer, api)
}

func (s *serverAPI) Login(ctx context.Context, in *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	if in.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	token, err := s.auth.Login(ctx, in.GetLogin(), in.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to login")
	}

	// DELETE
	log.Println("successful login")
	log.Printf("token: %s", token)

	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, in *ssov1.RegisterRequest) (*emptypb.Empty, error) {
	if in.Login == "" {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "login is required")
	}

	if in.Password == "" {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "password is required")
	}

	if in.Email == "" {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "email is required")
	}

	err := s.auth.Register(ctx, in.GetLogin(), in.GetPassword(), in.GetEmail())
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, fmt.Sprintf("failed to register user: %v", err))
	}

	// DELETE
	log.Println("successful registration")

	return &emptypb.Empty{}, nil
}
