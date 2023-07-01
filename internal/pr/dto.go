package pr

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
)

func (r prRepo) Create(instance *models.PublishRequest) (*models.PublishRequest, error) {
	results := r.db.Create(&instance)
	return instance, results.Error
}

func (r prRepo) UpdatePRInstance(prID uint64, query map[string]interface{}) error {
	err := r.db.Model(&models.CanvasBranch{}).Where("id = ?", prID).Updates(query).Error
	if err != nil {
		return err
	}
	return nil
}

func (r prRepo) GetPublishRequestsByStudio(query map[string]interface{}) (*[]models.PublishRequest, error) {
	fmt.Println(query)
	var instances []models.PublishRequest
	err := r.db.Model(&models.PublishRequest{}).Where(query).Find(&instances).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &instances, nil
}
