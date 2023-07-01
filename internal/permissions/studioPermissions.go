package permissions

import (
	"context"
	"encoding/json"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

/*
	CalculateStudioPermissions: Calculates the permissions of the studios that user has joined.

	Algorithm:
		- Get all the member instances based on userId.
		- Get roleIds & memberIds by looping members.
		- Here user can be part of multiple roles or a member with permission & each role will have different permissions.
          We need to get the max. weight permission for that user & studio.
		- Get studioPermissions based on memberIds & roleIds
		- Now we loop studioPermissions to find the max. Weight permission of a studio.

	Cache:
		- We use REDIS cache here. So first we check for the cache data in REDIS, if data is present we return the data
          or, we calculate the data and SET the data to REDIS and return the data.

	Args:
		userid uint64
	Returns:
		permissionList map[uint64]string
*/
func (s permissionService) CalculateStudioPermissions(userID uint64) (map[uint64]string, error) {
	permissionList := make(map[uint64]string)
	isOveriddenList := make(map[uint64]bool)

	// Getting the data from redis
	ctx := context.Background()
	redisKey := s.StudioPermissionsRedisKey(userID)
	value := s.cache.HGet(ctx, PermissionsHash+utils.String(userID), redisKey)
	err := json.Unmarshal([]byte(value), &permissionList)
	if err == nil {
		return permissionList, nil
	}

	members, err := App.Repo.getAllStudiosUserMemberOf(userID)
	if err != nil {
		return nil, err
	}

	var memberIds []uint64
	var roleIds []uint64
	for _, memb := range members {
		memberIds = append(memberIds, memb.ID)
		for _, role := range memb.Roles {
			roleIds = append(roleIds, role.ID)
		}
	}

	studioPerms, err := App.Repo.getStudioPermissionsByMemberIds(memberIds, roleIds)
	if err != nil {
		return nil, err
	}

	for _, perm := range studioPerms {

		if isOveriddenList[perm.StudioID] || (perm.MemberId != nil && perm.IsOverridden) {

			permissionList[perm.StudioID] = perm.PermissionGroup
			isOveriddenList[perm.StudioID] = true

		} else {

			var currPermissionWt = models.PermissionGroupWeightMap[permissiongroup.PGTYPESTUDIO][perm.PermissionGroup]
			var storedPermissionWt = models.PermissionGroupWeightMap[permissiongroup.PGTYPESTUDIO][permissionList[perm.StudioID]]

			if currPermissionWt > storedPermissionWt {
				permissionList[perm.StudioID] = perm.PermissionGroup
			}
		}
	}

	// Set the data to redis
	permissionListString, _ := json.Marshal(permissionList)
	s.cache.HSet(ctx, PermissionsHash+utils.String(userID), redisKey, permissionListString)

	return permissionList, nil
}

// StudioPermissionsRedisKey return the redis key by userID
func (s permissionService) StudioPermissionsRedisKey(userID uint64) string {
	return StudioPermissionRedisKey + utils.String(userID)
}

// StudioPermissionGroupByUserID return studios permissions by userID.
func (s permissionService) StudioPermissionGroupByUserID(userID uint64, studioID uint64) (*string, error) {
	var permissionGroup string
	var permissionsList map[uint64]string

	permissionsList, err := s.CalculateStudioPermissions(userID)
	if err != nil {
		return nil, err
	}

	permissionGroup = permissionsList[studioID]
	return &permissionGroup, nil
}

// CanUserDoThisOnStudio return the bool of whether user has a specific permission on this studio or not.
func (s permissionService) CanUserDoThisOnStudio(userID uint64, studioID uint64, permissionGroup string) (bool, error) {
	pg, err := s.StudioPermissionGroupByUserID(userID, studioID)
	if err != nil {
		return false, err
	}
	canDo := models.StudioPermissionsMap[*pg][permissionGroup]
	if canDo != 0 {
		return true, nil
	} else {
		return false, err
	}
}
