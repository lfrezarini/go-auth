package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents the data structure of a user in the MongoDB database
type User struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name          string             `json:"string" bson:"name,omitempty"`
	Email         string             `json:"email" bson:"email,omitempty"`
	Password      string             `json:"password" bson:"password,omitempty"`
	Roles         []string           `json:"roles" bson:"roles,omitempty"`
	Active        bool               `json:"active" bson:"active,omitempty"`
	RefreshTokens []RefreshToken     `json:"refresh_tokens" bson:"refresh_tokens,omitempty"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt     time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
}
