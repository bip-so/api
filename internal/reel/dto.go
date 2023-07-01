package reel

import (
	"fmt"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

// Get list of all Reels by StudioID
func (r reelRepo) GetPopular(skip, limit int) (*[]models.Reel, error) {
	// @todo: Branch should be public !!! CHirag
	var reels []models.Reel
	var err error
	err = r.db.Table("reels").
		Joins("LEFT JOIN canvas_branches ON reels.canvas_branch_id = canvas_branches.id").
		Where("reels.is_archived = ? and canvas_branches.public_access <> ?", false, "private").
		Preload("CanvasBranch").Preload("CreatedByUser").Preload("Studio").Order("created_at desc").Offset(skip).Limit(limit).Find(&reels).Error
	if err != nil {
		return nil, err
	}
	//publicReels := []models.Reel{}
	//for _, reel := range reels {
	//	if reel.CanvasBranch.PublicAccess != models.PRIVATE {
	//		publicReels = append(publicReels, reel)
	//	}
	//}
	return &reels, nil
}

func (r reelRepo) GetStudioPopular(skip, limit int, studioID uint64) (*[]models.Reel, error) {
	var reels []models.Reel
	var err error
	err = r.db.Table("reels").
		Joins("LEFT JOIN canvas_branches ON reels.canvas_branch_id = canvas_branches.id").
		Where("reels.is_archived = ? and canvas_branches.public_access <> ? and reels.studio_id = ?", false, "private", studioID).
		Preload("CanvasBranch").Preload("CreatedByUser").Preload("Studio").Order("created_at desc").Offset(skip).Limit(limit).Find(&reels).Error
	if err != nil {
		return nil, err
	}
	return &reels, nil
}

func (r reelRepo) Create(reel models.Reel) (*models.Reel, error) {
	err := r.db.Create(&reel).Error
	if err != nil {
		return nil, err
	}
	fmt.Println(reel.ID)
	_ = r.Manager.ReelCountPlus(models.BLOCK, reel.StartBlockID)
	return &reel, nil
}

///////////////////////////////// Reel Comment ///////////////////////
func (r reelRepo) CreateReelComment(data ReelCommentCreatePOST, cb models.CommentBase, reelID uint64, userID uint64) (*models.ReelComment, error) {
	var reelComment models.ReelComment
	reelComment.ReelID = reelID
	reelComment.CreatedByID = userID
	reelComment.Position = cb.Position
	reelComment.Data = cb.Data
	reelComment.IsReply = cb.IsReply
	reelComment.IsEdited = cb.IsEdited
	if data.ParentID != nil && *data.ParentID != 0 {
		reelComment.ParentID = data.ParentID
	}
	reelComment.UpdatedByID = userID

	err := r.db.Create(&reelComment).Error
	if err != nil {
		return nil, err
	}
	// Update comment count on reel when a new comment is created.
	if data.ParentID != nil && *data.ParentID != 0 {
		_ = r.Manager.CommentCountPlus(models.REEL_COMMENTS, *data.ParentID)
	} else {
		_ = r.Manager.CommentCountPlus(models.REEL, reelComment.ReelID)
	}
	//r.db.Model(models.Reel{}).Where("id = ?", reelComment.ReelID).UpdateColumn("comment_count", gorm.Expr("comment_count  + ?", 1))
	return &reelComment, nil
}

// Get list of all Reels by StudioID
// Root Level Comments
func (r reelRepo) GetAllReelComments(studioID uint64, reelID uint64) (*[]models.ReelComment, error) {
	var commendOnReels []models.ReelComment
	err := r.db.Model(&models.ReelComment{}).Where("reel_id = ? and parent_id IS NULL", reelID).Preload("CreatedByUser").Order("created_at desc").Find(&commendOnReels).Error
	if err != nil {
		return nil, err
	}
	return &commendOnReels, nil
}

// Get child commments
func (r reelRepo) GetChildReelComments(parentCommentId uint64) (*[]models.ReelComment, error) {
	var commendOnReels []models.ReelComment
	err := r.db.Model(&models.ReelComment{}).Where("parent_id = ?", parentCommentId).Preload("CreatedByUser").Order("created_at desc").Find(&commendOnReels).Error
	if err != nil {
		return nil, err
	}
	return &commendOnReels, nil
}

// Get list of all Reels by StudioID
// Will be deleted
//func (r reelRepo) GetAll(studioID uint64, withBranchID bool, canvasBranchID uint64) (*[]models.Reel, error) {
//	var reels []models.Reel
//	var err error
//	if withBranchID {
//		err = r.db.Model(&models.Reel{}).Where("studio_id = ? AND canvas_branch_id =?", studioID, canvasBranchID).Order("created_at desc").Find(&reels).Error
//	} else {
//		err = r.db.Model(&models.Reel{}).Where("studio_id = ?", studioID).Order("created_at desc").Find(&reels).Error
//	}
//	if err != nil {
//		return nil, err
//	}
//	return &reels, nil
//}

func (r reelRepo) GetReels(query map[string]interface{}) (*[]models.Reel, error) {
	var reels []models.Reel
	err := postgres.GetDB().Model(&models.Reel{}).Where(query).Preload("CreatedByUser").Preload("Studio").Order("created_at desc").Find(&reels).Error
	if err != nil {
		//log.Fatalln(err)
		return nil, err
	}
	return &reels, nil
}

func (r reelRepo) GetReel(query map[string]interface{}) (*models.Reel, error) {
	var reel models.Reel
	err := postgres.GetDB().Model(&models.Reel{}).Where(query).Preload("CanvasBranch").Preload("Studio").Preload("CreatedByUser").First(&reel).Error
	if err != nil {
		return nil, err
	}
	return &reel, nil
}

func (r reelRepo) GetReelsByIDs(reelIDs []uint64) (*[]models.Reel, error) {
	var reels []models.Reel
	err := postgres.GetDB().Model(&models.Reel{}).Where("id IN ? and is_archived = false", reelIDs).Preload("CreatedByUser").Preload("Studio").Order("created_at desc").Find(&reels).Error
	if err != nil {
		//log.Fatalln(err)
		return nil, err
	}
	return &reels, nil
}

func (r reelRepo) GetReelsByIDsForStudio(reelIDs []uint64, studioID uint64, skip, limit int) ([]models.Reel, error) {
	var reels []models.Reel
	err := postgres.GetDB().Model(&models.Reel{}).Where("id IN ? and studio_id = ?", reelIDs, studioID).Preload("CreatedByUser").Preload("Studio").Order("created_at desc").Offset(skip).Limit(limit).Find(&reels).Error
	if err != nil {
		//log.Fatalln(err)
		return nil, err
	}
	return reels, nil
}

func (r reelRepo) Delete(reelID uint64, userID uint64) error {
	err := r.Manager.SoftDeleteByID(models.REEL, reelID, userID)
	if err != nil {
		return err
	}

	reelInstance, _ := r.GetReel(map[string]interface{}{"id": reelID})
	_ = r.Manager.ReelCountMinus(models.BLOCK, reelInstance.StartBlockID)

	return nil
}

func (r reelRepo) GetMembersByUserID(userID uint64) ([]models.Member, error) {
	var members []models.Member
	result := r.db.Model(&models.Member{}).
		Where("user_id = ? AND has_left = false AND is_removed = false", userID).
		Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

func (r reelRepo) GetUserFollowings(userID uint64) (*[]models.FollowUser, error) {
	var userFollow []models.FollowUser
	err := r.db.Model(&models.FollowUser{}).Where("follower_id = ?", userID).Find(&userFollow).Error
	return &userFollow, err
}

// Get a Block Instance by UUID and Branch
func (r reelRepo) GetBlockByUUIDAndBranchID(query map[string]interface{}) (*models.Block, error) {
	var block models.Block
	err := r.db.Model(&models.Block{}).Where(query).First(&block).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &block, nil
}

func (r reelRepo) GetReelComment(query map[string]interface{}) (*models.ReelComment, error) {
	var reelComment models.ReelComment
	err := r.db.Model(&models.ReelComment{}).Where(query).First(&reelComment).Error
	if err != nil {
		return nil, err
	}
	return &reelComment, nil
}

func (r reelRepo) GetRepo(query map[string]interface{}) (*models.CanvasRepository, error) {
	var repo models.CanvasRepository
	err := r.db.Model(&models.CanvasRepository{}).Where(query).First(&repo).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &repo, nil
}

func (r reelRepo) GetBranchByID(branchID uint64) (*models.CanvasBranch, error) {
	var branch models.CanvasBranch
	err := r.db.Model(&models.CanvasBranch{}).Where("id = ?", branchID).First(&branch).Error
	if err != nil {
		return nil, err
	}
	return &branch, nil
}

func (r reelRepo) GetReelsPopulatedData(query map[string]interface{}) ([]models.Reel, error) {
	var reels []models.Reel
	err := postgres.GetDB().Model(&models.Reel{}).Where(query).Preload("Studio").Preload("CreatedByUser").Preload("CanvasBranch").Order("created_at desc").Find(&reels).Error
	if err != nil {
		//log.Fatalln(err)
		return nil, err
	}
	return reels, nil
}
