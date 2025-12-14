package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/shivamp1998/vpn_backend/internal/database"
	"github.com/shivamp1998/vpn_backend/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WireGuardKeysRepository struct {
	collection *mongo.Collection
}

func NewWireGuardKeysRepository() *WireGuardKeysRepository {
	return &WireGuardKeysRepository{
		collection: database.DB.Collection("wireguard_keys"),
	}
}

func (r *WireGuardKeysRepository) Create(ctx context.Context, keys *model.WireGuardKeys) error {
	keys.Id = primitive.NewObjectID()
	keys.CreatedAt = time.Now()
	keys.LastRotatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, keys)
	return err
}

func (r *WireGuardKeysRepository) GetByUserAndServer(ctx context.Context, userId, serverId primitive.ObjectID) (*model.WireGuardKeys, error) {
	var keys model.WireGuardKeys

	filter := bson.M{
		"user_id":   userId,
		"server_id": serverId,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&keys)
	fmt.Print(keys)

	if err == mongo.ErrNoDocuments {
		return nil, errors.New("keys not found")
	}

	return &keys, err
}

func (r *WireGuardKeysRepository) Update(ctx context.Context, keys *model.WireGuardKeys) error {
	keys.LastRotatedAt = time.Now()
	filter := bson.M{"_id": keys.Id}
	update := bson.M{"$set": keys}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *WireGuardKeysRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

func (r *WireGuardKeysRepository) GetAllByServer(ctx context.Context, serverId primitive.ObjectID) ([]*model.WireGuardKeys, error) {
	var keys []*model.WireGuardKeys

	filter := bson.M{
		"server_id": serverId,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &keys)
	return keys, err
}
