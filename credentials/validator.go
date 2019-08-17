package credentials

import (
	"fmt"

	"github.com/LucasFrezarini/go-auth-manager/dao"
	"github.com/LucasFrezarini/go-auth-manager/env"
	"github.com/LucasFrezarini/go-auth-manager/jsonwebtoken"
	"github.com/LucasFrezarini/go-auth-manager/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ValidateCredentials validate if a jwt claims is in the expected format, validating
// for instance: The issuer, if the user exists and if its active.
// If the claims are valid, a pointer to the validated user is returned
func ValidateCredentials(claims jsonwebtoken.Claims) (*models.User, error) {
	var user *models.User

	userDao := dao.UserDao{}

	if claims.Issuer != env.Config.ServerHost {
		return nil, fmt.Errorf("Error while validating claims credentials: Invalid issuer <%v>", claims.Issuer)
	}

	id, err := primitive.ObjectIDFromHex(claims.Subject)

	if err != nil {
		return nil, fmt.Errorf("Error while validating clams credentials: %v", err)
	}

	user, err = userDao.FindOne(models.User{
		ID:     id,
		Active: true,
	})

	if err != nil {
		return nil, fmt.Errorf("Error while validating clams credentials: %v", err)
	}

	return user, nil
}
