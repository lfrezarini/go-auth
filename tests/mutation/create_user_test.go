package mutation_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/LucasFrezarini/go-auth-manager/middlewares"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	srv := httptest.NewServer(middlewares.MakeHandlers())
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
					email: "test@email.com"
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

		require.NotEmpty(t, resp.CreateUser.User.ID)
		require.Equal(t, "test@email.com", resp.CreateUser.User.Email)
		require.Equal(t, []string{"user"}, resp.CreateUser.User.Roles)
		require.NotEmpty(t, resp.CreateUser.Token)
	})

	t.Run("Should not create a new user if email already exists", func(t *testing.T) {
		var errorResponse []struct {
			Message    string   `json:"message"`
			Path       []string `json:"path"`
			Extensions struct {
				Code string `json:"code"`
			} `json:"extensions"`
		}

		err := c.Post(`
			mutation {
				createUser(data:{
					email: "test@email.com"
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
		`, &errorResponse)

		json.Unmarshal([]byte(err.Error()), &errorResponse)
		require.NotZero(t, len(errorResponse))

		response := errorResponse[0]

		require.Equal(t, response.Message, "User already exists")
		require.Equal(t, response.Path[0], "createUser")
		require.Equal(t, response.Extensions.Code, "CONFLICT")
	})
}
