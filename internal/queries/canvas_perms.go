package queries

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

func (q *canvasBranchPermissionQuery) GetCanvasBranchPermissions(query map[string]interface{}) ([]models.CanvasBranchPermission, error) {
	var canvasBranchPerms []models.CanvasBranchPermission
	postgres.GetDB().Model(&canvasBranchPerms).Where(query).Preload("Role").Preload("Role.Members").Preload("Member").Preload("Member.User").Find(&canvasBranchPerms)
	return canvasBranchPerms, nil
}

func (q *canvasBranchPermissionQuery) CreateCanvasBranchPermission(collectionId uint64, permsgroup string, memberId *uint64, isOverridden bool, studioId uint64, canvasRepositoryId uint64, canvasBranchId uint64, parentCanvasRepoID *uint64, roleId *uint64) (*models.CanvasBranchPermission, error) {
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
	err := postgres.GetDB().Create(cbp).Error
	if err != nil {
		return nil, err
	}
	return cbp, nil
}
