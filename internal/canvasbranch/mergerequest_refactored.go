package canvasbranch

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/internal/xpcontribs"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
)

// @Summary Create a Merge Request, reequires 'CANVAS_BRANCH_CREATE_MERGE_REQUEST'
// @Tags		CanvasBranch
// @Accept       json
// @Produce      json
// @Param 		body body 		MergeRequestCreatePost true "Cancel"
// @Router 		/v1/canvas-branch/{canvasBranchID}/merge-request/create [post]
func (r *canvasBranchRoutes) CreateMergeRequest(c *gin.Context) {

	// Process Post and Bind
	var body MergeRequestCreatePost
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	isAutoMerge := c.Query("merge") == "true"
	// Build Needede vars
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)
	//thisCanvasBranchInstance :=
	authUser, errGettingUser := r.GetLoggedInUser(c)

	if errGettingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errGettingUser.Error(),
		})
		return
	}

	roughBranchInstance, _ := queries.App.BranchQuery.GetBranchWithRepoAndStudio(canvasBranchID)

	// Check Permissions upfront
	// Create a Merge Request if User Passes for permissiongroup.CANVAS_BRANCH_CREATE_MERGE_REQUEST Permissions
	if isAutoMerge {

		//permCheckBranch, _ := App.Repo.Get(map[string]interface{}{"id": canvasBranchID})
		//currentBranchInstance, _ := queries.App.BranchQuery.GetBranchByID.GetBranchWithRepoAndStudio(canvasBranchID)

		var permCheckBranchID uint64
		if roughBranchInstance.RoughFromBranchID != nil {
			permCheckBranchID = *roughBranchInstance.RoughFromBranchID
		} else {
			permCheckBranchID = *roughBranchInstance.FromBranchID
		}
		// We need to check is user has CANVAS_BRANCH_MANAGE_MERGE_REQUESTS permission
		canUserPerformMergeRequest, errGettingPermissions := permissions.App.Service.CanUserDoThisOnBranch(authUser.ID, permCheckBranchID, permissiongroup.CANVAS_BRANCH_MANAGE_MERGE_REQUESTS)
		if errGettingPermissions != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errGettingPermissions.Error(),
			})
			return
		}

		if !canUserPerformMergeRequest {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "User does not have permissions to create a merge request",
			})
			return
		}
	} else {
		canUserCreateMergeRequest, errGettingPermissions := permissions.App.Service.CanUserDoThisOnBranch(authUser.ID, canvasBranchID, permissiongroup.CANVAS_BRANCH_CREATE_MERGE_REQUEST)
		if errGettingPermissions != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errGettingPermissions.Error(),
			})
			return
		}
		if !canUserCreateMergeRequest {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "User does not have permissions to create a merge request",
			})
			return
		}
	}

	// Validate the branch instance
	// branchInstance, errValidation := App.Service.MergeRequestCreationValidation(canvasBranchID, authUser.ID)
	validatedBranchInstance, errValidation := App.Service.MergeRequestCreationValidation(roughBranchInstance, authUser.ID)
	if errValidation != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errValidation.Error(),
		})
		return
	}

	// Todo: CC Move after merge
	var mergeRequestInstance *models.MergeRequest
	var errCreatingMergeRequest error

	// Todo: FIX
	if validatedBranchInstance.RoughBranchCreatorID != nil {
		mergeRequestInstance, errCreatingMergeRequest = App.Service.CreateMergeRequest(validatedBranchInstance.ID, *validatedBranchInstance.RoughFromBranchID, authUser, body.CommitMessage)
	} else {
		mergeRequestInstance, errCreatingMergeRequest = App.Service.CreateMergeRequest(validatedBranchInstance.ID, *validatedBranchInstance.FromBranchID, authUser, body.CommitMessage)
	}
	if errCreatingMergeRequest != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errCreatingMergeRequest.Error(),
		})
		return
	}

	if isAutoMerge {
		var errMr error
		mergeRequestInstance, errMr = App.Repo.GetMergeRequestWithPreloads(map[string]interface{}{"id": mergeRequestInstance.ID})
		if errMr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Merge Instance not found",
			})
			return
		}
		_, errMergingBranches := App.Git.MergeMergeRequest(authUser, mergeRequestInstance.SourceBranchID, mergeRequestInstance.DestinationBranchID, MergeRequestAcceptPartialPost{
			MergeStatus:     models.MERGE_REQUEST_ACCEPTED,
			ChangesAccepted: nil,
			CommitMessage:   body.CommitMessage,
		}, mergeRequestInstance, isAutoMerge)
		if errMergingBranches != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errMergingBranches.Error(),
			})
			return
		}

		mergeRequestInstance.Status = models.MERGE_REQUEST_ACCEPTED

		// to run the notifications which are created on rough branch
		go func() {
			// on success on a merge we are calling XP Summarizer
			studioID, _ := r.GetStudioId(c)
			xpcontribs.App.Service.XPSummarizer(studioID)
			// For executing notifications which are created on rough branch.
			notifications.App.Service.ExecuteAllRoughBranchNotifications(roughBranchInstance.ID)
		}()

		c.JSON(http.StatusOK, gin.H{
			"message": "We have successfully merged branches!!!",
		})
		return

	}
	go func() {
		extraData := notifications.NotificationExtraData{
			CollectionID:   validatedBranchInstance.CanvasRepository.CollectionID,
			CanvasRepoID:   validatedBranchInstance.CanvasRepositoryID,
			CanvasBranchID: validatedBranchInstance.ID,
			Status:         models.MERGE_REQUEST_OPEN,
		}
		contentObject := models.MERGEREQUEST
		notifications.App.Service.PublishNewNotification(notifications.MergeRequested, authUser.ID, nil, &validatedBranchInstance.CanvasRepository.StudioID,
			nil, extraData, &mergeRequestInstance.ID, &contentObject)
	}()

	c.JSON(http.StatusOK, gin.H{
		"message":                "We have successfully created merge request!!!",
		"merge_request_instance": SimpleMergeRequestSerializer(mergeRequestInstance),
	})

	fmt.Println("CreateMergeRequest: TIME")
	defer utils.TimeTrack(time.Now())

	return
}

// @Summary Accept / Partially Accept a Merge request, Requires CANVAS_BRANCH_MANAGE_MERGE_REQUESTS
// @Tags		CanvasBranch
// @Accept       json
// @Produce      json
// @Param 		body body 		MergeRequestAcceptPartialPost true "Accept/Partial Accept"
// @Router 		/v1/canvas-branch/{canvasBranchID}/merge-request/:mergeRequestID/merge-accept [post]
func (r *canvasBranchRoutes) MergeRequestAcceptPartial(c *gin.Context) {
	// Process Post and Bind
	var body MergeRequestAcceptPartialPost
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Build Needede vars
	//canvasBranchIDStr := c.Param("canvasBranchID")
	//canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)
	authUser, errGettingUser := r.GetLoggedInUser(c)
	if errGettingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errGettingUser.Error(),
		})
		return
	}

	mergeRequestIDStr := c.Param("mergeRequestID")
	mergeRequestID, _ := strconv.ParseUint(mergeRequestIDStr, 10, 64)
	mergeRequestInstance, errMR := App.Repo.GetMergeRequestWithPreloads(map[string]interface{}{"id": mergeRequestID})
	if errMR != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Merge Instance not found",
		})
		return
	}

	branchInstance, errValidation := App.Service.MergeRequestAcceptValidation(mergeRequestInstance.SourceBranch, authUser.ID)
	fmt.Println(branchInstance)
	if errGettingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errValidation.Error(),
		})
		return
	}

	// Note: Permissions is checked on the Branch from which this branch was created (Branch -> RoughFromBranchID)
	canManageMergeRequest, errGettingManageMergeRequestPerms := permissions.App.Service.CanUserDoThisOnBranch(authUser.ID, mergeRequestInstance.DestinationBranchID, permissiongroup.CANVAS_BRANCH_MANAGE_MERGE_REQUESTS)
	if errGettingManageMergeRequestPerms != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errGettingManageMergeRequestPerms.Error(),
		})
		return
	}

	if !canManageMergeRequest {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User does not have permissions to Merge a Branch",
		})
		return
	}
	// start merge
	_, errMergingBranches := App.Git.MergeMergeRequest(authUser, mergeRequestInstance.SourceBranchID, mergeRequestInstance.DestinationBranchID, body, mergeRequestInstance, false)
	if errMergingBranches != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errMergingBranches.Error(),
		})
		return
	}
	// to run the notifications which are created on rough branch
	go func() {
		notifications.App.Service.ExecuteAllRoughBranchNotifications(branchInstance.ID)
	}()
	c.JSON(http.StatusOK, gin.H{
		"branch_id": mergeRequestInstance.DestinationBranchID,
		"message":   "Hurray Branches Are Merged Now!!!",
	})
	return

}

// @Summary Use this endpoint to Reject a MergeRequest needs to have CANVAS_BRANCH_MANAGE_MERGE_REQUESTS
// @Tags		CanvasBranch
// @Accept       json
// @Produce      json
// @Router 		/v1/canvas-branch/{canvasBranchID}/merge-request/:mergeRequestID/reject [post]
func (r *canvasBranchRoutes) MergeRequestRejected(c *gin.Context) {
	// Create a Merge Request
	// Process Post and Bind
	// var body MergeRequestMergeRequestRejectedPost
	// if err := c.ShouldBindJSON(&body); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	// Build Needede vars
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)
	authUser, errGettingUser := r.GetLoggedInUser(c)
	if errGettingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errGettingUser.Error(),
		})
		return
	}
	mergeRequestIDStr := c.Param("mergeRequestID")
	mergeRequestID, _ := strconv.ParseUint(mergeRequestIDStr, 10, 64)
	mergeRequestInstance, errMR := App.Repo.GetMergeRequestWithPreloads(map[string]interface{}{"id": mergeRequestID})
	if errMR != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Merge Instance not found",
		})
		return
	}
	// permission check
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUser.ID, mergeRequestInstance.DestinationBranchID, permissiongroup.CANVAS_BRANCH_MANAGE_MERGE_REQUESTS); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	_, errValidation := App.Service.MergeRequestRejectValidation(mergeRequestInstance.SourceBranch, authUser.ID)
	if errGettingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errValidation.Error(),
		})
		return
	}

	// Note: Permissions is checked on the Branch from which this branch was created (Branch -> RoughFromBranchID)
	canManageMergeRequest, errGettingManageMergeRequestPerms := permissions.App.Service.CanUserDoThisOnBranch(authUser.ID, mergeRequestInstance.DestinationBranchID, permissiongroup.CANVAS_BRANCH_MANAGE_MERGE_REQUESTS)
	if errGettingManageMergeRequestPerms != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errGettingManageMergeRequestPerms.Error(),
		})
		return
	}

	if !canManageMergeRequest {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User does not have permissions to Merge a Branch",
		})
		return
	}
	_, errMergingBranches := App.Git.RejectMergeRequest(authUser, canvasBranchID, mergeRequestInstance.DestinationBranchID, mergeRequestID)
	if errMergingBranches != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errMergingBranches.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"branch_id": mergeRequestInstance.DestinationBranchID,
		"message":   "Hurray Branches Are Merged Now!!!",
	})
	return

}
