package discord

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	ar "gitlab.com/phonepost/bip-be-platform/internal/accessrequest"
	"gitlab.com/phonepost/bip-be-platform/internal/canvasrepo"
	"gitlab.com/phonepost/bip-be-platform/internal/collection"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/internal/message"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	user2 "gitlab.com/phonepost/bip-be-platform/internal/user"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/integrations"
)

// func (impl *DiscordImpl) authorizeDiscord(c *gin.Context) {

// 	query := c.Request.URL.Query()
// 	token := query.Get("token")
// 	studioID, err := strconv.ParseUint(query.Get("id"), 10, 64)
// 	c.JSON(http.StatusBadRequest, gin.H{
// 		"error": err.Error(),
// 	})

// 	if token != "" || studioID != 0 {

// 		result, isValid, err := auth.ParseJWTToken(token)
// 		if err != nil {
// 			println("error while parsing token", err.Error())
// 		}
// 		println("token", result, isValid, " ", studioID)
// 		uid := result["uid"].(string)
// 		user, _ := user2.App.Repo.GetUser(map[string]interface{}{"id": uid})

// 		studio, err := studio2.App.StudioRepo.GetStudioByID(studioID)
// 		if err != nil || studio == nil {
// 			logger.Debug("connectDiscord: Error while parsing body")
// 			logger.Error(err.Error())
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"error": err.Error(),
// 			})
// 			return

// 		}
// 		println("studioID", studio != nil, studio.ID)
// 		activeDiscord, err := studioIntegration.App.Repo.GetActiveIntegrationForStudio(studio.ID, studioIntegration.DISCORD_INTEGRATION_TYPE)
// 		if err != nil {
// 			println("error while getting active discord integration for studio id", studioID)
// 		}
// 		println("activeDiscord:", activeDiscord != nil, " ID:", activeDiscord.ID == 0)
// 		//check if active discord integration already exist for this studio
// 		if activeDiscord != nil && activeDiscord.ID != 0 {
// 			erroMsg, err := json.Marshal(map[string]string{
// 				"error": "discord integration already exists for this studio",
// 			})
// 			c.Writer.Write(erroMsg)
// 			logger.Debug("connectDiscord: Error while parsing body")
// 			logger.Error(err.Error())

// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"error": err.Error(),
// 			})
// 			return
// 		}

// 		jwtPayload := map[string]interface{}{
// 			"studioId": studioID,
// 		}
// 		println("creating custom jwt tpoken")

// 		token, _ = auth.CreateCustomJWTToken(user.ID, jwtPayload)

// 	}

// 	authorizationURL := "https://discord.com/api/oauth2/authorize"
// 	scope := "bot email identify webhook.incoming applications.commands"

// 	clientId := configs.GetDiscordBotConfig().ClientID
// 	permission := configs.GetDiscordBotConfig().Permission
// 	var redirectUrl string
// 	if token == "" {
// 		redirectUrl = configs.GetAppInfoConfig().FrontendHost + "connect_integration?provider=discord"
// 	} else {
// 		redirectUrl = configs.GetAppInfoConfig().BackendHost + "/api/v1/integrations/discord/connect"
// 	}

// 	reqParams := url.Values{
// 		"client_id":     {clientId},
// 		"permissions":   {permission},
// 		"scope":         {scope},
// 		"state":         {token},
// 		"redirect_uri":  {redirectUrl},
// 		"response_type": {"code"},
// 	}
// 	redirectURL := authorizationURL + "?" + reqParams.Encode()
// 	println("redirect url", redirectURL)
// 	c.Redirect(http.StatusFound, redirectURL)
// }

// func (impl *DiscordImpl) connectDiscordLogin(c *gin.Context) {
// 	type discordConnectResponse struct {
// 		Code  string `json:"code"`
// 		State string `json:"state"`
// 	}
// 	var body discordConnectResponse
// 	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	endpoint := "https://discord.com/api/oauth2/token"
// 	contentType := "application/x-www-form-urlencoded"

// 	redirectUrl := configs.GetConfigString("SITEROOT") + "connect_integration?provider=discord"

// 	values := url.Values{
// 		"client_id":     {configs.GetDiscordBotConfig().ClientID},
// 		"client_secret": {configs.GetDiscordBotConfig().ClientSecret},
// 		"grant_type":    {"authorization_code"},
// 		"code":          {body.Code},
// 		"redirect_uri":  {redirectUrl},
// 	}
// 	reqBody := strings.NewReader(values.Encode())
// 	_response, err := http.Post(endpoint, contentType, reqBody)
// 	println("url", endpoint, " value", values.Encode(), " contetType:", contentType)
// 	_result, _ := ioutil.ReadAll(_response.Body)
// 	fmt.Println("\n", string(_result), "\n\n======")

// 	data := map[string]interface{}{}
// 	err = json.Unmarshal(_result, &data)
// 	if err != nil || data == nil {
// 		println("error while structToMap", err.Error(), data)
// 		return
// 	}
// 	if code, ok := data["code"]; ok {
// 		if code.(float64) == 30007 {
// 			response.RenderCustomResponse(c, map[string]interface{}{
// 				"error":   true,
// 				"message": "maximum_webhook_10",
// 			})
// 		} else {
// 			response.RenderCustomResponse(c, map[string]interface{}{
// 				"error":   true,
// 				"message": "there was some problem",
// 			})
// 		}
// 		return
// 	}

// 	tId, ok := data["webhook"].(map[string]interface{})
// 	if !ok {
// 		return
// 	}
// 	guildId, ok := tId["guild_id"].(string)
// 	if !ok {
// 		return
// 	}
// 	guild, err := integrations.GetDiscordTeam(guildId)
// 	if err != nil {
// 		fmt.Println("failed to get a guild")
// 		return
// 	}

// 	dataToEncrypt, err := json.Marshal(data)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	encryptedData := utils.Encrypt([]byte(configs.GetSecretShortKey()), string(dataToEncrypt))
// 	extras := map[string]interface{}{
// 		"data": encryptedData,
// 		"discordTeam": map[string]string{
// 			"name":  guild.Name,
// 			"image": guild.IconURL(),
// 		},
// 	}
// 	response.RenderCustomResponse(c, extras)
// }

func isSessionVerifed(c *gin.Context) bool {
	pkey, _ := hex.DecodeString(configs.GetDiscordBotConfig().PublicKey)
	isVerified := discordgo.VerifyInteraction(c.Request, ed25519.PublicKey(pkey))

	if !isVerified {
		//c.String(401, "invalid request signature")
		return false

	}
	return true
}

func (impl *DiscordImpl) interaction(c *gin.Context) {
	/*	pkey, _ := hex.DecodeString(configs.GetDiscordBotConfig().PublicKey)
		isVerified := discordgo.VerifyInteraction(c.Request, ed25519.PublicKey(pkey))

		if !isVerified {
			c.String(401, "invalid request signature")
			return
		}*/
	if !isSessionVerifed(c) {
		c.String(401, "invalid request signature")
		return
	}

	fmt.Println("⇾ Discord interaction API Called.")
	var body Interaction
	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		fmt.Printf("x : interactionDiscord: Error while parsing body %s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Ping Pong
	if body.Type == InteractionPing {
		response.RenderCustomResponse(c, map[string]interface{}{
			"type": 1,
		})
		return
	}

	if body.Type == InteractionApplicationCommandAutocomplete {
		guildID := body.GuildID
		activeDiscord, err := GetProductIntegrationByDiscordTeamId(guildID)
		if err != nil || len(activeDiscord) == 0 {
			response.RenderCustomResponse(c, map[string]interface{}{
				"type": 4,
				"data": map[string]interface{}{
					"tts":     false,
					"content": "Error capturing the message. Please try again.",
					"embeds":  []string{},
					"allowed_mentions": map[string]interface{}{
						"parse": []string{},
					},
					"flags": 1 << 6,
				},
			})
			return
		}
		var userInstance *models.User
		if body.Member != nil {
			socialInstance, _ := FindUsersByDiscordID(body.Member.User.ID)
			userInstance = socialInstance.User
		} else if body.User != nil {
			socialInstance, _ := FindUsersByDiscordID(body.User.ID)
			userInstance = socialInstance.User
		} else {
			response.RenderCustomResponse(c, map[string]interface{}{
				"type": 4,
				"data": map[string]interface{}{
					"tts":     false,
					"content": "Error capturing the message. Please try again.",
					"embeds":  []string{},
					"allowed_mentions": map[string]interface{}{
						"parse": []string{},
					},
					"flags": 1 << 6,
				},
			})
			return
		}
		fmt.Println("UserInstance logging for more information=====>id:", userInstance.ID, "data===?", userInstance)
		if body.Data.Name == "bip-new" {
			suggestions, err := AutoCompleteBipNewHandler(body, activeDiscord, userInstance)
			if err != nil {
				ErrorHandlerMessage(c)
				return
			}
			response.RenderCustomResponse(c, map[string]interface{}{
				"type": 8,
				"data": map[string]interface{}{
					"tts":     false,
					"content": "suggestions",
					"embeds":  []string{},
					"allowed_mentions": map[string]interface{}{
						"parse": []string{},
					},
					"flags":   1 << 6,
					"choices": suggestions,
				},
			})
			return
		} else if body.Data.Name == "bip-search" {
			suggestions, err := AutoCompleteBipSearchHandler(body, activeDiscord, userInstance)
			if err != nil {
				ErrorHandlerMessage(c)
				return
			}
			response.RenderCustomResponse(c, map[string]interface{}{
				"type": 8,
				"data": map[string]interface{}{
					"tts":     false,
					"content": "suggestions",
					"embeds":  []string{},
					"allowed_mentions": map[string]interface{}{
						"parse": []string{},
					},
					"flags":   1 << 6,
					"choices": suggestions,
				},
			})
			return
		}
	}

	if body.Type == InteractionApplicationCommand {
		// Bip Mark
		if body.Data.Type == 3 && body.Data.Name == "bip Mark" { // message command
			messages := []*models.Message{}
			for key, value := range body.Data.Resolved.Messages {
				fmt.Println("messageid:", key, "content:", value.Content)
				var foundAuthor *models.UserSocialAuth
				var foundUser *models.UserSocialAuth
				users, _ := FindUsersByDiscordIDs([]string{value.Author.ID, body.Member.User.ID})
				if len(users) > 0 {
					for i, user := range users {
						if value.Author.ID == user.ProviderID {
							foundAuthor = &users[i]
						}
						if body.Member.User.ID == user.ProviderID {
							foundUser = &users[i]
						}
					}
				}
				if foundAuthor == nil {
					discordUserName := value.Author.Username + value.Author.Discriminator
					//nickName := value.Author.Username
					user := CreateNewUser("", "", discordUserName, value.Author.AvatarURL(""))
					// TODO: pass metadata
					discordUser := NewDiscordUser(user.ID, value.Author.ID, nil)
					err := user2.App.Repo.CreateUserSocialAuth(discordUser)
					if err != nil {
						response.RenderCustomResponse(c, map[string]interface{}{
							"type": 4,
							"data": map[string]interface{}{
								"tts":     false,
								"content": "Error capturing message",
								"embeds":  []string{},
								"allowed_mentions": map[string]interface{}{
									"parse": []string{},
								},
								"flags": 1 << 6,
							},
						})
						return
					}
					foundAuthor = discordUser
				}
				if foundUser == nil {
					discordUserName := body.Member.User.Username + body.Member.User.Discriminator
					//nickName := body.Member.User.Username
					user := CreateNewUser("", "", discordUserName, body.Member.User.AvatarURL(""))
					// TODO: pass metadata
					discordUser := NewDiscordUser(user.ID, body.Member.User.ID, nil)
					err := user2.App.Repo.CreateUserSocialAuth(discordUser)
					if err != nil {
						response.RenderCustomResponse(c, map[string]interface{}{
							"type": 4,
							"data": map[string]interface{}{
								"tts":     false,
								"content": "Error capturing message",
								"embeds":  []string{},
								"allowed_mentions": map[string]interface{}{
									"parse": []string{},
								},
								"flags": 1 << 6,
							},
						})
						return
					}
					foundUser = discordUser
				}
				timestamp := value.Timestamp
				attachments := []string{}
				for _, att := range value.Attachments {
					attachments = append(attachments, att.URL)
				}
				guild, _ := integrations.GetDiscordTeam(body.GuildID)
				message := models.NewDiscordMessage(key, value.Content, foundAuthor.UserID, foundUser.UserID, timestamp, attachments, body.GuildID, guild.Name, guild.IconURL(), "")
				messages = append(messages, message)
			}
			err := message.CreateMessage(messages)
			if err != nil {
				response.RenderCustomResponse(c, map[string]interface{}{
					"type": 4,
					"data": map[string]interface{}{
						"tts":     false,
						"content": "Error capturing the message. Please try again.",
						"embeds":  []string{},
						"allowed_mentions": map[string]interface{}{
							"parse": []string{},
						},
						"flags": 1 << 6,
					},
				})
				return
			}
			response.RenderCustomResponse(c, map[string]interface{}{
				"type": 4,
				"data": map[string]interface{}{
					"tts":     false,
					"content": "Successfully captured the message. You can place it in an appropriate canvas by typing '//' on the canvas",
					"embeds":  []string{},
					"allowed_mentions": map[string]interface{}{
						"parse": []string{},
					},
					"flags": 1 << 6,
				},
			})
			for _, value := range body.Data.Resolved.Messages {
				go integrations.SendDiscordReaction(body.ChannelID, value.ID, []string{"✅"})
			}
			return
		} else if body.Data.Type == 1 && body.Data.Name == "bip-new" {
			guildID := body.GuildID
			activeDiscord, err := GetProductIntegrationByDiscordTeamId(guildID)
			if err != nil || len(activeDiscord) == 0 {
				response.RenderCustomResponse(c, map[string]interface{}{
					"type": 4,
					"data": map[string]interface{}{
						"tts":     false,
						"content": "Error capturing the message. Please try again.",
						"embeds":  []string{},
						"allowed_mentions": map[string]interface{}{
							"parse": []string{},
						},
						"flags": 1 << 6,
					},
				})
				return
			}
			var userInstance *models.User
			if body.Member != nil {
				socialInstance, _ := FindUsersByDiscordID(body.Member.User.ID)
				userInstance = socialInstance.User
			} else if body.User != nil {
				socialInstance, _ := FindUsersByDiscordID(body.User.ID)
				userInstance = socialInstance.User
			} else {
				response.RenderCustomResponse(c, map[string]interface{}{
					"type": 4,
					"data": map[string]interface{}{
						"tts":     false,
						"content": "Error capturing the message. Please try again.",
						"embeds":  []string{},
						"allowed_mentions": map[string]interface{}{
							"parse": []string{},
						},
						"flags": 1 << 6,
					},
				})
				return
			}
			title := body.Data.Options[0].Value
			collectionID := utils.Uint64(body.Data.Options[1].Value)
			var parentPageID uint64
			if len(body.Data.Options) > 2 {
				parentPageID = utils.Uint64(body.Data.Options[2].Value)
			}
			fmt.Println("title", title, "collectionid", collectionID, "parentPageID", parentPageID)
			parentPublicAccess := "private"
			newCanvasRepo := canvasrepo.NewCanvasRepoPost{
				CollectionID: collectionID,
				Name:         title,
				Position:     1,
			}
			if parentPageID != 0 {
				canvasRepo, _ := queries.App.RepoQuery.GetRepo(map[string]interface{}{"id": parentPageID})
				parentPublicAccess = canvasRepo.DefaultBranch.PublicAccess
				newCanvasRepo.ParentCanvasRepositoryID = parentPageID
			} else {
				collectionInstance, _ := collection.App.Repo.GetCollection(map[string]interface{}{"id": collectionID})
				parentPublicAccess = collectionInstance.PublicAccess
			}
			studio := activeDiscord[0].Studio
			canvasRepo, err := canvasrepo.App.Controller.CreateCanvasRepo(newCanvasRepo, userInstance.ID, studio.ID, *userInstance, parentPublicAccess)
			if err != nil {
				ErrorHandlerMessage(c)
				return
			}
			canvasUrl := notifications.App.Service.GenerateCanvasBranchUrl(canvasRepo.Key, canvasRepo.Name, studio.ID, *canvasRepo.DefaultBranchID)
			response.RenderCustomResponse(c, map[string]interface{}{
				"type": 4,
				"data": map[string]interface{}{
					"tts":     false,
					"content": "New Canvas Created\n" + canvasUrl,
					"embeds":  []string{},
					"allowed_mentions": map[string]interface{}{
						"parse": []string{},
					},
					"flags": 1 << 6,
				},
			})
			return
		} else if body.Data.Type == 1 && body.Data.Name == "bip-search" {
			guildID := body.GuildID
			activeDiscord, err := GetProductIntegrationByDiscordTeamId(guildID)
			if err != nil || len(activeDiscord) == 0 {
				response.RenderCustomResponse(c, map[string]interface{}{
					"type": 4,
					"data": map[string]interface{}{
						"tts":     false,
						"content": "Error capturing the message. Please try again.",
						"embeds":  []string{},
						"allowed_mentions": map[string]interface{}{
							"parse": []string{},
						},
						"flags": 1 << 6,
					},
				})
				return
			}
			var userInstance *models.User
			if body.Member != nil {
				socialInstance, _ := FindUsersByDiscordID(body.Member.User.ID)
				userInstance = socialInstance.User
			} else if body.User != nil {
				socialInstance, _ := FindUsersByDiscordID(body.User.ID)
				userInstance = socialInstance.User
			} else {
				response.RenderCustomResponse(c, map[string]interface{}{
					"type": 4,
					"data": map[string]interface{}{
						"tts":     false,
						"content": "Error capturing the message. Please try again.",
						"embeds":  []string{},
						"allowed_mentions": map[string]interface{}{
							"parse": []string{},
						},
						"flags": 1 << 6,
					},
				})
				return
			}
			studio := activeDiscord[0].Studio
			searchString := utils.Uint64(body.Data.Options[0].Value)
			var canvasRepo *models.CanvasRepository
			if searchString != 0 {
				canvasRepo, _ = queries.App.RepoQuery.GetRepo(map[string]interface{}{"id": searchString})
			}
			if canvasRepo == nil {
				embedMessage, err := BipSearchTreeBuilderHandler(body, activeDiscord, userInstance)
				if err != nil {
					ErrorHandlerMessage(c)
					return
				}
				fmt.Println("BipSearch event Message Multi Canvas returning", embedMessage)
				response.RenderCustomResponse(c, map[string]interface{}{
					"type": 4,
					"data": map[string]interface{}{
						"tts":     false,
						"content": "",
						"embeds":  []discordgo.MessageEmbed{*embedMessage},
						"allowed_mentions": map[string]interface{}{
							"parse": []string{},
						},
						"flags": 1 << 6,
					},
				})
				return
			} else {
				canvasUrl := notifications.App.Service.GenerateCanvasBranchUrl(canvasRepo.Key, canvasRepo.Name, studio.ID, *canvasRepo.DefaultBranchID)
				response.RenderCustomResponse(c, map[string]interface{}{
					"type": 4,
					"data": map[string]interface{}{
						"tts":     false,
						"content": "Canvas found\n" + canvasUrl,
						"embeds":  []string{},
						"allowed_mentions": map[string]interface{}{
							"parse": []string{},
						},
						"flags": 1 << 6,
					},
				})
				return
			}
		}
	}

	if body.Type == InteractionMessageComponent {
		if body.Data.CustomID == "accessrequest" && body.Data.ComponentType == 3 { // message command
			permission := body.Data.Values[0]
			// Get the notification instance by id
			// Notifier ID will be actor or auth User ID.
			// ObjectID will be AccessRequest ID
			// ManageAccessRequest method is called by access requestID, access/reject, permissionGroup

			notification, err := notifications.App.Repo.GetNotification(map[string]interface{}{"discord_dm_id": body.Message.ID})
			if err != nil {
				response.RenderCustomResponse(c, map[string]interface{}{
					"type": 4,
					"data": map[string]interface{}{
						"tts":     false,
						"content": "Error capturing the message. Please try again.",
						"embeds":  []string{},
						"allowed_mentions": map[string]interface{}{
							"parse": []string{},
						},
						"flags": 1 << 6,
					},
				})
				return
			}
			mergeRequestBody := ar.ManageAccessRequestPost{
				Status:                      models.ACCESS_REQUEST_ACCEPTED,
				CanvasBranchPermissionGroup: permission,
			}
			err = ar.App.Service.ManageAccessRequest(*notification.ObjectId, mergeRequestBody, notification.NotifierID)
			if err != nil {
				response.RenderCustomResponse(c, map[string]interface{}{
					"type": 4,
					"data": map[string]interface{}{
						"tts":     false,
						"content": "Error capturing the message. Please try again.",
						"embeds":  []string{},
						"allowed_mentions": map[string]interface{}{
							"parse": []string{},
						},
						"flags": 1 << 6,
					},
				})
				return
			}
			response.RenderCustomResponse(c, map[string]interface{}{
				"type": 6,
				"data": map[string]interface{}{
					"tts":     false,
					"content": "Successfully given access to the canvas",
					"embeds":  []string{},
					"allowed_mentions": map[string]interface{}{
						"parse": []string{},
					},
					"flags": 1 << 6,
				},
			})
			return
		}
		if body.Data.CustomID == "checkCanvasAccess" && body.Data.ComponentType == 2 {
			guildID := body.GuildID
			activeDiscord, err := GetProductIntegrationByDiscordTeamId(guildID)
			if err != nil || len(activeDiscord) == 0 {
				response.RenderCustomResponse(c, map[string]interface{}{
					"type": 4,
					"data": map[string]interface{}{
						"tts":     false,
						"content": "Error capturing the message. Please try again.",
						"embeds":  []string{},
						"allowed_mentions": map[string]interface{}{
							"parse": []string{},
						},
						"flags": 1 << 6,
					},
				})
				return
			}
			studio := activeDiscord[0].Studio
			var userInstance *models.User
			if body.Member != nil {
				socialInstance, _ := FindUsersByDiscordID(body.Member.User.ID)
				userInstance = socialInstance.User
			} else if body.User != nil {
				socialInstance, _ := FindUsersByDiscordID(body.User.ID)
				userInstance = socialInstance.User
			} else {
				response.RenderCustomResponse(c, map[string]interface{}{
					"type": 4,
					"data": map[string]interface{}{
						"tts":     false,
						"content": "Error capturing the message. Please try again.",
						"embeds":  []string{},
						"allowed_mentions": map[string]interface{}{
							"parse": []string{},
						},
						"flags": 1 << 6,
					},
				})
				return
			}
			embed := canvasrepo.App.Service.CreateStudioFileTree(studio, userInstance)
			response.RenderCustomResponse(c, map[string]interface{}{
				"type": 4,
				"data": map[string]interface{}{
					"tts":     false,
					"content": "Canvases List",
					"embeds":  []discordgo.MessageEmbed{*embed},
					"allowed_mentions": map[string]interface{}{
						"parse": []string{},
					},
					"flags": 1 << 6,
				},
			})
			return
		}
		//if body.Data.CustomID == "welcome" && body.Data.ComponentType == 3 { // message command
		//	value := body.Data.Values[0]
		//	fmt.Println("body member==========>", body.User)
		//	users, _ := FindUsersByDiscordIDs([]string{body.User.ID})
		//	user := users[0]
		//	fmt.Println("user===========>", user)
		//	if value == "off" {
		//		err := user.UpdateUserSettings(c.Request.Context(), models.NOTIFICATION_TYPE_DISCORD, false, false, false, false, false, false, false, false, false, false, false)
		//		if err != nil {
		//			response.RenderCustomResponse(c, map[string]interface{}{
		//				"type": 4,
		//				"data": map[string]interface{}{
		//					"tts":     false,
		//					"content": "Failed to update. Please try again!",
		//					"embeds":  []string{},
		//					"allowed_mentions": map[string]interface{}{
		//						"parse": []string{},
		//					},
		//					"flags": 1 << 6,
		//				},
		//			})
		//			return
		//		} else {
		//			response.RenderCustomResponse(c, map[string]interface{}{
		//				"type": 4,
		//				"data": map[string]interface{}{
		//					"tts":     false,
		//					"content": "Successfully updated the settings.",
		//					"embeds":  []string{},
		//					"allowed_mentions": map[string]interface{}{
		//						"parse": []string{},
		//					},
		//					"flags": 1 << 6,
		//				},
		//			})
		//			return
		//		}
		//	} else if value == "on" {
		//		err := user.UpdateWithDefaultUserSettings(c.Request.Context(), models.NOTIFICATION_TYPE_DISCORD)
		//		if err != nil {
		//			response.RenderCustomResponse(c, map[string]interface{}{
		//				"type": 4,
		//				"data": map[string]interface{}{
		//					"tts":     false,
		//					"content": "Failed to update. Please try again!",
		//					"embeds":  []string{},
		//					"allowed_mentions": map[string]interface{}{
		//						"parse": []string{},
		//					},
		//					"flags": 1 << 6,
		//				},
		//			})
		//			return
		//		} else {
		//			response.RenderCustomResponse(c, map[string]interface{}{
		//				"type": 4,
		//				"data": map[string]interface{}{
		//					"tts":     false,
		//					"content": "Successfully updated the settings.",
		//					"embeds":  []string{},
		//					"allowed_mentions": map[string]interface{}{
		//						"parse": []string{},
		//					},
		//					"flags": 1 << 6,
		//				},
		//			})
		//			return
		//		}
		//	}
		//
		//	response.RenderCustomResponse(c, map[string]interface{}{
		//		"type": 4,
		//		"data": map[string]interface{}{
		//			"tts":     false,
		//			"content": "Failed to update.Please try again! ",
		//			"embeds":  []string{},
		//			"allowed_mentions": map[string]interface{}{
		//				"parse": []string{},
		//			},
		//			"flags": 1 << 6,
		//		},
		//	})
		//	return
		//}

	}

	response.RenderCustomResponse(c, map[string]interface{}{
		"success": false,
	})
}
