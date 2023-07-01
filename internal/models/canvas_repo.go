package models

import (
	"time"
)

func (m *CanvasRepository) TableName() string {
	return "canvas_repositories"
}

const CanvasMaxPrivateReposAllowed = 25

// Parent CanvasRepository
type CanvasRepository struct {
	BaseModel
	// Future is product wants to make a Repo be part of many collections, we can handle that here.
	CollectionID uint64 // Parent collection
	StudioID     uint64 // This is more for query than any other use

	Name string `gorm:"type: varchar(500)"`
	// Deprecated
	Position    uint  // Talk to PW also when we do MOVE API -> Position Relative to this Collection (Has to be 1,2,3,4)
	Rank        int32 `gorm:"default:0"`
	Icon        string
	IsPublished bool
	// Cover Image
	CoverUrl string

	DefaultBranchID *uint64
	//self referential
	ParentCanvasRepositoryID *uint64

	CreatedByID uint64
	UpdatedByID uint64

	// Language Canvas Repo related
	DefaultLanguageCanvasRepoID *uint64
	Language                    *string `gorm:"default:en"`
	IsLanguageCanvas            bool    `gorm:"default:false"`
	AutoTranslated              bool    `gorm:"default:false"`
	IsProcessing                bool    `gorm:"default:false"` // used for translation pages

	// Soft Delete
	IsArchived   bool
	ArchivedAt   time.Time
	ArchivedByID *uint64
	Key          string

	SubCanvasCount  int  `gorm:"default:0"`
	HasPublicCanvas bool `gorm:"default:false"`

	Studio                 *Studio           `gorm:"foreignKey:StudioID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Collection             *Collection       `gorm:"foreignKey:CollectionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ParentCanvasRepository *CanvasRepository `gorm:"foreignKey:ParentCanvasRepositoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	DefaultBranch          *CanvasBranch     `gorm:"foreignkey:DefaultBranchID;constraint:OnDelete:CASCADE;"`
	CreatedByUser          *User             `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser          *User             `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ArchivedByUser         *User             `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

/*CanvasBranch
canvas_id

title
created_by is owner
is_draft : fasle
**edited - NR
is_merged ; False
from:
committed // git
created_from_commit_id
*/
