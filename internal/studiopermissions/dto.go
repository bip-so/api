package studiopermissions

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

// Interface
type StudioPermissionRepo interface {
	CreateStudioPermission(role *models.StudioPermission) (*models.StudioPermission, error)
	UpdateStudioPermissions(query map[string]interface{}, body CreateStudioPermissionsPost, studioID uint64) (*models.StudioPermission, error)
	GetStudioPermissions(query map[string]interface{}) ([]models.StudioPermission, error)
	GetStudioPermissionsByMemberIds(memberIds []uint64, roleIds []uint64) ([]models.StudioPermission, error)
}

// constructor
func NewStudioPermissionsRepo() StudioPermissionRepo {
	return &studioPermissionRepo{}
}

func (sr studioPermissionRepo) GetStudioPermission(query map[string]interface{}) (models.StudioPermission, error) {
	var studioPerms models.StudioPermission
	err := postgres.GetDB().Model(&studioPerms).Where(query).Find(&studioPerms).Error
	return studioPerms, err
}

func (sr studioPermissionRepo) GetStudioPermissions(query map[string]interface{}) ([]models.StudioPermission, error) {

	var studioPerms []models.StudioPermission
	err := postgres.GetDB().Model(&studioPerms).Where(query).Preload("Role").Preload("Studio").Preload("Member").Preload("Member.User").Preload("Role.Members").Find(&studioPerms).Error
	return studioPerms, err
}

func (sr studioPermissionRepo) UpdateStudioPermissions(query map[string]interface{}, body CreateStudioPermissionsPost, studioID uint64) (*models.StudioPermission, error) {
	var studioPerms *models.StudioPermission

	err := postgres.GetDB().Model(&studioPerms).Where(query).Preload("Role").Preload("Studio").Preload("Member.User").Preload("Role.Members").Preload("Member").First(&studioPerms).Error

	// create flow
	if err != nil {
		if body.RoleId == 0 {
			studioPerms.RoleId = nil
		} else {
			studioPerms.RoleId = &body.RoleId
		}

		if body.MemberId == 0 {
			studioPerms.MemberId = nil
		} else {
			studioPerms.MemberId = &body.MemberId
		}
	}

	// we are always updating these field if record if found or not found
	studioPerms.IsOverridden = body.IsOverriddenFlag
	studioPerms.PermissionGroup = body.PermsGroup
	studioPerms.StudioID = studioID

	// save will create a new record if it doesn't finds or update the existing record
	postgres.GetDB().Save(&studioPerms)

	if err != nil {
		// new record created we need to send preloaded data
		err = postgres.GetDB().Model(&studioPerms).Where(query).Preload("Role").Preload("Studio").Preload("Member.User").Preload("Role.Members").Preload("Member").First(&studioPerms).Error
		return studioPerms, err
	}

	return studioPerms, nil
}

func (sr studioPermissionRepo) CreateStudioPermission(stdperms *models.StudioPermission) (*models.StudioPermission, error) {
	result := postgres.GetDB().Create(stdperms)
	if result.Error != nil {
		return nil, result.Error
	}
	return stdperms, nil
}

func (sr studioPermissionRepo) GetStudioPermissionsByMemberIds(memberIds []uint64, roleIds []uint64) ([]models.StudioPermission, error) {
	var studioPerms []models.StudioPermission
	err := postgres.GetDB().Model(&studioPerms).Where("member_id IN ? OR role_id IN ?", memberIds, roleIds).Order("studio_id inc").Find(&studioPerms).Error
	return studioPerms, err
}
