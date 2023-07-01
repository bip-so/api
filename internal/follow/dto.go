package follow

import (
	"fmt"

	"gitlab.com/phonepost/bip-be-platform/internal/feed"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

type FollowRepo interface {
	UserCountFollowing(userId uint64) (uint64, error) // count following
	UserCountFollower(userId uint64) (uint64, error)  // count followers
	UserFollowUser(user models.User, userId uint64)   // follow u1, u2
	UserFollowUnfollowUser(user models.User, userId uint64)
	StudioCountFollowing(studioId uint64) (uint64, error)
	GetIsUserFollowingUser(userID uint64, followingUserIDs []uint64) (*[]models.FollowUser, error)
	GetIsUserFollowingStudios(userID uint64, followingStudioIDs []uint64) (*[]models.FollowStudio, error)
	UserFollowStudio(studioId uint64, userId uint64)
	UserUnFollowStudio(studioId uint64, userId uint64)
}

// Count of this user followings
func (fr followRepo) UserCountFollowing(userId uint64) (uint64, error) {
	var userFollow []models.FollowUser
	var result int64
	fr.db.Model(&userFollow).Where("user_id = ?", userId).Count(&result)
	return uint64(result), nil
}

// Count of users following this user.
func (fr followRepo) UserCountFollower(userId uint64) (uint64, error) {
	var userFollow []models.FollowUser
	var result int64
	fr.db.Model(&userFollow).Where("follower_id = ?", userId).Count(&result)
	return uint64(result), nil
}

// Logged in user follows another user
// UserId -> Other user
// FollowerId -> Logged user
func (fr followRepo) UserFollowUser(user models.User, userId uint64) {
	var userFollow models.FollowUser
	var count int64
	userFollow.FollowerId = user.ID
	userFollow.UserId = userId

	fr.db.Where("user_id = ? AND follower_id = ?", userId, user.ID).Find(&userFollow).Count(&count)
	fmt.Printf("Result : %v \n", count)
	if count == 0 {
		fmt.Printf("New Entry")
		fr.db.Model(&userFollow).Create(&userFollow)

	}
	deleteCahceFollowerCount(user.ID)
	deleteCahceFollowerCount(userId)

	// @todo later move to kafka
	go feed.App.Service.FollowUser(&userFollow)
	go func() {
		notifications.App.Service.PublishNewNotification(notifications.FollowUser, user.ID, []uint64{userId},
			nil, nil, notifications.NotificationExtraData{}, nil, nil)
	}()
}

// Unfollow a user
// Logged in user unfollows another user
// UserId -> Other user
// FollowerId -> Logged user
func (fr followRepo) UserFollowUnfollowUser(user models.User, userId uint64) {

	var userFollow models.FollowUser

	userFollow.FollowerId = user.ID
	userFollow.UserId = userId

	fr.db.Where("user_id = ? AND follower_id = ?", userId, user.ID).Delete(&userFollow)

	deleteCahceFollowerCount(user.ID)
	deleteCahceFollowerCount(userId)

	go feed.App.Service.UnfollowUser(user.ID, userId)
}

func (fr followRepo) StudioCountFollowing(studioId uint64) (uint64, error) {
	var studioFollow []models.FollowStudio
	var result int64
	fr.db.Model(&studioFollow).Where("studio_id = ?", studioId).Count(&result)
	return uint64(result), nil
}

// Follow a Studio
func (fr followRepo) UserFollowStudio(studioId uint64, userId uint64) {

	var studioFollow models.FollowStudio
	var count int64
	studioFollow.StudioId = studioId
	studioFollow.FollowerId = userId

	fr.db.Where("studio_id = ? AND follower_id = ?", studioId, userId).Find(&studioFollow).Count(&count)
	if count == 0 {
		fmt.Printf("New Entry")
		fr.db.Model(&studioFollow).Create(&studioFollow)
	}
	deleteCahceStudioFollowerCount(studioId)
}

func (fr followRepo) UserUnFollowStudio(studioId uint64, userId uint64) {
	var studioFollow models.FollowStudio
	studioFollow.StudioId = studioId
	studioFollow.FollowerId = userId

	fr.db.Where("studio_id = ? AND follower_id = ?", studioId, userId).Delete(&studioFollow)
	deleteCahceStudioFollowerCount(studioId)
}

func (fr followRepo) GetIsUserFollowingUser(userID uint64, followingUserIDs []uint64) (*[]models.FollowUser, error) {
	var userFollow []models.FollowUser
	err := fr.db.Model(&models.FollowUser{}).Where("user_id IN ? AND follower_id = ?", followingUserIDs, userID).Find(&userFollow).Error
	return &userFollow, err
}

func (fr followRepo) GetIsUserFollowingStudios(userID uint64, followingStudioIDs []uint64) (*[]models.FollowStudio, error) {
	var studioFollow []models.FollowStudio
	err := fr.db.Where("studio_id IN ? AND follower_id = ?", followingStudioIDs, userID).Find(&studioFollow).Error
	return &studioFollow, err
}

func (fr followRepo) GetFollowUser(query map[string]interface{}) (*models.FollowUser, error) {
	var userFollow *models.FollowUser
	err := fr.db.Where(query).First(&userFollow).Error
	return userFollow, err
}

// GetUserFollowing
func (fr followRepo) GetFollowUsers(query map[string]interface{}) ([]models.FollowUser, error) {
	var userFollows []models.FollowUser
	err := fr.db.Where(query).Preload("User").Preload("FollowerUser").First(&userFollows).Error
	return userFollows, err
}
