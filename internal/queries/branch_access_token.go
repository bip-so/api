package queries

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

//branchAccessTokenQuery

func (r branchAccessTokenQuery) GetBranchTokenInstance(query map[string]interface{}) (*models.BranchAccessToken, error) {
	var instance *models.BranchAccessToken
	err := postgres.GetDB().Model(&models.BranchAccessToken{}).Where(query).First(&instance).Error
	return instance, err
}

func (r branchAccessTokenQuery) GetAllBranchTokenInstance(query map[string]interface{}) (*[]models.BranchAccessToken, error) {
	var instances []models.BranchAccessToken
	err := postgres.GetDB().Model(&models.BranchAccessToken{}).Where(query).Find(&instances).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &instances, nil
}

func (r branchAccessTokenQuery) DeleteBranchToken(id uint64) error {
	err := r.Manager.HardDeleteByID(models.BRANCH_ACCESS_TOKEN, id)
	if err != nil {
		return err
	}
	return nil
}

func (r branchAccessTokenQuery) CreateBranchTokenInstance(instance *models.BranchAccessToken) (*models.BranchAccessToken, error) {
	results := postgres.GetDB().Create(&instance)
	return instance, results.Error
}
