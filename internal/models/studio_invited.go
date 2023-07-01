package models

import (
	"github.com/lib/pq"
)

type StudioInviteViaEmail struct {
	BaseModel
	Email           string
	PermissionGroup string // Deprecated
	CreatedByID     uint64
	StudioID        uint64
	// Roles           []uint64       `gorm:"type:integer[]"`
	Roles2        pq.StringArray `gorm:"type:text[]"`
	HasAccepted   bool
	CreatedByUser *User   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Studio        *Studio `gorm:"foreignKey:StudioID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
