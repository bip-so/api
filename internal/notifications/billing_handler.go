package notifications

func (s notificationService) CanvasLimitExceedHandler(notification *PostNotification) {
	event := SystemNotificationEntity.Events[notification.Event]
	notification.Entity = SystemNotifications
	notification.Activity = event.Activity
	studio, _ := App.Repo.GetStudioByID(*notification.StudioID)
	notification.Text = "You have exceeded the 25 private canvas limit of the free plan. Newly published canvases will be made public in 24 hours. Upgrade now to keep them private!"
	notification.Priority = event.Priority
	userIDs, _ := App.Repo.GetStudioModeratorsUserIDs(*notification.StudioID)
	billingUserIDs, _ := App.Repo.GetStudioBillingMemberUserIDs(*notification.StudioID)
	userIDs = append(userIDs, billingUserIDs...)
	notification.NotifierIDs = s.GetUniqueIDs(userIDs)
	notification.ExtraData.AppUrl = s.GenerateStudioBillingUrl(studio.Handle)
	notification.ExtraData.DiscordComponents = []interface{}{
		ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				MessageBtnComponent{
					Type:  2,
					Label: "Upgrade",
					Style: 5,
					Url:   notification.ExtraData.AppUrl,
				},
			},
		},
	}
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "Upgrade", notification.ExtraData.AppUrl)
	notification.ExtraData.EmailSubject = "[Limit Exceeded] Upgrade now to keep your canvases private!"
	s.handleNotificationCreation(notification)
}
