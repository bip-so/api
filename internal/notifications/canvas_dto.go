package notifications

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

func (r notificationRepo) GetCanvasBranchPerms(query map[string]interface{}) ([]models.CanvasBranchPermission, error) {
	var canvasBranchPerms []models.CanvasBranchPermission
	err := r.db.Model(models.CanvasBranchPermission{}).Where(query).Preload("Member").Preload("Role").Preload("Role.Members").Find(&canvasBranchPerms).Error
	if err != nil {
		return nil, err
	}
	return canvasBranchPerms, nil
}

func (r notificationRepo) GetCanvasBranchModeratorsUserIDs(branchID uint64) ([]uint64, error) {
	canvasPerms, err := r.GetCanvasBranchPerms(map[string]interface{}{"canvas_branch_id": branchID, "permission_group": models.PGCanvasModerateSysName})
	if err != nil {
		return nil, err
	}
	var userIDs []uint64
	for _, perm := range canvasPerms {
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

func (r notificationRepo) GetCanvasBranchModeratorsAndEditorsUserIDs(branchID uint64) ([]uint64, error) {
	var canvasBranchPerms []models.CanvasBranchPermission
	err := r.db.Model(models.CanvasBranchPermission{}).
		Where("canvas_branch_id = ? and permission_group in ?", branchID, []string{models.PGCanvasModerateSysName, models.PGCanvasEditSysName}).
		Preload("Member").
		Preload("Role").
		Preload("Role.Members").
		Find(&canvasBranchPerms).Error
	if err != nil {
		return nil, err
	}
	var userIDs []uint64
	for _, perm := range canvasBranchPerms {
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

func (r notificationRepo) GetStudioModeratorsUserIDs(studioID uint64) ([]uint64, error) {
	var studioPerms []models.StudioPermission
	err := r.db.Model(models.StudioPermission{}).
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

func (r notificationRepo) GetStudioBillingMemberUserIDs(studioID uint64) ([]uint64, error) {
	var role *models.Role
	err := r.db.Model(models.Role{}).
		Where("studio_id = ? and name = ?", studioID, models.SYSTEM_BILLING_ROLE).
		Preload("Members").
		First(&role).Error
	if err != nil {
		return nil, err
	}
	var userIDs []uint64
	for _, member := range role.Members {
		userIDs = append(userIDs, member.UserID)
	}
	return userIDs, err
}
