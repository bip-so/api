package canvasbranch

import (
	"context"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"net/http"
	"strconv"

	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"

	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/apiClient"
)

type InitMergeParamsOptional struct {
	OnlyRequest string `json:"onlyRequest"`
}

// @Summary Delete a Merge Request usually used by Person to delete Their own Request.
// @description This is based on persons permissions.
// @Tags		CanvasBranch
// @Accept       json
// @Produce      json
// @Param 		body body 		EmptyPost true "Delete"
// @Router 		/v1/canvas-branch/{canvasBranchID}/merge-request/:mergeRequestID/delete [post]
func (r *canvasBranchRoutes) DeleteMergeRequest(c *gin.Context) {
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)
	mergeRequestIDStr := c.Param("mergeRequestID")
	mergeRequestID, _ := strconv.ParseUint(mergeRequestIDStr, 10, 64)
	userID := r.RouteHelper.GetLoggedInUserId(c)

	// permission check
	mr, _ := App.Repo.GetMergeRequest(map[string]interface{}{"id": mergeRequestID})
	if mr.CreatedByID != userID {
		response.RenderPermissionError(c)
		return
	}

	if mr.Status != models.MERGE_REQUEST_OPEN {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "merge request is already updated",
		})
		return
	}

	errDeletingMergeRequest := App.Repo.DeleteMergeRequest(mergeRequestID)
	if errDeletingMergeRequest != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errDeletingMergeRequest.Error(),
		})
		return
	}
	// Adding a task to delete all the merge request notifications
	go apiClient.AddToQueue(apiClient.DeleteMergeRequestNotifications, []byte(mergeRequestIDStr), apiClient.DEFAULT, apiClient.CommonRetry)

	// Make the Branch Editable Again
	errUpdatingBranch := App.Repo.UpdateBranchInstance(canvasBranchID, map[string]interface{}{
		"committed":     false,
		"updated_by_id": userID,
	})

	if errUpdatingBranch != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errUpdatingBranch.Error(),
		})
		return
	}

	// delete rough branch notifications added to redis
	go App.Service.cache.Delete(context.Background(), notifications.App.Service.GetRoughBranchNotificationsRedisKey(canvasBranchID))

	c.JSON(http.StatusOK, gin.H{
		"message": "Merge request is cancelled.",
	})
}

// @Summary 	Get All the Merge Requests on this branch
// @Tags		CanvasBranch
// @Accept       json
// @Produce      json
// @Router 		/v1/canvas-branch/{canvasBranchID}/merge-request/list [get]
func (r *canvasBranchRoutes) ListMergeRequest(c *gin.Context) {
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)
	userID := r.GetLoggedInUserId(c)

	hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(userID, canvasBranchID, permissiongroup.CANVAS_BRANCH_MANAGE_MERGE_REQUESTS)
	if err != nil {
		response.RenderPermissionError(c)
		return
	}

	instances, err := App.Service.GetMergeRequestsByBranch(canvasBranchID, userID)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}

	if !hasPermission {
		var accessInstances []models.MergeRequest
		for _, instance := range *instances {
			if instance.CreatedByID == userID {
				accessInstances = append(accessInstances, instance)
			}
		}
		instances = &accessInstances
	}

	response.RenderResponse(c, MergeRequestSerializerMany(instances))
}

//:mergeRequestID/merge
func (r *canvasBranchRoutes) MergeMerge(c *gin.Context) {
	//canvasBranchIDStr := c.Param("canvasBranchID")
	//canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)
	mergeRequestIDStr := c.Param("mergeRequestID")
	mergeRequestID, _ := strconv.ParseUint(mergeRequestIDStr, 10, 64)
	userID := r.RouteHelper.GetLoggedInUserId(c)
	var body MergeMergePost
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Validation !!!
	valError := body.Validate()
	if valError != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": valError.Error(),
		})
		return
	}
	var query map[string]interface{}
	if body.Status == models.MERGE_REQUEST_PARTIALLY_ACCEPTED {
		query = map[string]interface{}{"status": body.Status, "commit_message": body.CommitMessage, "updated_by_id": userID, "changes_accepted": body.ChangesAccepted}
	} else {
		query = map[string]interface{}{"status": body.Status, "commit_message": body.CommitMessage, "updated_by_id": userID}

	}
	errUpdatingMergeRequest := App.Repo.UpdateMergeRequest(mergeRequestID, query)
	if errUpdatingMergeRequest != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errUpdatingMergeRequest.Error(),
		})
		return
	}

	// PW: Add rest of the Flow for actual Merge
}

// @Summary API will returns (Merge Request Object, Branch Object, Blocks on Rough Branch, Blocks on Parent Branch
// @Tags		CanvasBranch
// @Accept       json
// @Produce      json
// @Router 		/v1/canvas-branch/{canvasBranchID}/merge-request/{mergeRequestID}/response [get]
func (r *canvasBranchRoutes) MergeResponseRequest(c *gin.Context) {
	mergeRequestIDStr := c.Param("mergeRequestID")
	mergeRequestID, _ := strconv.ParseUint(mergeRequestIDStr, 10, 64)
	user, _ := r.GetLoggedInUser(c)
	mergeRequestInstance, errMR := App.Repo.GetMergeRequest(map[string]interface{}{"id": mergeRequestID})
	if errMR != nil {
		response.RenderCustomErrorResponse(c, errMR)
		return
	}
	authUserID := r.GetLoggedInUserId(c)
	if mergeRequestInstance.CreatedByID != authUserID {
		if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, mergeRequestInstance.DestinationBranchID, permissiongroup.CANVAS_BRANCH_MANAGE_MERGE_REQUESTS); err != nil || !hasPermission {
			response.RenderPermissionError(c)
			return
		}
	}

	resp, errMergeService := App.Service.MergeRequestService(mergeRequestInstance, user)
	if errMergeService != nil {
		response.RenderCustomErrorResponse(c, errMergeService)
		return
	}
	response.RenderResponse(c, resp)
}
