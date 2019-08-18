package jsonwebtoken_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/LucasFrezarini/go-auth-manager/jsonwebtoken"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
)

var jwtSecret = []byte("supersecretkey")

func TestEncodeToken(t *testing.T) {
	t.Run("Should encode a token with the passed claims", func(t *testing.T) {
		expirationTime := time.Now().UTC().Add(1 * time.Minute).Unix()

		claims := jsonwebtoken.Claims{
			StandardClaims: jwt.StandardClaims{
				Issuer:    "http://test.io",
				Subject:   "test",
				ExpiresAt: expirationTime,
			},
		}

		token, err := jsonwebtoken.Encode(claims)

		require.Empty(t, err)
		require.NotEmpty(t, token)

		// Reverse engineering token to check the saved claims
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return jwtSecret, nil
		})

		require.Empty(t, err)

		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
			require.Equal(t, "http://test.io", claims["iss"])
			require.Equal(t, "test", claims["sub"])
			require.Equal(t, float64(expirationTime), claims["exp"])
		} else {
			t.FailNow()
		}
	})
}

func TestDecodeToken(t *testing.T) {
	t.Run("Should decode a valid token and retrieve the claims", func(t *testing.T) {
		/*
			Generated with key: "supersecretkey"
			The claims saved in this token are:
			Issuer: "http://test.io"
			Subject: "5d470b3e98b0116d7d8ca48c"
		*/
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vdGVzdC5pbyIsInN1YiI6IjVkNDcwYjNlOThiMDExNmQ3ZDhjYTQ4YyJ9.R2bXULilY_kNrpg20uMxqZuy6UOCGNDH4Hrb9FU_aiQ"

		claims, err := jsonwebtoken.Decode(token)

		require.Empty(t, err)
		require.NotEmpty(t, claims)
		require.Equal(t, "http://test.io", claims.Issuer)
		require.Equal(t, "5d470b3e98b0116d7d8ca48c", claims.Subject)
	})

	t.Run("Should return error if the token is an empty string", func(t *testing.T) {
		claims, err := jsonwebtoken.Decode("")

		require.Empty(t, claims)
		require.NotEmpty(t, err)
	})

	t.Run("Should return error if the token is invalid", func(t *testing.T) {
		claims, err := jsonwebtoken.Decode("eyJhbGciOiJIUnR5cCI6IkpXVCJ9.pc3MiOiJod")

		require.Empty(t, claims)
		require.NotEmpty(t, err)
	})

	t.Run("Should return error if the token is expired", func(t *testing.T) {
		/*
			Generated with key: "supersecretkey"
			The claims saved in this token are:
			Issuer: "http://test.io"
			Subject: "5d470b3e98b0116d7d8ca48c"
			ExpiresAt: <some expired unix time>
		*/
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjYwODUxOTMsImlzcyI6Imh0dHA6Ly90ZXN0LmlvIiwic3ViIjoiNWQ0NzBiM2U5OGIwMTE2ZDdkOGNhNDhjIn0.ai6WbCdhbqRNgRZZZeqUyb-fmkVutPELRzQL6sVVP3M"

		claims, err := jsonwebtoken.Decode(token)

		require.Empty(t, claims)
		require.NotEmpty(t, err)
	})

	t.Run("Should return error if the token is signed with an unexpected algoritm (other than HMAC-SHA)", func(t *testing.T) {
		/*
			Generated with key: "supersecretkey"
			algoritm: RSA256
			The claims saved in this token are:
			Issuer: "http://test.io"
			Subject: "5d470b3e98b0116d7d8ca48c"
		*/
		token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vdGVzdC5pbyIsInN1YiI6IjVkNDcwYjNlOThiMDExNmQ3ZDhjYTQ4YyJ9.fY-dsbIWAHHEljM4UdG0ivgSZeTFRtD9KqQQF9rzrqHjhk57gpNdfTDv-Cohaou3ocuootHNEsOiaUJzbWV53EOS4hxCbI1YMx8dQQ0O2fgRHrGCcCjIvbB3xmKwkBYIx0s1yu6rfPA52iu7r32gUrBzfdtddN6ilYVAJffWYxBrdEnpunlNC70SGLwplVHekY7jgL580l0apa-onJwlK49hULAomOGEB8e_Naj9s1kh719zdeUDAkw8MzUN2I-IS1o6rnggV5OcADUaRDqD8vYwwPWea6ziJgNQkOO8v2KHTsHkUVYqH1cdYN3egIEyJExNiZD6zgpATbaSMFt60w"

		claims, err := jsonwebtoken.Decode(token)

		require.Empty(t, claims)
		require.NotEmpty(t, err)
	})
}
