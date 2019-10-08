package mutation_test

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/LucasFrezarini/go-auth-manager/jsonwebtoken"
	"github.com/LucasFrezarini/go-auth-manager/middlewares"
	"github.com/stretchr/testify/require"
)

func TestRefreshToken(t *testing.T) {
	srv := httptest.NewServer(middlewares.MakeHandlers())
	c := client.New(srv.URL)

	t.Run("Shouldn't be able to get a new token if a refresh token is valid, but not present on the database", func(t *testing.T) {
		var errorResponse []struct {
			Message    string   `json:"message"`
			Path       []string `json:"path"`
			Extensions struct {
				Code string `json:"code"`
			} `json:"extensions"`
		}

		// Will generate a valid refresh token, but not inserted on the database
		token, err := jsonwebtoken.Encode(jsonwebtoken.CreateRefreshTokenClaims("5d4a22e9587f3dbb8d33fd39"))

		if err != nil {
			t.FailNow()
		}

		err = c.Post(fmt.Sprintf(`
			mutation {
				refreshToken(refreshToken: "%s") {
					token
					refreshToken
					user {
					  id
					  email
					  roles
					}
				  }
			}
		`, token), &errorResponse)

		json.Unmarshal([]byte(err.Error()), &errorResponse)
		require.Equal(t, 1, len(errorResponse))

		resp := errorResponse[0]

		require.Equal(t, []string{"refreshToken"}, resp.Path)
		require.Equal(t, "UNAUTHORIZED", resp.Extensions.Code)
	})

	t.Run("Should be able to get a new token with an refresh token", func(t *testing.T) {
		t.SkipNow()
		var resp struct {
			RefreshToken struct {
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

		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Njg4NDMzNjcsImlhdCI6MTU2NjE2NDk2NywiaXNzIjoiaHR0cDovL3Rlc3QuaW8iLCJzdWIiOiI1ZDRhMjJlOTU4N2YzZGJiOGQzM2ZkMzkifQ.qqGHNJ5WJO5KjfylnW7klRvegeVlMVx8Fv4ehnuigpA"

		c.MustPost(fmt.Sprintf(`
			mutation {
				refreshToken(refreshToken: "%s") {
					token
					refreshToken
					user {
					  id
					  email
					  roles
					}
				  }
			}
		`, token), &resp)

		require.NotEmpty(t, resp.RefreshToken.Token)
		require.Equal(t, token, resp.RefreshToken.RefreshToken)
		require.Equal(t, "5d4a22e9587f3dbb8d33fd39", resp.RefreshToken.User.ID)
	})
}
