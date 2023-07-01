package permissions

import (
	"context"
	"encoding/json"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

/*
	CalculateCollectionPermissions: Calculates the permissions of the user on collections for a specific studio.

	Algorithm:
		- Get all the members of the studio based on userId & studioId
		- Get roleIds & memberIds by looping members.
		- Here user can be part of multiple roles or a member with permission & each role will have different permissions.
          We need to get the max. weight permission for that user & collection.
		- Get collectionPermissions based on memberIds & roleIds
		- Now we loop collectionPermissions to find the max. Weight permission of a collection.

	Cache:
		- We use REDIS cache here. So first we check for the cache data in REDIS, if data is present we return the data
          or, we calculate the data and SET the data to REDIS and return the data.

	Args:
		userid uint64
		studioID uint64
	Returns:
		permissionList map[uint64]string
*/
func (s permissionService) CalculateCollectionPermissions(userid, studioID uint64) (map[uint64]string, error) {
	permissionList := make(map[uint64]string)
	isOveriddenList := make(map[uint64]bool)

	// Getting the data from redis
	ctx := context.Background()
	redisKey := s.CollectionPermissionsRedisKey(userid, studioID)
	value := s.cache.HGet(ctx, PermissionsHash+utils.String(userid), redisKey)
	err := json.Unmarshal([]byte(value), &permissionList)
	if err == nil {
		return permissionList, nil
	}

	members, err := App.Repo.getMembersByUserIDs([]uint64{userid}, studioID)
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

	// find collection perm where memberId in members.IDS
	collectionPerms, err := App.Repo.getCollectionPermissionsForUser(memberIds, roleIds)
	if err != nil {
		return nil, err
	}

	for _, perm := range collectionPerms {
		if isOveriddenList[perm.CollectionId] || (perm.MemberId != nil && perm.IsOverridden) {

			isOveriddenList[perm.CollectionId] = true
			permissionList[perm.CollectionId] = perm.PermissionGroup

		} else {

			var currPermissionWt = models.PermissionGroupWeightMap[permissiongroup.PGTYPECOLLECTIION][perm.PermissionGroup]
			var storedPermissionWt = models.PermissionGroupWeightMap[permissiongroup.PGTYPECOLLECTIION][permissionList[perm.CollectionId]]

			if currPermissionWt > storedPermissionWt {
				permissionList[perm.CollectionId] = perm.PermissionGroup
			}

		}
	}

	// Set the data to redis
	permissionListString, _ := json.Marshal(permissionList)
	s.cache.HSet(ctx, PermissionsHash+utils.String(userid), redisKey, permissionListString)

	return permissionList, nil
}

// CollectionPermissionsRedisKey returns the collection permissions redis key
func (s permissionService) CollectionPermissionsRedisKey(userID uint64, studioID uint64) string {
	return CollectionPermissionRedisKey + utils.String(userID) + ":" + utils.String(studioID)
}

// CollectionPermissionGroupByUserID returns the permissions of collections present in the studio for a user
func (s permissionService) CollectionPermissionGroupByUserID(userID uint64, studioID uint64, collectionID uint64) (*string, error) {
	var permissionGroup string
	var permissionsList map[uint64]string

	permissionsList, err := s.CalculateCollectionPermissions(userID, studioID)
	if err != nil {
		return nil, err
	}

	permissionGroup = permissionsList[collectionID]
	return &permissionGroup, nil
}

// If you want to check if a User can do something on Collection, Pass UserID, BranchID, Permission to Check.
// This returns True/False or Error
func (s permissionService) CanUserDoThisOnCollection(userID uint64, studioID uint64, collectionID uint64, permissionGroup string) (bool, error) {
	collection, _ := App.Repo.GetCollection(map[string]interface{}{"id": collectionID})
	if userID == 0 && collection.PublicAccess == models.PRIVATE {
		return false, nil
	}

	var collectionAccess int
	if collection.PublicAccess != models.PRIVATE {
		if collection.PublicAccess == models.EDIT {
			collectionAccess = models.CollectionPermissionsMap[models.PGCollectionEditSysName][permissionGroup]
		} else if collection.PublicAccess == models.VIEW {
			collectionAccess = models.CollectionPermissionsMap[models.PGCollectionViewSysName][permissionGroup]
		} else if collection.PublicAccess == models.COMMENT {
			collectionAccess = models.CollectionPermissionsMap[models.PGCollectionCommentSysName][permissionGroup]
		}
	}

	if collectionAccess != 0 {
		return true, nil
	}
	pg, err := s.CollectionPermissionGroupByUserID(userID, studioID, collectionID)
	if err != nil {
		return false, err
	}
	canDo := models.CollectionPermissionsMap[*pg][permissionGroup]
	if canDo != 0 {
		return true, nil
	} else {
		return false, err
	}
}

/*
	AddMemberToCollectionIfNotPresent
	Args:
		authUserId, memberID, collectionID, studioID uint64
	Create the viewMetaData collectionPermission if user has no permission on the collection.
	We are also checking if user is having permission via role and if not present then only we are creating viewMetaData
	permission.
*/
func (s permissionService) AddMemberToCollectionIfNotPresent(authUserId, memberID, collectionID, studioID uint64) {
	collectionPerms, _ := App.Repo.getCollectionPermissions(map[string]interface{}{"member_id": memberID, "collection_id": collectionID})
	if len(collectionPerms) == 0 {
		// roles of the member
		roles, _ := App.Repo.GetMemberRolesByID(studioID, memberID)
		memberFound := false
		if len(roles) > 0 {
			// Get all collectionPerms
			collectionPerms, _ = App.Repo.getCollectionPermissions(map[string]interface{}{"collection_id": collectionID})
			for _, perm := range collectionPerms {
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
			App.Repo.UpdateCollectionPermissions(map[string]interface{}{"collection_id": collectionID, "member_id": memberID}, collectionID, studioID, memberID, 0, false, models.PGCollectionViewMetadataSysName, authUserId)
		}
		member, _ := App.Repo.GetMember(map[string]interface{}{"id": memberID})
		App.Service.InvalidateCollectionPermissionCache(member.UserID, studioID)
	}
}

// AddRoleToCollectionIfNotPresent Adding role to the collection only if it is not already present.
func (s permissionService) AddRoleToCollectionIfNotPresent(authUserId, roleID, collectionID, studioID uint64) {
	collectionPerms, _ := App.Repo.getCollectionPermissions(map[string]interface{}{"role_id": roleID, "collection_id": collectionID})
	if len(collectionPerms) == 0 {
		App.Repo.UpdateCollectionPermissions(map[string]interface{}{"collection_id": collectionID, "role_id": roleID}, collectionID, studioID, 0, roleID, false, models.PGCollectionViewMetadataSysName, authUserId)
		App.Service.InvalidateCollectionPermissionCacheByRole(roleID, studioID)
	}
}
