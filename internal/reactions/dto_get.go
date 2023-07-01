package reactions

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

func (r reactionRepo) GetBlockReaction(query map[string]interface{}) (*models.BlockReaction, error) {
	var reaction *models.BlockReaction
	err := postgres.GetDB().Model(&models.BlockReaction{}).Where(query).First(&reaction).Error

	return reaction, err
}

func (r reactionRepo) GetBlockReactions(query map[string]interface{}) (*[]models.BlockReaction, error) {
	var reactions *[]models.BlockReaction
	err := postgres.GetDB().Model(&models.BlockReaction{}).Where(query).Find(&reactions).Error
	return reactions, err
}

func (r reactionRepo) GetBlockThreadReaction(query map[string]interface{}) ([]models.BlockThreadReaction, error) {
	var reaction []models.BlockThreadReaction
	err := postgres.GetDB().Model(&models.BlockThreadReaction{}).Where(query).Find(&reaction).Error

	return reaction, err
}

func (r reactionRepo) GetBlockThreadCommentReaction(query map[string]interface{}) ([]models.BlockCommentReaction, error) {
	var reaction []models.BlockCommentReaction
	err := postgres.GetDB().Model(&models.BlockCommentReaction{}).Where(query).Find(&reaction).Error

	return reaction, err
}

func (r reactionRepo) GetReelReaction(query map[string]interface{}) ([]models.ReelReaction, error) {
	var reaction []models.ReelReaction
	err := postgres.GetDB().Model(&models.ReelReaction{}).Where(query).Find(&reaction).Error

	return reaction, err
}

func (r reactionRepo) GetOneReelReaction(query map[string]interface{}) ([]models.ReelReaction, error) {
	var reaction []models.ReelReaction
	err := postgres.GetDB().Model(&models.ReelReaction{}).Where(query).Find(&reaction).Error

	return reaction, err
}

func (r reactionRepo) GetReelReactionByIDs(reelIDs []uint64) ([]models.ReelReaction, error) {
	var reaction []models.ReelReaction
	err := r.db.Model(&models.ReelReaction{}).Where("reel_id IN ?", reelIDs).Preload("CreatedByUser").Find(&reaction).Error
	return reaction, err
}

func (r reactionRepo) GetUserReelReactionByIDs(reelIDs []uint64, userID uint64) ([]models.ReelReaction, error) {
	var reaction []models.ReelReaction
	err := r.db.Model(&models.ReelReaction{}).Where("reel_id IN ? and created_by_id = ?", reelIDs, userID).Preload("CreatedByUser").Find(&reaction).Error
	return reaction, err
}

func (r reactionRepo) GetReelCommentReaction(query map[string]interface{}) ([]models.ReelCommentReaction, error) {
	var reaction []models.ReelCommentReaction
	err := postgres.GetDB().Model(&models.ReelCommentReaction{}).Where(query).Find(&reaction).Error

	return reaction, err
}

func (r reactionRepo) GetBlockByID(blockId uint64) (models.Block, error) {
	var block models.Block
	err := r.db.Model(models.Block{}).Where("id = ?", blockId).First(&block).Error
	return block, err
}

func (r reactionRepo) GetReelByID(reelID uint64) (*models.Reel, error) {
	var reel models.Reel
	err := postgres.GetDB().Model(&models.Reel{}).Where("id = ?", reelID).First(&reel).Error
	if err != nil {
		return nil, err
	}
	return &reel, nil
}

func (r reactionRepo) UpdateBranchLastEdited(branchID uint64) error {
	err := r.db.Model(models.CanvasBranch{}).Where("id = ?", branchID).Updates(map[string]interface{}{}).Error
	return err
}
