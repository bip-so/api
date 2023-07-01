package blockthread

import "github.com/google/uuid"

type EmptyBlockThread struct {
}

// New blockThreads
type PostBlockThread struct {
	CanvasRepositoryID uint64    `json:"canvasRepositoryId" binding:"required"`
	CanvasBranchID     uint64    `json:"canvasBranchId" binding:"required"`
	StartBlockUUID     uuid.UUID `json:"startBlockUUID"`
	Position           uint      `json:"position" binding:"required"`
	TextRangeStart     uint      `json:"textRangeStart"`
	TextRangeEnd       uint      `json:"textRangeEnd"`
	Text               string    `json:"text" binding:"required"`
	HighlightedText    string    `json:"highlightedText"`
}

type PatchBlockThread struct {
	ID                 uint64 `json:"id" binding:"required"`
	CanvasRepositoryID uint64 `json:"canvasRepositoryID" binding:"required"`
	CanvasBranchID     uint64 `json:"canvasBranchID" binding:"required"`
	StartBlockID       uint64 `json:"startBlockID" binding:"required"`
	Position           uint   `json:"position" binding:"required"`
	TextRangeStart     uint   `json:"textRangeStart" binding:"required"`
	TextRangeEnd       uint   `json:"textRangeEnd" binding:"required"`
	Text               string `json:"text" binding:"required"`
	HighlightedText    string `json:"highlightedText" binding:"required"`
}
