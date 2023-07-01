package models

var PGCollectionNoneSysName = "pg_collection_none"
var PGCollectionViewSysName = "pg_collection_view"
var PGCollectionViewMetadataSysName = "pg_collection_view_metadata"
var PGCollectionCommentSysName = "pg_collection_comment"
var PGCollectionEditSysName = "pg_collection_edit"
var PGCollectionModerateSysName = "pg_collection_moderate"

var CollectionPermissionsMap = map[string]map[string]int{
	PGCollectionNoneSysName: {
		"COLLECTION_MEMBERSHIP_MANAGE":         0,
		"COLLECTION_PUBLIC_ACCESS_CHANGE":      0,
		"COLLECTION_MANAGE_PERMS":              0,
		"COLLECTION_DELETE":                    0,
		"COLLECTION_OVERRIDE_STUDIO_MODE_ROLE": 0,
		"COLLECTION_MANAGE_PUBLISH_REQUEST":    0,
		"COLLECTION_EDIT_NAME":                 0,
		"COLLECTION_VIEW_METADATA":             0,
	},
	PGCollectionViewMetadataSysName: {
		"COLLECTION_MEMBERSHIP_MANAGE":         0,
		"COLLECTION_PUBLIC_ACCESS_CHANGE":      0,
		"COLLECTION_MANAGE_PERMS":              0,
		"COLLECTION_DELETE":                    0,
		"COLLECTION_OVERRIDE_STUDIO_MODE_ROLE": 0,
		"COLLECTION_MANAGE_PUBLISH_REQUEST":    0,
		"COLLECTION_EDIT_NAME":                 0,
		"COLLECTION_VIEW_METADATA":             1,
	},
	PGCollectionViewSysName: {
		"COLLECTION_MEMBERSHIP_MANAGE":         0,
		"COLLECTION_PUBLIC_ACCESS_CHANGE":      1,
		"COLLECTION_MANAGE_PERMS":              0,
		"COLLECTION_DELETE":                    0,
		"COLLECTION_OVERRIDE_STUDIO_MODE_ROLE": 0,
		"COLLECTION_MANAGE_PUBLISH_REQUEST":    0,
		"COLLECTION_EDIT_NAME":                 0,
		"COLLECTION_VIEW_METADATA":             1,
	},
	PGCollectionCommentSysName: {
		"COLLECTION_MEMBERSHIP_MANAGE":         0,
		"COLLECTION_PUBLIC_ACCESS_CHANGE":      1,
		"COLLECTION_MANAGE_PERMS":              0,
		"COLLECTION_DELETE":                    0,
		"COLLECTION_OVERRIDE_STUDIO_MODE_ROLE": 1,
		"COLLECTION_MANAGE_PUBLISH_REQUEST":    0,
		"COLLECTION_EDIT_NAME":                 0,
		"COLLECTION_VIEW_METADATA":             1,
	},
	PGCollectionEditSysName: {
		"COLLECTION_MEMBERSHIP_MANAGE":         0,
		"COLLECTION_PUBLIC_ACCESS_CHANGE":      1,
		"COLLECTION_MANAGE_PERMS":              0,
		"COLLECTION_DELETE":                    0,
		"COLLECTION_OVERRIDE_STUDIO_MODE_ROLE": 1,
		"COLLECTION_MANAGE_PUBLISH_REQUEST":    0,
		"COLLECTION_EDIT_NAME":                 0,
		"COLLECTION_VIEW_METADATA":             1,
	},
	PGCollectionModerateSysName: {
		"COLLECTION_MEMBERSHIP_MANAGE":         1,
		"COLLECTION_PUBLIC_ACCESS_CHANGE":      1,
		"COLLECTION_MANAGE_PERMS":              1,
		"COLLECTION_DELETE":                    1,
		"COLLECTION_OVERRIDE_STUDIO_MODE_ROLE": 1,
		"COLLECTION_MANAGE_PUBLISH_REQUEST":    1,
		"COLLECTION_EDIT_NAME":                 1,
		"COLLECTION_VIEW_METADATA":             1,
	},
}
