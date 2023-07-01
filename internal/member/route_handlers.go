package member

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"strconv"
)

// LeaveStudio We'll get the Studio Object in the Context
// @Summary 	Leave Studio
// @Tags		Member
// @Accept 		json
// @Produce 	json
// @Param 		studioId 	path 	string	true "Studio Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/member/leave-studio/{studioId} [post]
func (r memberRoutes) LeaveStudio(c *gin.Context) {
	ctxUser, _ := c.Get("currentUser")
	if ctxUser == nil {
		response.RenderPermissionError(c)
		return
	}
	loggedInUser := ctxUser.(*models.User)
	rawStudioId := c.Param("studioId")
	studioIDInt := utils.Uint64(rawStudioId)
	err := App.Service.LeaveStudio([]uint64{loggedInUser.ID}, studioIDInt)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	response.RenderSuccessResponse(c, "successfully left the studio")

}

// GetCanvasBranchMembers All members or roles having access to canvas branch
// @Summary 	All members or roles having access to canvas branch
// @description we can get all members or roles having access to canvas branch
// @Tags		Member
// @Accept       json
// @Produce      json
// @Param 		canvasBranchID 	path 	string	true "Canvas Branch Id"
// @Router 		/v1/member/canvas-branch/{canvasBranchID} [get]
func (r *memberRoutes) GetCanvasBranchMembers(c *gin.Context) {
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID := utils.Uint64(canvasBranchIDStr)
	inviteCode := c.Query("inviteCode")

	// Permission check
	authUserID := r.GetLoggedInUserId(c)
	if authUserID == 0 && inviteCode != "" {
		branch, _ := permissions.App.Repo.GetBranchWithRepo(map[string]interface{}{"id": canvasBranchID})
		_, err := permissions.App.Service.CheckBranchAccessToken(inviteCode, branch.CanvasRepository.Key)
		if err != nil {
			response.RenderPermissionError(c)
			return
		}
		canvasBranchMembers, err := App.Controller.CanvasBranchMembersAndRoles(canvasBranchID)
		if err != nil {
			response.RenderErrorResponse(c, err)
			return
		}
		response.RenderResponse(c, canvasBranchMembers)
		return
	}
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, canvasBranchID, permissiongroup.CANVAS_BRANCH_VIEW); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	canvasBranchMembers, err := App.Controller.CanvasBranchMembersAndRoles(canvasBranchID)
	if err != nil {
		response.RenderErrorResponse(c, err)
		return
	}
	response.RenderResponse(c, canvasBranchMembers)
}

// GetRoleMembers All role members
// @Summary 	All role members
// @description we can get all members of a role
// @Tags		Member
// @Accept       json
// @Produce      json
// @Param 		roleID 	path 	string	true "roleID Id"
// @Router 		/v1/member/role/{roleID} [get]
func (r *memberRoutes) GetRoleMembers(c *gin.Context) {
	roleIDStr := c.Param("roleID")
	roleID := utils.Uint64(roleIDStr)

	skip := c.Query("skip")
	skipInt, _ := strconv.Atoi(skip)
	limit := c.Query("limit")
	limitInt, _ := strconv.Atoi(limit)
	if limitInt == 0 {
		limitInt = configs.PAGINATION_LIMIT
	}

	roleMembers, err := App.Controller.RoleMembers(roleID, skipInt, limitInt)
	if err != nil {
		response.RenderErrorResponse(c, err)
		return
	}
	next := "-1"
	if len(roleMembers) == limitInt {
		next = strconv.Itoa(skipInt + len(roleMembers))
	}
	response.RenderPaginatedResponse(c, roleMembers, next)
}

// SearchMembers Search studio members
// @Summary 	All studio members
// @description we get all members of a studio
// @Tags		Member
// @Accept       json
// @Produce      json
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Param 		search 	query 	string	true "search string"
// @Router 		/v1/member/search [get]
func (r *memberRoutes) SearchMembers(c *gin.Context) {
	search := c.Query("search")
	studioID, _ := r.GetStudioId(c)
	members, err := queries.App.MemberQuery.MembersSearch(search, studioID)
	if err != nil {
		response.RenderErrorResponse(c, err)
		return
	}
	response.RenderResponse(c, BulkSerializeMembers(members))
}

// GetRoleMembersSearch All role members by search
// @Summary 	All role members by search
// @description we can search for all members of a role
// @Tags		Member
// @Accept       json
// @Produce      json
// @Param 		roleID 	path 	string	true "roleID Id"
// @Param 		search 	query 	string	true "Search"
// @Router 		/v1/member/role/{roleID}/search-members [get]
func (r *memberRoutes) GetRoleMembersSearch(c *gin.Context) {
	roleIDStr := c.Param("roleID")
	roleID := utils.Uint64(roleIDStr)
	search := c.Query("search")

	skip := c.Query("skip")
	skipInt, _ := strconv.Atoi(skip)
	limit := c.Query("limit")
	limitInt, _ := strconv.Atoi(limit)
	if limitInt == 0 {
		limitInt = configs.PAGINATION_LIMIT
	}

	roleMembers, err := App.Controller.RoleMembersSearch(search, roleID, skipInt, limitInt)
	if err != nil {
		response.RenderErrorResponse(c, err)
		return
	}
	next := "-1"
	if len(roleMembers) == limitInt {
		next = strconv.Itoa(skipInt + len(roleMembers))
	}
	response.RenderPaginatedResponse(c, roleMembers, next)
}
