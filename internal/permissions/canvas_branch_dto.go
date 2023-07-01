package permissions

import (
	"fmt"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

func (r permissionRepo) createCanvasBranchPermission(collectionId uint64, permsgroup string, memberId *uint64, isOverridden bool, studioId uint64, canvasRepositoryId uint64, canvasBranchId uint64, parentCanvasRepoID *uint64, roleId *uint64) (*models.CanvasBranchPermission, error) {
	if parentCanvasRepoID != nil && *parentCanvasRepoID == 0 {
		parentCanvasRepoID = nil
	}
	cbp := &models.CanvasBranchPermission{
		StudioID:                    studioId,
		CollectionId:                collectionId,
		CanvasRepositoryID:          canvasRepositoryId,
		CanvasBranchID:              &canvasBranchId,
		PermissionGroup:             permsgroup,
		IsOverridden:                isOverridden,
		MemberId:                    memberId,
		RoleId:                      roleId,
		CbpParentCanvasRepositoryID: parentCanvasRepoID,
	}
	err := r.db.Create(cbp).Error
	if err != nil {
		return nil, err
	}
	return cbp, nil
}

func (r permissionRepo) getCanvasPermissionsByMemberIds(memberIds []uint64, roleIds []uint64, studioID uint64, collectionId uint64) ([]models.CanvasBranchPermission, error) {
	var canvasBranchPerms []models.CanvasBranchPermission
	err := r.db.Model(&canvasBranchPerms).Where("(member_id IN ? OR role_id IN ?) and studio_id = ? and collection_id = ? and cbp_parent_canvas_repository_id is null", memberIds, roleIds, studioID, collectionId).Order("studio_id ASC").Find(&canvasBranchPerms).Error
	return canvasBranchPerms, err
}

func (r permissionRepo) getCanvasBranchPermissionsByMemberIds(memberIds []uint64, roleIds []uint64, studioID uint64, collectionId uint64, canvasId uint64) ([]models.CanvasBranchPermission, error) {
	var canvasBranchPerms []models.CanvasBranchPermission
	err := r.db.Model(&canvasBranchPerms).Where("(member_id IN ? OR role_id IN ?) and studio_id = ? and cbp_parent_canvas_repository_id = ?", memberIds, roleIds, studioID, canvasId).Order("studio_id ASC").Find(&canvasBranchPerms).Error
	return canvasBranchPerms, err
}

func (r permissionRepo) GetCanvasBranch(query map[string]interface{}) (*models.CanvasBranch, error) {
	var canvasBranch models.CanvasBranch
	err := postgres.GetDB().Model(&models.CanvasBranch{}).Where(query).Preload("CanvasRepository").Preload("CanvasRepository.ParentCanvasRepository").First(&canvasBranch).Error
	if err != nil {
		return nil, err
	}
	return &canvasBranch, nil
}

func (r permissionRepo) GetCanvasPermissionsByID(canvasBranchID uint64) ([]models.CanvasBranchPermission, error) {
	var canvasBranchPerms []models.CanvasBranchPermission
	err := r.db.Model(&canvasBranchPerms).Where("canvas_branch_id = ?", canvasBranchID).Preload("Member").Find(&canvasBranchPerms).Error
	return canvasBranchPerms, err
}

func (r permissionRepo) UpdateCanvasBranchPermission(query, updates map[string]interface{}) error {
	err := postgres.GetDB().Model(&models.CanvasBranchPermission{}).Where(query).Updates(updates).Error
	if err != nil {
		return err
	}
	return nil
}

func (r permissionRepo) GetCanvasBranchPerms(query map[string]interface{}) ([]models.CanvasBranchPermission, error) {
	var canvasBranchPerms []models.CanvasBranchPermission
	err := r.db.Model(models.CanvasBranchPermission{}).Where(query).Preload("Member").Preload("Role").Preload("Role.Members").Find(&canvasBranchPerms).Error
	if err != nil {
		return nil, err
	}
	return canvasBranchPerms, nil
}

func (r permissionRepo) GetCanvasBranchPermission(query map[string]interface{}) (*models.CanvasBranchPermission, error) {
	var canvasBranchPermission models.CanvasBranchPermission
	err := postgres.GetDB().Model(&models.CanvasBranchPermission{}).Where(query).First(&canvasBranchPermission).Error
	if err != nil {
		return nil, err
	}
	return &canvasBranchPermission, nil
}

func (r permissionRepo) createCanvasPermissionByMemberID(collectionId uint64, studioId uint64, memberId *uint64, isOverridden bool, permissionGroup string, authUserId uint64, branchId uint64, canvasRepositoryID uint64, parentCanvasRepoID *uint64) (*models.CanvasBranchPermission, error) {
	canvasBranchPerms := &models.CanvasBranchPermission{}
	canvasBranchPerms.MemberId = memberId
	canvasBranchPerms.StudioID = studioId
	canvasBranchPerms.CollectionId = collectionId
	canvasBranchPerms.CanvasBranchID = &branchId
	canvasBranchPerms.CanvasRepositoryID = canvasRepositoryID
	canvasBranchPerms.PermissionGroup = permissionGroup
	canvasBranchPerms.IsOverridden = isOverridden
	canvasBranchPerms.CbpParentCanvasRepositoryID = parentCanvasRepoID

	err := r.db.Create(&canvasBranchPerms).Error
	if err != nil {
		return nil, err
	}

	extraData := notifications.NotificationExtraData{
		CanvasRepoID:   canvasBranchPerms.CanvasRepositoryID,
		CanvasBranchID: *canvasBranchPerms.CanvasBranchID,
	}
	contentObject := models.CANVAS_BRANCH
	member, _ := App.Repo.GetMember(map[string]interface{}{"id": memberId})
	notifications.App.Service.PublishNewNotification(notifications.CanvasInviteByName, authUserId, []uint64{member.UserID}, &canvasBranchPerms.StudioID,
		nil, extraData, canvasBranchPerms.CanvasBranchID, &contentObject)
	return canvasBranchPerms, nil
}

func (r permissionRepo) GetCanvasRepo(query map[string]interface{}) (*models.CanvasRepository, error) {
	var repo models.CanvasRepository
	err := r.db.Model(&models.CanvasRepository{}).Where(query).First(&repo).Error
	if err != nil {
		fmt.Println("Error in getting canvas repository", err)
		return nil, err
	}
	return &repo, nil
}

func (r permissionRepo) GetCanvasRepos(query map[string]interface{}) (*[]models.CanvasRepository, error) {
	var repos *[]models.CanvasRepository
	err := r.db.Model(&models.CanvasRepository{}).Where(query).Find(&repos).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return repos, nil
}

func (r permissionRepo) GetCanvasBranches(query map[string]interface{}) (*[]models.CanvasBranch, error) {
	var branches *[]models.CanvasBranch
	err := r.db.Model(&models.CanvasBranch{}).Where(query).Find(&branches).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return branches, nil
}

func (r permissionRepo) UpdateCanvasBranchPermissions(query map[string]interface{}, collectionId uint64, studioId uint64, memberId uint64, roleId uint64, isOverridden bool, permissionGroup string, authUserId uint64, branchId uint64, canvasRepositoryID uint64, parentCanvasRepoID *uint64) (*models.CanvasBranchPermission, error) {
	var canvasBranchPerms *models.CanvasBranchPermission

	err := postgres.GetDB().Model(&canvasBranchPerms).Where(query).First(&canvasBranchPerms).Error

	// create flow
	if err != nil {
		if roleId == 0 {
			canvasBranchPerms.RoleId = nil
		} else {
			canvasBranchPerms.RoleId = &roleId
		}

		if memberId == 0 {
			canvasBranchPerms.MemberId = nil
		} else {
			canvasBranchPerms.MemberId = &memberId
		}
		canvasBranchPerms.StudioID = studioId
		canvasBranchPerms.CollectionId = collectionId

		canvasBranchPerms.CanvasBranchID = &branchId
		canvasBranchPerms.CanvasRepositoryID = canvasRepositoryID
		canvasBranchPerms.CbpParentCanvasRepositoryID = parentCanvasRepoID
		canvasBranchPerms.PermissionGroup = permissionGroup
	}

	if err == nil {
		if roleId != 0 {
			var role *models.Role
			r.db.Model(&role).Where("id = ?", roleId).First(&role)
			if !(role.Name == models.SYSTEM_ADMIN_ROLE) {
				canvasBranchPerms.PermissionGroup = permissionGroup
			}
		} else {
			canvasBranchPerms.PermissionGroup = permissionGroup
		}
	}

	// we are always updating these field if record if found or not found
	canvasBranchPerms.IsOverridden = isOverridden
	canvasBranchPerms.CbpParentCanvasRepositoryID = parentCanvasRepoID

	// save will create a new record if it doesn't finds or update the existing record
	err2 := postgres.GetDB().Save(&canvasBranchPerms).Error

	if err2 != nil {
		return nil, err2
	}
	if err != nil {
		// new record created we need to send preloaded data
		//err = postgres.GetDB().Model(&canvasBranchPerms).Where(query).Preload("Member.User").First(&canvasBranchPerms).Error

		// Add a notification of collection invite
		// For auto adding permission notification is not needed to send so commenting here.
		//extraData := notifications.NotificationExtraData{
		//	CanvasRepoID:   canvasBranchPerms.CanvasRepositoryID,
		//	CanvasBranchID: *canvasBranchPerms.CanvasBranchID,
		//}
		//contentObject := models.CANVAS_BRANCH
		//if canvasBranchPerms.RoleId == nil {
		//	notifications.App.Service.PublishNewNotification(notifications.CanvasInviteByName, authUserId, []uint64{canvasBranchPerms.Member.UserID}, &canvasBranchPerms.StudioID,
		//		nil, extraData, canvasBranchPerms.CanvasBranchID, &contentObject)
		//} else {
		//	notifications.App.Service.PublishNewNotification(notifications.CanvasInviteByGroup, authUserId, nil, &canvasBranchPerms.StudioID,
		//		&[]uint64{*canvasBranchPerms.RoleId}, extraData, canvasBranchPerms.CanvasBranchID, &contentObject)
		//}
		//return canvasBranchPerms, err
	}

	return canvasBranchPerms, nil
}

func (r permissionRepo) GetCanvasBranchWithPreload(query map[string]interface{}) (*models.CanvasBranch, error) {
	var canvasBranch models.CanvasBranch
	err := postgres.GetDB().Model(&models.CanvasBranch{}).Where(query).
		Preload("CanvasRepository").
		Preload("CanvasRepository.Collection").
		Preload("CanvasRepository.ParentCanvasRepository").
		Preload("CanvasRepository.ParentCanvasRepository.DefaultBranch").
		First(&canvasBranch).Error
	if err != nil {
		return nil, err
	}
	return &canvasBranch, nil
}
