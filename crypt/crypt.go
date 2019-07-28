package crypt

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword encrypt the password using bcrypt with the default salt value, returning a string representating the hashed password
func HashPassword(plainPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Panic(err)
		return "", fmt.Errorf("Error while trying to encrypt the password: %v", err)
	}

	return string(hash), nil
}

// ComparePassword compares the hashedPasword against the plainText password, returning true if the password matches and false otherwise
func ComparePassword(hashedPasword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPasword), []byte(plainPassword))

	return err == nil
}
