package collectionpermissions

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
)

var (
	CollectionPermissionService collectionPermissionService
)

func (m collectionPermissionService) NewCollectionPermission(collectionId uint64, permsgroup string, roleId *uint64, memberId *uint64, isOverriddenFlag bool, studioId uint64) (*models.CollectionPermission, error) {

	sp := &models.CollectionPermission{
		CollectionId:    collectionId,
		StudioID:        studioId,
		PermissionGroup: permsgroup,
		IsOverridden:    isOverriddenFlag,
	}

	repo := NewCollectionPermissionsRepo()
	_, err := CollectionPermissionRepo.CreateCollectionPermission(repo, sp)
	if err != nil {
		return nil, err
	}

	return sp, nil
}

func (m collectionPermissionService) GetCollectionPermissions(query map[string]interface{}) ([]models.CollectionPermission, error) {

	repo := NewCollectionPermissionsRepo()
	collectionperms, err := repo.GetCollectionPermission(query)
	if err != nil {
		return nil, err
	}

	return collectionperms, nil
}

func (m collectionPermissionService) UpdateCollectionPermissions(query map[string]interface{}, body CollectionPermissionValidator, studioId uint64, authUserId uint64) (*models.CollectionPermission, error) {

	repo := NewCollectionPermissionsRepo()
	collperms, err := repo.UpdateCollectionPermissions(query, body, studioId, authUserId)
	if err != nil {
		return nil, err
	}

	return collperms, nil
}

func (s collectionPermissionService) UpdateCanvasBranchPermissions(query map[string]interface{}, body newCanvasBranchPermissionCreatePost, studioId uint64, authUserId uint64) (*models.CanvasBranchPermission, error) {
	collperms, err := App.Repo.UpdateCanvasBranchPermissions(query, body, studioId, authUserId)
	if err != nil {
		return nil, err
	}

	return collperms, nil
}

func (s collectionPermissionService) InheritUserPermsToCanvas(body CollectionPermissionValidator, studioId uint64, authUserId uint64) {
	subCanvases, _ := App.Repo.GetCanvasRepos(map[string]interface{}{"collection_id": body.CollectionId, "is_published": true})
	for _, canvas := range subCanvases {
		// update the canvas permissions
		// body.CanvasBranchId = *canvas.DefaultBranchID
		// body.CanvasRepositoryID = canvas.ID
		// body.CbpParentCanvasRepositoryID = *canvas.ParentCanvasRepositoryID
		newCanvasPerm := newCanvasBranchPermissionCreatePost{}
		newCanvasPerm.CanvasBranchId = *canvas.DefaultBranchID
		newCanvasPerm.CanvasRepositoryID = canvas.ID
		if canvas.ParentCanvasRepositoryID != nil {
			newCanvasPerm.CbpParentCanvasRepositoryID = *canvas.ParentCanvasRepositoryID
		}
		newCanvasPerm.MemberID = body.MemberID
		newCanvasPerm.RoleID = body.RoleID
		newCanvasPerm.CollectionId = body.CollectionId
		newCanvasPerm.PermGroup = permissiongroup.MapCollectionCanvasPerms[body.PermGroup]
		newCanvasPerm.IsOverridden = body.IsOverridden
		if body.RoleID != 0 {
			_, _ = App.Service.UpdateCanvasBranchPermissions(map[string]interface{}{"canvas_branch_id": canvas.DefaultBranchID, "role_id": body.RoleID}, newCanvasPerm, studioId, authUserId)
			permissions.App.Service.InvalidateRolePermissionCache(body.RoleID, studioId)
		} else {
			member, _ := App.Repo.GetMember(map[string]interface{}{"id": body.MemberID})
			if canvas.CreatedByID == member.UserID {
				continue
			}
			_, _ = App.Service.UpdateCanvasBranchPermissions(map[string]interface{}{"canvas_branch_id": canvas.DefaultBranchID, "member_id": body.MemberID}, newCanvasPerm, studioId, authUserId)
			permissions.App.Service.InvalidateUserPermissionCache(member.UserID, studioId)
		}
	}
}

func (s collectionPermissionService) DeleteInheritCollectionPermission(collectionPermission models.CollectionPermission) {
	query := map[string]interface{}{"collection_id": collectionPermission.CollectionId}
	if collectionPermission.RoleId != nil {
		query["role_id"] = *collectionPermission.RoleId
		role, _ := App.Repo.GetRoleMembers(map[string]interface{}{"id": *collectionPermission.RoleId})
		canvasRepos, _ := queries.App.CanvasRepoQuery.GetCanvasRepos(map[string]interface{}{"collection_id": collectionPermission.CollectionId})
		userIDsNeedViewMetadata := []models.Member{}
		for _, repo := range canvasRepos {
			for _, roleMember := range role.Members {
				if repo.CreatedByID == roleMember.UserID {
					//canvasBranchPerm, _ := queries.App.CanvasBranchPermissionQuery.GetCanvasBranchPermissions(map[string]interface{}{"member_id": roleMember.ID, "canvas_branch_id": *repo.DefaultBranchID, "permission_group": models.PGCanvasModerateSysName})
					//if len(canvasBranchPerm) > 0 {
					//perm := canvasBranchPerm[0]
					//perm.PermissionGroup = models.PGCanvasViewSysName
					//App.Repo.db.Save(&perm)
					//break
					//}
					userIDsNeedViewMetadata = append(userIDsNeedViewMetadata, roleMember)
				}
			}
		}
		canvasBranchPerms, _ := queries.App.CanvasBranchPermissionQuery.GetCanvasBranchPermissions(query)
		for _, branchPerm := range canvasBranchPerms {
			var col models.CanvasBranchPermission
			App.Repo.Manger.HardDeleteByID(col.TableName(), branchPerm.ID)
		}
		for _, mem := range userIDsNeedViewMetadata {
			collectionPerms, _ := App.Repo.GetCollectionPermission(map[string]interface{}{"collection_id": collectionPermission.CollectionId, "member_id": mem.ID, "permission_group": models.PGCollectionViewMetadataSysName})
			if len(collectionPerms) == 0 {
				_, err := permissions.App.Repo.CreateCollectionPermissionByMemberIDWithoutNotification(collectionPermission.CollectionId, collectionPermission.StudioID, mem.ID, false, models.PGCollectionViewMetadataSysName, mem.UserID)
				if err != nil {
					fmt.Println("Error in creating collection view metadata perm", err)
				}
			}
		}
		permissions.App.Service.InvalidateRolePermissionCache(*collectionPermission.RoleId, collectionPermission.StudioID)
	} else if collectionPermission.MemberId != nil {
		query["member_id"] = *collectionPermission.MemberId
		member, err := App.Repo.GetMember(map[string]interface{}{"id": *collectionPermission.MemberId})
		if err != nil {
			fmt.Println("Error in getting studio member", err)
			return
		}
		isUserCreatedCanvasPresent := false
		canvasBranchPerms, _ := App.Repo.GetCanvasBranchPermissions(query)
		for _, branchPerm := range canvasBranchPerms {
			if branchPerm.CanvasRepository.CreatedByID == member.UserID {
				isUserCreatedCanvasPresent = true
				//branchPerm.PermissionGroup = models.PGCanvasViewSysName
				//App.Repo.db.Save(&branchPerm)
				continue
			}
			var col models.CanvasBranchPermission
			App.Repo.Manger.HardDeleteByID(col.TableName(), branchPerm.ID)
		}
		if isUserCreatedCanvasPresent {
			collectionPerms, _ := App.Repo.GetCollectionPermission(map[string]interface{}{"collection_id": collectionPermission.CollectionId, "member_id": member.ID, "permission_group": models.PGCollectionViewMetadataSysName})
			if len(collectionPerms) == 0 {
				_, err = permissions.App.Repo.CreateCollectionPermissionByMemberIDWithoutNotification(collectionPermission.CollectionId, collectionPermission.StudioID, member.ID, false, models.PGCollectionViewMetadataSysName, member.UserID)
			}
		}
		permissions.App.Service.InvalidateUserPermissionCache(member.UserID, collectionPermission.StudioID)
	}
}

func (s collectionPermissionService) AddViewMetaDataPermOnCollection(collectionPermission models.CollectionPermission, userID uint64) {
	if collectionPermission.RoleId != nil {
		collectionPerms, _ := App.Repo.GetCollectionPermission(map[string]interface{}{"collection_id": collectionPermission.CollectionId, "role_id": *collectionPermission.RoleId, "permission_group": models.PGCollectionViewMetadataSysName})
		if len(collectionPerms) == 0 {
			_, err := permissions.App.Repo.CreateCollectionPermissionByRoleIDWithoutNotification(collectionPermission.CollectionId, collectionPermission.StudioID, *collectionPermission.RoleId, false, models.PGCollectionViewMetadataSysName, userID)
			if err != nil {
				fmt.Println("Error in creating role collection view metadata perm", err)
			}
		}
	} else if collectionPermission.MemberId != nil {
		collectionPerms, _ := App.Repo.GetCollectionPermission(map[string]interface{}{"collection_id": collectionPermission.CollectionId, "member_id": *collectionPermission.MemberId, "permission_group": models.PGCollectionViewMetadataSysName})
		if len(collectionPerms) == 0 {
			_, err := permissions.App.Repo.CreateCollectionPermissionByMemberIDWithoutNotification(collectionPermission.CollectionId, collectionPermission.StudioID, *collectionPermission.MemberId, false, models.PGCollectionViewMetadataSysName, userID)
			if err != nil {
				fmt.Println("Error in creating member collection view metadata perm", err)
			}
		}
	}
}
