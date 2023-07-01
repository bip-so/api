package permissiongroup

import "gitlab.com/phonepost/bip-be-platform/internal/models"

/*
Important Note If you do aanything here make sure to :
Please update : `internal/models/permisions.go` for any change here
*/
// Studio: None
func PGStudioNone() PermissionsTemplate {
	pg := PermissionsTemplate{}
	pg.DisplayName = "None"
	pg.SystemName = "pg_studio_none"
	pg.StudioID = 0
	pg.Type = PGTYPESTUDIO
	pg.Weight = 0
	pg.BetterPermissions = models.StudioPermissionsMap[models.PGStudioNoneSysName]

	pg.Permissions = []PermissionObject{
		{Key: STUDIO_CREATE_DELETE_ROLE, Value: 0},
		{Key: STUDIO_ADD_REMOVE_USER_TO_ROLE, Value: 0},
		{Key: STUDIO_MANAGE_INTEGRATION, Value: 0},
		{Key: STUDIO_CREATE_COLLECTION, Value: 0},
		{Key: STUDIO_METADATA_UPDATE, Value: 0},
		{Key: STUDIO_CHANGE_CANVAS_COLLECTION_POSITION, Value: 0},
		{Key: STUDIO_EDIT_STUDIO_PROFILE, Value: 0},
		{Key: STUDIO_DELETE, Value: 0},
		{Key: STUDIO_MANAGE_PERMS, Value: 0},
		{Key: STUDIO_CAN_MANAGE_BILLING, Value: 0},
	}
	return pg
}

// Studio: Admins
func PGStudioAdmin() PermissionsTemplate {
	pg := PermissionsTemplate{}
	pg.DisplayName = "Admin"
	pg.SystemName = "pg_studio_admin"
	pg.StudioID = 0
	pg.Type = PGTYPESTUDIO
	pg.Weight = 1000
	pg.BetterPermissions = models.StudioPermissionsMap[models.PGStudioAdminSysName]
	pg.Permissions = []PermissionObject{
		{Key: STUDIO_CREATE_DELETE_ROLE, Value: 1},
		{Key: STUDIO_ADD_REMOVE_USER_TO_ROLE, Value: 1},
		{Key: STUDIO_MANAGE_INTEGRATION, Value: 1},
		{Key: STUDIO_CREATE_COLLECTION, Value: 1},
		{Key: STUDIO_METADATA_UPDATE, Value: 1},
		{Key: STUDIO_CHANGE_CANVAS_COLLECTION_POSITION, Value: 1},
		{Key: STUDIO_EDIT_STUDIO_PROFILE, Value: 1},
		{Key: STUDIO_DELETE, Value: 1},
		{Key: STUDIO_MANAGE_PERMS, Value: 1},
		{Key: STUDIO_CAN_MANAGE_BILLING, Value: 1},
	}
	return pg
}
