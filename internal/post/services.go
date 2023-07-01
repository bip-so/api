package post

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gorm.io/gorm"
)

// get all posts from this studio
func (s postService) GetPostServices(studioID uint64, page int) (*[]models.Post, error) {
	return App.Repo.GetAllPosts(studioID, page)
}

func (s postService) GetStudioIDArrayUserID(userID uint64) []uint64 {
	return s.Manager.GetAllStudioIDsByUserID(userID)
}

func (s postService) GetHomepagePostServices(userID uint64, page int) (*[]models.Post, error) {
	studioIdArray := s.Manager.GetAllStudioIDsByUserID(userID)
	return App.Repo.GetAllHomedPagePosts(studioIdArray, page)
}
func (s postService) GetSinglePostService(postID uint64) (*models.Post, error) {
	return App.Repo.GetSinglePost(postID)
}

func (s postService) DeleteSinglePostService(postID uint64) error {
	return App.Repo.PostDelete(postID)
}

func (s postService) DeleteSinglePostCommentService(commentPostID uint64) error {
	return App.Repo.PostCommentDelete(commentPostID)
}

func (s postService) UpdatePostServices(id uint64, update map[string]interface{}) error {
	result := s.db.Model(&models.Post{}).Where("id = ?", id).Updates(update)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s postService) CreatePostComment(postID uint64, data CreatePostCommentValidation, currentUser *models.User) (*models.PostComment, error) {
	comment := &models.PostComment{}
	comment.Comment = data.Comment
	comment.IsEdited = false
	comment.PostID = postID
	comment.CreatedByID = currentUser.ID
	comment.UpdatedByID = currentUser.ID
	if data.ParentPostCommentID != 0 {
		comment.ParentPostCommentID = &data.ParentPostCommentID
	}

	err := s.db.Model(&models.PostComment{}).Create(comment).Error

	// We need to add to the comment counter : Increment the CommentCount on Post Model
	_ = s.db.Table("posts").Where("id", postID).UpdateColumn("comment_count", gorm.Expr("comment_count  + ?", 1)).Error
	// We also need to update the comment count on the PostComment on Parent Comment
	// If the parent comment is non-zero
	if data.ParentPostCommentID != 0 {
		_ = s.db.Table("post_comments").Where("id", &data.ParentPostCommentID).UpdateColumn("comment_count", gorm.Expr("comment_count  + ?", 1)).Error
	}

	return comment, err
}

func (s postService) GetPostCommentsServices(postID uint64) (*[]models.PostComment, error) {
	return App.Repo.GetAllPostComments(postID)
}

func (s postService) UpdatePostComment(id uint64, update map[string]interface{}) error {
	result := s.db.Model(&models.PostComment{}).Where("id = ?", id).Updates(update)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s postService) CreatePostReaction(postID uint64, data NewPostReaction, currentUser *models.User) error {
	// We need to check for Duplicate here for this USEta
	// Check Duplicate
	reactionDuplicate := &models.PostReaction{}
	// This user, this Emoji and this PostID
	err2 := s.db.Model(models.PostReaction{}).Where("post_id = ? and created_by_id = ? and emoji = ?", postID, currentUser.ID, data.Emoji).First(&reactionDuplicate).Error
	if err2 == gorm.ErrRecordNotFound {
		// Create Reaction
		reaction := &models.PostReaction{}
		reaction.CreatedByID = currentUser.ID
		reaction.PostID = postID
		reaction.Emoji = data.Emoji
		result := s.db.Model(&models.PostReaction{}).Create(&reaction).Error
		return result
	}
	return err2
}

func (s postService) CreatePostCommentReaction(postID uint64, postCommentID uint64, data NewPostCommentReaction, currentUser *models.User) error {
	// We need to check for Duplicate here for this USEta
	// Check Duplicate
	reactionDuplicate := &models.PostCommentReaction{}
	// This user, this Emoji and this PostID
	err2 := s.db.Model(models.PostCommentReaction{}).Where("post_id = ? and post_comment_id = ? and created_by_id = ? and emoji = ?", postID, postCommentID, currentUser.ID, data.Emoji).First(&reactionDuplicate).Error
	if err2 == gorm.ErrRecordNotFound {
		// Create Reaction
		reaction := &models.PostCommentReaction{}
		reaction.CreatedByID = currentUser.ID
		reaction.PostID = postID
		reaction.PostCommentID = postCommentID
		reaction.Emoji = data.Emoji
		result := s.db.Model(&models.PostCommentReaction{}).Create(&reaction).Error
		return result
	}
	return err2
}

func (s postService) RemovePostReaction(postID uint64, data RemovePostReaction, currentUser *models.User) error {
	reaction := &models.PostReaction{}
	err2 := s.db.Model(models.PostReaction{}).Where("post_id = ? and created_by_id = ? and emoji = ?", postID, currentUser.ID, data.Emoji).Delete(&reaction).Error
	return err2
}

func (s postService) RemovePostCommentReaction(postID uint64, postCommentID uint64, data RemovePostCommentReaction, currentUser *models.User) error {
	reaction := &models.PostCommentReaction{}
	err2 := s.db.Model(models.PostCommentReaction{}).Where("post_id = ? and post_comment_id = ? and created_by_id = ? and emoji = ?", postID, postCommentID, currentUser.ID, data.Emoji).Delete(&reaction).Error
	return err2
}

func (s postService) InvalidateStudioPostCache(studioID uint64) {
	BadInvalidationOfStudioPosts(studioID)
}
