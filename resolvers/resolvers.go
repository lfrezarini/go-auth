package resolvers

import (
	"github.com/LucasFrezarini/go-auth-manager/dao"
	"github.com/LucasFrezarini/go-auth-manager/generated"
	"github.com/LucasFrezarini/go-auth-manager/models"
)

var userDao dao.UserDao

func init() {
	userDao = dao.UserDao{}
}

// Resolver is the structure of the graphql root resolver
type Resolver struct {
	users []*models.User
}

// Mutation returns the root mutation resolver from GraphQL schema
func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}

// Query returns the root query resolver from GraphQL schema
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

// User returns the user resolver from GraphQL schema
func (r *Resolver) User() generated.UserResolver {
	return &userResolver{r}
}

//Claims returns the claims resolver from GraphQL schema
func (r *Resolver) Claims() generated.ClaimsResolver {
	return &claimsResolver{r}
}
