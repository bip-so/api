package blockthread

import (
	"fmt"
	"net/http"
	"strconv"

	"gitlab.com/phonepost/bip-be-platform/internal/canvasbranch"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"

	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
)

// Get Block Thread
// @Summary 	Get Block Thread
// @Description
// @Tags		BlockThread
// @Security 	bearerAuth
// @Param 		blockThreadID 	path 	string	true "Block Thread Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/block-thread/{blockThreadID} [get]
func (r *blockThreadRoutes) Get(c *gin.Context) {
	blockThreadID, _ := c.Params.Get("blockThreadID")
	parsedBlockThreadID, err := strconv.ParseUint(blockThreadID, 10, 64)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	authUser, _ := r.GetLoggedInUser(c)

	blockThread, err := App.Controller.Get(parsedBlockThreadID, authUser)
	if err != nil && err.Error() == response.NoPermissionError {
		response.RenderPermissionError(c)
		return
	}
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	response.RenderResponse(c, blockThread)
}

// GetByBranch Block Threads By BranchID
// @Summary 	Get Block Threads By BranchID
// @Tags		BlockThread
// @Security 	bearerAuth
// @Param 		canvasBranchID 	path 	string	true "Canvas Branch Id"
// @Param 		resolved 	query 	string	true "resolved"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/block-thread/branch/{canvasBranchID} [get]
func (r *blockThreadRoutes) GetByBranch(c *gin.Context) {
	canvasBranchID, _ := c.Params.Get("canvasBranchID")
	inviteCode := c.Query("inviteCode")
	showResolved := c.DefaultQuery("resolved", "false")
	parsedcanvasBranchID, err := strconv.ParseUint(canvasBranchID, 10, 64)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	authUser, _ := r.GetLoggedInUser(c)
	authUserID := r.GetLoggedInUserId(c)
	if authUser == nil && inviteCode != "" {
		branch, _ := permissions.App.Repo.GetBranchWithRepo(map[string]interface{}{"id": parsedcanvasBranchID})
		_, err := permissions.App.Service.CheckBranchAccessToken(inviteCode, branch.CanvasRepository.Key)
		if err != nil {
			response.RenderPermissionError(c)
			return
		}
		blockThreads, err := App.Controller.GetAllByBranch(parsedcanvasBranchID, authUser, showResolved)
		if err != nil {
			response.RenderCustomErrorResponse(c, err)
			return
		}
		response.RenderResponse(c, blockThreads)
		return
	}
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, parsedcanvasBranchID, permissiongroup.CANVAS_BRANCH_VIEW); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	blockThreads, err := App.Controller.GetAllByBranch(parsedcanvasBranchID, authUser, showResolved)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	response.RenderResponse(c, blockThreads)
}

// Create New Block Thread
// @Summary 	Create Block Thread
// @Description
// @Tags		BlockThread
// @Security 	bearerAuth
// @Param 		body  		body 		PostBlockThread true "Create BlockThread"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/block-thread/ [post]
func (r *blockThreadRoutes) Create(c *gin.Context) {
	var body PostBlockThread
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
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, body.CanvasBranchID, permissiongroup.CANVAS_BRANCH_ADD_COMMENT); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	blockThread, err := App.Controller.Create(body, authUser)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	canvasbranch.App.Repo.UpdateBranchLastEdited(body.CanvasBranchID)
	// invalidating the blocks cache
	go func() {
		canvasbranch.App.Service.InvalidateBranchBlocks(body.CanvasBranchID)
	}()
	response.RenderResponse(c, blockThread)
}

// Update Block Thread
// @Summary 	Update Block Thread
// @Description
// @Tags		BlockThread
// @Security 	bearerAuth
// @Param 		body  		body 		PatchBlockThread true "Create BlockThread"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/block-thread/ [patch]
func (r *blockThreadRoutes) Update(c *gin.Context) {
	var body PatchBlockThread
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
	blockThread, err := App.Repo.Get(map[string]interface{}{"id": body.ID})
	if blockThread.CreatedByID != authUserID {
		response.RenderPermissionError(c)
		return
	}

	err = App.Controller.Update(body, authUser)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	response.RenderSuccessResponse(c, "BlockThread Updated successfully")
}

// Delete Block Thread
// @Summary 	Delete Block Thread
// @Description
// @Tags		BlockThread
// @Security 	bearerAuth
// @Param 		blockThreadID 	path 	string	true "Block Thread Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/block-thread/{blockThreadID} [delete]
func (r *blockThreadRoutes) Delete(c *gin.Context) {

	blockThreadID, _ := c.Params.Get("blockThreadID")
	parsedBlockThreadID, err := strconv.ParseUint(blockThreadID, 10, 64)
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

	// Checking permissions for delete
	blockThread, err := App.Repo.Get(map[string]interface{}{"id": blockThreadID})
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	if blockThread.CreatedByID != authUser.ID {
		if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, blockThread.CanvasBranchID, permissiongroup.CANVAS_BRANCH_MANAGE_CONTENT); err != nil || !hasPermission {
			response.RenderPermissionError(c)
			return
		}
	}
	/* Workflow : Block Thread Comment Delete
	Get the Branch Instance with Repo
	Check if the comment is deleted on the Main Branch or Rough Branch
	if Main branch the get all the Rough Branches / Find the cloned instance of the Bloxk Delete them all
	If Branch is Rough Branch do Nothing.

	1. Delete all cloned instances of a given Comment when the comment is deleted on the Parent
	2. If the comment is deleted by Admin on the RB is should be deleted on the Parent Branch.
	*/
	// Flag
	isCommentOnMainBanch := false
	//get Branch Instance
	branch, _ := App.Repo.GetBranchAndRepPreload((map[string]interface{}{"id": blockThread.CanvasBranchID}))
	if branch.ID == *branch.CanvasRepository.DefaultBranchID {
		fmt.Println("yes deleted on main branch")
		isCommentOnMainBanch = true
	}
	// Comment is on Main Branch
	if isCommentOnMainBanch {
		_ = App.Controller.DeleteClonedCommentsOnRoughBranch(blockThread, branch, authUser)
	}
	// Soft delete the comment
	err = App.Controller.Delete(parsedBlockThreadID, authUser)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	canvasbranch.App.Repo.UpdateBranchLastEdited(branch.ID)
	// invalidating the blocks cache
	go func() {
		canvasbranch.App.Service.InvalidateBranchBlocks(branch.ID)
	}()
	response.RenderSuccessResponse(c, "BlockThread Deleted successfully")
}

// Resolve a Block Thread
// @Summary 	Resolve a Block Thread
// @Description
// @Tags		BlockThread
// @Security 	bearerAuth
// @Param 		body  		body 		EmptyBlockThread true "Resolve BlockThread"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/block-thread/:blockThreadID/resolve [post]
func (r *blockThreadRoutes) Resolve(c *gin.Context) {
	// @todo: Ask do wer need to check any permission
	var body EmptyBlockThread
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	authUserID := r.GetLoggedInUserId(c)

	blockThreadID, _ := c.Params.Get("blockThreadID")
	parsedBlockThreadID, err := strconv.ParseUint(blockThreadID, 10, 64)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	// Checking permissions
	blockThread, err := App.Repo.Get(map[string]interface{}{"id": blockThreadID})
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	if blockThread.CreatedByID != authUserID {
		if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, blockThread.CanvasBranchID, permissiongroup.CANVAS_BRANCH_MANAGE_CONTENT); err != nil || !hasPermission {
			response.RenderPermissionError(c)
			return
		}
	}

	errResolving := App.Controller.Resolve(parsedBlockThreadID, authUserID)
	if errResolving != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	canvasbranch.App.Repo.UpdateBranchLastEdited(blockThread.CanvasBranchID)
	// invalidating the blocks cache
	go func() {
		canvasbranch.App.Service.InvalidateBranchBlocks(blockThread.CanvasBranchID)
	}()
	c.JSON(http.StatusOK, gin.H{
		"message": "Thread is resolved",
	})
}
