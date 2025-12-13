package service

import (
	"context"
	"errors"
	"os"

	"github.com/shivamp1998/vpn_backend/internal/model"
	"github.com/shivamp1998/vpn_backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServerService struct {
	serverRepo *repository.ServerRepository
}

func NewServerService() *ServerService {
	return &ServerService{
		serverRepo: repository.NewServerRepository(),
	}
}

func (s *ServerService) CreateServer(ctx context.Context, name, endpoint, region string, maxClients int32) (*model.Server, error) {

	if name == "" || endpoint == "" || region == "" {
		return nil, errors.New("name, endpoint, region is required")
	}

	if maxClients <= 0 {
		return nil, errors.New("max_clients must not be greater than 0")
	}

	server := &model.Server{
		Name:                name,
		Endpoint:            endpoint,
		PublicKey:           os.Getenv("PUBLIC_KEY"),
		PrivateKeyEncrypted: os.Getenv("PRIVATE_KEY"),
		Region:              region,
		MaxClients:          maxClients,
	}

	err := s.serverRepo.Create(ctx, server)

	if err != nil {
		return nil, err
	}

	return server, nil
}

func (s *ServerService) GetServer(ctx context.Context, serverId string) (*model.Server, error) {
	id, err := primitive.ObjectIDFromHex(serverId)

	if err != nil {
		return nil, errors.New("invalid server id")
	}

	return s.serverRepo.GetById(ctx, id)
}

func (s *ServerService) ListServers(ctx context.Context) ([]*model.Server, error) {
	return s.serverRepo.ListAll(ctx)
}

func (s *ServerService) UpdateServer(ctx context.Context, server *model.Server) error {
	err := s.serverRepo.Update(ctx, server)
	return err
}
