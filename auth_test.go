package manager

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/handler"
	"github.com/LucasFrezarini/go-auth-manager/generated"
	"github.com/LucasFrezarini/go-auth-manager/resolvers"
	"github.com/stretchr/testify/require"
)

func TestAuth(t *testing.T) {
	srv := httptest.NewServer(handler.GraphQL(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers.Resolver{}})))
	c := client.New(srv.URL)

	t.Run("Should create a user", func(t *testing.T) {
		var resp struct {
			CreateUser struct {
				User struct {
					ID        string
					Email     string
					Roles     []string
					CreatedAt string
					UpdatedAt string
				}
				Token string
			}
		}

		c.MustPost(`
			mutation {
				createUser(data:{
					email: "test6@email.com"
					password: "12345"
					roles: ["user"]
				  }) {
					user {
					  id
					  email
					  roles
					  createdAt
					  updatedAt
					}
					token
				  }
			}
		`, &resp)

		fmt.Print(resp)

		require.Equal(t, "test6@email.com", resp.CreateUser.User.Email)

	})
}
