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
	AccessTokenLifetime = time.Minute * 15
)

// CreateDefaultClaims returns an default claims object for an subject
func CreateDefaultClaims(subject string) Claims {
	return Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    env.Config.ServerHost,
			Subject:   subject,
			ExpiresAt: time.Now().UTC().Add(AccessTokenLifetime).Unix(),
			IssuedAt:  time.Now().UTC().Unix(),
		},
	}
}
