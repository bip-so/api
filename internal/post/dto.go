package post

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/apiutil"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r postRepo) Create(instance *models.Post) (*models.Post, error) {
	results := r.db.Create(&instance)
	return instance, results.Error
}

func (r postRepo) GetOne(instance *models.Post) (*models.Post, error) {
	//results := r.db.Create(&instance)
	//return instance, results.Error
	return nil, nil
}

func (r postRepo) GetAllPosts(studioID uint64, page int) (*[]models.Post, error) {
	var posts *[]models.Post
	// Pagination
	perPage := apiutil.SharedPaginationPerPage
	offset := (page - 1) * perPage

	err := r.db.Model(&models.Post{}).
		Where("studio_id = ?", studioID).
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Studio").
		Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r postRepo) GetAllHomedPagePosts(studioIdArray []uint64, page int) (*[]models.Post, error) {
	//var posts2 *[]models.Post
	//r.db.Joins("")
	// Pagination
	perPage := apiutil.SharedPaginationPerPage
	offset := (page - 1) * perPage
	var posts *[]models.Post
	err := r.db.Model(&models.Post{}).
		Where("studio_id in ?", studioIdArray).
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Studio").
		Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&posts).Error

	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r postRepo) GetAllPostComments(id uint64) (*[]models.PostComment, error) {
	var comments *[]models.PostComment
	err := r.db.Model(&models.PostComment{}).Where("post_id = ?", id).Preload("CreatedByUser").Preload("UpdatedByUser").Order("created_at DESC").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r postRepo) GetRole(roleId uint64) (*models.Role, error) {
	var role *models.Role
	postgres.GetDB().Model(&models.Role{}).Where("id = ?", roleId).First(&role)
	return role, nil
}

func (r postRepo) GetSinglePost(postID uint64) (*models.Post, error) {
	var post *models.Post
	err := r.db.Model(&models.Post{}).Where("id = ?", postID).Preload("CreatedByUser").Preload("UpdatedByUser").Preload("Studio").First(&post).Error
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (r postRepo) PostDelete(id uint64) error {
	postObject := models.Post{BaseModel: models.BaseModel{ID: id}}
	err := postgres.GetDB().Unscoped().Select(clause.Associations).Delete(postObject).Error
	return err
}

func (r postRepo) PostCommentDelete(id uint64) error {
	var postCommentInstance *models.PostComment
	err1 := r.db.Model(&models.PostComment{}).Where("id = ?", id).First(&postCommentInstance).Error
	if err1 != nil {
		return nil
	}
	// There is an edge case here we can handle later.
	_ = r.db.Table("posts").Where("id", postCommentInstance.PostID).UpdateColumn("comment_count", gorm.Expr("comment_count  - ?", 1)).Error

	// We also need to reduce the comment count if parent exists
	if postCommentInstance.ParentPostCommentID != nil {
		_ = r.db.Table("post_comments").Where("id", postCommentInstance.ParentPostCommentID).UpdateColumn("comment_count", gorm.Expr("comment_count  - ?", 1)).Error
	}

	postCommentObject := models.PostComment{BaseModel: models.BaseModel{ID: id}}
	err := postgres.GetDB().Unscoped().Select(clause.Associations).Delete(postCommentObject).Error
	return err
}

func (r postRepo) GetPostReactions(query map[string]interface{}) (*[]models.PostReaction, error) {
	var reactions *[]models.PostReaction
	err := postgres.GetDB().Model(&models.PostReaction{}).Where(query).Find(&reactions).Error
	return reactions, err
}
