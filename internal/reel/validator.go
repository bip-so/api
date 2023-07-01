package reel

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// New Reel
type NewReelCreatePOST struct {
	CanvasRepositoryID uint64         `json:"canvasRepositoryID"`
	CanvasBranchID     uint64         `json:"canvasBranchID"`
	StartBlockUUID     uuid.UUID      `json:"startBlockUUID"`
	TextRangeStart     uint           `json:"textRangeStart"`
	TextRangeEnd       uint           `json:"textRangeEnd"`
	RangeStart         datatypes.JSON `json:"rangeStart"`
	RangeEnd           datatypes.JSON `json:"rangeEnd"`
	HighlightedText    datatypes.JSON `json:"highlightedText"`
	ContextData        datatypes.JSON `json:"contextData"`
	SelectedBlocks     datatypes.JSON `json:"selectedBlocks"`
}

type ReelCommentCreatePOST struct {
	// Comment Base
	Data     datatypes.JSON `json:"data"`
	IsEdited bool           `json:"isEdited"`
	IsReply  bool           `json:"isReply"`
	ParentID *uint64        `json:"parentId"`
}
