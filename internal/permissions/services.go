package permissions

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
)

// CreateDefaultCollectionPermission When a collection is created this can be called to create a default collectionPermission, by adding studio admin role to it.
func (s permissionService) CreateDefaultCollectionPermission(collectionId uint64, userId uint64, studioId uint64) error {
	studioPerm, err := App.Repo.getStudioPermission(map[string]interface{}{"studio_id": studioId, "permission_group": permissiongroup.PG_STUIDO_ADMIN})
	if err != nil {
		return err
	}

	_, err = App.Repo.createCollectionPermission(
		collectionId, studioId, *studioPerm.RoleId, false, models.PGCollectionModerateSysName)
	if err != nil {
		return err
	}

	member, _ := App.Repo.GetMember(map[string]interface{}{"user_id": userId, "studio_id": studioId})
	if member != nil {
		_, err = App.Repo.CreateCollectionPermissionByMemberID(collectionId, studioId, member.ID, false, models.PGCollectionModerateSysName, userId)
		s.InvalidateCollectionPermissionCache(userId, studioId)
	}

	s.InvalidateCollectionPermissionCacheByRole(*studioPerm.RoleId, studioId)
	return nil
}

func (s permissionService) CreateCustomCanvasBranchPermission(collectionId uint64, userId uint64, studioId uint64, canvasRepositoryId uint64, canvasBranchId uint64, parentCanvasRepoID *uint64, pg string, memberID uint64) error {
	canvasBranchPermission, _ := App.Repo.GetCanvasBranchPermission(map[string]interface{}{"canvas_branch_id": canvasBranchId, "member_id": memberID})
	if canvasBranchPermission != nil {
		return nil
	}

	_, err := App.Repo.createCanvasBranchPermission(
		collectionId, pg, &memberID, false, studioId, canvasRepositoryId, canvasBranchId, parentCanvasRepoID, nil)
	if err != nil {
		return err
	}

	if parentCanvasRepoID != nil && *parentCanvasRepoID != 0 {
		s.InvalidateSubCanvasPermissionCache(userId, studioId, collectionId, *parentCanvasRepoID)
	} else {
		s.InvalidateCanvasPermissionCache(userId, studioId, collectionId)
	}
	return nil
}

func (s permissionService) CreateCanvasBranchPermission(collectionId uint64, userId uint64, studioId uint64, canvasRepositoryId uint64, canvasBranchId uint64, parentCanvasRepoID *uint64, pg string) error {
	member, err := App.Repo.GetMember(map[string]interface{}{"user_id": userId, "studio_id": studioId})
	if err != nil {
		return err
	}
	_, err = App.Repo.createCanvasBranchPermission(
		collectionId, pg, &member.ID, false, studioId, canvasRepositoryId, canvasBranchId, parentCanvasRepoID, nil)
	if err != nil {
		return err
	}

	s.AddMemberToCollectionIfNotPresent(userId, member.ID, collectionId, studioId)
	if parentCanvasRepoID != nil && *parentCanvasRepoID != 0 {
		s.AddMemberToCanvasIfNotPresent(userId, member.ID, *parentCanvasRepoID, studioId)
		s.InvalidateSubCanvasPermissionCache(userId, studioId, collectionId, *parentCanvasRepoID)
	} else {
		s.InvalidateCanvasPermissionCache(userId, studioId, collectionId)
	}
	return nil
}

func (s permissionService) CheckBranchAccessToken(inviteCode, repoKey string) (*models.BranchAccessToken, error) {
	canvasBranchAccessToken, err := App.Repo.GetBranchAccessToken(map[string]interface{}{"invite_code": inviteCode, "is_active": true})
	if canvasBranchAccessToken == nil || err != nil {
		return nil, err
	}
	canvasBranch, _ := App.Repo.GetBranchWithRepo(map[string]interface{}{"id": canvasBranchAccessToken.BranchID})
	if canvasBranch.CanvasRepository.Key != repoKey {
		return nil, err
	}
	return canvasBranchAccessToken, nil
}
