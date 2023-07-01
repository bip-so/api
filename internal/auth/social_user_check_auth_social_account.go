package auth

import (
	"encoding/json"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	user2 "gitlab.com/phonepost/bip-be-platform/internal/user"
)

func (s *authService) UserNotFoundWithEmailMayHaveSocialAuthInstanceFlow(data SocialAuthPost) *models.User {
	userID := user2.App.Service.DoesUSAExitsWithOutUserID(data.Provider, data.ProviderID)
	if userID == 0 {
		return nil
	}
	metadata, _ := json.Marshal(map[string]interface{}{"accessToken": data.AccessToken})
	user2.App.Service.UpdateUSA(map[string]interface{}{"provider_id": data.ProviderID, "provider_name": data.Provider}, map[string]interface{}{"metadata": metadata})
	user, _ := App.Service.getUserWithUserID(userID)
	if user == nil {
		return nil
	} else {
		return user
	}
}
