package models

import "gorm.io/datatypes"

const (
	DISCORD_INTEGRATION_TYPE = "discord"
	SLACK_INTEGRATION_TYPE   = "slack"
)

const (
	StudioIntegrationPending = "PENDING"
	StudioIntegrationSuccess = "SUCCESS"
	StudioIntegrationFailed  = "FAILED"
)

type StudioIntegration struct {
	BaseModel

	CreatedByID uint64
	UpdatedByID uint64

	StudioID          uint64
	Type              string
	AccessKey         string
	Status            bool   `gorm:"default:false"`
	IntegrationStatus string `gorm:"type: varchar(20)"`
	Extra             datatypes.JSON
	MessagesData      *datatypes.JSON
	// ChannelsData  *datatypes.JSON
	TeamID        string
	Studio        *Studio `gorm:"foreignKey:StudioID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedByUser *User   `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedByUser *User   `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func newStudioIntegration(integrationType string, studioID uint64, createdByID uint64, accessToken string, extraStr []byte, teamID string) *StudioIntegration {
	return &StudioIntegration{
		Type:        integrationType,
		StudioID:    studioID,
		CreatedByID: createdByID,
		UpdatedByID: createdByID,
		AccessKey:   accessToken,
		Status:      false,
		Extra:       extraStr,
		TeamID:      teamID,
	}
}

func NewDiscordStudioIntegration(studioID uint64, createdByID uint64, accessToken string, extraStr []byte, teamID string) *StudioIntegration {
	return newStudioIntegration(DISCORD_INTEGRATION_TYPE, studioID, createdByID, accessToken, extraStr, teamID)
}

func NewSlackStudioIntegration(studioID uint64, createdByID uint64, accessToken string, extraStr []byte, teamID string) *StudioIntegration {
	return newStudioIntegration(SLACK_INTEGRATION_TYPE, studioID, createdByID, accessToken, extraStr, teamID)
}
