package notifications

import (
	"encoding/json"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gorm.io/datatypes"
)

/*
	BlockMentionHandler

	If block is of rough branch ignore the notification.
	If block is of main branch
		- get block mentions
		- check if each mention
*/
func (s notificationService) BlockMentionHandler(notification *PostNotification) {
	canvasBranch, _ := App.Repo.GetCanvasBranchByID(notification.ExtraData.CanvasBranchID)
	block, _ := App.Repo.GetBlockByID(*notification.ObjectID)
	canvasRepo, _ := App.Repo.GetCanvasRepoByID(canvasBranch.CanvasRepositoryID)
	if canvasBranch.IsRoughBranch {
		roughBranchID := notification.ExtraData.CanvasBranchID
		notification.ExtraData.CanvasBranchID = *canvasRepo.DefaultBranchID
		s.AddNotificationToRedis(roughBranchID, notification)
		return
	}

	event := MentionsEntity.Events[notification.Event]
	notification.Entity = Mentions
	notification.Activity = event.Activity

	userInstance, _ := App.Repo.GetUser(notification.CreatedByID)
	notification.Text = fmt.Sprintf(event.Text, userInstance.Username, "block")

	notification.Priority = event.Priority
	notification.IsPersonal = MentionsEntity.IsPersonal

	userIDs, roleIDs, _ := s.GetUserIDsFromMentions(*block.Mentions)
	notification.NotifierIDs = userIDs
	notification.RoleIDs = &roleIDs

	notification.StudioID = &block.CanvasBranch.CanvasRepository.StudioID
	notification.ExtraData.CanvasRepoID = canvasRepo.ID
	notification.ExtraData.CollectionID = block.CanvasBranch.CanvasRepository.CollectionID

	notification.ExtraData.AppUrl = s.GenerateCanvasBranchUrl(canvasRepo.Key, canvasRepo.Name, canvasRepo.StudioID, *canvasRepo.DefaultBranchID)
	notification.ExtraData.DiscordComponents = []interface{}{
		ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				MessageBtnComponent{
					Type:  2,
					Label: "💬 View",
					Style: 5,
					Url:   notification.ExtraData.AppUrl,
				},
			},
		},
	}
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	var blockChildren []BlockChildren
	json.Unmarshal(block.Children, &blockChildren)
	for _, children := range blockChildren {
		if children.Type == models.BlockTypeText {
			notification.ExtraData.DiscordMessage = append(notification.ExtraData.DiscordMessage, fmt.Sprintf("```%s```", children.Text))
			break
		}
	}
	notification.ExtraData.BlockUUID = block.UUID.String()
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "💬 View", notification.ExtraData.AppUrl)
	s.handleNotificationCreation(notification)
}

func (s notificationService) BlockThreadMentionHandler(notification *PostNotification) {
	canvasBranch, _ := App.Repo.GetCanvasBranchByID(notification.ExtraData.CanvasBranchID)
	blockThread, _ := App.Repo.GetBlockThreadByID(*notification.ObjectID)
	canvasRepo, _ := App.Repo.GetCanvasRepoByID(blockThread.CanvasRepositoryID)
	if canvasBranch.IsRoughBranch {
		roughBranchID := notification.ExtraData.CanvasBranchID
		notification.ExtraData.CanvasBranchID = *canvasRepo.DefaultBranchID
		s.AddNotificationToRedis(roughBranchID, notification)
		return
	}
	event := MentionsEntity.Events[notification.Event]
	notification.Entity = Mentions
	notification.Activity = event.Activity

	userInstance, _ := App.Repo.GetUser(notification.CreatedByID)
	notification.Text = fmt.Sprintf(event.Text, userInstance.Username, "block thread")

	notification.Priority = event.Priority
	notification.IsPersonal = MentionsEntity.IsPersonal

	userIDs, roleIDs, _ := s.GetUserIDsFromMentions(*blockThread.Mentions)
	notification.NotifierIDs = userIDs
	notification.RoleIDs = &roleIDs

	notification.StudioID = &blockThread.CanvasRepository.StudioID
	notification.ExtraData.CanvasRepoID = blockThread.CanvasRepositoryID
	notification.ExtraData.CollectionID = blockThread.CanvasRepository.CollectionID

	notification.ExtraData.AppUrl = s.GenerateBlockCommentUrl(canvasRepo.Key, canvasRepo.Name, blockThread.UUID.String(), canvasRepo.StudioID, *canvasRepo.DefaultBranchID)
	notification.ExtraData.DiscordComponents = []interface{}{
		ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				MessageBtnComponent{
					Type:  2,
					Label: "💬 View",
					Style: 5,
					Url:   notification.ExtraData.AppUrl,
				},
			},
		},
	}
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	notification.ExtraData.ActionOnText = blockThread.Text
	notification.ExtraData.DiscordMessage = append(notification.ExtraData.DiscordMessage, fmt.Sprintf("```%s```", notification.ExtraData.ActionOnText))
	notification.ExtraData.BlockUUID = blockThread.Block.UUID.String()
	notification.ExtraData.BlockThreadUUID = blockThread.UUID.String()
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "💬 View", notification.ExtraData.AppUrl)
	s.handleNotificationCreation(notification)
}

func (s notificationService) BlockThreadCommentMentionHandler(notification *PostNotification) {
	event := MentionsEntity.Events[notification.Event]
	notification.Entity = Mentions
	notification.Activity = event.Activity

	userInstance, _ := App.Repo.GetUser(notification.CreatedByID)
	notification.Text = fmt.Sprintf(event.Text, userInstance.Username, "block thread comment")

	notification.Priority = event.Priority
	notification.IsPersonal = MentionsEntity.IsPersonal

	blockThreadComment, _ := App.Repo.GetBlockThreadCommentByID(*notification.ObjectID)
	userIDs, roleIDs, _ := s.GetUserIDsFromMentions(*blockThreadComment.Mentions)
	notification.NotifierIDs = userIDs
	notification.RoleIDs = &roleIDs

	notification.StudioID = &blockThreadComment.Thread.CanvasRepository.StudioID
	notification.ExtraData.CanvasRepoID = blockThreadComment.Thread.CanvasRepositoryID
	notification.ExtraData.CollectionID = blockThreadComment.Thread.CanvasRepository.CollectionID

	canvasRepo, _ := App.Repo.GetCanvasRepoByID(blockThreadComment.Thread.CanvasRepositoryID)
	notification.ExtraData.AppUrl = s.GenerateCanvasBranchUrl(canvasRepo.Key, canvasRepo.Name, canvasRepo.StudioID, *canvasRepo.DefaultBranchID)
	notification.ExtraData.DiscordComponents = []interface{}{
		ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				MessageBtnComponent{
					Type:  2,
					Label: "💬 View",
					Style: 5,
					Url:   notification.ExtraData.AppUrl,
				},
			},
		},
	}
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	var data CommentsData
	json.Unmarshal(blockThreadComment.Data, &data)
	notification.ExtraData.ActionOnText = data.Text
	notification.ExtraData.DiscordMessage = append(notification.ExtraData.DiscordMessage, fmt.Sprintf("```%s```", notification.ExtraData.ActionOnText))
	notification.ExtraData.BlockUUID = blockThreadComment.Thread.Block.UUID.String()
	notification.ExtraData.BlockThreadUUID = blockThreadComment.Thread.UUID.String()
	notification.ExtraData.BlockThreadCommentUUID = blockThreadComment.UUID.String()
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "💬 View", notification.ExtraData.AppUrl)
	canvasBranch, _ := App.Repo.GetCanvasBranchByID(notification.ExtraData.CanvasBranchID)
	if canvasBranch.IsRoughBranch {
		roughBranchID := notification.ExtraData.CanvasBranchID
		notification.ExtraData.CanvasBranchID = *canvasRepo.DefaultBranchID
		s.AddNotificationToRedis(roughBranchID, notification)
		return
	}
	s.handleNotificationCreation(notification)
}

func (s notificationService) ReelMentionHandler(notification *PostNotification) {
	event := MentionsEntity.Events[notification.Event]
	notification.Entity = Mentions
	notification.Activity = event.Activity

	userInstance, _ := App.Repo.GetUser(notification.CreatedByID)
	notification.Text = fmt.Sprintf(event.Text, userInstance.Username, "reel")

	notification.Priority = event.Priority
	notification.IsPersonal = MentionsEntity.IsPersonal

	reel, _ := App.Repo.GetReelByID(*notification.ObjectID)
	userIDs, roleIDs, _ := s.GetUserIDsFromMentions(*reel.Mentions)
	notification.NotifierIDs = userIDs
	notification.RoleIDs = &roleIDs

	notification.StudioID = &reel.CanvasRepository.StudioID
	notification.ExtraData.CanvasRepoID = reel.CanvasRepositoryID
	notification.ExtraData.CollectionID = reel.CanvasRepository.CollectionID

	canvasRepo, _ := App.Repo.GetCanvasRepoByID(reel.CanvasRepositoryID)
	notification.ExtraData.AppUrl = s.GenerateCanvasBranchUrl(canvasRepo.Key, canvasRepo.Name, canvasRepo.StudioID, *canvasRepo.DefaultBranchID)
	notification.ExtraData.DiscordComponents = []interface{}{
		ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				MessageBtnComponent{
					Type:  2,
					Label: "💬 View",
					Style: 5,
					Url:   notification.ExtraData.AppUrl,
				},
			},
		},
	}
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	var data CommentsData
	json.Unmarshal(reel.ContextData, &data)
	notification.ExtraData.ActionOnText = data.Text
	notification.ExtraData.DiscordMessage = append(notification.ExtraData.DiscordMessage, fmt.Sprintf("```%s```", notification.ExtraData.ActionOnText))
	notification.ExtraData.BlockUUID = reel.Block.UUID.String()
	notification.ExtraData.ReelUUID = reel.UUID.String()
	notification.ExtraData.ReelID = reel.ID
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "💬 View", notification.ExtraData.AppUrl)
	canvasBranch, _ := App.Repo.GetCanvasBranchByID(notification.ExtraData.CanvasBranchID)
	if canvasBranch.IsRoughBranch {
		roughBranchID := notification.ExtraData.CanvasBranchID
		notification.ExtraData.CanvasBranchID = *canvasRepo.DefaultBranchID
		s.AddNotificationToRedis(roughBranchID, notification)
		return
	}
	s.handleNotificationCreation(notification)
}

func (s notificationService) ReelCommentMentionHandler(notification *PostNotification) {
	canvasBranch, _ := App.Repo.GetCanvasBranchByID(notification.ExtraData.CanvasBranchID)
	reelComment, _ := App.Repo.GetReelCommentByID(*notification.ObjectID)
	canvasRepo, _ := App.Repo.GetCanvasRepoByID(reelComment.Reel.CanvasRepositoryID)
	if canvasBranch.IsRoughBranch {
		roughBranchID := notification.ExtraData.CanvasBranchID
		notification.ExtraData.CanvasBranchID = *canvasRepo.DefaultBranchID
		s.AddNotificationToRedis(roughBranchID, notification)
		return
	}
	event := MentionsEntity.Events[notification.Event]
	notification.Entity = Mentions
	notification.Activity = event.Activity

	userInstance, _ := App.Repo.GetUser(notification.CreatedByID)
	notification.Text = fmt.Sprintf(event.Text, userInstance.Username, "reel comment")

	notification.Priority = event.Priority
	notification.IsPersonal = MentionsEntity.IsPersonal

	userIDs, roleIDs, _ := s.GetUserIDsFromMentions(*reelComment.Mentions)
	notification.NotifierIDs = userIDs
	notification.RoleIDs = &roleIDs

	notification.StudioID = &reelComment.Reel.CanvasRepository.StudioID
	notification.ExtraData.CanvasRepoID = reelComment.Reel.CanvasRepositoryID
	notification.ExtraData.CollectionID = reelComment.Reel.CanvasRepository.CollectionID

	notification.ExtraData.AppUrl = s.GenerateCanvasBranchUrl(canvasRepo.Key, canvasRepo.Name, canvasRepo.StudioID, *canvasRepo.DefaultBranchID)
	notification.ExtraData.DiscordComponents = []interface{}{
		ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				MessageBtnComponent{
					Type:  2,
					Label: "💬 View",
					Style: 5,
					Url:   notification.ExtraData.AppUrl,
				},
			},
		},
	}
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	var data CommentsData
	json.Unmarshal(reelComment.Data, &data)
	notification.ExtraData.ActionOnText = data.Text
	notification.ExtraData.DiscordMessage = append(notification.ExtraData.DiscordMessage, fmt.Sprintf("```%s```", notification.ExtraData.ActionOnText))
	notification.ExtraData.ReelUUID = reelComment.Reel.UUID.String()
	notification.ExtraData.ReelID = reelComment.Reel.ID
	notification.ExtraData.BlockUUID = reelComment.Reel.Block.UUID.String()
	notification.ExtraData.ReelCommentUUID = reelComment.UUID.String()
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "💬 View", notification.ExtraData.AppUrl)
	s.handleNotificationCreation(notification)
}

func (s notificationService) GetUserIDsFromMentions(mentions datatypes.JSON) ([]uint64, []uint64, []uint64) {
	mentionsData := []MentionsSerializer{}
	json.Unmarshal(mentions, &mentionsData)
	userIDs := []uint64{}
	roleIDS := []uint64{}
	canvasRepoIDs := []uint64{}
	if mentionsData == nil {
		return userIDs, roleIDS, canvasRepoIDs
	}
	for _, mention := range mentionsData {
		if mention.Type == "user" {
			userIDs = append(userIDs, mention.ID)
		} else if mention.Type == "role" {
			roleIDS = append(roleIDS, mention.ID)
		} else if mention.Type == "canvas" {
			canvasRepoIDs = append(canvasRepoIDs, mention.ID)
		}
	}
	return userIDs, roleIDS, canvasRepoIDs
}
