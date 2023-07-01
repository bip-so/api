package notifications

import (
	"encoding/json"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/pkg/integrations"
	"strings"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

// Get the notification settings of the user.
// based on the settings send the notifications.
func (s notificationService) handleNotificationCreation(notification *PostNotification) {
	var userIDs []uint64
	// If roleIDs is present. Fetching the UserIDs of members present in the role.
	if notification.RoleIDs != nil {
		roles, err := App.Repo.GetRolesByID(notification.RoleIDs)
		if err != nil {
			logger.Error(err.Error())
		}
		for _, role := range roles {
			for _, member := range role.Members {
				userIDs = append(userIDs, member.UserID)
			}
		}
	}
	notification.NotifierIDs = append(notification.NotifierIDs, userIDs...)

	extraData := s.GenerateExtraData(notification)
	notifierEmailIDs := []uint64{}
	for _, notifierID := range notification.NotifierIDs {
		if notifierID == notification.CreatedByID && notification.Event != BipMarkMessageAdded {
			continue
		}
		notificationInstance := s.NewNotification(notification, notifierID, extraData)
		err := App.Repo.AddNotificationToDbAndStream(notificationInstance)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		access, _ := s.CalculateNotificationEntityAccess(notificationInstance.NotifierID, notification.Entity)
		if access["app"] {
			// Increase the notification_count table
			s.updateNotificationCountTable(notificationInstance)
		}
		if access["email"] {
			// Send email
			notifierEmailIDs = append(notifierEmailIDs, notifierID)
		}
		if access["discord"] {
			// Send Discord notification
			s.sendDiscordNotification(notificationInstance, notification.ExtraData.DiscordComponents, notifierID, notification.ExtraData.DiscordMessage)
		}
		if access["slack"] {
			s.sendSlackNotification(notificationInstance, notification.ExtraData.SlackComponents, notifierID, notification.ExtraData.SlackMessage)
		}
	}

	// Send bulk emails to all the notifierIDs
	s.sendEmailNotification(notifierEmailIDs, &NotificationEmailData{
		Text:     notification.Text,
		Activity: notification.Activity,
		Subject:  notification.ExtraData.EmailSubject,
		AppUrl:   notification.ExtraData.AppUrl,
		Event:    notification.Event,
	})
}

func (s notificationService) NewNotification(notification *PostNotification, notifierID uint64, extraData []byte) *models.Notification {
	return &models.Notification{
		Entity:              notification.Entity,
		Event:               notification.Event,
		Activity:            notification.Activity,
		Text:                notification.Text,
		Priority:            notification.Priority,
		CreatedByID:         notification.CreatedByID,
		StudioID:            notification.StudioID,
		NotifierID:          notifierID,
		ExtraData:           extraData,
		ObjectId:            notification.ObjectID,
		ContentObject:       notification.ContentObject,
		TargetObjectId:      notification.TargetObjectID,
		TargetContentObject: notification.TargetContentObject,
		IsPersonal:          notification.IsPersonal,
		Version:             2,
	}
}

// This functions may be triggered from go routines or from kafka consumer
func (s notificationService) sendEmailNotification(notifierEmailIDs []uint64, data *NotificationEmailData) {
	// Sending emails if notifiers are present
	if len(notifierEmailIDs) > 0 {
		users, err := App.Repo.GetUserEmailsByIDs(notifierEmailIDs)
		if err != nil {
			logger.Error(err.Error())
		}
		toEmails := []string{}
		for _, user := range users {
			if user.Email.Valid {
				toEmails = append(toEmails, user.Email.String)
			}
		}
		mailer := pkg.BipMailer{}
		template := s.GenerateEmailTemplate(data.Text, data.AppUrl, data.Event)
		if data.Subject == "" {
			data.Subject = data.Activity
		}
		err = mailer.SendEmail(toEmails, nil, nil, data.Subject, template, template)
		if err != nil {
			logger.Error(err.Error())
		}
	}
}

func (s notificationService) sendDiscordNotification(notification *models.Notification, components []interface{}, notifierId uint64, discordMessage []string) {
	var studioIntegration *models.StudioIntegration
	var canSendNotification bool
	if notification.StudioID != nil {
		studioIntegration, canSendNotification = s.StudioDiscordIntegration(*notification.StudioID)
		if !canSendNotification {
			return
		}
	}

	userSocialAuth, _ := App.Repo.GetUserSocialAuth(notifierId, models.DISCORD_PROVIDER)
	if userSocialAuth != nil {
		isWelcomeMessage := App.Repo.IsUserDiscordNotificationSentAnytime(notifierId)
		if isWelcomeMessage && studioIntegration != nil {
			s.SendDiscordWelcomeMsg(userSocialAuth.ProviderID, studioIntegration.TeamID)
		}
		s.SendDiscordComponentsMsg(notification, components, userSocialAuth.ProviderID, discordMessage)
	}
}

func (s notificationService) sendDiscordDmNotification(notification *models.Notification, discordMessage []string) {
	var studioIntegration *models.StudioIntegration
	var canSendNotification bool
	if notification.StudioID != nil {
		studioIntegration, canSendNotification = s.StudioDiscordIntegration(*notification.StudioID)
		if !canSendNotification {
			return
		}
	}

	userSocialAuth, _ := App.Repo.GetUserSocialAuth(notification.NotifierID, models.DISCORD_PROVIDER)
	if userSocialAuth != nil {
		isWelcomeMessage := App.Repo.IsUserDiscordNotificationSentAnytime(notification.NotifierID)
		if isWelcomeMessage && studioIntegration != nil {
			s.SendDiscordWelcomeMsg(userSocialAuth.ProviderID, studioIntegration.TeamID)
		}
		s.SendDiscordDmComponentsMsg(notification, userSocialAuth.ProviderID, discordMessage)
	}
}

func (s notificationService) updateDiscordNotification(notification models.Notification, notificationEditText []string) {
	// Accept generic terms and update the discord notification.
	userSocialAuth, _ := App.Repo.GetUserSocialAuth(notification.NotifierID, models.DISCORD_PROVIDER)
	if userSocialAuth != nil {
		s.EditDiscordMsg(notification, userSocialAuth.ProviderID, notificationEditText)
	}
}

func (s notificationService) deleteDiscordNotification(notification models.Notification) {
	// Accept generic terms and update the discord notification.
	userSocialAuth, _ := App.Repo.GetUserSocialAuth(notification.NotifierID, models.DISCORD_PROVIDER)
	if userSocialAuth != nil {
		s.DeleteDiscordMsg(userSocialAuth.ProviderID, *notification.DiscordDmID)
	}
}

func (s notificationService) deleteSlackNotification(notification models.Notification) {
	// Accept generic terms and update the discord notification.
	userSocialAuth, _ := App.Repo.GetUserSocialAuth(notification.NotifierID, models.SLACK_PROVIDER)
	if userSocialAuth != nil {
		var metadata *SlackSocialAuthMetadata
		json.Unmarshal(userSocialAuth.Metadata, &metadata)
		if metadata != nil && metadata.AccessToken != "" {
			integrations.DeleteSlackMessage(metadata.AccessToken, *notification.SlackChannelID, *notification.SlackDmID)
			return
		} else {
			if notification.StudioID != nil {
				studioIntegration, _ := s.GetStudioIntegration(*notification.StudioID, models.SLACK_INTEGRATION_TYPE)
				if studioIntegration != nil {
					integrations.DeleteSlackMessage(studioIntegration.AccessKey, *notification.SlackChannelID, *notification.SlackDmID)
					return
				}
			}
		}
	}
}

// GenerateExtraData based on PostNotification
func (s notificationService) GenerateExtraData(data *PostNotification) []byte {
	extraData := map[string]interface{}{}
	if data.CreatedByID != 0 {
		user, _ := App.Repo.GetUser(data.CreatedByID)
		extraData["user"] = map[string]string{
			"id":        utils.String(user.ID),
			"name":      user.FullName,
			"avatarUrl": user.AvatarUrl,
			"username":  user.Username,
		}
	}
	if data.ReactorID != nil {
		reactor, _ := App.Repo.GetUser(*data.ReactorID)
		extraData["reactor"] = map[string]string{
			"id":        utils.String(reactor.ID),
			"name":      reactor.FullName,
			"avatarUrl": reactor.AvatarUrl,
			"username":  reactor.Username,
		}
	}
	if data.StudioID != nil {
		studio, _ := App.Repo.GetStudioByID(*data.StudioID)
		extraData["studio"] = map[string]string{
			"id":          utils.String(studio.ID),
			"displayName": studio.DisplayName,
			"handle":      studio.Handle,
			"imageUrl":    studio.ImageURL,
		}
	}
	if data.ExtraData.CanvasRepoID != 0 {
		canvasRepo, _ := App.Repo.GetCanvasRepoByID(data.ExtraData.CanvasRepoID)
		extraData["canvasRepo"] = map[string]string{
			"id":   utils.String(canvasRepo.ID),
			"name": canvasRepo.Name,
			"icon": canvasRepo.Icon,
			"key":  canvasRepo.Key,
		}
	}
	if data.ExtraData.CollectionID != 0 {
		collection, _ := App.Repo.GetCollectionByID(data.ExtraData.CollectionID)
		extraData["collection"] = map[string]string{
			"id":   utils.String(collection.ID),
			"name": collection.Name,
			"icon": collection.Icon,
		}
	}
	if data.ExtraData.CanvasBranchID != 0 {
		canvasBranch, _ := App.Repo.GetCanvasBranchByID(data.ExtraData.CanvasBranchID)
		extraData["canvasBranch"] = map[string]string{
			"id":   utils.String(canvasBranch.ID),
			"name": canvasBranch.Name,
			"key":  canvasBranch.Key,
		}
	}
	if data.ExtraData.Status != "" {
		extraData["actionStatus"] = data.ExtraData.Status
	}
	extraData["discordComponents"] = data.ExtraData.DiscordComponents
	extraData["slackComponents"] = data.ExtraData.SlackComponents
	extraData["discordMessage"] = data.ExtraData.DiscordMessage
	extraData["permissionGroup"] = data.ExtraData.PermissionGroup
	extraData["message"] = data.ExtraData.Message
	extraData["appUrl"] = data.ExtraData.AppUrl
	extraData["actionOnText"] = data.ExtraData.ActionOnText
	extraData["blockUUID"] = data.ExtraData.BlockUUID
	extraData["blockThreadUUID"] = data.ExtraData.BlockThreadUUID
	extraData["blockThreadCommentUUID"] = data.ExtraData.BlockThreadCommentUUID
	extraData["reelUUID"] = data.ExtraData.ReelUUID
	extraData["reelId"] = data.ExtraData.ReelID
	extraData["reelCommentUUID"] = data.ExtraData.ReelCommentUUID
	parsedData, _ := json.Marshal(extraData)
	return parsedData
}

func (s notificationService) GetUniqueIDs(userIDs []uint64) []uint64 {
	keys := make(map[uint64]bool)
	list := []uint64{}

	for _, userID := range userIDs {
		if !keys[userID] {
			keys[userID] = true
			list = append(list, userID)
		}
	}
	return list
}

func (s notificationService) CreateNewNotification(notificationInstance *models.Notification) {
	var extraData NotificationExtraData
	json.Unmarshal(notificationInstance.ExtraData, &extraData)
	notifierEmailIDs := []uint64{}
	notificationInstance.ID = 0
	err := App.Repo.AddNotificationToDbAndStream(notificationInstance)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	access, _ := s.CalculateNotificationEntityAccess(notificationInstance.NotifierID, notificationInstance.Entity)
	if access["app"] {
		// Increase the notification_count table
		s.updateNotificationCountTable(notificationInstance)
	}
	if access["email"] {
		// Send email
		notifierEmailIDs = append(notifierEmailIDs, notificationInstance.NotifierID)
	}
	if access["discord"] {
		// Send Discord notification
		s.sendDiscordNotification(notificationInstance, extraData.DiscordComponents, notificationInstance.NotifierID, extraData.DiscordMessage)
	}
	if access["slack"] {
		// Send Discord notification
		s.sendSlackNotification(notificationInstance, extraData.SlackComponents, notificationInstance.NotifierID, extraData.SlackMessage)
	}

	s.sendEmailNotification(notifierEmailIDs, &NotificationEmailData{
		Text:     notificationInstance.Text,
		Activity: notificationInstance.Activity,
		AppUrl:   extraData.AppUrl,
	})
}

func (s notificationService) GetDBNotificationForUser(user *models.User, skip int, limit int, notificationType string, studioID uint64, filter string) ([]models.Notification, error) {
	query := map[string]interface{}{"notifier_id": user.ID}
	if notificationType == "studio" {
		query["studio_id"] = studioID
	}
	if notificationType == "personal" {
		query["is_personal"] = true
	}
	filters := strings.Split(filter, ",")
	events := []string{}
	for _, key := range filters {
		if key == "unread" {
			query["seen"] = false
		}
		if key == "requests" {
			events = append(events, REQUEST_EVENTS...)
		}
		if key == "replies" {
			events = append(events, REPLIES_EVENTS...)
		}
		if key == "pr" {
			events = append(events, PR_EVENTS...)
		}
	}
	if filter != "" {
		query["event"] = events
	}
	notifications, err := App.Repo.GetNotifications(query, skip, limit)
	if err != nil {
		fmt.Println("Error in getting notifications", err)
		return nil, err
	}
	return notifications, nil
}
