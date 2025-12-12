package service

import (
	"context"
	"errors"

	"github.com/shivamp1998/vpn_backend/internal/auth"
	"github.com/shivamp1998/vpn_backend/internal/model"
	"github.com/shivamp1998/vpn_backend/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRespository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

func (s *UserService) Register(ctx context.Context, email, password string) (*model.User, error) {
	existing, _ := s.userRepo.GetByEmail(ctx, email)

	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := auth.HashPassword(password)

	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:        email,
		PasswordHash: hashedPassword,
	}

	err = s.userRepo.Create(ctx, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*model.User, error) {
	user, _ := s.userRepo.GetByEmail(ctx, email)

	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if !auth.CheckPasswordHash(user.PasswordHash, password) {
		return nil, errors.New("Passwords do not match")
	}

	return user, nil
}
