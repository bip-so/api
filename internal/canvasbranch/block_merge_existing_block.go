package canvasbranch

import (
	"encoding/json"
	"fmt"
	"strconv"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
)

// BlockArtifactsMerger Moves Thread Comments + Reactions + Reels
func (s canvasBranchService) MoveThreadCommentsReelsReactionsMainBranchExistingBlock(
	sourceBlockID uint64,
	destinationBlockID uint64,
	destinationBranchID uint64) {

	App.Repo.TransferBlockThreadComments(sourceBlockID, destinationBlockID, destinationBranchID)
	App.Repo.TransferReels(sourceBlockID, destinationBlockID, destinationBranchID)
	App.Repo.TransferReactions(sourceBlockID, destinationBlockID, destinationBranchID)
}

// Move reactions friom the rough block to
func (r canvasBranchRepo) TransferReactions(fromBlockID uint64, toBlockID uint64, branchID uint64) {
	var blockReactions *[]models.BlockReaction
	fmt.Println("---------Debug-----")
	fmt.Println(fromBlockID)
	fmt.Println(toBlockID)
	fmt.Println(branchID)
	// List of ids of cloned block threads

	//clonedblockThreadIDs := App.Repo.GetListOfClonedFromThreadIDs(fromBlockID)

	// We need to ONLY delete those BLOCK THREADS where were clone and not all of them.

	// We are getting all the Block Reactions on a Parent Block and Then deleting them.
	// Cleaning parent (Condition)
	//errDelete := r.db.Table(models.BLOCKREACTION).Delete("block_id = ?", toBlockID)
	query := "DELETE from block_reactions WHERE block_id = ?"
	err := r.db.Exec(query, toBlockID).Error
	if err != nil {
		fmt.Println(err)
	}

	// DELETE FROM "block_reactions" WHERE "block_reactions"."id" = 13351

	//if errDelete.Error != nil {
	//	fmt.Println("Deleted FFS")
	//}
	_ = r.db.Model(models.BlockReaction{}).Where("block_id = ?", fromBlockID).Find(&blockReactions).Error
	for _, blockReaction := range *blockReactions {
		fmt.Println(blockReaction.ID)
		// Updates the Block Reactions
		results := r.db.Model(&models.BlockReaction{}).Where("id = ?", blockReaction.ID).Updates(map[string]interface{}{
			"block_id":              toBlockID,
			"canvas_branch_id":      branchID,
			"cloned_block_reaction": 0,
		})
		if results.Error != nil {
			fmt.Println("Error while updating BlockReaction")
			fmt.Println(results.Error.Error())
		}
	}

	// @todo: Update the Blocks -> reactions on Parent BLOCK NOW
	query2 := "block_id = " + strconv.FormatUint(toBlockID, 10)
	counter, errGettingReactions := r.Manager.GetEmojiCounter(models.BLOCKREACTION, query2)
	fmt.Println(errGettingReactions)
	// Block should get updated
	reactionjson, err := json.Marshal(counter)
	_ = App.Repo.Manager.UpdateEntityByID(models.BLOCK, toBlockID, map[string]interface{}{"reactions": reactionjson})
	// EOF

}

// Move BlockThreads  from instances to new block
func (r canvasBranchRepo) TransferBlockThreadComments(fromBlockID uint64, toBlockID uint64, branchID uint64) {
	var threads *[]models.BlockThread
	// Get all the threads
	/// Loop Through the thread
	_ = r.db.Model(models.BlockThread{}).Where("start_block_id = ?", fromBlockID).Find(&threads).Error
	//count := len(*threads)
	for _, thread := range *threads {
		r.db.Model(models.BlockThreadReaction{}).Where("block_thread_id = ?", thread.ID).Updates(
			map[string]interface{}{"canvas_branch_id": branchID, "block_id": toBlockID, "cloned_block_thread_reaction": 0})

		r.db.Model(models.BlockCommentReaction{}).Where("block_thread_id = ?", thread.ID).Updates(
			map[string]interface{}{"canvas_branch_id": branchID, "block_id": toBlockID, "cloned_block_comment_reaction": 0})

		// Move the pointer mainly block and branch
		var thisThreadCloned uint64
		thisThreadCloned = thread.ClonedFromThread
		err := r.db.Model(&models.BlockThread{}).Where("id = ?", thread.ID).Updates(map[string]interface{}{
			"start_block_id":     toBlockID,
			"canvas_branch_id":   branchID,
			"cloned_from_thread": 0,
		}).Error
		if err != nil {
			fmt.Println("Error while updating BlockThread from TransferBlockThreadComments")
			fmt.Println(err)
		}
		// Delete the Original Thread.
		if thisThreadCloned != 0 {
			// search for orginal block on the parent and delete
			_ = r.db.Model(models.BlockThread{}).Where("cloned_from_thread = ?", thisThreadCloned).Updates(map[string]interface{}{
				"cloned_from_thread": thread.ID,
			}).Error
			_ = r.Manager.HardDeleteByID(models.BLOCK_THREAD, thisThreadCloned)
		} else {
			go func() {
				extraData := notifications.NotificationExtraData{
					CanvasRepoID:   thread.CanvasRepositoryID,
					CanvasBranchID: branchID,
				}
				contentObject := models.BLOCK_THREAD
				notifications.App.Service.PublishNewNotification(notifications.BlockComment,
					thread.CreatedByID, nil, nil, nil, extraData, &thread.ID, &contentObject)
			}()
		}
	}
}

func (r canvasBranchRepo) TransferReels(fromBlockID uint64, toBlockID uint64, branchID uint64) {
	var reels *[]models.Reel
	_ = r.db.Model(models.Reel{}).Where("start_block_id = ?", fromBlockID).Find(&reels).Error
	for _, reel := range *reels {
		r.db.Model(models.ReelReaction{}).Where("reel_id = ?", reel.ID).Updates(map[string]interface{}{
			"canvas_branch_id": branchID,
		})
		// Move the pointer mainly block and branch
		err := r.db.Model(&models.Reel{}).Where("id = ?", reel.ID).Updates(map[string]interface{}{
			"start_block_id":   toBlockID,
			"canvas_branch_id": branchID,
		}).Error
		if err != nil {
			fmt.Println("Error while updating BlockThread from TransferBlockThreadComments")
			fmt.Println(err)
		}
		// Delete the Original Thread.
		if reel.ClonedFromReel != 0 {
			// process a cloned reel.
			// search for cloned block if exists, delete
			_ = r.db.Model(models.Reel{}).Where("cloned_from_reel = ?", reel.ClonedFromReel).Updates(map[string]interface{}{
				"cloned_from_reel": reel.ID,
			}).Error
			r.db.Model(models.ReelReaction{}).Where("reel_id = ?", reel.ClonedFromReel).Updates(map[string]interface{}{
				"reel_id": reel.ID,
			})
			_ = r.Manager.HardDeleteByID(models.REEL, reel.ClonedFromReel)
		}
	}
}

//TransferReactions
