package models

var PGStudioNoneSysName = "pg_studio_none"
var PGStudioAdminSysName = "pg_studio_admin"

var StudioPermissionsMap = map[string]map[string]int{
	PGStudioNoneSysName: {
		"STUDIO_CREATE_DELETE_ROLE":                0,
		"STUDIO_ADD_REMOVE_USER_TO_ROLE":           0,
		"STUDIO_MANAGE_INTEGRATION":                0,
		"STUDIO_CREATE_COLLECTION":                 0,
		"STUDIO_METADATA_UPDATE":                   0,
		"STUDIO_CHANGE_CANVAS_COLLECTION_POSITION": 0,
		"STUDIO_EDIT_STUDIO_PROFILE":               0,
		"STUDIO_DELETE":                            0,
		"STUDIO_MANAGE_PERMS":                      0,
	},
	PGStudioAdminSysName: {
		"STUDIO_CREATE_DELETE_ROLE":                1,
		"STUDIO_ADD_REMOVE_USER_TO_ROLE":           1,
		"STUDIO_MANAGE_INTEGRATION":                1,
		"STUDIO_CREATE_COLLECTION":                 1,
		"STUDIO_METADATA_UPDATE":                   1,
		"STUDIO_CHANGE_CANVAS_COLLECTION_POSITION": 1,
		"STUDIO_EDIT_STUDIO_PROFILE":               1,
		"STUDIO_DELETE":                            1,
		"STUDIO_MANAGE_PERMS":                      1,
	},
}
