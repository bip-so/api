package permissions

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

/*
	CalculateCanvasRepoPermissions: Calculates the permissions of the user on canvasBranches which are on root level of a specific collection.

	Algorithm:
		- Get all the members of the studio based on userId, studioId, collectionId
		- Get roleIds & memberIds by looping members.
		- Here user can be part of multiple roles or a member with permission & each role will have different permissions.
          We need to get the max. weight permission for that user & canvas.
		- Get canvasBranchPermissions based on memberIds & roleIds
		- Now we loop canvasBranchPermissions to find the max. Weight permission of a collection.

	Cache:
		- We use REDIS cache here. So first we check for the cache data in REDIS, if data is present we return the data
          or, we calculate the data and SET the data to REDIS and return the data.

	Example of how data store is:
	canvasPermissions:userId:studioId:collectionId = {
		"canvasId1": {
			"branchId1": "permission1",
			"branchId2": "permission2",
		},
		"canvasId2": {
			"branchId1": "permission1",
			"branchId2": "permission2",
		},...
	}

	Args:
		userid uint64
		studioID uint64
		collectionID uint64
	Returns:
		permissionList map[uint64]map[uint64]string
*/
func (s permissionService) CalculateCanvasRepoPermissions(userID, studioID, collectionID uint64) (map[uint64]map[uint64]string, error) {

	permissionList := make(map[uint64]map[uint64]string)
	isOverriddenList := make(map[uint64]bool)

	// Getting the data from redis
	ctx := context.Background()
	redisKey := s.CanvasPermissionsRedisKey(userID, studioID, collectionID)
	value := s.cache.HGet(ctx, PermissionsHash+utils.String(userID), redisKey)
	err := json.Unmarshal([]byte(value), &permissionList)
	if err == nil {
		return permissionList, nil
	}

	studioMembers, err := App.Repo.getMembersByUserIDs([]uint64{userID}, studioID)
	if err != nil {
		return nil, err
	}

	var memberIds []uint64
	var roleIds []uint64
	for _, memb := range studioMembers {
		memberIds = append(memberIds, memb.ID)
		for _, role := range memb.Roles {
			roleIds = append(roleIds, role.ID)
		}
	}

	canvasRepoPerms, err := App.Repo.getCanvasPermissionsByMemberIds(memberIds, roleIds, studioID, collectionID)
	if err != nil {
		return nil, err
	}

	for _, perm := range canvasRepoPerms {
		if len(permissionList[perm.CanvasRepositoryID]) == 0 {
			permissionList[perm.CanvasRepositoryID] = make(map[uint64]string)
		}
		if isOverriddenList[*perm.CanvasBranchID] || (perm.MemberId != nil && perm.IsOverridden) {
			permissionList[perm.CanvasRepositoryID][*perm.CanvasBranchID] = perm.PermissionGroup
			isOverriddenList[*perm.CanvasBranchID] = true
		} else {
			var currPermissionWt = models.PermissionGroupWeightMap[permissiongroup.PGTYPECANVAS][perm.PermissionGroup]
			var storedPermissionWt = models.PermissionGroupWeightMap[permissiongroup.PGTYPECANVAS][permissionList[perm.CanvasRepositoryID][*perm.CanvasBranchID]]
			if currPermissionWt > storedPermissionWt {
				permissionList[perm.CanvasRepositoryID][*perm.CanvasBranchID] = perm.PermissionGroup
			}
		}
	}
	// Set the data to redis
	permissionListString, _ := json.Marshal(permissionList)
	s.cache.HSet(ctx, PermissionsHash+utils.String(userID), redisKey, permissionListString)

	return permissionList, nil
}

// CanvasPermissionsRedisKey returns the redis key for canvases
func (s permissionService) CanvasPermissionsRedisKey(userID uint64, studioID uint64, collectionID uint64) string {
	return CanvasPermissionRedisKey + utils.String(userID) + ":" + utils.String(studioID) + ":" + utils.String(collectionID)
}

/*
	SetCanvasPermissionsOnPublishRequest
	Args:
		branch models.CanvasBranch

	On publish request we are adding the parent permissions to the published branch.
	If parent is canvas repo then we consider the permissions of the canvasRepo.
	else we consider the collectionPermissions and add it to the branch
*/
func (s permissionService) SetCanvasPermissionsOnPublishRequest(branch models.CanvasBranch) {
	if branch.CanvasRepository.ParentCanvasRepositoryID != nil {
		// parent as canvas flow

		// Get canvasBranch permissions of parentCanvas default branch
		branchPerms, err := App.Repo.GetCanvasBranchPerms(map[string]interface{}{"canvas_branch_id": *branch.CanvasRepository.ParentCanvasRepository.DefaultBranchID, "permission_group": models.PGCanvasModerateSysName})
		//branchPerms, err := App.Repo.GetCanvasPermissionsByID(*branch.CanvasRepository.ParentCanvasRepository.DefaultBranchID)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		for _, perm := range branchPerms {
			if perm.MemberId != nil && perm.Member.UserID == branch.CanvasRepository.CreatedByID {
				continue
			}
			_, err = App.Repo.createCanvasBranchPermission(
				branch.CanvasRepository.CollectionID, perm.PermissionGroup, perm.MemberId, perm.IsOverridden,
				branch.CanvasRepository.StudioID, branch.CanvasRepositoryID, branch.ID, branch.CanvasRepository.ParentCanvasRepositoryID, perm.RoleId)
			if err != nil {
				logger.Error(err.Error())
			}
			// Invalidating the user permissions
			if perm.RoleId != nil {
				s.InvalidateSubCanvasPermissionCacheByRole(*perm.RoleId, branch.CanvasRepository.StudioID, branch.CanvasRepository.CollectionID, *branch.CanvasRepository.ParentCanvasRepositoryID)
			} else {
				s.InvalidateSubCanvasPermissionCache(perm.Member.UserID, branch.CanvasRepository.StudioID, branch.CanvasRepository.CollectionID, *branch.CanvasRepository.ParentCanvasRepositoryID)
			}
		}
	} else {
		// parent as collection flow
		collectionPerms, err := App.Repo.getCollectionPermissionsByID(map[string]interface{}{"collection_id": branch.CanvasRepository.CollectionID, "permission_group": models.PGCollectionModerateSysName})
		if err != nil {
			logger.Error(err.Error())
			return
		}
		for _, perm := range collectionPerms {
			if perm.Member != nil && perm.Member.UserID == branch.CanvasRepository.CreatedByID {
				continue
			}
			canvasPermGroup := permissiongroup.MapCollectionCanvasPerms[perm.PermissionGroup]
			_, err = App.Repo.createCanvasBranchPermission(
				branch.CanvasRepository.CollectionID, canvasPermGroup, perm.MemberId, perm.IsOverridden,
				branch.CanvasRepository.StudioID, branch.CanvasRepositoryID, branch.ID, nil, perm.RoleId)
			if err != nil {
				logger.Error(err.Error())
			}
			if perm.RoleId != nil {
				s.InvalidateCanvasPermissionCacheByRole(*perm.RoleId, branch.CanvasRepository.StudioID, branch.CanvasRepository.CollectionID)
			} else {
				fmt.Println(perm.Member, branch.CanvasRepository)
				s.InvalidateCanvasPermissionCache(perm.Member.UserID, branch.CanvasRepository.StudioID, branch.CanvasRepository.CollectionID)
			}
		}
	}
}

// CanvasBranchPermissionGroupByUserID returns the permissions of the canvases which parent as collection
func (s permissionService) CanvasBranchPermissionGroupByUserID(userID uint64, branchID uint64) (*string, error) {
	var permissionGroup string
	var permissionsList map[uint64]map[uint64]string

	canvasBranch, err := App.Repo.GetCanvasBranch(map[string]interface{}{"id": branchID})
	if err != nil {
		return nil, err
	}
	canvasRepository := canvasBranch.CanvasRepository

	if canvasRepository.ParentCanvasRepositoryID == nil {
		permissionsList, err = s.CalculateCanvasRepoPermissions(userID, canvasRepository.StudioID, canvasRepository.CollectionID)
	} else {
		permissionsList, err = s.CalculateSubCanvasRepoPermissions(userID, canvasRepository.StudioID, canvasRepository.CollectionID, *canvasRepository.ParentCanvasRepositoryID)
	}
	if err != nil {
		return nil, err
	}

	permissionGroup = permissionsList[canvasRepository.ID][canvasBranch.ID]
	return &permissionGroup, nil
}

// If you want to check if a User can do something on Branch, Pass UserID, BranchID, Permission to Check.
// This returns True/False or Error
func (s permissionService) CanUserDoThisOnBranch(userID uint64, branchID uint64, permissionGroup string) (bool, error) {
	var branchPermissionGroup string
	var permissionsList map[uint64]map[uint64]string
	canvasBranch, err := App.Repo.GetCanvasBranch(map[string]interface{}{"id": branchID})
	if err != nil {
		return false, err
	}

	if userID == 0 && canvasBranch.PublicAccess == models.PRIVATE {
		return false, nil
	}

	var branchAccess int
	if canvasBranch.PublicAccess != models.PRIVATE {
		if canvasBranch.PublicAccess == models.EDIT {
			branchAccess = models.CanvasPermissionsMap[models.PGCanvasEditSysName][permissionGroup]
		} else if canvasBranch.PublicAccess == models.VIEW {
			branchAccess = models.CanvasPermissionsMap[models.PGCanvasViewSysName][permissionGroup]
		} else if canvasBranch.PublicAccess == models.COMMENT {
			branchAccess = models.CanvasPermissionsMap[models.PGCanvasCommentSysName][permissionGroup]
		}
	}
	if branchAccess != 0 {
		return true, nil
	}

	canvasRepository := canvasBranch.CanvasRepository
	if canvasRepository.ParentCanvasRepositoryID == nil {
		permissionsList, err = s.CalculateCanvasRepoPermissions(userID, canvasRepository.StudioID, canvasRepository.CollectionID)
	} else {
		permissionsList, err = s.CalculateSubCanvasRepoPermissions(userID, canvasRepository.StudioID, canvasRepository.CollectionID, *canvasRepository.ParentCanvasRepositoryID)
	}
	if err != nil {
		return false, err
	}
	branchPermissionGroup = permissionsList[canvasRepository.ID][canvasBranch.ID]
	userBranchAccess := models.CanvasPermissionsMap[branchPermissionGroup][permissionGroup]
	if userBranchAccess != 0 {
		return true, nil
	} else {
		return false, err
	}
}

/*
	InheritParentPermissions
	Args:
		branch models.CanvasBranch

	Inherit the permissions of the parent in this method.
	If parent as canvas we consider the permissions of the canvas
	else if branch is a language branch we consider the main canvas of that branch
	else we consider the permission of the collection
*/
func (s permissionService) InheritParentPermissions(branchID uint64) {
	branch, _ := App.Repo.GetCanvasBranchWithPreload(map[string]interface{}{"id": branchID})
	if branch.CanvasRepository.ParentCanvasRepositoryID != nil && !branch.CanvasRepository.IsLanguageCanvas {
		// parent as canvas flow
		branchPerms, err := App.Repo.GetCanvasPermissionsByID(*branch.CanvasRepository.ParentCanvasRepository.DefaultBranchID)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		// current branch perms
		currentBranchPerms, err := App.Repo.GetCanvasPermissionsByID(*branch.CanvasRepository.DefaultBranchID)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		for _, perm := range branchPerms {
			if perm.PermissionGroup == models.PGCanvasViewMetadataSysName {
				continue
			}
			permissionPresent := false
			for _, currentBranchPerm := range currentBranchPerms {
				if (currentBranchPerm.RoleId != nil && perm.RoleId != nil && *currentBranchPerm.RoleId == *perm.RoleId && perm.PermissionGroup == currentBranchPerm.PermissionGroup) ||
					(currentBranchPerm.MemberId != nil && perm.MemberId != nil && *perm.MemberId == *currentBranchPerm.MemberId && perm.PermissionGroup == currentBranchPerm.PermissionGroup) {
					permissionPresent = true
					break
				}
				if (currentBranchPerm.RoleId != nil && perm.RoleId != nil && *currentBranchPerm.RoleId == *perm.RoleId) ||
					(perm.MemberId != nil && currentBranchPerm.MemberId != nil && *currentBranchPerm.MemberId == *perm.MemberId) {
					// if weightage is less we update the weightage of same permission
					currentBranchPermWeightage := models.PermissionGroupWeightMap["canvas"][currentBranchPerm.PermissionGroup]
					permWeightage := models.PermissionGroupWeightMap["canvas"][perm.PermissionGroup]
					if permWeightage > currentBranchPermWeightage {
						// update canvas branch permission group
						App.Repo.UpdateCanvasBranchPermission(map[string]interface{}{"id": currentBranchPerm.ID}, map[string]interface{}{"permission_group": perm.PermissionGroup})
					}
					permissionPresent = true
					break
				}
			}
			if !permissionPresent {
				_, err = App.Repo.createCanvasBranchPermission(
					branch.CanvasRepository.CollectionID, perm.PermissionGroup, perm.MemberId, perm.IsOverridden,
					branch.CanvasRepository.StudioID, branch.CanvasRepositoryID, branch.ID, branch.CanvasRepository.ParentCanvasRepositoryID, perm.RoleId)
				if err != nil {
					logger.Error(err.Error())
				}
			}
			// Invalidating the user permissions
			if perm.RoleId != nil {
				s.InvalidateSubCanvasPermissionCacheByRole(*perm.RoleId, branch.CanvasRepository.StudioID, branch.CanvasRepository.CollectionID, *branch.CanvasRepository.ParentCanvasRepositoryID)
			} else {
				s.InvalidateSubCanvasPermissionCache(perm.Member.UserID, branch.CanvasRepository.StudioID, branch.CanvasRepository.CollectionID, *branch.CanvasRepository.ParentCanvasRepositoryID)
			}
		}

		// Updating branch public access to parent
		queries.App.BranchQuery.UpdateBranchInstance(branch.ID,
			map[string]interface{}{"public_access": branch.CanvasRepository.ParentCanvasRepository.DefaultBranch.PublicAccess})
	} else if branch.CanvasRepository.IsLanguageCanvas {
		languageParentCanvasRepo, _ := App.Repo.GetCanvasRepo(map[string]interface{}{"id": *branch.CanvasRepository.DefaultLanguageCanvasRepoID})
		branchPerms, err := App.Repo.GetCanvasPermissionsByID(*languageParentCanvasRepo.DefaultBranchID)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		// current branch perms
		currentBranchPerms, err := App.Repo.GetCanvasPermissionsByID(*branch.CanvasRepository.DefaultBranchID)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		for _, perm := range branchPerms {
			if perm.PermissionGroup == models.PGCanvasViewMetadataSysName {
				continue
			}
			permissionPresent := false
			for _, currentBranchPerm := range currentBranchPerms {
				if (currentBranchPerm.RoleId != nil && perm.RoleId != nil && *currentBranchPerm.RoleId == *perm.RoleId && perm.PermissionGroup == currentBranchPerm.PermissionGroup) ||
					(currentBranchPerm.MemberId != nil && perm.MemberId != nil && *perm.MemberId == *currentBranchPerm.MemberId && perm.PermissionGroup == currentBranchPerm.PermissionGroup) {
					permissionPresent = true
					break
				}
				if (currentBranchPerm.RoleId != nil && perm.RoleId != nil && *currentBranchPerm.RoleId == *perm.RoleId) ||
					(perm.MemberId != nil && currentBranchPerm.MemberId != nil && *currentBranchPerm.MemberId == *perm.MemberId) {
					// if weightage is less we update the weightage of same permission
					currentBranchPermWeightage := models.PermissionGroupWeightMap["canvas"][currentBranchPerm.PermissionGroup]
					permWeightage := models.PermissionGroupWeightMap["canvas"][perm.PermissionGroup]
					if permWeightage > currentBranchPermWeightage {
						// update canvas branch permission group
						App.Repo.UpdateCanvasBranchPermission(map[string]interface{}{"id": currentBranchPerm.ID}, map[string]interface{}{"permission_group": perm.PermissionGroup})
					}
					permissionPresent = true
					break
				}
			}
			if !permissionPresent {
				_, err = App.Repo.createCanvasBranchPermission(
					branch.CanvasRepository.CollectionID, perm.PermissionGroup, perm.MemberId, perm.IsOverridden,
					branch.CanvasRepository.StudioID, branch.CanvasRepositoryID, branch.ID, branch.CanvasRepository.ParentCanvasRepositoryID, perm.RoleId)
				if err != nil {
					logger.Error(err.Error())
				}
			}
			// Invalidating the user permissions
			if branch.CanvasRepository.ParentCanvasRepositoryID != nil {
				if perm.RoleId != nil {
					s.InvalidateSubCanvasPermissionCacheByRole(*perm.RoleId, branch.CanvasRepository.StudioID, branch.CanvasRepository.CollectionID, *branch.CanvasRepository.ParentCanvasRepositoryID)
				} else {
					s.InvalidateSubCanvasPermissionCache(perm.Member.UserID, branch.CanvasRepository.StudioID, branch.CanvasRepository.CollectionID, *branch.CanvasRepository.ParentCanvasRepositoryID)
				}
			} else {
				if perm.RoleId != nil {
					s.InvalidateCanvasPermissionCacheByRole(*perm.RoleId, branch.CanvasRepository.StudioID, branch.CanvasRepository.CollectionID)
				} else {
					s.InvalidateCanvasPermissionCache(perm.Member.UserID, branch.CanvasRepository.StudioID, branch.CanvasRepository.CollectionID)
				}
			}
		}
		// Updating branch public access to parent
		langBranchInstance, _ := queries.App.BranchQuery.GetBranchByID(*languageParentCanvasRepo.DefaultBranchID)
		queries.App.BranchQuery.UpdateBranchInstance(branch.ID,
			map[string]interface{}{"public_access": langBranchInstance.PublicAccess})
	} else {
		// parent as collection flow
		collectionPerms, err := App.Repo.getCollectionPermissionsByID(map[string]interface{}{"collection_id": branch.CanvasRepository.CollectionID})
		if err != nil {
			logger.Error(err.Error())
			return
		}
		// current branch perms
		currentBranchPerms, err := App.Repo.GetCanvasPermissionsByID(*branch.CanvasRepository.DefaultBranchID)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		for _, perm := range collectionPerms {
			if perm.PermissionGroup == models.PGCollectionViewMetadataSysName {
				continue
			}
			permissionPresent := false
			for _, currentBranchPerm := range currentBranchPerms {
				collectionPermGroup := permissiongroup.MapCollectionCanvasPerms[perm.PermissionGroup]
				if (currentBranchPerm.RoleId != nil && perm.RoleId != nil && *currentBranchPerm.RoleId == *perm.RoleId && collectionPermGroup == currentBranchPerm.PermissionGroup) ||
					(perm.MemberId != nil && currentBranchPerm.MemberId != nil && *perm.MemberId == *currentBranchPerm.MemberId && collectionPermGroup == currentBranchPerm.PermissionGroup) {
					permissionPresent = true
					break
				}
				if (currentBranchPerm.RoleId != nil && perm.RoleId != nil && *currentBranchPerm.RoleId == *perm.RoleId) ||
					(perm.MemberId != nil && currentBranchPerm.MemberId != nil && *currentBranchPerm.MemberId == *perm.MemberId) {
					// if weightage is less we update the weightage of same permission
					currentBranchPermWeightage := models.PermissionGroupWeightMap["canvas"][currentBranchPerm.PermissionGroup]
					permWeightage := models.PermissionGroupWeightMap["canvas"][collectionPermGroup]
					if permWeightage > currentBranchPermWeightage {
						// update canvas branch permission group
						App.Repo.UpdateCanvasBranchPermission(map[string]interface{}{"id": currentBranchPerm.ID}, map[string]interface{}{"permission_group": collectionPermGroup})
					}
					permissionPresent = true
					break
				}
			}
			if !permissionPresent {
				canvasPermGroup := permissiongroup.MapCollectionCanvasPerms[perm.PermissionGroup]
				_, err = App.Repo.createCanvasBranchPermission(
					branch.CanvasRepository.CollectionID, canvasPermGroup, perm.MemberId, perm.IsOverridden,
					branch.CanvasRepository.StudioID, branch.CanvasRepositoryID, branch.ID, nil, perm.RoleId)
				if err != nil {
					logger.Error(err.Error())
				}
			}
			if perm.RoleId != nil {
				s.InvalidateCanvasPermissionCacheByRole(*perm.RoleId, branch.CanvasRepository.StudioID, branch.CanvasRepository.CollectionID)
			} else {
				s.InvalidateCanvasPermissionCache(perm.Member.UserID, branch.CanvasRepository.StudioID, branch.CanvasRepository.CollectionID)
			}
		}

		// Updating branch public access to parent
		queries.App.BranchQuery.UpdateBranchInstance(branch.ID,
			map[string]interface{}{"public_access": branch.CanvasRepository.Collection.PublicAccess})
	}
}

/*
	AddMemberToCanvasIfNotPresent

	We are adding viewMetaData Permission to the parent canvases when we add a permission to the user to some canvas
	where user doesn't have permission to their parents.

	Here we are calling this method recursively to add viewMetaDataPermission to all the parent canvases
*/
func (s permissionService) AddMemberToCanvasIfNotPresent(authUserId, memberID, canvasID, studioID uint64) {
	canvasRepo, _ := App.Repo.GetCanvasRepo(map[string]interface{}{"id": canvasID})
	if canvasRepo.ParentCanvasRepositoryID != nil {
		s.AddMemberToCanvasIfNotPresent(authUserId, memberID, *canvasRepo.ParentCanvasRepositoryID, studioID)
	}
	canvasPerms, _ := App.Repo.GetCanvasBranchPerms(map[string]interface{}{"member_id": memberID, "canvas_repository_id": canvasID})
	if len(canvasPerms) == 0 {
		// roles of the member
		memberFound := false
		roles, _ := App.Repo.GetMemberRolesByID(studioID, memberID)
		if len(roles) > 0 {
			// Get all collectionPerms
			canvasPerms, _ = App.Repo.GetCanvasBranchPerms(map[string]interface{}{"canvas_repository_id": canvasID})
			for _, perm := range canvasPerms {
				if perm.RoleId != nil {
					for _, role := range roles {
						if *perm.RoleId == role.Id {
							memberFound = true
							break
						}
					}
				}
				if memberFound {
					break
				}
			}
		}
		if !memberFound {
			App.Repo.UpdateCanvasBranchPermissions(map[string]interface{}{"canvas_branch_id": canvasRepo.DefaultBranchID, "member_id": memberID}, canvasRepo.CollectionID, studioID, memberID, 0, false, models.PGCanvasViewMetadataSysName, authUserId, *canvasRepo.DefaultBranchID, canvasID, canvasRepo.ParentCanvasRepositoryID)
		}
		member, _ := App.Repo.GetMember(map[string]interface{}{"id": memberID})
		if canvasRepo.ParentCanvasRepositoryID != nil {
			App.Service.InvalidateSubCanvasPermissionCache(member.UserID, studioID, canvasRepo.CollectionID, *canvasRepo.ParentCanvasRepositoryID)
		} else {
			App.Service.InvalidateCanvasPermissionCache(member.UserID, studioID, canvasRepo.CollectionID)
		}
	}
}

// Checking if the role is present or not and adding role to the canvas.
func (s permissionService) AddRoleToCanvasIfNotPresent(authUserId, roleID, canvasID, studioID uint64) {
	canvasRepo, _ := App.Repo.GetCanvasRepo(map[string]interface{}{"id": canvasID})
	fmt.Println("canvasRepo", canvasRepo.ParentCanvasRepositoryID)
	if canvasRepo.ParentCanvasRepositoryID != nil {
		s.AddRoleToCanvasIfNotPresent(authUserId, roleID, *canvasRepo.ParentCanvasRepositoryID, studioID)
	}
	canvasPerms, _ := App.Repo.GetCanvasBranchPerms(map[string]interface{}{"role_id": roleID, "canvas_repository_id": canvasID})
	// memberFound := false
	if len(canvasPerms) == 0 {
		// roles of the member
		// role, _ := App.Repo.GetRoleByID(roleID)
		// if len(role.Members) > 0 {
		// 	// Get all collectionPerms
		// 	canvasPerms, _ = App.Repo.GetCanvasBranchPerms(map[string]interface{}{"canvas_repository_id": canvasID})
		// 	for _, perm := range canvasPerms {
		// 		if perm.MemberId != nil {
		// 			for _, member := range role.Members {
		// 				if member.ID == *perm.MemberId {
		// 					memberFound = true
		// 					break
		// 				}
		// 			}
		// 		}
		// 		if memberFound {
		// 			break
		// 		}
		// 	}
		// }
		// if !memberFound {
		App.Repo.UpdateCanvasBranchPermissions(map[string]interface{}{"canvas_branch_id": canvasRepo.DefaultBranchID, "role_id": roleID}, canvasRepo.CollectionID, studioID, 0, roleID, false, models.PGCanvasViewMetadataSysName, authUserId, *canvasRepo.DefaultBranchID, canvasID, canvasRepo.ParentCanvasRepositoryID)
		// }
		if canvasRepo.ParentCanvasRepositoryID != nil {
			App.Service.InvalidateSubCanvasPermissionCacheByRole(roleID, studioID, canvasRepo.CollectionID, *canvasRepo.ParentCanvasRepositoryID)
		} else {
			App.Service.InvalidateCanvasPermissionCacheByRole(roleID, studioID, canvasRepo.CollectionID)
		}
	}
}

/*
RemoveViewMetadataPermissionOnParents


you will get branchID, userID here
get canvasRepo of branchID
if canvasRepo.Parent
// Do some stuff
else
canvasBranchPerms by memberID, collectionID
if canvasBranchPerms == 0
delete viewMetadataPerm on collection if present
*/
func (s permissionService) RemoveMemberViewMetadataPermissionOnParents(canvasBranchID, memberID, userID, studioID uint64) {
	canvasBranch, _ := App.Repo.GetCanvasBranch(map[string]interface{}{"id": canvasBranchID})
	if canvasBranch.CanvasRepository.ParentCanvasRepositoryID != nil {
		canvasBranchPerms, _ := App.Repo.GetCanvasBranchPerms(map[string]interface{}{"member_id": memberID, "cbp_parent_canvas_repository_id": canvasBranch.CanvasRepository.ParentCanvasRepositoryID})
		if len(canvasBranchPerms) == 0 {
			branchPerms, _ := App.Repo.GetCanvasBranchPerms(map[string]interface{}{
				"member_id":        memberID,
				"canvas_branch_id": canvasBranch.CanvasRepository.ParentCanvasRepository.DefaultBranchID,
				"permission_group": models.PGCanvasViewMetadataSysName,
			})
			for _, perm := range branchPerms {
				err := App.Repo.Manger.HardDeleteByID(models.CANVAS_BRANCH_PERMISSION, perm.ID)
				if err != nil {
					fmt.Println("Error on removing canvas view metadata permission", err)
				}
			}
			App.Service.InvalidateSubCanvasPermissionCache(userID, studioID, canvasBranch.CanvasRepository.CollectionID, *canvasBranch.CanvasRepository.ParentCanvasRepositoryID)
		}
		s.RemoveMemberViewMetadataPermissionOnParents(*canvasBranch.CanvasRepository.ParentCanvasRepository.DefaultBranchID, memberID, userID, studioID)
	} else {
		// Collection perms update
		canvasBranchPerms, _ := App.Repo.GetCanvasBranchPerms(map[string]interface{}{"member_id": memberID, "collection_id": canvasBranch.CanvasRepository.CollectionID})
		if len(canvasBranchPerms) == 0 {
			collectionPerms, _ := App.Repo.getCollectionPermissions(map[string]interface{}{
				"collection_id":    canvasBranch.CanvasRepository.CollectionID,
				"permission_group": models.PGCollectionViewMetadataSysName,
				"member_id":        memberID,
			})
			for _, collectionPerm := range collectionPerms {
				err := App.Repo.Manger.HardDeleteByID(models.COLLECTION_PERMISSION, collectionPerm.ID)
				if err != nil {
					fmt.Println("Error on removing collection view metadata permission", err)
				}
			}
			App.Service.InvalidateCanvasPermissionCache(userID, studioID, canvasBranch.CanvasRepository.CollectionID)
			App.Service.InvalidateCollectionPermissionCache(userID, studioID)
		}
	}
}

/*
	RemoveRoleViewMetadataPermissionOnParents

	Removes the viewMetaData permissions on the parent canvas if user doesn't have any other permissions on that same
	level or under parent canvas level.
*/
func (s permissionService) RemoveRoleViewMetadataPermissionOnParents(canvasBranchID, roleID, studioID uint64) {
	canvasBranch, _ := App.Repo.GetCanvasBranch(map[string]interface{}{"id": canvasBranchID})
	if canvasBranch.CanvasRepository.ParentCanvasRepositoryID != nil {
		canvasBranchPerms, _ := App.Repo.GetCanvasBranchPerms(map[string]interface{}{"role_id": roleID, "cbp_parent_canvas_repository_id": canvasBranch.CanvasRepository.ParentCanvasRepositoryID})
		if len(canvasBranchPerms) == 0 {
			branchPerms, _ := App.Repo.GetCanvasBranchPerms(map[string]interface{}{
				"role_id":          roleID,
				"canvas_branch_id": canvasBranch.CanvasRepository.ParentCanvasRepository.DefaultBranchID,
				"permission_group": models.PGCanvasViewMetadataSysName,
			})
			for _, perm := range branchPerms {
				err := App.Repo.Manger.HardDeleteByID(models.CANVAS_BRANCH_PERMISSION, perm.ID)
				if err != nil {
					fmt.Println("Error on removing canvas view metadata permission", err)
				}
			}
			App.Service.InvalidateSubCanvasPermissionCacheByRole(roleID, studioID, canvasBranch.CanvasRepository.CollectionID, *canvasBranch.CanvasRepository.ParentCanvasRepositoryID)
		}
		s.RemoveRoleViewMetadataPermissionOnParents(*canvasBranch.CanvasRepository.ParentCanvasRepository.DefaultBranchID, roleID, studioID)
	} else {
		// Collection perms update
		canvasBranchPerms, _ := App.Repo.GetCanvasBranchPerms(map[string]interface{}{"role_id": roleID, "collection_id": canvasBranch.CanvasRepository.CollectionID})
		if len(canvasBranchPerms) == 0 {
			collectionPerms, _ := App.Repo.getCollectionPermissions(map[string]interface{}{
				"collection_id":    canvasBranch.CanvasRepository.CollectionID,
				"permission_group": models.PGCollectionViewMetadataSysName,
				"role_id":          roleID,
			})
			for _, collectionPerm := range collectionPerms {
				err := App.Repo.Manger.HardDeleteByID(models.COLLECTION_PERMISSION, collectionPerm.ID)
				if err != nil {
					fmt.Println("Error on removing collection view metadata permission", err)
				}
			}
			App.Service.InvalidateCanvasPermissionCacheByRole(roleID, studioID, canvasBranch.CanvasRepository.CollectionID)
			App.Service.InvalidateCollectionPermissionCacheByRole(roleID, studioID)
		}
	}
}

// Inherits the mainCanvasRepo permissions for the lanuage page repo.
func (s permissionService) InheritLanguageRepoParentPerms(languageRepo *models.CanvasRepository) {
	canvasRepo, _ := App.Repo.GetCanvasRepo(map[string]interface{}{"id": languageRepo.DefaultLanguageCanvasRepoID})
	// parent as canvas flow
	branchPerms, err := App.Repo.GetCanvasPermissionsByID(*canvasRepo.DefaultBranchID)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	for _, perm := range branchPerms {

		_, err = App.Repo.createCanvasBranchPermission(
			canvasRepo.CollectionID, perm.PermissionGroup, perm.MemberId, perm.IsOverridden,
			canvasRepo.StudioID, languageRepo.ID, *languageRepo.DefaultBranchID, nil, perm.RoleId)
		if err != nil {
			logger.Error(err.Error())
		}

		// Invalidating the user permissions
		if perm.RoleId != nil {
			if canvasRepo.ParentCanvasRepositoryID != nil {
				s.InvalidateSubCanvasPermissionCacheByRole(*perm.RoleId, canvasRepo.StudioID, canvasRepo.CollectionID, *canvasRepo.ParentCanvasRepositoryID)
			} else {
				s.InvalidateCanvasPermissionCacheByRole(*perm.RoleId, canvasRepo.StudioID, canvasRepo.CollectionID)
			}
		} else {
			if canvasRepo.ParentCanvasRepository != nil {
				s.InvalidateSubCanvasPermissionCache(perm.Member.UserID, canvasRepo.StudioID, canvasRepo.CollectionID, *canvasRepo.ParentCanvasRepositoryID)
			} else {
				s.InvalidateCanvasPermissionCache(perm.Member.UserID, canvasRepo.StudioID, canvasRepo.CollectionID)
			}
		}
	}
}

// If we make the canvas as public then we need to make the parent has_public_canvas key true if parent is private.
func (s permissionService) UpdateParentHasCanvasRepoOnPublic(canvasBranchID uint64) {
	canvasBranch, _ := App.Repo.GetCanvasBranch(map[string]interface{}{"id": canvasBranchID})
	if canvasBranch.CanvasRepository.ParentCanvasRepositoryID != nil {
		canvasRepo, _ := App.Repo.GetCanvasRepo(map[string]interface{}{"id": canvasBranch.CanvasRepository.ParentCanvasRepositoryID})
		if !canvasRepo.HasPublicCanvas {
			App.Repo.Manger.UpdateEntityByID(models.CANVAS_REPO, canvasRepo.ID, map[string]interface{}{"has_public_canvas": true})
		}
		s.UpdateParentHasCanvasRepoOnPublic(*canvasRepo.DefaultBranchID)
	} else {
		App.Repo.Manger.UpdateEntityByID(models.COLLECTION, canvasBranch.CanvasRepository.CollectionID, map[string]interface{}{"has_public_canvas": true})
	}
}

// If we make the canvas as private then we need to check the parent and make has_public_canvas to false if no other
// public canvases are present.
func (s permissionService) UpdateParentHasCanvasRepoOnPrivate(canvasBranchID uint64) {
	canvasBranch, _ := App.Repo.GetCanvasBranch(map[string]interface{}{"id": canvasBranchID})
	if canvasBranch.CanvasRepository.ParentCanvasRepositoryID != nil {
		canvasRepos, _ := App.Repo.GetCanvasRepos(map[string]interface{}{
			"parent_canvas_repository_id": canvasBranch.CanvasRepository.ParentCanvasRepositoryID,
			"has_public_canvas":           true,
		})
		allCanvasRepos, _ := App.Repo.GetCanvasRepos(map[string]interface{}{
			"parent_canvas_repository_id": canvasBranch.CanvasRepository.ParentCanvasRepositoryID,
		})
		canvasRepoIDs := []uint64{}
		for _, repo := range *allCanvasRepos {
			canvasRepoIDs = append(canvasRepoIDs, repo.ID)
		}
		canvasBranches, _ := App.Repo.GetCanvasBranches(map[string]interface{}{
			"canvas_repository_id": canvasRepoIDs,
			"public_access":        []string{models.VIEW, models.COMMENT, models.EDIT},
			"is_default":           true,
		})
		if len(*canvasRepos) == 0 && len(*canvasBranches) == 0 {
			App.Repo.Manger.UpdateEntityByID(models.CANVAS_REPO, *canvasBranch.CanvasRepository.ParentCanvasRepositoryID, map[string]interface{}{"has_public_canvas": false})
		}
		s.UpdateParentHasCanvasRepoOnPrivate(*canvasBranch.CanvasRepository.ParentCanvasRepository.DefaultBranchID)
	} else {
		canvasRepos, _ := App.Repo.GetCanvasRepos(map[string]interface{}{
			"collection_id":     canvasBranch.CanvasRepository.CollectionID,
			"has_public_canvas": true,
		})
		allCanvasRepos, _ := App.Repo.GetCanvasRepos(map[string]interface{}{
			"collection_id": canvasBranch.CanvasRepository.CollectionID,
		})
		canvasRepoIDs := []uint64{}
		for _, repo := range *allCanvasRepos {
			canvasRepoIDs = append(canvasRepoIDs, repo.ID)
		}
		canvasBranches, _ := App.Repo.GetCanvasBranches(map[string]interface{}{
			"canvas_repository_id": canvasRepoIDs,
			"public_access":        []string{models.VIEW, models.COMMENT, models.EDIT},
			"is_default":           true,
		})
		if len(*canvasRepos) == 0 && len(*canvasBranches) == 0 {
			App.Repo.Manger.UpdateEntityByID(models.COLLECTION, canvasBranch.CanvasRepository.CollectionID, map[string]interface{}{"has_public_canvas": false})
		}
	}
}

// If we make the canvas on remove then we need to check the parent and make has_public_canvas to false if no other
// public canvases are present.
func (s permissionService) UpdateParentHasCanvasRepoOnRemove(canvasBranchID uint64, removeBranchID uint64) {
	canvasBranch, _ := App.Repo.GetCanvasBranch(map[string]interface{}{"id": canvasBranchID})
	if canvasBranch.CanvasRepository.ParentCanvasRepositoryID != nil {
		canvasRepos := []models.CanvasRepository{}
		canvasReposData, _ := App.Repo.GetCanvasRepos(map[string]interface{}{
			"parent_canvas_repository_id": canvasBranch.CanvasRepository.ParentCanvasRepositoryID,
			"has_public_canvas":           true,
		})
		canvasRepos = *canvasReposData
		for i, repo := range canvasRepos {
			if *repo.DefaultBranchID == removeBranchID {
				canvasRepos = append(canvasRepos[:i], canvasRepos[i+1:]...)
			}
		}
		allCanvasRepos, _ := App.Repo.GetCanvasRepos(map[string]interface{}{
			"parent_canvas_repository_id": canvasBranch.CanvasRepository.ParentCanvasRepositoryID,
		})
		canvasRepoIDs := []uint64{}
		for _, repo := range *allCanvasRepos {
			if *repo.DefaultBranchID == removeBranchID {
				continue
			}
			canvasRepoIDs = append(canvasRepoIDs, repo.ID)
		}
		canvasBranches, _ := App.Repo.GetCanvasBranches(map[string]interface{}{
			"canvas_repository_id": canvasRepoIDs,
			"public_access":        []string{models.VIEW, models.COMMENT, models.EDIT},
			"is_default":           true,
		})
		fmt.Println(len(canvasRepos), len(*canvasBranches))
		if len(canvasRepos) == 0 && len(*canvasBranches) == 0 {
			App.Repo.Manger.UpdateEntityByID(models.CANVAS_REPO, *canvasBranch.CanvasRepository.ParentCanvasRepositoryID, map[string]interface{}{"has_public_canvas": false})
		}
		s.UpdateParentHasCanvasRepoOnRemove(*canvasBranch.CanvasRepository.ParentCanvasRepository.DefaultBranchID, canvasBranchID)
	} else {
		canvasRepos := []models.CanvasRepository{}
		canvasReposData, _ := App.Repo.GetCanvasRepos(map[string]interface{}{
			"collection_id":     canvasBranch.CanvasRepository.CollectionID,
			"has_public_canvas": true,
		})
		canvasRepos = *canvasReposData
		for i, repo := range canvasRepos {
			if *repo.DefaultBranchID == removeBranchID {
				canvasRepos = append(canvasRepos[:i], canvasRepos[i+1:]...)
			}
		}
		allCanvasRepos, _ := App.Repo.GetCanvasRepos(map[string]interface{}{
			"collection_id": canvasBranch.CanvasRepository.CollectionID,
		})
		canvasRepoIDs := []uint64{}
		for _, repo := range *allCanvasRepos {
			if *repo.DefaultBranchID == removeBranchID {
				fmt.Println("skipping this canavs", repo.ID, repo.Name)
				continue
			}
			canvasRepoIDs = append(canvasRepoIDs, repo.ID)
		}
		canvasBranches, _ := App.Repo.GetCanvasBranches(map[string]interface{}{
			"canvas_repository_id": canvasRepoIDs,
			"public_access":        []string{models.VIEW, models.COMMENT, models.EDIT},
			"is_default":           true,
		})
		fmt.Println(len(canvasRepos), len(*canvasBranches))
		if len(canvasRepos) == 0 && len(*canvasBranches) == 0 {
			App.Repo.Manger.UpdateEntityByID(models.COLLECTION, canvasBranch.CanvasRepository.CollectionID, map[string]interface{}{"has_public_canvas": false})
		}
	}
}
