package notifications

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

func (s notificationService) JoinedStudioHandler(notification *PostNotification) {
	event := FollowedMyStudioEntity.Events[notification.Event]
	notification.Entity = SystemNotifications
	notification.Activity = event.Activity
	user, _ := App.Repo.GetUser(notification.CreatedByID)
	studio, _ := App.Repo.GetStudioByID(*notification.StudioID)
	notification.Text = fmt.Sprintf(event.Text, user.Username, studio.DisplayName)
	notification.Priority = event.Priority
	userIDs, _ := App.Repo.GetStudioModeratorsUserIDs(*notification.StudioID)
	notification.NotifierIDs = s.GetUniqueIDs(userIDs)
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
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "💬 View", notification.ExtraData.AppUrl)
	s.handleNotificationCreation(notification)
}

func (s notificationService) GetStudioIntegrationWithStatus(studioId uint64) (*models.StudioIntegration, bool) {
	studio, _ := App.Repo.GetStudioByID(studioId)
	studioIntegration, err := s.GetStudioIntegration(studioId, models.SLACK_INTEGRATION_TYPE)
	if err != nil {
		if err.Error() == "record not found" {
			fmt.Println("studio is not integrated")
			return nil, true
		} else {
			return nil, false
		}
	} else {
		if studioIntegration != nil {
			fmt.Println("Studio has integration", studioIntegration)
			if !studio.SlackNotificationsEnabled {
				fmt.Println("Studio notification not allowed in product, returning..")
				return studioIntegration, false
			}
		}
	}
	return studioIntegration, true
}
