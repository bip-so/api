package user

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"gitlab.com/phonepost/bip-be-platform/pkg/logger"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/kafka"
)

type UserRepo interface {
	GetUser(query map[string]interface{}) (*models.User, error)
	GetPopularUsers(skip int) (*[]models.User, error)
	CreateUser(user *models.User) error
	CreateUserSocialAuth(user *models.UserSocialAuth) error
	CreateNewUserSettings(userSetting *models.UserSettings) bool
	GetUpdatedContactsList(since time.Time) (contacts []models.UserContact, err error)
	CreateContacts(userContacts []models.UserContact) error
	GetContacts(user models.User) (contacts []models.UserContact, err error)
	GetUserSettings(user models.User, notificationType string) (*models.UserSettings, error)
	UpdateUserSettings(user models.User, userSettings *UpdateUserSettingsValidator) error
	UpdateUser(user models.User, updates UpdateUserValidator) (*models.User, error)
	CreateUserProfile(user *models.UserProfile) error
	GetUserProfile(userID uint64) (*models.UserProfile, error)
	CheckIfUserPresentWithEmail(email string) bool
	UpdateUserProfile(userProfile models.UserProfile, bio string, location *string, website *string, twitter_url *string) (*models.UserProfile, error)
}

func (ur userRepo) GetUser(query map[string]interface{}) (*models.User, error) {
	var user models.User
	err := ur.db.Model(&models.User{}).Where(query).Preload("UserProfile").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur userRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := ur.db.Model(&models.User{}).Where("LOWER(email) = ?", strings.ToLower(email)).Preload("UserProfile").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur userRepo) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := ur.db.Model(&models.User{}).Where("LOWER(username) = ?", strings.ToLower(username)).Preload("UserProfile").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur userRepo) CheckIfUserPresentWithEmail(email string) bool {
	var user models.User
	var result int64
	err := ur.db.Model(&user).Where("email = ?", email).Count(&result).Error
	if err != nil {
		return false
	}
	if result == 1 {
		return true
	} else {
		return false
	}
}

func (ur userRepo) GetPopularUsers(skip int) (*[]models.User, error) {
	var users []models.User
	// TODO: add a logic for ordering by followers
	err := ur.db.Model(&models.User{}).Offset(skip).Limit(configs.PAGINATION_LIMIT).Preload("UserProfile").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func (ur userRepo) GetUsersByUUIDs(userUUIDs []string) (*[]models.User, error) {
	var users []models.User
	err := ur.db.Model(&models.User{}).Where("uuid IN ?", userUUIDs).Preload("UserProfile").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func (ur userRepo) GetAllUsers() (*[]models.User, error) {
	var users []models.User
	err := ur.db.Model(&models.User{}).Preload("UserProfile").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func (ur userRepo) CreateUserSocialAuth(user *models.UserSocialAuth) error {
	err := ur.db.Create(user).Error
	return err
}

func (ur userRepo) CreateNewUserSettings(userSetting *models.UserSettings) bool {
	result := ur.db.Create(userSetting)
	return result.RowsAffected == 1
}

func (ur userRepo) GetUpdatedContactsList(since time.Time) (contacts []models.UserContact, err error) {
	err = ur.db.Where("updated_at > ?", since).Preload("ContactUser").Find(&contacts).Error
	return
}

func (ur userRepo) CreateContacts(userContacts []models.UserContact) error {
	err := ur.db.Model(&userContacts).Create(userContacts).Error
	return err
}

func (ur userRepo) GetContacts(user models.User) (contacts []models.UserContact, err error) {
	err = ur.db.Where("user_id = ? AND deleted = false AND (email != ? OR contact_user_id IS NOT NULL )", user.ID, "").Preload("ContactUser").Find(&contacts).Error
	return
}

func (ur userRepo) GetUserSettings(userID uint64) (*[]models.UserSettings, error) {
	var userSettings *[]models.UserSettings
	err := ur.db.Model(&models.UserSettings{}).Where("user_id = ?", userID).Preload("User").Find(&userSettings).Error
	return userSettings, err
}

func (ur userRepo) GetUserDiscordSetting(userID uint64) (*models.UserSettings, error) {
	var userSetting *models.UserSettings
	err := ur.db.Model(&models.UserSettings{}).Where("user_id = ? and type = ?", userID, models.NOTIFICATION_TYPE_DISCORD).First(&userSetting).Error
	return userSetting, err
}

func (ur userRepo) GetUserSetting(userID uint64, integrationType string) (*models.UserSettings, error) {
	var userSetting *models.UserSettings
	err := ur.db.Model(&models.UserSettings{}).Where("user_id = ? and type = ?", userID, integrationType).First(&userSetting).Error
	return userSetting, err
}

func (ur userRepo) UpdateUserSettings(user models.User, userSettings *UpdateUserSettingsValidator) int64 {
	updates := map[string]interface{}{
		"all_comments":               userSettings.AllComments,
		"replies_to_me":              userSettings.RepliesToMe,
		"mentions":                   userSettings.Mentions,
		"reactions":                  userSettings.Reactions,
		"invite":                     userSettings.Invite,
		"followed_me":                userSettings.FollowedMe,
		"followed_my_studio":         userSettings.FollowedMyStudio,
		"publish_and_merge_requests": userSettings.PublishAndMergeRequests,
		"response_to_my_requests":    userSettings.ResponseToMyRequests,
		"system_notifications":       userSettings.SystemNotifications,
		"dark_mode":                  userSettings.DarkMode,
	}
	res := ur.db.Model(&models.UserSettings{}).Where(&models.UserSettings{UserID: user.ID, Type: userSettings.Type}).Updates(updates)
	return res.RowsAffected
}

func (ur userRepo) UpdateUser(user models.User, updates UpdateUserValidator) (*models.User, error) {
	// Assigning values to the userfield to populate the exact data in the response.
	user.Username = updates.Username
	user.FullName = updates.FullName

	// Updated fields JSON array.
	userUpdates := map[string]interface{}{
		"username":  updates.Username,
		"full_name": updates.FullName,
	}
	if updates.IsSetupDone {
		userUpdates["is_setup_done"] = updates.IsSetupDone
		user.IsSetupDone = updates.IsSetupDone
	}
	if updates.AvatarUrl != "" {
		userUpdates["avatar_url"] = updates.AvatarUrl
		user.AvatarUrl = updates.AvatarUrl
	}
	err := ur.db.Model(&models.User{}).Where("id = ?", user.ID).Updates(&userUpdates).First(&user).Error
	if err != nil {
		return nil, err
	}
	go func() {
		std, _ := ur.GetUser(map[string]interface{}{"id": user.ID})
		stdData, _ := json.Marshal(std)
		ur.kafka.Publish(configs.KAFKA_TOPICS_UPDATE_USER, strconv.FormatUint(std.ID, 10), stdData)
	}()
	return &user, nil
}

func (ur userRepo) UpdatePassword(user models.User, password string) error {
	return App.Repo.Manager.UpdateEntityByID(models.USER, user.ID, map[string]interface{}{"password": password, "has_password": true})
}

func (ur userRepo) CreateUserProfile(userProfile *models.UserProfile) error {
	err := ur.db.Create(userProfile).Error

	if err == nil {
		go func() {
			userProfileData, _ := json.Marshal(userProfile)
			kafkaClient := kafka.GetKafkaClient()
			kafkaClient.Publish(configs.KAFKA_TOPICS_NEW_USER, strconv.FormatUint(userProfile.UserID, 10), userProfileData)
		}()
	}

	return err
}

func (ur userRepo) GetUserProfile(userID uint64) (*models.UserProfile, error) {
	var userProfile *models.UserProfile
	err := ur.db.Model(&models.UserProfile{}).Where("user_id = ?", userID).First(&userProfile).Error
	if err != nil {
		return nil, err
	}
	return userProfile, nil
}

func (ur userRepo) UpdateUserProfile(userProfile models.UserProfile, bio string, location *string, website *string, twitter_url *string) (*models.UserProfile, error) {
	userProfile.Bio = bio
	userProfile.Location = location
	userProfile.Website = website
	userProfile.TwitterUrl = twitter_url
	err := ur.db.Where("user_id = ?", userProfile.UserID).Save(userProfile).Error
	if err != nil {
		return nil, err
	}
	return &userProfile, nil
}

func (fr userFollowRepo) UserCountFollowing(userId uint64) (uint64, error) {
	var userFollow []models.FollowUser
	var result int64
	App.Repo.db.Model(&userFollow).Where("follower_id = ?", userId).Count(&result)
	return uint64(result), nil
}
func (fr userFollowRepo) UserCountFollower(userId uint64) (uint64, error) {
	var userFollow []models.FollowUser
	var result int64
	App.Repo.db.Model(&userFollow).Where("user_id = ?", userId).Count(&result)
	return uint64(result), nil
}

func (ur userRepo) GetUserSocialAuth(query map[string]interface{}) (*models.UserSocialAuth, error) {
	var usa models.UserSocialAuth
	err := ur.db.Model(&models.UserSocialAuth{}).Where(query).First(&usa).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &usa, nil
}

func (ur userRepo) GetUserFollows(userID, followerID uint64) (*models.FollowUser, error) {
	var userFollow *models.FollowUser
	err := ur.db.Model(&models.FollowUser{}).Where("user_id = ? AND follower_id = ?", userID, followerID).Find(&userFollow).Error
	return userFollow, err
}

// People following this User
func (ur userRepo) GetUserFollowersList(user *models.User) *[]models.FollowUser {
	var followers []models.FollowUser
	_ = ur.db.Where(map[string]interface{}{"user_id": user.ID}).Preload("User").Preload("FollowerUser").Find(&followers).Error
	return &followers
}

// Returns followings for this USER.
func (ur userRepo) GetUserFollowingList(user *models.User) *[]models.FollowUser {
	var followings []models.FollowUser
	_ = ur.db.Where(map[string]interface{}{"follower_id": user.ID}).Preload("User").Preload("FollowerUser").Find(&followings).Error
	return &followings
}
