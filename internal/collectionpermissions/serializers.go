package collectionpermissions

import (
	"gitlab.com/phonepost/bip-be-platform/internal/collection"
	"gitlab.com/phonepost/bip-be-platform/internal/member"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/role"
	"gitlab.com/phonepost/bip-be-platform/internal/studio"
)

type CollectionPermissionsSerializer struct {
	ID              uint64                          `json:"id"`
	UUID            string                          `json:"uuid"`
	StudioID        uint64                          `json:"studioID"`
	CollectionID    uint64                          `json:"collectionID"`
	PermissionGroup string                          `json:"permissionGroup"`
	RoleID          *uint64                         `json:"roleID"`
	MemberID        *uint64                         `json:"memberID"`
	IsOverridden    bool                            `json:"isOverridden"`
	Role            *role.RoleSerializer            `json:"role"`
	Collection      collection.CollectionSerializer `json:"collection"`
	Member          *member.MemberSerializer        `json:"member"`
	Studio          *studio.StudioSerializer        `json:"studio"`
}

func SerializeCollectionPermission(collectionperm *models.CollectionPermission) *CollectionPermissionsSerializer {
	view := CollectionPermissionsSerializer{
		ID:              collectionperm.ID,
		UUID:            collectionperm.UUID.String(),
		StudioID:        collectionperm.StudioID,
		CollectionID:    collectionperm.CollectionId,
		PermissionGroup: collectionperm.PermissionGroup,
		RoleID:          collectionperm.RoleId,
		MemberID:        collectionperm.MemberId,
		IsOverridden:    collectionperm.IsOverridden,
	}
	if collectionperm.Member != nil {
		view.Member = member.SerializeMember(collectionperm.Member)
	}
	if collectionperm.Role != nil {
		view.Role = role.SerializeRole(collectionperm.Role)
	}
	if collectionperm.Studio != nil {
		view.Studio = studio.SerializeStudio(collectionperm.Studio)
	}

	view.Collection = collection.CollectionSerializerData(&collectionperm.Collection)

	return &view
}

// post canvas perms
type newCanvasBranchPermissionCreatePost struct {
	CollectionId                uint64 `json:"collectionId" binding:"required"`
	CanvasBranchId              uint64 `json:"canvasBranchId" binding:"required"`
	CanvasRepositoryID          uint64 `json:"canvasRepositoryId" binding:"required"`
	CbpParentCanvasRepositoryID uint64 `json:"parentCanvasRepositoryId"`
	PermGroup                   string `json:"permGroup" binding:"required"`
	RoleID                      uint64 `json:"roleID"`
	MemberID                    uint64 `json:"memberID"`
	UserID                      uint64 `json:"userID"`
	IsOverridden                bool   `json:"isOverridden"`
}
