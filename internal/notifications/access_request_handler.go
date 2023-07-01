package notifications

import (
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/integrations"
)

func (s notificationService) AccessRequestedHandler(notification *PostNotification) {
	event := PublishAndMergeRequestsEntity.Events[notification.Event]
	notification.Entity = PublishAndMergeRequests
	notification.Activity = event.Activity

	user, _ := App.Repo.GetUser(notification.CreatedByID)
	canvasRepo, _ := App.Repo.GetCanvasRepoByID(notification.ExtraData.CanvasRepoID)
	notification.Text = fmt.Sprintf(event.Text, user.Username, canvasRepo.Name)

	notification.Priority = event.Priority
	notification.IsPersonal = InvitesEntity.IsPersonal

	userIDs, _ := App.Repo.GetStudioModeratorsUserIDs(*notification.StudioID)
	notification.NotifierIDs = s.GetUniqueIDs(userIDs)

	notification.ExtraData.AppUrl = s.GenerateCanvasBranchUrl(canvasRepo.Key, canvasRepo.Name, *notification.StudioID, *canvasRepo.DefaultBranchID)
	notification.ExtraData.DiscordComponents = []interface{}{
		ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				discordgo.SelectMenu{
					CustomID:    "accessrequest",
					Placeholder: "Grant Permission",
					MaxValues:   1,
					Options: []discordgo.SelectMenuOption{
						{
							Label:       "Moderate",
							Value:       models.PGCanvasModerateSysName,
							Description: "",
						},
						{
							Label:       "Edit",
							Value:       models.PGCanvasEditSysName,
							Description: "",
						},
						{
							Label:       "Reply",
							Value:       models.PGCanvasCommentSysName,
							Description: "",
						},
						{
							Label:       "View",
							Value:       models.PGCanvasViewSysName,
							Description: "",
						},
					},
				},
			},
		},
		integrations.ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				integrations.MessageBtnComponent{
					Type:  2,
					Label: "📄 View Canvas",
					Style: 5,
					Url:   notification.ExtraData.AppUrl,
				},
			},
		},
	}
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	notification.ExtraData.SlackComponents = s.SlackAccessRequestMessageBuilder(notification.ExtraData.DiscordMessage, "📄 View Canvas", notification.ExtraData.AppUrl)

	s.handleNotificationCreation(notification)
}

func (s notificationService) AccessRequestedUpdateHandler(notification *PostNotification) {
	event := PublishAndMergeRequestsEntity.Events[notification.Event]
	notification.Entity = PublishAndMergeRequests
	notification.Activity = event.Activity

	user, _ := App.Repo.GetUser(notification.CreatedByID)
	canvasRepo, _ := App.Repo.GetCanvasRepoByID(notification.ExtraData.CanvasRepoID)
	if notification.ExtraData.Status == models.ACCESS_REQUEST_ACCEPTED {
		notification.Text = fmt.Sprintf(event.Text, user.Username, PermissionsTextMap[notification.ExtraData.PermissionGroup], canvasRepo.Name)
	} else {
		notification.Text = fmt.Sprintf("%s has rejected your access in %s", user.Username, canvasRepo.Name)
	}
	notification.ExtraData.AppUrl = s.GenerateCanvasBranchUrl(canvasRepo.Key, canvasRepo.Name, *notification.StudioID, *canvasRepo.DefaultBranchID)

	notification.Priority = event.Priority
	notification.IsPersonal = InvitesEntity.IsPersonal
	notification.ReactorID = &notification.CreatedByID

	query := map[string]interface{}{"object_id": notification.ObjectID, "content_object": notification.ContentObject, "event": AccessRequested}
	notifications, err := App.Repo.Get(query)
	if err != nil {
		fmt.Println(err)
	}
	// update the query to get all the notifications to update
	var requesterID uint64
	for _, notificationInstance := range *notifications {
		var extraData map[string]interface{}
		var discordMessage []string
		json.Unmarshal(notificationInstance.ExtraData, &extraData)
		actor, _ := App.Repo.GetUser(notificationInstance.CreatedByID)
		discordMessage = append(discordMessage, fmt.Sprintf("@%s is requesting access to `📄 %s`", actor.Username, canvasRepo.Name))
		extraData["actionStatus"] = notification.ExtraData.Status
		extraData["permissionGroup"] = notification.ExtraData.PermissionGroup
		if notification.ExtraData.Status == models.ACCESS_REQUEST_REJECTED {
			discordMessage = append(discordMessage, "")
			discordMessage = append(discordMessage, fmt.Sprintf("**%s** by @%s", models.ACCESS_REQUEST_REJECTED, user.Username))
		} else {
			discordMessage = append(discordMessage, "")
			discordMessage = append(discordMessage, fmt.Sprintf("✅ **Granted** %s access by @%s", PermissionsTextMap[notification.ExtraData.PermissionGroup], user.Username))
		}
		extraData["discordMessage"] = discordMessage
		extraData["slackComponents"] = s.SlackNotificationBlockBuilder(discordMessage, "📄 View", extraData["appUrl"].(string))
		extraData["discordComponents"] = []interface{}{
			ActionRowsComponent{
				Type: 1,
				Components: []interface{}{
					integrations.MessageBtnComponent{
						Type:  2,
						Label: "📄 View Canvas",
						Style: 5,
						Url:   notification.ExtraData.AppUrl,
					},
				},
			},
		}
		notificationInstance.ExtraData, _ = json.Marshal(extraData)
		requesterID = notificationInstance.CreatedByID
		notificationInstance.ReactorID = notification.ReactorID

		s.RemoveActivity(&notificationInstance)
		//s.UpdateActivity(&notificationInstance)
		if notificationInstance.DiscordDmID != nil {
			s.deleteDiscordNotification(notificationInstance)
			//s.sendDiscordDmNotification(&notificationInstance, discordMessage)
		}
		if notificationInstance.SlackDmID != nil && notificationInstance.SlackChannelID != nil {
			s.deleteSlackNotification(notificationInstance)
		}
		App.Repo.Manger.HardDeleteByID(models.NOTIFICATIONS, notificationInstance.ID)
		s.CreateNewNotification(&notificationInstance)
		requesterID = notificationInstance.CreatedByID
	}
	//updates := map[string]interface{}{
	//	"extra_data": extraData,
	//	"reactor_id": notification.ReactorID,
	//}
	//notifications, err = App.Repo.Update(query, updates)
	//if err != nil {
	//	fmt.Println(err)
	//}
	// Send the notification to the requester
	notification.NotifierIDs = []uint64{requesterID}
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "📄 View Canvas", notification.ExtraData.AppUrl)

	s.handleNotificationCreation(notification)
}
