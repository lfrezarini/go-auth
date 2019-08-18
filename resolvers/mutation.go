package resolvers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/LucasFrezarini/go-auth-manager/credentials"
	"github.com/LucasFrezarini/go-auth-manager/crypt"
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
		log.Printf("Error while trying to create user: %v\n", err)
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
		log.Printf("Error while trying to create user: %v\n", err)
		return nil, gqlerrors.CreateInternalServerError("Error while trying to create user")
	}

	user.ID = insertedID
	token, err := jsonwebtoken.Encode(jsonwebtoken.CreateDefaultClaims(user.ID.Hex()))

	if err != nil {
		log.Printf("Error while trying to create user: %v\n", err)
		return nil, gqlerrors.CreateInternalServerError("Error while trying to create user")
	}

	refreshToken, err := jsonwebtoken.Encode(jsonwebtoken.CreateRefreshTokenClaims(user.ID.Hex()))

	if err != nil {
		log.Printf("Error while trying to create user: %v\n", err)
		return nil, gqlerrors.CreateInternalServerError("Error while trying to create user")
	}

	updatedUser, err := refreshTokenDao.CreateOne(user.ID, models.RefreshToken{
		Token:      refreshToken,
		Identifier: "unknown",
	})

	if err != nil {
		log.Printf("Error while trying to create user: %v\n", err)
		return nil, gqlerrors.CreateInternalServerError("Error while trying to create user")
	}

	return &gqlmodels.AuthUserPayload{
		RefreshToken: refreshToken,
		Token:        token,
		User:         &updatedUser,
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
		log.Printf("Error while trying to login: %v\n", err)
		return nil, gqlerrors.CreateAuthorizationError()
	}

	if !crypt.ComparePassword(user.Password, data.Password) {
		log.Println("Error while trying to login: Invalid Password")
		return nil, gqlerrors.CreateAuthorizationError()
	}

	token, err := jsonwebtoken.Encode(jsonwebtoken.CreateDefaultClaims(user.ID.Hex()))

	if err != nil {
		log.Printf("Error while trying to login: %v\n", err)
		return nil, gqlerrors.CreateInternalServerError("Error while trying to login")
	}

	refreshToken, err := jsonwebtoken.Encode(jsonwebtoken.CreateRefreshTokenClaims(user.ID.Hex()))

	if err != nil {
		log.Printf("Error while trying to login: %v\n", err)
		return nil, gqlerrors.CreateInternalServerError("Error while trying to login")
	}

	return &gqlmodels.AuthUserPayload{
		User:         user,
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (r *mutationResolver) ValidateToken(ctx context.Context, token string) (*gqlmodels.ValidateTokenPayload, error) {
	claims, err := jsonwebtoken.Decode(token)

	if err != nil {
		return &gqlmodels.ValidateTokenPayload{
			Claims: nil,
			User:   nil,
			Valid:  false,
		}, nil
	}

	user, err := credentials.ValidateCredentials(claims)

	if err != nil {
		return &gqlmodels.ValidateTokenPayload{
			Claims: nil,
			User:   nil,
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

func (r *mutationResolver) RefreshToken(ctx context.Context, refreshToken string) (*gqlmodels.AuthUserPayload, error) {
	claims, err := jsonwebtoken.Decode(refreshToken)

	if err != nil {
		return nil, gqlerrors.CreateAuthorizationError()
	}

	objectID, err := primitive.ObjectIDFromHex(claims.Subject)

	if err != nil {
		log.Printf("Error while trying to convert userID to objectID: %v\n", err)
		return nil, gqlerrors.CreateInternalServerError("Error while trying to refresh token the user")
	}

	user, err := userDao.FindByID(objectID)

	if err != nil {
		log.Printf("Error while trying to get a refreshed token: %v", err)
		return nil, gqlerrors.CreateInternalServerError("Error while trying to refresh token")
	}

	refreshTokenExists := false

	for _, value := range user.RefreshTokens {
		refreshTokenExists = value.Token == refreshToken

		if refreshTokenExists {
			break
		}
	}

	if !refreshTokenExists {
		return nil, gqlerrors.CreateAuthorizationError()
	}

	token, err := jsonwebtoken.Encode(jsonwebtoken.CreateRefreshTokenClaims(user.ID.Hex()))

	if err != nil {
		return nil, gqlerrors.CreateInternalServerError("Error while trying to refresh token")
	}

	return &gqlmodels.AuthUserPayload{
		User:         user,
		RefreshToken: refreshToken,
		Token:        token,
	}, nil
}
