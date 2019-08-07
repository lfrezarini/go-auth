package resolvers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/LucasFrezarini/go-auth-manager/crypt"
	"github.com/LucasFrezarini/go-auth-manager/env"
	"github.com/LucasFrezarini/go-auth-manager/gqlerrors"
	"github.com/LucasFrezarini/go-auth-manager/gqlmodels"
	"github.com/LucasFrezarini/go-auth-manager/jsonwebtoken"
	"github.com/LucasFrezarini/go-auth-manager/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateUser(ctx context.Context, data gqlmodels.CreateUserInput) (*gqlmodels.AuthUserPayload, error) {

	if emailAlreadyExists(data.Email) {
		return nil, gqlerrors.CreateConflictError("User already exists")
	}

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
		Iss: env.Config.ServerHost,
		Sub: user.ID.Hex(),
	})

	return &gqlmodels.AuthUserPayload{
		Token: token,
		User:  &user,
	}, nil
}

func emailAlreadyExists(email string) bool {
	registered, _ := userDao.FindOne(models.User{
		Email: email,
	})

	return registered != nil
}

func (r *mutationResolver) Login(ctx context.Context, data gqlmodels.LoginUserInput) (*gqlmodels.AuthUserPayload, error) {
	user, err := userDao.FindOne(models.User{
		Email:  data.Email,
		Active: true,
	})

	if err != nil {
		return nil, gqlerrors.CreateAuthorizationError()
	}

	if !crypt.ComparePassword(user.Password, data.Password) {
		return nil, gqlerrors.CreateAuthorizationError()
	}

	token, err := jsonwebtoken.Encode(jsonwebtoken.Claims{
		Iss: env.Config.ServerHost,
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
	var user *models.User

	claims, err := jsonwebtoken.Decode(token)

	if err != nil {
		return &gqlmodels.ValidateTokenPayload{
			Claims: &claims,
			User:   user,
			Valid:  false,
		}, nil
	}

	if claims.Iss != env.Config.ServerHost {
		return &gqlmodels.ValidateTokenPayload{
			Claims: &jsonwebtoken.Claims{},
			User:   user,
			Valid:  false,
		}, nil
	}

	id, err := primitive.ObjectIDFromHex(claims.Sub)

	if err != nil {
		return nil, gqlerrors.CreateInternalServerError("Error while trying to validate the token")
	}

	user, err = userDao.FindOne(models.User{
		ID:     id,
		Active: true,
	})

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

func (r *mutationResolver) UpdateUser(ctx context.Context, data gqlmodels.UpdateUserInput) (*models.User, error) {
	userID := ctx.Value("userID")

	objectID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userID))

	if err != nil {
		log.Printf("Error while trying to convert userID to objectID: %v\n", err)
		return nil, gqlerrors.CreateInternalServerError("Error while trying to update user")
	}

	hash, err := crypt.HashPassword(data.Password)

	if err != nil {
		log.Printf("Error while trying to crypt password: %v", err)
		return nil, gqlerrors.CreateInternalServerError("Error while trying to update user")
	}

	user, err := userDao.UpdateByID(objectID, models.User{
		Password: hash,
	})

	if err != nil {
		log.Printf("Error while trying to update user: %v", err)
		return nil, gqlerrors.CreateInternalServerError("Error while trying to update user")
	}

	return user, nil
}

func (r *mutationResolver) DeactivateUser(ctx context.Context) (*models.User, error) {
	userID := ctx.Value("userID")

	objectID, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userID))

	if err != nil {
		log.Printf("Error while trying to convert userID to objectID: %v\n", err)
		return nil, gqlerrors.CreateInternalServerError("Error while trying to deactivate the user user")
	}

	user, err := userDao.UpdateByID(objectID, models.User{
		Active: false,
	})

	if err != nil {
		log.Printf("Error while trying to deactivate the user: %v", err)
		return nil, gqlerrors.CreateInternalServerError("Error while trying to deactivate the user user")
	}

	return user, nil
}
