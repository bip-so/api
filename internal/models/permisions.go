package models

const PermissionRedisNSPrefix = "permissions:"

// Patterns of redis keys
// studio: str(userid)+"-"+(studioid)
// collection : str(userid)+"-"+(collectionid)
// studio: str(userid)+"-"+(canvasBranchid)
// Example : PermissionRedisNSPrefix+str(userID)+"-"+(studioID)

// Remember to update  `internal/permission_groups/pg_studio.go`

var PermissionGroupWeightMap = map[string]map[string]int{
	"studio": {
		"pg_studio_none":  0,
		"pg_studio_admin": 1000,
	},
	"collection": {
		"pg_collection_none":          0,
		"pg_collection_view_metadata": 10,
		"pg_collection_view":          20,
		"pg_collection_comment":       30,
		"pg_collection_edit":          40,
		"pg_collection_moderate":      1000,
	},
	"canvas": {
		"pg_canvas_branch_none":          0,
		"pg_canvas_branch_view_metadata": 5,
		"pg_canvas_branch_view":          10,
		"pg_canvas_branch_comment":       30,
		"pg_canvas_branch_edit":          100,
		"pg_canvas_branch_moderate":      1000,
	},
}

// AnonymousUserPerms : If user is Anonymous then use this data
var AnonymousUserPerms = map[string]string{
	"studio":        "pg_studio_none",
	"collection":    "pg_collection_none",
	"canvas_branch": "pg_canvas_branch_none",
}

const PRIVATE = "private"
const EDIT = "edit"
const VIEW = "view"
const COMMENT = "comment"
