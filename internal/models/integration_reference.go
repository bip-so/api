package models

import (
	"gorm.io/datatypes"
)

type IntegrationReference struct {
	BaseModel
	ExternalID         string
	ExternalSourceType string         //e.g. slack, discord, linkedin etc.
	InternalID         uint64         // .e.g. acitivtyId
	InternalObjectType string         // e.g. reel , comment etc.
	Extra              datatypes.JSON `gorm:"type:text[]"`
}
