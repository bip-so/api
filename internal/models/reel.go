package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

func (m Reel) TableName() string {
	return "reels"
}

// Enable followers to consume important content from the studio with low mental load
// Reels is sub entity which is created on top of the canvas to extraxt certain parts of canvas
// Which can be then embdeed into a Reels
// Reel {{Comment}} {{Block(s)}}
// Only show in branch created in and carry forward on merge. For now, no option to NOT carry forward on merge.
// Will always carry forward as long as relevant blocks / content is accepted in merge. In the future, will be optional.
type Reel struct {
	StudioID           uint64         // Added for Query Mainly
	CanvasRepositoryID uint64         // Canvas Repo this belongs to "null"
	CanvasBranchID     uint64         // Canvas Branch this belongs to originally {id name}
	StartBlockID       uint64         //Block this belongs to "null"
	StartBlockUUID     uuid.UUID      `gorm:"type:uuid"`           // Incase orignal is deleted
	SelectedBlocks     datatypes.JSON `gorm:"default:'{}'::jsonb"` // This will be List of UUID's for a selected blocks.
	// ID of all block

	// Only for Legacy - Not needed for new cool implementation
	TextRangeStart uint // Legacy:  Slate Specific be:0
	TextRangeEnd   uint // Legacy: Slate Specific be:12

	RangeStart datatypes.JSON // Slate Specific be:0
	RangeEnd   datatypes.JSON // Slate Specific be:12

	// Need a way to get
	HighlightedText datatypes.JSON // Selected Text
	// This is json for allowing lighter version of slate to be used.
	// This is the Actual Comment
	ContextData  datatypes.JSON // Message or Context to thge reel.
	CommentCount uint           // This is just a count so we can show is UI //- Everything
	AuthorID     uint64         // Person who has commented // This would be userID nor MemberID

	// FK
	CanvasBranch     *CanvasBranch     `gorm:"foreignkey:CanvasBranchID;constraint:OnDelete:CASCADE;"`
	CanvasRepository *CanvasRepository `gorm:"foreignKey:CanvasRepositoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Block            *Block            `gorm:"foreignKey:StartBlockID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Author           *User             `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Studio           *Studio           `gorm:"foreignKey:StudioID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// Boiler Plate
	BaseModel      // Boilerplate Stuff
	CreatedByID    uint64
	UpdatedByID    uint64
	ArchivedByID   *uint64
	IsArchived     bool
	ArchivedAt     time.Time
	CreatedByUser  *User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser  *User `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ArchivedByUser *User `gorm:"foreignKey:ArchivedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Reactions      datatypes.JSON  // Reactions
	Mentions       *datatypes.JSON // Mentions
	ClonedFromReel uint64          `gorm:"default:0"`
}

func (m Reel) NewReel(studioID uint64, canvasRepositoryID uint64, canvasBranchID uint64, startBlockID uint64, startBlockUUID uuid.UUID, textRangeStart uint, textRangeEnd uint, rangeStart datatypes.JSON, rangeEnd datatypes.JSON, highlightedText datatypes.JSON, contextData datatypes.JSON, commentCount uint, authorID uint64, createdByID uint64, updatedByID uint64, selectedblocks datatypes.JSON) *Reel {
	return &Reel{StudioID: studioID, CanvasRepositoryID: canvasRepositoryID, CanvasBranchID: canvasBranchID, StartBlockID: startBlockID, StartBlockUUID: startBlockUUID, TextRangeStart: textRangeStart, TextRangeEnd: textRangeEnd, RangeStart: rangeStart, RangeEnd: rangeEnd, HighlightedText: highlightedText, ContextData: contextData, CommentCount: commentCount, AuthorID: authorID, CreatedByID: createdByID, UpdatedByID: updatedByID, SelectedBlocks: selectedblocks}
}
