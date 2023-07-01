package notifications

import (
	"encoding/json"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

func (s notificationService) NotionImportHandler(notification *PostNotification) {
	event := SystemNotificationEntity.Events[notification.Event]
	notification.Entity = SystemNotifications
	notification.Activity = event.Activity
	notification.Text = fmt.Sprintf(event.Text)
	notification.Priority = event.Priority
	//userIDs, _ := App.Repo.GetStudioModeratorsUserIDs(*notification.StudioID)
	//notification.NotifierIDs = s.GetUniqueIDs(userIDs)
	canvasRepo, _ := App.Repo.GetCanvasRepoByID(notification.ExtraData.CanvasRepoID)
	notification.ExtraData.AppUrl = s.GenerateCanvasBranchUrl(canvasRepo.Key, canvasRepo.Name, *notification.StudioID, *canvasRepo.DefaultBranchID)
	notification.ExtraData.DiscordComponents = []interface{}{
		ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				MessageBtnComponent{
					Type:  2,
					Label: "🚀 View",
					Style: 5,
					Url:   notification.ExtraData.AppUrl,
				},
			},
		},
	}
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "🚀 View", notification.ExtraData.AppUrl)
	s.handleNotificationCreation(notification)
}

func (s notificationService) FileImportHandler(notification *PostNotification) {
	event := SystemNotificationEntity.Events[notification.Event]
	notification.Entity = SystemNotifications
	notification.Activity = event.Activity
	//studio, _ := App.Repo.GetStudioByID(*notification.StudioID)
	notification.Text = fmt.Sprintf(event.Text)
	notification.Priority = event.Priority
	//userIDs, _ := App.Repo.GetStudioModeratorsUserIDs(*notification.StudioID)
	//notification.NotifierIDs = s.GetUniqueIDs(userIDs)
	//notification.ExtraData.AppUrl = s.GenerateStudioUrl(studio.Handle)
	canvasRepo, _ := App.Repo.GetCanvasRepoByID(notification.ExtraData.CanvasRepoID)
	notification.ExtraData.AppUrl = s.GenerateCanvasBranchUrl(canvasRepo.Key, canvasRepo.Name, *notification.StudioID, *canvasRepo.DefaultBranchID)
	notification.ExtraData.DiscordComponents = []interface{}{
		ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				MessageBtnComponent{
					Type:  2,
					Label: "🚀 View",
					Style: 5,
					Url:   notification.ExtraData.AppUrl,
				},
			},
		},
	}
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "🚀 View", notification.ExtraData.AppUrl)
	s.handleNotificationCreation(notification)
}

func (s notificationService) TranslateCanvasHandler(notification *PostNotification) {
	event := SystemNotificationEntity.Events[notification.Event]
	notification.Entity = SystemNotifications
	notification.Activity = event.Activity
	studio, _ := App.Repo.GetStudioByID(*notification.StudioID)
	notification.Text = fmt.Sprintf(event.Text)
	notification.Priority = event.Priority
	//userIDs, _ := App.Repo.GetStudioModeratorsUserIDs(*notification.StudioID)
	//notification.NotifierIDs = s.GetUniqueIDs(userIDs)
	notification.ExtraData.AppUrl = s.GenerateStudioUrl(studio.Handle)
	notification.ExtraData.DiscordComponents = []interface{}{
		ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				MessageBtnComponent{
					Type:  2,
					Label: "🚀 View",
					Style: 5,
					Url:   notification.ExtraData.AppUrl,
				},
			},
		},
	}
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "🚀 View", notification.ExtraData.AppUrl)
	s.handleNotificationCreation(notification)
}

func (s notificationService) DiscordIntegrationTaskHandler(notification *PostNotification) {
	event := SystemNotificationEntity.Events[notification.Event]
	notification.Entity = SystemNotifications
	notification.Activity = event.Activity
	studio, _ := App.Repo.GetStudioByID(*notification.StudioID)
	notification.Text = fmt.Sprintf(event.Text)
	notification.Priority = event.Priority
	userIDs, _ := App.Repo.GetStudioModeratorsUserIDs(*notification.StudioID)
	notification.NotifierIDs = s.GetUniqueIDs(userIDs)
	notification.ExtraData.AppUrl = s.GenerateStudioIntegrationSettingsUrl(studio.Handle)
	notification.ExtraData.DiscordComponents = []interface{}{
		ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				MessageBtnComponent{
					Type:  2,
					Label: "🚀 View",
					Style: 5,
					Url:   notification.ExtraData.AppUrl,
				},
			},
		},
	}
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "🚀 View", notification.ExtraData.AppUrl)
	s.handleNotificationCreation(notification)
}

func (s notificationService) DiscordIntegrationTaskFailedHandler(notification *PostNotification) {
	event := SystemNotificationEntity.Events[notification.Event]
	notification.Entity = SystemNotifications
	notification.Activity = event.Activity
	studio, _ := App.Repo.GetStudioByID(*notification.StudioID)
	notification.Text = fmt.Sprintf(event.Text)
	notification.Priority = event.Priority
	userIDs, _ := App.Repo.GetStudioModeratorsUserIDs(*notification.StudioID)
	notification.NotifierIDs = s.GetUniqueIDs(userIDs)
	notification.ExtraData.AppUrl = s.GenerateStudioIntegrationSettingsUrl(studio.Handle)
	notification.ExtraData.DiscordComponents = []interface{}{
		ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				MessageBtnComponent{
					Type:  2,
					Label: "🚀 View",
					Style: 5,
					Url:   notification.ExtraData.AppUrl,
				},
			},
		},
	}
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "🚀 View", notification.ExtraData.AppUrl)
	s.handleNotificationCreation(notification)
}

func (s notificationService) CreateRequestToJoinStudioHandler(notification *PostNotification) {
	event := SystemNotificationEntity.Events[notification.Event]
	notification.Entity = SystemNotifications
	notification.Activity = event.Activity
	requestedUser, _ := App.Repo.GetUser(notification.CreatedByID)
	studio, _ := App.Repo.GetStudioByID(*notification.StudioID)
	notification.Text = fmt.Sprintf(event.Text, requestedUser.FullName, studio.DisplayName)
	notification.Priority = event.Priority
	userIDs, _ := App.Repo.GetStudioModeratorsUserIDs(*notification.StudioID)
	notification.NotifierIDs = s.GetUniqueIDs(userIDs)
	notification.ExtraData.AppUrl = s.GenerateStudioPendingRequestsUrl(studio.Handle)
	notification.ExtraData.DiscordComponents = []interface{}{
		ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				MessageBtnComponent{
					Type:  2,
					Label: "🚀 View",
					Style: 5,
					Url:   notification.ExtraData.AppUrl,
				},
			},
		},
	}
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "🚀 View", notification.ExtraData.AppUrl)
	s.handleNotificationCreation(notification)
}

func (s notificationService) RejectRequestToJoinStudioHandler(notification *PostNotification) {
	event := SystemNotificationEntity.Events[notification.Event]
	notification.Entity = SystemNotifications
	notification.Activity = event.Activity
	studio, _ := App.Repo.GetStudioByID(*notification.StudioID)
	rejectedUser, _ := App.Repo.GetUser(notification.CreatedByID)
	notification.Text = fmt.Sprintf(event.Text, rejectedUser.FullName)
	notification.Priority = event.Priority
	notification.ExtraData.AppUrl = s.GenerateStudioUrl(studio.Handle)
	notification.ExtraData.DiscordComponents = []interface{}{
		ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				MessageBtnComponent{
					Type:  2,
					Label: "🚀 View",
					Style: 5,
					Url:   notification.ExtraData.AppUrl,
				},
			},
		},
	}
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "🚀 View", notification.ExtraData.AppUrl)

	notification.ReactorID = &notification.CreatedByID
	requesterID := notification.NotifierIDs[0]
	query := map[string]interface{}{"created_by_id": requesterID, "event": CreateRequestToJoinStudio, "studio_id": notification.StudioID}
	notifications, err := App.Repo.Get(query)
	if err != nil {
		fmt.Println(err)
	}
	for _, notificationInstance := range *notifications {
		var extraData map[string]interface{}
		var discordMessage []string
		json.Unmarshal(notificationInstance.ExtraData, &extraData)
		actor, _ := App.Repo.GetUser(notificationInstance.CreatedByID)
		discordMessage = append(discordMessage, fmt.Sprintf("%s requesting to join `%s` studio", actor.Username, studio.DisplayName))
		discordMessage = append(discordMessage, fmt.Sprintf("**%s** by %s", "Rejected", rejectedUser.Username))
		extraData["discordMessage"] = discordMessage
		extraData["actionStatus"] = notification.ExtraData.Status
		extraData["slackComponents"] = s.SlackNotificationBlockBuilder(discordMessage, "📄 View", notification.ExtraData.AppUrl)
		extraData["discordComponents"] = []interface{}{
			ActionRowsComponent{
				Type: 1,
				Components: []interface{}{
					MessageBtnComponent{
						Type:  2,
						Label: "🚀 View",
						Style: 5,
						Url:   notification.ExtraData.AppUrl,
					},
				},
			},
		}
		notificationInstance.ExtraData, _ = json.Marshal(extraData)
		notificationInstance.ReactorID = notification.ReactorID
		s.RemoveActivity(&notificationInstance)
		//s.UpdateActivity(&notificationInstance)
		if notificationInstance.DiscordDmID != nil {
			//s.updateDiscordNotification(notificationInstance, discordMessage)
			s.deleteDiscordNotification(notificationInstance)
		}
		if notificationInstance.SlackDmID != nil && notificationInstance.SlackChannelID != nil {
			s.deleteSlackNotification(notificationInstance)
		}
		App.Repo.Manger.HardDeleteByID(models.NOTIFICATIONS, notificationInstance.ID)
		s.CreateNewNotification(&notificationInstance)
	}

	s.handleNotificationCreation(notification)
}

func (s notificationService) AcceptRequestToJoinStudioHandler(notification *PostNotification) {
	event := SystemNotificationEntity.Events[notification.Event]
	notification.Entity = SystemNotifications
	notification.Activity = event.Activity
	studio, _ := App.Repo.GetStudioByID(*notification.StudioID)
	acceptedUser, _ := App.Repo.GetUser(notification.CreatedByID)
	notification.Text = fmt.Sprintf(event.Text, acceptedUser.FullName)
	notification.Priority = event.Priority
	notification.ExtraData.AppUrl = s.GenerateStudioUrl(studio.Handle)
	notification.ExtraData.DiscordComponents = []interface{}{
		ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				MessageBtnComponent{
					Type:  2,
					Label: "🚀 View",
					Style: 5,
					Url:   notification.ExtraData.AppUrl,
				},
			},
		},
	}
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "🚀 View", notification.ExtraData.AppUrl)

	notification.ReactorID = &notification.CreatedByID
	requesterID := notification.NotifierIDs[0]
	query := map[string]interface{}{"created_by_id": requesterID, "event": CreateRequestToJoinStudio, "studio_id": notification.StudioID}
	notifications, err := App.Repo.Get(query)
	if err != nil {
		fmt.Println(err)
	}
	for _, notificationInstance := range *notifications {
		var extraData map[string]interface{}
		var discordMessage []string
		json.Unmarshal(notificationInstance.ExtraData, &extraData)
		actor, _ := App.Repo.GetUser(notificationInstance.CreatedByID)
		discordMessage = append(discordMessage, fmt.Sprintf("%s requesting to join `%s` studio", actor.Username, studio.DisplayName))
		discordMessage = append(discordMessage, fmt.Sprintf("**%s** by %s", "Accepted", acceptedUser.Username))
		extraData["discordMessage"] = discordMessage
		extraData["actionStatus"] = notification.ExtraData.Status
		extraData["slackComponents"] = s.SlackNotificationBlockBuilder(discordMessage, "📄 View", notification.ExtraData.AppUrl)
		extraData["discordComponents"] = []interface{}{
			ActionRowsComponent{
				Type: 1,
				Components: []interface{}{
					MessageBtnComponent{
						Type:  2,
						Label: "🚀 View",
						Style: 5,
						Url:   notification.ExtraData.AppUrl,
					},
				},
			},
		}
		notificationInstance.ExtraData, _ = json.Marshal(extraData)
		notificationInstance.ReactorID = notification.ReactorID
		s.RemoveActivity(&notificationInstance)
		//s.UpdateActivity(&notificationInstance)
		if notificationInstance.DiscordDmID != nil {
			//s.updateDiscordNotification(notificationInstance, discordMessage)
			s.deleteDiscordNotification(notificationInstance)
		}
		if notificationInstance.SlackDmID != nil && notificationInstance.SlackChannelID != nil {
			s.deleteSlackNotification(notificationInstance)
		}
		App.Repo.Manger.HardDeleteByID(models.NOTIFICATIONS, notificationInstance.ID)
		s.CreateNewNotification(&notificationInstance)
	}

	s.handleNotificationCreation(notification)
}
