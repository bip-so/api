package notifications

import (
	"encoding/json"
	"fmt"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

func (s notificationService) CanvasMergedHandler(notification *PostNotification) {
	event := SystemNotificationEntity.Events[notification.Event]
	notification.Entity = SystemNotifications
	notification.Activity = event.Activity

	user, _ := App.Repo.GetUser(notification.CreatedByID)
	canvasRepo, _ := App.Repo.GetCanvasRepoByID(notification.ExtraData.CanvasRepoID)
	notification.Text = fmt.Sprintf(event.Text, user.Username, canvasRepo.Name)

	notification.Priority = event.Priority
	notification.IsPersonal = FollowedMyStudioEntity.IsPersonal

	mergeRequest := models.MergeRequest{}
	err := json.Unmarshal([]byte(notification.ExtraData.Data), &mergeRequest)
	if err != nil {
		fmt.Println("Error on unmarshal merge request", err)
		return
	}
	// Get all the moderators of destination branch except requester
	destinationBranchModUserIDs, _ := App.Repo.GetCanvasBranchModeratorsUserIDs(mergeRequest.DestinationBranchID)
	// Get all the moderator & editors of the source branch except requester
	sourceBranchModAndEditorUserIDs, _ := App.Repo.GetCanvasBranchModeratorsAndEditorsUserIDs(mergeRequest.SourceBranchID)
	userIDs := append(destinationBranchModUserIDs, sourceBranchModAndEditorUserIDs...)
	notification.NotifierIDs = s.GetUniqueIDs(userIDs)
	notification.ExtraData.AppUrl = s.GenerateCanvasBranchUrl(canvasRepo.Key, canvasRepo.Name, canvasRepo.StudioID, *canvasRepo.DefaultBranchID)
	notification.ExtraData.DiscordComponents = []interface{}{
		ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				MessageBtnComponent{
					Type:  2,
					Label: "📄 View",
					Style: 5,
					Url:   notification.ExtraData.AppUrl,
				},
			},
		},
	}
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	notification.ExtraData.EmailSubject = fmt.Sprintf("Canvas Merged - %s", canvasRepo.Name)
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "📄 View", notification.ExtraData.AppUrl)

	s.handleNotificationCreation(notification)
}

func (s notificationService) PublishRequestedHandler(notification *PostNotification) {
	event := PublishAndMergeRequestsEntity.Events[notification.Event]
	notification.Entity = PublishAndMergeRequests
	notification.Activity = event.Activity

	user, _ := App.Repo.GetUser(notification.CreatedByID)
	canvasRepo, _ := App.Repo.GetCanvasRepoByID(notification.ExtraData.CanvasRepoID)
	notification.Text = fmt.Sprintf(event.Text, user.Username, canvasRepo.Name, canvasRepo.Studio.DisplayName)

	notification.Priority = event.Priority
	notification.IsPersonal = PublishAndMergeRequestsEntity.IsPersonal

	collectionModUserIDs, _ := App.Repo.GetCollectionModeratorUserIDs(notification.ExtraData.CollectionID)
	branchModUserIDs, _ := App.Repo.GetCanvasBranchModeratorsUserIDs(notification.ExtraData.CanvasBranchID) // objectID is branchID here
	userIDs := append(collectionModUserIDs, branchModUserIDs...)
	notification.NotifierIDs = s.GetUniqueIDs(userIDs)
	notification.ExtraData.AppUrl = s.GenerateCanvasBranchUrl(canvasRepo.Key, canvasRepo.Name, canvasRepo.StudioID, *canvasRepo.DefaultBranchID)
	notification.ExtraData.DiscordComponents = []interface{}{
		ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				MessageBtnComponent{
					Type:  2,
					Label: "🆕 View",
					Style: 5,
					Url:   notification.ExtraData.AppUrl,
				},
			},
		},
	}
	notification.ExtraData.DiscordMessage = []string{
		notification.Text,
	}
	if notification.ExtraData.Message != "" {
		notification.ExtraData.DiscordMessage = append(notification.ExtraData.DiscordMessage,
			fmt.Sprintf("```%s```", notification.ExtraData.Message))
	}
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "📄 View", notification.ExtraData.AppUrl)

	s.handleNotificationCreation(notification)
}

func (s notificationService) MergeRequestedHandler(notification *PostNotification) {
	event := PublishAndMergeRequestsEntity.Events[notification.Event]
	notification.Entity = PublishAndMergeRequests
	notification.Activity = event.Activity

	user, _ := App.Repo.GetUser(notification.CreatedByID)
	canvasRepo, _ := App.Repo.GetCanvasRepoByID(notification.ExtraData.CanvasRepoID)
	notification.Text = fmt.Sprintf(event.Text, user.Username, canvasRepo.Name)

	notification.Priority = event.Priority
	notification.IsPersonal = PublishAndMergeRequestsEntity.IsPersonal

	mergeRequest, _ := App.Repo.GetMergeRequestByID(*notification.ObjectID)
	branchModUserIDs, _ := App.Repo.GetCanvasBranchModeratorsUserIDs(mergeRequest.DestinationBranchID) // objectID is branch id here
	notification.NotifierIDs = s.GetUniqueIDs(branchModUserIDs)
	notification.ExtraData.AppUrl = s.GenerateMergeRequestUrl(canvasRepo.Key, canvasRepo.Name, canvasRepo.StudioID, *canvasRepo.DefaultBranchID, mergeRequest.ID)
	notification.ExtraData.DiscordComponents = []interface{}{
		ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				MessageBtnComponent{
					Type:  2,
					Label: "🔀 View",
					Style: 5,
					Url:   notification.ExtraData.AppUrl,
				},
			},
		},
	}
	notification.ExtraData.DiscordMessage = []string{
		notification.Text,
	}
	if notification.ExtraData.Message != "" {
		notification.ExtraData.DiscordMessage = append(notification.ExtraData.DiscordMessage,
			fmt.Sprintf("```%s```", notification.ExtraData.Message))
	}
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "📄 View", notification.ExtraData.AppUrl)

	s.handleNotificationCreation(notification)
}

func (s notificationService) PublishRequestedUpdateHandler(notification *PostNotification) {
	// Question should we have to update the activity and event names after the notification is updated with reactor??
	event := ResponseToMyRequestsEntity.Events[notification.Event]
	notification.Entity = ResponseToMyRequests
	notification.Activity = event.Activity

	user, _ := App.Repo.GetUser(notification.CreatedByID)
	canvasRepo, _ := App.Repo.GetCanvasRepoByID(notification.ExtraData.CanvasRepoID)
	notification.Text = fmt.Sprintf(event.Text, user.Username, notification.ExtraData.Status, canvasRepo.Name)

	notification.Priority = event.Priority
	notification.IsPersonal = ResponseToMyRequestsEntity.IsPersonal
	// Get the notification extra data
	// GEt the notifications based on objectid, contentobject, event.
	// update the notifications with new extra data.
	notification.ReactorID = &notification.CreatedByID
	query := map[string]interface{}{"object_id": notification.ObjectID, "content_object": notification.ContentObject, "event": PublishRequested}
	notifications, err := App.Repo.Get(query)
	if err != nil {
		fmt.Println(err)
	}
	var requesterID uint64
	for _, notificationInstance := range *notifications {
		var extraData map[string]interface{}
		var discordMessage []string
		json.Unmarshal(notificationInstance.ExtraData, &extraData)
		message, _ := extraData["message"].(string)
		actor, _ := App.Repo.GetUser(notificationInstance.CreatedByID)
		discordMessage = append(discordMessage, fmt.Sprintf("@%s is requesting to publish `📄 %s` in `🎨 %s`", actor.Username, canvasRepo.Name, canvasRepo.Studio.DisplayName))
		if message != "" {
			discordMessage = append(discordMessage, fmt.Sprintf("```%s```", message))
		}
		discordMessage = append(discordMessage, fmt.Sprintf("**%s** by @%s", notification.ExtraData.Status, user.Username))
		extraData["discordMessage"] = discordMessage
		extraData["actionStatus"] = notification.ExtraData.Status
		extraData["slackComponents"] = s.SlackNotificationBlockBuilder(discordMessage, "📄 View", extraData["appUrl"].(string))
		notificationInstance.ExtraData, _ = json.Marshal(extraData)
		requesterID = notificationInstance.CreatedByID
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
	//updates := map[string]interface{}{
	//	"extra_data": extraData,
	//	"reactor_id": notification.ReactorID,
	//}
	//err = App.Repo.Delete(query)
	//if err != nil {
	//	fmt.Println(err)
	//}
	// Send the notification to the requester
	notification.NotifierIDs = []uint64{requesterID}
	notification.ExtraData.AppUrl = s.GenerateCanvasBranchUrl(canvasRepo.Key, canvasRepo.Name, canvasRepo.StudioID, *canvasRepo.DefaultBranchID)
	notification.ExtraData.DiscordComponents = []interface{}{
		ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				MessageBtnComponent{
					Type:  2,
					Label: "🆕 View",
					Style: 5,
					Url:   notification.ExtraData.AppUrl,
				},
			},
		},
	}
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "📄 View", notification.ExtraData.AppUrl)

	s.handleNotificationCreation(notification)
}

func (s notificationService) MergeRequestedUpdateHandler(notification *PostNotification) {
	// Question should we have to update the activity and event names after the notification is updated with reactor??
	event := ResponseToMyRequestsEntity.Events[notification.Event]
	notification.Entity = ResponseToMyRequests
	notification.Activity = event.Activity

	user, _ := App.Repo.GetUser(notification.CreatedByID)
	canvasRepo, _ := App.Repo.GetCanvasRepoByID(notification.ExtraData.CanvasRepoID)
	notification.Text = fmt.Sprintf(event.Text, user.Username, notification.ExtraData.Status, canvasRepo.Name)

	notification.Priority = event.Priority
	notification.IsPersonal = ResponseToMyRequestsEntity.IsPersonal
	// Get the notification extra data
	// GEt the notifications based on objectid, contentobject, event.
	// update the notifications with new extra data.
	notification.ReactorID = &notification.CreatedByID
	query := map[string]interface{}{"object_id": notification.ObjectID, "content_object": notification.ContentObject, "event": MergeRequested}
	notifications, err := App.Repo.Get(query)
	if err != nil {
		fmt.Println(err)
	}
	var requesterID uint64
	fmt.Println("length of notifications", len(*notifications))
	for _, notificationInstance := range *notifications {
		var discordMessage []string
		var extraData map[string]interface{}
		fmt.Println("notifier id ", notificationInstance.NotifierID, notificationInstance.ContentObject)
		json.Unmarshal(notificationInstance.ExtraData, &extraData)
		actor, _ := App.Repo.GetUser(notificationInstance.CreatedByID)
		discordMessage = append(discordMessage, fmt.Sprintf("@%s is requesting to merge changes to `📄 %s`", actor.Username, canvasRepo.Name))
		if notification.ExtraData.Message != "" {
			discordMessage = append(discordMessage, fmt.Sprintf("```%s```", notification.ExtraData.Message))
		}
		discordMessage = append(discordMessage, fmt.Sprintf("**%s** by @%s", notification.ExtraData.Status, user.Username))
		extraData["discordMessage"] = discordMessage
		extraData["actionStatus"] = notification.ExtraData.Status
		extraData["slackComponents"] = s.SlackNotificationBlockBuilder(discordMessage, "📄 View", extraData["appUrl"].(string))
		notificationInstance.ExtraData, _ = json.Marshal(extraData)
		requesterID = notificationInstance.CreatedByID
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
	notification.ExtraData.AppUrl = s.GenerateCanvasBranchUrl(canvasRepo.Key, canvasRepo.Name, canvasRepo.StudioID, *canvasRepo.DefaultBranchID)
	notification.ExtraData.DiscordComponents = []interface{}{
		ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				MessageBtnComponent{
					Type:  2,
					Label: "📄 View",
					Style: 5,
					Url:   notification.ExtraData.AppUrl,
				},
			},
		},
	}
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "📄 View", notification.ExtraData.AppUrl)

	s.handleNotificationCreation(notification)
}

func (s notificationService) BipMarkMessageAddedHandler(notification *PostNotification) {
	event := SystemNotificationEntity.Events[notification.Event]
	notification.Entity = SystemNotifications
	notification.Activity = event.Activity
	canvasRepo, _ := App.Repo.GetCanvasRepoByID(notification.ExtraData.CanvasRepoID)
	notification.Priority = event.Priority
	notification.IsPersonal = FollowedMyStudioEntity.IsPersonal
	blocks, _ := App.Repo.GetBlocks(map[string]interface{}{"canvas_branch_id": notification.ExtraData.CanvasBranchID})
	for _, block := range blocks {
		notification.ExtraData.AppUrl = s.GenerateCanvasBranchBlockUrl(canvasRepo.Key, canvasRepo.Name, canvasRepo.StudioID, *canvasRepo.DefaultBranchID, block.UUID.String())
		var attributes BlockAttributes
		var user models.User
		json.Unmarshal(block.Attributes, &attributes)
		// if messageId is not present it will not be a bip mark block. so ignoring other blocks.
		if attributes.MessageID == "" {
			continue
		}
		contentObject := models.BLOCK
		notifications, _ := App.Repo.Get(map[string]interface{}{"object_id": block.ID, "content_object": contentObject, "event": notification.Event})
		// if notification is already sent for this block id then we don't send again.
		if len(*notifications) > 0 {
			continue
		}
		if block.UpdatedByID != 0 {
			user, _ = App.Repo.GetUser(block.UpdatedByID)
		} else {
			user, _ = App.Repo.GetUser(block.CreatedByID)
		}
		notification.Text = fmt.Sprintf(event.Text, canvasRepo.Name, user.Username)
		notification.ExtraData.DiscordMessage = []string{notification.Text}
		message, _ := App.Repo.GetMessage(map[string]interface{}{"uuid": attributes.MessageID})
		notification.ExtraData.DiscordMessage = append(notification.ExtraData.DiscordMessage, fmt.Sprintf("```%s```", message.Text))
		notification.NotifierIDs = []uint64{message.AuthorID}
		notification.ObjectID = &block.ID
		notification.ContentObject = &contentObject
		notification.ExtraData.DiscordComponents = []interface{}{
			ActionRowsComponent{
				Type: 1,
				Components: []interface{}{
					MessageBtnComponent{
						Type:  2,
						Label: "📄 View",
						Style: 5,
						Url:   notification.ExtraData.AppUrl,
					},
				},
			},
		}
		slackMessage := fmt.Sprintf("Your slack message was added to `📄 %s` by `@%s`", canvasRepo.Name, user.Username) + "/n" + fmt.Sprintf("```%s```", message.Text)
		notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder([]string{slackMessage}, "📄 View", notification.ExtraData.AppUrl)
		s.handleNotificationCreation(notification)
	}
}

func (s notificationService) MergeRequestDeleteHandler(mrID uint64) {
	notifications, err := App.Repo.Get(map[string]interface{}{"object_id": mrID, "content_object": models.MERGEREQUEST})
	if err != nil {
		s.logg.Error("[MergeRequestDeleteHandler] Error on getting notifications with merge requestID")
		s.logg.Error(err.Error())
	}
	for _, notification := range *notifications {
		App.Service.RemoveActivity(&notification)
		if notification.DiscordDmID != nil {
			App.Service.deleteDiscordNotification(notification)
		} else if notification.SlackDmID != nil && notification.SlackChannelID != nil {
			App.Service.deleteSlackNotification(notification)
		}
	}
	App.Repo.Delete(map[string]interface{}{"object_id": mrID, "content_object": models.MERGEREQUEST})
}

func (s notificationService) PublishRequestDeleteHandler(prInstance *models.PublishRequest) {
	notifications, err := App.Repo.Get(map[string]interface{}{"object_id": prInstance.ID, "content_object": models.PUBLISHREQUEST})
	if err != nil {
		s.logg.Error("[MergeRequestDeleteHandler] Error on getting notifications with merge requestID")
		s.logg.Error(err.Error())
	}
	for _, notification := range *notifications {
		App.Service.RemoveActivity(&notification)
		if notification.DiscordDmID != nil {
			App.Service.deleteDiscordNotification(notification)
		} else if notification.SlackDmID != nil && notification.SlackChannelID != nil {
			App.Service.deleteSlackNotification(notification)
		}
	}
	App.Repo.Delete(map[string]interface{}{"object_id": prInstance.ID, "content_object": models.PUBLISHREQUEST})
}

func (s notificationService) DeleteModUsersOnCanvas(prInstance *models.PublishRequest) ([]uint64, []uint64) {
	canvasBranchPerms, _ := App.Repo.GetCanvasBranchPerms(map[string]interface{}{"canvas_branch_id": prInstance.CanvasBranchID, "permission_group": models.PGCanvasModerateSysName})
	canvasRepo, _ := App.Repo.GetCanvasRepoByID(prInstance.CanvasRepositoryID)
	repoCreatorMember, _ := App.Repo.GetMemberByUserAndStudioID(canvasRepo.CreatedByID, canvasRepo.StudioID)
	deleteCanvasPerms := []models.CanvasBranchPermission{}
	userIDs := []uint64{}
	roleIDs := []uint64{}
	for _, perm := range canvasBranchPerms {
		if perm.MemberId != nil && repoCreatorMember.ID != *perm.MemberId {
			userIDs = append(userIDs, perm.Member.UserID)
			deleteCanvasPerms = append(deleteCanvasPerms, perm)
		} else if perm.RoleId != nil {
			roleIDs = append(roleIDs, *perm.RoleId)
			deleteCanvasPerms = append(deleteCanvasPerms, perm)
		}
	}
	App.Repo.db.Delete(deleteCanvasPerms)
	return userIDs, roleIDs
}
