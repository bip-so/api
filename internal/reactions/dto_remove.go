package reactions

import (
	"encoding/json"
	"fmt"
	"strconv"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

func (r reactionRepo) RemoveBlockReaction(obj CreateMentionPost, studioID uint64, userID uint64) error {
	blockInstance, err := App.Repo.GetBlockByUUIDAndBranchID(map[string]interface{}{"uuid": obj.BlockUUID, "canvas_branch_id": obj.CanvasBranchID})

	var instance models.BlockReaction
	instance.Emoji = obj.Emoji
	if obj.CanvasBranchID != 0 {
		instance.CanvasBranchID = &obj.CanvasBranchID
	}
	instance.CreatedByID = userID
	instance.BlockID = blockInstance.ID
	results := r.db.Where(&instance).Delete(&instance)

	// Create now updating counters
	query := "block_id = " + strconv.FormatUint(blockInstance.ID, 10)
	counter, _ := r.Manager.GetEmojiCounter(models.BLOCKREACTION, query)
	reactionjson, err := json.Marshal(counter)
	if err != nil {
		return err
	}
	err = App.Repo.Manager.UpdateEntityByID(models.BLOCK, blockInstance.ID, map[string]interface{}{"reactions": reactionjson})
	if err != nil {
		return err
	}
	return results.Error
}
func (r reactionRepo) RemoveBlockThreadReaction(obj CreateMentionPost, studioID uint64, userID uint64) error {
	blockInstance, err := App.Repo.GetBlockByUUIDAndBranchID(map[string]interface{}{"uuid": obj.BlockUUID, "canvas_branch_id": obj.CanvasBranchID})

	var instance models.BlockThreadReaction
	instance.Emoji = obj.Emoji
	if obj.CanvasBranchID != 0 {
		instance.CanvasBranchID = &obj.CanvasBranchID
	}
	instance.CreatedByID = userID
	instance.BlockID = blockInstance.ID
	instance.BlockThreadID = obj.BlockThreadID
	results := r.db.Where(&instance).Delete(&instance)

	query := "block_thread_id = " + strconv.FormatUint(obj.BlockThreadID, 10)
	counter, _ := r.Manager.GetEmojiCounter(models.BLOCKTHREADREACTION, query)
	reactionjson, err := json.Marshal(counter)
	if err != nil {
		return err
	}
	err = App.Repo.Manager.UpdateEntityByID(models.BLOCK_THREAD, obj.BlockThreadID, map[string]interface{}{"reactions": reactionjson})
	if err != nil {
		return err
	}
	return results.Error
}
func (r reactionRepo) RemoveBlockThreadCommentReaction(obj CreateMentionPost, studioID uint64, userID uint64) error {
	blockInstance, err := App.Repo.GetBlockByUUIDAndBranchID(map[string]interface{}{"uuid": obj.BlockUUID, "canvas_branch_id": obj.CanvasBranchID})

	var instance models.BlockCommentReaction
	instance.Emoji = obj.Emoji
	if obj.CanvasBranchID != 0 {
		instance.CanvasBranchID = &obj.CanvasBranchID
	}
	instance.CreatedByID = userID

	instance.BlockID = blockInstance.ID
	instance.BlockThreadID = obj.BlockThreadID
	instance.BlockCommentID = obj.BlockCommentID
	results := r.db.Where(&instance).Delete(&instance)

	query := "block_comment_id = " + strconv.FormatUint(obj.BlockCommentID, 10)
	counter, _ := r.Manager.GetEmojiCounter(models.BLOCKCOMMENTREACTION, query)
	reactionjson, err := json.Marshal(counter)
	if err != nil {
		return err
	}
	err = App.Repo.Manager.UpdateEntityByID(models.BLOCK_THREAD_COMMENT, obj.BlockCommentID, map[string]interface{}{"reactions": reactionjson})
	if err != nil {
		return err
	}
	return results.Error
}
func (r reactionRepo) RemoveReelReaction(obj CreateMentionPost, studioID uint64, userID uint64) error {
	var instance models.ReelReaction
	instance.Emoji = obj.Emoji
	if obj.CanvasBranchID != 0 {
		instance.CanvasBranchID = &obj.CanvasBranchID
	}
	instance.CreatedByID = userID
	instance.ReelID = obj.ReelID
	results := r.db.Where(&instance).Delete(&instance)
	query := "reel_id = " + strconv.FormatUint(obj.ReelID, 10)
	counter, errGettingReactions := r.Manager.GetEmojiCounter(models.REELREACTION, query)
	fmt.Println(counter)
	fmt.Println(errGettingReactions)
	reactionjson, err := json.Marshal(counter)
	if err != nil {
		return err
	}
	err = App.Repo.Manager.UpdateEntityByID(models.REEL, obj.ReelID, map[string]interface{}{"reactions": reactionjson})
	if err != nil {
		return err
	}
	return results.Error
}
func (r reactionRepo) RemoveReelCommentReaction(obj CreateMentionPost, studioID uint64, userID uint64) error {
	var instance models.ReelCommentReaction
	instance.Emoji = obj.Emoji
	if obj.CanvasBranchID != 0 {
		instance.CanvasBranchID = &obj.CanvasBranchID
	}
	instance.CreatedByID = userID
	instance.ReelID = obj.ReelID
	instance.ReelCommentID = obj.ReelCommentID
	results := r.db.Where(&instance).Delete(&instance)

	query := "reel_comment_id = " + strconv.FormatUint(obj.ReelCommentID, 10)
	counter, errGettingReactions := r.Manager.GetEmojiCounter(models.REELCOMMENTREACTION, query)
	fmt.Println(counter)
	fmt.Println(errGettingReactions)
	reactionjson, err := json.Marshal(counter)
	if err != nil {
		return err
	}
	err = App.Repo.Manager.UpdateEntityByID(models.REEL_COMMENTS, obj.ReelCommentID, map[string]interface{}{"reactions": reactionjson})
	if err != nil {
		return err
	}
	return results.Error
}
