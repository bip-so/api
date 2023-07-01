package auth

import (
	"encoding/json"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	user2 "gitlab.com/phonepost/bip-be-platform/internal/user"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"strings"
)

// Ok We didn't find the any user with this email
// First check if the the
func (s *authService) UserNotFoundWithEmailFlow(data SocialAuthPost) (*models.User, error) {
	var username string
	if data.Email == "" {
		username = data.UserName
	} else {
		if data.UserName == "" {
			components := strings.Split(data.Email, "@")
			username, _ = components[0], components[1]
		} else {
			username = data.UserName
		}
	}

	userNameExisits, _ := App.Service.getUserWithUserName(username)
	if userNameExisits != nil {
		username = username + "-" + utils.NewNanoid()
	}
	// Create User Instane
	user, errCreatingSocialUser := user2.App.Service.CreateNewSocialUser(data.Email, username, data.FullName, data.Image, data.ClientReferenceId)
	if errCreatingSocialUser != nil {
		return nil, errCreatingSocialUser
	}
	// Create User Autrh Instane
	metadata, _ := json.Marshal(map[string]interface{}{"accessToken": data.AccessToken})
	errSocialAuthInstancr := user2.App.Service.CreateUserSocialAuth(user.ID, data.Provider, data.ProviderID, metadata)
	if errSocialAuthInstancr != nil {
		return nil, errSocialAuthInstancr
	}

	return user, nil
}
