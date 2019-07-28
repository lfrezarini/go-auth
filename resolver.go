//go:generate go run github.com/99designs/gqlgen

package auth_manager

import (
	"context"
	"errors"
	"time"

	"github.com/LucasFrezarini/go-auth-manager/crypt"
	"github.com/LucasFrezarini/go-auth-manager/dao"
	"github.com/LucasFrezarini/go-auth-manager/jsonwebtoken"
	"github.com/LucasFrezarini/go-auth-manager/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Resolver struct {
	users []*models.User
}

var userDao dao.UserDao

func init() {
	userDao = dao.UserDao{}
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateUser(ctx context.Context, data CreateUserInput) (*AuthUserPayload, error) {
	hash, err := crypt.HashPassword(data.Password)

	if err != nil {
		return nil, errors.New("Error while trying to save the user")
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
		return nil, errors.New("Error while trying to save the user")
	}

	user.ID = insertedID
	token, err := jsonwebtoken.Encode(jsonwebtoken.Claims{
		Iss: "http://localhost:8080",
		Sub: user.ID.Hex(),
	})

	return &AuthUserPayload{
		Token: token,
		User:  &user,
	}, nil
}

func (r *mutationResolver) Login(ctx context.Context, data LoginUserInput) (*AuthUserPayload, error) {
	user, err := userDao.FindOne(models.User{
		Email: data.Email,
	})

	if err != nil {
		return nil, errors.New("Error while trying to login")
	}

	if !crypt.ComparePassword(user.Password, data.Password) {
		return nil, errors.New("Wrong email or password")
	}

	token, err := jsonwebtoken.Encode(jsonwebtoken.Claims{
		Iss: "http://localhost:8080",
		Sub: user.ID.Hex(),
	})

	if err != nil {
		return nil, errors.New("Error while trying to login")
	}

	return &AuthUserPayload{
		User:  user,
		Token: token,
	}, nil
}

func (r *mutationResolver) ValidateToken(ctx context.Context, token string) (*ValidateTokenPayload, error) {
	claims, err := jsonwebtoken.Decode(token)
	var user *models.User

	if err != nil {
		return &ValidateTokenPayload{
			Claims: &claims,
			User:   user,
			Valid:  false,
		}, nil
	}

	id, err := primitive.ObjectIDFromHex(claims.Sub)

	if err != nil {
		return nil, errors.New("Error while trying to validate the token")
	}

	user, err = userDao.FindByID(id)

	if err != nil {
		return &ValidateTokenPayload{
			Claims: &jsonwebtoken.Claims{},
			User:   user,
			Valid:  false,
		}, nil
	}

	return &ValidateTokenPayload{
		Claims: &claims,
		User:   user,
		Valid:  true,
	}, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	users := userDao.GetAll()
	return users, nil
}

type userResolver struct{ *Resolver }

func (r *userResolver) ID(ctx context.Context, obj *models.User) (string, error) {
	return obj.ID.Hex(), nil
}

func (r *userResolver) CreatedAt(ctx context.Context, obj *models.User) (string, error) {
	return obj.CreatedAt.Format("2006-01-02 15:04:05"), nil
}

func (r *userResolver) UpdatedAt(ctx context.Context, obj *models.User) (string, error) {
	return obj.UpdatedAt.Format("2006-01-02 15:04:05"), nil
}
