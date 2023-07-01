package permissiongroup

import "gitlab.com/phonepost/bip-be-platform/internal/models"

/*
Important Note If you do aanything here make sure to :
Please update : `internal/models/permisions.go` for any change here
*/

// @todo Add PGCanvasViewMetadata
// CANVAS_BRANCH_VIEW_METADATA = 1 for all other than NONE

// None
func PGCanvasNone() PermissionsTemplate {
	pg := PermissionsTemplate{}
	pg.DisplayName = "None"
	pg.SystemName = models.PGCanvasNoneSysName
	pg.StudioID = 0
	pg.Type = PGTYPECANVAS
	pg.Weight = 0
	pg.BetterPermissions = models.CanvasPermissionsMap[models.PGCanvasNoneSysName]
	pg.Permissions = []PermissionObject{
		{Key: CANVAS_BRANCH_VIEW, Value: 0},
		{Key: CANVAS_BRANCH_EDIT, Value: 0},
		{Key: CANVAS_BRANCH_EDIT_NAME, Value: 0},
		{Key: CANVAS_BRANCH_DELETE, Value: 0},
		{Key: CANVAS_BRANCH_ADD_COMMENT, Value: 0},
		{Key: CANVAS_BRANCH_ADD_REACTION, Value: 0},
		{Key: CANVAS_BRANCH_CREATE_REEL, Value: 0},
		{Key: CANVAS_BRANCH_COMMENT_ON_REEL, Value: 0},
		{Key: CANVAS_BRANCH_REACT_TO_REEL, Value: 0},
		{Key: CANVAS_BRANCH_MANAGE_PERMS, Value: 0},
		{Key: CANVAS_BRANCH_MANAGE_MERGE_REQUESTS, Value: 0},
		{Key: CANVAS_BRANCH_MANAGE_PUBLISH_REQUESTS, Value: 0},
		{Key: CANVAS_BRANCH_CREATE_MERGE_REQUEST, Value: 0},
		{Key: CANVAS_BRANCH_CREATE_PUBLISH_REQUEST, Value: 0},
		{Key: CANVAS_BRANCH_CLONE, Value: 0},
		{Key: CANVAS_BRANCH_RESOLVE_COMMENTS, Value: 0},
		{Key: CANVAS_BRANCH_MANAGE_CONTENT, Value: 0},
		{Key: CANVAS_BRANCH_VIEW_METADATA, Value: 0},
	}
	return pg
}

// AuthNone
func PGCanvasViewMetaData() PermissionsTemplate {
	pg := PermissionsTemplate{}
	pg.DisplayName = "ViewMetadata"
	pg.SystemName = models.PGCanvasViewMetadataSysName
	pg.StudioID = 0
	pg.Type = PGTYPECANVAS
	pg.Weight = 5
	//pg.Hidden = True
	pg.BetterPermissions = models.CanvasPermissionsMap[models.PGCanvasViewMetadataSysName]
	pg.Permissions = []PermissionObject{
		{Key: CANVAS_BRANCH_VIEW, Value: 1},
		{Key: CANVAS_BRANCH_EDIT, Value: 0},
		{Key: CANVAS_BRANCH_EDIT_NAME, Value: 0},
		{Key: CANVAS_BRANCH_DELETE, Value: 0},
		{Key: CANVAS_BRANCH_ADD_COMMENT, Value: 0},
		{Key: CANVAS_BRANCH_ADD_REACTION, Value: 0},
		{Key: CANVAS_BRANCH_CREATE_REEL, Value: 0},
		{Key: CANVAS_BRANCH_COMMENT_ON_REEL, Value: 0},
		{Key: CANVAS_BRANCH_REACT_TO_REEL, Value: 0},
		{Key: CANVAS_BRANCH_MANAGE_PERMS, Value: 0},
		{Key: CANVAS_BRANCH_MANAGE_MERGE_REQUESTS, Value: 0},
		{Key: CANVAS_BRANCH_MANAGE_PUBLISH_REQUESTS, Value: 0},
		{Key: CANVAS_BRANCH_CREATE_MERGE_REQUEST, Value: 0},
		{Key: CANVAS_BRANCH_CREATE_PUBLISH_REQUEST, Value: 0},
		{Key: CANVAS_BRANCH_CLONE, Value: 0},
		{Key: CANVAS_BRANCH_RESOLVE_COMMENTS, Value: 0},
		{Key: CANVAS_BRANCH_MANAGE_CONTENT, Value: 0},
		{Key: CANVAS_BRANCH_VIEW_METADATA, Value: 1},
	}
	return pg
}

func PGCanvasView() PermissionsTemplate {
	pg := PermissionsTemplate{}
	pg.DisplayName = "View"
	pg.SystemName = models.PGCanvasViewSysName
	pg.StudioID = 0
	pg.Type = PGTYPECANVAS
	pg.Weight = 10
	pg.BetterPermissions = models.CanvasPermissionsMap[models.PGCanvasViewSysName]
	pg.Permissions = []PermissionObject{
		{Key: CANVAS_BRANCH_VIEW, Value: 1},
		{Key: CANVAS_BRANCH_EDIT, Value: 0},
		{Key: CANVAS_BRANCH_EDIT_NAME, Value: 0},
		{Key: CANVAS_BRANCH_DELETE, Value: 0},
		{Key: CANVAS_BRANCH_ADD_COMMENT, Value: 0},
		{Key: CANVAS_BRANCH_ADD_REACTION, Value: 0},
		{Key: CANVAS_BRANCH_CREATE_REEL, Value: 0},
		{Key: CANVAS_BRANCH_COMMENT_ON_REEL, Value: 1},
		{Key: CANVAS_BRANCH_REACT_TO_REEL, Value: 1},
		{Key: CANVAS_BRANCH_MANAGE_PERMS, Value: 0},
		{Key: CANVAS_BRANCH_MANAGE_MERGE_REQUESTS, Value: 0},
		{Key: CANVAS_BRANCH_MANAGE_PUBLISH_REQUESTS, Value: 0},
		{Key: CANVAS_BRANCH_CREATE_MERGE_REQUEST, Value: 0},
		{Key: CANVAS_BRANCH_CREATE_PUBLISH_REQUEST, Value: 0},
		{Key: CANVAS_BRANCH_CLONE, Value: 1},
		{Key: CANVAS_BRANCH_RESOLVE_COMMENTS, Value: 0},
		{Key: CANVAS_BRANCH_MANAGE_CONTENT, Value: 0},
		{Key: CANVAS_BRANCH_VIEW_METADATA, Value: 1},
	}
	return pg
}

func PGCanvasComment() PermissionsTemplate {
	pg := PermissionsTemplate{}
	pg.DisplayName = "Comment"
	pg.SystemName = models.PGCanvasCommentSysName
	pg.StudioID = 0
	pg.Type = PGTYPECANVAS
	pg.Weight = 30
	pg.BetterPermissions = models.CanvasPermissionsMap[models.PGCanvasCommentSysName]
	pg.Permissions = []PermissionObject{
		{Key: CANVAS_BRANCH_VIEW, Value: 1},
		{Key: CANVAS_BRANCH_EDIT, Value: 0},
		{Key: CANVAS_BRANCH_EDIT_NAME, Value: 0},
		{Key: CANVAS_BRANCH_DELETE, Value: 0},
		{Key: CANVAS_BRANCH_ADD_COMMENT, Value: 1},
		{Key: CANVAS_BRANCH_ADD_REACTION, Value: 1},
		{Key: CANVAS_BRANCH_CREATE_REEL, Value: 0},
		{Key: CANVAS_BRANCH_COMMENT_ON_REEL, Value: 1},
		{Key: CANVAS_BRANCH_REACT_TO_REEL, Value: 1},
		{Key: CANVAS_BRANCH_MANAGE_PERMS, Value: 0},
		{Key: CANVAS_BRANCH_MANAGE_MERGE_REQUESTS, Value: 0},
		{Key: CANVAS_BRANCH_MANAGE_PUBLISH_REQUESTS, Value: 0},
		{Key: CANVAS_BRANCH_CREATE_MERGE_REQUEST, Value: 0},
		{Key: CANVAS_BRANCH_CREATE_PUBLISH_REQUEST, Value: 0},
		{Key: CANVAS_BRANCH_CLONE, Value: 1},
		{Key: CANVAS_BRANCH_RESOLVE_COMMENTS, Value: 0},
		{Key: CANVAS_BRANCH_MANAGE_CONTENT, Value: 0},
		{Key: CANVAS_BRANCH_VIEW_METADATA, Value: 1},
	}
	return pg
}

func PGCanvasEdit() PermissionsTemplate {
	pg := PermissionsTemplate{}
	pg.DisplayName = "Edit"
	pg.SystemName = models.PGCanvasEditSysName
	pg.StudioID = 0
	pg.Type = PGTYPECANVAS
	pg.Weight = 100
	pg.BetterPermissions = models.CanvasPermissionsMap[models.PGCanvasEditSysName]
	pg.Permissions = []PermissionObject{
		{Key: CANVAS_BRANCH_VIEW, Value: 1},
		{Key: CANVAS_BRANCH_EDIT, Value: 1},
		{Key: CANVAS_BRANCH_EDIT_NAME, Value: 0},
		{Key: CANVAS_BRANCH_DELETE, Value: 0},
		{Key: CANVAS_BRANCH_ADD_COMMENT, Value: 1},
		{Key: CANVAS_BRANCH_ADD_REACTION, Value: 1},
		{Key: CANVAS_BRANCH_CREATE_REEL, Value: 1},
		{Key: CANVAS_BRANCH_COMMENT_ON_REEL, Value: 1},
		{Key: CANVAS_BRANCH_REACT_TO_REEL, Value: 1},
		{Key: CANVAS_BRANCH_MANAGE_PERMS, Value: 0},
		{Key: CANVAS_BRANCH_MANAGE_MERGE_REQUESTS, Value: 0},
		{Key: CANVAS_BRANCH_MANAGE_PUBLISH_REQUESTS, Value: 0},
		{Key: CANVAS_BRANCH_CREATE_MERGE_REQUEST, Value: 0},
		{Key: CANVAS_BRANCH_CREATE_PUBLISH_REQUEST, Value: 0},
		{Key: CANVAS_BRANCH_CLONE, Value: 1},
		{Key: CANVAS_BRANCH_RESOLVE_COMMENTS, Value: 0},
		{Key: CANVAS_BRANCH_MANAGE_CONTENT, Value: 0},
		{Key: CANVAS_BRANCH_VIEW_METADATA, Value: 1},
	}
	return pg
}

func PGCanvasModerate() PermissionsTemplate {
	pg := PermissionsTemplate{}
	pg.DisplayName = "Moderate"
	pg.SystemName = models.PGCanvasModerateSysName
	pg.StudioID = 0
	pg.Type = PGTYPECANVAS
	pg.Weight = 1000
	pg.BetterPermissions = models.CanvasPermissionsMap[models.PGCanvasModerateSysName]
	pg.Permissions = []PermissionObject{
		{Key: CANVAS_BRANCH_VIEW, Value: 1},
		{Key: CANVAS_BRANCH_EDIT, Value: 1},
		{Key: CANVAS_BRANCH_EDIT_NAME, Value: 1},
		{Key: CANVAS_BRANCH_DELETE, Value: 1},
		{Key: CANVAS_BRANCH_ADD_COMMENT, Value: 1},
		{Key: CANVAS_BRANCH_ADD_REACTION, Value: 1},
		{Key: CANVAS_BRANCH_CREATE_REEL, Value: 1},
		{Key: CANVAS_BRANCH_COMMENT_ON_REEL, Value: 1},
		{Key: CANVAS_BRANCH_REACT_TO_REEL, Value: 1},
		{Key: CANVAS_BRANCH_MANAGE_PERMS, Value: 1},
		{Key: CANVAS_BRANCH_MANAGE_MERGE_REQUESTS, Value: 1},
		{Key: CANVAS_BRANCH_MANAGE_PUBLISH_REQUESTS, Value: 1},
		{Key: CANVAS_BRANCH_CREATE_MERGE_REQUEST, Value: 1},
		{Key: CANVAS_BRANCH_CREATE_PUBLISH_REQUEST, Value: 1},
		{Key: CANVAS_BRANCH_CLONE, Value: 1},
		{Key: CANVAS_BRANCH_RESOLVE_COMMENTS, Value: 1},
		{Key: CANVAS_BRANCH_MANAGE_CONTENT, Value: 1},
		{Key: CANVAS_BRANCH_VIEW_METADATA, Value: 1},
	}
	return pg
}
