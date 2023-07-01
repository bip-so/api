package permissions

import (
	"context"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

func (s permissionService) InvalidatePermissionCache(ctx context.Context, hash string, key string) {
	if key != "" {
		s.cache.HDelete(ctx, hash, key)
	} else {
		s.cache.Delete(ctx, hash)
	}
}

// InvalidateUserPermissionCache invalidates all the data for that user & studio
func (s permissionService) InvalidateUserPermissionCache(userID uint64, studioID uint64) {
	userIDStr := utils.String(userID)
	studioIDStr := utils.String(studioID)
	s.cache.HDeleteMatching(context.Background(), PermissionsHash+utils.String(userID), StudioPermissionRedisKey+userIDStr+"*")
	s.cache.HDeleteMatching(context.Background(), PermissionsHash+utils.String(userID), CollectionPermissionRedisKey+userIDStr+":"+studioIDStr+"*")
	s.cache.HDeleteMatching(context.Background(), PermissionsHash+utils.String(userID), CanvasPermissionRedisKey+userIDStr+":"+studioIDStr+"*")
}

// InvalidateUserPermissionCache invalidates all the data for that user & studio
func (s permissionService) InvalidateRolePermissionCache(roleID uint64, studioID uint64) {
	role, _ := App.Repo.GetRoleByID(roleID)
	for _, member := range role.Members {
		userIDStr := utils.String(member.UserID)
		studioIDStr := utils.String(studioID)
		s.cache.HDeleteMatching(context.Background(), PermissionsHash+utils.String(member.UserID), StudioPermissionRedisKey+userIDStr+"*")
		s.cache.HDeleteMatching(context.Background(), PermissionsHash+utils.String(member.UserID), CollectionPermissionRedisKey+userIDStr+":"+studioIDStr+"*")
		s.cache.HDeleteMatching(context.Background(), PermissionsHash+utils.String(member.UserID), CanvasPermissionRedisKey+userIDStr+":"+studioIDStr+"*")
	}
}

// InvalidateCollectionMatchingPermissionCache invalidates all the data for that user & studio
func (s permissionService) InvalidateCollectionMatchingPermissionCache(userID uint64, studioID uint64, collectionID uint64) {
	userIDStr := utils.String(userID)
	studioIDStr := utils.String(studioID)
	collectionIDStr := utils.String(collectionID)
	s.cache.HDeleteMatching(context.Background(), PermissionsHash+utils.String(userID), CanvasPermissionRedisKey+userIDStr+":"+studioIDStr+":"+collectionIDStr+"*")
}

// InvalidateCollectionMatchingPermissionCacheByRole invalidates all the data for that user & studio
func (s permissionService) InvalidateCollectionMatchingPermissionCacheByRole(roleID uint64, studioID uint64, collectionID uint64) {
	role, err := App.Repo.GetRoleByID(roleID)
	if err != nil {
		logger.Error(err.Error())
	}
	studioIDStr := utils.String(studioID)
	collectionIDStr := utils.String(collectionID)
	for _, member := range role.Members {
		userIDStr := utils.String(member.UserID)
		s.cache.HDeleteMatching(context.Background(), PermissionsHash+userIDStr, CanvasPermissionRedisKey+userIDStr+":"+studioIDStr+":"+collectionIDStr+"*")
	}
}

// InvalidateStudioPermissionCache invalidates only the studio permissions cache
func (s permissionService) InvalidateStudioPermissionCache(userID uint64) {
	redisKey := s.StudioPermissionsRedisKey(userID)
	s.cache.HDelete(context.Background(), PermissionsHash+utils.String(userID), redisKey)
}

func (s permissionService) InvalidateCollectionPermissionCacheByRole(roleID uint64, studioID uint64) {
	role, err := App.Repo.GetRoleByID(roleID)
	if err != nil {
		logger.Error(err.Error())
	}
	for _, member := range role.Members {
		redisKey := s.CollectionPermissionsRedisKey(member.UserID, studioID)
		s.cache.HDelete(context.Background(), PermissionsHash+utils.String(member.UserID), redisKey)
	}
}

func (s permissionService) InvalidateCanvasPermissionCacheByRole(roleID uint64, studioID uint64, collectionID uint64) {
	role, err := App.Repo.GetRoleByID(roleID)
	if err != nil {
		logger.Error(err.Error())
	}
	for _, member := range role.Members {
		redisKey := s.CanvasPermissionsRedisKey(member.UserID, studioID, collectionID)
		s.cache.HDelete(context.Background(), PermissionsHash+utils.String(member.UserID), redisKey)
	}
}

func (s permissionService) InvalidateSubCanvasPermissionCacheByRole(roleID uint64, studioID uint64, collectionID uint64, canvasID uint64) {
	role, err := App.Repo.GetRoleByID(roleID)
	if err != nil {
		logger.Error(err.Error())
	}
	for _, member := range role.Members {
		redisKey := s.SubCanvasPermissionsRedisKey(member.UserID, studioID, collectionID, canvasID)
		s.cache.HDelete(context.Background(), PermissionsHash+utils.String(member.UserID), redisKey)
	}
}

func (s permissionService) InvalidateCollectionPermissionCache(userID uint64, studioID uint64) {
	redisKey := s.CollectionPermissionsRedisKey(userID, studioID)
	s.cache.HDelete(context.Background(), PermissionsHash+utils.String(userID), redisKey)
}

func (s permissionService) InvalidateCanvasPermissionCache(userID uint64, studioID uint64, collectionID uint64) {
	redisKey := s.CanvasPermissionsRedisKey(userID, studioID, collectionID)
	s.cache.HDelete(context.Background(), PermissionsHash+utils.String(userID), redisKey)
}

func (s permissionService) InvalidateSubCanvasPermissionCache(userID uint64, studioID uint64, collectionID uint64, canvasRepositoryID uint64) {
	redisKey := s.SubCanvasPermissionsRedisKey(userID, studioID, collectionID, canvasRepositoryID)
	s.cache.HDelete(context.Background(), PermissionsHash+utils.String(userID), redisKey)
}

// InvalidatePermissions is a generic method to invalidate the cache permissions.
func (s permissionService) InvalidatePermissions(data *InvalidatePermissions) error {
	if data.MemberID != nil && *data.MemberID != 0 {
		member, err := App.Repo.GetMember(map[string]interface{}{"id": data.MemberID})
		if err != nil {
			return err
		}
		if data.InvalidationOn == "studio" {
			s.InvalidateStudioPermissionCache(member.UserID)
		} else if data.InvalidationOn == "collection" {
			s.InvalidateCollectionPermissionCache(member.UserID, member.StudioID)
		} else if data.InvalidationOn == "canvas" {
			s.InvalidateCanvasPermissionCache(member.UserID, member.StudioID, data.CollectionID)
		} else if data.InvalidationOn == "subCanvas" {
			s.InvalidateSubCanvasPermissionCache(member.UserID, member.StudioID, data.CollectionID, *data.ParentCanvasID)
		}
	} else {
		role, err := App.Repo.GetRoleByID(*data.RoleID)
		var members []MemberUserStudio
		if err != nil {
			return err
		}
		for _, member := range role.Members {
			memberUserData := MemberUserStudio{
				UserID:   member.UserID,
				StudioID: member.StudioID,
			}
			members = append(members, memberUserData)
		}
		for _, member := range members {
			if data.InvalidationOn == "studio" {
				s.InvalidateStudioPermissionCache(member.UserID)
			} else if data.InvalidationOn == "collection" {
				s.InvalidateCollectionPermissionCache(member.UserID, member.StudioID)
			} else if data.InvalidationOn == "canvas" {
				s.InvalidateCanvasPermissionCache(member.UserID, member.StudioID, data.CollectionID)
			} else if data.InvalidationOn == "subCanvas" {
				s.InvalidateSubCanvasPermissionCache(member.UserID, member.StudioID, data.CollectionID, *data.ParentCanvasID)
			}
		}
	}
	return nil
}
