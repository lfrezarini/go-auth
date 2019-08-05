package mutation_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/LucasFrezarini/go-auth-manager/jsonwebtoken"
	"github.com/LucasFrezarini/go-auth-manager/middlewares"
	tests "github.com/LucasFrezarini/go-auth-manager/tests/helpers"
	"github.com/stretchr/testify/require"
)

func TestUpdateUser(t *testing.T) {
	srv := httptest.NewServer(middlewares.MakeHandlers())
	c := client.New(srv.URL)

	t.Run("Should verify if token is present before update user", func(t *testing.T) {
		var resp tests.ErrorResponse

		err := c.Post(`
			mutation {
				updateUser(data:{ password: "1234" }) {
					id
					createdAt
				}
			}
		`, &resp)

		json.Unmarshal([]byte(err.Error()), &resp)

		require.Equal(t, 1, len(resp))

		errorResponse := resp[0]

		require.Equal(t, "Unauthorized", errorResponse.Message)
		require.Equal(t, "updateUser", errorResponse.Path[0])
		require.Equal(t, "UNAUTHORIZED", errorResponse.Extensions.Code)
	})

	t.Run("Should allow the user to update his password", func(t *testing.T) {
		var expectedResponse struct {
			Data struct {
				UpdateUser struct {
					ID        string   `json:"id"`
					Email     string   `json:"email"`
					Roles     []string `json:"roles"`
					CreatedAt string   `json:"createdAt"`
					UpdatedAt string   `json:"updatedAt"`
				} `json:"updateUser"`
			} `json:"data"`
		}

		httpClient := tests.HTTPClient{}

		query := `
			mutation {
				updateUser(data:{
					password: "changed"
				}) {
					id
					email
					roles
					createdAt
					updatedAt
				}
			}
		`

		token, err := jsonwebtoken.Encode(jsonwebtoken.Claims{
			Iss: "http://test.io",
			Sub: "5d470b3e98b0116d7d8ca48c",
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

		data := expectedResponse.Data.UpdateUser

		require.Equal(t, "5d470b3e98b0116d7d8ca48c", data.ID)
		require.Equal(t, "test1@test.com", data.Email)
	})
}
