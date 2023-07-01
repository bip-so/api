package models

type BranchInviteViaEmail struct {
	BaseModel
	Email           string
	BranchID        uint64
	StudioID        uint64
	PermissionGroup string // Canvas PermissionsR
	CreatedByID     uint64
	Branch          *CanvasBranch `gorm:"foreignkey:BranchID;constraint:OnDelete:CASCADE;"`
	CreatedByUser   *User         `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	//RepositoryID    uint64
	//Repository   *CanvasRepository `gorm:"foreignKey:RepositoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
