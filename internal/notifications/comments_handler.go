package notifications

import (
	"encoding/json"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

type CommentData struct {
	Text string `json:"text"`
}

func (s notificationService) BlockCommentHandler(notification *PostNotification) {
	// Changing the canvasBranch id to main. if it is a rough branch
	canvasBranch, _ := App.Repo.GetCanvasBranchByID(notification.ExtraData.CanvasBranchID)
	if canvasBranch.IsRoughBranch {
		return
	}
	event := AllCommentsEntity.Events[notification.Event]
	notification.Entity = AllComments
	notification.Activity = event.Activity

	userInstance, _ := App.Repo.GetUser(notification.CreatedByID)
	canvasRepo, _ := App.Repo.GetCanvasRepoByID(notification.ExtraData.CanvasRepoID)
	notification.Text = fmt.Sprintf(event.Text, userInstance.Username, canvasRepo.Name)

	notification.Priority = event.Priority
	notification.IsPersonal = AllCommentsEntity.IsPersonal

	blockThread, _ := App.Repo.GetBlockThreadByID(*notification.ObjectID)
	notification.StudioID = &blockThread.CanvasRepository.StudioID

	blockInstance, _ := App.Repo.GetBlockByID(blockThread.StartBlockID)
	var user models.User
	if blockInstance.UpdatedByID != 0 {
		user, _ = App.Repo.GetUser(blockInstance.UpdatedByID)
	} else {
		user, _ = App.Repo.GetUser(blockInstance.CreatedByID)
	}

	userIds := []uint64{user.ID}
	// GEt moderator userIDs of the block branch.
	branchModUserIDs, _ := App.Repo.GetCanvasBranchModeratorsUserIDs(blockThread.CanvasBranchID)
	userIds = append(userIds, branchModUserIDs...)

	if blockThread.Mentions != nil {
		mentionUserIDs, roleIds, _ := s.GetUserIDsFromMentions(*blockThread.Mentions)
		roleUserIDs := s.GetUserIdsFromRoleIds(&roleIds)
		mentionUserIDs = append(mentionUserIDs, roleUserIDs...)
		userIds = s.RemoveMentionedUserIDs(userIds, mentionUserIDs)
	}

	notification.NotifierIDs = s.GetUniqueIDs(userIds)

	// Extra Data
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
	blockText := App.Service.GetBlockText(&blockInstance)
	notification.ExtraData.DiscordMessage = append(notification.ExtraData.DiscordMessage, blockText)
	notification.ExtraData.DiscordMessage = append(notification.ExtraData.DiscordMessage, fmt.Sprintf("```%s```", blockThread.Text))
	notification.ExtraData.BlockUUID = blockInstance.UUID.String()
	notification.ExtraData.BlockThreadUUID = blockThread.UUID.String()
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "💬 View", notification.ExtraData.AppUrl)
	s.handleNotificationCreation(notification)
}

func (s notificationService) CommentReplyHandler(notification *PostNotification) {
	event := RepliesToMeEntity.Events[notification.Event]
	notification.Entity = RepliesToMe
	notification.Activity = event.Activity

	userInstance, _ := App.Repo.GetUser(notification.CreatedByID)
	canvasRepo, _ := App.Repo.GetCanvasRepoByID(notification.ExtraData.CanvasRepoID)
	notification.Text = fmt.Sprintf(event.Text, userInstance.Username, canvasRepo.Name)

	notification.Priority = event.Priority
	notification.IsPersonal = RepliesToMeEntity.IsPersonal

	// Notifiers
	// person on whose reply this action is done
	// person who previously replied on the thread
	// person who mentioned in the reply_to message or anywhere in the thread
	// Mod of the canvas
	//
	// we have commentId -> parentCommentId
	// GET all comments for that thread and -> list UserIDs
	// thread has canvasRepo -> get mods of the canvasRepo Default Main Branch
	// thread mentions -> list userIds
	blockComment, _ := App.Repo.GetBlockThreadCommentByID(*notification.ObjectID)
	blockThreadComments, _ := App.Repo.GetBlockCommentsByThreadID(blockComment.ThreadID)
	userIDs := []uint64{}
	roleIDs := []uint64{}
	mentionUserIDs := []uint64{}
	mentionRoleIDs := []uint64{}
	for _, comment := range blockThreadComments {
		userIDs = append(userIDs, comment.CreatedByID)
		if comment.Mentions != nil {
			mentionUserIDs, mentionRoleIDs, _ = s.GetUserIDsFromMentions(*comment.Mentions)
			userIDs = append(userIDs, mentionUserIDs...)
			userIDs = append(userIDs, mentionRoleIDs...)
		}
	}
	branchModUserIDs, _ := App.Repo.GetCanvasBranchModeratorsUserIDs(*blockComment.Thread.CanvasRepository.DefaultBranchID)
	userIDs = append(userIDs, branchModUserIDs...)
	if blockComment.Thread.Mentions != nil {
		mentionUserIDs, mentionRoleIDs, _ = s.GetUserIDsFromMentions(*blockComment.Thread.Mentions)
		userIDs = append(userIDs, mentionUserIDs...)
		roleIDs = append(roleIDs, mentionRoleIDs...)
	}

	if blockComment.Mentions != nil {
		mentionUserIDs, roleIds, _ := s.GetUserIDsFromMentions(*blockComment.Mentions)
		roleUserIDs := s.GetUserIdsFromRoleIds(&roleIds)
		mentionUserIDs = append(mentionUserIDs, roleUserIDs...)
		modMentionUserIDs := s.GetModMentionedUserIDs(branchModUserIDs, mentionUserIDs)
		userIDs = s.RemoveMentionedUserIDs(userIDs, modMentionUserIDs)
	}

	notification.NotifierIDs = s.GetUniqueIDs(userIDs)
	uniqueRoleIDs := s.GetUniqueIDs(roleIDs)
	notification.RoleIDs = &uniqueRoleIDs

	// Extra Data
	notification.ExtraData.CanvasRepoID = blockComment.Thread.CanvasRepositoryID
	notification.ExtraData.CollectionID = blockComment.Thread.CanvasRepository.CollectionID
	notification.ExtraData.AppUrl = s.GenerateBlockCommentUrl(canvasRepo.Key, canvasRepo.Name, blockComment.Thread.UUID.String(), canvasRepo.StudioID, *canvasRepo.DefaultBranchID)
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
	var blockCommentData CommentData
	json.Unmarshal(blockComment.Data, &blockCommentData)
	notification.ExtraData.DiscordMessage = append(notification.ExtraData.DiscordMessage, fmt.Sprintf("```%s```", blockCommentData.Text))
	notification.ExtraData.BlockUUID = blockComment.Thread.Block.UUID.String()
	notification.ExtraData.BlockThreadUUID = blockComment.Thread.UUID.String()
	notification.ExtraData.BlockThreadCommentUUID = blockComment.UUID.String()
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "💬 View", notification.ExtraData.AppUrl)
	// Changing the canvasBranch id to main. if it is a rough branch
	canvasBranch, _ := App.Repo.GetCanvasBranchByID(notification.ExtraData.CanvasBranchID)
	if canvasBranch.IsRoughBranch {
		roughBranchID := notification.ExtraData.CanvasBranchID
		notification.ExtraData.CanvasBranchID = *canvasRepo.DefaultBranchID
		s.AddNotificationToRedis(roughBranchID, notification)
		return
	}
	s.handleNotificationCreation(notification)
}

func (s notificationService) ReelCommentHandler(notification *PostNotification) {
	event := AllCommentsEntity.Events[notification.Event]
	notification.Entity = AllComments
	notification.Activity = event.Activity
	notification.Priority = event.Priority
	notification.IsPersonal = AllCommentsEntity.IsPersonal

	reelComment, _ := App.Repo.GetReelCommentByID(*notification.ObjectID)

	userIds := []uint64{reelComment.Reel.CreatedByID}
	// GEt moderator userIDs of the block.
	branchModUserIDs, _ := App.Repo.GetCanvasBranchModeratorsUserIDs(reelComment.Reel.CanvasBranchID)
	userIds = append(userIds, branchModUserIDs...)

	if reelComment.Mentions != nil {
		mentionUserIDs, roleIds, _ := s.GetUserIDsFromMentions(*reelComment.Mentions)
		roleUserIDs := s.GetUserIdsFromRoleIds(&roleIds)
		mentionUserIDs = append(mentionUserIDs, roleUserIDs...)
		userIds = s.RemoveMentionedUserIDs(userIds, mentionUserIDs)
	}

	userInstance, _ := App.Repo.GetUser(notification.CreatedByID)
	notification.Text = fmt.Sprintf(event.Text, userInstance.Username)
	notification.NotifierIDs = s.GetUniqueIDs(userIds)

	canvasRepo, _ := App.Repo.GetCanvasRepoByID(notification.ExtraData.CanvasRepoID)
	notification.ExtraData.AppUrl = s.GenerateReelCommentUrl(canvasRepo.Key, canvasRepo.Name, canvasRepo.StudioID, *canvasRepo.DefaultBranchID, reelComment.Reel.UUID.String())
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
	var reelCommentData CommentData
	json.Unmarshal(reelComment.Data, &reelCommentData)
	notification.ExtraData.DiscordMessage = append(notification.ExtraData.DiscordMessage, fmt.Sprintf("```%s```", reelCommentData.Text))
	notification.ExtraData.BlockUUID = reelComment.Reel.Block.UUID.String()
	notification.ExtraData.ReelUUID = reelComment.Reel.UUID.String()
	notification.ExtraData.ReelID = reelComment.Reel.ID
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "💬 View", notification.ExtraData.AppUrl)
	// Changing the canvasBranch id to main. if it is a rough branch
	canvasBranch, _ := App.Repo.GetCanvasBranchByID(notification.ExtraData.CanvasBranchID)
	if canvasBranch.IsRoughBranch {
		roughBranchID := notification.ExtraData.CanvasBranchID
		notification.ExtraData.CanvasBranchID = *canvasRepo.DefaultBranchID
		s.AddNotificationToRedis(roughBranchID, notification)
		return
	}
	s.handleNotificationCreation(notification)
}

func (s notificationService) ReelCommentReplyHandler(notification *PostNotification) {
	event := RepliesToMeEntity.Events[notification.Event]
	notification.Entity = RepliesToMe
	notification.Activity = event.Activity
	userInstance, _ := App.Repo.GetUser(notification.CreatedByID)
	notification.Text = fmt.Sprintf(event.Text, userInstance.Username)
	notification.Priority = event.Priority
	notification.IsPersonal = RepliesToMeEntity.IsPersonal

	// Notifiers
	// person on whose reply this action is done
	// person who previously replied on the thread
	// person who mentioned in the reply_to message or anywhere in the thread
	// Mod of the canvas
	//
	// We have ReelCommentId -> parentId
	// we get comments based on parentId
	// list userIDs
	// Get Reel comments mentions in all comments of thread -> add to list userIDs
	reelComment, _ := App.Repo.GetReelCommentByID(*notification.ObjectID)
	allReelComments, _ := App.Repo.GetReelComments(map[string]interface{}{"parent_id": *reelComment.ParentID})
	userIDs := []uint64{}
	roleIDs := []uint64{}
	for _, comment := range allReelComments {
		userIDs = append(userIDs, comment.CreatedByID)
		if comment.Mentions != nil {
			mentionUserIDs, mentionRoleIDs, _ := s.GetUserIDsFromMentions(*comment.Mentions)
			userIDs = append(userIDs, mentionUserIDs...)
			roleIDs = append(roleIDs, mentionRoleIDs...)
		}
	}

	if reelComment.Mentions != nil {
		branchModUserIDs, _ := App.Repo.GetCanvasBranchModeratorsUserIDs(reelComment.Reel.CanvasBranchID)
		mentionUserIDs, roleIds, _ := s.GetUserIDsFromMentions(*reelComment.Mentions)
		roleUserIDs := s.GetUserIdsFromRoleIds(&roleIds)
		mentionUserIDs = append(mentionUserIDs, roleUserIDs...)
		modMentionUserIDs := s.GetModMentionedUserIDs(branchModUserIDs, mentionUserIDs)
		userIDs = s.RemoveMentionedUserIDs(userIDs, modMentionUserIDs)
	}

	notification.NotifierIDs = s.GetUniqueIDs(userIDs)
	roleIDs = s.GetUniqueIDs(roleIDs)
	notification.RoleIDs = &roleIDs

	// Extra Data
	notification.ExtraData.CanvasRepoID = reelComment.Reel.CanvasRepositoryID
	canvasRepo, _ := App.Repo.GetCanvasRepoByID(notification.ExtraData.CanvasRepoID)
	notification.ExtraData.AppUrl = s.GenerateReelCommentUrl(canvasRepo.Key, canvasRepo.Name, canvasRepo.StudioID, *canvasRepo.DefaultBranchID, reelComment.Reel.UUID.String())
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
	var reelCommentData CommentData
	json.Unmarshal(reelComment.Data, &reelCommentData)
	notification.ExtraData.DiscordMessage = append(notification.ExtraData.DiscordMessage, fmt.Sprintf("```%s```", reelCommentData.Text))
	notification.ExtraData.ReelCommentUUID = reelComment.UUID.String()
	notification.ExtraData.ReelUUID = reelComment.Reel.UUID.String()
	notification.ExtraData.ReelID = reelComment.Reel.ID
	notification.ExtraData.BlockUUID = reelComment.Reel.Block.UUID.String()
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "💬 View", notification.ExtraData.AppUrl)
	// Changing the canvasBranch id to main. if it is a rough branch
	canvasBranch, _ := App.Repo.GetCanvasBranchByID(notification.ExtraData.CanvasBranchID)
	if canvasBranch.IsRoughBranch {
		roughBranchID := notification.ExtraData.CanvasBranchID
		notification.ExtraData.CanvasBranchID = *canvasRepo.DefaultBranchID
		s.AddNotificationToRedis(roughBranchID, notification)
		return
	}
	s.handleNotificationCreation(notification)
}

func (s notificationService) BlockReactHandler(notification *PostNotification) {
	event := ReactionsEntity.Events[notification.Event]
	notification.Entity = Reactions
	notification.Activity = event.Activity
	notification.Priority = event.Priority
	notification.IsPersonal = ReactionsEntity.IsPersonal

	blockReaction, _ := App.Repo.GetBlockReactionByID(*notification.ObjectID)
	notification.StudioID = &blockReaction.CanvasBranch.CanvasRepository.StudioID

	blockInstance, _ := App.Repo.GetBlockByID(blockReaction.BlockID)
	user, _ := App.Repo.GetUser(blockInstance.UpdatedByID)

	// Changing the canvasBranch id to main. if it is a rough branch
	canvasRepo, _ := App.Repo.GetCanvasRepoByID(notification.ExtraData.CanvasRepoID)
	canvasBranch, _ := App.Repo.GetCanvasBranchByID(notification.ExtraData.CanvasBranchID)
	if canvasBranch.IsRoughBranch {
		roughBranchID := notification.ExtraData.CanvasBranchID
		notification.ExtraData.CanvasBranchID = *canvasRepo.DefaultBranchID
		s.AddNotificationToRedis(roughBranchID, notification)
		return
	}

	userIds := []uint64{user.ID}
	// GEt moderator userIDs of the block branch.
	branchModUserIDs, _ := App.Repo.GetCanvasBranchModeratorsUserIDs(*blockReaction.CanvasBranchID)
	userIds = append(userIds, branchModUserIDs...)

	notification.NotifierIDs = s.GetUniqueIDs(userIds)
	// Extra Data
	notification.ExtraData.CanvasRepoID = blockReaction.CanvasBranch.CanvasRepositoryID
	notification.ExtraData.CollectionID = blockReaction.CanvasBranch.CanvasRepository.CollectionID

	notification.ExtraData.AppUrl = s.GenerateBlockReactionUrl(canvasRepo.Key, canvasRepo.Name, blockInstance.UUID.String(), canvasRepo.StudioID, *canvasRepo.DefaultBranchID)
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

	userInstance, _ := App.Repo.GetUser(notification.CreatedByID)
	notification.Text = fmt.Sprintf(event.Text, userInstance.Username, blockReaction.Emoji, canvasRepo.Name)
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	var blockChildren []BlockChildren
	json.Unmarshal(blockInstance.Children, &blockChildren)
	var blockText string
	for _, children := range blockChildren {
		blockText = fmt.Sprintf("%s", children.Text)
	}
	notification.ExtraData.DiscordMessage = append(notification.ExtraData.DiscordMessage, fmt.Sprintf("```%s```", blockText))
	notification.ExtraData.BlockUUID = blockInstance.UUID.String()
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "💬 View", notification.ExtraData.AppUrl)
	s.handleNotificationCreation(notification)
}

func (s notificationService) ReelReactHandler(notification *PostNotification) {
	// Changing the canvasBranch id to main. if it is a rough branch
	canvasBranch, _ := App.Repo.GetCanvasBranchByID(notification.ExtraData.CanvasBranchID)
	canvasRepo, _ := App.Repo.GetCanvasRepoByID(notification.ExtraData.CanvasRepoID)
	if canvasBranch.IsRoughBranch {
		roughBranchID := notification.ExtraData.CanvasBranchID
		notification.ExtraData.CanvasBranchID = *canvasRepo.DefaultBranchID
		s.AddNotificationToRedis(roughBranchID, notification)
		return
	}

	event := ReactionsEntity.Events[notification.Event]
	notification.Entity = Reactions
	notification.Activity = event.Activity
	notification.Text = event.Text
	notification.Priority = event.Priority
	notification.IsPersonal = ReactionsEntity.IsPersonal

	reelReaction, _ := App.Repo.GetReelReactionByID(*notification.ObjectID)

	userIds := []uint64{reelReaction.Reel.CreatedByID}
	// GEt moderator userIDs of the block.
	branchModUserIDs, _ := App.Repo.GetCanvasBranchModeratorsUserIDs(reelReaction.Reel.CanvasBranchID)
	userIds = append(userIds, branchModUserIDs...)

	notification.NotifierIDs = s.GetUniqueIDs(userIds)

	// Extra data
	notification.ExtraData.CanvasRepoID = reelReaction.Reel.CanvasRepositoryID
	notification.ExtraData.AppUrl = s.GenerateReelCommentUrl(canvasRepo.Key, canvasRepo.Name, canvasRepo.StudioID, *canvasRepo.DefaultBranchID, reelReaction.Reel.UUID.String())
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
	userInstance, _ := App.Repo.GetUser(notification.CreatedByID)
	notification.Text = fmt.Sprintf(event.Text, userInstance.Username, reelReaction.Emoji, canvasRepo.Name)
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	var data CommentsData
	json.Unmarshal(reelReaction.Reel.ContextData, &data)
	notification.ExtraData.DiscordMessage = append(notification.ExtraData.DiscordMessage, fmt.Sprintf("```%s```", data.Text))
	notification.ExtraData.ReelUUID = reelReaction.Reel.UUID.String()
	notification.ExtraData.ReelID = reelReaction.Reel.ID
	notification.ExtraData.BlockUUID = reelReaction.Reel.Block.UUID.String()
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "💬 View", notification.ExtraData.AppUrl)
	s.handleNotificationCreation(notification)
}

func (s notificationService) BlockCommentReactHandler(notification *PostNotification) {
	// Changing the canvasBranch id to main. if it is a rough branch
	canvasBranch, _ := App.Repo.GetCanvasBranchByID(notification.ExtraData.CanvasBranchID)
	canvasRepo, _ := App.Repo.GetCanvasRepoByID(notification.ExtraData.CanvasRepoID)
	if canvasBranch.IsRoughBranch {
		roughBranchID := notification.ExtraData.CanvasBranchID
		notification.ExtraData.CanvasBranchID = *canvasRepo.DefaultBranchID
		s.AddNotificationToRedis(roughBranchID, notification)
		return
	}
	event := ReactionsEntity.Events[notification.Event]
	notification.Entity = Reactions
	notification.Activity = event.Activity
	notification.Text = event.Text
	notification.Priority = event.Priority
	notification.IsPersonal = ReactionsEntity.IsPersonal

	// Notifiers
	// person on whose reply this action is done
	// person who previously replied on the thread
	// person who mentioned in the reply_to message or anywhere in the thread
	// Mod of the canvas
	//
	// BlockThreadReaction
	// BlockThreadCommentReaction
	//
	// we have commentId -> threadId
	// GET all comments for that thread and -> list UserIDs
	// thread has canvasRepo -> get mods of the canvasRepo Default Main Branch
	// thread mentions -> list userIds
	var blockThread models.BlockThread
	var reactEmoji string
	var commentText string
	switch *notification.ContentObject {
	case models.BLOCKTHREADREACTION:
		blockThreadReaction, _ := App.Repo.GetBlockThreadReactionByID(*notification.ObjectID)
		blockThread = *blockThreadReaction.BlockThread
		reactEmoji = blockThreadReaction.Emoji
		commentText = blockThreadReaction.BlockThread.Text
	case models.BLOCKCOMMENTREACTION:
		blockCommentReaction, _ := App.Repo.GetBlockThreadCommentReactionByID(*notification.ObjectID)
		blockThread = *blockCommentReaction.BlockThread
		reactEmoji = blockCommentReaction.Emoji
		var data CommentsData
		json.Unmarshal(blockCommentReaction.BlockComment.Data, &data)
		commentText = data.Text
	}
	userIDs := []uint64{blockThread.CreatedByID}
	mentionUserIDs := []uint64{}
	mentionRoleIDs := []uint64{}
	blockThreadComments, _ := App.Repo.GetBlockCommentsByThreadID(blockThread.ID)
	roleIDs := []uint64{}
	for _, comment := range blockThreadComments {
		userIDs = append(userIDs, comment.CreatedByID)
		if comment.Mentions != nil {
			mentionUserIDs, mentionRoleIDs, _ = s.GetUserIDsFromMentions(*comment.Mentions)
			userIDs = append(userIDs, mentionUserIDs...)
			roleIDs = append(roleIDs, mentionRoleIDs...)
		}
	}
	branchModUserIDs, _ := App.Repo.GetCanvasBranchModeratorsUserIDs(*blockThread.CanvasRepository.DefaultBranchID)
	userIDs = append(userIDs, branchModUserIDs...)

	if blockThread.Mentions != nil {
		mentionUserIDs, mentionRoleIDs, _ = s.GetUserIDsFromMentions(*blockThread.Mentions)
		userIDs = append(userIDs, mentionUserIDs...)
		roleIDs = append(roleIDs, mentionRoleIDs...)
	}

	notification.NotifierIDs = s.GetUniqueIDs(userIDs)
	roleIDs = s.GetUniqueIDs(roleIDs)
	notification.RoleIDs = &roleIDs

	// Extra data
	notification.ExtraData.CanvasRepoID = blockThread.CanvasRepositoryID

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
	userInstance, _ := App.Repo.GetUser(notification.CreatedByID)
	notification.Text = fmt.Sprintf(event.Text, userInstance.Username, reactEmoji, canvasRepo.Name)
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	notification.ExtraData.DiscordMessage = append(notification.ExtraData.DiscordMessage, fmt.Sprintf("```%s```", commentText))
	notification.ExtraData.BlockUUID = blockThread.Block.UUID.String()
	notification.ExtraData.BlockThreadUUID = blockThread.UUID.String()
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "💬 View", notification.ExtraData.AppUrl)
	s.handleNotificationCreation(notification)
}

func (s notificationService) ReelCommentReactHandler(notification *PostNotification) {
	// Changing the canvasBranch id to main. if it is a rough branch
	canvasBranch, _ := App.Repo.GetCanvasBranchByID(notification.ExtraData.CanvasBranchID)
	canvasRepo, _ := App.Repo.GetCanvasRepoByID(notification.ExtraData.CanvasRepoID)
	if canvasBranch.IsRoughBranch {
		//roughBranchID := notification.ExtraData.CanvasBranchID
		//notification.ExtraData.CanvasBranchID = *canvasRepo.DefaultBranchID
		//s.AddNotificationToRedis(roughBranchID, notification)
		return
	}
	event := ReactionsEntity.Events[notification.Event]
	notification.Entity = Reactions
	notification.Activity = event.Activity
	notification.Text = event.Text
	notification.Priority = event.Priority
	notification.IsPersonal = ReactionsEntity.IsPersonal
	// Notifiers
	// person on whose reply this action is done
	// person who previously replied on the thread
	// person who mentioned in the reply_to message or anywhere in the thread
	// Mod of the canvas
	//
	// We have ReelCommentId -> createdByID
	// if parentId is present we get reel comments based on parentID
	// list userIDs
	// Get Reel comments mentions -> add to list userIDs
	reelCommentReaction, _ := App.Repo.GetReelCommentReactionByID(*notification.ObjectID)
	if reelCommentReaction.ID == 0 {
		return
	}
	var allReelComments []models.ReelComment
	if reelCommentReaction.ReelComment.ParentID != nil {
		allReelComments, _ = App.Repo.GetReelComments(map[string]interface{}{"parent_id": *reelCommentReaction.ReelComment.ParentID})
	} else {
		allReelComments, _ = App.Repo.GetReelComments(map[string]interface{}{"reel_id": reelCommentReaction.ReelID})
	}

	userIDs := []uint64{}
	roleIDs := []uint64{}
	for _, comment := range allReelComments {
		userIDs = append(userIDs, comment.CreatedByID)
		if comment.Mentions != nil {
			mentionUserIDs, mentionRoleIDs, _ := s.GetUserIDsFromMentions(*comment.Mentions)
			userIDs = append(userIDs, mentionUserIDs...)
			roleIDs = append(roleIDs, mentionRoleIDs...)
		}
	}

	notification.NotifierIDs = s.GetUniqueIDs(userIDs)
	roleIDs = s.GetUniqueIDs(roleIDs)
	notification.RoleIDs = &roleIDs

	// Extra Data
	notification.ExtraData.CanvasRepoID = reelCommentReaction.Reel.CanvasRepositoryID

	notification.ExtraData.AppUrl = s.GenerateReelCommentUrl(canvasRepo.Key, canvasRepo.Name, canvasRepo.StudioID, *canvasRepo.DefaultBranchID, reelCommentReaction.Reel.UUID.String())
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
	userInstance, _ := App.Repo.GetUser(notification.CreatedByID)
	notification.Text = fmt.Sprintf(event.Text, userInstance.Username, reelCommentReaction.Emoji, canvasRepo.Name)
	notification.ExtraData.DiscordMessage = []string{notification.Text}
	var data CommentsData
	json.Unmarshal(reelCommentReaction.ReelComment.Data, &data)
	notification.ExtraData.DiscordMessage = append(notification.ExtraData.DiscordMessage, fmt.Sprintf("```%s```", data.Text))
	notification.ExtraData.ReelUUID = reelCommentReaction.Reel.UUID.String()
	notification.ExtraData.ReelID = reelCommentReaction.Reel.ID
	notification.ExtraData.BlockUUID = reelCommentReaction.Reel.Block.UUID.String()
	notification.ExtraData.SlackComponents = s.SlackNotificationBlockBuilder(notification.ExtraData.DiscordMessage, "💬 View", notification.ExtraData.AppUrl)
	s.handleNotificationCreation(notification)
}
