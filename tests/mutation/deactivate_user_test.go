package mutation_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/99designs/gqlgen/client"
	"github.com/LucasFrezarini/go-auth-manager/jsonwebtoken"
	"github.com/LucasFrezarini/go-auth-manager/middlewares"
	tests "github.com/LucasFrezarini/go-auth-manager/tests/helpers"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
)

func TestDeactivateUser(t *testing.T) {
	srv := httptest.NewServer(middlewares.MakeHandlers())
	c := client.New(srv.URL)

	t.Run("Should verify if token is present before deactivate user", func(t *testing.T) {
		var resp tests.ErrorResponse

		err := c.Post(`
			mutation {
				deactivateUser() {
					id
					active
					createdAt
				}
			}
		`, &resp)

		json.Unmarshal([]byte(err.Error()), &resp)

		require.Equal(t, 1, len(resp))

		errorResponse := resp[0]

		require.Equal(t, "Unauthorized", errorResponse.Message)
		require.Equal(t, "deactivateUser", errorResponse.Path[0])
		require.Equal(t, "UNAUTHORIZED", errorResponse.Extensions.Code)
	})

	t.Run("Should not allow the deactivate if the user is using an token with different issuer", func(t *testing.T) {
		var expectedResponse struct {
			Data   interface{} `json:"data"`
			Errors []struct {
				Message    string   `json:"message"`
				Path       []string `json:"path"`
				Extensions struct {
					Code string `json:"code"`
				} `json:"extensions"`
			} `json:"errors"`
		}

		httpClient := tests.HTTPClient{}

		query := `
			mutation {
				deactivateUser {
					id
					email
					roles
					active
					createdAt
					updatedAt
				}
			}
		`

		token, err := jsonwebtoken.Encode(jsonwebtoken.Claims{
			StandardClaims: jwt.StandardClaims{
				Issuer:    "http://another.site.io",
				Subject:   "5d4a22e9587f3dbb8d33fd39",
				IssuedAt:  time.Now().UTC().Unix(),
				ExpiresAt: time.Now().UTC().Add(15 * time.Minute).Unix(),
			},
		})

		if err != nil {
			t.Fatalf("Error while trying to get the token for test: %v", err)
		}

		headers := map[string]string{
			"Authorization": token,
			"Content-Type":  "application/json",
		}

		response, err := httpClient.DoRequest(srv.URL, query, headers)

		if err != nil {
			t.Fatalf("Error while doing the request: %v", err)
		}

		err = json.Unmarshal(response, &expectedResponse)

		if err != nil {
			t.Fatalf("Error while trying to Unmarshal response: %v", err)
		}

		require.Empty(t, expectedResponse.Data)
		require.NotZero(t, len(expectedResponse.Errors))

		errResponse := expectedResponse.Errors[0]

		require.Equal(t, errResponse.Message, "Unauthorized")
		require.Equal(t, errResponse.Path, []string{"deactivateUser"})
		require.Equal(t, errResponse.Extensions.Code, "UNAUTHORIZED")
	})

	t.Run("Should not allow the deactivate if the user is using an expired token", func(t *testing.T) {
		var expectedResponse struct {
			Data   interface{} `json:"data"`
			Errors []struct {
				Message    string   `json:"message"`
				Path       []string `json:"path"`
				Extensions struct {
					Code string `json:"code"`
				} `json:"extensions"`
			} `json:"errors"`
		}

		httpClient := tests.HTTPClient{}

		query := `
			mutation {
				deactivateUser {
					id
					email
					roles
					active
					createdAt
					updatedAt
				}
			}
		`

		token, err := jsonwebtoken.Encode(jsonwebtoken.Claims{
			StandardClaims: jwt.StandardClaims{
				Issuer:    "http://test.io",
				Subject:   "5d4a22e9587f3dbb8d33fd39",
				ExpiresAt: time.Now().UTC().Add(-15 * time.Minute).Unix(),
				IssuedAt:  time.Now().UTC().Add(-time.Hour).Unix(),
			},
		})

		if err != nil {
			t.Fatalf("Error while trying to get the token for test: %v", err)
		}

		headers := map[string]string{
			"Authorization": token,
			"Content-Type":  "application/json",
		}

		response, err := httpClient.DoRequest(srv.URL, query, headers)

		if err != nil {
			t.Fatalf("Error while doing the request: %v", err)
		}

		err = json.Unmarshal(response, &expectedResponse)

		if err != nil {
			t.Fatalf("Error while trying to Unmarshal response: %v", err)
		}

		require.Empty(t, expectedResponse.Data)
		require.NotZero(t, len(expectedResponse.Errors))

		errResponse := expectedResponse.Errors[0]

		require.Equal(t, errResponse.Message, "Unauthorized")
		require.Equal(t, errResponse.Path, []string{"deactivateUser"})
		require.Equal(t, errResponse.Extensions.Code, "UNAUTHORIZED")
	})

	t.Run("Should deactivate user if token is present", func(t *testing.T) {
		var expectedResponse struct {
			Data struct {
				DeactivateUser struct {
					ID        string   `json:"id"`
					Email     string   `json:"email"`
					Roles     []string `json:"roles"`
					Active    bool     `json:"active"`
					CreatedAt string   `json:"createdAt"`
					UpdatedAt string   `json:"updatedAt"`
				} `json:"deactivateUser"`
			} `json:"data"`
		}

		httpClient := tests.HTTPClient{}

		query := `
			mutation {
				deactivateUser {
					id
					email
					roles
					active
					createdAt
					updatedAt
				}
			}
		`

		token, err := jsonwebtoken.Encode(jsonwebtoken.Claims{
			StandardClaims: jwt.StandardClaims{
				Issuer:    "http://test.io",
				Subject:   "5d4a22e9587f3dbb8d33fd39",
				IssuedAt:  time.Now().UTC().Unix(),
				ExpiresAt: time.Now().UTC().Add(15 * time.Minute).Unix(),
			},
		})

		if err != nil {
			t.Fatalf("Error while trying to get the token for test: %v", err)
		}

		headers := map[string]string{
			"Authorization": token,
			"Content-Type":  "application/json",
		}

		response, err := httpClient.DoRequest(srv.URL, query, headers)

		if err != nil {
			t.Fatalf("Error while doing the request: %v", err)
		}

		err = json.Unmarshal(response, &expectedResponse)

		if err != nil {
			t.Fatalf("Error while trying to Unmarshal response: %v", err)
		}

		data := expectedResponse.Data.DeactivateUser

		require.Equal(t, "5d4a22e9587f3dbb8d33fd39", data.ID)
		require.Equal(t, "test4@test.com", data.Email)
		require.False(t, data.Active)
	})
}
