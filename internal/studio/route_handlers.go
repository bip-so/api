package studio

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/payments"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/internal/shared"
	"net/http"
	"strconv"

	"gitlab.com/phonepost/bip-be-platform/internal/feed"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/supabase"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"

	"gitlab.com/phonepost/bip-be-platform/internal/member"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/role"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/apiutil"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
)

// Create New Studio
// @Summary 	Creates a new studio for the user
// @Description
// @Tags		Studio
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		body 		body 		CreateStudioValidator true "Create Studio Data"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studio/create [post]
func createStudioRouteHandler(c *gin.Context) {
	authUsr, _ := c.Get("currentUser")
	authUser := authUsr.(*models.User)
	var body CreateStudioValidator

	if err := apiutil.Bind(c, &body); err != nil {
		fmt.Println(err, body)
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}

	studio, err := App.Controller.CreateStudioController(&body, authUser)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}

	response.RenderResponse(c, SerializeStudio(studio))
}

// Get Studio
// @Summary 	Gets a studio by id
// @Description
// @Tags		Studio
// @Accept 		json
// @Produce 	json
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studio/:studioid [get]
func (r *studioRoutes) getStudioRouteHandler(c *gin.Context) {
	authUsr, _ := c.Get("currentUser")
	var authUser *models.User
	if authUsr != nil {
		authUser = authUsr.(*models.User)
	}

	studioIDStr := c.Param("studioId")
	studioID, err := strconv.ParseUint(studioIDStr, 10, 64)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	studio, member, err := App.Controller.GetStudioController(studioID, authUser)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.RenderErrorResponse(c, response.NotFoundError(c.Request.Context()))
			return
		}
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	if authUsr != nil {
		members := []models.Member{}
		if member != nil {
			members = []models.Member{*member}
		}
		response.RenderResponse(c, SerializeStudioForUser(studio, authUser, &members))
		return
	}
	response.RenderResponse(c, SerializeStudio(studio))
}

// Edit a Studio
// @Summary 	Edit a studio
// @Description
// @Tags		Studio
// @Accept 		json
// @Produce 	json
// @Param 		body 		body 		UpdateStudioValidator true "Update Studio Data"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studio/edit [post]
func (r *studioRoutes) editRouteHandler(c *gin.Context) {
	authUser, _ := r.GetLoggedInUser(c)
	rawStudioId, _ := c.Get("currentStudio")
	studioID := rawStudioId.(uint64)
	var body UpdateStudioValidator
	if err := apiutil.Bind(c, &body); err != nil {
		fmt.Println(err, body)
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}

	// permission check
	//authUserID := r.GetLoggedInUserId(c)
	//if hasPermission, err := permissions.App.Service.CanUserDoThisOnStudio(authUserID, studioID, permissiongroup.STUDIO_EDIT_STUDIO_PROFILE); err != nil || !hasPermission {
	//	response.RenderPermissionError(c)
	//	return
	//}

	studio, err := App.Controller.editStudioControler(studioID, &body)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.NotFoundError(c)
			return
		}
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	serialized := SerializeStudio(studio)
	if authUser == nil {
		serialized.Permission = models.PGStudioNoneSysName
	} else {
		permissionList, _ := permissions.App.Service.CalculateStudioPermissions(authUser.ID)
		serialized.Permission = permissionList[studio.ID]
	}
	response.RenderResponse(c, serialized)
}

// Update Studio Image
// @Summary 	Upload & Update a studio image
// @Description
// @Tags		Studio
// @Accept 		json
// @Produce 	json
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Param 		file		formData file false "File"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studio/image [post]
func (r *studioRoutes) imageRouteHandler(c *gin.Context) {
	studioId, exists := c.Get("currentStudio")
	fmt.Println(studioId, exists)
	if !exists {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": "studio id not found",
		})
		return
	}
	// permission check
	studioID, _ := r.GetStudioId(c)
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnStudio(authUserID, studioID, permissiongroup.STUDIO_EDIT_STUDIO_PROFILE); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": "file not found",
		})
		return
	}
	object, _ := file.Open()
	studio, err := App.Controller.UpdateStudioImage(studioId.(uint64), object, file.Filename)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	response.RenderResponse(c, SerializeStudio(studio))
}

// Get Popular Studios
// @Summary 	Gets popular studios
// @Description
// @Tags		Studio
// @Accept 		json
// @Produce 	json
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studio/popular [get]
func (r *studioRoutes) popularStudiosRouteHandler(c *gin.Context) {
	authUsr, _ := c.Get("currentUser")
	var authUser *models.User
	if authUsr != nil {
		authUser = authUsr.(*models.User)
	}
	skip := c.Query("skip")
	skipInt, _ := strconv.Atoi(skip)
	limit := c.Query("limit")
	limitInt, _ := strconv.Atoi(limit)
	if limitInt == 0 {
		limitInt = configs.PAGINATION_LIMIT
	}
	studios, members, err := App.Controller.GetPopularStudioController(authUser, skipInt, limitInt)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}

	var responseStudios []*StudioSerializer
	for _, studio := range studios {
		membersCount := App.StudioRepo.studioMembersCount(studio.ID)
		if authUser != nil {
			studioData := SerializeStudioForUser(&studio, authUser, members)
			studioData.MembersCount = membersCount
			responseStudios = append(responseStudios, studioData)
		} else {
			studioData := SerializeStudio(&studio)
			studioData.MembersCount = membersCount
			responseStudios = append(responseStudios, studioData)
		}
	}
	next := "-1"
	if len(responseStudios) == configs.PAGINATION_LIMIT {
		next = strconv.Itoa(skipInt + len(responseStudios))
	}
	response.RenderPaginatedResponse(c, responseStudios, next)
}

// @Summary 	Toggle studio Membership
// @Tags		Studio
// @Accept 		json
// @Produce 	json
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studio/toggle-membership [get]
func (r *studioRoutes) toggleStudioMembershipRouteHandler(c *gin.Context) {
	authUsr, _ := c.Get("currentUser")
	var authUser *models.User
	if authUsr != nil {
		authUser = authUsr.(*models.User)
	}
	rawStudioId, _ := c.Get("currentStudio")
	studioID := rawStudioId.(uint64)

	App.Controller.ToggleStudioMembershipController(studioID, authUser.ID)
	response.RenderSuccessResponse(c, "successfully toggled")
}

func (r *studioRoutes) studioListRouteHandler(c *gin.Context) {

}

// Delete Studio
// @Summary 	Deletes a studio by id
// @Description
// @Tags		Studio
// @Accept 		json
// @Produce 	json
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studio/:studioid [delete]
func (r *studioRoutes) deleteRouteHandler(c *gin.Context) {
	studioIDStr := c.Param("studioId")
	authUsr, _ := c.Get("currentUser")
	authUser := authUsr.(*models.User)
	studioID, err := strconv.ParseUint(studioIDStr, 10, 64)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}

	// permission check
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnStudio(authUserID, studioID, permissiongroup.STUDIO_DELETE); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	err = App.Controller.deleteStudioController(studioID, authUser.ID)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	response.RenderSuccessResponse(c, "successfully deleted studio")
}

// We'll get the Studio Object in the Context
// We'll get user as Param (:userId) from URL
// func addMemberToStudio(c *gin.Context) {
// userIDStr := c.Param("userId")
// userID, _ := strconv.ParseUint(userIDStr, 10, 64)
// rawStudioId, _ := c.Get("currentStudio")
// studioIDInt := rawStudioId.(uint64)
// err := member.MemberService.AddMemberToStudio(member.NewDiscordMember{UserId: userID, StudioId: studioIDInt})
// fmt.Println(err)
// }

// Get Studio Roles
// @Summary 	Gets roles by studio Id
// @Description
// @Tags		Studio
// @Accept 		json
// @Produce 	json
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studio/roles [get]
func (r *studioRoutes) getStudioRolesRouteHandler(c *gin.Context) {
	rawStudioId, _ := c.Get("currentStudio")
	studioID := rawStudioId.(uint64)

	roles, err := role.App.Service.GetRolesByStudio(studioID)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}

	var result []role.RoleSerializer
	for _, rol := range roles {
		result = append(result, *role.SerializeRole(&rol))
	}
	response.RenderResponse(c, result)
}

// Get Studio Members
// @Summary 	Gets members by studioId
// @Description
// @Tags		Studio
// @Accept 		json
// @Produce 	json
// @Param 		skip 	query 		string		 		false "next page"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studio/members [get]
func (r *studioRoutes) getStudioMembersRouteHandler(c *gin.Context) {
	authUserID := r.GetLoggedInUserId(c)
	rawStudioId, _ := c.Get("currentStudio")
	studioID := rawStudioId.(uint64)

	skip := c.Query("skip")
	skipInt, _ := strconv.Atoi(skip)
	members, err := member.App.Service.GetMembersByStudio(studioID, skipInt)

	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	var serializedMembers []member.MemberSerializer
	for _, memb := range members {
		memberData := member.SerializeMember(&memb)
		isFollowing := false
		if authUserID != memb.UserID {
			follower, _ := App.StudioRepo.GetUserFollows(memb.UserID, authUserID)
			if follower.ID != 0 {
				isFollowing = true
			}
		}
		memberData.User.IsFollowing = isFollowing
		serializedMembers = append(serializedMembers, *memberData)
	}
	next := "-1"
	if len(members) == configs.PAGINATION_LIMIT {
		next = strconv.Itoa(skipInt + len(members))
	}
	response.RenderPaginatedResponse(c, serializedMembers, next)

}

// Ban User
// @Summary
// @Description
// @Tags		Studio
// @Accept 		json
// @Produce 	json
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Param 		body body 		BanUserValidator true "Ban User"
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Router 		/v1/studio/ban [post]
func (r *studioRoutes) banUser(c *gin.Context) {

	rawStudioId, _ := c.Get("currentStudio")
	studioIDInt := rawStudioId.(uint64)

	removedByUsr, _ := c.Get("currentUser")
	removedByUser := removedByUsr.(*models.User)

	var body BanUserValidator
	if err := apiutil.Bind(c, &body); err != nil {
		fmt.Println(err, body)
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}

	// permission check
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnStudio(authUserID, studioIDInt, permissiongroup.STUDIO_CREATE_DELETE_ROLE); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	err := member.App.Service.BanUser(body.UserID, studioIDInt, body.BanReason, removedByUser.ID)
	fmt.Println(err)
	response.RenderSuccessResponse(c, "User banned successfully")
	// @todo Later move to kafka on member join we need to invalidate the user associated studios and send event to supabase
	go func() {
		queries.App.StudioQueries.DeleteUserAssociatedStudioDataByUserID(body.UserID)
		supabase.UpdateUserSupabase(body.UserID, true)

		// unfollow in stream
		feed.App.Service.LeaveStudio(studioIDInt, body.UserID)
	}()
}

// Get List Request to Join Studio List
// @Tags		Studio
// @Accept 		json
// @Produce 	json
// @Param       studioId  path    string  true  "Studio Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studio/{studioId}/membership-request/list [get]
func (r *studioRoutes) requestToJoinStudioList(c *gin.Context) {
	rawStudioId := c.Param("studioId")
	studioIDInt := utils.Uint64(rawStudioId)
	//authUsr, _ := c.Get("currentUser")
	//authUser := authUsr.(*models.User)
	requestsList, err := App.Controller.GetRequestToJoinStudioListController(studioIDInt)
	requestsSerialized := ManySerializeStudioMembersRequest(requestsList)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	c.JSON(http.StatusOK, shared.GetSimpleListData(requestsSerialized))
}

// Request membership to Join Studio
// @Tags		Studio
// @Accept 		json
// @Produce 	json
// @Param       studioId  path    string  true  "Studio Id"
// @Param 		body 		body 		StudioMembershipRequestNew true "New Request Object"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studio/{studioId}/membership-request/new [post]
func (r *studioRoutes) requestToJoinStudioCreate(c *gin.Context) {
	rawStudioId := c.Param("studioId")
	studioIDInt := utils.Uint64(rawStudioId)
	authUsr, _ := c.Get("currentUser")
	authUser := authUsr.(*models.User)
	var body StudioMembershipRequestNew
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err := App.Controller.CreateRequestToJoinStudioController(studioIDInt, authUser.ID)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	c.JSON(http.StatusOK, shared.GetStringMessageData("Request has been added."))
	go func() {
		notifications.App.Service.PublishNewNotification(notifications.CreateRequestToJoinStudio, authUser.ID, []uint64{}, &studioIDInt,
			nil, notifications.NotificationExtraData{Status: "Pending"}, nil, nil)
	}()
}

// Request Reject To Join a Studio
// @Tags		Studio
// @Accept 		json
// @Produce 	json
// @Param       studioId  path    string  true  "Studio Id"
// @Param 		body 		body 		StudioMembershipRequestReject true "New Request Object"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studio/{studioId}/membership-request/:membershipRequestID/reject [post]
func (r *studioRoutes) requestToJoinStudioReject(c *gin.Context) {
	membershipRequestIDStr := c.Param("membershipRequestID")
	membershipRequestID, _ := strconv.ParseUint(membershipRequestIDStr, 10, 64)
	rawStudioId := c.Param("studioId")
	studioIDInt := utils.Uint64(rawStudioId)
	authUsr, _ := c.Get("currentUser")
	authUser := authUsr.(*models.User)

	var body StudioMembershipRequestReject
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	membershipRequest, err := App.Controller.RejectRequestToJoinStudioController(membershipRequestID, authUser.ID)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	c.JSON(http.StatusOK, shared.GetStringMessageData("Request has been reject."))
	go func() {
		notifications.App.Service.PublishNewNotification(notifications.RejectRequestToJoinStudio, authUser.ID, []uint64{membershipRequest.UserID}, &studioIDInt,
			nil, notifications.NotificationExtraData{Status: "Rejected"}, nil, nil)
	}()
}

// Request Accept To Join a Studio
// @Tags		Studio
// @Accept 		json
// @Produce 	json
// @Param       studioId  path    string  true  "Studio Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studio/{studioId}/membership-request/:membershipRequestID/accept [post]
func (r *studioRoutes) acceptToJoinStudioReject(c *gin.Context) {
	membershipRequestIDStr := c.Param("membershipRequestID")
	membershipRequestID, _ := strconv.ParseUint(membershipRequestIDStr, 10, 64)
	rawStudioId := c.Param("studioId")
	studioIDInt := utils.Uint64(rawStudioId)
	authUsr, _ := c.Get("currentUser")
	authUser := authUsr.(*models.User)

	membershipRequest, err := App.Controller.AcceptRequestToJoinStudioController(membershipRequestID, authUser.ID)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	//membershipRequestUser, _ := user2.App.Repo.GetUser(map[string]interface{}{"id": membershipRequest.UserID})
	membershipRequestUser, _ := queries.App.UserQueries.GetUserByID(membershipRequest.UserID)

	err = App.Controller.JoinStudioController(membershipRequestUser, studioIDInt)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	mem, _ := queries.App.MemberQuery.GetMemberInstanceWithdUserInstance(map[string]interface{}{"user_id": membershipRequestUser.ID, "studio_id": studioIDInt})
	resp := member.SerializeMember(mem)
	go func() {
		notifications.App.Service.PublishNewNotification(notifications.AcceptRequestToJoinStudio, authUser.ID, []uint64{membershipRequest.UserID}, &studioIDInt,
			nil, notifications.NotificationExtraData{Status: "Accepted"}, nil, nil)
	}()
	response.RenderResponse(c, resp)
}

// Request to Join Studio
// @Summary 	adding user to a studio
// @Description If user had already left, he can rejoin, else if he was banned then it throws error, otherwise it will join
// @Tags		Studio
// @Accept 		json
// @Produce 	json
// @Param       studioId  path    string  true  "Studio Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studio/{studioId}/join [post]
func (r *studioRoutes) joinStudio(c *gin.Context) {
	rawStudioId := c.Param("studioId")
	studioIDInt := utils.Uint64(rawStudioId)
	authUsr, _ := c.Get("currentUser")
	authUser := authUsr.(*models.User)
	err := App.Controller.JoinStudioController(authUser, studioIDInt)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}

	go func() {
		notifications.App.Service.PublishNewNotification(notifications.JoinedStudio, authUser.ID, []uint64{authUser.ID}, &studioIDInt,
			nil, notifications.NotificationExtraData{}, nil, nil)
	}()
	response.RenderSuccessResponse(c, "successfully joined the studio")
}

// Get Studio Stats
// @Summary 	Get Studio Stats
// @Tags		Studio
// @Accept 		json
// @Produce 	json
// @Param       studioId  path    string  true  "Studio Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studio/{studioId}/stats [get]
func (r *studioRoutes) getStudioStats(c *gin.Context) {
	rawStudioId := c.Param("studioId")
	studioIDInt := utils.Uint64(rawStudioId)
	//authUsr, _ := c.Get("currentUser")
	//authUser := authUsr.(*models.User)
	data := App.Controller.StudioStats(studioIDInt)
	c.JSON(http.StatusOK, data)

}

// MemberCount
// @Summary
// @Description counts all members which belong to the studio excluding those who have left or banned
// @Tags		Studio
// @Accept 		json
// @Produce 	json
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studio/memberCount [post]
func (r *studioRoutes) memberCount(c *gin.Context) {
	rawStudioId, _ := c.Get("currentStudio")
	studioIDInt := rawStudioId.(uint64)
	count, err := App.Controller.MemberCountController(studioIDInt)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	response.RenderOkWithData(c, SerializeMemberCount(count))
	return
}

// Join Bulk
// @Summary 	add users in bulk to the studio
// @Description for all users, If user had already left, he can rejoin, else if he was banned then it throws error, otherwise it will join
// @Tags		Studio
// @Accept 		json
// @Produce 	json
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Param 		body 		body 		JoinStudioBulkPost true "Join users in bulk"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studio/join/bulk [post]
func (r *studioRoutes) joinStudioInBulk(c *gin.Context) {
	rawStudioId, _ := c.Get("currentStudio")
	studioIDInt := rawStudioId.(uint64)

	var body JoinStudioBulkPost
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

	members, err := App.Controller.JoinStudioInBulkController(body, studioIDInt, authUserID)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	var serializedMembers []member.MemberSerializer
	for _, memb := range members {
		serializedMembers = append(serializedMembers, *member.SerializeMember(&memb))
	}

	response.RenderResponse(c, serializedMembers)
}

// @Summary 	Get Studio Admins
// @Tags		Studio
// @Security 	bearerAuth
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Router 		/v1/studio/admins [get]
func (r studioRoutes) getStudioAdminMembers(c *gin.Context) {
	studioID, _ := r.GetStudioId(c)
	adminMembers, err := App.Controller.StudioAdminMembers(studioID)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	response.RenderResponse(c, adminMembers)
	return
}

// Invite Via Email Flow
// @Summary 	Invite many users via emails flow
// @Tags		Studio
// @Accept 		json
// @Produce 	json
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Param 		body 		body 		[]NewInvitePostOne true "Invite Users in Bulk"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studio/invite-via-email [post]
func (r studioRoutes) InviteWithEmailFlow(c *gin.Context) {
	rawStudioId, _ := c.Get("currentStudio")
	studioIDInt := rawStudioId.(uint64)
	var body []NewInvitePostOne
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	authUsr, _ := c.Get("currentUser")
	var authUser *models.User
	// we can also throw an error here
	if authUsr != nil {
		authUser = authUsr.(*models.User)
	}
	data := App.Controller.CreateStudioInvites(body, studioIDInt, authUser)

	c.JSON(http.StatusOK, gin.H{
		"message": "Invites are send",
		"context": data,
	})

}

// Get Payment Portal Link
// @Summary 	Get Payment Portal Link
// @Tags		Studio
// @Accept 		json
// @Produce 	json
// @Param       studioId  path    string  true  "Studio Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studio/{studioId}/customer-portal-session [get]
func (r *studioRoutes) NewCustomerPortalSession(c *gin.Context) {
	rawStudioId := c.Param("studioId")
	studioIDInt := utils.Uint64(rawStudioId)
	url := c.Query("url")

	authUsr, _ := c.Get("currentUser")
	var authUser *models.User
	// we can also throw an error here
	if authUsr != nil {
		authUser = authUsr.(*models.User)
	}

	data := payments.App.Service.PortalLink(studioIDInt, url, authUser)
	c.JSON(http.StatusOK, data)
}

// Get Payment Link
// @Summary 	Get Payment Link
// @Tags		Studio
// @Accept 		json
// @Produce 	json
// @Param       studioId  path    string  true  "Studio Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studio/{studioId}/customer-portal-session [get]
func (r *studioRoutes) NewPaymentLinkSession(c *gin.Context) {
	rawStudioId := c.Param("studioId")
	studioIDInt := utils.Uint64(rawStudioId)
	url := c.Query("url")

	//authUsr, _ := c.Get("currentUser")
	//var authUser *models.User
	//// we can also throw an error here
	//if authUsr != nil {
	//	authUser = authUsr.(*models.User)
	//}

	data := payments.App.Service.CheckoutSession(studioIDInt, url)
	c.JSON(http.StatusOK, data)
}
