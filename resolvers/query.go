package resolvers

import (
	"context"

	"github.com/LucasFrezarini/go-auth-manager/gqlerrors"
	"github.com/LucasFrezarini/go-auth-manager/models"
)

// QueryResolver defines the root resolver from the query in GraphQL schema
type queryResolver struct{ *Resolver }

// Users is the resolver of Users on graphql schema
func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	users, err := userDao.GetAll()

	if err != nil {
		return nil, gqlerrors.CreateInternalServerError("Error while trying to fetch all users")
	}

	return users, nil
}
