package mutation_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/LucasFrezarini/go-auth-manager/middlewares"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	srv := httptest.NewServer(middlewares.MakeHandlers())
	c := client.New(srv.URL)

	t.Run("Should be able to login with an existing user", func(t *testing.T) {
		var resp struct {
			Login struct {
				User struct {
					ID        string
					Email     string
					Roles     []string
					CreatedAt string
					UpdatedAt string
				}
				Token        string
				RefreshToken string
			}
		}

		c.MustPost(`
			mutation {
				login(data:{
					email:"test1@test.com"
					password:"12345"
				  }) {
					token
					refreshToken
					user {
					  id
					  email
					  roles
					}
				  }
			}
		`, &resp)

		require.NotEmpty(t, resp.Login.User.ID)
		require.Equal(t, "test1@test.com", resp.Login.User.Email)
		require.Equal(t, []string{"user"}, resp.Login.User.Roles)
		require.NotEmpty(t, resp.Login.Token)
		require.NotEmpty(t, resp.Login.RefreshToken)
	})

	t.Run("Should not allow the user to login if password is invalid", func(t *testing.T) {
		var errorResponse []struct {
			Message    string   `json:"message"`
			Path       []string `json:"path"`
			Extensions struct {
				Code string `json:"code"`
			} `json:"extensions"`
		}

		err := c.Post(`
			mutation {
				login(data:{
					email: "test1@test.com"
					password: "1234"
				}) {
					user {
						id
					}
					token
				}
			}
		`, &errorResponse)

		json.Unmarshal([]byte(err.Error()), &errorResponse)

		require.Equal(t, len(errorResponse), 1)
		require.Equal(t, errorResponse[0].Message, "Unauthorized")

		require.Equal(t, len(errorResponse[0].Path), 1)
		require.Equal(t, errorResponse[0].Path[0], "login")
		require.Equal(t, errorResponse[0].Extensions.Code, "UNAUTHORIZED")
	})

	t.Run("Should not allow the user to login if email is invalid", func(t *testing.T) {
		var errorResponse []struct {
			Message    string   `json:"message"`
			Path       []string `json:"path"`
			Extensions struct {
				Code string `json:"code"`
			} `json:"extensions"`
		}

		err := c.Post(`
			mutation {
				login(data:{
					email: "test@random.com"
					password: "12345"
				}) {
					user {
						id
					}
					token
				}
			}
		`, &errorResponse)

		json.Unmarshal([]byte(err.Error()), &errorResponse)

		require.Equal(t, len(errorResponse), 1)
		require.Equal(t, errorResponse[0].Message, "Unauthorized")

		require.Equal(t, len(errorResponse[0].Path), 1)
		require.Equal(t, errorResponse[0].Path[0], "login")
		require.Equal(t, errorResponse[0].Extensions.Code, "UNAUTHORIZED")
	})

	t.Run("Should not allow the user to login if user is deactivated", func(t *testing.T) {
		var errorResponse []struct {
			Message    string   `json:"message"`
			Path       []string `json:"path"`
			Extensions struct {
				Code string `json:"code"`
			} `json:"extensions"`
		}

		err := c.Post(`
			mutation {
				login(data:{
					email: "test2@test.com"
					password: "12345"
				}) {
					user {
						id
					}
					token
				}
			}
		`, &errorResponse)

		json.Unmarshal([]byte(err.Error()), &errorResponse)

		require.Equal(t, 1, len(errorResponse))
		require.Equal(t, "Unauthorized", errorResponse[0].Message)

		require.Equal(t, 1, len(errorResponse[0].Path))
		require.Equal(t, "login", errorResponse[0].Path[0])
		require.Equal(t, "UNAUTHORIZED", errorResponse[0].Extensions.Code)
	})
}
