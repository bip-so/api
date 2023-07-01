package models

import (
	"gorm.io/datatypes"
	"time"
)

func (m BlockComment) TableName() string {
	return "block_comments"
}

// Comment is only on Thread and relevence is only to Threads
type BlockComment struct {
	// Comment Base
	CommentBase
	// Query
	ThreadID uint64
	ParentID *uint64

	// FK
	Thread                  *BlockThread  `gorm:"foreignKey:ThreadID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Parent                  *BlockComment `gorm:"foreignkey:ParentID;constraint:OnDelete:CASCADE;"`
	ClonedFromThreadComment uint64        `gorm:"default:0"`

	// Boiler Plate
	BaseModel    // Boilerplate Stuff
	CreatedByID  uint64
	UpdatedByID  uint64
	ArchivedByID uint64
	IsArchived   bool
	ArchivedAt   time.Time

	CreatedByUser  *User           `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser  *User           `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ArchivedByUser *User           `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Reactions      datatypes.JSON  // Reactions
	Mentions       *datatypes.JSON // Mentions
}
