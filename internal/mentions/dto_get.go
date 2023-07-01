package mentions

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
)

func (r mentionsRepo) GetBlock(id uint64) (*models.Block, error) {
	var instance models.Block
	err := r.db.Model(&models.Block{}).Where("id = ?", id).First(&instance).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &instance, nil
}
func (r mentionsRepo) GetBlockByUUID(uuid string) (*models.Block, error) {
	var instance models.Block
	err := r.db.Model(&models.Block{}).Where("uuid = ?", uuid).First(&instance).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &instance, nil
}

func (r mentionsRepo) GetBlockThread(id uint64) (*models.BlockThread, error) {
	var instance models.BlockThread
	err := r.db.Model(&models.BlockThread{}).Where("id = ?", id).First(&instance).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &instance, nil
}

func (r mentionsRepo) GetBlockComment(id uint64) (*models.BlockComment, error) {
	var instance models.BlockComment
	err := r.db.Model(&models.BlockComment{}).Where("id = ?", id).Preload("Thread").First(&instance).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &instance, nil
}

func (r mentionsRepo) GetReel(id uint64) (*models.Reel, error) {
	var instance models.Reel
	err := r.db.Model(&models.Reel{}).Where("id = ?", id).First(&instance).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &instance, nil
}

func (r mentionsRepo) GetReelComment(id uint64) (*models.ReelComment, error) {
	var instance models.ReelComment
	err := r.db.Model(&models.ReelComment{}).Where("id = ?", id).Preload("Reel").First(&instance).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &instance, nil
}
