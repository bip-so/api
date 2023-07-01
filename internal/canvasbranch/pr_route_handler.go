package canvasbranch

import (
	"encoding/json"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/canvasrepo"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/internal/xpcontribs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/apiClient"
)

// This API is overloaded
// It has 2 succes responses
// Create a Publish Request or Directly do a PR
// Till PR a Branch does not gett any Git functions.

// @Summary 	Publish a branch or Request a publish request (based on your permissions)
// @description This will eitehr create a PR or Directy you can check the `"published": false,` in the json response
// @description if the publish was done or a request was created
// @Tags		CanvasBranch
// @Accept       json
// @Produce      json
// @Param 		body body 		InitPRPost true "Publish"
// @Router 		/v1/canvas-branch/{canvasBranchID}/publish-request/init [post]
func (r *canvasBranchRoutes) InitPr(c *gin.Context) {
	var body InitPRPost
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	studioID, _ := r.GetStudioId(c)
	nudgeFlag, repoPrivateCount, maxRepoPrivateCount := canvasrepo.App.Repo.StudioPlanCheck(studioID)
	// We need to check if repoPrivateCount > maxRepoPrivateCount
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)
	authUser, _ := r.RouteHelper.GetLoggedInUser(c)
	instance, canCreatePR, err := App.Service.InitPrRequest(canvasBranchID, authUser.ID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if canCreatePR {
		fmt.Println("Direct Publish ")
		errDirectPR := App.Service.DirectPublishRequest(*instance, body.Message, *authUser)
		if errDirectPR != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errDirectPR.Error(),
			})
			return
		}
		go func() {
			canvasBranch, _ := App.Repo.Get(map[string]interface{}{"id": canvasBranchID})
			payload, _ := json.Marshal(map[string]uint64{"collectionId": canvasBranch.CanvasRepository.CollectionID})
			apiClient.AddToQueue(apiClient.UpdateDiscordTreeMessage, payload, apiClient.DEFAULT, apiClient.CommonRetry)
			if nudgeFlag && int(repoPrivateCount) > maxRepoPrivateCount {
				// We need to add this Repo / Branch to Redis for next day.
				canvasrepo.App.PlansAutoPublish.AddToAutoPublishQueue(canvasBranchID, instance.CanvasRepositoryID, studioID, authUser.ID)
			}
			if nudgeFlag && (repoPrivateCount == int64(maxRepoPrivateCount+1) || (repoPrivateCount/5 > 5 && repoPrivateCount%5 == 0)) {
				notifications.App.Service.PublishNewNotification(notifications.CanvasLimitExceed, 0, []uint64{}, &studioID,
					nil, notifications.NotificationExtraData{}, nil, nil)
			}
		}()
		xpcontribs.App.Service.XPSummarizer(studioID)
		c.JSON(http.StatusOK, gin.H{
			//"branch_id": parentBranchID,
			"published":             true,
			"message":               "Hurray Branch is Published!!!",
			"nudge":                 nudgeFlag,
			"privateRepoCount":      repoPrivateCount,
			"maxAllowedPrivateRepo": maxRepoPrivateCount,
		})
		return
	}
	fmt.Println("Publish Request  ")
	pr, errCreatePR := App.Service.CreatePublishRequest(*instance, body.Message, authUser.ID)
	if errCreatePR != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errCreatePR.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		//"branch_id": parentBranchID,
		"publishRequestId":      pr.ID,
		"published":             false,
		"message":               "We have created publish request, awaiting for moderators to approve",
		"nudge":                 nudgeFlag,
		"privateRepoCount":      repoPrivateCount,
		"maxAllowedPrivateRepo": maxRepoPrivateCount,
	})

}

// @Summary 	Get All the Publish Requests on a Branch
// @Tags		CanvasBranch
// @Accept       json
// @Produce      json
// @Router 		/v1/canvas-branch/{canvasBranchID}/publish-request/list [get]
func (r *canvasBranchRoutes) ListPublishRequests(c *gin.Context) {
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)
	studioID, _ := r.GetStudioId(c)
	userID := r.GetLoggedInUserId(c)
	instances, err := App.Service.GetPublishRequestsByBranch(studioID, canvasBranchID, userID)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}

	response.RenderResponse(c, PublishRequestSerializerMany(instances))
}

// @Summary 	Accept and Reject a Publish Request (Check Perms)
// @Tags		CanvasBranch
// @Accept       json
// @Produce      json
// @Param 		body body 		ManagePublishRequest true "Publish Req"
// @Router 		/v1/canvas-branch/{canvasBranchID}/publish-request/:publishRequestID/manage [post]
func (r *canvasBranchRoutes) ManagePR(c *gin.Context) {
	var body ManagePublishRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)
	publishRequestIDStr := c.Param("publishRequestID")
	publishRequestID, _ := strconv.ParseUint(publishRequestIDStr, 10, 64)
	authUser, _ := r.RouteHelper.GetLoggedInUser(c)

	// permission check
	studioID, _ := r.GetStudioId(c)
	nudgeFlag, repoPrivateCount, maxRepoPrivateCount := canvasrepo.App.Repo.StudioPlanCheck(studioID)

	authUserID := r.GetLoggedInUserId(c)
	branchInstance, _ := queries.App.BranchQuery.GetBranchByID(canvasBranchID)
	if branchInstance.CanvasRepository.ParentCanvasRepositoryID != nil {
		if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, *branchInstance.CanvasRepository.ParentCanvasRepository.DefaultBranchID, permissiongroup.CANVAS_BRANCH_MANAGE_PUBLISH_REQUESTS); err != nil || !hasPermission {
			response.RenderPermissionError(c)
			return
		}
	} else {
		if hasPermission, err := permissions.App.Service.CanUserDoThisOnCollection(authUserID, studioID, branchInstance.CanvasRepository.CollectionID, permissiongroup.COLLECTION_MANAGE_PUBLISH_REQUEST); err != nil || !hasPermission {
			response.RenderPermissionError(c)
			return
		}
	}

	err := App.Service.ManagePublishRequest(canvasBranchID, publishRequestID, authUser, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if body.Accept {
		xpcontribs.App.Service.XPSummarizer(studioID)
		c.JSON(http.StatusOK, gin.H{
			"published":             true,
			"message":               "Publish Request is updated.",
			"nudge":                 nudgeFlag,
			"privateRepoCount":      repoPrivateCount,
			"maxAllowedPrivateRepo": maxRepoPrivateCount,
		})
		go func() {
			if nudgeFlag && int(repoPrivateCount) > maxRepoPrivateCount {
				// We need to add this Repo / Branch to Redis for next day.
				canvasrepo.App.PlansAutoPublish.AddToAutoPublishQueue(canvasBranchID, branchInstance.CanvasRepositoryID, studioID, authUser.ID)
			}
			if nudgeFlag && (repoPrivateCount == int64(maxRepoPrivateCount+1) || (repoPrivateCount/5 > 5 && repoPrivateCount%5 == 0)) {
				notifications.App.Service.PublishNewNotification(notifications.CanvasLimitExceed, 0, []uint64{}, &studioID,
					nil, notifications.NotificationExtraData{}, nil, nil)
			}
		}()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":               "Publish Request is updated.",
		"nudge":                 nudgeFlag,
		"privateRepoCount":      repoPrivateCount,
		"maxAllowedPrivateRepo": maxRepoPrivateCount,
	})

}

// @Summary 	API is used by SAME person to delete the PR
// @Tags		CanvasBranch
// @Accept       json
// @Produce      json
// @Router 		/v1/canvas-branch/:canvasBranchID/publish-request/:publishRequestID/delete [delete]
func (r *canvasBranchRoutes) DeletePR(c *gin.Context) {
	// Check me and Delete PR and check Alrealdy published
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)
	publishRequestIDStr := c.Param("publishRequestID")
	publishRequestID, _ := strconv.ParseUint(publishRequestIDStr, 10, 64)
	authUser, _ := r.RouteHelper.GetLoggedInUser(c)

	prInstance, errGettingInstance := queries.App.PublishRequestQuery.PublishRequestGetter(canvasBranchID, authUser.ID, publishRequestID)
	// Check if prInstance
	if errGettingInstance != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Instance not found!",
			"error":   errGettingInstance.Error(),
		})
		return
	}

	// ONLY User can delete user's PR
	if prInstance.CreatedByID != authUser.ID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Only user can deleted their own PRs",
		})
		return
	}
	// if status is not pending
	if prInstance.Status == models.PUBLISH_REQUEST_ACCEPTED || prInstance.Status == models.PUBLISH_REQUEST_REJECTED {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "PR has already been " + string(prInstance.Status),
		})
		return
	}
	// Delete PR
	errDeleting := queries.App.PublishRequestQuery.DeletePR(prInstance.ID)
	if errDeleting != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   errDeleting.Error(),
			"message": "Error deleting PR",
		})
		return
	}

	// Adding a task to delete all the merge request notifications
	prInstanceStr, _ := json.Marshal(prInstance)
	go apiClient.AddToQueue(apiClient.DeletePublishRequestNotifications, []byte(prInstanceStr), apiClient.DEFAULT, apiClient.CommonRetry)
	go apiClient.AddToQueue(apiClient.DeleteModsOnCanvas, []byte(prInstanceStr), apiClient.DEFAULT, apiClient.CommonRetry)

	c.JSON(http.StatusOK, gin.H{
		"message": "Publish Request is Deleted.",
	})
}
