package permissiongroup

const PGTYPESTUDIO = "studio"
const PGTYPECOLLECTIION = "collection"
const PGTYPECANVAS = "canvas"

type PermissionObject struct {
	Key   string `json:"key"`   // uppercase no space CAN_INVITE
	Value uint   `json:"value"` // value can be 0, 1
}
type PermissionsTemplate struct {
	StudioID          int64              `json:"studioID"`
	DisplayName       string             `json:"displayName"`
	SystemName        string             `json:"systemName"`
	Type              string             `json:"type"` // Can be studio/Canvas or /Collections
	Weight            uint               `json:"weight"`
	BetterPermissions map[string]int     `json:"betterPermissions"`
	Permissions       []PermissionObject `json:"permissions"` // Has many permisions
}
