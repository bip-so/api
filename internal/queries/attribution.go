package queries

// attributionQuery

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

func (r attributionQuery) GetAllAttributionsForBranch(branchID uint64) ([]models.Attribution, error) {
	var attributions []models.Attribution
	err := postgres.GetDB().Where(models.Attribution{CanvasBranchID: branchID}).Preload("User").Find(&attributions).Error
	return attributions, err
}

func (r attributionQuery) ReplaceOrCreateNewAttributions(attrs []models.Attribution, branchID uint64) error {
	err := r.DeleteAllAttributionsForBranch(branchID)
	if err != nil {
		return err
	}
	return r.CreateNewAttributions(attrs)
}

func (r attributionQuery) DeleteAllAttributionsForBranch(canvasbranchID uint64) error {
	return postgres.GetDB().Where(models.Attribution{CanvasBranchID: canvasbranchID}).Delete(models.Attribution{}).Error
}

func (r attributionQuery) CreateNewAttributions(attrs []models.Attribution) error {
	return postgres.GetDB().Create(&attrs).Error
}
