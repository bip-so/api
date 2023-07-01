package models

type Attribution struct {
	CollectionID   string
	CanvasRepoID   uint64
	CanvasBranchID uint64
	UserID         uint64
	Edits          int

	User         User             `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CanvasRepo   CanvasRepository `gorm:"foreignKey:CanvasRepoID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CanvasBranch CanvasBranch     `gorm:"foreignKey:CanvasBranchID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func NewAttribution(canvasRepositoryID, canvasBranchID, userID uint64, edits int) *Attribution {
	return &Attribution{
		CanvasRepoID:   canvasRepositoryID,
		CanvasBranchID: canvasBranchID,
		UserID:         userID,
		Edits:          edits,
	}
}
