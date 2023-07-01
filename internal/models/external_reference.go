package models

import "github.com/lib/pq"

type ExternalReference struct {
	ID                 string `gorm:"primaryKey"`
	ExternalID         string
	ExternalSourceType string         //e.g. slack, discord, linkedin etc.
	InternalID         uint64         // .e.g. acitivtyId
	InternalObjectType string         // e.g. reel , comment etc.
	Extra              pq.StringArray `gorm:"type:text[]"`
}
