package resolvers

import (
	"context"

	"github.com/LucasFrezarini/go-auth-manager/models"
)

// QueryResolver defines the root resolver from the query in GraphQL schema
type queryResolver struct{ *Resolver }

// Users is the resolver of Users on graphql schema
func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	users := userDao.GetAll()
	return users, nil
}
