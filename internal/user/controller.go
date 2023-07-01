package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gorm.io/gorm"
	"path/filepath"
	"strconv"
	"time"

	"gitlab.com/phonepost/bip-be-platform/pkg/utils"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/s3"
)

func (c userController) UserInfoController(userID string, userName string, authUserID uint64) (*UserGetSerializer, error) {
	var user *models.User
	var err error

	if userName != "" {
		user, err = App.Repo.GetUser(map[string]interface{}{"username": userName})
		if err != nil {
			return nil, err
		}
	} else if userID != "" {
		// Commented
		// user, err = App.Repo.GetUser(map[string]interface{}{"id": userID})
		user, err = queries.App.UserQueries.GetUserByID(utils.Uint64(userID))

		if err != nil {
			return nil, err
		}
	}

	// Get user followers
	followService := NewFollowService()
	userFollowers := FollowUserService.FollowerCountUse(followService, *user)

	userData := UserProfileGetSerializerData(user, userFollowers)
	isFollowing := false
	if authUserID != utils.Uint64(userID) {
		follower, _ := App.Repo.GetUserFollows(utils.Uint64(userID), authUserID)
		if follower.ID != 0 {
			isFollowing = true
		}
	}
	userData.IsFollowing = &isFollowing

	return &userData, nil
}

func (c userController) UpdateUserController(requestBody *UpdateUserValidator, user *models.User) (*UserGetSerializer, error) {
	if requestBody.File.Size > 0 {
		file, err := requestBody.File.Open()
		if err != nil {
			logger.Error(fmt.Sprintf("Error on reading user avatar url %s", err.Error()))
		} else {

			//userImagePath := fmt.Sprintf("user/%s/%s", user.UUID, utils.NewNanoid())
			//userImagePath := fmt.Sprintf("user/%s/%s", user.UUID, requestBody.File.Filename)
			Fext := filepath.Ext(requestBody.File.Filename)
			userImagePath := fmt.Sprintf("user/%s/%s", user.UUID, utils.NewNanoid()+Fext)
			response, err := s3.UploadImageToBucket(userImagePath, file, true, true)
			if err != nil {
				logger.Error(err.Error())
			} else {
				requestBody.AvatarUrl = response.URL
			}
		}
	}

	user, err := App.Repo.UpdateUser(*user, *requestBody)
	if err != nil {
		return nil, err
	}

	var userProfile *models.UserProfile
	userProfile, err = App.Repo.GetUserProfile(user.ID)
	if err == gorm.ErrRecordNotFound {
		App.Repo.CreateUserProfile(&models.UserProfile{UserID: user.ID})
		userProfile, err = App.Repo.GetUserProfile(user.ID)
		if err != nil {
			fmt.Println("Error in getting userProfile", err)
			return nil, err
		}
	}
	userProfile, err = App.Repo.UpdateUserProfile(*userProfile, requestBody.Bio, &requestBody.Location, &requestBody.Website, &requestBody.TwitterUrl)
	if err != nil {
		return nil, err
	}

	// Get user followers
	followService := NewFollowService()
	userFollowers := FollowUserService.FollowerCountUse(followService, *user)

	user.UserProfile = userProfile
	userData := UserProfileGetSerializerData(user, userFollowers)

	return &userData, nil
}

func (c userController) GetUserSettingsController(userID uint64) (*[]models.UserSettings, error) {
	userSettings, err := App.Repo.GetUserSettings(userID)
	if err != nil {
		return nil, err
	}
	userSocialAuth, _ := App.Repo.GetUserSocialAuth(map[string]interface{}{"user_id": userID, "provider_name": models.DISCORD_PROVIDER})

	if len(*userSettings) == 0 {
		defaultNotificationTypes := models.DefaultNotificationTypes
		if userSocialAuth != nil {
			defaultNotificationTypes = append(defaultNotificationTypes, models.NOTIFICATION_TYPE_DISCORD)
		}
		for _, notificationType := range defaultNotificationTypes {
			userSettings := models.NewDefaultUserSettings(userID, notificationType)
			done := App.Repo.CreateNewUserSettings(userSettings)
			if !done {
				return nil, errors.New("can't fetch user settings")
			}
		}
	}

	// Extra protection to ensure discord setting are present or created
	if userSocialAuth != nil {
		// Meaning user has discord auth/login.
		userDiscordSetting, _ := App.Repo.GetUserDiscordSetting(userID)
		if userDiscordSetting.ID == 0 {
			// user has discord login but no settings.
			// create the setting now.
			userDiscordSettingStruct := models.NewDefaultUserSettings(userID, models.NOTIFICATION_TYPE_DISCORD)
			_ = App.Repo.CreateNewUserSettings(userDiscordSettingStruct)
		}
	}

	// Extra protection to ensure slack setting are present or created
	slackSocialAuth, _ := App.Repo.GetUserSocialAuth(map[string]interface{}{"user_id": userID, "provider_name": models.SLACK_PROVIDER})
	if slackSocialAuth != nil {
		userSlackSetting, _ := App.Repo.GetUserSetting(userID, models.NOTIFICATION_TYPE_SLACK)
		if userSlackSetting.ID == 0 {
			userSlackSettingStruct := models.NewDefaultUserSettings(userID, models.NOTIFICATION_TYPE_SLACK)
			_ = App.Repo.CreateNewUserSettings(userSlackSettingStruct)
		}
	}

	userSettings, err = App.Repo.GetUserSettings(userID)
	if err != nil {
		return nil, err
	}

	return userSettings, nil
}

func (c userController) UpdateUserSettingsController(requestBody *PatchUserSettingsValidator, user *models.User) (*[]models.UserSettings, error) {

	for _, userSetting := range requestBody.Data {
		if userSetting.Type == "" {
			continue
		}
		isUpdated := App.Repo.UpdateUserSettings(*user, &userSetting)
		if isUpdated == 0 {
			userSettings := models.NewDefaultUserSettings(user.ID, userSetting.Type)
			App.Repo.CreateNewUserSettings(userSettings)
		}
	}
	userSettings, err := App.Repo.GetUserSettings(user.ID)
	if err != nil {
		return nil, err
	}
	return userSettings, nil
}

func (c userController) GetUserContactsController(since string, user *models.User) ([]UserContactSerializer, error) {

	contacts := []models.UserContact{}
	layout := "2006-01-02T15:04:05Z07:00"
	sinceDate, err := time.Parse(layout, since)
	if err != nil {
		contacts, err = App.Repo.GetContacts(*user)
	} else {
		contacts, err = App.Repo.GetUpdatedContactsList(sinceDate)
	}

	if err != nil {
		return nil, err
	}

	userContactViews := []UserContactSerializer{}
	for _, contact := range contacts {
		userContactViews = append(userContactViews, UserContactSerializerData(contact))
	}
	return userContactViews, nil
}

func (c userController) CreateUserContactsController(requestBody map[string]interface{}, user *models.User) error {

	contacts, isOK := requestBody["contacts"].([]interface{})
	if !isOK {
		return errors.New("Contacts Not Found")
	}

	var userContacts []models.UserContact
	for _, contact := range contacts {
		_contact := contact.(map[string]interface{})
		userId, isOk := _contact["userId"].(string)
		if isOk {
			parsedUserId, err := strconv.ParseUint(userId, 10, 64)
			if err != nil {
				logger.Error(err.Error())
				continue
			}
			userContacts = append(userContacts, user.NewUserContact(&parsedUserId, "", "", "", ""))
		} else {
			name, _ := _contact["name"].(string)
			email, _ := _contact["email"].(string)
			phone, _ := _contact["phone"].(string)
			photo, _ := _contact["photo"].(string)
			userContacts = append(userContacts, user.NewUserContact(nil, phone, email, name, photo))
		}
	}

	_ = App.Repo.CreateContacts(userContacts)
	return nil
}

func (c userController) UserSearchController() {

}

func (c userController) PopularUsersController(skip int) (*[]models.User, error) {
	cache := redis.NewCache()
	cacheKey := "popular:users"
	if skip == 0 {
		if value := cache.Get(context.Background(), cacheKey); value != nil {
			var users []models.User
			json.Unmarshal([]byte(value.(string)), &users)
			return &users, nil
		}
	}

	users, err := App.Repo.GetPopularUsers(skip)
	if err != nil {
		return nil, err
	}
	if skip == 0 {
		go func() {
			data, _ := json.Marshal(users)
			cache.Set(context.Background(), cacheKey, data, &redis.Options{
				Expiration: 24 * time.Hour,
			})
		}()
	}
	return users, nil
}

func (c userController) FollowerListController() {

}

func (c userController) UpdateFollowController() {

}

func (c userController) SetupUserController(requestBody *UpdateUserValidator, user *models.User) (*models.User, error) {
	_, err := App.Service.CreateNewUserProfile(user.ID, requestBody.Bio, &requestBody.TwitterUrl, &requestBody.Website, &requestBody.Location)
	if err != nil {
		return nil, err
	}

	// Adding IsSetupDone for the user
	requestBody.IsSetupDone = true

	if requestBody.File.Size > 0 {
		file, err := requestBody.File.Open()
		if err != nil {
			logger.Error(fmt.Sprintf("Error on reading user avatar url %s", err.Error()))
		} else {
			//			userImagePath := fmt.Sprintf("user/%s/%s", user.UUID, requestBody.File.Filename)
			Fext := filepath.Ext(requestBody.File.Filename)
			userImagePath := fmt.Sprintf("user/%s/%s", user.UUID, utils.NewNanoid()+Fext)
			response, err := s3.UploadImageToBucket(userImagePath, file, true, true)
			if err != nil {
				logger.Error(err.Error())
			} else {
				requestBody.AvatarUrl = response.URL
			}
		}
	}

	user, err = App.Repo.UpdateUser(*user, *requestBody)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (sc *userController) GetUserByHandleController(username string) (*models.User, error) {

	user, err := App.Repo.GetUserByUsername(username)
	return user, err
}

// We'll get the userID and serilize right here.
func (c userController) GetUsersFollowers(user *models.User, authUser *models.User) []UserMiniSerializer {
	followers := App.Repo.GetUserFollowersList(user)
	//following := App.Repo.GetUserFollowingList(user)
	return VarFollowersDataMaker(followers, authUser)
	//followingsSerialData := VarFollowersFollowingDataMaker(following)
	//return VarFollowersFollowing{
	//	Followers: followersSerialData,
	//	Following: followingsSerialData,
	//}
}

func (c userController) GetUsersFollowing(user *models.User, authUser *models.User) []UserMiniSerializer {
	//followers := App.Repo.GetUserFollowersList(user)
	following := App.Repo.GetUserFollowingList(user)
	//followersSerialData := VarFollowersFollowingDataMaker(followers)
	return VarFollowingDataMaker(following, authUser)
	//return VarFollowersFollowing{
	//	Followers: followersSerialData,
	//	Following: followingsSerialData,
	//}
}
