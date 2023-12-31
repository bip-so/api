package models

var PGCanvasNoneSysName = "pg_canvas_branch_none"
var PGCanvasViewMetadataSysName = "pg_canvas_branch_view_metadata"
var PGCanvasViewSysName = "pg_canvas_branch_view"
var PGCanvasCommentSysName = "pg_canvas_branch_comment"
var PGCanvasEditSysName = "pg_canvas_branch_edit"
var PGCanvasModerateSysName = "pg_canvas_branch_moderate"

var CanvasPermissionsMap = map[string]map[string]int{
	PGCanvasNoneSysName: {
		"CANVAS_BRANCH_VIEW":                    0,
		"CANVAS_BRANCH_EDIT":                    0,
		"CANVAS_BRANCH_EDIT_NAME":               0,
		"CANVAS_BRANCH_DELETE":                  0,
		"CANVAS_BRANCH_ADD_COMMENT":             0,
		"CANVAS_BRANCH_ADD_REACTION":            0,
		"CANVAS_BRANCH_CREATE_REEL":             0,
		"CANVAS_BRANCH_COMMENT_ON_REEL":         0,
		"CANVAS_BRANCH_REACT_TO_REEL":           0,
		"CANVAS_BRANCH_MANAGE_PERMS":            0,
		"CANVAS_BRANCH_MANAGE_MERGE_REQUESTS":   0,
		"CANVAS_BRANCH_MANAGE_PUBLISH_REQUESTS": 0,
		"CANVAS_BRANCH_CREATE_MERGE_REQUEST":    0,
		"CANVAS_BRANCH_CREATE_PUBLISH_REQUEST":  0,
		"CANVAS_BRANCH_CLONE":                   0,
		"CANVAS_BRANCH_RESOLVE_COMMENTS":        0,
		"CANVAS_BRANCH_MANAGE_CONTENT":          0,
		"CANVAS_BRANCH_VIEW_METADATA":           0,
	},
	PGCanvasViewMetadataSysName: {
		"CANVAS_BRANCH_VIEW":                    0,
		"CANVAS_BRANCH_EDIT":                    0,
		"CANVAS_BRANCH_EDIT_NAME":               0,
		"CANVAS_BRANCH_DELETE":                  0,
		"CANVAS_BRANCH_ADD_COMMENT":             0,
		"CANVAS_BRANCH_ADD_REACTION":            0,
		"CANVAS_BRANCH_CREATE_REEL":             0,
		"CANVAS_BRANCH_COMMENT_ON_REEL":         0,
		"CANVAS_BRANCH_REACT_TO_REEL":           0,
		"CANVAS_BRANCH_MANAGE_PERMS":            0,
		"CANVAS_BRANCH_MANAGE_MERGE_REQUESTS":   0,
		"CANVAS_BRANCH_MANAGE_PUBLISH_REQUESTS": 0,
		"CANVAS_BRANCH_CREATE_MERGE_REQUEST":    0,
		"CANVAS_BRANCH_CREATE_PUBLISH_REQUEST":  0,
		"CANVAS_BRANCH_CLONE":                   0,
		"CANVAS_BRANCH_RESOLVE_COMMENTS":        0,
		"CANVAS_BRANCH_MANAGE_CONTENT":          0,
		"CANVAS_BRANCH_VIEW_METADATA":           1,
	},
	PGCanvasViewSysName: {"CANVAS_BRANCH_VIEW": 1,
		"CANVAS_BRANCH_EDIT":                    0,
		"CANVAS_BRANCH_EDIT_NAME":               0,
		"CANVAS_BRANCH_DELETE":                  0,
		"CANVAS_BRANCH_ADD_COMMENT":             0,
		"CANVAS_BRANCH_ADD_REACTION":            0,
		"CANVAS_BRANCH_CREATE_REEL":             0,
		"CANVAS_BRANCH_COMMENT_ON_REEL":         1,
		"CANVAS_BRANCH_REACT_TO_REEL":           1,
		"CANVAS_BRANCH_MANAGE_PERMS":            0,
		"CANVAS_BRANCH_MANAGE_MERGE_REQUESTS":   0,
		"CANVAS_BRANCH_MANAGE_PUBLISH_REQUESTS": 0,
		"CANVAS_BRANCH_CREATE_MERGE_REQUEST":    0,
		"CANVAS_BRANCH_CREATE_PUBLISH_REQUEST":  0,
		"CANVAS_BRANCH_CLONE":                   1,
		"CANVAS_BRANCH_RESOLVE_COMMENTS":        0,
		"CANVAS_BRANCH_MANAGE_CONTENT":          0,
		"CANVAS_BRANCH_VIEW_METADATA":           1,
	},
	PGCanvasCommentSysName: {
		"CANVAS_BRANCH_VIEW":                    1,
		"CANVAS_BRANCH_EDIT":                    0,
		"CANVAS_BRANCH_EDIT_NAME":               0,
		"CANVAS_BRANCH_DELETE":                  0,
		"CANVAS_BRANCH_ADD_COMMENT":             1,
		"CANVAS_BRANCH_ADD_REACTION":            1,
		"CANVAS_BRANCH_CREATE_REEL":             0,
		"CANVAS_BRANCH_COMMENT_ON_REEL":         1,
		"CANVAS_BRANCH_REACT_TO_REEL":           1,
		"CANVAS_BRANCH_MANAGE_PERMS":            0,
		"CANVAS_BRANCH_MANAGE_MERGE_REQUESTS":   0,
		"CANVAS_BRANCH_MANAGE_PUBLISH_REQUESTS": 0,
		"CANVAS_BRANCH_CREATE_MERGE_REQUEST":    0,
		"CANVAS_BRANCH_CREATE_PUBLISH_REQUEST":  0,
		"CANVAS_BRANCH_CLONE":                   1,
		"CANVAS_BRANCH_RESOLVE_COMMENTS":        0,
		"CANVAS_BRANCH_MANAGE_CONTENT":          0,
		"CANVAS_BRANCH_VIEW_METADATA":           1,
	},
	PGCanvasEditSysName: {
		"CANVAS_BRANCH_VIEW":                    1,
		"CANVAS_BRANCH_EDIT":                    1,
		"CANVAS_BRANCH_EDIT_NAME":               0,
		"CANVAS_BRANCH_DELETE":                  0,
		"CANVAS_BRANCH_ADD_COMMENT":             1,
		"CANVAS_BRANCH_ADD_REACTION":            1,
		"CANVAS_BRANCH_CREATE_REEL":             1,
		"CANVAS_BRANCH_COMMENT_ON_REEL":         1,
		"CANVAS_BRANCH_REACT_TO_REEL":           1,
		"CANVAS_BRANCH_MANAGE_PERMS":            0,
		"CANVAS_BRANCH_MANAGE_MERGE_REQUESTS":   0,
		"CANVAS_BRANCH_MANAGE_PUBLISH_REQUESTS": 0,
		"CANVAS_BRANCH_CREATE_MERGE_REQUEST":    0,
		"CANVAS_BRANCH_CREATE_PUBLISH_REQUEST":  0,
		"CANVAS_BRANCH_CLONE":                   1,
		"CANVAS_BRANCH_RESOLVE_COMMENTS":        0,
		"CANVAS_BRANCH_MANAGE_CONTENT":          0,
		"CANVAS_BRANCH_VIEW_METADATA":           1,
	},
	PGCanvasModerateSysName: {
		"CANVAS_BRANCH_VIEW":                    1,
		"CANVAS_BRANCH_EDIT":                    1,
		"CANVAS_BRANCH_EDIT_NAME":               1,
		"CANVAS_BRANCH_DELETE":                  1,
		"CANVAS_BRANCH_ADD_COMMENT":             1,
		"CANVAS_BRANCH_ADD_REACTION":            1,
		"CANVAS_BRANCH_CREATE_REEL":             1,
		"CANVAS_BRANCH_COMMENT_ON_REEL":         1,
		"CANVAS_BRANCH_REACT_TO_REEL":           1,
		"CANVAS_BRANCH_MANAGE_PERMS":            1,
		"CANVAS_BRANCH_MANAGE_MERGE_REQUESTS":   1,
		"CANVAS_BRANCH_MANAGE_PUBLISH_REQUESTS": 1,
		"CANVAS_BRANCH_CREATE_MERGE_REQUEST":    1,
		"CANVAS_BRANCH_CREATE_PUBLISH_REQUEST":  1,
		"CANVAS_BRANCH_CLONE":                   1,
		"CANVAS_BRANCH_RESOLVE_COMMENTS":        1,
		"CANVAS_BRANCH_MANAGE_CONTENT":          1,
		"CANVAS_BRANCH_VIEW_METADATA":           1,
	},
}
