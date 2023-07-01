package models

// https://tinyurl.com/27r7l62o
type BranchAccessToken struct {
	BaseModel
	InviteCode      string
	BranchID        uint64
	Branch          *CanvasBranch `gorm:"foreignkey:BranchID;constraint:OnDelete:CASCADE;"`
	IsActive        bool
	PermissionGroup string // Canvas Permissions Group
	CreatedByID     uint64
	CreatedByUser   *User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	//RepositoryID    uint64
	//Repository   *CanvasRepository `gorm:"foreignKey:RepositoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
