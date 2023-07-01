package permissiongroup

var UserAccessCanvasPermissionsList = []string{
	"pg_canvas_branch_view_metadata",
	"pg_canvas_branch_view",
	"pg_canvas_branch_comment",
	"pg_canvas_branch_edit",
	"pg_canvas_branch_moderate",
}

var UserAccessViewCanvasPermissionsList = []string{
	"pg_canvas_branch_view",
	"pg_canvas_branch_comment",
	"pg_canvas_branch_edit",
	"pg_canvas_branch_moderate",
}

const (
	PG_STUIDO_ADMIN = "pg_studio_admin"
)

var MapCollectionCanvasPerms = map[string]string{
	"pg_collection_none":          "pg_canvas_branch_none",
	"pg_collection_view_metadata": "pg_canvas_branch_view_metadata",
	"pg_collection_view":          "pg_canvas_branch_view",
	"pg_collection_comment":       "pg_canvas_branch_comment",
	"pg_collection_edit":          "pg_canvas_branch_edit",
	"pg_collection_moderate":      "pg_canvas_branch_moderate",
}
