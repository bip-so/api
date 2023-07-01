package canvasbranch

import (
	"encoding/json"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"net/http"
	"strconv"

	"gitlab.com/phonepost/bip-be-platform/internal/canvasrepo"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/apiClient"

	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

// Get Branch Metadata
// @Summary 	Get branch metadata with MR if present.
// @Tags		CanvasBranch
// @Router 		/v1/canvas-branch/{canvasBranchID}  [get]
func (r *canvasBranchRoutes) Get(c *gin.Context) {
	inviteCode := c.Query("inviteCode")
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)

	// userID := r.GetLoggedInUserId(c)
	userInstance, _ := r.GetLoggedInUser(c)
	authUserID := r.GetLoggedInUserId(c)

	canUserViewBranch, errGettingPermissions := permissions.App.Service.CanUserDoThisOnBranch(authUserID, canvasBranchID, permissiongroup.CANVAS_BRANCH_VIEW)
	if errGettingPermissions != nil && errGettingPermissions == gorm.ErrRecordNotFound {
		response.RenderNotFoundResponse(c, errGettingPermissions.Error())
		return
	}
	if errGettingPermissions != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errGettingPermissions.Error(),
		})
		return
	}

	branch, err := App.Repo.GetBranchWithRepo(map[string]interface{}{"id": canvasBranchID})
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	if userInstance == nil && inviteCode != "" {
		canvasBranchAccessToken, err := permissions.App.Service.CheckBranchAccessToken(inviteCode, branch.CanvasRepository.Key)
		if err != nil || canvasBranchAccessToken == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Anonymous user does not have permissions to view this.",
			})
			return
		}
	} else if branch.PublicAccess == "private" && !canUserViewBranch {
		// We are adding a requested is User already have Access Request
		exists := queries.App.AccessRequestQuery.AccessRequestExistsSimple(canvasBranchID, authUserID)

		c.JSON(http.StatusForbidden, gin.H{
			"error":            "User does not have permissions to view branch",
			"access_requested": exists,
		})
		return
	}

	attributions, _ := App.Service.GetBranchAttributions(userInstance, canvasBranchID)
	response.RenderResponse(c, BranchMetaWithMRSerializer(branch, canvasBranchID, authUserID, attributions))
}

// GetRepoBranch Canvas Branch & repo Metadata
// @Summary 	Get branch metadata with repo and MR if present.
// @Tags		CanvasBranch
// @Router 		/v1/canvas-branch/repo/{canvasBranchID}  [get]
func (r *canvasBranchRoutes) GetRepoBranch(c *gin.Context) {
	inviteCode := c.Query("inviteCode")
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)
	// userID := r.GetLoggedInUserId(c)
	userInstance, _ := r.GetLoggedInUser(c)
	branch, err := App.Repo.GetBranchWithRepo(map[string]interface{}{"id": canvasBranchID})
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	canvasRepoData, repoErr := canvasrepo.App.Service.GetCanvasRepoWithKey(branch.CanvasRepository.Key, userInstance)
	if repoErr != nil {
		response.RenderPermissionErrorDataResponse(c, repoErr)
		return
	}
	canvasBranchData, branchErr := App.Service.GetCanvasBranchData(canvasBranchID, userInstance, inviteCode)
	if userInstance != nil {
		canvasrepo.App.UserRepoHistory.AddToUserRepoToSet(userInstance.ID, canvasRepoData.ID)
	}
	resp := map[string]interface{}{
		"canvasRepo":      canvasRepoData,
		"canvasBranch":    canvasBranchData,
		"canvasBranchErr": branchErr,
	}
	response.RenderCustomResponse(c, resp)
}

// Get invited users (Invited)
// @Summary 	Get Emails which were invited on a Branch
// @Tags		CanvasBranch
// @Router 		/v1/canvas-branch/{canvasBranchID}/invited  [get]
func (r *canvasBranchRoutes) Invited(c *gin.Context) {
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)
	emails := queries.App.BranchQuery.GetPendingEmails(canvasBranchID)
	c.JSON(http.StatusOK, gin.H{
		"pending_emails": emails,
	})
	return
}

// Create
// @Summary 	Create a new canvasBranch for the canvasRepo
// @Tags		CanvasBranch
// @Security 	bearerAuth
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Router 		/v1/canvas-branch/create  [post]
// @Param 		body body 		newCanvasBranchPost true "Create Canvas Repo"
func (r *canvasBranchRoutes) Create(c *gin.Context) {
	var body newCanvasBranchPost
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctxUser, _ := c.Get("currentUser")
	loggedInUser := ctxUser.(*models.User)
	studioID, _ := r.GetStudioId(c)
	authUserID := r.GetLoggedInUserId(c)

	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, body.FromCanvasBranchID, permissiongroup.CANVAS_BRANCH_EDIT); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	// Call the controller
	data, err := App.Controller.Create(*loggedInUser, body, studioID)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	response.RenderResponse(c, SerializeCanvasBranch(data))
}

// @Summary 	Update Branch Visibility
// @Tags		CanvasBranch
// @Accept       json
// @Produce      json
// @Param 		body body 		CanvasBranchVisibilityPost true "Update"
// @Router 		/v1/canvas-branch/branch-ops/{canvasBranchID}/visibility [post]
func (r *canvasBranchRoutes) UpdateVisibility(c *gin.Context) {
	// Get visibility from the POST
	var body CanvasBranchVisibilityPost
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if !utils.SliceContainsItem([]string{"private", "view", "comment", "edit"}, body.Visibility) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Allowed Scopes: private, view, comment, edit",
		})
		return
	}
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)

	// Permission check
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, canvasBranchID, permissiongroup.CANVAS_BRANCH_MANAGE_PERMS); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	errUpdatingBranh := App.Service.UpdateCanvasBranchVisibility(canvasBranchID, authUserID, body.Visibility)
	if errUpdatingBranh != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Saving Failed",
		})
		return
	}

	// updating canvas language pages visibility
	App.Service.UpdateCanvasLanguageBranchesVisibility(canvasBranchID, authUserID, body.Visibility)

	go func() {
		inherit := c.Query("inherit")
		canvasBranch, _ := queries.App.BranchQuery.GetBranchByID(canvasBranchID)
		if inherit == "true" {
			canvasRepos, _ := App.Repo.GetCanvasRepos(map[string]interface{}{"parent_canvas_repository_id": canvasBranch.CanvasRepositoryID})
			for _, repo := range *canvasRepos {
				App.Service.UpdateCanvasBranchVisibility(*repo.DefaultBranchID, authUserID, body.Visibility)
			}
		}
		if body.Visibility == models.PRIVATE {
			permissions.App.Service.UpdateParentHasCanvasRepoOnPrivate(canvasBranchID)
		} else {
			permissions.App.Service.UpdateParentHasCanvasRepoOnPublic(canvasBranchID)
		}
		payload, _ := json.Marshal(map[string]uint64{"collectionId": canvasBranch.CanvasRepository.CollectionID})
		apiClient.AddToQueue(apiClient.UpdateDiscordTreeMessage, payload, apiClient.DEFAULT, apiClient.CommonRetry)
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "Updated",
	})
	return

}

// @Summary 	Delete Branch
// @Tags		CanvasBranch
// @Accept       json
// @Produce      json
// @Router 		/v1/canvas-branch/{canvasBranchID} [delete]
func (r *canvasBranchRoutes) Delete(c *gin.Context) {
	canvasBranchIDStr := c.Param("canvasBranchID")

	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, utils.Uint64(canvasBranchIDStr), permissiongroup.CANVAS_BRANCH_MANAGE_CONTENT); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	// Refactor
	canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)
	_ = App.Repo.DeleteBranch(canvasBranchID)
	App.Repo.UpdateBranchLastEdited(canvasBranchID)
	go func() {
		App.Service.InvalidateBranchBlocks(canvasBranchID)
	}()
	c.JSON(http.StatusOK, gin.H{
		"message": "Deleted",
	})
	return
}

// Bulk Block API - Crete and Update (WIP)
// @Summary 	Bulk Create / Update
// @description You can call this api to update ALL the blocks on a given branch
// @description This is an Idempotent endpoint and you need to specfiy the "scope" with each block
// @description If the scope on a block json is Empty "" it will be ignore from the operations.
// @description Allowed Scope: "create" - New Block, "update" - Update Block, "delete" - Delete Block (Softdeleteonly)
// @Tags		CanvasBranch
// @Accept       json
// @Produce      json
// @Param 		body body 		CanvasBlockPost true "Bulk"
// @Router 		/v1/canvas-branch/{canvasBranchID}/blocks [post]
func (r *canvasBranchRoutes) UpdateBlocks(c *gin.Context) {
	var body CanvasBlockPost
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Get UserID
	userInstance, userNotFound := r.GetLoggedInUser(c)
	if userNotFound != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": userNotFound.Error(),
		})
		return
	}
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, err := strconv.ParseUint(canvasBranchIDStr, 10, 64)

	//branchInstance, _ := App.Repo.Get(map[string]interface{}{"id": canvasBranchID})
	// @todo add the perm check to return 403 if it is not rough branch or not published branch.
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, canvasBranchID, permissiongroup.CANVAS_BRANCH_EDIT); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	studioID, _ := r.GetStudioId(c)
	studio, err := queries.App.StudioQueries.GetStudioQuery(map[string]interface{}{
		"id": studioID,
	})
	// Get StudioID
	meta, err1 := App.Controller.BlocksManager(*userInstance, canvasBranchID, body, false, studio)
	fmt.Println(err1, meta)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err1.Error(),
		})
		return
	}

	/* Commented 16 Aug: Chirag
	blocks, errGettingBlocks := App.Service.GetAllBlockByBranchID(canvasBranchID)
	if errGettingBlocks != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errGettingBlocks.Error(),
		})
		return
	}

	response.RenderResponse(c, BulkBlocksSerializerData(blocks, meta))
	*/
	blocks, errGettingBlocks := App.Service.GetAllBlockByBranchID(canvasBranchID)
	if errGettingBlocks != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errGettingBlocks.Error(),
		})
		return
	}
	blocksData, err := App.Service.addReactionsToBlocks(blocks, canvasBranchID, userInstance)
	if errGettingBlocks != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errGettingBlocks.Error(),
		})
		return
	}
	go func() {
		App.Service.CacheBranchBlocks(canvasBranchID, authUserID, blocksData)
	}()
	response.RenderResponse(c, blocksData)
}

// Bulk Block Association API
// @Tags		CanvasBranch
// @Accept       json
// @Produce      json
// @Param 		body body 		CanvasBlockAssociationPost true "Bulk"
// @Router 		/v1/canvas-branch/{canvasBranchID}/blocks/associations [post]
func (r *canvasBranchRoutes) AssociatesBlock(c *gin.Context) {
	// update the block with what data is being sent
	// Expect json of Blocks do Excat -> PONLY UPDATE (Add Perms)
	//
	var body CanvasBlockAssociationPost
	if err := c.ShouldBindJSON(&body); err != nil {
		fmt.Println("Error ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Get UserID
	LoggedInUserInstance, userNotFound := r.GetLoggedInUser(c)
	if userNotFound != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": userNotFound.Error(),
		})
		return
	}
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, err := strconv.ParseUint(canvasBranchIDStr, 10, 64)
	// Permission check
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(LoggedInUserInstance.ID, canvasBranchID, permissiongroup.CANVAS_BRANCH_ADD_COMMENT); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	errValidation := body.Validate()
	if errValidation != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	var x CanvasBlockPost
	x.Blocks = body.Blocks
	// Get StudioID
	studioID, _ := r.GetStudioId(c)
	studio, err := queries.App.StudioQueries.GetStudioQuery(map[string]interface{}{
		"id": studioID,
	})
	meta, err1 := App.Controller.BlocksManager(*LoggedInUserInstance, canvasBranchID, x, true, studio)
	fmt.Println(meta)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err1.Error(),
		})
		return
	}
	/*
		blocks, errGettingBlocks := App.Service.GetAllBlockByBranchID(canvasBranchID)
		if errGettingBlocks != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errGettingBlocks.Error(),
			})
			return
		}
		response.RenderResponse(c, BulkBlocksSerializerData(blocks, meta))
		*
	*/

	blocks, errGettingBlocks := App.Service.GetAllBlockByBranchID(canvasBranchID)
	if errGettingBlocks != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errGettingBlocks.Error(),
		})
		return
	}
	blocksData, err := App.Service.addReactionsToBlocks(blocks, canvasBranchID, LoggedInUserInstance)
	if errGettingBlocks != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errGettingBlocks.Error(),
		})
		return
	}
	response.RenderResponse(c, blocksData)
}

// Get All Blocks on a Given Branch
// @Summary 	Get Many Blocks on a Branch
// @description You can call this api to Get ALL the blocks on a given branch
// @Tags		CanvasBranch
// @Accept       json
// @Produce      json
// @Router 		/v1/canvas-branch/{canvasBranchID}/blocks [get]
func (r *canvasBranchRoutes) GetBlocks(c *gin.Context) {
	// Get UserID
	userInstance, _ := r.GetLoggedInUser(c)
	authUserID := r.GetLoggedInUserId(c)
	// if userNotFound != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": userNotFound.Error(),
	// 	})
	// 	return
	// }
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, err := strconv.ParseUint(canvasBranchIDStr, 10, 64)
	fmt.Println(err)
	// Permission check
	// if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(userInstance.ID, canvasBranchID, permissiongroup.CANVAS_BRANCH_VIEW); err != nil || !hasPermission {
	// 	response.RenderPermissionError(c)
	// }
	// if err != nil {
	// 	response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
	// 	return
	// }

	// Chirag @todo: Need to implement Redis for This.
	/*
		cachedBranchKeyName := blocks2.App.Service.MakeBranchRedisKey(canvasBranchIDStr)
		cacheExists := blocks2.App.Service.DoesKeyExists(cachedBranchKeyName)
		fmt.Println("The Cache Key Exists ", cacheExists)
		if cacheExists {
			fmt.Println("get the cache ")
		} else {
			fmt.Println("buid the cache.")
		}*/
	blocksCachedData, err := App.Service.GetBranchBlocksCachedData(canvasBranchID, authUserID)
	if blocksCachedData != nil {
		response.RenderResponse(c, blocksCachedData)
		return
	}
	blocks, errGettingBlocks := App.Service.GetAllBlockByBranchID(canvasBranchID)
	if errGettingBlocks != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errGettingBlocks.Error(),
		})
		return
	}
	blocksData, err := App.Service.addReactionsToBlocks(blocks, canvasBranchID, userInstance)
	if errGettingBlocks != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errGettingBlocks.Error(),
		})
		return
	}
	go func() {
		App.Service.CacheBranchBlocks(canvasBranchID, authUserID, blocksData)
	}()
	response.RenderResponse(c, blocksData)
}

// Get history of commits of the Main branch
// @Summary 	Get history of commits of the Main branch
// @description Get history of commits of the Main branch
// @Tags		CanvasBranch
// @Accept       json
// @Produce      json
// @Router 		/v1/canvas-branch/branch-ops/{canvasBranchID}/history [get]
func (r *canvasBranchRoutes) BranchHistory(c *gin.Context) {
	// FetchBranchHistoryFromGit
	userInstance, userNotFound := r.GetLoggedInUser(c)
	if userNotFound != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": userNotFound.Error(),
		})
		return
	}
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)

	// permission check
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(userInstance.ID, canvasBranchID, permissiongroup.CANVAS_BRANCH_VIEW); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	startCommitID := c.Query("startCommitID")

	logs, users, next, err := App.Service.BranchHistory(userInstance, canvasBranchID, startCommitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs": GitLogsSerializerData(logs, users),
		"next": next,
	})

}

// @Summary 	Get Branches which are Rough and Unpublished by Studio
// @Tags		CanvasBranch
// @Accept       json
// @Produce      json
// @Router 		/v1/canvas-branch/branch-ops/drafts [get]
func (r *canvasBranchRoutes) RoughAndUnpublishedByStudio(c *gin.Context) {
	userInstance, userNotFound := r.GetLoggedInUser(c)
	if userNotFound != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": userNotFound.Error(),
		})
		return
	}
	// Get StudioID
	studioID, studioNotFoundErr := r.GetStudioId(c)
	if studioNotFoundErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": studioNotFoundErr.Error(),
		})
		return
	}

	instances, err := App.Service.draftBranches(studioID, userInstance.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	response.RenderResponse(c, MultiSerializeCanvasBranch(*instances))
	return

}

// BlockHistoryByCommitID: Get blocks from git history by commit id
// @Summary 	Get blocks from git history by commit id
// @Tags		CanvasBranch
// @Accept       json
// @Produce      json
// @Router 		/v1/canvas-branch/{canvasBranchID}/blocks/{commitID}/blocks-history [get]
func (r *canvasBranchRoutes) BlockHistoryByCommitID(c *gin.Context) {

	userInstance, userNotFound := r.GetLoggedInUser(c)
	if userNotFound != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": userNotFound.Error(),
		})
		return
	}
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)

	// permission check
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(userInstance.ID, canvasBranchID, permissiongroup.CANVAS_BRANCH_VIEW); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	commitID := c.Param("commitID")

	blocks, err := App.Service.GetBlocksFromHistoryByCommitId(userInstance, commitID, canvasBranchID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	blocksData := BulkGitBlocksSerializerData(blocks)
	response.RenderResponse(c, blocksData)
}

// BuildRoughBranch: Create a Rough branch for this user.
// @Summary 	Create a rough branch for this user.
// @description When you call call this API it will create a copy of all the blocks
// @description This is a rough branch but is only for a Logged in user
// @Tags		CanvasBranch
// @Accept       json
// @Produce      json
// @Param 		body body 		NewDraftBranchPost true "Bulk"
// @Router 		/v1/canvas-branch/branch-ops/{canvasBranchID}/rough-branch [post]
func (r *canvasBranchRoutes) BuildRoughBranch(c *gin.Context) {
	var body NewDraftBranchPost
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Get UserID
	userInstance, userNotFound := r.GetLoggedInUser(c)
	if userNotFound != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": userNotFound.Error(),
		})
		return
	}
	// Get StudioID
	studioID, studioNotFoundErr := r.GetStudioId(c)
	if studioNotFoundErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": studioNotFoundErr.Error(),
		})
		return
	}

	// query 1 : GetCanvasRepoInstance
	canvasRepo, _ := App.Repo.GetCanvasRepoInstance(map[string]interface{}{"id": body.CanvasRepoID})
	// collectionId may be wrong sometimes when user doesn't refresh the screen. so we are taking it from BE
	body.CollectionID = canvasRepo.CollectionID
	if canvasRepo.ParentCanvasRepositoryID != nil {
		body.ParentCanvasRepoID = *canvasRepo.ParentCanvasRepositoryID
	} else {
		body.ParentCanvasRepoID = 0
	}

	// Get the canvas branch
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)
	// permission check
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(userInstance.ID, canvasBranchID, permissiongroup.CANVAS_BRANCH_EDIT); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	// send the branch and current user for new draft branch
	RoughBranchInstance, errCreatingRoughBranch := App.Controller.RoughBranchBuilder(userInstance, canvasBranchID, studioID, body)
	if errCreatingRoughBranch != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errCreatingRoughBranch.Error(),
		})
		return
	}

	//branch, err := App.Repo.Get(map[string]interface{}{"id": RoughBranchID})
	//branch, err := queries.App.BranchQuery.GetBranchByID.GetBranchWithRepoAndStudio(RoughBranchID)
	//if err != nil {
	//	response.RenderErrorResponse(c, err)
	//	return
	//}
	response.RenderResponse(c, SerializeCanvasBranch(RoughBranchInstance))
	return

}

// GetCanvasBranches All branches on a Given Canvas
// @Summary 	Get branches of a block.
// @description You can call this api to Get ALL the branches on a given canvas
// @Tags		CanvasBranch
// @Accept      json
// @Produce     json
// @Param 		body 	body 		GetCanvasBranches true "Data"
// @Router 		/v1/canvas-branch/branch-ops/nav/get-branches [post]
func (r *canvasBranchRoutes) GetCanvasBranches(c *gin.Context) {
	var body GetCanvasBranches
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, _ := r.GetLoggedInUser(c)
	studioID, err := r.GetStudioId(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	branches := &[]GetCanvasBranchSerializer{}
	if user == nil {
		branches, err = App.Controller.AnonymousBranchesController(body)
	} else {
		branches, err = App.Controller.AuthGetCanvasBranchesController(body, user, studioID)
	}
	fmt.Println(branches)
	if err != nil {
		response.RenderErrorResponse(c, err)
		return
	}

	response.RenderResponse(c, branches)
	return
}

// @Summary 	Get Root of a Given Branch
// @description Given a BranchID we will return Collection -> Repos -> BranchID
// @Tags		CanvasBranch
// @Accept      json
// @Produce     json
// @Param 		body 	body 		GetCanvasBranches true "Data"
// @Router 		/v1/canvas-branch/branch-ops/nav/:canvasBranchID/root [get]
func (r *canvasBranchRoutes) GetCanvasBranchesRoots(c *gin.Context) {
	// user, _ := r.GetLoggedInUser(c)
	// Get the canvas branch
	authUserID := r.GetLoggedInUserId(c)
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)
	public := c.Query("public")
	resp, err := App.Service.BuildRootNavObject(canvasBranchID, authUserID, public)
	if err != nil {
		response.RenderErrorResponse(c, err)
		return
	}
	response.RenderResponse(c, resp)
}

// @Summary 	Get Children of a Given Branch
// @description Given a BranchID we will Children of a Given Node
// @Tags		CanvasBranch
// @Accept      json
// @Produce     json
// @Param 		body 	body 		GetCanvasBranches true "Data"
// @Router 		/v1/canvas-branch/branch-ops/nav/:canvasBranchID/node [get]
func (r *canvasBranchRoutes) GetCanvasBranchesNodes(c *gin.Context) {
	// user, _ := r.GetLoggedInUser(c)
	// Get the canvas branch
	authUserID := r.GetLoggedInUserId(c)
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)
	resp, err := App.Service.BuildRootNodeObject(canvasBranchID, authUserID)
	if err != nil {
		response.RenderErrorResponse(c, err)
		return
	}
	response.RenderResponse(c, resp)
}

// Search Repos and Branches.
// @Summary 	Search Repos and Branches.
// @description Search by query
// @Tags		CanvasBranch
// @Accept      json
// @Produce     json
// @Param 		body 	body 		SearchBranchRepos true "Data"
// @Router 		/v1/canvas-branch/branch-ops/nav/search [post]
func (r *canvasBranchRoutes) SearchCanvasBranches(c *gin.Context) {
	var body SearchBranchRepos
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	userID := r.GetLoggedInUserId(c)
	studioID, _ := r.GetStudioId(c)
	repoAndCollectionRows := App.Repo.QueryDB(body.Query, studioID)
	fmt.Println(repoAndCollectionRows)
	public := c.Query("public")
	resp := App.Service.ProcessSearchDump(repoAndCollectionRows, userID, studioID, public)
	c.JSON(http.StatusOK, resp)
	//
	//user, err := r.GetLoggedInUser(c)
	//if err != nil {
	//	response.RenderErrorResponse(c, err)
	//	return
	//}
	//studioID, err := r.GetStudioId(c)
	//if err != nil {
	//	response.RenderErrorResponse(c, err)
	//	return
	//}
	//fmt.Println(body)
	//branches := &[]GetCanvasBranchSerializer{}
	//if user == nil {
	//	branches, err = App.Controller.AnonymousBranchesController(body)
	//} else {
	//	branches, err = App.Controller.AuthGetCanvasBranchesController(body, user, studioID)
	//}
	//fmt.Println(branches)
	//if err != nil {
	//	response.RenderErrorResponse(c, err)
	//	return
	//}
	//
	//response.RenderResponse(c, branches)
	return
}

// Get All Attributions of a canvasbranch
// @Summary 	Get CanvasBranch All Attributions
// @Tags		CanvasBranch
// @Accept       json
// @Produce      json
// @Router 		/v1/canvas-branch/attributions/{canvasBranchID} [get]
func (r *canvasBranchRoutes) Attributions(c *gin.Context) {
	userInstance, _ := r.GetLoggedInUser(c)
	// if userNotFound != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": userNotFound.Error(),
	// 	})
	// 	return
	// }
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)

	// authUserID := r.GetLoggedInUserId(c)
	// if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, canvasBranchID, permissiongroup.CANVAS_BRANCH_VIEW); err != nil || !hasPermission {
	// 	response.RenderPermissionError(c)
	// 	return
	// }

	attributions, err := App.Service.GetBranchAttributions(userInstance, canvasBranchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": GitAttributionsSerializerData(attributions),
	})
}

// @Summary API will returns Merge Request reponse before merge request is created.
// @Tags		CanvasBranch
// @Accept       json
// @Produce      json
// @Router 		/v1/canvas-branch/{canvasBranchID}/diffblocks [get]
func (r *canvasBranchRoutes) DiffBlocks(c *gin.Context) {
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)
	user, err := r.GetLoggedInUser(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	// permission check
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, canvasBranchID, permissiongroup.CANVAS_BRANCH_CREATE_MERGE_REQUEST); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	branch, errGettingBranch := App.Service.GetCanvasBranchInstance(canvasBranchID)
	if errGettingBranch != nil {
		response.RenderCustomErrorResponse(c, errGettingBranch)
		return
	}
	resp, errMergeService := App.Service.BlockBeforeMergeService(branch, user)
	if errMergeService != nil {
		response.RenderCustomErrorResponse(c, errMergeService)
		return
	}
	response.RenderResponse(c, resp)

}

// GetRepoBranch Canvas Branch & repo Metadata
// @Summary 	Get branch lastupdated at.
// @Tags		CanvasBranch
// @Router 		/v1/canvas-branch/{canvasBranchID}/last-updated  [get]
func (r *canvasBranchRoutes) GetBranchLastUpdated(c *gin.Context) {
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)
	branch, err := App.Repo.GetBranchWithRepo(map[string]interface{}{"id": canvasBranchID})
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	resp := map[string]interface{}{
		"id":            branch.ID,
		"lastUpdatedAt": branch.UpdatedAt,
	}
	response.RenderCustomResponse(c, resp)
}
