package canvasbranchpermissions

import (
	"errors"
	"fmt"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/internal/studio"
)

/// Check if Canvas and collections
// Non existing
// main check

func (c canvasBranchPermissionsController) CreateCanvasBranchPermissionController(body NewCanvasBranchPermissionCreatePost, studioId uint64, authUserId uint64, inherit string) (*models.CanvasBranchPermission, error) {
	var err error
	var canvasBranchPerm *models.CanvasBranchPermission

	// Neither Role or Member iD found.
	// Todo: this should be on route handler
	if body.RoleID != 0 && body.MemberID != 0 {
		err = errors.New("provide either roleId or memberId")
	}

	// Role is attached to the CB
	if body.RoleID != 0 {
		canvasBranchPerm, err = App.Service.UpdateCanvasBranchPermissions(map[string]interface{}{"canvas_branch_id": body.CanvasBranchId, "role_id": body.RoleID}, body, studioId, authUserId)

	} else if body.MemberID != 0 {
		// Member is attached to the CB
		canvasRepo, _ := App.Repo.GetCanvasRepo(map[string]interface{}{"id": body.CanvasRepositoryID})
		member, _ := App.Repo.GetMember(map[string]interface{}{"id": body.MemberID})
		if canvasRepo.CreatedByID == member.UserID {
			return nil, errors.New("cannot update canvas creator")
		}
		canvasBranchPerm, err = App.Service.UpdateCanvasBranchPermissions(map[string]interface{}{"canvas_branch_id": body.CanvasBranchId, "member_id": body.MemberID}, body, studioId, authUserId)

	} else {
		err = errors.New("roleId or memberId not provided")
	}

	if inherit == "true" {
		// branchInstance, _ := App.Repo.GetCanvasBranch(map[string]interface{}{"id": body.CanvasBranchId})
		// permissions.App.Service.InheritParentPermissions(*branchInstance)
		App.Service.InheritUserPermsToSubCanvas(body, studioId, authUserId, body.CanvasRepositoryID)
	}

	// clearing redis cache
	// @todo Later move this to kafka
	go func() {
		inv := &permissions.InvalidatePermissions{
			MemberID:     &body.MemberID,
			RoleID:       &body.RoleID,
			CollectionID: body.CollectionId,
		}
		if body.CbpParentCanvasRepositoryID != 0 {
			inv.InvalidationOn = "subCanvas"
			inv.ParentCanvasID = &body.CbpParentCanvasRepositoryID
		} else {
			inv.InvalidationOn = "canvas"
		}
		err = permissions.App.Service.InvalidatePermissions(inv)
		if err != nil {
			fmt.Println(err)
		}

		if body.MemberID != 0 {
			permissions.App.Service.AddMemberToCollectionIfNotPresent(authUserId, body.MemberID, body.CollectionId, studioId)
			if body.CbpParentCanvasRepositoryID != 0 {
				permissions.App.Service.AddMemberToCanvasIfNotPresent(authUserId, body.MemberID, body.CbpParentCanvasRepositoryID, studioId)
			}
		} else {
			permissions.App.Service.AddRoleToCollectionIfNotPresent(authUserId, body.RoleID, body.CollectionId, studioId)
			if body.CbpParentCanvasRepositoryID != 0 {
				permissions.App.Service.AddMemberToCanvasIfNotPresent(authUserId, body.RoleID, body.CbpParentCanvasRepositoryID, studioId)
			}
		}
	}()

	return canvasBranchPerm, err
}

func (c canvasBranchPermissionsController) getCanvasBranchPermissionController(query map[string]interface{}) ([]models.CanvasBranchPermission, error) {

	canvasBranchPerm, err := App.Service.Get(query)
	return canvasBranchPerm, err
}

func (c canvasBranchPermissionsController) BulkCreateCanvasBranchPermissionController(bulkMembers []NewCanvasBranchPermissionCreatePost, studioId uint64, authUserId uint64, inheritPerms string) (*[]CanvasBranchPermissionSerializer, error) {

	var err error
	var canvasBranchPerm *models.CanvasBranchPermission
	canvasBranchPermData := []CanvasBranchPermissionSerializer{}
	// Looping to get the userIDs
	userIDs := []uint64{}
	for _, member := range bulkMembers {
		if member.UserID != 0 {
			userIDs = append(userIDs, member.UserID)
		}
	}

	if len(userIDs) > 0 {
		joinStudioBulk := studio.JoinStudioBulkPost{UsersAdded: userIDs}
		_, err = studio.App.Controller.JoinStudioInBulkController(joinStudioBulk, studioId, authUserId)
		if err != nil {
			return nil, err
		}
		members, err := App.Repo.GetMembersByUserIDs(userIDs, studioId)
		if err != nil {
			return nil, err
		}

		// Looping to add memberId to the data
		for _, member := range members {
			for i, data := range bulkMembers {
				if member.UserID == data.UserID {
					bulkMembers[i].MemberID = member.ID
				}
			}
		}
	}

	// Looping to update the permissions
	for _, body := range bulkMembers {

		if body.RoleID != 0 && body.MemberID != 0 {
			err = errors.New("provide either roleId or memberId")
		}
		if body.RoleID != 0 {
			canvasBranchPerm, err = App.Service.UpdateCanvasBranchPermissions(map[string]interface{}{"canvas_branch_id": body.CanvasBranchId, "role_id": body.RoleID}, body, studioId, authUserId)

		} else if body.MemberID != 0 {
			canvasBranchPerm, err = App.Service.UpdateCanvasBranchPermissions(map[string]interface{}{"canvas_branch_id": body.CanvasBranchId, "member_id": body.MemberID}, body, studioId, authUserId)
		} else {
			err = errors.New("roleId or memberId not provided")
		}
		canvasBranchPermData = append(canvasBranchPermData, *SerializeCanvasBranchPermissionsPermission(canvasBranchPerm))
	}

	// clearing redis cache
	// @todo Later move this to kafka
	go func() {
		for _, body := range bulkMembers {
			fmt.Println("Inherit permission check", inheritPerms, inheritPerms == "true")
			if inheritPerms == "true" {
				App.Service.InheritUserPermsToSubCanvas(body, studioId, authUserId, body.CanvasRepositoryID)
			}

			inv := &permissions.InvalidatePermissions{
				MemberID:     &body.MemberID,
				RoleID:       &body.RoleID,
				CollectionID: body.CollectionId,
			}
			if body.CbpParentCanvasRepositoryID != 0 {
				inv.InvalidationOn = "subCanvas"
				inv.ParentCanvasID = &body.CbpParentCanvasRepositoryID
			} else {
				inv.InvalidationOn = "canvas"
			}
			err = permissions.App.Service.InvalidatePermissions(inv)
			if err != nil {
				fmt.Println(err)
			}

			if body.MemberID != 0 {
				permissions.App.Service.AddMemberToCollectionIfNotPresent(authUserId, body.MemberID, body.CollectionId, studioId)
				if body.CbpParentCanvasRepositoryID != 0 {
					permissions.App.Service.AddMemberToCanvasIfNotPresent(authUserId, body.MemberID, body.CbpParentCanvasRepositoryID, studioId)
				}
			} else {
				permissions.App.Service.AddRoleToCollectionIfNotPresent(authUserId, body.RoleID, body.CollectionId, studioId)
				if body.CbpParentCanvasRepositoryID != 0 {
					permissions.App.Service.AddRoleToCanvasIfNotPresent(authUserId, body.RoleID, body.CbpParentCanvasRepositoryID, studioId)
				}
			}
		}
	}()
	return &canvasBranchPermData, nil
}

func (c canvasBranchPermissionsController) inheritParentPermissionController(branchID uint64) ([]models.CanvasBranchPermission, error) {
	branchInstance, err := App.Repo.GetCanvasBranch(map[string]interface{}{"id": branchID})
	if err != nil {
		return nil, err
	}
	permissions.App.Service.InheritParentPermissions(branchInstance.ID)
	canvasBranchPerm, err := App.Service.Get(map[string]interface{}{"canvas_branch_id": branchID})
	return canvasBranchPerm, err
}
