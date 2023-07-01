package ar

import (
	"github.com/gin-gonic/gin"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
	"net/http"
	"strconv"
)

// @Summary 	Create a new access request
// @Tags		CanvasBranch
// @Accept       json
// @Produce      json
// @Param 		body body 		queries.CreateAccessRequestPost true "New AR Req"
// @Router 		/v1/canvas-branch/:canvasBranchID/access-request/create [post]
func (r *arRoutes) CreateAccessRequest(c *gin.Context) {
	var body queries.CreateAccessRequestPost
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":            err.Error(),
			"access_requested": false,
		})
		return
	}
	userID := r.GetLoggedInUserId(c)
	if body.CanvasRepositoryID == 0 || body.CollectionID == 0 {
		canvasBranch, _ := queries.App.BranchQuery.GetBranchByID(body.CanvasBranchID)
		body.CanvasRepositoryID = canvasBranch.CanvasRepositoryID
		body.CollectionID = canvasBranch.CanvasRepository.CollectionID
	}
	exists, errCreateAccessRequest := App.Service.CreateAccessRequest(body, userID)
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":            errCreateAccessRequest.Error(),
			"access_requested": true,
		})
		return
	}

	if errCreateAccessRequest != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":            errCreateAccessRequest.Error(),
			"access_requested": false,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":          "Access Request is created.",
		"access_requested": false,
	})

}

// @Summary 	get all access requests by studio
// @Tags		CanvasBranch
// @Accept       json
// @Produce      json
// @Router 		/v1/canvas-branch/:canvasBranchID/access-request/list [get]
func (r *arRoutes) ListAccessRequest(c *gin.Context) {
	studioID, err := r.GetStudioId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)

	// permission check
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, canvasBranchID, permissiongroup.CANVAS_BRANCH_MANAGE_PERMS); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	instances, errList := App.Service.GetAllAccessRequest(studioID)
	if errList != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), errList))
		return
	}
	response.RenderResponse(c, AccessRequestSerializerMany(instances))

}

// @Summary 	Create a new access request
// @Tags		CanvasBranch
// @Accept       json
// @Produce      json
// @Param 		body body 		ManageAccessRequestPost true "Manage AR Req"
// @Router 		/v1/canvas-branch/:canvasBranchID/access-request/:accessRequestID/manage [post]
func (r *arRoutes) ManageAccessRequest(c *gin.Context) {
	var body ManageAccessRequestPost
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	user, _ := r.GetLoggedInUser(c)
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)
	accessRequestStr := c.Param("accessRequestID")
	accessRequestID, _ := strconv.ParseUint(accessRequestStr, 10, 64)
	validationErr := body.Validate()
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validationErr.Error(),
		})
		return
	}

	// permission check
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, canvasBranchID, permissiongroup.CANVAS_BRANCH_MANAGE_PERMS); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	err := App.Service.ManageAccessRequest(accessRequestID, body, user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Access Request is Updated.",
	})
}
