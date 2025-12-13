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

type ServerRepository struct {
	collection *mongo.Collection
}

func NewServerRepository() *ServerRepository {
	return &ServerRepository{
		collection: database.DB.Collection("servers"),
	}
}

func (r *ServerRepository) Create(ctx context.Context, server *model.Server) error {
	server.Id = primitive.NewObjectID()
	server.CreatedAt = time.Now()
	server.UpdatedAt = time.Now()
	server.CurrentClients = 0

	_, err := r.collection.InsertOne(ctx, server)
	return err
}

func (r *ServerRepository) GetById(ctx context.Context, id primitive.ObjectID) (*model.Server, error) {
	var server model.Server
	filter := bson.M{"_id": id}

	err := r.collection.FindOne(ctx, filter).Decode(&server)

	if err == mongo.ErrNoDocuments {
		return nil, errors.New("server not found")
	}

	return &server, err
}

func (r *ServerRepository) ListAll(ctx context.Context) ([]*model.Server, error) {
	var servers []*model.Server

	cursor, err := r.collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &servers)
	return servers, err
}

func (r *ServerRepository) Update(ctx context.Context, server *model.Server) error {
	server.UpdatedAt = time.Now()
	filter := bson.M{"_id": server.Id}
	update := bson.M{"$set": server}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
