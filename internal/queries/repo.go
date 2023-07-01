package queries

import (
	"encoding/json"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"strconv"
)

func (q repoQuery) EmptyRepoObject() *models.CanvasRepository {
	return &models.CanvasRepository{}
}

func (q repoQuery) CreateRepo(parentCanvasRepositoryID uint64, userID uint64, studioID uint64, collectionID uint64, name string, icon string, position uint) (*models.CanvasRepository, error) {
	instance := q.EmptyRepoObject()
	instance.CreatedByID = userID
	instance.UpdatedByID = userID
	instance.Name = name
	instance.Icon = icon
	instance.CollectionID = collectionID
	instance.StudioID = studioID
	instance.Position = position
	if parentCanvasRepositoryID != 0 {
		instance.ParentCanvasRepositoryID = &parentCanvasRepositoryID
	}
	instance.IsPublished = false
	instance.DefaultBranchID = nil // Generated later
	instance.Key = utils.NewNanoid()
	results := postgres.GetDB().Create(&instance)
	go func() {
		canvasRepo, _ := json.Marshal(instance)
		q.kafka.Publish(configs.KAFKA_TOPICS_NEW_CANVAS, strconv.FormatUint(instance.ID, 10), canvasRepo)
	}()
	return instance, results.Error
}

//
func (q repoQuery) UpdateRepo(repoID uint64, updates map[string]interface{}) (*models.CanvasRepository, error) {
	var instance *models.CanvasRepository
	err := postgres.GetDB().Model(models.CanvasRepository{}).Where(models.CanvasRepository{BaseModel: models.BaseModel{ID: repoID}}).Updates(updates).Find(&instance).Error
	return instance, err
}

// branchId uint64
func (q repoQuery) GetRepo(query map[string]interface{}) (*models.CanvasRepository, error) {
	var repo models.CanvasRepository
	err := postgres.GetDB().Model(&models.CanvasRepository{}).Where(query).Preload("Studio").Preload("DefaultBranch").First(&repo).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &repo, nil
}

func (q repoQuery) GetReposLanguages(studioID uint64) ([]models.CanvasRepository, error) {
	var repo []models.CanvasRepository
	err := postgres.GetDB().Model(&models.CanvasRepository{}).
		Joins("LEFT OUTER JOIN canvas_branches on canvas_repositories.id = canvas_branches.canvas_repository_id").
		Where("studio_id = ? and is_published = true and canvas_repositories.is_archived = false and canvas_branches.public_access <> 'private'", studioID).
		Distinct("language").
		Find(&repo).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return repo, nil
}
