package blockthread

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"time"
)

func (r blockThreadRepo) Get(query map[string]interface{}) (*models.BlockThread, error) {
	var repo models.BlockThread
	err := r.db.Model(&models.BlockThread{}).Where(query).Preload("CreatedByUser").First(&repo).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &repo, nil
}
func (r blockThreadRepo) GetAllThread(query map[string]interface{}) (*[]models.BlockThread, error) {
	var repos []models.BlockThread
	err := r.db.Model(&models.BlockThread{}).Where(query).Preload("CreatedByUser").Order("position DESC").Find(&repos).Error

	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &repos, nil
}

func (r blockThreadRepo) GetAllReels(query map[string]interface{}) ([]models.Reel, error) {
	var reels []models.Reel
	err := r.db.Model(&models.Reel{}).Where(query).Find(&reels).Error

	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return reels, nil
}

func (r blockThreadRepo) Create(instance *models.BlockThread) (*models.BlockThread, error) {
	results := r.db.Create(&instance)
	fmt.Println("Block ID", instance.StartBlockID)
	_ = r.Manager.CommentCountPlus(models.BLOCK, instance.StartBlockID)
	return instance, results.Error
}

func (r blockThreadRepo) Update(blockThreadID uint64, updates map[string]interface{}) error {
	err := r.db.Model(&models.BlockThread{}).Where("id = ?", blockThreadID).Updates(updates).Error
	if err != nil {
		return err
	}
	return nil
}

func (r blockThreadRepo) Delete(blockThreadID uint64, userId uint64) error {
	err := r.db.Model(&models.BlockThread{}).Where("id = ?", blockThreadID).Updates(map[string]interface{}{
		"is_archived":    true,
		"archived_at":    time.Now(),
		"archived_by_id": userId,
	}).Error

	if err != nil {
		return err
	}
	blockThreadInstance, _ := r.Get(map[string]interface{}{"id": blockThreadID})
	if !blockThreadInstance.IsResolved {
		_ = r.Manager.CommentCountMinus(models.BLOCK, blockThreadInstance.StartBlockID)
	}
	return nil
}

func (r blockThreadRepo) Resolve(blockThreadID uint64, userId uint64) error {
	err := r.db.Model(&models.BlockThread{}).Where("id = ?", blockThreadID).Updates(map[string]interface{}{
		"is_resolved":    true,
		"resolved_at":    time.Now(),
		"resolved_by_id": userId,
	}).Error
	if err != nil {
		return err
	}

	blockThreadInstance, _ := r.Get(map[string]interface{}{"id": blockThreadID})
	_ = r.Manager.CommentCountMinus(models.BLOCK, blockThreadInstance.StartBlockID)

	return nil
}

func (r blockThreadRepo) GetBranchAndRepPreload(query map[string]interface{}) (*models.CanvasBranch, error) {
	var branch models.CanvasBranch
	err := r.db.Model(&models.CanvasBranch{}).Where(query).Preload("CanvasRepository").First(&branch).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &branch, nil
}

// get all rough branches from this branch
func (r blockThreadRepo) AllRoughBranchesForGivenBranch(mainBranchID uint64) []uint64 {
	var ids []uint64
	r.db.Model(&models.CanvasBranch{}).Where("rough_from_branch_id = ?", mainBranchID).Pluck("id", &ids)
	return ids
}

// get and delete BlockThread instances on this branch (branchID, nodeCommentID)
func (r blockThreadRepo) DeleteClonedCommentOnRoughBranch(branchID uint64, nodeCommentID uint64) {
	err := r.db.Model(models.BlockThread{}).Delete(map[string]interface{}{"canvas_branch_id": branchID, "cloned_from_thread": nodeCommentID}).Error
	if err != nil {
		fmt.Println(err.Error())
	}
}
