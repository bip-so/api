package models

import (
	"gorm.io/datatypes"
	"time"
)

func (m ReelComment) TableName() string {
	return "reel_comments"
}

// Comment is only on Reel Threads and relevence is only to Threads
// We can init and thread on Reels to enable comment.
type ReelComment struct {
	CommentBase
	// Query
	ReelID   uint64
	ParentID *uint64

	// FK
	Reel   *Reel        `gorm:"foreignKey:ReelID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Parent *ReelComment `gorm:"foreignkey:ParentID;constraint:OnDelete:CASCADE;"`
	// Boiler Plate
	BaseModel    // Boilerplate Stuff
	CreatedByID  uint64
	UpdatedByID  uint64
	ArchivedByID uint64
	IsArchived   bool
	ArchivedAt   time.Time

	CreatedByUser  *User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser  *User `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ArchivedByUser *User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Reactions             datatypes.JSON  // Reactions
	Mentions              *datatypes.JSON // Mentions
	ClonedFromReelComment uint64          `gorm:"default:0"`

	CommentCount uint `gorm:"default:0"` // This is just a count so we can show is UI //- Everything
}
