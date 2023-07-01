package queries

import (
	"context"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

const (
	StudioPermissionRedisKey     = "studioPermissions:"
	CollectionPermissionRedisKey = "collectionPermissions:"
	CanvasPermissionRedisKey     = "canvasPermissions:"
	PermissionsHash              = "permissions:"
)

func (q permsQuery) createCanvasBranchPermission(collectionId uint64, permsgroup string, memberId *uint64, isOverridden bool, studioId uint64, canvasRepositoryId uint64, canvasBranchId uint64, parentCanvasRepoID *uint64, roleId *uint64) (*models.CanvasBranchPermission, error) {
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

// CreateDefaultCanvasBranchPermission When a canvas is created this can be called to create default canvasBranchPermission by memberId as the creator.
func (q permsQuery) CreateDefaultCanvasBranchPermission(collectionId uint64, userId uint64, studioId uint64, canvasRepositoryId uint64, canvasBranchId uint64, parentCanvasRepoID *uint64) error {
	member, err := App.MemberQuery.GetMember(map[string]interface{}{"user_id": userId, "studio_id": studioId})
	if err != nil {
		return err
	}
	_, err = q.createCanvasBranchPermission(
		collectionId, "pg_canvas_branch_moderate", &member.ID, false, studioId, canvasRepositoryId, canvasBranchId, parentCanvasRepoID, nil)
	if err != nil {
		return err
	}

	if parentCanvasRepoID != nil && *parentCanvasRepoID != 0 {
		q.WorkflowHelperInvalidateSubCanvasPermissionCache(userId, studioId, collectionId, *parentCanvasRepoID)
	} else {
		q.WorkflowHelperInvalidateCanvasPermissionCache(userId, studioId, collectionId)
	}
	return nil
}

// 4 functions below are repeated to remove the permissions cyclic dependency

func (q permsQuery) WorkflowHelperSubCanvasPermissionsRedisKey(userID uint64, studioID uint64, collectionID uint64, canvasID uint64) string {
	return CanvasPermissionRedisKey + utils.String(userID) + ":" + utils.String(studioID) + ":" + utils.String(collectionID) + ":" + utils.String(canvasID)
}
func (q permsQuery) WorkflowHelperCanvasPermissionsRedisKey(userID uint64, studioID uint64, collectionID uint64) string {
	return CanvasPermissionRedisKey + utils.String(userID) + ":" + utils.String(studioID) + ":" + utils.String(collectionID)
}

func (q permsQuery) WorkflowHelperInvalidateSubCanvasPermissionCache(userID uint64, studioID uint64, collectionID uint64, canvasRepositoryID uint64) {
	redisKey := q.WorkflowHelperSubCanvasPermissionsRedisKey(userID, studioID, collectionID, canvasRepositoryID)
	q.cache.HDelete(context.Background(), PermissionsHash+utils.String(userID), redisKey)
}

func (q permsQuery) WorkflowHelperInvalidateCanvasPermissionCache(userID uint64, studioID uint64, collectionID uint64) {
	redisKey := q.WorkflowHelperCanvasPermissionsRedisKey(userID, studioID, collectionID)
	q.cache.HDelete(context.Background(), PermissionsHash+utils.String(userID), redisKey)
}
