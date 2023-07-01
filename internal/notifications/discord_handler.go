package notifications

import (
	"encoding/json"
	"fmt"
	"gorm.io/datatypes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	studiointegration "gitlab.com/phonepost/bip-be-platform/internal/studio_integration"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/integrations"
)

func (s notificationService) SendDiscordWelcomeMsg(userDiscordId string, teamID string) {
	guild, _ := integrations.GetDiscordTeam(teamID)
	resp := []string{"Hey there!"}
	resp = append(resp, "")

	resp = append(resp, fmt.Sprintf(" `%s` server was integrated with bip.so.. Relevant notifications (when you are invited to a document, mentioned etc. on bip.so) will be sent to you as a DM here for your convenience. You can change switch this off anytime using the dropdown below!", guild.Name))

	resp = append(resp, "")

	resp = append(resp, "Happy collaborating!")

	welcomecomponents := []interface{}{
		integrations.ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				discordgo.SelectMenu{
					CustomID:    "welcome",
					Placeholder: "DM Notifications Enabled",
					MaxValues:   1,
					Options: []discordgo.SelectMenuOption{
						{
							Label:       "Enable",
							Value:       "on",
							Description: "",
						},
						{
							Label:       "Disable",
							Value:       "off",
							Description: "",
						},
					},
				},
			},
		},
	}
	_, err := integrations.SendDiscordUserComponentsDM(userDiscordId, resp, welcomecomponents)
	if err != nil {
		println("error while sending welcome message ", err.Error())
	}
}

func (s notificationService) SendDiscordComponentsMsg(notification *models.Notification, components []interface{}, userDiscordId string, discordMessage []string) {
	msg, err := integrations.SendDiscordUserComponentsDM(userDiscordId, discordMessage, components)
	if err != nil {
		println("error while sending message ", err.Error())
		return
	}
	fmt.Println("Send Discord component msg successfully", msg, msg.ID)
	App.Repo.UpdateNotifications(map[string]interface{}{"id": notification.ID}, map[string]interface{}{"discord_dm_id": msg.ID})
}

func (s notificationService) SendDiscordDmComponentsMsg(notification *models.Notification, userDiscordId string, discordMessage []string) {
	msgID, err := integrations.SendDiscordUserDM(userDiscordId, discordMessage)
	if err != nil {
		println("error while sending message ", err.Error())
		return
	}
	fmt.Println("Send Discord component msg successfully", msgID)
	App.Repo.UpdateNotifications(map[string]interface{}{"id": notification.ID}, map[string]interface{}{"discord_dm_id": msgID})
}

func (s notificationService) EditDiscordMsg(notification models.Notification, userDiscordId string, notificationEditText []string) {
	msgID, err := integrations.EditDiscordUserDM(userDiscordId, *notification.DiscordDmID, notificationEditText)
	if err != nil {
		println("error while sending message ", err.Error())
	}
	App.Repo.UpdateNotifications(map[string]interface{}{"id": notification.ID}, map[string]interface{}{"discord_dm_id": msgID})
}

func (s notificationService) DeleteDiscordMsg(userDiscordId string, discordMsgId string) {
	err := integrations.DeleteDiscordMessage(userDiscordId, discordMsgId)
	if err != nil {
		println("error while sending message ", err.Error())
	}
}

func (s notificationService) GetDiscordStudioIntegration(studioId uint64) (integration *models.StudioIntegration, err error) {
	err = postgres.GetDB().Model(&models.StudioIntegration{}).Where("studio_id = ? and type = ?", studioId, studiointegration.DISCORD_INTEGRATION_TYPE).First(&integration).Error
	return
}

func (s notificationService) GetStudioIntegration(studioId uint64, integrationType string) (integration *models.StudioIntegration, err error) {
	err = postgres.GetDB().Model(&models.StudioIntegration{}).Where("studio_id = ? and type = ?", studioId, integrationType).First(&integration).Error
	return
}

func (s notificationService) StudioDiscordIntegration(studioId uint64) (*models.StudioIntegration, bool) {
	studio, _ := App.Repo.GetStudioByID(studioId)
	studioIntegration, err := s.GetDiscordStudioIntegration(studioId)
	if err != nil {
		if err.Error() == "record not found" {
			fmt.Println("studio is not integrated")
			return nil, true
		} else {
			return nil, false
		}
	} else {
		if studioIntegration != nil {
			fmt.Println("Studio has discord integration", studioIntegration)
			if !studio.DiscordNotificationsEnabled {
				fmt.Println("Discord notification not allowed in product, returning..")
				return studioIntegration, false
			}
		}
	}
	return studioIntegration, true
}

func (s notificationService) DiscordEventHandler(reel *models.Reel, integration *models.StudioIntegration) {
	extra := integration.Extra
	if len(extra) < 1 {
		return
	}
	discordData := map[string]interface{}{}
	err := json.Unmarshal(extra, &discordData)
	if err != nil {
		fmt.Println("Unmarshal extra data error:", err)
		return
	}
	contextData := map[string]interface{}{}
	err = json.Unmarshal(reel.ContextData, &contextData)
	if err != nil {
		fmt.Println("Unmarshal context data error:", err)
		return
	}

	reelText := contextData["text"].(string)

	selectedBlocks := map[string]map[string]BlockData{}
	err = json.Unmarshal(reel.SelectedBlocks, &selectedBlocks)
	blocksData := selectedBlocks["blocksData"]

	text := "**" + reel.CreatedByUser.Username + "**: " + reelText + "\n\n>>> "
	for _, block := range blocksData {
		if block.Type == models.BlockSimpleTableV1 {
			continue
		}
		format := "%s \n"
		if block.Type == models.BlockTypeHeading1 {
			for _, children := range block.Children {
				text = fmt.Sprintf(format+"%s", text, "**#"+children.Text+"**")
			}
		} else if block.Type == models.BlockTypeHeading2 {
			for _, children := range block.Children {
				text = fmt.Sprintf(format+"%s", text, "**##"+children.Text+"**")
			}
		} else if block.Type == models.BlockTypeHeading3 {
			for _, children := range block.Children {
				text = fmt.Sprintf(format+"%s", text, "**###"+children.Text+"**")
			}
		} else if block.Type == models.BlockTypeHeading4 {
			for _, children := range block.Children {
				text = fmt.Sprintf(format+"%s", text, "**####"+children.Text+"*")
			}
		} else if block.Type == models.BlockTypeHeading5 {
			for _, children := range block.Children {
				text = fmt.Sprintf(format+"%s", text, "**#####"+children.Text+"**")
			}
		} else if block.Type == models.BlockTypeHeading6 {
			for _, children := range block.Children {
				text = fmt.Sprintf(format+"%s", text, "**######"+children.Text+"**")
			}
		} else if block.Type == models.BlockTypeImage {
			for _, children := range block.Children {
				if children.Type == models.BlockTypeImage {
					text = fmt.Sprintf(format+"Image: %s", text, children.Url)
				}
			}
		} else if block.Type == models.BlockTypeAttachment {
			for _, children := range block.Children {
				if children.Type == models.BlockTypeAttachment {
					text = fmt.Sprintf(format+"File: %s", text, children.Url)
				}
			}
		} else if utils.Contains([]string{models.BlockTypeVideo, models.BlockTypeTweet}, block.Type) {
			for _, children := range block.Children {
				text = fmt.Sprintf(format+"Embed: %s", text, children.Url)
			}
		} else {
			for _, children := range block.Children {
				if children.Type == models.BlockTypeImage {
					text = fmt.Sprintf(format+"Image: %s", text, children.Url)
				} else if children.Type == models.BlockTypeAttachment {
					text = fmt.Sprintf(format+"File: %s", text, children.Url)
				} else if children.Type == models.BlockDrawIO {
					continue
				} else if children.Url != "" {
					text = fmt.Sprintf(format+"Embed: %s", text, children.Url)
				} else {
					text = fmt.Sprintf(format+"%s", text, children.Text)
				}
			}
		}
	}

	canvasRepo, _ := App.Repo.GetCanvasRepoByID(reel.CanvasRepositoryID)
	reelURL := configs.GetAppInfoConfig().FrontendHost
	if canvasRepo.ID != 0 {
		reelURL = App.Service.GenerateReelUUIDUrl(canvasRepo.Key, canvasRepo.Name, reel.StudioID, reel.CanvasBranchID, reel.UUID.String())
	}
	text = fmt.Sprintf("%s\n\n%s", text, reelURL)

	messageData := map[string]string{}
	postsChannelID := ""
	if integration.MessagesData != nil {
		json.Unmarshal(*integration.MessagesData, &messageData)
		postsChannelID = messageData["postsChannelId"]
		if postsChannelID == "" {
			postsChannelID = s.CreateBipPostsChannel(integration.TeamID, integration.ID)
		}
	}

	if postsChannelID != "" {
		msg, err := integrations.SendDiscordDMMessageToChannel(postsChannelID, []string{text}, []interface{}{})
		if err != nil {
			fmt.Println("Error in sending reel message to discord", msg)
			return
		}
		_result, _ := json.Marshal(msg)
		err = App.Repo.AddIntegrationReference(msg.ID, models.DISCORD_INTEGRATION_TYPE, reel.ID, models.REEL, string(_result))
		if err != nil {
			fmt.Println("Error on adding integration reference", err)
		}
	} else {
		webhookData := discordData["webhook"].(map[string]interface{})
		endpoint := webhookData["url"].(string)
		contentType := "application/x-www-form-urlencoded"
		values := url.Values{
			"content": {text},
		}
		reqBody := strings.NewReader(values.Encode())
		_response, err := http.Post(endpoint+"?wait=true", contentType, reqBody)
		if err != nil {
			fmt.Println("error on sending reel event to discord", err, endpoint)
			return
		}
		_result, _ := ioutil.ReadAll(_response.Body)
		responseData := map[string]interface{}{}
		err = json.Unmarshal(_result, &discordData)
		if err != nil || responseData == nil {
			println("error while structToMap", err)
			return
		}
		id := discordData["id"].(string)
		err = App.Repo.AddIntegrationReference(id, models.DISCORD_INTEGRATION_TYPE, reel.ID, models.REEL, string(_result))
		if err != nil {
			fmt.Println("Error on adding integration reference", err)
		}
	}

}

type PostChildrenBlocks struct {
	Type     string `json:"type"`
	Children []struct {
		Text string `json:"text"`
	} `json:"children"`
}

func (s notificationService) DiscordNewPostEventHandler(post *models.Post, integration *models.StudioIntegration) {
	extra := integration.Extra
	if len(extra) < 1 {
		return
	}
	discordData := map[string]interface{}{}
	err := json.Unmarshal(extra, &discordData)
	if err != nil {
		fmt.Println("Unmarshal extra data error:", err)
		return
	}
	contextData := map[string]interface{}{}
	err = json.Unmarshal(post.Children, &contextData)
	if err != nil {
		fmt.Println("Unmarshal context data error:", err)
		return
	}

	postBlocks := contextData["blocks"].(string)
	postBlocksData := []PostChildrenBlocks{}
	var postText string
	json.Unmarshal([]byte(postBlocks), &postBlocksData)
	for _, child := range postBlocksData {
		for _, message := range child.Children {
			postText += message.Text + "\n"
		}
	}
	text := "**" + post.CreatedByUser.Username + "**: " + postText
	// /feed?postId=16
	studioInstance, _ := App.Repo.GetStudioByID(post.StudioID)
	ENVMAILER := configs.GetConfigString("ENV")
	postLink := models.MailerRouterPaths[ENVMAILER]["BASE_URL"] + studioInstance.Handle + "/feed?postId=" + utils.String(post.ID)
	text = fmt.Sprintf("%s\n> %s", text, postLink)

	messageData := map[string]string{}
	postsChannelID := ""
	if integration.MessagesData != nil {
		json.Unmarshal(*integration.MessagesData, &messageData)
		postsChannelID = messageData["postsChannelId"]
		if postsChannelID == "" {
			postsChannelID = s.CreateBipPostsChannel(integration.TeamID, integration.ID)
		}
	}
	if postsChannelID != "" {
		msg, err := integrations.SendDiscordDMMessageToChannel(postsChannelID, []string{text}, []interface{}{})
		if err != nil {
			fmt.Println("Error in sending post message to discord", msg)
			return
		}
		_result, _ := json.Marshal(msg)
		err = App.Repo.AddIntegrationReference(msg.ID, models.DISCORD_INTEGRATION_TYPE, post.ID, models.REEL, string(_result))
		if err != nil {
			fmt.Println("Error on adding integration reference", err)
		}
	} else {
		webhookData := discordData["webhook"].(map[string]interface{})
		endpoint := webhookData["url"].(string)
		contentType := "application/x-www-form-urlencoded"
		values := url.Values{
			"content": {text},
		}
		reqBody := strings.NewReader(values.Encode())
		_response, err := http.Post(endpoint+"?wait=true", contentType, reqBody)
		if err != nil {
			fmt.Println("error on sending post event to discord", err, endpoint)
			return
		}
		_result, _ := ioutil.ReadAll(_response.Body)
		responseData := map[string]interface{}{}
		err = json.Unmarshal(_result, &discordData)
		if err != nil || responseData == nil {
			println("error while structToMap", err)
			return
		}
		id := discordData["id"].(string)
		err = App.Repo.AddIntegrationReference(id, models.DISCORD_INTEGRATION_TYPE, post.ID, models.POST, string(_result))
		if err != nil {
			fmt.Println("Error on adding integration reference", err)
		}
	}
}

func (s notificationService) CreateBipPostsChannel(guildID string, integrationID uint64) string {
	channel, err := integrations.CreateBipPostsCanvasChannel(guildID)
	if err != nil {
		fmt.Println("Error in creating bip posts channel", err)
		return ""
	}
	studioIntegration, _ := App.Repo.GetStudioIntegrationByID(integrationID)
	discordMessagesData := map[string]interface{}{}
	json.Unmarshal(*studioIntegration.MessagesData, &discordMessagesData)
	discordMessagesData["postsChannelId"] = channel.ID
	messagesDataStr, _ := json.Marshal(discordMessagesData)
	App.Repo.StudioIntegrationUpdate(studioIntegration.ID, map[string]interface{}{"messages_data": datatypes.JSON(messagesDataStr)})
	return channel.ID
}
