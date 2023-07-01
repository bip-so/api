package permissiongroup

import "gitlab.com/phonepost/bip-be-platform/internal/models"

/*
Important Note If you do aanything here make sure to :
Please update : `internal/models/permisions.go` for any change here
*/

func PGCollectionNone() PermissionsTemplate {
	pg := PermissionsTemplate{}
	pg.DisplayName = "None"
	pg.SystemName = models.PGCollectionNoneSysName
	pg.StudioID = 0
	pg.Type = PGTYPECOLLECTIION
	pg.Weight = 0
	pg.BetterPermissions = models.CollectionPermissionsMap[models.PGCollectionNoneSysName]
	pg.Permissions = []PermissionObject{
		{Key: COLLECTION_MEMBERSHIP_MANAGE, Value: 0},
		{Key: COLLECTION_PUBLIC_ACCESS_CHANGE, Value: 0},
		{Key: COLLECTION_MANAGE_PERMS, Value: 0},
		{Key: COLLECTION_DELETE, Value: 0},
		{Key: COLLECTION_OVERRIDE_STUDIO_MODE_ROLE, Value: 0},
		{Key: COLLECTION_MANAGE_PUBLISH_REQUEST, Value: 0},
		{Key: COLLECTION_EDIT_NAME, Value: 0},
		{Key: COLLECTION_VIEW_METADATA, Value: 0},
	}
	return pg
}

func PGCollectionViewMetadata() PermissionsTemplate {
	pg := PermissionsTemplate{}
	pg.DisplayName = "ViewMetadata"
	pg.SystemName = models.PGCollectionViewMetadataSysName
	pg.StudioID = 0
	pg.Type = PGTYPECOLLECTIION
	pg.Weight = 10
	pg.BetterPermissions = models.CollectionPermissionsMap[models.PGCollectionViewMetadataSysName]
	pg.Permissions = []PermissionObject{
		{Key: COLLECTION_MEMBERSHIP_MANAGE, Value: 0},
		{Key: COLLECTION_PUBLIC_ACCESS_CHANGE, Value: 0},
		{Key: COLLECTION_MANAGE_PERMS, Value: 0},
		{Key: COLLECTION_DELETE, Value: 0},
		{Key: COLLECTION_OVERRIDE_STUDIO_MODE_ROLE, Value: 0},
		{Key: COLLECTION_MANAGE_PUBLISH_REQUEST, Value: 0},
		{Key: COLLECTION_EDIT_NAME, Value: 0},
		{Key: COLLECTION_VIEW_METADATA, Value: 1},
	}
	return pg
}

// View
func PGCollectionView() PermissionsTemplate {
	pg := PermissionsTemplate{}
	pg.DisplayName = "View"
	pg.SystemName = models.PGCollectionViewSysName
	pg.StudioID = 0
	pg.Type = PGTYPECOLLECTIION
	pg.Weight = 20
	pg.BetterPermissions = models.CollectionPermissionsMap[models.PGCollectionViewSysName]

	pg.Permissions = []PermissionObject{
		{Key: COLLECTION_MEMBERSHIP_MANAGE, Value: 0},
		{Key: COLLECTION_PUBLIC_ACCESS_CHANGE, Value: 1},
		{Key: COLLECTION_MANAGE_PERMS, Value: 0},
		{Key: COLLECTION_DELETE, Value: 0},
		{Key: COLLECTION_OVERRIDE_STUDIO_MODE_ROLE, Value: 0},
		{Key: COLLECTION_MANAGE_PUBLISH_REQUEST, Value: 0},
		{Key: COLLECTION_EDIT_NAME, Value: 0},
		{Key: COLLECTION_VIEW_METADATA, Value: 1},
	}
	return pg
}

// Comment
func PGCollectionComment() PermissionsTemplate {
	pg := PermissionsTemplate{}
	pg.DisplayName = "Comment"
	pg.SystemName = models.PGCollectionCommentSysName
	pg.StudioID = 0
	pg.Type = PGTYPECOLLECTIION
	pg.Weight = 30
	pg.BetterPermissions = models.CollectionPermissionsMap[models.PGCollectionCommentSysName]
	pg.Permissions = []PermissionObject{
		{Key: COLLECTION_MEMBERSHIP_MANAGE, Value: 0},
		{Key: COLLECTION_PUBLIC_ACCESS_CHANGE, Value: 1},
		{Key: COLLECTION_MANAGE_PERMS, Value: 0},
		{Key: COLLECTION_DELETE, Value: 0},
		{Key: COLLECTION_OVERRIDE_STUDIO_MODE_ROLE, Value: 1},
		{Key: COLLECTION_MANAGE_PUBLISH_REQUEST, Value: 0},
		{Key: COLLECTION_EDIT_NAME, Value: 0},
		{Key: COLLECTION_VIEW_METADATA, Value: 1},
	}
	return pg
}

// Edit
func PGCollectionEdit() PermissionsTemplate {
	pg := PermissionsTemplate{}
	pg.DisplayName = "Edit"
	pg.StudioID = 0
	pg.Type = PGTYPECOLLECTIION
	pg.Weight = 40
	pg.SystemName = models.PGCollectionEditSysName
	pg.BetterPermissions = models.CollectionPermissionsMap[models.PGCollectionEditSysName]
	pg.Permissions = []PermissionObject{
		{Key: COLLECTION_MEMBERSHIP_MANAGE, Value: 0},
		{Key: COLLECTION_PUBLIC_ACCESS_CHANGE, Value: 1},
		{Key: COLLECTION_MANAGE_PERMS, Value: 0},
		{Key: COLLECTION_DELETE, Value: 0},
		{Key: COLLECTION_OVERRIDE_STUDIO_MODE_ROLE, Value: 1},
		{Key: COLLECTION_MANAGE_PUBLISH_REQUEST, Value: 0},
		{Key: COLLECTION_EDIT_NAME, Value: 0},
		{Key: COLLECTION_VIEW_METADATA, Value: 1},
	}
	return pg
}

// None
func PGCollectionModerate() PermissionsTemplate {
	pg := PermissionsTemplate{}
	pg.DisplayName = "Moderate"
	pg.SystemName = models.PGCollectionModerateSysName
	pg.StudioID = 0
	pg.Type = PGTYPECOLLECTIION
	pg.Weight = 1000
	pg.BetterPermissions = models.CollectionPermissionsMap[models.PGCollectionModerateSysName]
	pg.Permissions = []PermissionObject{
		{Key: COLLECTION_MEMBERSHIP_MANAGE, Value: 1},
		{Key: COLLECTION_PUBLIC_ACCESS_CHANGE, Value: 1},
		{Key: COLLECTION_MANAGE_PERMS, Value: 1},
		{Key: COLLECTION_DELETE, Value: 1},
		{Key: COLLECTION_OVERRIDE_STUDIO_MODE_ROLE, Value: 1},
		{Key: COLLECTION_MANAGE_PUBLISH_REQUEST, Value: 1},
		{Key: COLLECTION_EDIT_NAME, Value: 1},
		{Key: COLLECTION_VIEW_METADATA, Value: 1},
	}
	return pg
}
