package models

const (
	PUBLISH_REQUEST_PENDING  = "PENDING"
	PUBLISH_REQUEST_ACCEPTED = "ACCEPTED"
	PUBLISH_REQUEST_REJECTED = "REJECTED"
)

type PublishRequest struct {
	BaseModel
	StudioID           uint64
	CanvasRepositoryID uint64
	CanvasBranchID     uint64
	Status             string
	Message            string
	ReviewedByUserID   uint64
	CreatedByID        uint64
	CreatedByUser      *User             `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ReviewedByUser     *User             `gorm:"foreignKey:ReviewedByUserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Studio             *Studio           `gorm:"foreignKey:StudioID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CanvasRepository   *CanvasRepository `gorm:"foreignKey:CanvasRepositoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CanvasBranch       *CanvasBranch     `gorm:"foreignkey:CanvasBranchID;constraint:OnDelete:CASCADE;"`
}

func (m PublishRequest) NewPublishRequest(studioID uint64, canvasRepositoryID uint64, canvasBranchID uint64, status string, message string, createdByID uint64, reviewedByUserID uint64) *PublishRequest {
	return &PublishRequest{StudioID: studioID, CanvasRepositoryID: canvasRepositoryID, CanvasBranchID: canvasBranchID, Status: status, Message: message, CreatedByID: createdByID, ReviewedByUserID: reviewedByUserID}
}
