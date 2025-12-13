package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WireGuardKeys struct {
	Id                  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId              primitive.ObjectID `bson:"user_id" json:"user_id"`
	ServerId            primitive.ObjectID `bson:"server_id" json:"server_id"`
	PrivateKeyEncrypted string             `bson:"private_key_encrypted" json:"private_key_encrypted"`
	PublicKey           string             `bson:"public_key" json:"public_key"`
	IpAddress           string             `bson:"ip_address" json:"ip_address"`
	CreatedAt           time.Time          `bson:"created_at" json:"created_at"`
	LastRotatedAt       time.Time          `bson:"last_rotated_at" json:"last_rotated_at"`
}
