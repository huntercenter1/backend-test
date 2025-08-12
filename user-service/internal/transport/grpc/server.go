package grpcsvr

import (
	"context"
	"time"

	userpb "github.com/huntercenter1/backend-test/proto"
	"github.com/huntercenter1/backend-test/user-service/internal/models"
	"github.com/huntercenter1/backend-test/user-service/internal/service"
)

type Server struct {
	userpb.UnimplementedUserServiceServer
	svc service.UserService
}

func NewServer(svc service.UserService) *Server {
	return &Server{svc: svc}
}

func toProto(u *models.User) *userpb.User {
	return &userpb.User{
		Id:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *Server) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.User, error) {
	u, err := s.svc.Create(ctx, req.GetUsername(), req.GetEmail(), req.GetPassword())
	if err != nil { return nil, err }
	return toProto(u), nil
}

func (s *Server) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.User, error) {
	u, err := s.svc.Get(ctx, req.GetId())
	if err != nil { return nil, err }
	return toProto(u), nil
}

func (s *Server) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.User, error) {
	u, err := s.svc.Update(ctx, req.GetId(), req.GetUsername(), req.GetEmail(), req.GetPassword())
	if err != nil { return nil, err }
	return toProto(u), nil
}

func (s *Server) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	if err := s.svc.Delete(ctx, req.GetId()); err != nil { return nil, err }
	return &userpb.DeleteUserResponse{Ok: true}, nil
}

func (s *Server) AuthenticateUser(ctx context.Context, req *userpb.AuthRequest) (*userpb.AuthResponse, error) {
	id, err := s.svc.Authenticate(ctx, req.GetUsername(), req.GetPassword())
	if err != nil { return &userpb.AuthResponse{Ok: false, Message: "invalid credentials"}, nil }
	return &userpb.AuthResponse{Ok: true, UserId: id, Message: "ok"}, nil
}

func (s *Server) ValidateUser(ctx context.Context, req *userpb.ValidateUserRequest) (*userpb.ValidateUserResponse, error) {
	ok, err := s.svc.Validate(ctx, req.GetUserId())
	if err != nil { return nil, err }
	return &userpb.ValidateUserResponse{Valid: ok, Message: ""}, nil
}
