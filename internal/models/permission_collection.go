package models

func (m *CollectionPermission) TableName() string {
	return "collection_permissions"
}

type CollectionPermission struct {
	BaseModel
	StudioID     uint64
	CollectionId uint64

	PermissionGroup string
	RoleId          *uint64
	MemberId        *uint64
	IsOverridden    bool `gorm:"default:false;"`

	Collection Collection `gorm:"foreignKey:CollectionId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Role       *Role      `gorm:"foreignKey:RoleId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Studio     *Studio    `gorm:"foreignKey:StudioID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Member     *Member    `gorm:"foreignKey:MemberId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
