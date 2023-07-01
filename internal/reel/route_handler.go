package reel

import (
	"fmt"
	"net/http"
	"strconv"

	"gitlab.com/phonepost/bip-be-platform/internal/canvasbranch"
	"gitlab.com/phonepost/bip-be-platform/internal/shared"

	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/reactions"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
)

// Get All Reels for your Studio
// @Summary 	Get All Reels for your Studio - Also supports ?canvasBranchID=45 QP
// @description You can call this api to Get ALL the reels for your studio
// @description :  {{ .BasePath }}/v1/reels/ - All reels for that studio
// @description :  {{ .BasePath }}/v1/reels/?canvasBranchID=256 - All reels for canvasBranchID=256
// @description :  {{ .BasePath }}/v1/reels/?blockUUID=e92942a8-de88-43e6-ba20-fb2c634129b0 - All reels for blockUUID=<>
// @Tags		Reels
// @Accept       json
// @Produce      json
// @Router 		/v1/reels/ [get]
func (r *reelRoutes) Get(c *gin.Context) {
	studioID, _ := r.GetStudioId(c)
	branchID := c.Query("canvasBranchID")
	blockUUIDReq := c.Query("blockUUID")
	canvasBranchID, _ := strconv.ParseUint(branchID, 10, 64)
	uuidAs, _ := uuid.Parse(blockUUIDReq)
	authUser, _ := r.GetLoggedInUser(c)
	instances, err := App.Service.GetAll(studioID, canvasBranchID, uuidAs, authUser)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	reelReactions := []models.ReelReaction{}
	if authUser != nil {
		reelReactions, _ = reactions.App.Repo.GetReelReaction(map[string]interface{}{"canvas_branch_id": canvasBranchID, "created_by_id": authUser.ID})
	}
	response.RenderResponse(c, SerializeDefaultManyReelsWithReactionsForUser(instances, reelReactions, authUser, nil, nil))
}

// Get Single Reel
// @Summary 	Get Single Reel
// @description You can call this api to Get One Reel By ID
// @Tags		Reels
// @Accept       json
// @Produce      json
// @Router 		/v1/reels/:reelID [get]
func (r *reelRoutes) GetOneReel(c *gin.Context) {
	studioID, _ := r.GetStudioId(c)
	reelIDStr := c.Param("reelID")
	reelID, _ := strconv.ParseUint(reelIDStr, 10, 64)
	authUser, _ := r.GetLoggedInUser(c)

	instance, err := App.Service.GetOne(studioID, reelID)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	reelReactions := []models.ReelReaction{}
	if authUser != nil {
		reelReactions, _ = reactions.App.Repo.GetReelReaction(map[string]interface{}{"canvas_branch_id": instance.CanvasBranchID, "created_by_id": authUser.ID})
	}
	response.RenderResponse(c, SerializeDefaultSingleReelsWithReactionsForUser(instance, reelReactions, authUser, nil, nil))
}

// Get All Reels for USER
// @Summary 	Get Popular Reels (POV User)
// @description It sends all REELS for now
// @Tags		Reels
// @Accept       json
// @Produce      json
// @Router 		/v1/reels/popular [get]
func (r *reelRoutes) GetPopularReels(c *gin.Context) {
	user, _ := r.GetLoggedInUser(c)
	skip, _ := strconv.Atoi(c.Query("skip"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit == 0 {
		limit = configs.PAGINATION_LIMIT
	}
	var instances *[]ReelsSerialData
	var err error
	if user != nil {
		instances, err = App.Service.GetLoggedInPopular(user, skip, limit)
	} else {
		instances, err = App.Service.GetAnonymousPopular(skip, limit)
	}
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	next := "-1"
	if len(*instances) == configs.PAGINATION_LIMIT {
		next = strconv.Itoa(skip + len(*instances))
	}
	response.RenderPaginatedResponse(c, instances, next)
}

// @Summary 	Create a reel
// @description Create reeel for your studio
// @Tags		Reels
// @Accept       json
// @Produce      json
// @Param 		body body 		NewReelCreatePOST true "Create Reel"
// @Router 		/v1/reels/ [post]
func (r *reelRoutes) Create(c *gin.Context) {
	var body NewReelCreatePOST
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	_, status := middlewares.StudioHeaderRequiredCheck(c)
	if !status {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": "Studio Header Missing",
		})
		return
	}
	user, _ := r.GetLoggedInUser(c)
	studioID, _ := r.GetStudioId(c)
	// permission check
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, body.CanvasBranchID, permissiongroup.CANVAS_BRANCH_CREATE_REEL); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	instance, errCreating := App.Controller.CreateReel(body, studioID, user.ID)
	if errCreating != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), errCreating))
		return
	}
	instance.CreatedByUser = user
	canvasbranch.App.Repo.UpdateBranchLastEdited(body.CanvasBranchID)
	go func() {
		App.Caching.InvalidateReelsCachingViaStudio(studioID)
		canvasbranch.App.Service.InvalidateBranchBlocks(body.CanvasBranchID)
	}()
	response.RenderResponse(c, SerializeDefaultReel(instance, authUserID))
}

// @Summary 	Add comment on a Reel
// @description create a new comment on a Reel
// @Tags		Reels
// @Accept       json
// @Produce      json
// @Param 		body body 		ReelCommentCreatePOST true "Create Comment on Reel"
// @Router 		/v1/reels/:reelID/comments/ [post]
func (r *reelRoutes) NewCommentReel(c *gin.Context) {
	var body ReelCommentCreatePOST
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	_, status := middlewares.StudioHeaderRequiredCheck(c)
	if !status {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": "Studio Header Missing",
		})
		return
	}
	user, _ := r.GetLoggedInUser(c)
	studioID, _ := r.GetStudioId(c)

	reelIDStr := c.Param("reelID")
	reelID, _ := strconv.ParseUint(reelIDStr, 10, 64)

	// permission check
	reelInstance, _ := App.Repo.GetReel(map[string]interface{}{"id": reelID})
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, reelInstance.CanvasBranchID, permissiongroup.CANVAS_BRANCH_COMMENT_ON_REEL); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	instance, errCreating := App.Controller.CreateReelComment(body, studioID, user.ID, reelID)
	instance.CreatedByUser = user
	if errCreating != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), errCreating))
		return
	}
	go App.Caching.InvalidateReelsCachingViaStudio(studioID)

	response.RenderResponse(c, SerializeDefaultReelComment(instance))
}

// @Summary 	Get All Comments on a Reels for your Studio
// @Tags		Reels
// @Accept       json
// @Param 		parentId 	query 	string	false "parent comment id"
// @Produce      json
// @Router 		/v1/reels/:reelID/comments/ [get]
func (r *reelRoutes) GetCommentsReel(c *gin.Context) {
	studioID, _ := r.GetStudioId(c)
	reelIDStr := c.Param("reelID")
	reelID, _ := strconv.ParseUint(reelIDStr, 10, 64)
	authUser, _ := r.GetLoggedInUser(c)
	parentIDStr := c.Query("parentId")
	parentID := utils.Uint64(parentIDStr)
	// permission check
	reelInstance, _ := App.Repo.GetReel(map[string]interface{}{"id": reelID})
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, reelInstance.CanvasBranchID, permissiongroup.CANVAS_BRANCH_VIEW); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}
	var instances *[]models.ReelComment
	var err error
	if parentID == 0 {
		instances, err = App.Service.GetAllComments(studioID, reelID)
	} else {
		instances, err = App.Service.GetChildComments(parentID)
	}
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	reelReactions := []models.ReelCommentReaction{}
	if authUser != nil {
		reelReactions, _ = reactions.App.Repo.GetReelCommentReaction(map[string]interface{}{"reel_id": reelID, "created_by_id": authUser.ID})
	}
	fmt.Println(reelReactions)
	response.RenderResponse(c, SerializeDefaultManyReelCommentsWithReactionsForUser(instances, reelReactions, authUser))

}

// GetReelsFeed Reels feed
// @Summary 	Get Reels feed from getStream. This is a paginated API
// @description This API is specific to the user. This will return all the reels which user have access in the system.
// @description To get studio specific reels => ?filter=studio and studio id as bip-studio-id in header
// @Tags		Reels
// @Accept       json
// @Produce      json
// @Router 		/v1/reels/feed [get]
func (r *reelRoutes) GetReelsFeed(c *gin.Context) {
	authUser, _ := r.GetLoggedInUser(c)
	skip, _ := strconv.Atoi(c.Query("skip"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	filter := c.Query("filter")
	studioID, _ := r.GetStudioId(c)
	var reelsFeed *[]ReelsSerialData
	var err error
	next := "-1"

	if authUser == nil && studioID != 0 {
		reelsFeed, err = App.Service.GetAnonymousStudioPopular(skip, limit, studioID)
		response.RenderResponse(c, reelsFeed)
		return
	} else if filter != "studio" {
		// Check if cache Exists?
		// For now, we are disabling caching on User+Studio
		/*
			checkCachedUserData := CheckIfNonStudioDataRedisWithUserIDKeyExists(authUser.ID)
			if checkCachedUserData && skip == 0 {
				finalResponse := GetReelDataViaRedisWithUserID(authUser.ID)
				//reelsFeed = GetReelDataViaRedisWithUserID(authUser.ID)
				c.JSON(http.StatusOK, finalResponse)
				return
			} else {
				reelsFeed, err = App.Controller.GetReelsFeed(authUser, skip, limit)
				if len(*reelsFeed) == 15 {
					next = strconv.Itoa(skip + len(*reelsFeed))
				}
				if skip == 0 {
					SetReelsNonStudioDataRedisWithUserID(authUser.ID, shared.NewGenericResponseV1{Data: reelsFeed, Next: next})
				}
			}
		*/
		// Getting direct data
		reelsFeed, err = App.Controller.GetReelsFeed(authUser, skip, limit)

	} else {
		// Cached based on StudioID + UserID
		checkCachedUserStudioData := App.Caching.CheckIfReelsWithStudioDataKeyExists(authUser.ID, studioID)
		if checkCachedUserStudioData && skip == 0 {
			finalResponse := App.Caching.GetReelDataWithStudioAndStudioViaRedisWithUserID(authUser.ID, studioID)
			//reelsFeed = GetReelDataViaRedisWithUserID(authUser.ID)
			c.JSON(http.StatusOK, finalResponse)
			return
		} else {
			reelsFeed, err = App.Controller.GetReelsFeedForStudio(authUser, studioID, skip, limit)
			if len(*reelsFeed) == 15 {
				next = strconv.Itoa(skip + len(*reelsFeed))
			}
			if skip == 0 {
				App.Caching.SetReelsWithStudioDataRedis(authUser.ID, studioID, shared.NewGenericResponseV1{Data: reelsFeed, Next: next})
			}
		}
	}
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}

	if len(*reelsFeed) == 15 {
		next = strconv.Itoa(skip + len(*reelsFeed))
	}

	c.JSON(http.StatusOK, shared.NewGenericResponseV1{Data: reelsFeed, Next: next})
	return
	//response.RenderPaginatedResponse(c, reelsFeed, next)
}

// @Summary 	Delete a reel
// @Tags		Reels
// @Accept       json
// @Produce      json
// @Router 		/v1/reels/:reelID [delete]
func (r *reelRoutes) DeleteReel(c *gin.Context) {
	userID := r.GetLoggedInUserId(c)
	//studioID, _ := r.GetStudioId(c)
	reelIDStr := c.Param("reelID")
	reelID, _ := strconv.ParseUint(reelIDStr, 10, 64)

	// permission check
	reelInstance, _ := App.Repo.GetReel(map[string]interface{}{"id": reelID})
	studioID := reelInstance.StudioID
	authUserID := r.GetLoggedInUserId(c)
	if !(reelInstance.CreatedByID == authUserID) {
		if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, reelInstance.CanvasBranchID, permissiongroup.CANVAS_BRANCH_MANAGE_CONTENT); err != nil || !hasPermission || reelInstance.CreatedByID != authUserID {
			response.RenderPermissionError(c)
			return
		}
	}

	errDeleting := App.Controller.DeleteReel(reelID, userID)
	if errDeleting != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), errDeleting))
		return
	}
	canvasbranch.App.Repo.UpdateBranchLastEdited(reelInstance.CanvasBranchID)
	// invalidating the blocks cache
	go func() {
		App.Caching.InvalidateReelsCachingViaStudio(studioID)
		canvasbranch.App.Service.InvalidateBranchBlocks(reelInstance.CanvasBranchID)
	}()
	c.JSON(http.StatusOK, gin.H{
		"message": "Reel is Deleted",
	})
}

// DeleteReelComment
// @Summary 	Delete a reel
// @Tags		Reels
// @Param 		reelID 	path 	string	true "Reel Id"
// @Param 		reelCommentID 	path 	string	true "reelCommentID"
// @Param 		resolved 	query 	string	true "resolved"
// @Accept       json
// @Produce      json
// @Router 		/v1/reels/{reelID}/comments/{reelCommentID} [delete]
func (r *reelRoutes) DeleteReelComment(c *gin.Context) {
	userID := r.GetLoggedInUserId(c)
	//studioID, _ := r.GetStudioId(c)
	reelCommentIDStr := c.Param("reelCommentID")
	reelCommentID := utils.Uint64(reelCommentIDStr)

	reelComment, _ := App.Repo.GetReelComment(map[string]interface{}{"id": reelCommentID})
	if userID != reelComment.CreatedByID {
		response.RenderPermissionError(c)
		return
	}

	errDeleting := App.Controller.DeleteReelComment(reelComment)
	if errDeleting != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), errDeleting))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reel Comment is Deleted",
	})
}
