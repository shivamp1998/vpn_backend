package server

import (
	"context"
	"errors"
	"log"

	"github.com/shivamp1998/vpn_backend/internal/auth"
	"github.com/shivamp1998/vpn_backend/internal/service"
	pb "github.com/shivamp1998/vpn_backend/proto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	pb.UnimplementedServerServiceServer
	pb.UnimplementedConfigServiceServer
	userService   *service.UserService
	serverService *service.ServerService
	configService *service.ConfigService
}

func NewServer() *Server {
	return &Server{
		userService:   service.NewUserService(),
		serverService: service.NewServerService(),
		configService: service.NewConfigService(),
	}
}

func (s *Server) GenerateConfig(ctx context.Context, req *pb.GenerateConfigRequest) (*pb.GenerateConfigResponse, error) {
	userId, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	result, err := s.configService.GenerateConfig(ctx, userId, req.ServerId)

	if err != nil {
		return &pb.GenerateConfigResponse{
			Message: "error",
		}, err
	}
	return &pb.GenerateConfigResponse{
		ConfigContent: result.ConfigContent,
		QrCodeBase64:  result.QRCodeBase64,
		ConfigData: &pb.ConfigData{
			PrivateKey:      result.ConfigData.PrivateKey,
			PublicKey:       result.ConfigData.PublicKey,
			ServerPublicKey: result.ConfigData.ServerPublicKey,
			ServerEndpoint:  result.ConfigData.ServerEndpoint,
			ServerAddress:   result.ConfigData.ServerAddress,
			ServerPort:      result.ConfigData.ServerPort,
			ClientIp:        result.ConfigData.ClientIp,
			Dns:             result.ConfigData.DNS,
		},
		Message: "success",
	}, nil
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
	server, err := s.serverService.CreateServer(ctx, req.Name, req.Endpoint, req.Region, req.PublicKey, req.MaxClients)

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
