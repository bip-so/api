package permissions

import "gitlab.com/phonepost/bip-be-platform/internal/models"

type PermissionList struct {
	ID         uint64 `json:"id"`
	Permission string `json:"permission"`
}

func SerializeStudioPermission(studioPermissions *models.StudioPermission) *PermissionList {
	return &PermissionList{
		ID:         studioPermissions.StudioID,
		Permission: studioPermissions.PermissionGroup,
	}
}

func SerializeCollectionPermission(collPermissions *models.CollectionPermission) *PermissionList {
	return &PermissionList{
		ID:         collPermissions.CollectionId,
		Permission: collPermissions.PermissionGroup,
	}
}

type InvalidatePermissions struct {
	MemberID       *uint64 `json:"memberID"`
	RoleID         *uint64 `json:"roleID"`
	UserID         uint64  `json:"userID"`
	StudioID       uint64  `json:"studioID"`
	CollectionID   uint64  `json:"CollectionID"`
	ParentCanvasID *uint64 `json:"parentCanvasID"`
	InvalidationOn string  `json:"invalidationOn"`
}

type MemberUserStudio struct {
	UserID   uint64
	StudioID uint64
}
