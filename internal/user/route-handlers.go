package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/apiutil"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/context"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
	"net/http"
	"strconv"
)

// Get User details
// @Summary 	Get User
// @Description
// @Tags		User
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Security 	bearerAuth
// @Param 		user_id 	query 		string		 		false "User ID"
// @Param 		username 	query		string		 		false "Username"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/user/info [get]
func (r userRoutes) userInfoRoute(c *gin.Context) {
	query := c.Request.URL.Query()
	userID := query.Get("user_id")
	userName := query.Get("username")
	authUserID := r.GetLoggedInUserId(c)

	userData, err := App.Controller.UserInfoController(userID, userName, authUserID)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	response.RenderResponse(c, userData)
}

// Update User details
// @Summary 	Update User
// @Description
// @Tags		User
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		file		formData file false "Body with image file"
// @Param 		fullName		formData string true "Full Name"
// @Param 		username		formData string true "Username"
// @Param 		bio		formData string true "Bio"
// @Param 		twitterUrl		formData string true "Twitter Url"
// @Param 		website		formData string true "Website"
// @Param 		location		formData string true "Location"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/user/update [put]
func (r userRoutes) updateUserRoute(c *gin.Context) {
	var body *UpdateUserValidator
	if err := apiutil.Bind(c, &body); err != nil {
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}

	user, err := context.GetAuthUser(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	updateUser, err := App.Controller.UpdateUserController(body, user)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}

	response.RenderResponse(c, updateUser)
	return
}

func (r userRoutes) UserSearchRoute(c *gin.Context) {

}

// @Summary 	Get Users Followers and Following List (Send the ?userId=0 to the get)
// @Tags		User
// @Router 		/v1/user/followers-list [get]
func (r userRoutes) followerListRoute(c *gin.Context) {
	_, err := context.GetAuthUser(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		fmt.Println(err.Error())
		return
	}
	authUser, _ := r.GetLoggedInUser(c)

	// UserID will come from query Param
	userIDString := c.Query("userId")
	userID, _ := strconv.ParseUint(userIDString, 10, 64)
	//userInstanceRequested, erruserInstanceRequested := App.Repo.GetUser(map[string]interface{}{"id": userID})
	userInstanceRequested, erruserInstanceRequested := queries.App.UserQueries.GetUserByID(userID)
	if erruserInstanceRequested != nil {
		response.RenderCustomErrorResponse(c, erruserInstanceRequested)
		fmt.Println(erruserInstanceRequested.Error())
		return
	}

	finalData := App.Controller.GetUsersFollowers(userInstanceRequested, authUser)
	c.JSON(http.StatusOK, finalData)
	return
}

// @Summary 	Get Users Following (Send the ?userId=0 to the get)
// @Tags		User
// @Router 		/v1/user/following-list [get]
func (r userRoutes) followingListRoute(c *gin.Context) {
	_, err := context.GetAuthUser(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		fmt.Println(err.Error())
		return
	}
	authUser, _ := r.GetLoggedInUser(c)

	// UserID will come from query Param
	userIDString := c.Query("userId")
	userID, _ := strconv.ParseUint(userIDString, 10, 64)
	//userInstanceRequested, erruserInstanceRequested := App.Repo.GetUser(map[string]interface{}{"id": userID})
	userInstanceRequested, erruserInstanceRequested := queries.App.UserQueries.GetUserByID(userID)

	if erruserInstanceRequested != nil {
		response.RenderCustomErrorResponse(c, erruserInstanceRequested)
		fmt.Println(erruserInstanceRequested.Error())
		return
	}

	finalData := App.Controller.GetUsersFollowing(userInstanceRequested, authUser)
	c.JSON(http.StatusOK, finalData)
	return
}

//
//func (r userRoutes) updateFollowerRoute(c *gin.Context) {
//	var body *UserUpdateFollowers
//	if err := apiutil.Bind(c, &body); err != nil {
//		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
//		return
//	}
//
//	user := c.MustGet("currentUser").(*models.User)
//	success, err := follower.FollowerService.FollowUser(user.ID, body.FollowerID)
//	if err != nil {
//		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
//		return
//	}
//
//	response.RenderResponse(c, success)
//}

// Get User Settings
// @Summary 	Get User settings
// @Description
// @Tags		User Settings
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/user/settings [get]
func (r userRoutes) GetUserSettingsRoute(c *gin.Context) {
	user, err := context.GetAuthUser(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	userSettings, err := App.Controller.GetUserSettingsController(user.ID)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}

	resp := CustomUserSettingsResponse{}
	for _, userSetting := range *userSettings {
		if userSetting.Type == models.NOTIFICATION_TYPE_APP {
			resp.App = CustomUserSettingsSerializerData(userSetting)
		} else if userSetting.Type == models.NOTIFICATION_TYPE_EMAIL {
			resp.Email = CustomUserSettingsSerializerData(userSetting)
		} else if userSetting.Type == models.NOTIFICATION_TYPE_DISCORD {
			resp.Discord = CustomUserSettingsSerializerData(userSetting)
		} else if userSetting.Type == models.NOTIFICATION_TYPE_SLACK {
			resp.Slack = CustomUserSettingsSerializerData(userSetting)
		}
	}
	response.RenderResponse(c, resp)
}

// Update User Settings
// @Summary 	Update User settings
// @Description
// @Tags		User Settings
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
//@Param 		body 		body 		PatchUserSettingsValidator true "Update User Settings"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/user/settings [patch]
func (r userRoutes) UpdateUserSettingsRoute(c *gin.Context) {
	var body *PatchUserSettingsValidator
	if err := apiutil.Bind(c, &body); err != nil {
		fmt.Println(err)
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}

	user, err := context.GetAuthUser(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	userSettings, err := App.Controller.UpdateUserSettingsController(body, user)
	if err != nil {
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}

	resp := CustomUserSettingsResponse{}
	for _, userSetting := range *userSettings {
		if userSetting.Type == models.NOTIFICATION_TYPE_APP {
			resp.App = CustomUserSettingsSerializerData(userSetting)
		} else if userSetting.Type == models.NOTIFICATION_TYPE_EMAIL {
			resp.Email = CustomUserSettingsSerializerData(userSetting)
		} else if userSetting.Type == models.NOTIFICATION_TYPE_DISCORD {
			resp.Discord = CustomUserSettingsSerializerData(userSetting)
		} else if userSetting.Type == models.NOTIFICATION_TYPE_SLACK {
			resp.Slack = CustomUserSettingsSerializerData(userSetting)
		}
	}
	response.RenderResponse(c, resp)
}

// Get User Contacts
// @Summary 	Get User contacts
// @Description
// @Tags		User Contacts
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		since 	query 		string		 		true "Since"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/user/contacts [get]
func (r userRoutes) GetUserContactsRoute(c *gin.Context) {
	query := c.Request.URL.Query()
	since := query.Get("since")

	user, err := context.GetAuthUser(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	userContacts, err := App.Controller.GetUserContactsController(since, user)
	if err != nil {
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}

	response.RenderResponse(c, userContacts)
}

// Create User Contacts
// @Summary 	Create User contacts
// @Description
// @Tags		User Contacts
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		body 		body 		map[string]interface{} true "Create User Contacts"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/user/contacts [post]
func (r userRoutes) UpdateUserContactsRoute(c *gin.Context) {
	var body map[string]interface{}
	if err := apiutil.Bind(c, &body); err != nil {
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}

	user, err := context.GetAuthUser(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	err = App.Controller.CreateUserContactsController(body, user)
	if err != nil {
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}

	response.RenderSuccessResponse(c, "User Contacts Created Successfully")
}

// Setup User details
// @Summary 	Setup User
// @Description
// @Tags		User
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		file		formData file false "Body with image file"
// @Param 		fullName		formData string true "Last Name"
// @Param 		username		formData string true "Username"
// @Param 		bio		formData string true "Bio"
// @Param 		twitterUrl		formData string true "Twitter Url"
// @Param 		website		formData string true "Website"
// @Param 		location		formData string true "Location"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/user/setup [post]
func (r userRoutes) setupUserRoute(c *gin.Context) {
	var body *UpdateUserValidator
	if err := apiutil.Bind(c, &body); err != nil {
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}

	user, err := context.GetAuthUser(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	user, err = App.Controller.SetupUserController(body, user)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	response.RenderResponse(c, UserGetSerializerData(user, nil))
	return
}
