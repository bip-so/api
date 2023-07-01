package bootstrap

import (
	"fmt"
	"net/http"

	"gitlab.com/phonepost/bip-be-platform/internal/canvasrepo"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"

	"gitlab.com/phonepost/bip-be-platform/pkg/utils"

	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/internal/follow"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/internal/studio"
	"gitlab.com/phonepost/bip-be-platform/internal/user"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
	"gorm.io/gorm"
)

var (
	bootstrapRouteHandler BootstrapRoutes
)

type BootstrapRoutes struct{}

// Special API: Bootstrap handle API (Public API)
// @Summary 	This API gets a "string" and then will return Studio or User
// @Tags		Bootstrap
// @Accept 		json
// @Produce 	json
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/bootstrap/handle [get]
func (br *BootstrapRoutes) Handle(c *gin.Context) {
	// handle
	handle := c.Param("handle")
	authUsr, _ := c.Get("currentUser")
	var authUser *models.User
	if authUsr != nil {
		authUser = authUsr.(*models.User)
	}

	studioInstance, members, userInstance, userFollows, err := bootstrapController.HandleController(handle, authUser)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "No Studio or User found with " + handle,
			})
			return
		}
		response.RenderErrorResponse(c, response.ServerError(c, err))
		return
	}
	if userInstance != nil {
		// USER FOUND
		//UserGetSerializerData
		serialzed := user.UserGetSerializerData(userInstance, userFollows)
		resp, _ := follow.App.Controller.GetUserFollowFollowCountHandler(userInstance.ID)
		serialzed.Followers = resp.Followers
		serialzed.Following = resp.Following
		response.RenderResponse(c, map[string]interface{}{
			"data":    serialzed,
			"context": "user",
		})
		fmt.Println(userInstance)
		return
	}
	// STUDIO FOUND
	serialzed := studio.SerializeStudioForUser(studioInstance, authUser, members)
	if authUser == nil {
		serialzed.Permission = models.PGStudioNoneSysName
		canvasRepo, _ := canvasrepo.App.Service.GetFirstPublicCanvasOfStudio(studioInstance.ID)
		if canvasRepo != nil {
			serialzed.DefaultCanvasRepoID = canvasRepo.ID
			serialzed.DefaultCanvasBranchID = *canvasRepo.DefaultBranchID
			serialzed.DefaultCanvasRepoKey = canvasRepo.Key
			serialzed.DefaultCanvasRepoName = canvasRepo.Name
		}
	} else {
		permissionList, _ := permissions.App.Service.CalculateStudioPermissions(authUser.ID)
		serialzed.Permission = permissionList[serialzed.ID]
		canvasRepo, _ := canvasrepo.App.Service.GetFirstUserCanvasOfStudio(studioInstance.ID, authUser)

		if canvasRepo != nil {
			serialzed.DefaultCanvasRepoID = canvasRepo.ID
			serialzed.DefaultCanvasBranchID = *canvasRepo.DefaultBranchID
			serialzed.DefaultCanvasRepoKey = canvasRepo.Key
			serialzed.DefaultCanvasRepoName = canvasRepo.Name
		}
	}
	memberCount, _ := queries.App.MemberQuery.GetMemberCountForStudio(studioInstance.ID)
	serialzed.MembersCount = memberCount
	response.RenderResponse(c, map[string]interface{}{
		"data":    serialzed,
		"context": "studio",
	})
}

// Bootstrap Get
// @Summary 	Gets all data required
// @Description
// @Tags		Bootstrap
// @Accept 		json
// @Produce 	json
// @Param 		userId 	query 	string	false "user Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/bootstrap/get [get]
func (br *BootstrapRoutes) Get(c *gin.Context) {
	var userID uint64
	ctxUser, _ := c.Get("currentUser")
	if ctxUser != nil {
		loggedInUser := ctxUser.(*models.User)
		userID = loggedInUser.ID
	}
	userIDstr := c.Query("userId")
	if userIDstr != "" {
		userID = utils.Uint64(userIDstr)
	}

	userAssociatedStudio, err := bootstrapController.GetUserAssociatedStudios(userID)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	studios := SerializeAssociatedStudios(userAssociatedStudio.StudiosData)
	response.RenderResponse(c, map[string]interface{}{
		"userStudios": studios,
	})
}

func (br *BootstrapRoutes) Getuser(c *gin.Context) {
	useridstr := c.Param("userid")
	userID := utils.Uint64(useridstr)
	result,_:=bootstrapController.GetUserAssociatedStudios(userID)

	response.RenderResponse(c, result)
	return
}

