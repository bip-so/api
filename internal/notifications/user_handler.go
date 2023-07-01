package notifications

import "fmt"

func (s notificationService) FollowUserHandler(notification *PostNotification) {
	// X person followed you
	// Send notification to the followee
	event := FollowedMeEntity.Events[notification.Event]
	notification.Entity = FollowedMe
	notification.Activity = event.Activity
	user, _ := App.Repo.GetUser(notification.CreatedByID)
	notification.Text = fmt.Sprintf(event.Text, user.Username)
	notification.Priority = event.Priority
	notification.IsPersonal = FollowedMeEntity.IsPersonal
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	notification.ExtraData.SlackMessage = notification.Text
	s.handleNotificationCreation(notification)
}
