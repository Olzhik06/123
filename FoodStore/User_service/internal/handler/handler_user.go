package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"user/internal/auth"
	"user/internal/model"
	"user/internal/service"
	pb "user/proto"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     "user",
	}
	id, err := h.svc.Register(ctx, user)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterResponse{Id: id}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := h.svc.GetUserByEmail(ctx, req.Email)

	token, err := auth.GenerateToken(user.ID, user.Role) // передаем роль в токен
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &pb.LoginResponse{
		Message: "Login successful",
		Token:   token,
		UserId:  user.ID,
		Role:    user.Role,
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := h.svc.GetUserByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserResponse{
		User: &pb.User{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Password: "",
			Role:     user.Role,
		},
	}, nil
}
