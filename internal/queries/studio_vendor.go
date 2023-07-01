package queries

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

func (q *studioVendorQuery) EmptyStudioVendorObject() *models.StudioVendor {
	return &models.StudioVendor{}
}

func (q *studioVendorQuery) CreateStudioVendor(guildId string, guildName string, partnerName string) (*models.StudioVendor, error) {
	sv := q.EmptyStudioVendorObject()
	sv.GuildId = guildId
	sv.IntegrationStatus = "Started"
	sv.PartnerName = partnerName
	sv.GuildName = guildName
	results := postgres.GetDB().Create(&sv)
	return sv, results.Error
}

func (m *studioVendorQuery) GetStudioVendor(guildId string) (*models.StudioVendor, error) {
	var vendor *models.StudioVendor
	postgres.GetDB().Model(&models.StudioVendor{}).Where("guild_id = ?", guildId).First(&vendor)
	return vendor, nil
}
