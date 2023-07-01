package shared

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"gorm.io/gorm"
)

func CheckIsUserFollowing(userID uint64, userFollowings *[]models.FollowUser) bool {

	if userFollowings == nil {
		return false
	}
	for _, userFollowing := range *userFollowings {
		if userFollowing.UserId == userID {
			return true
		}
	}
	return false
}

func IsUserStudioMember(userID, studioID uint64) bool {
	row := &models.Member{}
	err2 := postgres.GetDB().Model(models.Member{}).Where("user_id = ? and studio_id = ? and has_left = false AND is_removed = false", userID, studioID).First(&row).Error
	if err2 == gorm.ErrRecordNotFound {
		return false
	}
	if row.ID == 0 {
		return false
	}
	fmt.Println(row)

	return true
}

func IsUserStudioAdmin(userID, studioID uint64) bool {
	// Studio Admin Role Object
	userIDs, err2 := GetStudioModeratorsUserIDs(studioID)
	if err2 != nil {
		return false
	}
	return utils.SliceContainsInt(userIDs, userID)
}
