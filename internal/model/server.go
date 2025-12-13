package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Server struct {
	Id                  primitive.ObjectID `bson:"_id,omit_empty" json:"id"`
	Name                string             `bson:"name" json:"name"`
	Endpoint            string             `bson:"endpoint" json:"endpoint"`
	PublicKey           string             `bson:"public_key" json:"public_key"`
	PrivateKeyEncrypted string             `bson:"private_key_encrypted" json:"private_key_encrypted"`
	Region              string             `bson:"region" json:"region"`
	MaxClients          int32              `bson:"max_clients" json:"max_clients"`
	CurrentClients      int32              `bson:"current_clients" json:"current_clients"`
	CreatedAt           time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt           time.Time          `bson:"updated_at" json:"updated_at"`
}
