package global

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"gitlab.com/phonepost/bip-be-platform/internal/follow"
	"gitlab.com/phonepost/bip-be-platform/pkg/core"

	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/internal/auth"
	"gitlab.com/phonepost/bip-be-platform/internal/discord"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/studio"
	"gitlab.com/phonepost/bip-be-platform/internal/user"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/apiutil"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/search"
)

var (
	GlobalRouteHandler globalRoutes
)

type globalRoutes struct{ core.RouteHelper }

// Create New Stuido
// @Summary 	Creates a new studio for the user
// @Description This API will also do the following
// @Description Create 2 Roles - Member and Admin
// @Description Create a Member Instace for this USER for this STUDIO
// @Description Adds this users ADMIN Role as a Memeber
// @Description We create a new collection "My First Collection"
// @Description Inside this collection we create "Default Canvas"
// @Tags		Studio
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		body 		body 		studio.CreateStudioValidator true "Create Studio Data"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studio/create [post]
func (gr globalRoutes) createStudioRouteHandler(c *gin.Context) {
	authUsr, _ := c.Get("currentUser")
	if authUsr == nil {
		response.RenderPermissionError(c)
		return
	}
	authUser := authUsr.(*models.User)
	var body studio.CreateStudioValidator
	if err := apiutil.Bind(c, &body); err != nil {
		fmt.Println(err, body)
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}

	studioObj, err := GlobalController.CreateStudioController(&body, authUser)
	if err != nil {
		fmt.Println(err)
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	canvasRepo, _ := FirstCanvasRepo(studioObj.ID)

	fmt.Println("Reached End \n")
	fmt.Println(err)
	resp := studio.SerializeStudio(studioObj)
	resp.Permission = models.PGStudioAdminSysName
	if canvasRepo != nil {
		resp.DefaultCanvasRepoID = canvasRepo.ID
		resp.DefaultCanvasRepoName = canvasRepo.Name
		resp.DefaultCanvasRepoKey = canvasRepo.Key
		resp.DefaultCanvasBranchID = *canvasRepo.DefaultBranchID
	}
	response.RenderResponse(c, resp)
}

// Check handle available
// @Summary 	Check user/studio username/handle available
// @Description
// @Tags		Global
// @Accept 		json
// @Produce 	json
// @Param 		handle 	query 		string		 		false "handle to check"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/global/check-handle [get]
func (gr globalRoutes) checkHandleAvailableRouteHandler(c *gin.Context) {
	handle := c.Query("handle")
	available, err := GlobalController.CheckHandleAvailable(handle)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	response.RenderCustomResponse(c, map[string]interface{}{
		"available": available,
	})
}

// Get Popular Users
// @Summary 	Gets popular users
// @Description
// @Tags		User
// @Accept 		json
// @Produce 	json
// @Param 		skip 	query 		string		 		false "next page"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/user/popular [get]
func (r globalRoutes) popularUsersRoute(c *gin.Context) {
	skip := c.Query("skip")
	skipInt, _ := strconv.Atoi(skip)
	ctxUser, _ := c.Get("currentUser")
	var loggedInUser *models.User
	if ctxUser != nil {
		loggedInUser = ctxUser.(*models.User)
	}
	users, followUsers, err := GlobalController.PopularUsersController(loggedInUser, skipInt)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	var serializedUsers []user.UserGetSerializer
	for _, usr := range *users {
		resp, _ := follow.App.Controller.GetUserFollowFollowCountHandler(usr.ID)
		userData := user.UserGetSerializerData(&usr, followUsers)
		userData.Followers = resp.Followers
		userData.Following = resp.Following
		serializedUsers = append(serializedUsers, userData)
	}
	next := "-1"
	if len(*users) == configs.PAGINATION_LIMIT {
		next = strconv.Itoa(skipInt + len(*users))
	}
	response.RenderPaginatedResponse(c, serializedUsers, next)
}

// Search
// @Summary 	Search Studios & Users
// @Description
// @Tags		Global
// @Accept 		json
// @Produce 	json
// @Param 		skip 	query 		int		 		false "next page"
// @Param 		query 	query 		string		 		false "query to search"
// @Param 		type 	query 		string		 		false "type of objects: studios/users/pages/reels, empty for all"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/global/search [get]
func (r globalRoutes) searchRouteHandler(c *gin.Context) {
	skip := c.Query("skip")
	skipInt, _ := strconv.Atoi(skip)
	query := c.Query("query")
	objectType := c.Query("type")
	ctxUser, _ := c.Get("currentUser")
	var loggedInUser *models.User
	if ctxUser != nil {
		loggedInUser = ctxUser.(*models.User)
	}
	studioID, _ := r.GetStudioId(c)

	studioDocs, userDocs, roleDocs, reelDocs, err := GlobalController.SearchController(query, objectType, loggedInUser, skipInt, studioID)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	var studioDocsLen float64 = 0
	var userDocsLen float64 = 0
	var reelDocsLen float64 = 0
	if studioDocs != nil {
		studioDocsLen = float64(len(studioDocs))
	}
	if userDocs != nil {
		userDocsLen = float64(len(userDocs))
	}
	if reelDocs != nil {
		reelDocsLen = float64(len(userDocs))
	}
	maxOfTwo := math.Max(userDocsLen, studioDocsLen)
	maximum := math.Max(maxOfTwo, reelDocsLen)
	next := "-1"
	if maximum == configs.PAGINATION_LIMIT {
		next = strconv.Itoa(skipInt + configs.PAGINATION_LIMIT)
	}
	data := map[string]interface{}{
		"users":   userDocs,
		"studios": studioDocs,
		"roles":   roleDocs,
		"reels":   reelDocs,
	}
	response.RenderPaginatedResponse(c, data, next)
}

func (r globalRoutes) addSearch(c *gin.Context) {

	search.GetIndex(search.UserDocumentIndexName).SetupIndex([]string{"fullName", "handle"})
	search.GetIndex(search.StudioDocumentIndexName).SetupIndex([]string{"displayName", "description", "handle"})

	err := studio.App.StudioService.AddAllStudiosToAlgolia()
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	err = user.App.Service.AddAllUsersToAlgolia()
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	response.RenderSuccessResponse(c, "success")
}

// Generic Upload Image
// @Summary 	Global generic upload image API
// @Description
// @Tags		Global
// @Security 	bearerAuth
// @Param 		file		formData file true "File"
// @Param 		model		formData string true "Model Name"
// @Param 		uuid		formData string true "Model uuid"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/global/upload-file [post]
func (r *globalRoutes) imageRouteHandler(c *gin.Context) {

	var body ImageUpload
	if err := apiutil.Bind(c, &body); err != nil {
		fmt.Println(err, body)
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}
	ctxStudio, _ := c.Get("currentStudio")
	ctxStudioAsInt := ctxStudio.(uint64)

	file, err := c.FormFile("file")
	if err != nil {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": "file not found",
		})
		return
	}
	object, _ := file.Open()
	resp, err := GlobalController.UpdateImage(object, strings.ToLower(body.Model), body.UUID, file.Filename, ctxStudioAsInt, body.RepoID)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	response.RenderResponse(c, resp)
}

func (r *globalRoutes) connectDiscord(c *gin.Context) {
	query := c.Request.URL.Query()
	userIdStr := query.Get("userId")
	studioIdStr := query.Get("studioId")
	partnerIntegrationID := query.Get("partnerIntegrationId")
	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	studioId, err := strconv.ParseUint(studioIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//user, err := user.App.Repo.GetUser(map[string]interface{}{"id": userId})
	user, err := queries.App.UserQueries.GetUserByID(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//product, err := studio.App.StudioRepo.GetStudioByID(studioId)
	//if err != nil || product == nil {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"error": err.Error(),
	//	})
	//	return
	//}
	if studioId != 0 {
		activeDiscord, err := discord.GetDiscordStudioIntegration(studioId)
		if err != nil {
			println("error while getting active discord integration for studio id", studioId)
		}
		//check if active discord integration already exist for this studio
		if len(activeDiscord) != 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "discord integration already exists for this studio",
			})
			return
		}
	}
	jwtPayload := map[string]interface{}{
		"studioId":             studioId,
		"partnerIntegrationID": partnerIntegrationID,
	}
	token, _ := auth.CreateCustomJWTToken(user.ID, jwtPayload)

	authorizationURL := "https://discord.com/api/oauth2/authorize"
	scope := "bot email identify applications.commands"
	//apihost := configs.AppConfig.APIHost
	// user, err := middlewares.AuthenticatedUser(r.Context())
	// if err != nil {
	// 	println("user authorisation error")
	// }

	clientId := configs.GetDiscordBotConfig().ClientID
	permission := configs.GetDiscordBotConfig().Permission
	// var redirectUrl string
	// if token == "" {
	// 	redirectUrl = configs.GetAppInfoConfig().FrontendHost + "connect_integration?provider=discord"
	// redirectUrl := configs.GetAppInfoConfig().FrontendHost + "/@" + product.Handle + "/settings?provider=discord&status=success"
	// } else {
	// 	redirectUrl = configs.GetAppInfoConfig().BackendHost + "/api/discord/connect_discord"
	// }
	redirectUrl := configs.GetDiscordLambdaServerEndpoint() + "/discord/connect"

	reqParams := url.Values{
		"client_id":     {clientId},
		"permissions":   {permission},
		"scope":         {scope},
		"state":         {token},
		"redirect_uri":  {redirectUrl},
		"response_type": {"code"},
	}
	redirectURL := authorizationURL + "?" + reqParams.Encode()
	c.Redirect(http.StatusFound, redirectURL)
}

// Get Messages API
// @Summary 	Get Messages
// @Description Gets Dicord Messages captured from bipmark for the user
// @Tags		Messages
// @Security 	bearerAuth
// @Param 		skip 	query 		string		 		false "next page"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/message/get [get]
func (r *globalRoutes) getMessages(c *gin.Context) {
	skip, err := strconv.Atoi(c.Query("skip"))
	if err != nil {
		skip = 0
	}
	userID := r.GetLoggedInUserId(c)

	serializeMessages, err := GlobalController.GetMessages(userID, skip)
	if err != nil {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	response.RenderResponse(c, serializeMessages)
}

// delete Message API
// @Summary 	Delete Message
// @Description Deletes Discord Message captured from bipmark for the user
// @Tags		Messages
// @Security 	bearerAuth
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/message/{messageID} [delete]
func (r *globalRoutes) deleteMessage(c *gin.Context) {
	messageID, err := strconv.ParseUint(c.Param("messageID"), 10, 64)
	if err != nil {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	userID := r.GetLoggedInUserId(c)

	err = GlobalController.DeleteMessage(userID, messageID)
	if err != nil {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	response.RenderSuccessResponse(c, "message deleted")
}
