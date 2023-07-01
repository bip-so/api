package global

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"strings"
	"time"
)

func BlocksCloner(fromBranchId uint64, toBranchId uint64, user *models.User, newCanvasRepoID uint64) error {
	// Get the blocks
	postgres.GetDB().Where("canvas_branch_id = ?", toBranchId).Delete(models.Block{})
	ogblocks, err := queries.App.BlockQuery.GetBlocksByBranchID(fromBranchId)
	if err != nil {
		return err
	}
	// If Size of og blocks is 0 we return.
	if len(*ogblocks) == 0 {
		return nil
	}
	var newBlocks []models.Block
	for _, block := range *ogblocks {
		newBlock := CloneBlockInstance(block, toBranchId, user, newCanvasRepoID)
		newBlocks = append(newBlocks, newBlock)
	}
	errBulkCreating := postgres.GetDB().Create(&newBlocks).Error
	if errBulkCreating != nil {
		return errBulkCreating
	}
	// All the Blocks are created now, We need to move things
	_, errGettingClonedBlock := queries.App.BlockQuery.GetBlocksByBranchID(toBranchId)
	if errGettingClonedBlock != nil {
		return errGettingClonedBlock
	}
	// Get newly created blocks and create things on it.
	// commenting as comments no need to be duplicated
	// for _, newBlock := range *newblocks {
	// 	 CloneBlockThreadAndComments(&newBlock, toBranchId, newCanvasRepoID)
	// }
	return nil
}

func CloneBlockInstance(ogBlockInstance models.Block, branchID uint64, user *models.User, newCanvasRepoID uint64) models.Block {
	var block models.Block
	block = ogBlockInstance
	block.ClonedFromBlockID = ogBlockInstance.ID
	block.ID = 0 // Reset the PK
	block.UUID = uuid.New()
	block.CanvasBranchID = &branchID
	block.CanvasRepositoryID = newCanvasRepoID
	fmt.Println(user.ID, user.FullName)
	block.CreatedByUser = nil
	block.UpdatedByUser = nil
	block.CreatedByID = 0
	block.UpdatedByID = 0
	block.ArchivedByID = 0
	block.CreatedByID = user.ID
	block.UpdatedByID = user.ID
	block.CreatedAt = time.Now()
	block.UpdatedAt = time.Now()
	contributors := []map[string]interface{}{
		{
			"id":        user.ID,
			"uuid":      user.UUID.String(),
			"repoID":    newCanvasRepoID,
			"branchID":  branchID,
			"fullName":  user.FullName,
			"username":  user.Username,
			"avatarUrl": user.AvatarUrl,
			"timestamp": time.Now(),
		},
	}
	contributorsStr, _ := json.Marshal(contributors)
	block.Contributors = contributorsStr
	block.CommentCount = 0
	return block
}

func CloneBlockThreadAndComments(newBlock *models.Block, branchID uint64, newCanvasRepoID uint64) {
	var blockThreads *[]models.BlockThread
	_ = postgres.GetDB().Model(&models.BlockThread{}).Where("start_block_id = ? and is_archived = false", newBlock.ClonedFromBlockID).Find(&blockThreads).Error
	// loop all the block threads and create a copy.
	for _, blockThread := range *blockThreads {
		CreateClonedBlockThread(&blockThread, newBlock, branchID, newCanvasRepoID)
	}
}

func CreateClonedBlockThread(og *models.BlockThread, newBlock *models.Block, branchID uint64, newCanvasRepoID uint64) {
	// Copy OG to Instance with new block id and branch
	fmt.Println("Cloning Block Thread")
	fmt.Println(og.ID)
	blockThreadClonedID := og.ID
	blockThreadUUID := og.UUID.String()

	instance := og
	instance.StartBlockID = newBlock.ID
	instance.StartBlockUUID = newBlock.UUID
	instance.CanvasBranchID = branchID
	instance.CanvasRepositoryID = newCanvasRepoID
	// Saving the reference to the original BlockThread.
	instance.ClonedFromThread = og.ID
	instance.UUID = uuid.New()
	instance.ID = 0 // Presumably we are setting ID to NIL
	_ = postgres.GetDB().Create(&instance)
	// Now update the block with the comment uuid
	var blockChildren []map[string]interface{}
	json.Unmarshal(newBlock.Children, &blockChildren)
	for i, children := range blockChildren {
		containsComment := false
		for key, _ := range children {
			if strings.Contains(key, "commentThread_") {
				containsComment = true
				break
			}
		}
		if containsComment {
			delete(blockChildren[i], fmt.Sprintf("commentThread_%s", blockThreadUUID))
			blockChildren[i][fmt.Sprintf("commentThread_%s", instance.UUID.String())] = true
		}
	}
	blockChildrenStr, _ := json.Marshal(blockChildren)
	newBlock.Children = blockChildrenStr
	postgres.GetDB().Save(&newBlock)
	// Now we have a created a new BlockThread we also need to create Comments copy
	var comments *[]models.BlockComment
	// Get all comments from the parent thread id
	_ = postgres.GetDB().Model(&models.BlockComment{}).Where("thread_id = ? and is_archived = false", blockThreadClonedID).Find(&comments).Error
	fmt.Println(comments)
	for _, blockThreadComment := range *comments {
		queries.App.BlockQuery.CreateBlockThreadComment(&blockThreadComment, instance.ID, branchID)
	}
}

func UpdateMentionsInBlocks(branchID uint64, canvasBranchMap map[uint64]uint64, user *models.User, studio *models.Studio) {
	newblocks, err := GetBlocksByBranchID(branchID)
	if err != nil {
		fmt.Println("Error in getting blocks of branch", err)
		return
	}
	for blockIndex, block := range newblocks {
		var blockChildren []map[string]interface{}
		json.Unmarshal(block.Children, &blockChildren)
		pageMentionPresent := false
		for childIndex, children := range blockChildren {
			if children["type"] != nil && children["type"].(string) == "pageMention" {
				oldMention := children["mention"].(map[string]interface{})
				newBranchID := canvasBranchMap[uint64(oldMention["id"].(float64))]
				newBranch, _ := queries.App.BranchQuery.GetBranchByID(newBranchID)
				newMention := map[string]interface{}{
					"id":                    newBranch.ID,
					"key":                   newBranch.Key,
					"name":                  newBranch.Name,
					"type":                  "branch",
					"uuid":                  newBranch.UUID.String(),
					"repoID":                newBranch.CanvasRepositoryID,
					"repoKey":               newBranch.CanvasRepository.Key,
					"repoName":              strings.TrimSuffix(newBranch.CanvasRepository.Name, "\n"),
					"repoUUID":              newBranch.CanvasRepository.UUID.String(),
					"studioID":              newBranch.CanvasRepository.StudioID,
					"createdByUserID":       user.ID,
					"createdByUserFullName": user.FullName,
					"createdByUserUsername": user.Username,
				}
				blockChildren[childIndex]["mention"] = newMention
				blockChildren[childIndex]["uuid"] = uuid.New().String()
				pageMentionPresent = true
			} else if children["type"] != nil && children["type"].(string) == "userMention" {
				newMention := map[string]interface{}{
					"id":                    user.ID,
					"type":                  "user",
					"uuid":                  user.UUID.String(),
					"fullName":              user.FullName,
					"studioID":              studio.ID,
					"username":              user.Username,
					"avatarUrl":             user.AvatarUrl,
					"createdByUserID":       user.ID,
					"createdByUserFullName": user.FullName,
					"createdByUserUsername": user.FullName,
				}
				blockChildren[childIndex]["mention"] = newMention
				blockChildren[childIndex]["uuid"] = uuid.New().String()
				blockChildren[childIndex]["children"] = []map[string]string{
					{
						"text": fmt.Sprintf("<@%s>", user.FullName),
					},
				}
			}
		}
		if pageMentionPresent {
			blockChildrenStr, _ := json.Marshal(blockChildren)
			newblocks[blockIndex].Children = blockChildrenStr
		}
	}
	postgres.GetDB().Save(newblocks)
}

func GetBlocksByBranchID(branchID uint64) ([]models.Block, error) {
	var blocks []models.Block
	err := postgres.GetDB().Model(&models.Block{}).Where("canvas_branch_id = ?", branchID).Preload("CreatedByUser").Preload("UpdatedByUser").Order("rank ASC").Find(&blocks).Error
	if err != nil {
		return nil, err
	}
	return blocks, nil
}
