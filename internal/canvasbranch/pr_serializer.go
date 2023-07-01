package canvasbranch

import (
	"github.com/google/uuid"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"time"
)

type PublishRequestSerializerDefault struct {
	ID                 uint64    `json:"id"`
	UUID               uuid.UUID `json:"uuid"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
	CreatedByID        uint64    `json:"createdByID"`
	StudioID           uint64    `json:"studioID"`
	CanvasRepositoryID uint64    `json:"canvasRepositoryID"`
	CanvasBranchID     uint64    `json:"canvasBranchID"`
	Status             string    `json:"status"`
	Message            string    `json:"message"`
	ReviewedByUserID   uint64    `json:"reviewedByUserID"`
}

type ManyPublishRequests struct {
	Data []PublishRequestSerializerDefault
}

func PublishRequestSerializerMany(instances *[]models.PublishRequest) ManyPublishRequests {
	var publishRequests []PublishRequestSerializerDefault
	for _, pr := range *instances {
		publishRequests = append(publishRequests, PublishRequestSerializerDefault{
			ID:                 pr.ID,
			UUID:               pr.UUID,
			CreatedAt:          pr.CreatedAt,
			UpdatedAt:          pr.UpdatedAt,
			CreatedByID:        pr.CreatedByID,
			StudioID:           pr.StudioID,
			CanvasRepositoryID: pr.CanvasRepositoryID,
			CanvasBranchID:     pr.CanvasBranchID,
			Status:             pr.Status,
			Message:            pr.Message,
			ReviewedByUserID:   pr.ReviewedByUserID,
		})
	}
	Many := ManyPublishRequests{}
	Many.Data = publishRequests
	return Many
}
