package blockThreadCommentcomment

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"time"
)

func (r blockThreadCommentRepo) Get(query map[string]interface{}) (*models.BlockComment, error) {
	var repo models.BlockComment
	err := r.db.Model(&models.BlockComment{}).Where(query).Preload("CreatedByUser").First(&repo).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &repo, nil
}

func (r blockThreadCommentRepo) GetPreloadThread(query map[string]interface{}) (*models.BlockComment, error) {
	var repo models.BlockComment
	err := r.db.Model(&models.BlockComment{}).Where(query).Preload("Thread").First(&repo).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &repo, nil
}

func (r blockThreadCommentRepo) GetAllComments(query map[string]interface{}) (*[]models.BlockComment, error) {
	var repos []models.BlockComment
	err := r.db.Model(&models.BlockComment{}).Where(query).Preload("CreatedByUser").Order("position ASC").Find(&repos).Error

	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &repos, nil
}

func (r blockThreadCommentRepo) Create(instance *models.BlockComment) (*models.BlockComment, error) {
	results := r.db.Create(&instance)
	// Update comment count on Thread when a new comment is created.
	_ = r.Manager.CommentCountPlus(models.BLOCK_THREAD, instance.ThreadID)
	return instance, results.Error
}

func (r blockThreadCommentRepo) Update(blockThreadCommentID uint64, updates map[string]interface{}) error {
	err := r.db.Model(&models.BlockComment{}).Where("id = ?", blockThreadCommentID).Updates(updates).Error
	if err != nil {
		return err
	}
	return nil
}

func (r blockThreadCommentRepo) Delete(blockThreadCommentID uint64, userId uint64) error {

	err := r.db.Model(&models.BlockComment{}).Where("id = ?", blockThreadCommentID).Updates(map[string]interface{}{
		"is_archived":    true,
		"archived_at":    time.Now(),
		"archived_by_id": userId,
	}).Error
	if err != nil {
		return err
	}

	blockThreadCommentInstance, _ := r.GetPreloadThread(map[string]interface{}{"id": blockThreadCommentID})
	//threadInstance := blockThreadCommentInstance.Thread
	if blockThreadCommentInstance.Thread.CommentCount > 0 {
		_ = r.Manager.CommentCountMinus(models.BLOCK_THREAD, blockThreadCommentInstance.ThreadID)
	}
	return nil
}

func (r blockThreadCommentRepo) GetBlockThread(query map[string]interface{}) (*models.BlockThread, error) {
	var repo models.BlockThread
	err := r.db.Model(&models.BlockThread{}).Where(query).Preload("CreatedByUser").First(&repo).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &repo, nil
}
