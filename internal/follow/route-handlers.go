package follow

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/user"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

/*
	Follw User Route handlers.
*/

// @Summary 	Followers Count for loggedin user
// @Tags		Followings
// @Router 		/v1/follow/user/follow-count [get]
// @Success 	200 		{object} 	FollowUserFollowCountResponse
func (r followRoutes) followUserFollowCountRoute(c *gin.Context) {
	ctxUser, _ := c.Get("currentUser")
	loggedInUser := ctxUser.(*models.User)
	resp, err := App.Controller.GetUserFollowFollowCountHandler(loggedInUser.ID)
	if err != nil {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": "Generic, Error while fetching data",
		})
		return
	}

	serializer := FollowUserSerializer{c}
	response.RenderOkWithData(c, serializer.GetUserFollowCounts(resp))
	return
}

// @Summary 	Follow a User
// @Tags		Followings
// @Router 		/v1/follow/user/follow [post]
// @Param 		body body 		PostFollowUserRequest true "Follow User"
func (r followRoutes) followUserRoute(c *gin.Context) {
	var body PostFollowUserRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctxUser, _ := c.Get("currentUser")
	loggedInUser := ctxUser.(*models.User)

	App.Repo.UserFollowUser(*loggedInUser, body.UserId)

	response.RenderSuccessResponse(c, "Done.")
	return
}

// @Summary 	Unfollow a User
// @Tags		Followings
// @Router 		/v1/follow/user/unfollow [post]
// @Param 		body body 		PostUnFollowUserRequest true "unFollow User"
func (r followRoutes) unfollowUserRoute(c *gin.Context) {
	var body PostUnFollowUserRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	authUser, err := r.GetLoggedInUser(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	App.Repo.UserFollowUnfollowUser(*authUser, body.UserId)

	response.RenderSuccessResponse(c, "Done.")
	return
}

/*
	Follw Studio Route handlers.
*/

// @Summary 	Followers Count for Studio
// @Tags		Followings
// @Router 		/v1/follow/studio/follower [get]
// @Success 	200 		{object} 	FollowUserStudioCountResponse
func (r followRoutes) followStudioCountRoute(c *gin.Context) {
	studioId, status := middlewares.StudioHeaderRequiredCheck(c)
	if !status {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": "Studio Header Missing",
		})
		return
	}

	resp, err := App.Controller.GetStudioFollowersCountHandler(studioId)

	if err != nil {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": "Generic, Error while fetching data",
		})
		return
	}

	serializer := FollowStudioSerializer{c}
	response.RenderOkWithData(c, serializer.GetStudioFollowCounts(resp))
	return
}

// @Summary 	Follow a Studio
// @Tags		Followings
// @Router 		/v1/follow/studio/follow [post]
// @Param 		body body 		PostFollowStudioRequest true "Follow studio"
func (r followRoutes) followUserStudioRoute(c *gin.Context) {
	var body PostFollowStudioRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	studioId, status := middlewares.StudioHeaderRequiredCheck(c)
	if !status {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": "Studio Header Missing",
		})
		return
	}

	App.Repo.UserFollowStudio(studioId, body.UserId)

	response.RenderSuccessResponse(c, "Done.")
	return
}

// @Summary 	Unfollow a Studio
// @Tags		Followings
// @Router 		/v1/follow/studio/unfollow [post]
// @Param 		body body 		PostUnFollowStudioRequest true "unFollow User"
func (r followRoutes) unfollowStudioRoute(c *gin.Context) {
	var body PostUnFollowStudioRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// STUDIO ID AS UINT64
	studioId, status := middlewares.StudioHeaderRequiredCheck(c)
	if !status {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": "Studio Header Missing",
		})
		return
	}

	App.Repo.UserUnFollowStudio(studioId, body.UserId)

	response.RenderSuccessResponse(c, "Done.")
	return
}

// @Summary 	Get user followers or following a Studio
// @Tags		Followings
// @Router 		/v1/follow/user/list [get]
func (r followRoutes) FollowList(c *gin.Context) {
	followType := c.Query("type")
	userIDstr := c.Query("userId")
	userID := utils.Uint64(userIDstr)
	var err error
	var data []user.UserMiniSerializer
	if followType == "following" {
		// Get user following users
		data, err = App.Controller.GetUserFollowing(userID)
		if err != nil {
			response.RenderCustomErrorResponse(c, err)
		}
	} else {
		// Get user followers
		data, err = App.Controller.GetUserFollowers(userID)
		if err != nil {
			response.RenderCustomErrorResponse(c, err)
		}
	}
	response.RenderResponse(c, data)
	return
}
