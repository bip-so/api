package models

import (
	"time"
)

func (m *Collection) TableName() string {
	return "collections"
}

type Collection struct {
	BaseModel

	CreatedByID   uint64
	UpdatedByID   uint64
	CreatedByUser *User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	IsArchived     bool
	ArchivedAt     time.Time
	ArchivedByID   *uint64
	ArchivedByUser *User `gorm:"foreignKey:ArchivedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Name string `gorm:"type: varchar(100)"`

	Rank int32 `gorm:"default:0"`
	// Deprecated
	Position uint // Talk to PW also when we do MOVE API ->
	Icon     string

	StudioID uint64

	// todo: Add ENUM
	PublicAccess string

	ParentCollectionID *uint64

	// do we need both ?
	// @todo Set defausudo rm -rf /usr/local/golts where ever possible

	ComputedRootCanvasCount int `gorm:"default:0"`
	ComputedAllCanvasCount  int `gorm:"default:0"`

	// @todo: Ask NR Do we delete child on parent delete ? Do dwe delete chiuldere on parent
	ParentCollection *Collection `gorm:"foreignKey:ParentCollectionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Studio           *Studio     `gorm:"foreignKey:StudioID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	HasPublicCanvas bool `gorm:"default:false"`
}
