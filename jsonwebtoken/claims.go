package jsonwebtoken

import (
	"time"

	"github.com/LucasFrezarini/go-auth-manager/env"
	"github.com/dgrijalva/jwt-go"
)

// Claims represents the claims that can be passed on the creation of jsonwebtoken
type Claims struct {
	jwt.StandardClaims
}

const (
	// AccessTokenLifetime represents the lifetime, in minutes, of an access token
	AccessTokenLifetime = time.Minute * 300

	// RefreshTokenLifetimeInMonths represents the lifetime, in months, of an refresh token
	RefreshTokenLifetimeInMonths = 1
)

// CreateDefaultClaims returns an default claims object for an subject
func CreateDefaultClaims(subject string) (claims Claims) {
	claims = createCommonClains(subject)
	claims.ExpiresAt = time.Now().UTC().Add(AccessTokenLifetime).Unix()

	return
}

func createCommonClains(subject string) Claims {
	return Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:   env.Config.ServerHost,
			Subject:  subject,
			IssuedAt: time.Now().UTC().Unix(),
		},
	}
}

// CreateRefreshTokenClaims returns a claims object for an subject for an refresh token
func CreateRefreshTokenClaims(subject string) (claims Claims) {
	claims = createCommonClains(subject)
	claims.ExpiresAt = time.Now().UTC().AddDate(0, RefreshTokenLifetimeInMonths, 0).Unix()

	return
}
