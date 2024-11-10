package handlers

import (
	"context"
	"github.com/vlks-dev/sibavto/auth-service/internal/domain"
	"github.com/vlks-dev/sibavto/auth-service/internal/services"
	"github.com/vlks-dev/sibavto/auth-service/proto/authpb"
)

type AuthServiceServer struct {
	authpb.UnimplementedAuthServiceServer
	authService *services.AuthService
}

func NewAuthServiceServer(authService *services.AuthService) *AuthServiceServer {
	return &AuthServiceServer{authService: authService}
}

func (a *AuthServiceServer) RegisterUser(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	newUser := &domain.User{
		Name:       req.User.Name,
		Surname:    req.User.Surname,
		Patronymic: req.User.Patronymic,
		Email:      req.User.Email,
		Password:   req.User.Password,
	}

	id, err := a.authService.Register(ctx, newUser)
	if err != nil {
		return nil, err
	}
	return &authpb.RegisterResponse{
		Id:      id.String(),
		Message: "User registered successfully",
	}, nil
}

func (a *AuthServiceServer) LoginUser(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	jwtToken, err := a.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &authpb.LoginResponse{
		Token:   jwtToken,
		Message: "User logged successfully",
	}, nil
}
