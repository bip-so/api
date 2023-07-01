package role

import (
	"net/http"
	"strconv"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"

	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
)

// Role CRUD
var (
	RoleCrudRouteHandler roleRoutes
)

// @Summary 	Create role for a studio
// @Tags		Roles
// @Router 		/v1/role/create [post]
// @Param 		body body 		CreateRolePost true "Create Role"
func (r roleRoutes) createRole(c *gin.Context) {
	// Check if Studio Header Key is present
	studioId, status := middlewares.StudioHeaderRequiredCheck(c)
	if !status {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": "Studio Header Missing",
		})
		return
	}
	// Get the Post Data in Body !
	var body CreateRolePost
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// permission check
	if body.Name == models.SYSTEM_ADMIN_ROLE || body.Name == models.SYSTEM_ROLE_MEMBER {
		response.RenderPermissionError(c)
		return
	} else {
		studioID, _ := r.GetStudioId(c)
		authUserID := r.GetLoggedInUserId(c)
		if hasPermission, err := permissions.App.Service.CanUserDoThisOnStudio(authUserID, studioID, permissiongroup.STUDIO_CREATE_DELETE_ROLE); err != nil || !hasPermission {
			response.RenderPermissionError(c)
			return
		}
	}

	// Call Create Role Service
	role, err := App.Service.CreateNewRole(studioId, body)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	serializer := RoleGenericSerializer{c}
	response.RenderOkWithData(c, serializer.GetStudioSerializer(role))
	return
}

// @Summary 	Edit role for a studio
// @Tags		Roles
// @Router 		/v1/role/edit [post]
// @Param 		body body 		UpdateRolePost true "Edit Role"
func (r roleRoutes) editRole(c *gin.Context) {

	// Get the Post Data in Body !
	var body UpdateRolePost
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// permission check
	if body.Name == models.SYSTEM_ADMIN_ROLE || body.Name == models.SYSTEM_ROLE_MEMBER {
		response.RenderPermissionError(c)
		return
	} else {
		studioID, _ := r.GetStudioId(c)
		authUserID := r.GetLoggedInUserId(c)
		if hasPermission, err := permissions.App.Service.CanUserDoThisOnStudio(authUserID, studioID, permissiongroup.STUDIO_CREATE_DELETE_ROLE); err != nil || !hasPermission {
			response.RenderPermissionError(c)
			return
		}
	}

	// Call Create Role Service
	err := App.Service.UpdateRole(body)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	response.RenderSuccessResponse(c, "successfully edited role")
	return
}

// @Summary 	Delete role for a studio
// @Tags		Roles
// @Router 		/v1/role/:roleId [delete]
func (r roleRoutes) deleteRole(c *gin.Context) {
	// Check if Studio Header Key is present - 	//RoleID
	roleId, _ := c.Params.Get("roleId")
	roleIdInt, err := strconv.ParseUint(roleId, 10, 64)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	_, status := middlewares.StudioHeaderRequiredCheck(c)
	if !status {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": "Studio Header Missing",
		})
		return
	}

	role, _ := App.Repo.GetRole(roleIdInt)
	// permission check
	if role.Name == models.SYSTEM_ADMIN_ROLE || role.Name == models.SYSTEM_ROLE_MEMBER {
		response.RenderPermissionError(c)
		return
	} else {
		studioID, _ := r.GetStudioId(c)
		authUserID := r.GetLoggedInUserId(c)
		if hasPermission, err := permissions.App.Service.CanUserDoThisOnStudio(authUserID, studioID, permissiongroup.STUDIO_CREATE_DELETE_ROLE); err != nil || !hasPermission {
			response.RenderPermissionError(c)
			return
		}
	}
	members := role.Members
	err = App.Service.DeleteRole(roleIdInt)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	go func() {
		for _, member := range members {
			permissions.App.Service.InvalidateUserPermissionCache(member.UserID, member.StudioID)
		}
	}()
	response.RenderSuccessResponse(c, "successfully deleted role")
	return
}

// @Summary 	Update Role Membership
// @Tags		Roles
// @Router 		/v1/role/membership [post]
// @Param 		body body 		UpdateManagementPost true "Update Members In Role"
func (r roleRoutes) updateMembership(c *gin.Context) {
	var body UpdateManagementPost
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// permission check
	studioID, _ := r.GetStudioId(c)
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnStudio(authUserID, studioID, permissiongroup.STUDIO_ADD_REMOVE_USER_TO_ROLE); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	addMemberIDs, err := App.Service.UpdateMembershipRole(body)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	role, err := App.Repo.GetRoleByID(body.RoleId)

	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	members, err := App.Repo.GetMembersByIDs(addMemberIDs)
	role.Members = members

	response.RenderOkWithData(c, SerializeRole(role))

}

// @Summary 	Get Member roles
// @Tags		Roles
// @Security 	bearerAuth
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Param       memberId  path    string  true  "Member Id"
// @Router 		/v1/role/member/{memberId} [get]
func (r roleRoutes) getMemberRoles(c *gin.Context) {
	memberID := c.Param("memberId")
	studioID, _ := r.GetStudioId(c)
	memberRoles, err := App.Controller.GetMemberRoles(studioID, utils.Uint64(memberID))
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	response.RenderResponse(c, memberRoles)
	return
}
