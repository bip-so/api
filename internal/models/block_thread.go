package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"time"
)

//Mail
//	Blocks for Main Branch
//B1 :
//	Br1Blk1 Help
//	Br2Blk1 Hello

// BlockThread Create a Thread
// Pagination : 50?
// All the unresolved comments should be sent
func (m BlockThread) TableName() string {
	return "block_threads"
}

type BlockThread struct {

	// Query
	CanvasRepositoryID uint64 // Canvas Repo this belongs to "null"
	CanvasBranchID     uint64 // Canvas Branch this belongs to originally

	StartBlockID   uint64    // Block this belongs to "null"
	StartBlockUUID uuid.UUID `gorm:"type:uuid"`

	//"Can be across Blocks"

	Position uint // May be useful may comment in same

	// Only for Legacy - Not needed for new cool implementation
	TextRangeStart uint   // Slate Specific be:0
	TextRangeEnd   uint   // Slate Specific be:12
	Text           string // Temp may be deleted in future

	HighlightedText string // Selected Text

	CommentCount uint // This is just a count so we can show is UI //- Everything +/0
	// Deleted
	// AuthorID     uint64 // Person who has commented // This would be userID nor MemberID
	// Author           *User             `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ClonedFromThread uint64 `gorm:"default:0"`

	//Resolution
	IsResolved   bool
	ResolvedByID *uint64
	ResolvedAt   time.Time

	// FK
	CanvasBranch     *CanvasBranch     `gorm:"foreignkey:CanvasBranchID;constraint:OnDelete:CASCADE;"`
	CanvasRepository *CanvasRepository `gorm:"foreignKey:CanvasRepositoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Block            *Block            `gorm:"foreignKey:StartBlockID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ResolvedByUser   *User             `gorm:"foreignKey:ResolvedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	// Boiler Plate
	BaseModel      // Boilerplate Stuff
	CreatedByID    uint64
	UpdatedByID    uint64
	ArchivedByID   uint64
	IsArchived     bool
	ArchivedAt     time.Time
	CreatedByUser  *User           `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser  *User           `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ArchivedByUser *User           `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Reactions      datatypes.JSON  // Reactions
	Mentions       *datatypes.JSON // Mentions
}
