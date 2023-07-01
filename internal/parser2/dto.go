package parser2

import "gitlab.com/phonepost/bip-be-platform/internal/models"

func (r parser2Repo) Get(branchID uint64) ([]models.Block, error) {
	var blocks []models.Block
	err := r.db.Model(&models.Block{}).Where("canvas_branch_id = ?", branchID).Order("rank").Find(&blocks).Error
	if err != nil {
		return nil, err
	}
	return blocks, nil
}

func (r parser2Repo) GetCanvasRepos(repoID uint64) ([]models.CanvasRepository, error) {
	var repos []models.CanvasRepository
	err := r.db.Model(&models.CanvasRepository{}).Where("parent_canvas_repository_id = ?", repoID).Order("position").Find(&repos).Error
	if err != nil {
		return nil, err
	}
	return repos, nil
}

func (r parser2Repo) GetCanvasRepo(query map[string]interface{}) (*models.CanvasRepository, error) {
	var repo *models.CanvasRepository
	err := r.db.Model(&models.CanvasRepository{}).Where(query).First(&repo).Error
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (r parser2Repo) GetBranchByID(branchID uint64) (*models.CanvasBranch, error) {
	var branch models.CanvasBranch
	err := r.db.Model(&models.CanvasBranch{}).Where("id = ?", branchID).Preload("CanvasRepository").Preload("CanvasRepository.Studio").First(&branch).Error
	if err != nil {
		return nil, err
	}
	return &branch, nil
}

func (r parser2Repo) GetStudioByID(studioID uint64) (*models.Studio, error) {
	var studio models.Studio
	err := r.db.Model(models.Studio{}).Where("id = ?", studioID).First(&studio).Error
	return &studio, err
}
