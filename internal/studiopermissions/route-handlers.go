package studiopermissions

import (
	"errors"
	"fmt"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
)

// Get Studio Permissions
// @Summary 	Gets studio permissions by studio id
// @Description
// @Tags		Studio Permission
// @Accept 		json
// @Produce 	json
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studiopermission/getAll [get]
func (r *studioPermissionsRoutes) getStudioPermissionsRouteHandler(c *gin.Context) {

	rawStudioId, _ := c.Get("currentStudio")
	studioID := rawStudioId.(uint64)

	studioPerms, err := StudioPermissionService.GetStudioPermissions(map[string]interface{}{"studio_id": studioID})

	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	var result []StudioPermissionsSerializer
	for _, perm := range studioPerms {
		result = append(result, *SerializeStudioPermission(&perm))
	}
	response.RenderResponse(c, result)

}

// Create Studio Permissions
// @Summary 	Create/Update studio permissions
// @Description
// @Tags		Studio Permission
// @Accept 		json
// @Produce 	json
// @Param 		body  	body 		CreateStudioPermissionsPost true "Update permissions Data"
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studiopermission/update [post]
func (r *studioPermissionsRoutes) updateStudioPermissionsRouteHandler(c *gin.Context) {

	rawStudioId, _ := c.Get("currentStudio")
	studioIDInt := rawStudioId.(uint64)

	var err error
	var studioPerm *models.StudioPermission
	var body CreateStudioPermissionsPost
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// permission check
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnStudio(authUserID, studioIDInt, permissiongroup.STUDIO_MANAGE_PERMS); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	if body.RoleId != 0 && body.MemberId != 0 {
		err = errors.New("provide either roleId or memberId")
	}
	if body.RoleId != 0 {
		studioPerm, err = StudioPermissionService.UpdateStudioPermissions(map[string]interface{}{"studio_id": studioIDInt, "role_id": body.RoleId}, body, studioIDInt)

	} else if body.MemberId != 0 {
		studioPerm, err = StudioPermissionService.UpdateStudioPermissions(map[string]interface{}{"studio_id": studioIDInt, "member_id": body.MemberId}, body, studioIDInt)

	} else {
		err = errors.New("roleId or memberId not provided")
	}

	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}

	// clearing redis cache
	// @todo Later move this to kafka
	go func() {
		inv := &permissions.InvalidatePermissions{
			MemberID:       &body.MemberId,
			RoleID:         &body.RoleId,
			InvalidationOn: "studio",
		}
		err = permissions.App.Service.InvalidatePermissions(inv)
		if err != nil {
			fmt.Println(err)
		}
	}()

	response.RenderResponse(c, *SerializeStudioPermission(studioPerm))

}

// Delete Studio Permissions
// @Summary 	deletes studio permissions by studio id
// @Description
// @Tags		Studio Permission
// @Accept 		json
// @Produce 	json
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studiopermission/:studioPermissionID [delete]
func (r *studioPermissionsRoutes) deleteStudioPermissionsRouteHandler(c *gin.Context) {

	studioPermIDStr := c.Param("studioPermissionID")
	studioPermID, err := strconv.ParseUint(studioPermIDStr, 10, 64)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}

	// permission check
	studioID, _ := r.GetStudioId(c)
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnStudio(authUserID, studioID, permissiongroup.STUDIO_MANAGE_PERMS); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	var col models.StudioPermission
	studioPermission, _ := App.Repo.GetStudioPermission(map[string]interface{}{"id": studioPermID})

	err = App.Repo.Manger.HardDeleteByID(col.TableName(), studioPermID)

	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}

	// clearing redis cache
	// @todo Later move this to kafka
	go func() {
		inv := &permissions.InvalidatePermissions{
			MemberID:       studioPermission.MemberId,
			RoleID:         studioPermission.RoleId,
			InvalidationOn: "studio",
		}
		err = permissions.App.Service.InvalidatePermissions(inv)
		if err != nil {
			fmt.Println(err)
		}
	}()

	response.RenderSuccessResponse(c, "deleted permission successfully.")

}
