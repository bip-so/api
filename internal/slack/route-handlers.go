package slack2

import (
	"encoding/json"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/apiutil"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/apiClient"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"gitlab.com/phonepost/bip-be-platform/internal/auth"
	"gitlab.com/phonepost/bip-be-platform/internal/studio"
	studio2 "gitlab.com/phonepost/bip-be-platform/internal/studio"
	"gitlab.com/phonepost/bip-be-platform/internal/studio_integration"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/integrations"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

// Slack Integration
// @Summary 	Connect studio to slack
// @Description
// @Tags		slack
// @Accept 		json
// @Produce 	json
// @Param 		userId 	query 		string		 		false "userId"
// @Param 		studioId 	query 		string		 		false "userId"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/slack/connect [get]
func (r *slackRoutes) connectSlack(c *gin.Context) {
	query := c.Request.URL.Query()
	userIdStr := query.Get("userId")
	studioIdStr := query.Get("studioId")
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
	product, err := studio.App.StudioRepo.GetStudioByID(studioId)
	if err != nil || product == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	activeSlack, err := GetSlackStudioIntegration(studioId)
	if err != nil {
		println("error while getting active slack integration for studio id", studioId)
	}
	if len(activeSlack) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "slack integration already exists for this studio",
		})
		return
	}
	jwtPayload := map[string]interface{}{
		"studioId": studioId,
	}
	state, _ := auth.CreateCustomJWTToken(user.ID, jwtPayload)
	authorizationURL := "https://slack.com/oauth/authorize"
	clientId := configs.GetSlackConfig().ClientID
	scope := "files:read reactions:write usergroups:read users.profile:read users:read users:read.email team:read channels:read channels:history groups:read groups:history im:read im:history mpim:read mpim:history incoming-webhook bot chat:write:bot commands"
	redirectUri := configs.GetSlackLambdaServerEndpoint() + "/slack/connect"
	reqParams := url.Values{
		"client_id":    {clientId},
		"scope":        {scope},
		"state":        {state},
		"redirect_uri": {redirectUri},
	}
	redirectURL := authorizationURL + "?" + reqParams.Encode()
	c.Redirect(http.StatusFound, redirectURL)
}

func (r *slackRoutes) authorize(c *gin.Context) {
	query := c.Request.URL.Query()
	token := query.Get("token")
	studioID, err := strconv.ParseUint(query.Get("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// next := query.Get("next")

	println("body: token:", token, " =>prodID:", studioID)

	var state string

	if token == "" && studioID == 0 {
		// Flow without passing studioId (from slack directory)
		state = ""
	} else {
		result, isValid, err := auth.ParseJWTToken(token)
		if err != nil {
			println("error while parsing token", err.Error())
		}
		println("token", result, isValid, " ", studioID)
		uid := result["uid"].(string)
		//user, err := user2.App.Repo.GetUser(map[string]interface{}{"id": uid})
		user, err := queries.App.UserQueries.GetUserByID(utils.Uint64(uid))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		studio, err := studio2.App.StudioRepo.GetStudioByID(studioID)
		if err != nil || studio == nil {
			logger.Debug("connectSlack: Error while parsing body")
			logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		println("studioID", studio != nil, studio)
		activeSlack, err := studio_integration.App.Repo.GetActiveIntegrationForStudio(studioID, studio_integration.SLACK_INTEGRATION_TYPE)
		if err != nil {
			println("error while getting active slack integration for studio id", studioID)
		}
		println("activeSlack:", activeSlack != nil, " ID:", activeSlack != nil && activeSlack.ID == 0)
		//check if active slack integration already exist for this studio
		if activeSlack != nil && activeSlack.ID != 0 {
			erroMsg, err := json.Marshal(map[string]string{
				"error": "slack integration already exists for this studio",
			})
			c.Writer.Write(erroMsg)
			logger.Debug("connectSlack: Error while parsing body")
			logger.Error(err.Error())
			response.RenderErrorResponse(c, response.BadDataError(c.Request.Context()))
			return
		}

		jwtPayload := map[string]interface{}{
			"studioId": studio.ID,
		}
		println("creating custom jwt tpoken")
		jwtToken, err := auth.CreateCustomJWTToken(user.ID, jwtPayload)
		if err != nil {
			response.RenderErrorResponse(c, response.BadDataError(c.Request.Context()))
			return
		}
		state = jwtToken
	}

	authorizationURL := "https://slack.com/oauth/authorize"
	scope := "users.profile:read users:read users:read.email team:read channels:read channels:history groups:read groups:history im:read im:history mpim:read mpim:history incoming-webhook bot chat:write:bot commands"
	clientId := configs.GetSlackConfig().ClientID

	var redirectUri string
	if state == "" {
		redirectUri = configs.GetAppInfoConfig().FrontendHost + "connect_integration?provider=slack"
	} else {
		redirectUri = configs.GetAppInfoConfig().BackendHost + "/api/slack/connect_slack"
	}

	reqParams := url.Values{
		"client_id":    {clientId},
		"scope":        {scope},
		"state":        {state},
		"redirect_uri": {redirectUri},
	}
	redirectURL := authorizationURL + "?" + reqParams.Encode()
	c.Redirect(http.StatusFound, redirectURL)
}

func (r *slackRoutes) connectSlackLogin(c *gin.Context) {
	type slackConnectResponse struct {
		Code           string `json:"code"`
		responseSecret string `json:"response_secret"`
	}
	var body slackConnectResponse
	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	slackConfig := configs.GetSlackConfig()
	httpClient := http.Client{}
	oathRedirect := configs.GetAppInfoConfig().FrontendHost + "connect_integration?provider=slack"
	response, _err := slack.GetOAuthResponse(&httpClient, slackConfig.ClientID,
		slackConfig.ClientSecret, body.Code, oathRedirect)

	if _err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": _err.Error(),
		})
		return
	}

	dataToEncrypt, err := json.Marshal(response)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	fmt.Println("line1:", response, response.UserID)
	slackProfile, err := integrations.GetSlackProfile(response.AccessToken, response.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	slackWorkspace, err := integrations.GetSlackTeam(response.AccessToken, response.TeamID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	fmt.Println("line2:", slackProfile)
	userName := strings.ReplaceAll(slackProfile.UserProfile.Email, "@", "-")
	user, err := UpsertSlackUser(response.UserID, slackProfile.UserProfile.DisplayName,
		userName, slackProfile.UserProfile.Email, slackProfile.UserProfile.Avatar, response.AccessToken, response.TeamID, body.responseSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	encryptedData := utils.Encrypt([]byte(configs.GetSecretShortKey()), string(dataToEncrypt))
	extras := map[string]interface{}{
		"data": encryptedData,
		"slackTeam": map[string]string{
			"name":  slackWorkspace.Team.Name,
			"image": slackWorkspace.Team.Icon.Image132,
		},
	}
	RenderAccountWithExtras(c, user, extras)
}

func (r *slackRoutes) SlackShortcutsHandler(c *gin.Context) {
	buf, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read request body: %s", err)
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}
	jsonStr, err := url.QueryUnescape(string(buf)[8:])
	if err != nil {
		log.Printf("[ERROR] Failed to unescape request body: %s", err)
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}

	var message SlackAppMentionPayload
	if err := json.Unmarshal([]byte(jsonStr), &message); err != nil {
		log.Printf("[ERROR] Failed to decode json message from slack: %s", jsonStr)
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}

	messageStr, _ := json.Marshal(message)
	apiClient.AddToQueue(apiClient.SlackBipMarkAction, messageStr, apiClient.DEFAULT, apiClient.CommonRetry)
	response.RenderSuccessResponse(c, "Captured the message successfully")
	return
}

func (r *slackRoutes) SlackEventsHandler(c *gin.Context) {
	var body map[string]interface{}
	if err := apiutil.Bind(c, &body); err != nil {
		fmt.Println(err, body)
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}
	if body["challenge"] != nil && len(body["challenge"].(string)) > 0 {
		c.Writer.Write([]byte(body["challenge"].(string)))
		return
	}
	bodyStr, _ := json.Marshal(body)
	apiClient.AddToQueue(apiClient.SlackEventSubscriptions, bodyStr, apiClient.DEFAULT, apiClient.CommonRetry)
	response.RenderSuccessResponse(c, "Event received successfully")
	return
}

func (r *slackRoutes) SlashCommandsHandler(c *gin.Context) {
	s, err := slack.SlashCommandParse(c.Request)
	if err != nil {
		fmt.Println("Error in parsing slash command  ==> ", err)
		response.RenderCustomResponse(c, map[string]interface{}{"text": err.Error()})
		return
	}
	bodyStr, _ := json.Marshal(s)
	apiClient.AddToQueue(apiClient.SlackSlashCommands, bodyStr, apiClient.DEFAULT, apiClient.CommonRetry)
	return
}
