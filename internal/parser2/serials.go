package parser2

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"mime/multipart"
)

type ImportNotionValidator struct {
	File multipart.FileHeader `json:"file" form:"file"`
}

type ImportTask struct {
	File     multipart.FileHeader `json:"file,omitempty" form:"file,omitempty"`
	StudioID uint64               `json:"studioId,omitempty" form:"studioId,omitempty"`
	User     *models.User         `json:"user,omitempty" form:"user,omitempty"`
}
