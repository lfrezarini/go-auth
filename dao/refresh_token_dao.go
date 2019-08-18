package dao

import (
	"context"
	"fmt"

	"github.com/LucasFrezarini/go-auth-manager/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RefreshTokenDao is a representation of a RefreshToken DAO
type RefreshTokenDao struct{}

// CreateOne pushs a new refresh token on refresh_tokens user array
func (r *RefreshTokenDao) CreateOne(userID primitive.ObjectID, token models.RefreshToken) (models.User, error) {
	updatedUser := models.User{}
	collection := db.Collection(UserCollection)

	returnDocument := options.After

	options := options.FindOneAndUpdateOptions{
		ReturnDocument: &returnDocument,
	}

	err := collection.FindOneAndUpdate(context.Background(), bson.M{"_id": userID}, bson.M{
		"$push": bson.M{
			"refresh_tokens": bson.M{
				"token":      token.Token,
				"identifier": token.Identifier,
			},
		},
	}, &options).Decode(&updatedUser)

	if err != nil {
		return models.User{}, fmt.Errorf("Error while trying to create an refresh token: %v", err)
	}

	return updatedUser, nil
}
