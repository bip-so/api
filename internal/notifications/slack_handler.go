package notifications

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/integrations"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"io/ioutil"
	"net/http"
	"strings"
)

func (s notificationService) sendSlackNotification(notification *models.Notification, components []map[string]interface{}, notifierId uint64, slackMessage string) {
	var studioIntegration *models.StudioIntegration
	var canSendNotification bool
	if notification.StudioID != nil {
		studioIntegration, canSendNotification = s.GetStudioIntegrationWithStatus(*notification.StudioID)
		if !canSendNotification {
			return
		}
	}

	userSocialAuth, _ := App.Repo.GetUserSocialAuth(notifierId, models.SLACK_PROVIDER)
	if userSocialAuth != nil {
		isWelcomeMessage := App.Repo.IsUserSlackNotificationSentAnytime(notifierId)
		if isWelcomeMessage && studioIntegration != nil {
			s.SlackPostWelcomeMessage(studioIntegration.AccessKey, studioIntegration.TeamID, userSocialAuth.ProviderID, studioIntegration.StudioID, userSocialAuth.UserID)
		}
		var slackAccessKey string
		if studioIntegration != nil {
			slackAccessKey = studioIntegration.AccessKey
		} else {
			var metadata *SlackSocialAuthMetadata
			json.Unmarshal(userSocialAuth.Metadata, &metadata)
			if metadata != nil && metadata.AccessToken == "" {
				return
			}
			slackAccessKey = metadata.AccessToken
		}
		s.SlackPostMessage(notification, components, slackAccessKey, userSocialAuth.ProviderID, slackMessage)
	}
}

func (s notificationService) SlackPostWelcomeMessage(accessToken, teamID, slackUserID string, studioID, userID uint64) {
	slackTeam, err := integrations.GetSlackTeam(accessToken, teamID)
	studio, err := App.Repo.GetStudioByID(studioID)
	bipUser, err := App.Repo.GetUser(userID)
	teamName := ""
	if err != nil {
		fmt.Println("Error in getting slack Team", err)
	}
	if slackTeam != nil {
		teamName = slackTeam.Team.Name
	}
	blocks := []map[string]interface{}{
		{
			"type": "section",
			"text": map[string]interface{}{
				"type": "plain_text",
				"text": fmt.Sprintf(
					"Hey %s!\n\n%s workspace is integrated with %s.\n\nYou will now receive DMs of relevant notifications\n-Canvas invites\n-Merge requests\n-Mentions\n\nHappy collaborating!", bipUser.FullName, teamName, studio.DisplayName),
				"emoji": true,
			},
		},
		{
			"type": "section",
			"text": map[string]interface{}{
				"type": "mrkdwn",
				"text": "DM Notifications",
			},
			"accessory": map[string]interface{}{
				"type": "static_select",
				"placeholder": map[string]interface{}{
					"type":  "plain_text",
					"text":  "DM Notifications",
					"emoji": true,
				},
				"options": []map[string]interface{}{
					{
						"text": map[string]interface{}{
							"type":  "plain_text",
							"text":  "Enable",
							"emoji": true,
						},
						"value": "enable",
					},
					{
						"text": map[string]interface{}{
							"type":  "plain_text",
							"text":  "Disable",
							"emoji": true,
						},
						"value": "disable",
					},
				},
				"action_id": "dmNotifications",
			},
		},
	}
	_, err = integrations.SendMessageToSlackChannel(context.Background(), "", accessToken, "", slackUserID, blocks)
	if err != nil {
		fmt.Println("Error in sending slack messages", err)
	}
}

func (s notificationService) SlackPostMessage(notification *models.Notification, slackMessage []map[string]interface{}, accessToken, slackUserID, slackTextMessage string) {
	slackResponse, err := integrations.SendMessageToSlackChannel(context.Background(), slackTextMessage, accessToken, "", slackUserID, slackMessage)
	if err != nil {
		fmt.Println("Error in sending slack messages", err)
		return
	}
	fmt.Println("Send slack component msg successfully")
	App.Repo.UpdateNotifications(map[string]interface{}{"id": notification.ID}, map[string]interface{}{"slack_dm_id": slackResponse.Ts, "slack_channel_id": slackResponse.Channel})
}

func (s notificationService) SlackNotificationBlockBuilder(messageText []string, buttonText, url string) []map[string]interface{} {
	text := ""
	for _, message := range messageText {
		message = strings.ReplaceAll(message, "**", "*")
		text += message + "\n"
	}
	blocks := []map[string]interface{}{
		{
			"type": "section",
			"text": map[string]interface{}{
				"type": "mrkdwn",
				"text": text,
			},
		},
		{
			"type": "actions",
			"elements": []map[string]interface{}{
				{
					"type": "button",
					"text": map[string]interface{}{
						"type":  "plain_text",
						"text":  buttonText,
						"emoji": true,
					},
					"url": url,
				},
			},
		},
	}
	return blocks
}

func (s notificationService) SlackAccessRequestMessageBuilder(messageText []string, buttonText, appUrl string) []map[string]interface{} {
	text := ""
	for _, message := range messageText {
		text += message + "\n"
	}
	blocks := []map[string]interface{}{
		{
			"type": "section",
			"text": map[string]interface{}{
				"type":  "plain_text",
				"text":  text,
				"emoji": true,
			},
		},
		{
			"type": "section",
			"text": map[string]interface{}{
				"type": "mrkdwn",
				"text": "Grant Permission",
			},
			"accessory": map[string]interface{}{
				"type": "static_select",
				"placeholder": map[string]interface{}{
					"type":  "plain_text",
					"text":  "Select Permission",
					"emoji": true,
				},
				"options": []map[string]interface{}{
					{
						"text": map[string]interface{}{
							"type":  "plain_text",
							"text":  "Moderate",
							"emoji": true,
						},
						"value": models.PGCanvasModerateSysName,
					},
					{
						"text": map[string]interface{}{
							"type":  "plain_text",
							"text":  "Edit",
							"emoji": true,
						},
						"value": models.PGCanvasEditSysName,
					},
					{
						"text": map[string]interface{}{
							"type":  "plain_text",
							"text":  "Reply",
							"emoji": true,
						},
						"value": models.PGCanvasCommentSysName,
					},
					{
						"text": map[string]interface{}{
							"type":  "plain_text",
							"text":  "View",
							"emoji": true,
						},
						"value": models.PGCanvasViewSysName,
					},
				},
				"action_id": "grantPermission",
			},
		},
		{
			"type": "actions",
			"elements": []map[string]interface{}{
				{
					"type": "button",
					"text": map[string]interface{}{
						"type":  "plain_text",
						"text":  buttonText,
						"emoji": true,
					},
					"url": appUrl,
				},
			},
		},
	}
	return blocks
}

func (s notificationService) SlackReelEventHandler(reel *models.Reel, integration *models.StudioIntegration) {
	extra := integration.Extra
	if len(extra) < 1 {
		return
	}
	slackData := map[string]interface{}{}
	err := json.Unmarshal(extra, &slackData)
	if err != nil {
		fmt.Println("Unmarshal extra data error:", err)
		return
	}
	webhookData := slackData["incoming_webhook"].(map[string]interface{})
	endpoint := webhookData["url"].(string)
	contentType := "application/json"

	messageText := s.GetReelBlocksText(reel)
	data := map[string]interface{}{
		"blocks": []map[string]interface{}{
			{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": messageText,
				},
			},
		},
	}
	valuesStr, _ := json.Marshal(data)
	payload := strings.NewReader(string(valuesStr))
	response, err := http.Post(endpoint, contentType, payload)
	if err != nil {
		fmt.Println("error on sending reel event to slack", err, endpoint)
		return
	}
	_result, _ := ioutil.ReadAll(response.Body)
	fmt.Println("Send reel message", string(_result))
}

func (s notificationService) SlackPostEventHandler(post *models.Post, integration *models.StudioIntegration) {
	extra := integration.Extra
	if len(extra) < 1 {
		return
	}
	slackData := map[string]interface{}{}
	err := json.Unmarshal(extra, &slackData)
	if err != nil {
		fmt.Println("Unmarshal extra data error:", err)
		return
	}
	webhookData := slackData["incoming_webhook"].(map[string]interface{})
	endpoint := webhookData["url"].(string)
	contentType := "application/json"
	contextData := map[string]interface{}{}
	err = json.Unmarshal(post.Children, &contextData)
	if err != nil {
		fmt.Println("Unmarshal context data error:", err)
		return
	}
	var messageText string
	postBlocks := contextData["blocks"].(string)
	postBlocksData := []PostChildrenBlocks{}
	json.Unmarshal([]byte(postBlocks), &postBlocksData)
	for _, child := range postBlocksData {
		for _, message := range child.Children {
			messageText += message.Text + "\n"
		}
	}
	ENVMAILER := configs.GetConfigString("ENV")
	studio, _ := App.Repo.GetStudioByID(post.StudioID)
	postLink := models.MailerRouterPaths[ENVMAILER]["BASE_URL"] + studio.Handle + "/feed?postId=" + utils.String(post.ID)
	messageText = fmt.Sprintf("%s\n\n%s", messageText, postLink)
	data := map[string]interface{}{
		"blocks": []map[string]interface{}{
			{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": messageText,
				},
			},
		},
	}
	valuesStr, _ := json.Marshal(data)
	payload := strings.NewReader(string(valuesStr))
	response, err := http.Post(endpoint, contentType, payload)
	if err != nil {
		fmt.Println("error on sending reel event to slack", err, endpoint)
		return
	}
	_result, _ := ioutil.ReadAll(response.Body)
	fmt.Println("Send reel message", string(_result))
}
