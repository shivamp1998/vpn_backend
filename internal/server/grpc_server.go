package server

import (
	"context"
	"log"

	"github.com/shivamp1998/vpn_backend/internal/service"
	pb "github.com/shivamp1998/vpn_backend/proto/gen"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	userService *service.UserService
}

func NewServer() *Server {
	return &Server{
		userService: service.NewUserService(),
	}
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthenticationResponse, error) {
	log.Printf("Register request for email: %s", req.Email)

	user, err := s.userService.Register(ctx, req.Email, req.Password)

	if err != nil {
		return &pb.AuthenticationResponse{
			Status: "error",
			Token:  "",
		}, err
	}

	return &pb.AuthenticationResponse{
		Status: "success",
		Token:  user.Id.Hex(),
	}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthenticationResponse, error) {
	user, err := s.userService.Login(ctx, req.Email, req.Password)

	if err != nil {
		return &pb.AuthenticationResponse{
			Status: "Error",
			Token:  "",
		}, err
	}

	return &pb.AuthenticationResponse{
		Status: "Success",
		Token:  user.Id.Hex(),
	}, nil
}
