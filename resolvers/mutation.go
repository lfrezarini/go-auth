package resolvers

import (
	"context"
	"time"

	"github.com/LucasFrezarini/go-auth-manager/crypt"
	"github.com/LucasFrezarini/go-auth-manager/gqlmodels"
	"github.com/LucasFrezarini/go-auth-manager/jsonwebtoken"
	"github.com/LucasFrezarini/go-auth-manager/models"
	"github.com/vektah/gqlparser/gqlerror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateUser(ctx context.Context, data gqlmodels.CreateUserInput) (*gqlmodels.AuthUserPayload, error) {
	hash, err := crypt.HashPassword(data.Password)

	if err != nil {
		return nil, &gqlerror.Error{
			Message: "Error while trying to save the user",
			Extensions: map[string]interface{}{
				"code": "INTERNAL_SERVER_ERROR",
			},
		}
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
		return nil, &gqlerror.Error{
			Message: "Error while trying to save the user",
			Extensions: map[string]interface{}{
				"code": "INTERNAL_SERVER_ERROR",
			},
		}
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
		return nil, &gqlerror.Error{
			Message: "Wrong email or password",
			Extensions: map[string]interface{}{
				"code": "UNAUTHORIZED",
			},
		}
	}

	if !crypt.ComparePassword(user.Password, data.Password) {
		return nil, &gqlerror.Error{
			Message: "Wrong email or password",
			Extensions: map[string]interface{}{
				"code": "UNAUTHORIZED",
			},
		}
	}

	token, err := jsonwebtoken.Encode(jsonwebtoken.Claims{
		Iss: "http://localhost:8080",
		Sub: user.ID.Hex(),
	})

	if err != nil {
		return nil, &gqlerror.Error{
			Message: "Error while trying to login",
			Extensions: map[string]interface{}{
				"code": "INTERNAL_SERVER_ERROR",
			},
		}
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
		return nil, &gqlerror.Error{
			Message: "Error while trying to validate the token",
			Extensions: map[string]interface{}{
				"code": "INTERNAL_SERVER_ERROR",
			},
		}
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
