package notifications

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gorm.io/datatypes"
)

func (r notificationRepo) GetStudioIntegrations(studioId uint64) ([]models.StudioIntegration, error) {
	var integrations []models.StudioIntegration
	err := r.db.Model(models.StudioIntegration{}).Where("studio_id = ?", studioId).Find(&integrations).Error
	if err != nil {
		return nil, err
	}
	return integrations, nil
}

func (r notificationRepo) AddIntegrationReference(externalID string, externalSource string, internalID uint64, internalSource, extra string) error {
	integrationReference := &models.IntegrationReference{
		ExternalID:         externalID,
		ExternalSourceType: externalSource,
		InternalID:         internalID,
		InternalObjectType: internalSource,
		Extra:              datatypes.JSON(extra),
	}
	err := r.db.Model(models.IntegrationReference{}).Create(integrationReference).Error
	return err
}
