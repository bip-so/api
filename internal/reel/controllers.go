package reel

import (
	"encoding/json"

	"gitlab.com/phonepost/bip-be-platform/internal/feed"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/internal/reactions"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/apiClient"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

//var AllowedScope = []string{"", "create", "update", "delete"}

func (c reelController) CreateReel(reqData NewReelCreatePOST, studioID uint64, userID uint64) (*models.Reel, error) {
	blockInstance, err := App.Repo.GetBlockByUUIDAndBranchID(map[string]interface{}{"uuid": reqData.StartBlockUUID, "canvas_branch_id": reqData.CanvasBranchID})
	if err != nil {
		return nil, err
	}
	var reel models.Reel
	reel.StudioID = studioID
	reel.CanvasRepositoryID = reqData.CanvasRepositoryID
	reel.CanvasBranchID = reqData.CanvasBranchID
	reel.StartBlockID = blockInstance.ID
	reel.StartBlockUUID = reqData.StartBlockUUID
	reel.TextRangeStart = reqData.TextRangeStart
	reel.TextRangeEnd = reqData.TextRangeEnd
	reel.RangeStart = reqData.RangeStart
	reel.RangeEnd = reqData.RangeEnd
	reel.HighlightedText = reqData.HighlightedText
	reel.ContextData = reqData.ContextData
	reel.CreatedByID = userID
	reel.UpdatedByID = userID
	reel.AuthorID = userID
	reel.SelectedBlocks = reqData.SelectedBlocks
	created, errCreatng := App.Repo.Create(reel)
	if errCreatng != nil {
		return nil, errCreatng
	}
	// @todo move this kafka later
	go func() {
		feed.App.Service.AddReelActivity(created)
		//notifications.App.Service.PublishIntegrationEvent(created, models.REEL)
		reelStr, _ := json.Marshal(created)
		apiClient.AddToQueue(apiClient.SendToIntegration, reelStr, apiClient.DEFAULT, apiClient.CommonRetry)
		apiClient.AddToQueue(apiClient.AddReelToAlgolia, reelStr, apiClient.DEFAULT, apiClient.CommonRetry)
	}()
	return created, nil
}

func (c reelController) CreateReelComment(reqData ReelCommentCreatePOST, studioID uint64, userID uint64, reelID uint64) (*models.ReelComment, error) {
	var reelComment *models.ReelComment
	cb := models.CommentBase{
		0,
		reqData.Data,
		reqData.IsEdited,
		reqData.IsReply,
	}
	reelComment, err := App.Repo.CreateReelComment(reqData, cb, reelID, userID)
	if err != nil {
		return nil, err
	}

	go func() {
		reel, _ := App.Repo.GetReel(map[string]interface{}{"id": reelComment.ReelID})
		contentObject := models.REEL_COMMENTS
		extraData := notifications.NotificationExtraData{
			CanvasRepoID:   reel.CanvasRepositoryID,
			CanvasBranchID: reel.CanvasBranchID,
		}
		if reqData.ParentID == nil {
			notifications.App.Service.PublishNewNotification(notifications.ReelComment, userID, nil, &studioID, nil, extraData, &reelComment.ID, &contentObject)
		} else {
			notifications.App.Service.PublishNewNotification(notifications.ReelCommentReply, userID, nil, &studioID, nil, extraData, &reelComment.ID, &contentObject)
		}
	}()

	return reelComment, nil
}

func (c reelController) DeleteReel(reelID uint64, userID uint64) error {
	reel, err := App.Repo.GetReel(map[string]interface{}{"id": reelID})
	if err != nil {
		return err
	}
	err = App.Repo.Delete(reelID, userID)
	if err != nil {
		return err
	}

	// removing from feed activity
	go func() {
		feed.App.Service.RemoveReelActivity(reel, userID)
		reelStr, _ := json.Marshal(reel)
		apiClient.AddToQueue(apiClient.DeleteReelFromAlgolia, reelStr, apiClient.DEFAULT, apiClient.CommonRetry)
	}()
	return err
}

func (c reelController) GetReelsFeed(user *models.User, offset int, limit int) (*[]ReelsSerialData, error) {
	reelsActivityFeed, err := feed.App.Service.GetReelActivities(user.ID, offset, limit)
	if err != nil {
		return nil, err
	}

	var reelIDs []uint64
	for _, reelActivity := range reelsActivityFeed.Results {
		if reelActivity.ForeignID != "" {
			reelIDs = append(reelIDs, utils.Uint64(reelActivity.ForeignID))
		}
	}
	reels, err := App.Repo.GetReelsByIDs(reelIDs)
	if err != nil {
		return nil, err
	}

	reelReactions := []models.ReelReaction{}
	var userFollowings *[]models.FollowUser
	var members []models.Member
	if user != nil {
		reelReactions, _ = reactions.App.Repo.GetUserReelReactionByIDs(reelIDs, user.ID)
		userFollowings, _ = App.Repo.GetUserFollowings(user.ID)
		members, _ = App.Repo.GetMembersByUserID(user.ID)
	}
	reelsFeed := SerializeDefaultManyReelsWithReactionsForUser(reels, reelReactions, user, members, userFollowings)
	return reelsFeed, err
}

func (c reelController) GetReelsFeedForStudio(user *models.User, studioID uint64, offset int, limit int) (*[]ReelsSerialData, error) {
	reels := []models.Reel{}
	reels, err := App.Service.StudioUserReelsFromDb(user, studioID, offset, limit)
	if err != nil {
		return nil, err
	}
	var reelIDs []uint64
	for _, reel := range reels {
		reelIDs = append(reelIDs, reel.ID)
	}

	reelReactions := []models.ReelReaction{}
	var userFollowings *[]models.FollowUser
	var members []models.Member
	if user != nil {
		reelReactions, _ = reactions.App.Repo.GetUserReelReactionByIDs(reelIDs, user.ID)
		userFollowings, _ = App.Repo.GetUserFollowings(user.ID)
		members, _ = App.Repo.GetMembersByUserID(user.ID)
	}
	reelsFeed := SerializeDefaultManyReelsWithReactionsForUser(&reels, reelReactions, user, members, userFollowings)
	return reelsFeed, err
}

func (c reelController) DeleteReelComment(reelComment *models.ReelComment) error {

	errDeleting := App.Repo.Manager.HardDeleteByID(models.REEL_COMMENTS, reelComment.ID)
	if errDeleting != nil {
		return errDeleting
	}
	// Update comment count on reel when a comment is deleted.
	if reelComment.ParentID != nil {
		_ = App.Repo.Manager.CommentCountMinus(models.REEL_COMMENTS, *reelComment.ParentID)
	} else {
		_ = App.Repo.Manager.CommentCountMinus(models.REEL, reelComment.ReelID)
	}
	return nil
}
