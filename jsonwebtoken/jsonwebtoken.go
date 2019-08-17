package jsonwebtoken

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("supersecretkey")

// Decode validates and decodes the claims of a token
func Decode(tokenString string) (Claims, error) {
	if tokenString == "" {
		return Claims{}, errors.New("Authorization token cannot be empty")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return jwtSecret, nil
	})

	if err != nil {
		return Claims{}, errors.New("Invalid authorization token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return Claims{
			StandardClaims: jwt.StandardClaims{
				Issuer:  claims["iss"].(string),
				Subject: claims["sub"].(string),
			},
		}, nil
	}

	return Claims{}, errors.New("Invalid authorization token")

}

// Encode creates a new jwt token with the provided claims
func Encode(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": claims.StandardClaims.Issuer,
		"sub": claims.StandardClaims.Subject,
	})

	tokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		return "", fmt.Errorf("Error while signing the jwt token: %v", err)
	}

	return tokenString, nil
}
