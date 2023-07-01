package collectionpermissions

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

// Interface
type CollectionPermissionRepo interface {
	CreateCollectionPermission(role *models.CollectionPermission) (*models.CollectionPermission, error)
	GetCollectionPermission(query map[string]interface{}) ([]models.CollectionPermission, error)
	UpdateCollectionPermissions(query map[string]interface{}, body CollectionPermissionValidator, studioID uint64, authUserId uint64) (*models.CollectionPermission, error)
	GetCollectionPermissionsForUser(memberIds []uint64, roleIds []uint64) ([]models.CollectionPermission, error)
}

// constructor
func NewCollectionPermissionsRepo() CollectionPermissionRepo {
	return &collectionPermissionRepo{}
}

func (sr collectionPermissionRepo) Get(query map[string]interface{}) (models.CollectionPermission, error) {
	var collectionPerm models.CollectionPermission
	postgres.GetDB().Model(&collectionPerm).Where(query).Find(&collectionPerm)
	return collectionPerm, nil
}

func (sr collectionPermissionRepo) GetCollectionPermission(query map[string]interface{}) ([]models.CollectionPermission, error) {
	var collectionPerms []models.CollectionPermission
	postgres.GetDB().Model(&collectionPerms).Where(query).Preload("Studio").Preload("Role").Preload("Collection").Preload("Member").Preload("Member.User").Preload("Role.Members").Preload("Collection").Find(&collectionPerms)
	return collectionPerms, nil
}

func (sr collectionPermissionRepo) CreateCollectionPermission(stdperms *models.CollectionPermission) (*models.CollectionPermission, error) {
	result := postgres.GetDB().Create(stdperms)
	if result.Error != nil {
		return nil, result.Error
	}
	return stdperms, nil
}

func (cr collectionPermissionRepo) UpdateCollectionPermissions(query map[string]interface{}, body CollectionPermissionValidator, studioID uint64, authUserId uint64) (*models.CollectionPermission, error) {
	var collectionPerms *models.CollectionPermission

	err := postgres.GetDB().Model(&collectionPerms).Where(query).Preload("Role").Preload("Collection").Preload("Member.User").Preload("Role.Members").Preload("Member").Preload("Studio").First(&collectionPerms).Error

	// create flow
	if err != nil {
		if body.RoleID == 0 {
			collectionPerms.RoleId = nil
		} else {
			collectionPerms.RoleId = &body.RoleID
		}

		if body.MemberID == 0 {
			collectionPerms.MemberId = nil
		} else {
			collectionPerms.MemberId = &body.MemberID
		}
		collectionPerms.StudioID = studioID
		collectionPerms.CollectionId = body.CollectionId
		collectionPerms.PermissionGroup = body.PermGroup
	}

	if err == nil {
		if body.RoleID != 0 {
			var role *models.Role
			postgres.GetDB().Model(&role).Where("id = ?", body.RoleID).First(&role)
			if !(role.Name == models.SYSTEM_ADMIN_ROLE) {
				collectionPerms.PermissionGroup = body.PermGroup
			}
		} else {
			collectionPerms.PermissionGroup = body.PermGroup
		}
	}

	// we are always updating these field if record if found or not found
	collectionPerms.IsOverridden = body.IsOverridden

	// save will create a new record if it doesn't finds or update the existing record
	err2 := postgres.GetDB().Save(&collectionPerms).Error

	if err2 != nil {
		return nil, err2
	}
	if err != nil {
		// new record created we need to send preloaded data
		err = postgres.GetDB().Model(&collectionPerms).Where(query).Preload("Role").Preload("Studio").Preload("Collection").Preload("Member.User").Preload("Role.Members").Preload("Member").First(&collectionPerms).Error
		// Add a notification of collection invite
		go func() {
			extraData := notifications.NotificationExtraData{
				CollectionID: collectionPerms.CollectionId,
			}
			contentObject := models.COLLECTION
			if collectionPerms.RoleId == nil {
				notifications.App.Service.PublishNewNotification(notifications.CollectionInviteByName, authUserId, []uint64{collectionPerms.Member.UserID}, &collectionPerms.StudioID,
					nil, extraData, &collectionPerms.CollectionId, &contentObject)
			} else {
				notifications.App.Service.PublishNewNotification(notifications.CollectionInviteByGroup, authUserId, nil, &collectionPerms.StudioID,
					&[]uint64{*collectionPerms.RoleId}, extraData, &collectionPerms.CollectionId, &contentObject)
			}
		}()
		return collectionPerms, err
	}

	return collectionPerms, nil
}

func (cr collectionPermissionRepo) GetCollectionPermissionsForUser(memberIds, roleIds []uint64) ([]models.CollectionPermission, error) {
	var collectionPerms []models.CollectionPermission
	err := postgres.GetDB().Model(&collectionPerms).Where("member_id IN ? OR role_id IN ?", memberIds, roleIds).Preload("Studio").Preload("Role").Preload("Collection").Preload("Member").Preload("Member.User").Preload("Role.Members").Preload("Collection").Find(&collectionPerms).Error
	return collectionPerms, err
}

func (sr collectionPermissionRepo) GetCanvasRepos(query map[string]interface{}) ([]models.CanvasRepository, error) {
	var canvasRepos []models.CanvasRepository
	err := postgres.GetDB().Model(&canvasRepos).Where(query).Find(&canvasRepos).Error
	return canvasRepos, err
}

// canvas permissions
func (cr collectionPermissionRepo) UpdateCanvasBranchPermissions(query map[string]interface{}, body newCanvasBranchPermissionCreatePost, studioID uint64, authUserId uint64) (*models.CanvasBranchPermission, error) {
	var canvasBranchPerms *models.CanvasBranchPermission

	err := postgres.GetDB().Model(&canvasBranchPerms).Where(query).Preload("Studio").Preload("Role").Preload("CanvasBranch").Preload("Member").Preload("Member.User").Preload("Role.Members").Preload("CanvasRepository").Preload("CbpParentCanvasRepository").Preload("Collection").Preload("CanvasRepository.DefaultBranch").Preload("CbpParentCanvasRepository.DefaultBranch").First(&canvasBranchPerms).Error

	// create flow
	if err != nil {
		if body.RoleID == 0 {
			canvasBranchPerms.RoleId = nil
		} else {
			canvasBranchPerms.RoleId = &body.RoleID
		}

		if body.MemberID == 0 {
			canvasBranchPerms.MemberId = nil
		} else {
			canvasBranchPerms.MemberId = &body.MemberID
		}
		canvasBranchPerms.StudioID = studioID
		canvasBranchPerms.CollectionId = body.CollectionId

		canvasBranchPerms.CanvasBranchID = &body.CanvasBranchId
		canvasBranchPerms.CanvasRepositoryID = body.CanvasRepositoryID
		canvasBranchPerms.CbpParentCanvasRepositoryID = &body.CbpParentCanvasRepositoryID
		canvasBranchPerms.PermissionGroup = body.PermGroup
	}

	if err == nil {
		if body.RoleID != 0 {
			var role *models.Role
			postgres.GetDB().Model(&role).Where("id = ?", body.RoleID).First(&role)
			if !(role.Name == models.SYSTEM_ADMIN_ROLE || role.Name == models.SYSTEM_ROLE_MEMBER) {
				canvasBranchPerms.PermissionGroup = body.PermGroup
			}
		} else {
			canvasBranchPerms.PermissionGroup = body.PermGroup
		}
	}

	// we are always updating these field if record if found or not found
	canvasBranchPerms.IsOverridden = body.IsOverridden
	if body.CbpParentCanvasRepositoryID == 0 {
		canvasBranchPerms.CbpParentCanvasRepositoryID = nil
	} else {
		canvasBranchPerms.CbpParentCanvasRepositoryID = &body.CbpParentCanvasRepositoryID
	}

	// save will create a new record if it doesn't finds or update the existing record
	err2 := postgres.GetDB().Save(&canvasBranchPerms).Error

	if err2 != nil {
		return nil, err2
	}
	if err != nil {
		// new record created we need to send preloaded data
		err = postgres.GetDB().Model(&canvasBranchPerms).Where(query).Preload("Studio").Preload("Role").Preload("CanvasBranch").Preload("Member").Preload("Member.User").Preload("Role.Members").Preload("CanvasRepository").Preload("CbpParentCanvasRepository").Preload("Collection").Preload("CanvasRepository.DefaultBranch").Preload("CbpParentCanvasRepository.DefaultBranch").First(&canvasBranchPerms).Error

		// Add a notification of collection invite
		// go func() {
		// 	extraData := notifications.NotificationExtraData{
		// 		CanvasRepoID:   canvasBranchPerms.CanvasRepositoryID,
		// 		CanvasBranchID: *canvasBranchPerms.CanvasBranchID,
		// 	}
		// 	contentObject := models.CANVAS_BRANCH
		// 	if canvasBranchPerms.RoleId == nil {
		// 		fmt.Println("canvasbranch perms", canvasBranchPerms)
		// 		notifications.App.Service.PublishNewNotification(notifications.CanvasInviteByName, authUserId, []uint64{canvasBranchPerms.Member.UserID}, &canvasBranchPerms.StudioID,
		// 			nil, extraData, canvasBranchPerms.CanvasBranchID, &contentObject)
		// 	} else {
		// 		notifications.App.Service.PublishNewNotification(notifications.CanvasInviteByGroup, authUserId, nil, &canvasBranchPerms.StudioID,
		// 			&[]uint64{*canvasBranchPerms.RoleId}, extraData, canvasBranchPerms.CanvasBranchID, &contentObject)
		// 	}
		// }()
		return canvasBranchPerms, err
	}

	return canvasBranchPerms, nil
}

func (cr collectionPermissionRepo) GetMember(query map[string]interface{}) (*models.Member, error) {
	var member models.Member
	err := postgres.GetDB().Model(&models.Member{}).Where(query).First(&member).Error
	if err != nil {
		//log.Fatalln(err)
		return nil, err
	}
	return &member, nil
}

func (cr collectionPermissionRepo) GetCanvasBranchPermissions(query map[string]interface{}) ([]models.CanvasBranchPermission, error) {
	var canvasBranchPerms []models.CanvasBranchPermission
	err := postgres.GetDB().Model(models.CanvasBranchPermission{}).Where(query).Preload("CanvasRepository").Find(&canvasBranchPerms).Error
	return canvasBranchPerms, err
}

func (cr collectionPermissionRepo) GetRoleMembers(query map[string]interface{}) (*models.Role, error) {
	var role *models.Role
	postgres.GetDB().Model(&role).Where(query).Preload("Members").First(&role)
	return role, nil
}
