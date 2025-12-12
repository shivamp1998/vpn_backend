package repository

import (
	"context"
	"errors"
	"time"

	"github.com/shivamp1998/vpn_backend/internal/database"
	"github.com/shivamp1998/vpn_backend/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRespository struct {
	collection *mongo.Collection
}

func NewUserRepository() *UserRespository {
	return &UserRespository{
		collection: database.DB.Collection("users"),
	}
}

func (r *UserRespository) Create(ctx context.Context, user *model.User) error {
	user.Id = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsActive = true

	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *UserRespository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	query := bson.M{
		"email": email,
	}

	err := r.collection.FindOne(ctx, query).Decode(&user)

	if err == mongo.ErrNoDocuments {
		return nil, errors.New("user not found")
	}

	return &user, err
}
