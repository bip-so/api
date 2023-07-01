package queries

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

// publishRequestQuery

func (r publishRequestQuery) DeletePR(id uint64) error {
	err := postgres.GetDB().Delete(&models.PublishRequest{}, "id = ?", id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r publishRequestQuery) CreatePublishRequest(instance *models.PublishRequest) (*models.PublishRequest, error) {
	results := postgres.GetDB().Create(&instance)
	return instance, results.Error
}

func (r publishRequestQuery) GetAllPublishRequests(query map[string]interface{}) (*[]models.PublishRequest, error) {
	var instances []models.PublishRequest
	err := postgres.GetDB().Model(&models.PublishRequest{}).Where(query).Find(&instances).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &instances, nil
}

func (r publishRequestQuery) UpdatePublishRequest(id uint64, query map[string]interface{}) error {
	//var instances []models.PublishRequest
	err := postgres.GetDB().Model(&models.PublishRequest{}).Where("id = ?", id).Updates(query).Error
	if err != nil {
		logger.Debug(err.Error())
		return err
	}
	return nil
}

func (r publishRequestQuery) PublishRequestGetter(branchID uint64, userID uint64, prID uint64) (*models.PublishRequest, error) {
	var instances models.PublishRequest
	err := postgres.GetDB().Model(&models.PublishRequest{}).Where("id = ?", prID).First(&instances).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &instances, nil
}

func (r publishRequestQuery) PublishRequestExists(branchID uint64, userID uint64) bool {
	var count int64
	// Added status query condition to return only active publish requests.
	_ = postgres.GetDB().Model(&models.PublishRequest{}).Where("canvas_branch_id = ? and created_by_id = ? and status <> ?", branchID, userID, models.PUBLISH_REQUEST_REJECTED).Count(&count).Error
	if count == 0 {
		return false
	}
	return true
}
