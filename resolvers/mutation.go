package resolvers

import (
	"context"
	"time"

	"github.com/LucasFrezarini/go-auth-manager/crypt"
	"github.com/LucasFrezarini/go-auth-manager/gqlerrors"
	"github.com/LucasFrezarini/go-auth-manager/gqlmodels"
	"github.com/LucasFrezarini/go-auth-manager/jsonwebtoken"
	"github.com/LucasFrezarini/go-auth-manager/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateUser(ctx context.Context, data gqlmodels.CreateUserInput) (*gqlmodels.AuthUserPayload, error) {
	hash, err := crypt.HashPassword(data.Password)

	if err != nil {
		return nil, gqlerrors.CreateInternalServerError("Error while trying to create user")
	}

	user := models.User{
		Email:     data.Email,
		Password:  hash,
		Roles:     []string{"user"},
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	insertedID, err := userDao.CreateOne(user)

	if err != nil {
		return nil, gqlerrors.CreateInternalServerError("Error while trying to create user")
	}

	user.ID = insertedID
	token, err := jsonwebtoken.Encode(jsonwebtoken.Claims{
		Iss: "http://localhost:8080",
		Sub: user.ID.Hex(),
	})

	return &gqlmodels.AuthUserPayload{
		Token: token,
		User:  &user,
	}, nil
}

func (r *mutationResolver) Login(ctx context.Context, data gqlmodels.LoginUserInput) (*gqlmodels.AuthUserPayload, error) {
	user, err := userDao.FindOne(models.User{
		Email: data.Email,
	})

	if err != nil {
		return nil, gqlerrors.CreateAuthorizationError()
	}

	if !crypt.ComparePassword(user.Password, data.Password) {
		return nil, gqlerrors.CreateAuthorizationError()
	}

	token, err := jsonwebtoken.Encode(jsonwebtoken.Claims{
		Iss: "http://localhost:8080",
		Sub: user.ID.Hex(),
	})

	if err != nil {
		return nil, gqlerrors.CreateInternalServerError("Error while trying to login")
	}

	return &gqlmodels.AuthUserPayload{
		User:  user,
		Token: token,
	}, nil
}

func (r *mutationResolver) ValidateToken(ctx context.Context, token string) (*gqlmodels.ValidateTokenPayload, error) {
	claims, err := jsonwebtoken.Decode(token)
	var user *models.User

	if err != nil {
		return &gqlmodels.ValidateTokenPayload{
			Claims: &claims,
			User:   user,
			Valid:  false,
		}, nil
	}

	id, err := primitive.ObjectIDFromHex(claims.Sub)

	if err != nil {
		return nil, gqlerrors.CreateInternalServerError("Error while trying to validate the token")
	}

	user, err = userDao.FindByID(id)

	if err != nil {
		return &gqlmodels.ValidateTokenPayload{
			Claims: &jsonwebtoken.Claims{},
			User:   user,
			Valid:  false,
		}, nil
	}

	return &gqlmodels.ValidateTokenPayload{
		Claims: &claims,
		User:   user,
		Valid:  true,
	}, nil
}
