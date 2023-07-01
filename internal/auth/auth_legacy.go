package auth

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	user2 "gitlab.com/phonepost/bip-be-platform/internal/user"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func TryLegacySignup(email string, password string, clientReferenceId string) (*models.User, error) {
	var err error
	var user *models.User
	var username string

	// Unique user name
	components := strings.Split(email, "@")
	username, _ = components[0], components[1]
	userExisits, _ := App.Service.getUserWithUserName(username)
	if userExisits != nil {
		username = username + "-" + utils.NewNanoid()
	}
	// Unique user name

	// user object created
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user = user2.App.Service.CreateNewUser(email, string(hashedPassword), username, "", clientReferenceId)
	if err != nil {
		return nil, err
	}

	return user, nil

}
