package auth

import (
	"encoding/json"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"strconv"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	user2 "gitlab.com/phonepost/bip-be-platform/internal/user"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

// Auth describes the service.
func (s *authService) getUserWithUserID(userID uint64) (*models.User, error) {
	var user *models.User
	var err error
	//user, err = user2.App.Repo.GetUser(map[string]interface{}{"id": userID})
	user, err = queries.App.UserQueries.GetUser(map[string]interface{}{"id": userID})

	if err != nil {
		return nil, err
	}

	return user, nil
}
func (s *authService) getUserWithUserName(username string) (*models.User, error) {
	var user *models.User
	var err error
	user, err = user2.App.Repo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) getUserWithEmail(email string) (*models.User, error) {
	var user *models.User
	var err error
	user, err = user2.App.Repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (s *authService) isValidProvider(socialProvider string) bool {
	if !utils.SliceContainsItem([]string{models.TWITTER_PROVIDER, models.DISCORD_PROVIDER, models.SLACK_PROVIDER}, socialProvider) {
		return false
	}
	return true
}

func (s *authService) IsValidOtp(email string, otp string) bool {
	// todo
	// Get user with the email
	user, _ := s.getUserWithEmail(email)
	// Check cache with the UserId string
	UserIdStr := strconv.FormatUint(user.ID, 10)
	val := getOtpDataFromCache(UserIdStr)
	if val == "" {
		return false
	}
	// If userkey is present "Check Token"
	data := ResponseNewOtp{}
	_ = json.Unmarshal([]byte(val), &data)
	// if they match exit True
	// else exit false.
	if data.Otp != otp {
		return false
	}
	return true
}

func (s *authService) CheckIsNewUser(user *models.User) bool {
	user, _ = queries.App.UserQueries.GetUser(map[string]interface{}{"id": user.ID})
	fmt.Println("UserId and user default studioId", user.ID, user.DefaultStudioID)
	studioPerms, err := permissions.App.Service.CalculateStudioPermissions(user.ID)
	if err != nil {
		fmt.Println("Error in getting studio perms", err)
		return false
	}
	// Get all user studios
	// if user has mod perm on one studio then return false
	// else if user doesn't have mod perm on any of his studios then return true
	studios, _ := notifications.App.Repo.GetUserStudiosByID(user.ID)
	for _, studio := range studios {
		if studio.ID == user.DefaultStudioID {
			continue
		}
		if studioPerms[studio.ID] == models.PGStudioAdminSysName {
			return false
		}
	}
	return true
}
