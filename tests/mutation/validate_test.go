package mutation_test

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/LucasFrezarini/go-auth-manager/jsonwebtoken"
	"github.com/LucasFrezarini/go-auth-manager/middlewares"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
)

type validateResponse struct {
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

func requireResponseIsInvalid(t *testing.T, resp validateResponse) {
	validateToken := resp.ValidateToken

	require.Empty(t, validateToken.User.ID)
	require.Empty(t, validateToken.User.Email)
	require.Empty(t, validateToken.User.Roles)
	require.Empty(t, validateToken.Claims.Iss)
	require.Empty(t, validateToken.Claims.Sub)
	require.False(t, validateToken.Valid)
}

func TestValidate(t *testing.T) {
	srv := httptest.NewServer(middlewares.MakeHandlers())
	c := client.New(srv.URL)

	t.Run("Should return true in validate if token is valid", func(t *testing.T) {
		var resp validateResponse

		token, err := jsonwebtoken.Encode(jsonwebtoken.Claims{
			StandardClaims: jwt.StandardClaims{
				Issuer:  "http://test.io",
				Subject: "5d470b3e98b0116d7d8ca48c",
			},
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
		var resp validateResponse

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

		requireResponseIsInvalid(t, resp)
	})

	t.Run("Should return false and empty fields if the owner of the token is inactive", func(t *testing.T) {
		var resp validateResponse

		token, err := jsonwebtoken.Encode(jsonwebtoken.Claims{
			StandardClaims: jwt.StandardClaims{
				Issuer:  "http://test.io",
				Subject: "5d4a22b1106eded67d47c02e", // test2@test.com is inactive on seed.js
			},
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

		requireResponseIsInvalid(t, resp)
	})

	t.Run("Token should be invalid if the sub id doesn't exists on the database", func(t *testing.T) {
		var resp validateResponse

		token, err := jsonwebtoken.Encode(jsonwebtoken.Claims{
			StandardClaims: jwt.StandardClaims{
				Issuer:  "http://test.io",
				Subject: "5d5035c42dd89da99833d081", // Valid objectID but no user is registered with it
			},
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

		requireResponseIsInvalid(t, resp)
	})

	t.Run("Token should be invalid if the issuer from token is different than the auth manager server", func(t *testing.T) {
		var resp validateResponse

		token, err := jsonwebtoken.Encode(jsonwebtoken.Claims{
			StandardClaims: jwt.StandardClaims{
				Issuer:  "http://another.site.io",
				Subject: "5d470b3e98b0116d7d8ca48c",
			},
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

		requireResponseIsInvalid(t, resp)
	})
}
