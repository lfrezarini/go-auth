package jsonwebtoken

import (
	"github.com/dgrijalva/jwt-go"
)

// Claims represents the claims that can be passed on the creation of jsonwebtoken
type Claims struct {
	jwt.StandardClaims
}
