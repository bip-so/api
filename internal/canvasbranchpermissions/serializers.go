package canvasbranchpermissions

import (
	"gitlab.com/phonepost/bip-be-platform/internal/canvasrepo"
	"gitlab.com/phonepost/bip-be-platform/internal/collection"
	"gitlab.com/phonepost/bip-be-platform/internal/member"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/role"
	"gitlab.com/phonepost/bip-be-platform/internal/studio"
)

type CanvasBranchPermissionSerializer struct {
	ID                       uint64                                 `json:"id"`
	UUID                     string                                 `json:"uuid"`
	StudioID                 uint64                                 `json:"studioID"`
	CollectionID             uint64                                 `json:"collectionID"`
	CanvasRepositoryID       uint64                                 `json:"canvasRepositoryID"`
	CanvasBranchID           *uint64                                `json:"canvasBranchID"`
	ParentCanvasRepositoryID *uint64                                `json:"parentCanvasRepositoryID"`
	PermissionGroup          string                                 `json:"permissionGroup"`
	RoleID                   *uint64                                `json:"roleID"`
	MemberID                 *uint64                                `json:"memberID"`
	IsOverridden             bool                                   `json:"isOverridden"`
	Role                     *role.RoleSerializer                   `json:"role"`
	Member                   *member.MemberSerializer               `json:"member"`
	Studio                   *studio.StudioSerializer               `json:"studio"`
	Collection               collection.CollectionSerializer        `json:"collection"`
	CanvasRepository         canvasrepo.CanvasRepoDefaultSerializer `json:"canvasRepository"`
	ParentCanvasRepository   canvasrepo.CanvasRepoDefaultSerializer `json:"parentCanvasRepository"`
	Branch                   *canvasrepo.CanvasBranchMiniSerializer `json:"canvasBranch"`
}

func SerializeCanvasBranchPermissionsPermission(collectionperm *models.CanvasBranchPermission) *CanvasBranchPermissionSerializer {
	view := CanvasBranchPermissionSerializer{
		ID:                       collectionperm.ID,
		UUID:                     collectionperm.UUID.String(),
		CollectionID:             collectionperm.CollectionId,
		CanvasRepositoryID:       collectionperm.CanvasRepositoryID,
		CanvasBranchID:           collectionperm.CanvasBranchID,
		ParentCanvasRepositoryID: collectionperm.CbpParentCanvasRepositoryID,
		StudioID:                 collectionperm.StudioID,
		PermissionGroup:          collectionperm.PermissionGroup,
		RoleID:                   collectionperm.RoleId,
		MemberID:                 collectionperm.MemberId,
		IsOverridden:             collectionperm.IsOverridden,
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
	if collectionperm.CanvasBranch != nil {
		view.Branch = canvasrepo.SerializeCanvasBranchMini(collectionperm.CanvasBranch)
	}
	if collectionperm.CanvasRepository != nil {
		view.CanvasRepository = *canvasrepo.SerializeDefaultCanvasRepo(collectionperm.CanvasRepository)
	}
	if collectionperm.CbpParentCanvasRepository != nil {
		view.ParentCanvasRepository = *canvasrepo.SerializeDefaultCanvasRepo(collectionperm.CbpParentCanvasRepository)
	}

	view.Collection = collection.CollectionSerializerData(&collectionperm.Collection)

	return &view
}
