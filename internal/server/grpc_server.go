package server

import (
	"context"
	"errors"
	"log"

	"github.com/shivamp1998/vpn_backend/internal/auth"
	"github.com/shivamp1998/vpn_backend/internal/service"
	pb "github.com/shivamp1998/vpn_backend/proto/gen"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	pb.UnimplementedServerServiceServer
	userService   *service.UserService
	serverService *service.ServerService
}

func NewServer() *Server {
	return &Server{
		userService:   service.NewUserService(),
		serverService: service.NewServerService(),
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

	token, err := auth.GenerateToken(user.Id, user.Email)

	if err != nil {
		return &pb.AuthenticationResponse{
			Status: "error",
			Token:  "",
		}, err
	}

	return &pb.AuthenticationResponse{
		Status: "success",
		Token:  token,
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

	token, err := auth.GenerateToken(user.Id, user.Email)

	if err != nil {
		return &pb.AuthenticationResponse{
			Status: "Error",
			Token:  "",
		}, err
	}

	return &pb.AuthenticationResponse{
		Status: "Success",
		Token:  token,
	}, nil
}

func (s *Server) CreateServer(ctx context.Context, req *pb.CreateServerRequest) (*pb.CreateServerResponse, error) {
	server, err := s.serverService.CreateServer(ctx, req.Name, req.Endpoint, req.Region, req.MaxClients)

	if err != nil {
		return nil, errors.New("error creating server")
	}

	return &pb.CreateServerResponse{
		Server: &pb.Server{
			Id:             server.Id.Hex(),
			Name:           server.Name,
			Endpoint:       server.Endpoint,
			PublicKey:      server.PublicKey,
			Region:         server.Region,
			MaxClients:     server.MaxClients,
			CurrentClients: server.CurrentClients,
		},
		Message: "server created successfully!",
	}, nil
}

func (s *Server) ListServers(ctx context.Context, req *pb.ListServerRequest) (*pb.ListServerResponse, error) {
	log.Println("ListServer request")

	servers, err := s.serverService.ListServers(ctx)

	if err != nil {
		return nil, err
	}

	pbServers := make([]*pb.Server, len(servers))

	for i, server := range servers {
		pbServers[i] = &pb.Server{
			Id:             server.Id.Hex(),
			Name:           server.Name,
			Endpoint:       server.Endpoint,
			PublicKey:      server.PublicKey,
			Region:         server.Region,
			MaxClients:     server.MaxClients,
			CurrentClients: server.CurrentClients,
		}
	}

	return &pb.ListServerResponse{
		Servers: pbServers,
	}, nil
}

func (s *Server) GetServer(ctx context.Context, req *pb.GetServerRequest) (*pb.GetServerResponse, error) {
	server, err := s.serverService.GetServer(ctx, req.ServerId)

	if err != nil {
		return nil, err
	}

	return &pb.GetServerResponse{
		Server: &pb.Server{
			Id:             server.Id.Hex(),
			Name:           server.Name,
			Endpoint:       server.Endpoint,
			Region:         server.Region,
			PublicKey:      server.PublicKey,
			MaxClients:     server.MaxClients,
			CurrentClients: server.CurrentClients,
		},
	}, nil
}
