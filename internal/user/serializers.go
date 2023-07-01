package user

import (
	"database/sql"
	"fmt"
	"time"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

type UserGetSerializer struct {
	Id              uint64                `json:"id"`
	UUID            string                `json:"uuid"`
	FullName        string                `json:"fullName"`
	Username        string                `json:"username"`
	HasEmail        bool                  `json:"hasEmail"`
	IsSuperuser     bool                  `json:"isSuperuser"`
	IsSetupDone     bool                  `json:"isSetupDone"`
	IsEmailVerified bool                  `json:"isEmailVerified"`
	AvatarUrl       string                `json:"avatarUrl"`
	Followers       uint64                `json:"followers"`
	Following       uint64                `json:"following"`
	IsFollowing     *bool                 `json:"isFollowing"`
	UserProfile     UserProfileSerializer `json:"userProfile"`
	CreatedAt       time.Time             `json:"createdAt"`
	UpdatedAt       time.Time             `json:"updatedAt"`
	HasPassword     bool                  `json:"hasPassword"`
	DefaultStudioID uint64                `json:"defaultStudioID"`
}

type UserProfileSerializer struct {
	UserID     uint64  `json:"userId"`
	Bio        string  `json:"bio"`
	Website    *string `json:"website"`
	TwitterUrl *string `json:"twitterUrl"`
	Location   *string `json:"location"`
}

func SerializeUserProfileData(userProfile *models.UserProfile) UserProfileSerializer {
	return UserProfileSerializer{
		UserID:     userProfile.UserID,
		Bio:        userProfile.Bio,
		Website:    userProfile.Website,
		TwitterUrl: userProfile.TwitterUrl,
		Location:   userProfile.Location,
	}
}

func UserGetSerializerData(user *models.User, followUsers *[]models.FollowUser) UserGetSerializer {
	isFollowing := checkIsFollowing(user, followUsers)
	userData := UserGetSerializer{
		Id:              user.ID,
		UUID:            user.UUID.String(),
		Username:        user.Username,
		FullName:        user.FullName,
		HasEmail:        user.Email.Valid,
		AvatarUrl:       user.AvatarUrl,
		IsSuperuser:     user.IsSuperuser,
		IsSetupDone:     user.IsSetupDone,
		IsEmailVerified: user.IsEmailVerified,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
		IsFollowing:     isFollowing,
		HasPassword:     user.HasPassword,
		DefaultStudioID: user.DefaultStudioID,
	}
	if user.UserProfile != nil {
		userData.UserProfile = SerializeUserProfileData(user.UserProfile)
	}
	return userData
}

func UserProfileGetSerializerData(user *models.User, userFollowers FollowUserFollowCountResponse) UserGetSerializer {
	userData := UserGetSerializer{
		Id:              user.ID,
		UUID:            user.UUID.String(),
		Username:        user.Username,
		FullName:        user.FullName,
		HasEmail:        user.Email.Valid,
		AvatarUrl:       user.AvatarUrl,
		IsSuperuser:     user.IsSuperuser,
		IsSetupDone:     user.IsSetupDone,
		IsEmailVerified: user.IsEmailVerified,
		Followers:       userFollowers.Followers,
		Following:       userFollowers.Following,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
		HasPassword:     user.HasPassword,
		DefaultStudioID: user.DefaultStudioID,
	}
	if user.UserProfile != nil {
		userData.UserProfile = SerializeUserProfileData(user.UserProfile)
	}
	return userData
}

type UserSettingSerializer struct {
	Type                    string `json:"type"`
	ID                      uint64 `json:"id"`
	UserID                  uint64 `json:"userId"`
	AllComments             bool   `json:"allComments"`
	RepliesToMe             bool   `json:"repliesToMe"`
	Mentions                bool   `json:"mentions"`
	Reactions               bool   `json:"reactions"`
	Invite                  bool   `json:"invite"`
	FollowedMe              bool   `json:"followedMe"`
	FollowedMyStudio        bool   `json:"followedMyStudio"`
	PublishAndMergeRequests bool   `json:"publishAndMergeRequests"`
	ResponseToMyRequests    bool   `json:"responseToMyRequests"`
	SystemNotifications     bool   `json:"systemNotifications"`
	DarkMode                bool   `json:"darkMode"`
}

type CustomUserSettingsResponse struct {
	App     UserSettingSerializer `json:"app"`
	Discord UserSettingSerializer `json:"discord"`
	Email   UserSettingSerializer `json:"email"`
	Slack   UserSettingSerializer `json:"slack"`
}

func CustomUserSettingsSerializerData(setting models.UserSettings) UserSettingSerializer {
	return UserSettingSerializer{
		Type:                    setting.Type,
		ID:                      setting.ID,
		UserID:                  setting.UserID,
		AllComments:             setting.AllComments,
		RepliesToMe:             setting.RepliesToMe,
		Mentions:                setting.Mentions,
		Reactions:               setting.Reactions,
		Invite:                  setting.Invite,
		FollowedMe:              setting.FollowedMe,
		FollowedMyStudio:        setting.FollowedMyStudio,
		PublishAndMergeRequests: setting.PublishAndMergeRequests,
		ResponseToMyRequests:    setting.ResponseToMyRequests,
		SystemNotifications:     setting.SystemNotifications,
		DarkMode:                setting.DarkMode,
	}
}

func MultiUserSettingsSerializerData(userSettings []models.UserSettings) []UserSettingSerializer {
	userSettingsData := &[]UserSettingSerializer{}
	for _, setting := range userSettings {
		data := &UserSettingSerializer{
			Type:                    setting.Type,
			ID:                      setting.ID,
			UserID:                  setting.UserID,
			AllComments:             setting.AllComments,
			RepliesToMe:             setting.RepliesToMe,
			Mentions:                setting.Mentions,
			Reactions:               setting.Reactions,
			Invite:                  setting.Invite,
			FollowedMe:              setting.FollowedMe,
			FollowedMyStudio:        setting.FollowedMyStudio,
			PublishAndMergeRequests: setting.PublishAndMergeRequests,
			ResponseToMyRequests:    setting.ResponseToMyRequests,
			SystemNotifications:     setting.SystemNotifications,
			DarkMode:                setting.DarkMode,
		}
		*userSettingsData = append(*userSettingsData, *data)
	}
	return *userSettingsData
}

type UserContactSerializer struct {
	ID          uint64            `json:"id"`
	ContactUser UserGetSerializer `json:"user"`
	Email       string            `json:"email"`
	Deleted     bool              `json:"deleted"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

func UserContactSerializerData(contactUser models.UserContact) UserContactSerializer {
	var contactUserView UserGetSerializer
	email := contactUser.Email
	if contactUser.ContactUserID != nil {
		contactUserView = UserGetSerializerData(contactUser.ContactUser, nil)
		email = contactUser.ContactUser.Email.String
	} else {
		contactUserView = UserGetSerializerData(&models.User{
			Username:  contactUser.Name,
			Email:     sql.NullString{String: contactUser.Email, Valid: true},
			AvatarUrl: contactUser.Photo,
		}, nil)
	}
	return UserContactSerializer{
		ID:          contactUser.ID,
		ContactUser: contactUserView,
		Email:       email,
		Deleted:     contactUser.Deleted,
		CreatedAt:   contactUser.CreatedAt,
		UpdatedAt:   contactUser.UpdatedAt,
	}
}

type PaginatedResponse struct {
	Data interface{}
	Next int
}

func PaginatedResponseData(data interface{}, next int) PaginatedResponse {
	return PaginatedResponse{
		Data: data,
		Next: next,
	}
}

func checkIsFollowing(user *models.User, followUsers *[]models.FollowUser) *bool {
	if followUsers == nil {
		return nil
	}
	following := false
	for _, followUser := range *followUsers {
		if followUser.UserId == user.ID {
			following = true
			return &following
		}
	}
	return &following
}

type UserMiniSerializer struct {
	Id          uint64 `json:"id"`
	UUID        string `json:"uuid"`
	FullName    string `json:"fullName"`
	Username    string `json:"username"`
	AvatarUrl   string `json:"avatarUrl"`
	Followers   uint64 `json:"followers"`
	Following   uint64 `json:"following"`
	IsFollowing bool   `json:"isFollowing"`
}

func UserMiniSerializerData(user *models.User) UserMiniSerializer {
	userData := UserMiniSerializer{
		Id:        user.ID,
		UUID:      user.UUID.String(),
		Username:  user.Username,
		FullName:  user.FullName,
		AvatarUrl: user.AvatarUrl,
	}
	return userData
}

func VarFollowersDataMaker(users *[]models.FollowUser, authUser *models.User) []UserMiniSerializer {
	serilizedUserList := &[]UserMiniSerializer{}
	fmt.Println("Userssssss", users)
	if users == nil {
		return *serilizedUserList
	}
	authUserFollowers := App.Repo.GetUserFollowingList(authUser)
	followMap := map[uint64]*models.FollowUser{}
	for _, fUser := range *authUserFollowers {
		followMap[fUser.UserId] = &fUser
	}
	for _, fUser := range *users {
		user := fUser.FollowerUser
		data := &UserMiniSerializer{
			Id:        user.ID,
			UUID:      user.UUID.String(),
			Username:  user.Username,
			FullName:  user.FullName,
			AvatarUrl: user.AvatarUrl,
		}
		if followUser, exists := followMap[user.ID]; exists {
			data.IsFollowing = followUser.FollowerId == authUser.ID
		}
		followingCount, _ := App.FollowRepo.UserCountFollowing(user.ID)
		followerCount, _ := App.FollowRepo.UserCountFollower(user.ID)
		data.Following = followingCount
		data.Followers = followerCount
		*serilizedUserList = append(*serilizedUserList, *data)
	}
	return *serilizedUserList
}

func VarFollowingDataMaker(users *[]models.FollowUser, authUser *models.User) []UserMiniSerializer {
	serilizedUserList := &[]UserMiniSerializer{}
	if users == nil {
		return *serilizedUserList
	}
	authUserFollowers := App.Repo.GetUserFollowingList(authUser)
	followMap := map[uint64]*models.FollowUser{}
	for _, fUser := range *authUserFollowers {
		followMap[fUser.UserId] = &fUser
	}
	for _, fUser := range *users {
		user := fUser.User
		data := &UserMiniSerializer{
			Id:        user.ID,
			UUID:      user.UUID.String(),
			Username:  user.Username,
			FullName:  user.FullName,
			AvatarUrl: user.AvatarUrl,
		}
		if followUser, exists := followMap[user.ID]; exists {
			data.IsFollowing = followUser.FollowerId == authUser.ID
		}
		followingCount, _ := App.FollowRepo.UserCountFollowing(user.ID)
		followerCount, _ := App.FollowRepo.UserCountFollower(user.ID)
		data.Following = followingCount
		data.Followers = followerCount
		*serilizedUserList = append(*serilizedUserList, *data)
	}
	return *serilizedUserList
}
