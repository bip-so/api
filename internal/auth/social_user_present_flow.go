package auth

import (
	"encoding/json"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	user2 "gitlab.com/phonepost/bip-be-platform/internal/user"
)

// If yes: Do we have a UserSocialAuth Instance
// If yes: (Login the intended user with email)
// If no: Create a UserSocialAuth Instance and (Login the intented user with email)
func (s *authService) UserFoundWithEmailFlow(user *models.User, data SocialAuthPost) (*models.User, error) {
	flag := user2.App.Service.DoesUSAExits(user.ID, data.Provider, data.ProviderID)
	if flag {
		return user, nil
	} else {
		metadata, _ := json.Marshal(map[string]interface{}{"accessToken": data.AccessToken})
		errSocialAuthInstancr := user2.App.Service.CreateUserSocialAuth(user.ID, data.Provider, data.ProviderID, metadata)
		if errSocialAuthInstancr != nil {
			return nil, errSocialAuthInstancr
		}
	}
	return user, nil
}
