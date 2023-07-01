package canvasbranchpermissions

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
)

type ServiceInterface interface {
	Create()
}

func (s canvasBranchPermissionsService) Create(collectionId uint64, canvasBranchID uint64, canvasRepositoryID uint64, parentCanvasRepositoryId *uint64, permsgroup string, roleId *uint64, memberId *uint64, isOverriddenFlag bool, studioId uint64) (*models.CanvasBranchPermission, error) {
	cbp := &models.CanvasBranchPermission{
		CollectionId:                collectionId,
		StudioID:                    studioId,
		CanvasBranchID:              &canvasBranchID,
		CanvasRepositoryID:          canvasRepositoryID,
		CbpParentCanvasRepositoryID: parentCanvasRepositoryId,
		RoleId:                      roleId,
		MemberId:                    memberId,
		PermissionGroup:             permsgroup,
		IsOverridden:                isOverriddenFlag,
	}

	_, err := App.Repo.CreateCanvasBranchPermission(cbp)

	if err != nil {
		return nil, err
	}

	return cbp, nil
}

func (s canvasBranchPermissionsService) Get(query map[string]interface{}) ([]models.CanvasBranchPermission, error) {

	canvasBranchPerms, err := App.Repo.GetCanvasBranchPermission(query)

	if err != nil {
		return nil, err
	}

	return canvasBranchPerms, nil
}

func (s canvasBranchPermissionsService) UpdateCanvasBranchPermissions(query map[string]interface{}, body NewCanvasBranchPermissionCreatePost, studioId uint64, authUserId uint64) (*models.CanvasBranchPermission, error) {
	repo := NewCanvasBranchPermissionsRepo()
	collperms, err := repo.UpdateCanvasBranchPermissions(query, body, studioId, authUserId)
	if err != nil {
		return nil, err
	}

	return collperms, nil
}

func (s canvasBranchPermissionsService) InheritUserPermsToSubCanvas(body NewCanvasBranchPermissionCreatePost, studioId uint64, authUserId uint64, canvasRepoId uint64) {
	subCanvases, _ := App.Repo.GetCanvasRepos(map[string]interface{}{"parent_canvas_repository_id": canvasRepoId, "is_published": true})
	for _, canvas := range subCanvases {
		if canvas.ParentCanvasRepositoryID != nil {
			s.InheritUserPermsToSubCanvas(body, studioId, authUserId, canvas.ID)
		}
		// update the sub canvas permissions
		body.CanvasBranchId = *canvas.DefaultBranchID
		body.CanvasRepositoryID = canvas.ID
		body.CbpParentCanvasRepositoryID = *canvas.ParentCanvasRepositoryID
		if body.RoleID != 0 {
			_, _ = App.Service.UpdateCanvasBranchPermissions(map[string]interface{}{"canvas_branch_id": canvas.DefaultBranchID, "role_id": body.RoleID}, body, studioId, authUserId)
			permissions.App.Service.InvalidateRolePermissionCache(body.RoleID, studioId)
		} else {
			var userID uint64
			if body.UserID == 0 {
				member, _ := App.Repo.GetMember(map[string]interface{}{"id": body.MemberID})
				userID = member.UserID
			} else {
				userID = body.UserID
			}
			if canvas.CreatedByID == userID {
				continue
			}
			_, _ = App.Service.UpdateCanvasBranchPermissions(map[string]interface{}{"canvas_branch_id": canvas.DefaultBranchID, "member_id": body.MemberID}, body, studioId, authUserId)
			permissions.App.Service.InvalidateUserPermissionCache(userID, studioId)
		}
	}
}
