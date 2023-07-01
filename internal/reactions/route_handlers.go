package reactions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/apiutil"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
)

// Create New Mention
// @Summary 	Create Mention
// @Description This API will let you create a reaction on following
//@Description models Blocks, Block Thread, Block Comment, Reel, Reel Comment
// @Description Requires a scope: scope can be "block","block_comment","reel","reel_comment", "block_thread"
// @Tags		Reactions
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		body 		body 		CreateMentionPost true "Create Mention"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/reactions/create [post]
func (r reactionRoutes) Create(c *gin.Context) {
	var body CreateMentionPost
	if err := apiutil.Bind(c, &body); err != nil {
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}
	// Validation
	validateError := body.Validate()
	if validateError != nil {
		response.RenderCustomErrorResponse(c, validateError)
		return
	}
	userID := r.GetLoggedInUserId(c)
	studioID, _ := r.GetStudioId(c)

	err := App.Controller.Create(body, studioID, userID)

	if err != nil && err.Error() == response.NoPermissionError {
		response.RenderPermissionError(c)
		return
	}

	App.Repo.UpdateBranchLastEdited(body.CanvasBranchID)
	go func() {
		App.Service.InvalidateBranchBlocks(body.CanvasBranchID)
		App.Service.InvalidateReelsCachingViaStudio(studioID)

	}()
	c.JSON(http.StatusOK, gin.H{
		"message": "Reaction added to " + body.Scope,
	})
	return
}

// Remove Mention
// @Summary 	Removes a Reaction
// @Description This API will let you remove a reaction on following
//@Description models Blocks, Block Thread, Block Comment, Reel, Reel Comment
// @Description Requires a scope: scope can be "block","block_comment","reel","reel_comment", "block_thread"
// @Tags		Reactions
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		body 		body 		CreateMentionPost true "Remove Reaction"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/reactions/remove [post]
func (r reactionRoutes) Remove(c *gin.Context) {
	var body CreateMentionPost
	if err := apiutil.Bind(c, &body); err != nil {
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}
	// Validation
	validateError := body.Validate()
	if validateError != nil {
		response.RenderCustomErrorResponse(c, validateError)
		return
	}
	userID := r.GetLoggedInUserId(c)
	studioID, _ := r.GetStudioId(c)

	err := App.Controller.Remove(body, studioID, userID)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	App.Repo.UpdateBranchLastEdited(body.CanvasBranchID)
	go func() {
		App.Service.InvalidateBranchBlocks(body.CanvasBranchID)
	}()
	c.JSON(http.StatusOK, gin.H{
		"message": "Reaction removed from  " + body.Scope,
	})
	return
}
