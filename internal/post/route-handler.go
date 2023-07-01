package post

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/shared"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/apiClient"
	"net/http"
	"strconv"
)

// Create New Post
// @Summary 	Creates a new POST
// @Tags		Post
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		body 		body 		NewPostThread true "Create Studio Post"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/posts/create [post]
func (r *postRoutes) CreatePost(c *gin.Context) {
	var body NewPostThread
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Studio ID
	rawStudioId, _ := c.Get("currentStudio")
	studioID := rawStudioId.(uint64)
	// User Auth
	authUser, userNotFound := r.RouteHelper.GetLoggedInUser(c)
	if userNotFound != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": userNotFound.Error(),
		})
		return
	}

	// Check is User is Studio Member of not
	isStudioMember := shared.IsUserStudioMember(authUser.ID, studioID)
	if !isStudioMember {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User is not studio member.",
		})
		return
	}

	// Create Post Instance
	postInstance, errCreatingPostInstance := App.Controller.CreatePost(body, studioID, *authUser)
	postInstance1, _ := App.Controller.GetOnePost(postInstance.ID)
	fmt.Println(postInstance)
	if errCreatingPostInstance != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errCreatingPostInstance.Error(),
		})
		return
	}

	go func() {
		App.Service.InvalidateStudioPostCache(studioID)
		postStr, _ := json.Marshal(postInstance1)
		apiClient.AddToQueue(apiClient.SendPostToIntegration, postStr, apiClient.DEFAULT, apiClient.CommonRetry)
	}()

	response.RenderResponse(c, SinglePostSerializerData(postInstance1, authUser.ID))
	return
}

// @Summary Get All Posts for this Studio
// @Tags		Post
// @Router 		/v1/posts/  [get]
func (r *postRoutes) GetAllPost(c *gin.Context) {

	page := 1
	pageStr := c.Query("page")
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}

	// User Auth
	authUser, userNotFound := r.RouteHelper.GetLoggedInUser(c)
	if userNotFound != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": userNotFound.Error(),
		})
		return
	}

	// get studio
	studioID, _ := r.GetStudioId(c)
	if studioID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Studio header not found.",
		})
		return
	}
	// Check Cache before DB call
	getFromRedis := CheckIfPostStudioKeyExists(studioID)
	if page == 1 && getFromRedis {
		again := GetPostDataViaRedis(studioID)
		c.JSON(http.StatusOK, again)
		return
	}

	postInstances, errorGettingPosts := App.Controller.GetAllPosts(studioID, page)
	if errorGettingPosts != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errorGettingPosts.Error(),
		})
		return
	}
	posts := ManyPostSerializerData(postInstances, authUser.ID)
	var postInterface interface{}
	postInterface = posts
	finalResponse := shared.GetPaginationData(page, map[string]interface{}{"studio_id": studioID}, &models.Post{}, postInterface)
	// This one populates the data on redis is page is 1 and the data was not picked from cache.
	if page == 1 && !getFromRedis {
		SetPostDataRedis(studioID, finalResponse)
	}
	c.JSON(http.StatusOK, finalResponse)
	return
}

// @Summary Get All Posts for this Studio
// @Tags		Post
// @Router 		/v1/posts/homepage  [get]
func (r *postRoutes) GetPostHomepage(c *gin.Context) {

	page := 1
	pageStr := c.Query("page")
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}

	// get studio
	// User Auth
	authUser, userNotFound := r.RouteHelper.GetLoggedInUser(c)
	if userNotFound != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": userNotFound.Error(),
		})
		return
	}

	postInstances, errorGettingPosts := App.Controller.GetPostHomepage(authUser.ID, page)
	if errorGettingPosts != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errorGettingPosts.Error(),
		})
		return
	}

	posts := ManyPostSerializerData(postInstances, authUser.ID)
	studioIds := App.Service.GetStudioIDArrayUserID(authUser.ID)
	var postInterface interface{}
	postInterface = posts

	c.JSON(http.StatusOK, shared.GetPaginationData(page, map[string]interface{}{"studio_id": studioIds}, &models.Post{}, postInterface))

	return
}

// @Summary Get Single Post
// @Tags		Post
// @Router 		/v1/posts/:postID  [get]
func (r *postRoutes) GetOnePost(c *gin.Context) {

	authUser, userNotFound := r.RouteHelper.GetLoggedInUser(c)
	if userNotFound != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": userNotFound.Error(),
		})
		return
	}
	postIDStr := c.Param("postID")
	postID, _ := strconv.ParseUint(postIDStr, 10, 64)
	postInstance, errorGettingPost := App.Controller.GetOnePost(postID)
	if errorGettingPost != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errorGettingPost.Error(),
		})
		return
	}
	response.RenderResponse(c, OnePostSerializerData(postInstance, authUser.ID))
	return
}

// @Summary Patch a single post
// @Tags		Post
// @Param 		body 		body 		UpdatePostThread true "Update Studio Post"
// @Router 		/v1/posts/:postID/edit  [patch]
func (r *postRoutes) UpdatePost(c *gin.Context) {
	var body UpdatePostThread
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// User Auth
	authUser, userNotFound := r.RouteHelper.GetLoggedInUser(c)
	if userNotFound != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": userNotFound.Error(),
		})
		return
	}

	postIDStr := c.Param("postID")
	postID, _ := strconv.ParseUint(postIDStr, 10, 64)

	// Update the post
	errorUpdatingPost := App.Controller.UpdatePost(postID, body, *authUser)
	if errorUpdatingPost != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errorUpdatingPost.Error(),
		})
		return
	}

	// get the post instance and send to serializer
	postInstance, errorGettingPost := App.Controller.GetOnePost(postID)
	if errorGettingPost != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errorGettingPost.Error(),
		})
		return
	}
	App.Service.InvalidateStudioPostCache(postInstance.StudioID)
	response.RenderResponse(c, OnePostSerializerData(postInstance, authUser.ID))
	return

}

// @Summary Delete Single Post
// @Tags		Post
// @Router 		/v1/posts/:postID  [delete]
func (r *postRoutes) DeletePost(c *gin.Context) {
	postIDStr := c.Param("postID")
	postID, _ := strconv.ParseUint(postIDStr, 10, 64)

	rawStudioId, _ := c.Get("currentStudio")
	studioID := rawStudioId.(uint64)

	error1 := App.Controller.DeletePost(postID)
	if error1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": error1.Error(),
		})
		return
	}

	App.Service.InvalidateStudioPostCache(studioID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Deleted",
	})

	return
}

// @Summary Get All Comments on a POST
// @Tags		Post
// @Router 		/v1/posts/:postID/comments  [get]
func (r *postRoutes) GetAllPostComments(c *gin.Context) {
	postIDStr := c.Param("postID")
	postID, _ := strconv.ParseUint(postIDStr, 10, 64)
	currentUser, userNotFound := r.RouteHelper.GetLoggedInUser(c)
	if userNotFound != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": userNotFound.Error(),
		})
		return
	}
	comments, err := App.Controller.GetAllPostComments(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	response.RenderResponse(c, ManyPostCommentsSerializerData(comments, currentUser.ID))
	return
}

// Create New Post Comment
// @Summary 	Creates a new POST
// @Tags		Post
// @Param 		body 		body 		CreatePostCommentValidation true "Create Studio Post Comment"
// @Router 		/v1/posts/:postID/comments/create  [post]
func (r *postRoutes) CreateCommentPost(c *gin.Context) {
	// get POST data
	postIDStr := c.Param("postID")
	postID, _ := strconv.ParseUint(postIDStr, 10, 64)
	var body CreatePostCommentValidation
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	currentUser, userNotFound := r.RouteHelper.GetLoggedInUser(c)
	if userNotFound != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": userNotFound.Error(),
		})
		return
	}
	postcomment, err := App.Controller.CreatePostComment(postID, body, currentUser)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	response.RenderResponse(c, SinglePostCommentSerializerData(postcomment, currentUser.ID))
	return
}

// @Summary 	Delete a comment on a POST
// @Tags		Post
// @Router 		/v1/posts/:postID/comments/:postCommentID  [delete]
func (r *postRoutes) DeleteCommentPost(c *gin.Context) {
	commentPostIDStr := c.Param("postCommentID")
	commentPostID, _ := strconv.ParseUint(commentPostIDStr, 10, 64)

	error1 := App.Controller.DeletePostComment(commentPostID)
	if error1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": error1.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Deleted",
	})
	return

}

// @Summary Patch a single comment
// @Tags		Post
// @Param 		body 		body 		UpdatePostCommentValidation true "Update Post Comment"
// @Router 		/v1/posts/:postID/comments/:postCommentID/edit  [patch]
func (r *postRoutes) UpdateCommentPost(c *gin.Context) {
	var body UpdatePostCommentValidation
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// User Auth
	authUser, userNotFound := r.RouteHelper.GetLoggedInUser(c)
	if userNotFound != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": userNotFound.Error(),
		})
		return
	}

	commentPostIDStr := c.Param("postCommentID")
	commentPostID, _ := strconv.ParseUint(commentPostIDStr, 10, 64)

	// Update the post
	errorUpdatingComment := App.Controller.UpdatePostComment(commentPostID, body, authUser)
	if errorUpdatingComment != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errorUpdatingComment.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Update",
	})
}

// @Summary New Post Reaction
// @Tags		Post
// @Param 		body 		body 		NewPostReaction true "New post reaction"
// @Router 		/v1/posts/:postID/add-reaction  [post]
func (r *postRoutes) AddReactionPost(c *gin.Context) {
	var body NewPostReaction
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// User Auth
	authUser, userNotFound := r.RouteHelper.GetLoggedInUser(c)
	if userNotFound != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": userNotFound.Error(),
		})
		return
	}

	postIDStr := c.Param("postID")
	postID, _ := strconv.ParseUint(postIDStr, 10, 64)
	errCreatingReaction := App.Controller.CreatePostReaction(postID, body, authUser)

	if errCreatingReaction != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errCreatingReaction.Error(),
		})
		return
	}

	rawStudioId, _ := c.Get("currentStudio")
	studioID := rawStudioId.(uint64)
	App.Service.InvalidateStudioPostCache(studioID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Reaction Added",
	})
	return
}

// @Summary New Post Reaction
// @Tags		Post
// @Param 		body 		body 		RemovePostReaction true "New post reaction"
// @Router 		/v1/posts/:postID/remove-reaction  [post]
func (r *postRoutes) RemoveReactionPost(c *gin.Context) {
	var body RemovePostReaction
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// User Auth
	authUser, userNotFound := r.RouteHelper.GetLoggedInUser(c)
	if userNotFound != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": userNotFound.Error(),
		})
		return
	}

	postIDStr := c.Param("postID")
	postID, _ := strconv.ParseUint(postIDStr, 10, 64)
	errCreatingReaction := App.Controller.RemovePostReaction(postID, body, authUser)

	if errCreatingReaction != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errCreatingReaction.Error(),
		})
		return
	}

	rawStudioId, _ := c.Get("currentStudio")
	studioID := rawStudioId.(uint64)
	App.Service.InvalidateStudioPostCache(studioID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Reaction Deleted",
	})
	return
}

// @Summary New Post Comment Reaction
// @Tags		Post
// @Param 		body 		body 		NewPostCommentReaction true "New post comment reaction"
// @Router 		/v1/posts/:postID/comments/:postCommentID/add-reaction  [post]
func (r *postRoutes) AddReactionPostComment(c *gin.Context) {
	var body NewPostCommentReaction
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// User Auth
	authUser, userNotFound := r.RouteHelper.GetLoggedInUser(c)
	if userNotFound != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": userNotFound.Error(),
		})
		return
	}

	postIDStr := c.Param("postID")
	postID, _ := strconv.ParseUint(postIDStr, 10, 64)

	postCommentStr := c.Param("postCommentID")
	postCommentID, _ := strconv.ParseUint(postCommentStr, 10, 64)

	errCreatingReaction := App.Controller.CreatePostCommentReaction(postID, postCommentID, body, authUser)

	if errCreatingReaction != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errCreatingReaction.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reaction Added To Comment",
	})
	return
}

// @Summary New Post Comment Reaction
// @Tags		Post
// @Param 		body 		body 		RemovePostCommentReaction true "Remove post comment reaction"
// @Router 		/v1/posts/:postID/comments/:postCommentID/remove-reaction  [post]
func (r *postRoutes) RemoveReactionPostComment(c *gin.Context) {
	var body RemovePostCommentReaction
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// User Auth
	authUser, userNotFound := r.RouteHelper.GetLoggedInUser(c)
	if userNotFound != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": userNotFound.Error(),
		})
		return
	}

	postIDStr := c.Param("postID")
	postID, _ := strconv.ParseUint(postIDStr, 10, 64)

	postCommentStr := c.Param("postCommentID")
	postCommentID, _ := strconv.ParseUint(postCommentStr, 10, 64)

	errCreatingReaction := App.Controller.RemovePostCommentReaction(postID, postCommentID, body, authUser)

	if errCreatingReaction != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errCreatingReaction.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reaction Added To Comment",
	})
	return
}
