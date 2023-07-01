package mentions

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
)

func (r mentionsRepo) GetRepo(query map[string]interface{}) (*models.CanvasRepository, error) {
	var repo models.CanvasRepository
	err := r.db.Model(&models.CanvasRepository{}).Where(query).Preload("DefaultBranch").First(&repo).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &repo, nil
}
