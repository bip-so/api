package models

func (m *StudioPermission) TableName() string {
	return "studio_permissions"
}

type StudioPermission struct {
	BaseModel
	StudioID uint64

	PermissionGroup string
	RoleId          *uint64
	MemberId        *uint64
	IsOverridden    bool `gorm:"default:false;"`

	Role   *Role   `gorm:"foreignKey:RoleId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Studio *Studio `gorm:"foreignKey:StudioID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Member *Member `gorm:"foreignKey:MemberId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// uniquetogether?
