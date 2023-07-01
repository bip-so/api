package notifications

import (
	"errors"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

// CalculateNotificationEntityAccess returns the notification media that user wants the notifications want to send.
// returns : {
// 		"app": true,
// 		"email": true,
// 		"discord": false,
//	}
// We can calculate and store it in redis.
// Invalidation is done on updating the userSettings.
func (s notificationService) CalculateNotificationEntityAccess(userID uint64, entity string) (map[string]bool, error) {
	userSettings, err := s.GetUserSettingsController(userID)
	if err != nil {
		return nil, err
	}
	access := map[string]bool{
		"app":     false,
		"email":   false,
		"discord": false,
	}
	for _, userSetting := range *userSettings {
		entityValue := utils.GetValue(userSetting, entity)
		access[userSetting.Type] = entityValue.Bool()
	}
	return access, nil
}

func (s notificationService) GetUserSettingsController(userID uint64) (*[]models.UserSettings, error) {

	userSettings, err := App.Repo.GetUserSettings(userID)
	if err != nil {
		return nil, err
	}

	if len(*userSettings) == 0 {
		for _, notificationType := range models.DefaultNotificationTypes {
			userSettings := models.NewDefaultUserSettings(userID, notificationType)
			done := App.Repo.CreateNewUserSettings(userSettings)
			if !done {
				return nil, errors.New("can't fetch user settings")
			}
		}
		userSettings, err = App.Repo.GetUserSettings(userID)
		if err != nil {
			return nil, err
		}
	}

	return userSettings, nil
}
