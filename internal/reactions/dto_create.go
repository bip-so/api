package reactions

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"strconv"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

func (r reactionRepo) GetBlockByUUIDAndBranchID(query map[string]interface{}) (*models.Block, error) {
	var block models.Block
	err := r.db.Model(&models.Block{}).Where(query).First(&block).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &block, nil
}

func (r reactionRepo) CreateBlockReaction(obj CreateMentionPost, studioID uint64, userID uint64) error {
	blockInstance, err := App.Repo.GetBlockByUUIDAndBranchID(map[string]interface{}{"uuid": obj.BlockUUID, "canvas_branch_id": obj.CanvasBranchID})
	if err != nil {
		return err
	}

	// Todo: permission check
	//if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(userID, *blockInstance.CanvasBranchID, permissiongroup.CANVAS_BRANCH_ADD_REACTION); err != nil || !hasPermission {
	//	return errors.New(response.NoPermissionError)
	//}

	var instance models.BlockReaction
	instance.Emoji = obj.Emoji
	if obj.CanvasBranchID != 0 {
		instance.CanvasBranchID = &obj.CanvasBranchID
	}
	instance.CreatedByID = userID
	instance.UpdatedByID = userID

	instance.BlockID = blockInstance.ID
	results := r.db.Create(&instance)
	// Block Reaction is created we need to update the "Reactions" now.
	query := "block_id = " + strconv.FormatUint(instance.BlockID, 10)
	counter, errGettingReactions := r.Manager.GetEmojiCounter(models.BLOCKREACTION, query)
	fmt.Println(counter)
	fmt.Println(errGettingReactions)
	// Block should get updated
	reactionjson, err := json.Marshal(counter)
	if err != nil {
		return err
	}
	err = App.Repo.Manager.UpdateEntityByID(models.BLOCK, instance.BlockID, map[string]interface{}{"reactions": reactionjson})
	if err != nil {
		return err
	}

	go func() {
		blockInstance, _ := App.Repo.GetBlockByID(instance.BlockID)
		extraData := notifications.NotificationExtraData{
			CanvasRepoID:   blockInstance.CanvasRepositoryID,
			CanvasBranchID: *blockInstance.CanvasBranchID,
		}
		contentObject := models.BLOCKREACTION
		notifications.App.Service.PublishNewNotification(notifications.BlockReact, userID, nil,
			&studioID, nil, extraData, &instance.ID, &contentObject)
	}()
	return results.Error
}
func (r reactionRepo) CreateBlockThreadReaction(obj CreateMentionPost, studioID uint64, userID uint64) error {
	blockInstance, err := App.Repo.GetBlockByUUIDAndBranchID(map[string]interface{}{"uuid": obj.BlockUUID, "canvas_branch_id": obj.CanvasBranchID})

	// permission check
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(userID, *blockInstance.CanvasBranchID, permissiongroup.CANVAS_BRANCH_ADD_REACTION); err != nil || !hasPermission {
		return errors.New(response.NoPermissionError)
	}

	var instance models.BlockThreadReaction
	instance.Emoji = obj.Emoji
	if obj.CanvasBranchID != 0 {
		instance.CanvasBranchID = &obj.CanvasBranchID
	}
	instance.CreatedByID = userID
	instance.UpdatedByID = userID
	instance.BlockID = blockInstance.ID
	instance.BlockThreadID = obj.BlockThreadID
	results := r.db.Create(&instance)

	query := "block_thread_id = " + strconv.FormatUint(instance.BlockThreadID, 10)
	counter, errGettingReactions := r.Manager.GetEmojiCounter(models.BLOCKTHREADREACTION, query)
	fmt.Println(counter)
	fmt.Println(errGettingReactions)
	reactionjson, err := json.Marshal(counter)
	if err != nil {
		return err
	}
	err = App.Repo.Manager.UpdateEntityByID(models.BLOCK_THREAD, instance.BlockThreadID, map[string]interface{}{"reactions": reactionjson})
	if err != nil {
		return err
	}

	go func() {
		blockInstance, _ := App.Repo.GetBlockByID(instance.BlockID)
		extraData := notifications.NotificationExtraData{
			CanvasRepoID:   blockInstance.CanvasRepositoryID,
			CanvasBranchID: *blockInstance.CanvasBranchID,
		}
		contentObject := models.BLOCKTHREADREACTION
		notifications.App.Service.PublishNewNotification(notifications.BlockCommentReact, userID, nil,
			&studioID, nil, extraData, &instance.ID, &contentObject)
	}()

	return results.Error
}
func (r reactionRepo) CreateBlockThreadCommentReaction(obj CreateMentionPost, studioID uint64, userID uint64) error {
	blockInstance, err := App.Repo.GetBlockByUUIDAndBranchID(map[string]interface{}{"uuid": obj.BlockUUID, "canvas_branch_id": obj.CanvasBranchID})

	// permission check
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(userID, *blockInstance.CanvasBranchID, permissiongroup.CANVAS_BRANCH_ADD_REACTION); err != nil || !hasPermission {
		return errors.New(response.NoPermissionError)
	}

	var instance models.BlockCommentReaction
	instance.Emoji = obj.Emoji
	if obj.CanvasBranchID != 0 {
		instance.CanvasBranchID = &obj.CanvasBranchID
	}
	instance.CreatedByID = userID
	instance.UpdatedByID = userID

	instance.BlockID = blockInstance.ID
	fmt.Println("Block ID", blockInstance.ID)
	instance.BlockThreadID = obj.BlockThreadID
	fmt.Println("Block Thread ID", obj.BlockThreadID)
	instance.BlockCommentID = obj.BlockCommentID
	fmt.Println("Block Comment ID", obj.BlockCommentID)
	results := r.db.Create(&instance)
	query := "block_comment_id = " + strconv.FormatUint(instance.BlockCommentID, 10)
	counter, errGettingReactions := r.Manager.GetEmojiCounter(models.BLOCKCOMMENTREACTION, query)
	fmt.Println(counter)
	fmt.Println(errGettingReactions)
	reactionjson, err := json.Marshal(counter)
	if err != nil {
		return err
	}
	err = App.Repo.Manager.UpdateEntityByID(models.BLOCK_THREAD_COMMENT, instance.BlockCommentID, map[string]interface{}{"reactions": reactionjson})
	if err != nil {
		return err
	}

	go func() {
		blockInstance2, _ := App.Repo.GetBlockByID(instance.BlockID)
		extraData := notifications.NotificationExtraData{
			CanvasRepoID:   blockInstance2.CanvasRepositoryID,
			CanvasBranchID: *blockInstance2.CanvasBranchID,
		}
		contentObject := models.BLOCKCOMMENTREACTION
		notifications.App.Service.PublishNewNotification(notifications.BlockCommentReact, userID, nil,
			&studioID, nil, extraData, &instance.ID, &contentObject)
	}()

	return results.Error
}
func (r reactionRepo) CreateReelReaction(obj CreateMentionPost, studioID uint64, userID uint64) error {
	reelInstance, _ := App.Repo.GetReelByID(obj.ReelID)

	// permission check
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(userID, reelInstance.CanvasBranchID, permissiongroup.CANVAS_BRANCH_REACT_TO_REEL); err != nil || !hasPermission {
		return errors.New(response.NoPermissionError)
	}

	var instance models.ReelReaction
	instance.Emoji = obj.Emoji
	if obj.CanvasBranchID != 0 {
		instance.CanvasBranchID = &obj.CanvasBranchID
	}
	instance.CreatedByID = userID
	instance.UpdatedByID = userID
	instance.ReelID = obj.ReelID
	results := r.db.Create(&instance)
	query := "reel_id = " + strconv.FormatUint(instance.ReelID, 10)
	counter, errGettingReactions := r.Manager.GetEmojiCounter(models.REELREACTION, query)
	fmt.Println(counter)
	fmt.Println(errGettingReactions)
	reactionjson, err := json.Marshal(counter)
	if err != nil {
		return err
	}
	err = App.Repo.Manager.UpdateEntityByID(models.REEL, instance.ReelID, map[string]interface{}{"reactions": reactionjson})
	if err != nil {
		return err
	}

	go func() {
		extraData := notifications.NotificationExtraData{
			CanvasRepoID:   reelInstance.CanvasRepositoryID,
			CanvasBranchID: reelInstance.CanvasBranchID,
		}
		contentObject := models.REELREACTION
		notifications.App.Service.PublishNewNotification(notifications.ReelReact, userID, nil,
			&studioID, nil, extraData, &instance.ID, &contentObject)
	}()

	return results.Error
}
func (r reactionRepo) CreateReelCommentReaction(obj CreateMentionPost, studioID uint64, userID uint64) error {
	reelInstance, _ := App.Repo.GetReelByID(obj.ReelID)
	// permission check
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(userID, reelInstance.CanvasBranchID, permissiongroup.CANVAS_BRANCH_REACT_TO_REEL); err != nil || !hasPermission {
		return errors.New(response.NoPermissionError)
	}

	var instance models.ReelCommentReaction
	instance.Emoji = obj.Emoji
	if obj.CanvasBranchID != 0 {
		instance.CanvasBranchID = &obj.CanvasBranchID
	}
	instance.CreatedByID = userID
	instance.UpdatedByID = userID
	instance.ReelID = obj.ReelID
	instance.ReelCommentID = obj.ReelCommentID
	results := r.db.Create(&instance)
	query := "reel_comment_id = " + strconv.FormatUint(instance.ReelCommentID, 10)
	counter, errGettingReactions := r.Manager.GetEmojiCounter(models.REELCOMMENTREACTION, query)
	fmt.Println(counter)
	fmt.Println(errGettingReactions)
	reactionjson, err := json.Marshal(counter)
	if err != nil {
		return err
	}
	err = App.Repo.Manager.UpdateEntityByID(models.REEL_COMMENTS, instance.ReelCommentID, map[string]interface{}{"reactions": reactionjson})
	if err != nil {
		return err
	}

	go func() {
		extraData := notifications.NotificationExtraData{
			CanvasRepoID:   reelInstance.CanvasRepositoryID,
			CanvasBranchID: reelInstance.CanvasBranchID,
		}
		contentObject := models.REELCOMMENTREACTION
		notifications.App.Service.PublishNewNotification(notifications.ReelCommentReact, userID, nil,
			&studioID, nil, extraData, &instance.ID, &contentObject)
	}()

	return results.Error
}
