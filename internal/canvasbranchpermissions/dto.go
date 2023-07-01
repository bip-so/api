package canvasbranchpermissions

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

// Interface
type CanvasBranchPermissionRepo interface {
	CreateCanvasBranchPermission(role *models.CanvasBranchPermission) (*models.CanvasBranchPermission, error)
	GetCanvasBranchPermission(query map[string]interface{}) ([]models.CanvasBranchPermission, error)
	UpdateCanvasBranchPermissions(query map[string]interface{}, body NewCanvasBranchPermissionCreatePost, studioID uint64, authUserId uint64) (*models.CanvasBranchPermission, error)
}

// constructor
func NewCanvasBranchPermissionsRepo() CanvasBranchPermissionRepo {
	return &canvasBranchPermissionRepo{}
}

func (sr canvasBranchPermissionRepo) Get(query map[string]interface{}) (models.CanvasBranchPermission, error) {
	var canvasBranchPerms models.CanvasBranchPermission
	postgres.GetDB().Model(&canvasBranchPerms).Where(query).Preload("CanvasRepository").Preload("Member").First(&canvasBranchPerms)
	return canvasBranchPerms, nil
}

func (sr canvasBranchPermissionRepo) GetCanvasBranchPermission(query map[string]interface{}) ([]models.CanvasBranchPermission, error) {
	var canvasBranchPerms []models.CanvasBranchPermission
	postgres.GetDB().Model(&canvasBranchPerms).Where(query).Preload("Studio").Preload("Role").Preload("CanvasBranch").Preload("Member").Preload("Member.User").Preload("Role.Members").Preload("CanvasRepository").Preload("CbpParentCanvasRepository").Preload("Collection").Preload("CanvasRepository.DefaultBranch").Preload("CbpParentCanvasRepository.DefaultBranch").Find(&canvasBranchPerms)
	return canvasBranchPerms, nil
}

func (sr canvasBranchPermissionRepo) CreateCanvasBranchPermission(stdperms *models.CanvasBranchPermission) (*models.CanvasBranchPermission, error) {
	result := postgres.GetDB().Create(stdperms)
	if result.Error != nil {
		return nil, result.Error
	}
	return stdperms, nil
}

func (cr canvasBranchPermissionRepo) UpdateCanvasBranchPermissions(query map[string]interface{}, body NewCanvasBranchPermissionCreatePost, studioID uint64, authUserId uint64) (*models.CanvasBranchPermission, error) {
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
			if !(role.Name == models.SYSTEM_ADMIN_ROLE) {
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
		go func() {
			extraData := notifications.NotificationExtraData{
				CanvasRepoID:   canvasBranchPerms.CanvasRepositoryID,
				CanvasBranchID: *canvasBranchPerms.CanvasBranchID,
			}
			contentObject := models.CANVAS_BRANCH
			if canvasBranchPerms.RoleId == nil {
				fmt.Println("canvasbranch perms", canvasBranchPerms)
				notifications.App.Service.PublishNewNotification(notifications.CanvasInviteByName, authUserId, []uint64{canvasBranchPerms.Member.UserID}, &canvasBranchPerms.StudioID,
					nil, extraData, canvasBranchPerms.CanvasBranchID, &contentObject)
			} else {
				notifications.App.Service.PublishNewNotification(notifications.CanvasInviteByGroup, authUserId, nil, &canvasBranchPerms.StudioID,
					&[]uint64{*canvasBranchPerms.RoleId}, extraData, canvasBranchPerms.CanvasBranchID, &contentObject)
			}
		}()
		return canvasBranchPerms, err
	}

	return canvasBranchPerms, nil
}

func (cr canvasBranchPermissionRepo) GetMembersByUserIDs(userIds []uint64, studioID uint64) ([]models.Member, error) {
	var members []models.Member
	result := postgres.GetDB().Model(&models.Member{}).Where("user_id IN ? and studio_id = ?", userIds, studioID).Find(&members)

	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

func (sr canvasBranchPermissionRepo) GetCollectionPermission(query map[string]interface{}) ([]models.CollectionPermission, error) {
	var collectionPerm []models.CollectionPermission
	err := postgres.GetDB().Model(&collectionPerm).Where(query).Find(&collectionPerm).Error
	return collectionPerm, err
}

func (sr canvasBranchPermissionRepo) createCollectionPermissionByMemberID(collectionId uint64, studioId uint64, memberId uint64, isOverridden bool, permissionGroup string, authUserId uint64) (*models.CollectionPermission, error) {
	collectionPerms := &models.CollectionPermission{}
	collectionPerms.CollectionId = collectionId
	if memberId != 0 {
		collectionPerms.MemberId = &memberId
	}
	collectionPerms.StudioID = studioId
	collectionPerms.IsOverridden = isOverridden
	collectionPerms.PermissionGroup = permissionGroup
	err := sr.db.Create(collectionPerms).Error
	if err != nil {
		return nil, err
	}

	go func() {
		var collectionPerm models.CollectionPermission
		err = postgres.GetDB().Model(&collectionPerm).Where("collection_id = ? and member_id = ?", collectionPerms.CollectionId, collectionPerms.MemberId).Preload("Member").Find(&collectionPerm).Error
		if err != nil {
			fmt.Println("Error on fetching collectionpermission", err.Error())
		}
		extraData := notifications.NotificationExtraData{
			CanvasRepoID: collectionPerms.CollectionId,
		}
		contentObject := models.COLLECTION
		notifications.App.Service.PublishNewNotification(notifications.CollectionInviteByName, authUserId, []uint64{collectionPerms.Member.UserID}, &collectionPerms.StudioID,
			nil, extraData, &collectionPerms.ID, &contentObject)
	}()

	return collectionPerms, nil
}

type RoleMembersSerializer struct {
	MemberID   uint64 `json:"memberId"`
	Id         uint64 `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	IsSystem   bool   `json:"isSystem"`
	IsNonPerms bool   `json:"isNonPerms"`
}

func (sr canvasBranchPermissionRepo) GetMemberRolesByID(studioID uint64, memberID uint64) ([]RoleMembersSerializer, error) {
	var roles []RoleMembersSerializer
	result := postgres.GetDB().Raw(`
		select * from roles 
		left join role_members on role_members.role_id = roles.id
		where studio_id = ? and role_members.member_id = ?
	`, studioID, memberID).Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}

func (sr canvasBranchPermissionRepo) GetCanvasBranch(query map[string]interface{}) (*models.CanvasBranch, error) {
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

func (sr canvasBranchPermissionRepo) GetMember(query map[string]interface{}) (*models.Member, error) {
	var member models.Member
	err := postgres.GetDB().Model(&models.Member{}).Where(query).First(&member).Error
	if err != nil {
		//log.Fatalln(err)
		return nil, err
	}
	return &member, nil
}

func (sr canvasBranchPermissionRepo) GetCanvasRepos(query map[string]interface{}) ([]models.CanvasRepository, error) {
	var canvasRepos []models.CanvasRepository
	err := postgres.GetDB().Model(&canvasRepos).Where(query).Find(&canvasRepos).Error
	return canvasRepos, err
}

func (sr canvasBranchPermissionRepo) GetCanvasRepo(query map[string]interface{}) (*models.CanvasRepository, error) {
	var repo models.CanvasRepository
	err := sr.db.Model(&models.CanvasRepository{}).Where(query).First(&repo).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &repo, nil
}
