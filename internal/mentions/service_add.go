package mentions

import (
	"encoding/json"
	"errors"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
)

func (s mentionsService) AddMentionToBlock(obj MentionPost, userObjects *[]MentionedUserSerializer, canvasObjects *[]MentionedCanvasSerializer, roleObjects *[]MentionedRolesSerializer, user *models.User, studioID uint64) ([]map[string]interface{}, error) {
	var blockID = obj.ObjectUUID
	blockInstance, err := App.Repo.GetBlockByUUID(blockID)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	// permission check
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(user.ID, *blockInstance.CanvasBranchID, permissiongroup.CANVAS_BRANCH_EDIT); err != nil || !hasPermission {
		return nil, errors.New(response.NoPermissionError)
	}

	mappedData := s.MentionMapMaker(userObjects, canvasObjects, roleObjects, user, studioID)
	mappedDataJson, err := json.Marshal(mappedData)
	err = App.Repo.Manager.UpdateEntityByID(models.BLOCK, blockInstance.ID, map[string]interface{}{"mentions": mappedDataJson})

	go func() {
		extraData := notifications.NotificationExtraData{
			CanvasBranchID: obj.CanvasBranchID,
		}
		contentObject := models.BLOCK
		notifications.App.Service.PublishNewNotification(notifications.BlockMention, user.ID, nil, nil,
			nil, extraData, &blockInstance.ID, &contentObject)
	}()
	return mappedData, nil
}

func (s mentionsService) AddMentionToBlockThread(obj MentionPost, userObjects *[]MentionedUserSerializer, canvasObjects *[]MentionedCanvasSerializer, roleObjects *[]MentionedRolesSerializer, user *models.User, studioID uint64) ([]map[string]interface{}, error) {
	var thingID = obj.ObjectID
	blockThread, err := App.Repo.GetBlockThread(thingID)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	// permission check
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(user.ID, blockThread.CanvasBranchID, permissiongroup.CANVAS_BRANCH_ADD_COMMENT); err != nil || !hasPermission {
		return nil, errors.New(response.NoPermissionError)
	}

	mappedData := s.MentionMapMaker(userObjects, canvasObjects, roleObjects, user, studioID)
	mappedDataJson, err := json.Marshal(mappedData)
	err = App.Repo.Manager.UpdateEntityByID(models.BLOCK_THREAD, thingID, map[string]interface{}{"mentions": mappedDataJson})

	go func() {
		blockInstance, _ := App.Repo.GetBlock(blockThread.StartBlockID)
		extraData := notifications.NotificationExtraData{
			CanvasRepoID:   blockInstance.CanvasRepositoryID,
			CanvasBranchID: *blockInstance.CanvasBranchID,
		}
		contentObject := models.BLOCK_THREAD
		notifications.App.Service.PublishNewNotification(notifications.BlockThreadMention, user.ID, nil, nil,
			nil, extraData, &thingID, &contentObject)
	}()

	return mappedData, nil
}

func (s mentionsService) AddMentionToBlockThreadComment(obj MentionPost, userObjects *[]MentionedUserSerializer, canvasObjects *[]MentionedCanvasSerializer, roleObjects *[]MentionedRolesSerializer, user *models.User, studioID uint64) ([]map[string]interface{}, error) {
	var thingID = obj.ObjectID
	blockComment, err := App.Repo.GetBlockComment(thingID)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	// permission check
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(user.ID, blockComment.Thread.CanvasBranchID, permissiongroup.CANVAS_BRANCH_ADD_COMMENT); err != nil || !hasPermission {
		return nil, errors.New(response.NoPermissionError)
	}

	mappedData := s.MentionMapMaker(userObjects, canvasObjects, roleObjects, user, studioID)
	mappedDataJson, err := json.Marshal(mappedData)
	err = App.Repo.Manager.UpdateEntityByID(models.BLOCK_THREAD_COMMENT, thingID, map[string]interface{}{"mentions": mappedDataJson})

	go func() {
		blockInstance, _ := App.Repo.GetBlock(blockComment.Thread.StartBlockID)
		extraData := notifications.NotificationExtraData{
			CanvasRepoID:   blockInstance.CanvasRepositoryID,
			CanvasBranchID: *blockInstance.CanvasBranchID,
		}
		contentObject := models.BLOCK_THREAD_COMMENT
		notifications.App.Service.PublishNewNotification(notifications.BlockThreadCommentMention, user.ID, nil, nil,
			nil, extraData, &thingID, &contentObject)
	}()

	return mappedData, nil
}

func (s mentionsService) AddMentionToReel(obj MentionPost, userObjects *[]MentionedUserSerializer, canvasObjects *[]MentionedCanvasSerializer, roleObjects *[]MentionedRolesSerializer, user *models.User, studioID uint64) ([]map[string]interface{}, error) {
	var thingID = obj.ObjectID
	reelInstance, err := App.Repo.GetReel(thingID)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	// permission check
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(user.ID, reelInstance.CanvasBranchID, permissiongroup.CANVAS_BRANCH_CREATE_REEL); err != nil || !hasPermission {
		return nil, errors.New(response.NoPermissionError)
	}

	mappedData := s.MentionMapMaker(userObjects, canvasObjects, roleObjects, user, studioID)
	mappedDataJson, err := json.Marshal(mappedData)
	err = App.Repo.Manager.UpdateEntityByID(models.REEL, thingID, map[string]interface{}{"mentions": mappedDataJson})

	go func() {
		reelInstance, _ := App.Repo.GetReel(obj.ObjectID)
		extraData := notifications.NotificationExtraData{
			CanvasRepoID:   reelInstance.CanvasRepositoryID,
			CanvasBranchID: reelInstance.CanvasBranchID,
		}
		contentObject := models.REEL
		notifications.App.Service.PublishNewNotification(notifications.ReelMention, user.ID, nil, nil,
			nil, extraData, &thingID, &contentObject)
	}()

	return mappedData, nil
}

func (s mentionsService) AddMentionToReelComment(obj MentionPost, userObjects *[]MentionedUserSerializer, canvasObjects *[]MentionedCanvasSerializer, roleObjects *[]MentionedRolesSerializer, user *models.User, studioID uint64) ([]map[string]interface{}, error) {
	var thingID = obj.ObjectID
	reelComment, err := App.Repo.GetReelComment(thingID)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	// permission check
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(user.ID, reelComment.Reel.CanvasBranchID, permissiongroup.CANVAS_BRANCH_COMMENT_ON_REEL); err != nil || !hasPermission {
		return nil, errors.New(response.NoPermissionError)
	}

	mappedData := s.MentionMapMaker(userObjects, canvasObjects, roleObjects, user, studioID)
	mappedDataJson, err := json.Marshal(mappedData)
	err = App.Repo.Manager.UpdateEntityByID(models.REEL_COMMENTS, thingID, map[string]interface{}{"mentions": mappedDataJson})

	go func() {
		reelInstance, _ := App.Repo.GetReel(reelComment.ReelID)
		extraData := notifications.NotificationExtraData{
			CanvasRepoID:   reelInstance.CanvasRepositoryID,
			CanvasBranchID: reelInstance.CanvasBranchID,
		}
		contentObject := models.REEL_COMMENTS
		notifications.App.Service.PublishNewNotification(notifications.ReelComment, user.ID, nil, nil,
			nil, extraData, &thingID, &contentObject)
	}()
	return mappedData, nil
}
