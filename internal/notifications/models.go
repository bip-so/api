package notifications

import "gitlab.com/phonepost/bip-be-platform/internal/models"

type NotificationEntity struct {
	Entity     string
	App        bool
	Email      bool
	Discord    bool
	IsPersonal bool
	Events     map[string]NotificationEvent
}

type NotificationEvent struct {
	Activity string
	Text     string
	Priority string // enum types will be ["high", "medium", "low"]
}

var PermissionsTextMap = map[string]string{
	models.PGCanvasNoneSysName:         "None",
	models.PGCanvasViewSysName:         "View",
	models.PGCanvasViewMetadataSysName: "ViewMetadata",
	models.PGCanvasEditSysName:         "Edit",
	models.PGCanvasCommentSysName:      "Comment",
	models.PGCanvasModerateSysName:     "Moderate",
}
