package models

func (m *CanvasBranchPermission) TableName() string {
	return "canvas_branch_permissions"
}

type CanvasBranchPermission struct {
	BaseModel
	StudioID                    uint64
	CollectionId                uint64
	CanvasBranchID              *uint64 // Canvas Branch this belongs to _>  CB1
	CanvasRepositoryID          uint64  // Canvas Repository this belongs to
	CbpParentCanvasRepositoryID *uint64

	PermissionGroup string
	RoleId          *uint64
	MemberId        *uint64
	IsOverridden    bool `gorm:"default:false;"`

	CanvasBranch              *CanvasBranch     `gorm:"foreignkey:CanvasBranchID;constraint:OnDelete:CASCADE;"`
	Collection                Collection        `gorm:"foreignKey:CollectionId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CanvasRepository          *CanvasRepository `gorm:"foreignKey:CanvasRepositoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CbpParentCanvasRepository *CanvasRepository `gorm:"foreignKey:CbpParentCanvasRepositoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Role                      *Role             `gorm:"foreignKey:RoleId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Studio                    *Studio           `gorm:"foreignKey:StudioID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Member                    *Member           `gorm:"foreignKey:MemberId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
