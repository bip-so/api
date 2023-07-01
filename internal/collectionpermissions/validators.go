package collectionpermissions

type CollectionPermissionValidator struct {
	CollectionId uint64 `json:"collectionId" binding:"required"`
	PermGroup    string `json:"permGroup" binding:"required"`
	RoleID       uint64 `json:"roleID"`
	MemberID     uint64 `json:"memberID"`
	IsOverridden bool   `json:"isOverridden" default:"false"`
}
