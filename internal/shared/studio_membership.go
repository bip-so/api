package shared

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

func GetStudioAdminRole(studioID uint64) (role *models.Role, err error) {
	err = postgres.GetDB().Model(models.Role{}).
		Where("studio_id = ? and name = ?", studioID, models.SYSTEM_ADMIN_ROLE).
		Preload("Members").
		Preload("Members.User").
		First(&role).Error
	return
}

// Get list of all the UserID's of Studio Admins
func GetStudioModeratorsUserIDs(studioID uint64) ([]uint64, error) {
	var studioPerms []models.StudioPermission
	err := postgres.GetDB().Model(models.StudioPermission{}).
		Where("studio_id = ? and permission_group in ?", studioID, []string{models.PGStudioAdminSysName}).
		Preload("Member").
		Preload("Role").
		Preload("Role.Members").
		Find(&studioPerms).Error
	if err != nil {
		return nil, err
	}
	var userIDs []uint64
	for _, perm := range studioPerms {
		if perm.MemberId != nil {
			userIDs = append(userIDs, perm.Member.UserID)
		} else if perm.RoleId != nil {
			for _, member := range perm.Role.Members {
				userIDs = append(userIDs, member.UserID)
			}
		}
	}
	return userIDs, err
}
