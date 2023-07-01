package models

type StudioMembersRequest struct {
	BaseModel

	UserID uint64
	User   *User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// Pending // Rejected // Accepted
	Action string `gorm:"default:Pending"`

	StudioID uint64
	Studio   *Studio `gorm:"foreignKey:StudioID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	ActionByID   uint64
	ActionByUser *User `gorm:"foreignKey:ActionByID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
