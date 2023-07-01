package models

import (
	"time"
)

func (m *Member) TableName() string {
	return "members"
}

/*
	Studio Has Many Members
*/
type Member struct {
	BaseModel
	// Todo: We need to add `gorm:"uniqueIndex:user_studio_idx"`
	UserID   uint64 // User id for this Member
	StudioID uint64 // Studio to which Member belongs

	CreatedByID  uint64
	UpdatedByID  uint64
	IsArchived   bool
	ArchivedAt   time.Time
	ArchivedByID uint64

	Roles    []Role    `gorm:"many2many:role_members;"` // Member has many roles.
	JoinedAt time.Time `gorm:"autoCreateTime"`

	// Removed and Unremoved API
	// IS removed
	IsRemoved     bool `gorm:"default:False"` // Removed by Admins and All Permissions Gone
	RemovedByID   uint64
	HasLeft       bool   `gorm:"default:False"`      // Left on -> Add to DB // Left self and we keep permissions in case decides to join in future
	RemovedReason string `gorm:"type: varchar(100)"` // Ban//Removed/+++

	User           *User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser  *User   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	RemovedByUser  *User   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser  *User   `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ArchivedByUser *User   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Studio         *Studio `gorm:"foreignKey:StudioID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func NewMember(userId uint64, studioId uint64) *Member {
	return &Member{
		UserID:      userId,
		StudioID:    studioId,
		CreatedByID: userId,
		UpdatedByID: userId,
	}
}
