package post

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

func (c postController) CreatePost(postRequest NewPostThread, studioID uint64, authUser models.User) (*models.Post, error) {
	newPost := &models.Post{}
	var roles []models.Role
	fmt.Println("StudioID", studioID)

	newPost.StudioID = studioID
	newPost.CreatedByID = authUser.ID
	newPost.UpdatedByID = authUser.ID
	newPost.Children = postRequest.Children
	newPost.Attributes = postRequest.Attributes
	newPost.IsPublic = postRequest.IsPublic
	// Looping the role ID create newPost (Roles) instance
	for _, roleID := range postRequest.Roles {
		role, _ := App.Repo.GetRole(roleID)
		roles = append(roles, *role)
	}
	newPost.Roles = roles
	postInstance, err := App.Repo.Create(newPost)
	if err != nil {
		return nil, err
	}
	return postInstance, nil
}

func (c postController) UpdatePost(postID uint64, postRequest UpdatePostThread, authUser models.User) error {
	return App.Service.UpdatePostServices(postID, map[string]interface{}{
		"children":      postRequest.Children,
		"attributes":    postRequest.Attributes,
		"is_public":     postRequest.IsPublic,
		"updated_by_id": authUser.ID,
	})
}

func (c postController) GetAllPosts(studioID uint64, page int) (*[]models.Post, error) {
	// if page is 1
	// we'll try and get from cache
	// if cache is found send that
	// else
	return App.Service.GetPostServices(studioID, page)
}

func (c postController) GetPostHomepage(userID uint64, page int) (*[]models.Post, error) {
	return App.Service.GetHomepagePostServices(userID, page)
}

func (c postController) GetOnePost(postID uint64) (*models.Post, error) {
	return App.Service.GetSinglePostService(postID)
}
func (c postController) DeletePost(postID uint64) error {
	return App.Service.DeleteSinglePostService(postID)
}

func (c postController) DeletePostComment(postCommentID uint64) error {
	return App.Service.DeleteSinglePostCommentService(postCommentID)
}

func (c postController) CreatePostComment(postID uint64, data CreatePostCommentValidation, currentUser *models.User) (*models.PostComment, error) {
	return App.Service.CreatePostComment(postID, data, currentUser)
}

func (c postController) GetAllPostComments(postID uint64) (*[]models.PostComment, error) {
	return App.Service.GetPostCommentsServices(postID)
}

func (c postController) UpdatePostComment(commentID uint64, data UpdatePostCommentValidation, currentUser *models.User) error {
	return App.Service.UpdatePostComment(commentID, map[string]interface{}{
		"is_edited":     data.IsEdited,
		"comment":       data.Comment,
		"updated_by_id": currentUser.ID,
	})
}

func (c postController) CreatePostReaction(postID uint64, data NewPostReaction, currentUser *models.User) error {
	return App.Service.CreatePostReaction(postID, data, currentUser)
}

func (c postController) RemovePostReaction(postID uint64, data RemovePostReaction, currentUser *models.User) error {
	return App.Service.RemovePostReaction(postID, data, currentUser)
}

func (c postController) CreatePostCommentReaction(postID uint64, postCommentID uint64, data NewPostCommentReaction, currentUser *models.User) error {
	return App.Service.CreatePostCommentReaction(postID, postCommentID, data, currentUser)
}

func (c postController) RemovePostCommentReaction(postID uint64, postCommentID uint64, data RemovePostCommentReaction, currentUser *models.User) error {
	return App.Service.RemovePostCommentReaction(postID, postCommentID, data, currentUser)
}
