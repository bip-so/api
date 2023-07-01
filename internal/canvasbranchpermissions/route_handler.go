package canvasbranchpermissions

import (
	"errors"
	"fmt"
	"strconv"

	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"

	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/apiutil"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
)

var (
	CanvasBranchPermissionsRouteHandler canvasBranchPermissionsRoutes
)

// getCanvasBranchPermissionsRouteHandler
// @Summary 	Get canvas branch permissions
// @Description
// @Tags		CanvasBranch Permissions
// @Security 	bearerAuth
// @Param 		canvasBranchId 	path 	string	true "Canvas Branch Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/canvasbranchpermission/{canvasBranchId} [get]
func (r *canvasBranchPermissionsRoutes) getCanvasBranchPermissionsRouteHandler(c *gin.Context) {

	canvasBranchIDStr := c.Param("canvasBranchId")
	canvasBranchID, err := strconv.ParseUint(canvasBranchIDStr, 10, 64)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, canvasBranchID, permissiongroup.CANVAS_BRANCH_VIEW); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	// get id
	canvasBranchPerm, err := App.Controller.getCanvasBranchPermissionController(map[string]interface{}{"canvas_branch_id": canvasBranchID})

	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}

	var result []CanvasBranchPermissionSerializer
	for _, perm := range canvasBranchPerm {
		result = append(result, *SerializeCanvasBranchPermissionsPermission(&perm))
	}
	response.RenderResponse(c, result)
}

// createCanvasBranchPermissionsRouteHandler New canvas branch permission
// @Summary 	Create New canvas branch permission
// @Description
// @Tags		CanvasBranch Permissions
// @Security 	bearerAuth
// @Param 		body  		body 		NewCanvasBranchPermissionCreatePost true "Create Canvas Branch permission"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/canvasbranchpermission/update [post]
func (r *canvasBranchPermissionsRoutes) createCanvasBranchPermissionsRouteHandler(c *gin.Context) {
	inherit := c.Query("inherit")
	user, _ := r.GetLoggedInUser(c)
	rawStudioId, _ := c.Get("currentStudio")
	studioID := rawStudioId.(uint64)
	var body NewCanvasBranchPermissionCreatePost
	if err := apiutil.Bind(c, &body); err != nil {
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}

	// permission check
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, body.CanvasBranchId, permissiongroup.CANVAS_BRANCH_MANAGE_PERMS); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	canvasBranchPerm, err := App.Controller.CreateCanvasBranchPermissionController(body, studioID, user.ID, inherit)

	if err != nil && err.Error() == "cannot update canvas creator" {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}

	response.RenderResponse(c, SerializeCanvasBranchPermissionsPermission(canvasBranchPerm))

}

// deleteCanvasBranchPermissionsRouteHandler Delete canvas Branch
// @Summary 	Hard delete canvas branch permission
// @Description
// @Tags		CanvasBranch Permissions
// @Security 	bearerAuth
// @Param 		canvasBranchPermissionId 	path 	string	true "Canvas Branch ID"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/canvasbranchpermission/{canvasBranchPermissionId} [delete]
func (r *canvasBranchPermissionsRoutes) deleteCanvasBranchPermissionsRouteHandler(c *gin.Context) {

	// get id
	canvasBranchIDStr := c.Param("canvasBranchPermissionId")
	canvasBranchPermissionID, err := strconv.ParseUint(canvasBranchIDStr, 10, 64)

	// Permission check
	canvasBranchPermission, _ := App.Repo.Get(map[string]interface{}{"id": canvasBranchPermissionID})
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, *canvasBranchPermission.CanvasBranchID, permissiongroup.CANVAS_BRANCH_MANAGE_PERMS); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	if (canvasBranchPermission.Member != nil && canvasBranchPermission.Member.UserID == canvasBranchPermission.CanvasRepository.CreatedByID) ||
		(canvasBranchPermission.Role != nil && canvasBranchPermission.Role.Name == models.SYSTEM_ADMIN_ROLE) {
		response.RenderCustomErrorResponse(c, errors.New("Cannot delete the canvas creator or bip admin"))
		return
	}
	var col models.CanvasBranchPermission
	err = App.Repo.Manager.HardDeleteByID(col.TableName(), canvasBranchPermissionID)

	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}

	// clearing redis cache
	// @todo Later move this to kafka
	go func() {
		inv := &permissions.InvalidatePermissions{
			MemberID:     canvasBranchPermission.MemberId,
			RoleID:       canvasBranchPermission.RoleId,
			CollectionID: canvasBranchPermission.CollectionId,
		}
		if canvasBranchPermission.CbpParentCanvasRepositoryID != nil && *canvasBranchPermission.CbpParentCanvasRepositoryID != 0 {
			inv.InvalidationOn = "subCanvas"
			inv.ParentCanvasID = canvasBranchPermission.CbpParentCanvasRepositoryID
		} else {
			inv.InvalidationOn = "canvas"
		}
		err = permissions.App.Service.InvalidatePermissions(inv)
		if err != nil {
			fmt.Println(err)
		}

		if canvasBranchPermission.MemberId != nil {
			permissions.App.Service.RemoveMemberViewMetadataPermissionOnParents(*canvasBranchPermission.CanvasBranchID, *canvasBranchPermission.MemberId, canvasBranchPermission.Member.UserID, canvasBranchPermission.StudioID)
		} else {
			permissions.App.Service.RemoveRoleViewMetadataPermissionOnParents(*canvasBranchPermission.CanvasBranchID, *canvasBranchPermission.RoleId, canvasBranchPermission.StudioID)
		}
	}()

	response.RenderSuccessResponse(c, "deleted successfully")

}

// BulkCreateCanvasBranchPermissionsRouteHandler New canvas branch permission
// @Summary 	Create Bulk New canvas branch permission
// @Description
// @Tags		CanvasBranch Permissions
// @Security 	bearerAuth
// @Param 		body  		body 		[]NewCanvasBranchPermissionCreatePost true "Create bulk Canvas Branch permission"
// @Param 		inherit 	query 	string	false "inherit permissions"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/canvasbranchpermission/bulk-update [post]
func (r *canvasBranchPermissionsRoutes) BulkCreateCanvasBranchPermissionsRouteHandler(c *gin.Context) {

	user, _ := r.GetLoggedInUser(c)
	rawStudioId, _ := c.Get("currentStudio")
	studioID := rawStudioId.(uint64)
	inheritPerms := c.Query("inherit")
	var body []NewCanvasBranchPermissionCreatePost
	if err := apiutil.Bind(c, &body); err != nil {
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}

	// permission check
	authUserID := r.GetLoggedInUserId(c)
	if len(body) > 0 {
		if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, body[0].CanvasBranchId, permissiongroup.CANVAS_BRANCH_MANAGE_PERMS); err != nil || !hasPermission {
			response.RenderPermissionError(c)
			return
		}
	}

	canvasBranchPerm, err := App.Controller.BulkCreateCanvasBranchPermissionController(body, studioID, user.ID, inheritPerms)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}

	response.RenderResponse(c, canvasBranchPerm)
}

// inheritParentPermissionsRouteHandler
// @Summary 	Inherits the parent permissions.
// @Description
// @Tags		CanvasBranch Permissions
// @Security 	bearerAuth
// @Param 		canvasBranchId 	path 	string	true "Canvas Branch Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/canvasbranchpermission/inherit/{canvasBranchId} [post]
func (r *canvasBranchPermissionsRoutes) inheritParentPermissionsRouteHandler(c *gin.Context) {

	canvasBranchIDStr := c.Param("canvasBranchId")
	canvasBranchID, err := strconv.ParseUint(canvasBranchIDStr, 10, 64)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	// Permission check
	//authUserID := r.GetLoggedInUserId(c)
	//if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, canvasBranchID, permissiongroup.CANVAS_BRANCH_MANAGE_PERMS); err != nil || !hasPermission {
	//	response.RenderPermissionError(c)
	//	return
	//}
	// get id
	canvasBranchPerm, err := App.Controller.inheritParentPermissionController(canvasBranchID)

	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}

	var result []CanvasBranchPermissionSerializer
	for _, perm := range canvasBranchPerm {
		result = append(result, *SerializeCanvasBranchPermissionsPermission(&perm))
	}
	response.RenderResponse(c, result)
}
