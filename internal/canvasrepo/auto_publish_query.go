package canvasrepo

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

//

func UpdateCanvasLanguageBranchesVisibility(branchID uint64, userID uint64, visibility string) error {
	canvasBranch, _ := GetCanvasBranch(map[string]interface{}{"id": branchID})
	canvasRepos, _ := GetCanvasRepos(map[string]interface{}{"default_language_canvas_repo_id": canvasBranch.CanvasRepositoryID})
	for _, repo := range *canvasRepos {
		UpdateBranchInstance(*repo.DefaultBranchID, map[string]interface{}{"updated_by_id": userID, "public_access": visibility})
	}
	return nil
}

// UpdateCanvasBranchVisibility Start
func UpdateCanvasBranchVisibility(branchID uint64, userID uint64, visibility string) error {
	err := UpdateBranchInstance(branchID, map[string]interface{}{"updated_by_id": userID, "public_access": visibility})
	if err != nil {
		return err
	}
	return nil
}

func UpdateBranchInstance(branchID uint64, query map[string]interface{}) error {
	err := postgres.GetDB().Model(&models.CanvasBranch{}).Where("id = ?", branchID).Updates(query).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateCanvasBranchVisibility ENd

// UpdateCanvasLanguageBranchesVisibility Start
func GetCanvasBranch(query map[string]interface{}) (*models.CanvasBranch, error) {
	var branch models.CanvasBranch
	err := postgres.GetDB().Model(&models.CanvasBranch{}).Where(query).Preload("CanvasRepository").Preload("CanvasRepository.Studio").First(&branch).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &branch, nil
}

func GetCanvasRepos(query map[string]interface{}) (*[]models.CanvasRepository, error) {
	var canvasRepos []models.CanvasRepository
	err := postgres.GetDB().Model(&models.CanvasRepository{}).Where(query).Find(&canvasRepos).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &canvasRepos, nil
}
