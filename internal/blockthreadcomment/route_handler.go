package blockThreadCommentcomment

import (
	"net/http"
	"strconv"

	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"

	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
)

// Get Block Thread Comment
// @Summary 	Get Block Thread Comment
// @Description
// @Tags		BlockThreadComment
// @Security 	bearerAuth
// @Param 		blockThreadID 	path 	string	true "Block Thread Comment ID"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/block-thread-comment/{blockThreadID} [get]
func (r *blockThreadCommentRoutes) Get(c *gin.Context) {
	blockThreadID, _ := c.Params.Get("blockThreadID")
	parsedBlockThreadID, err := strconv.ParseUint(blockThreadID, 10, 64)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	authUser, _ := r.GetLoggedInUser(c)
	authUserID := r.GetLoggedInUserId(c)

	// checking permissions
	blockThread, err := App.Repo.GetBlockThread(map[string]interface{}{"id": parsedBlockThreadID})
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, blockThread.CanvasBranchID, permissiongroup.CANVAS_BRANCH_VIEW); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	blockThreadComment, err := App.Controller.Get(parsedBlockThreadID, authUser)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	response.RenderResponse(c, blockThreadComment)
}

// GetReply Block Thread Comment Replies
// @Summary 	Get Block Thread Comment Replies
// @Description
// @Tags		BlockThreadComment
// @Security 	bearerAuth
// @Param 		blockThreadID 		path 	string	true "Block Thread Comment ID"
// @Param 		parentCommentID 	path 	string	true "Parent Comment ID"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/block-thread-comment/reply/{blockThreadID}/{parentCommentID} [get]
func (r *blockThreadCommentRoutes) GetReply(c *gin.Context) {
	blockThreadID, _ := c.Params.Get("blockThreadID")
	parsedBlockThreadID, err := strconv.ParseUint(blockThreadID, 10, 64)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	parentCommentID, _ := c.Params.Get("parentCommentID")
	parsedParentCommentID, err := strconv.ParseUint(parentCommentID, 10, 64)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	blockThreadComment, err := App.Controller.GetReply(parsedBlockThreadID, parsedParentCommentID)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	response.RenderResponse(c, blockThreadComment)
}

// Create New Block Thread Comment
// @Summary 	Create Block Thread Comment
// @Description
// @Tags		BlockThreadComment
// @Security 	bearerAuth
// @Param 		body  		body 		PostBlockThreadComment true "Create BlockThreadComment"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/block-thread-comment/ [post]
func (r *blockThreadCommentRoutes) Create(c *gin.Context) {
	var body PostBlockThreadComment
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	authUser, err := r.GetLoggedInUser(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	authUserID := r.GetLoggedInUserId(c)
	studioID, _ := r.GetStudioId(c)

	// checking permissions
	blockThread, err := App.Repo.GetBlockThread(map[string]interface{}{"id": body.ThreadID})
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, blockThread.CanvasBranchID, permissiongroup.CANVAS_BRANCH_ADD_COMMENT); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	blockThreadComment, err := App.Controller.Create(body, authUser, studioID)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	response.RenderResponse(c, blockThreadComment)
}

// Update Block Thread Comment
// @Summary 	Update Block Thread Comment
// @Description
// @Tags		BlockThreadComment
// @Security 	bearerAuth
// @Param 		body  		body 		PatchBlockThreadComment true "Create BlockThreadComment"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/block-thread-comment/ [patch]
func (r *blockThreadCommentRoutes) Update(c *gin.Context) {
	var body PatchBlockThreadComment
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	authUser, err := r.GetLoggedInUser(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	authUserID := r.GetLoggedInUserId(c)

	// checking permissions
	blockThreadComment, err := App.Repo.GetPreloadThread(map[string]interface{}{"id": body.ID})
	if blockThreadComment.CreatedByID != authUserID {
		response.RenderPermissionError(c)
		return
	}

	err = App.Controller.Update(body, authUser)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	response.RenderSuccessResponse(c, "BlockThreadComment Updated successfully")
}

// Delete Block Thread Comment
// @Summary 	Delete Block Thread Comment
// @Description
// @Tags		BlockThreadComment
// @Security 	bearerAuth
// @Param 		blockThreadCommentID 	path 	string	true "Block Thread Comment ID"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/block-thread-comment/{blockThreadCommentID} [delete]
func (r *blockThreadCommentRoutes) Delete(c *gin.Context) {

	blockThreadCommentID, _ := c.Params.Get("blockThreadCommentID")
	parsedBlockThreadCommentID, err := strconv.ParseUint(blockThreadCommentID, 10, 64)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	authUser, err := r.GetLoggedInUser(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	authUserID := r.GetLoggedInUserId(c)

	// Checking permissions
	blockThreadComment, err := App.Repo.GetPreloadThread(map[string]interface{}{"id": blockThreadCommentID})
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	if blockThreadComment.CreatedByID != authUser.ID {
		if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, blockThreadComment.Thread.CanvasBranchID, permissiongroup.CANVAS_BRANCH_MANAGE_CONTENT); err != nil || !hasPermission {
			response.RenderPermissionError(c)
			return
		}
	}

	err = App.Controller.Delete(parsedBlockThreadCommentID, authUser)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	response.RenderSuccessResponse(c, "BlockThreadComment Deleted successfully")
}
