package studio

import (
	"time"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

type StudioDocument struct {
	ID                    uint64            `json:"id"`
	UUID                  string            `json:"uuid"`
	ObjectID              string            `json:"objectID"`
	CreatedAt             time.Time         `json:"createdAt"`
	UpdatedAt             time.Time         `json:"updatedAt"`
	DisplayName           string            `json:"displayName"`
	Handle                string            `json:"handle"`
	Description           string            `json:"description"`
	ImageURL              string            `json:"imageUrl"`
	IsJoined              *bool             `json:"isJoined"`
	IsRequested           bool              `json:"isRequested"`
	AllowPublicMembership bool              `json:"allowPublicMembership"`
	MembersCount          int64             `json:"membersCount"`
	Topics                []TopicSerializer `json:"topics"`
	CreatedByID           uint64            `json:"createdById"`
}

func StudioModelToStudioDocument(stdio *models.Studio) *StudioDocument {
	return &StudioDocument{
		ID:                    stdio.ID,
		UUID:                  stdio.UUID.String(),
		ObjectID:              stdio.UUID.String(),
		CreatedAt:             stdio.CreatedAt,
		UpdatedAt:             stdio.UpdatedAt,
		DisplayName:           stdio.DisplayName,
		Handle:                stdio.Handle,
		Description:           stdio.Description,
		ImageURL:              stdio.ImageURL,
		IsJoined:              nil,
		Topics:                SerializeTopicData(stdio.Topics),
		AllowPublicMembership: stdio.AllowPublicMembership,
		CreatedByID:           stdio.CreatedByID,
	}
}
