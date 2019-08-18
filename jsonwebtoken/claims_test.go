package jsonwebtoken_test

import (
	"testing"
	"time"

	"github.com/LucasFrezarini/go-auth-manager/jsonwebtoken"
	"github.com/stretchr/testify/require"
)

func TestCreateDefaultClaims(t *testing.T) {
	t.Run("Should create claims with the subject passed as parameter", func(t *testing.T) {
		claims := jsonwebtoken.CreateDefaultClaims("321")

		require.Equal(t, "321", claims.Subject)
	})

	t.Run("Should create claims with the expiration time of 15 minutes", func(t *testing.T) {
		claims := jsonwebtoken.CreateDefaultClaims("321")

		require.Equal(t, time.Now().UTC().Add(15*time.Minute).Unix(), claims.ExpiresAt)
	})

	t.Run("Should create claims with an issue date", func(t *testing.T) {
		claims := jsonwebtoken.CreateDefaultClaims("321")

		require.NotEmpty(t, claims.IssuedAt)
	})
}

func TestCreateRefreshTokenClaims(t *testing.T) {
	t.Run("Should create claims with the subject passed as parameter", func(t *testing.T) {
		claims := jsonwebtoken.CreateRefreshTokenClaims("321")

		require.Equal(t, "321", claims.Subject)
	})

	t.Run("Should create claims with 1 month of expiration time", func(t *testing.T) {
		claims := jsonwebtoken.CreateRefreshTokenClaims("321")

		require.Equal(t, time.Now().UTC().AddDate(0, 1, 0).Unix(), claims.ExpiresAt)
	})

	t.Run("Should create claims with an issue date", func(t *testing.T) {
		claims := jsonwebtoken.CreateRefreshTokenClaims("321")

		require.NotEmpty(t, claims.IssuedAt)
	})
}
