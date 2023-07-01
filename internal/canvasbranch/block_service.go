package canvasbranch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/blocks"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"

	"github.com/google/uuid"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/reactions"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

// Clones Blocks from one Branch to Another Branch
func (s canvasBranchService) CloneBlocks(fromBranchId, tofromBranchId, userID uint64) error {
	err := queries.App.BlockQuery.BlocksCloner(fromBranchId, tofromBranchId, userID)
	if err != nil {
		return err
	}
	return nil
}

func uniqueUUIDList(s []uuid.UUID) []uuid.UUID {
	inResult := make(map[uuid.UUID]bool)
	var result []uuid.UUID
	for _, item := range s {
		if _, ok := inResult[item]; !ok {
			inResult[item] = true
			result = append(result, item)
		}
	}
	return result
}

func hasDuplicateUUID(uuidSlice []uuid.UUID) bool {
	// @todo: Add check for duplicate values in UUID
	cleanedList := uniqueUUIDList(uuidSlice)
	if len(uuidSlice) == len(cleanedList) {
		return false
	} else {
		return true
	}

}

// ValidateBlocksData will check if post has proper Positions, BlockTyp, Scope
func (s canvasBranchService) ValidateBlocksData(branchID uint64, data CanvasBlockPost, isAssociationAPI bool) error {
	// 1. Position
	// 2. Block Type
	// 3. Scope
	// 5. Check for Duplicate UUID's
	// Looping through each block item and checking.
	// We can also add more in future
	// If the Branch is "Committed" is True
	branch, errGettingBranch := queries.App.BranchQuery.GetBranchByID(branchID)
	if errGettingBranch != nil {
		return errors.New("Branch not found!!!")
	}
	if !isAssociationAPI {
		if branch.Committed {
			return errors.New("Please cancel merge request before editing the document.")
		}
	}

	var listUUIDS = []uuid.UUID{}

	for _, blockItem := range data.Blocks {
		if blockItem.UUID.String() == "" {
			return errors.New("UUID Field is empty.")
		}

		listUUIDS = append(listUUIDS, blockItem.UUID)
		// rank := blockItem.Rank
		blocktype := blockItem.Type
		scope := blockItem.Scope

		if !utils.SliceContainsItem(models.AllowedBlockTypes, blocktype) {
			return errors.New("Incorrect block type found")
		}
		// Scope
		if !utils.SliceContainsItem(AllowedScope, scope) {
			return errors.New("Incorrect scope found")
		}

	}

	if hasDuplicateUUID(listUUIDS) {
		return errors.New("Found duplicate UUID's ")
	}

	return nil
}

func (s canvasBranchService) EmptyBlockInstance() *models.Block {
	return &models.Block{}
}

func (s canvasBranchService) ProcessCreateBlock(user models.User, branchId uint64, data models.PostBlocks) (uint64, error) {
	firstContrib := blocks.App.Service.BlockContributorFirst(user, branchId)
	id, err := queries.App.BlockQuery.CreateBlock2(user, branchId, data, firstContrib)
	if err != nil {
		return 0, err
	}
	var attributes map[string]interface{}
	err = json.Unmarshal(data.Attributes, &attributes)
	if err == nil {
		if messageId, exists := attributes["messageId"]; exists {
			go App.Git.CommitMessageBlockToGit(id, messageId.(string))
		}
	}
	return id, nil
}

func (s canvasBranchService) ProcessUpdateBlock(user models.User, branchId uint64, data models.PostBlocks) (uint64, error) {
	contrib := blocks.App.Service.BlockContributorNext(user, branchId)
	id, err := queries.App.BlockQuery.UpdateBlock(user.ID, branchId, data, contrib)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s canvasBranchService) ProcessUpdateClonedBlock(userID uint64, branchId uint64, data models.PostBlocks) (uint64, error) {
	id, err := queries.App.BlockQuery.UpdateClonedBlock(userID, branchId, data)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s canvasBranchService) ProcessDeleteBlock(userID uint64, branchId uint64, data models.PostBlocks) error {
	if data.ID != 0 {
		err := queries.App.BlockQuery.DeleteBlock(data.ID, userID)
		if err != nil {
			return err
		}
	} else {
		var query = map[string]interface{}{"uuid": data.UUID, "canvas_branch_id": branchId}
		blockInstance, _ := queries.App.BlockQuery.GetBlock(query)
		_ = queries.App.BlockQuery.DeleteBlock(blockInstance.ID, userID)
	}
	return nil
}

func (s canvasBranchService) ProcessDeleteClonedBlock(userID uint64, branchId uint64, data models.PostBlocks) error {
	err := queries.App.BlockQuery.DeleteClonedBlock(data.UUID, userID, branchId)
	if err != nil {
		return err
	}
	return nil
}

func (s canvasBranchService) GetAllBlockByBranchID(branchId uint64) (*[]models.Block, error) {
	blocks, err := queries.App.BlockQuery.GetBlocksByBranchID(branchId)
	if err != nil {
		return nil, err
	}
	return blocks, nil
}

func (s canvasBranchService) addReactionsToBlocks(blocks *[]models.Block, branchID uint64, user *models.User) ([]BulkBlocks, error) {
	var blockReactions *[]models.BlockReaction
	var err error
	if user != nil {
		blockReactions, err = reactions.App.Repo.GetBlockReactions(map[string]interface{}{"canvas_branch_id": branchID, "created_by_id": user.ID})
		if err != nil {
			return nil, err
		}
	} else {
		blockReactions, err = reactions.App.Repo.GetBlockReactions(map[string]interface{}{"canvas_branch_id": branchID})
		if err != nil {
			return nil, err
		}
	}
	blocksData := BulkBlocksReactionsSerializerData(blocks, blockReactions, user)
	return blocksData, nil
}

func (s canvasBranchService) GetBlocksFromHistoryByCommitId(user *models.User, commitID string, canvasBranchID uint64) ([]*models.Block, error) {

	branch, err := queries.App.BranchQuery.GetBranchByID(canvasBranchID)
	if err != nil {
		return nil, err
	}

	return App.Git.FetchCommitFromGit(user, commitID, branch.CanvasRepositoryID)
}

func (s canvasBranchService) GetBranchBlocksCachedData(canvasBranchID, userID uint64) ([]BulkBlocks, error) {
	var cachedBranchKey string
	//canvasBranch, _ := App.Repo.Get(map[string]interface{}{"id": canvasBranchID})
	canvasBranch, _ := queries.App.BranchQuery.GetBranchWithRepoAndStudio(canvasBranchID)

	if userID == 0 || canvasBranch.Name == "main" {
		cachedBranchKey = redis.GenerateCacheKey(models.CANVAS_BRANCH, []string{utils.String(canvasBranchID)})
	} else {
		cachedBranchKey = redis.GenerateCacheKey(models.CANVAS_BRANCH, []string{utils.String(canvasBranchID), utils.String(userID)})
	}
	data := s.cache.Get(context.Background(), cachedBranchKey)
	var blocks []BulkBlocks
	if data != nil {
		err := json.Unmarshal([]byte(data.(string)), &blocks)
		if err != nil {
			fmt.Println("GetBlocks Cache data", err)
			return nil, err
		}
	}
	return blocks, nil
}

func (s canvasBranchService) CacheBranchBlocks(canvasBranchID, userID uint64, blocks []BulkBlocks) {
	var cachedBranchKey string
	//canvasBranch, _ := App.Repo.Get(map[string]interface{}{"id": canvasBranchID})
	canvasBranch, _ := queries.App.BranchQuery.GetBranchWithRepoAndStudio(canvasBranchID)

	if userID == 0 || canvasBranch.Name == "main" {
		cachedBranchKey = redis.GenerateCacheKey(models.CANVAS_BRANCH, []string{utils.String(canvasBranchID)})
	} else {
		cachedBranchKey = redis.GenerateCacheKey(models.CANVAS_BRANCH, []string{utils.String(canvasBranchID), utils.String(userID)})
	}
	blocksStr, _ := json.Marshal(blocks)
	s.cache.Set(context.Background(), cachedBranchKey, blocksStr, nil)
}

func (s canvasBranchService) InvalidateBranchBlocks(canvasBranchID uint64) {
	cachedBranchKey := models.CANVAS_BRANCH + ":" + utils.String(canvasBranchID) + ":*"
	s.cache.DeleteMatching(context.Background(), cachedBranchKey)
	cachedBranchKey = models.CANVAS_BRANCH + ":" + utils.String(canvasBranchID)
	s.cache.Delete(context.Background(), cachedBranchKey)
}
