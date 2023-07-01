package notifications

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

// Get Notifications
// @Summary 	Get Notifications
// @Description
// @Tags		Notifications
// @Security 	bearerAuth
// @Param 		type 	path 	string	true "Type Eg. all, studio, personal"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/notifications/ [get]
func (r notificationRoutes) getNotifications(c *gin.Context) {
	user, err := r.GetLoggedInUser(c)
	studioID, err := r.GetStudioId(c)
	skip := c.Query("skip")
	skipInt, _ := strconv.Atoi(skip)
	limit := c.Query("limit")
	limitInt, _ := strconv.Atoi(limit)
	notificationType := c.Query("type")
	getStreamKey := utils.String(user.ID)
	filter := c.Query("filter")
	if notificationType == "studio" {
		getStreamKey += "-" + utils.String(studioID)
	}
	if notificationType == "personal" {
		getStreamKey += "-personal"
	}
	//resp, err := App.Service.GetStreamNotificationForUser(getStreamKey, skipInt, limitInt)
	resp, err := App.Service.GetDBNotificationForUser(user, skipInt, limitInt, notificationType, studioID, filter)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	respSerialized := MultiSerializeNotification(resp)
	next := "-1"
	if len(resp) == configs.PAGINATION_LIMIT {
		next = strconv.Itoa(skipInt + len(respSerialized))
	}
	response.RenderPaginatedResponse(c, respSerialized, next)
	go func() {
		App.Controller.MarkNotificationsAsRead(user.ID)
	}()
	return
}

// Notifications Mark as read
// @Summary 	Mark as read
// @Description
// @Tags		Notifications
// @Security 	bearerAuth
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/notifications/mark-as-seen [post]
func (r notificationRoutes) markAsSeen(c *gin.Context) {
	user, err := r.GetLoggedInUser(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	err = App.Controller.MarkNotificationsAsRead(user.ID)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	response.RenderSuccessResponse(c, "Done")
	return
}

// Get Notifications count of user
// @Summary 	All Notifications Count of user
// @Description
// @Tags		Notifications
// @Security 	bearerAuth
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/notifications/count [get]
func (r notificationRoutes) notificationCount(c *gin.Context) {
	user, err := r.GetLoggedInUser(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	notificationCount, err := App.Repo.GetNotificationCount(map[string]interface{}{"user_id": user.ID})
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	response.RenderResponse(c, SerializeNotificationCount(*notificationCount))
	return
}

func (r notificationRoutes) updateNotification(c *gin.Context) {
	// Need notification ID, GetStream Activity ID,
	// Need to figure out what all the data needs to be updated
	return
}
