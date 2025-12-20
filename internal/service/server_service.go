package service

import (
	"context"
	"errors"

	"github.com/shivamp1998/vpn_backend/internal/model"
	"github.com/shivamp1998/vpn_backend/internal/repository"
	"github.com/shivamp1998/vpn_backend/internal/wireguard"
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

func (s *ServerService) CreateServer(ctx context.Context, name, endpoint, region, publicKey string, maxClients int32) (*model.Server, error) {

	if name == "" || endpoint == "" || region == "" {
		return nil, errors.New("name, endpoint, region is required")
	}

	if maxClients <= 0 {
		return nil, errors.New("max_clients must not be greater than 0")
	}

	server := &model.Server{
		Name:       name,
		Endpoint:   endpoint,
		Region:     region,
		MaxClients: maxClients,
	}

	if publicKey == "" {
		server.PublicKey = publicKey
		server.PrivateKeyEncrypted = ""
	} else {
		privateKey, publicKey, err := wireguard.GenerateKeyPair()
		if err != nil {
			return nil, err
		}
		server.PublicKey = publicKey
		server.PrivateKeyEncrypted = privateKey
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
