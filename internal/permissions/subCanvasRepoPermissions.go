package permissions

import (
	"context"
	"encoding/json"
	"fmt"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

/*
	CalculateSubCanvasRepoPermissions: Calculates the permissions of the user on canvasBranches which are under another canvas.

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
func (s permissionService) CalculateSubCanvasRepoPermissions(userID, studioID, collectionID, canvasID uint64) (map[uint64]map[uint64]string, error) {
	permissionList := make(map[uint64]map[uint64]string)
	isOverriddenList := make(map[uint64]bool)

	// Getting the data from redis
	ctx := context.Background()
	redisKey := s.SubCanvasPermissionsRedisKey(userID, studioID, collectionID, canvasID)
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
	fmt.Println(memberIds, roleIds, studioID, collectionID, canvasID)
	subCanvasRepoPerms, err := App.Repo.getCanvasBranchPermissionsByMemberIds(memberIds, roleIds, studioID, collectionID, canvasID)
	if err != nil {
		return nil, err
	}
	fmt.Println("subcanvasRepoPerms===>", subCanvasRepoPerms)
	for _, perm := range subCanvasRepoPerms {
		fmt.Println(perm.ID, perm.CanvasBranchID)
		if len(permissionList[perm.CanvasRepositoryID]) == 0 {
			permissionList[perm.CanvasRepositoryID] = make(map[uint64]string)
		}
		if isOverriddenList[*perm.CanvasBranchID] || (perm.MemberId != nil && perm.IsOverridden) {
			permissionList[perm.CanvasRepositoryID][*perm.CanvasBranchID] = perm.PermissionGroup
			isOverriddenList[*perm.CanvasBranchID] = true
		} else {
			var currPermissionWt = models.PermissionGroupWeightMap[permissiongroup.PGTYPECANVAS][perm.PermissionGroup]
			var storedPermissionWt = models.PermissionGroupWeightMap[permissiongroup.PGTYPECANVAS][permissionList[perm.CanvasRepositoryID][*perm.CanvasBranchID]]
			fmt.Println("ar 90", perm.ID, perm.CanvasBranchID, currPermissionWt, storedPermissionWt, permissionList)
			if currPermissionWt > storedPermissionWt {
				permissionList[perm.CanvasRepositoryID][*perm.CanvasBranchID] = perm.PermissionGroup
			}
		}
		fmt.Println("at llast", perm.ID, perm.CanvasBranchID, permissionList)
	}

	// Set the data to redis
	permissionListString, _ := json.Marshal(permissionList)
	s.cache.HSet(ctx, PermissionsHash+utils.String(userID), redisKey, permissionListString)

	return permissionList, nil
}

func (s permissionService) SubCanvasPermissionsRedisKey(userID uint64, studioID uint64, collectionID uint64, canvasID uint64) string {
	return CanvasPermissionRedisKey + utils.String(userID) + ":" + utils.String(studioID) + ":" + utils.String(collectionID) + ":" + utils.String(canvasID)
}
