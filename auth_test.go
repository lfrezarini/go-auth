package manager

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/handler"
	"github.com/LucasFrezarini/go-auth-manager/generated"
	"github.com/LucasFrezarini/go-auth-manager/jsonwebtoken"
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
				Token string
			}
		}

		c.MustPost(`
			mutation {
				login(data:{
					email: "test@email.com"
					password: "12345"
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

		require.NotEmpty(t, resp.Login.User.ID)
		require.Equal(t, "test@email.com", resp.Login.User.Email)
		require.Equal(t, []string{"user"}, resp.Login.User.Roles)
		require.NotEmpty(t, resp.Login.Token)
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
					email: "test@email.com"
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

	t.Run("Should return true in validate if token is valid", func(t *testing.T) {
		var resp struct {
			ValidateToken struct {
				User struct {
					ID        string
					Email     string
					Roles     []string
					CreatedAt string
					UpdatedAt string
				}
				Claims struct {
					Iss string
					Sub string
				}
				Valid bool
			}
		}

		token, err := jsonwebtoken.Encode(jsonwebtoken.Claims{
			Iss: "http://test.io",
			Sub: "5d470b3e98b0116d7d8ca48c",
		})

		if err != nil {
			t.FailNow()
		}

		c.MustPost(fmt.Sprintf(`
			mutation {
				validateToken(token: "%s") {
					user {
					  id
					  email
					  roles
					  createdAt
					  updatedAt
					}
					claims {
						iss
						sub
					}
					valid
				  }
			}
		`, token), &resp)

		validateToken := resp.ValidateToken

		require.Equal(t, validateToken.User.ID, "5d470b3e98b0116d7d8ca48c")
		require.Equal(t, "test1@test.com", validateToken.User.Email)
		require.Equal(t, []string{"user"}, validateToken.User.Roles)

		require.Equal(t, validateToken.Claims.Iss, "http://test.io")
		require.Equal(t, validateToken.Claims.Sub, "5d470b3e98b0116d7d8ca48c")
		require.True(t, validateToken.Valid)
	})

	t.Run("Should return false in validate if token is invalid and no infos", func(t *testing.T) {
		var resp struct {
			ValidateToken struct {
				User struct {
					ID        string
					Email     string
					Roles     []string
					CreatedAt string
					UpdatedAt string
				}
				Claims struct {
					Iss string
					Sub string
				}
				Valid bool
			}
		}

		c.MustPost(fmt.Sprintf(`
			mutation {
				validateToken(token: "invalid") {
					user {
					  id
					  email
					  roles
					  createdAt
					  updatedAt
					}
					claims {
						iss
						sub
					}
					valid
				  }
			}
		`), &resp)

		validateToken := resp.ValidateToken

		require.Empty(t, validateToken.User.ID)
		require.Empty(t, validateToken.User.Email)
		require.Empty(t, validateToken.User.Roles)
		require.Empty(t, validateToken.Claims.Iss)
		require.Empty(t, validateToken.Claims.Sub)
		require.False(t, validateToken.Valid)
	})
}
