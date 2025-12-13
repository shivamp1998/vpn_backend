package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/shivamp1998/vpn_backend/internal/model"
	"github.com/shivamp1998/vpn_backend/internal/repository"
	"github.com/shivamp1998/vpn_backend/internal/wireguard"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConfigService struct {
	userRepo      *repository.UserRespository
	serverRepo    *repository.ServerRepository
	keysRepo      *repository.WireGuardKeysRepository
	serverNetwork string
}

func NewConfigService() *ConfigService {
	return &ConfigService{
		userRepo:      repository.NewUserRepository(),
		serverRepo:    repository.NewServerRepository(),
		keysRepo:      repository.NewWireGuardKeysRepository(),
		serverNetwork: "10.0.0.0/24",
	}
}

type ConfigResult struct {
	ConfigContent string
	QRCodeBase64  string
	ConfigData    ConfigData
}

type ConfigData struct {
	PrivateKey      string
	PublicKey       string
	ServerPublicKey string
	ServerEndpoint  string
	ClientIp        string
	DNS             string
}

func (s *ConfigService) GenerateConfig(ctx context.Context, userId primitive.ObjectID, serverId string) (*ConfigResult, error) {
	serverObjId, err := primitive.ObjectIDFromHex(serverId)
	if err != nil {
		return nil, errors.New("invalid server ID")
	}

	server, err := s.serverRepo.GetById(ctx, serverObjId)
	if err != nil {
		return nil, errors.New("server not found")
	}
	existingKeys, err := s.keysRepo.GetByUserAndServer(ctx, userId, serverObjId)
	fmt.Print(err)
	var keys *model.WireGuardKeys

	if err != nil {
		privateKey, publicKey, err := wireguard.GenerateKeyPair()
		if err != nil {
			return nil, fmt.Errorf("failed to generate keys: %v", err)
		}

		clientIp, err := s.assignClientIp(ctx, serverObjId)
		if err != nil {
			return nil, fmt.Errorf("failed to generate keys: %v", err)
		}

		keys = &model.WireGuardKeys{
			UserId:              userId,
			ServerId:            serverObjId,
			PrivateKeyEncrypted: privateKey,
			PublicKey:           publicKey,
			IpAddress:           clientIp,
		}

		err = s.keysRepo.Create(ctx, keys)

		if err != nil {
			return nil, fmt.Errorf("failed to save keys: %v", err)
		}

	} else {
		keys = existingKeys
	}

	configContent := wireguard.GenerateClientConfig(keys.PrivateKeyEncrypted, server.PublicKey, server.Endpoint, keys.IpAddress, "8.8.8.8")

	qrCode, err := wireguard.GeneateQRCode(configContent)
	if err != nil {
		qrCode = ""
	}

	result := &ConfigResult{
		ConfigContent: configContent,
		QRCodeBase64:  qrCode,
		ConfigData: ConfigData{
			PrivateKey:      keys.PrivateKeyEncrypted,
			PublicKey:       keys.PublicKey,
			ServerPublicKey: server.PublicKey,
			ServerEndpoint:  server.Endpoint,
			ClientIp:        keys.IpAddress,
			DNS:             "8.8.8.8",
		},
	}

	return result, nil
}

func (s *ConfigService) assignClientIp(ctx context.Context, serverId primitive.ObjectID) (string, error) {
	baseIp := 2

	ip := fmt.Sprintf("10.0.0.%d/32", baseIp)

	return ip, nil
}
