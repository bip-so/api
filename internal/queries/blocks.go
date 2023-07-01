package queries

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/s3"
	"gorm.io/datatypes"
	"strings"
	"time"
)

type UserBlockContributor struct {
	Id        uint64    `json:"id"`
	UUID      string    `json:"uuid"`
	FullName  string    `json:"fullName"`
	Username  string    `json:"username"`
	AvatarUrl string    `json:"avatarUrl"`
	Timestamp time.Time `json:"timestamp"`
	BranchID  uint64    `json:"branchID"`
	RepoID    uint64    `json:"repoID"`
}

func (q blockQuery) BlockContributorFirst(user models.User, branchId uint64) datatypes.JSON {
	contribution := UserBlockContributor{
		Id:        user.ID,
		UUID:      user.UUID.String(),
		FullName:  user.FullName,
		Username:  user.Username,
		AvatarUrl: user.AvatarUrl,
		Timestamp: time.Now(),
		BranchID:  branchId,
	}
	j := []UserBlockContributor{contribution}
	singleContrib, _ := json.Marshal(j)
	first := datatypes.JSON(singleContrib)
	return first
}

//CreateBlock : Create a new block
func (q blockQuery) CreateBlock(instance *models.Block) (*models.Block, error) {
	results := postgres.GetDB().Create(&instance)
	return instance, results.Error
}

func (q blockQuery) CreateFirstBlock(user models.User, branchInstance *models.CanvasBranch, repoInstance *models.CanvasRepository) error {
	var block models.Block
	firstContrib := q.BlockContributorFirst(user, branchInstance.ID)
	newBlock, errCreatingBlockInst := block.NewBlock(
		uuid.New(),
		user.ID,
		repoInstance.ID,
		branchInstance.ID,
		nil,
		1,
		2,
		models.BlockTypeText,
		models.MyFirstBlockJson(),
		models.MyFirstEmptyBlockJson(),
		firstContrib,
	)
	if errCreatingBlockInst != nil {
		return errCreatingBlockInst
	}
	_, errNewBlock := q.CreateBlock(newBlock)
	if errNewBlock != nil {
		return errNewBlock
	}
	return nil
}

func (q blockQuery) UpdateDiscordImageURL(block *models.Block) error {
	if block.Type != models.BlockTypeText {
		return nil
	}
	var childrenData []map[string]interface{}
	err := json.Unmarshal(block.Children, &childrenData)
	if err != nil {
		return err
	}
	if len(childrenData) == 3 && childrenData[1]["type"] == models.BlockTypeImage {
		url := childrenData[1]["url"].(string)
		if url != "" && strings.HasPrefix(url, "https://cdn.discordapp.com") {
			imgResponse, err := s3.UploadImageFromURLToS3(url, "blocks/"+block.UUID.String()+".jpg", true, false)
			if err != nil {
				return err
			}
			childrenData[1]["url"] = imgResponse.URL
			data, _ := json.Marshal(childrenData)
			block.Children = data
		}
	}
	return nil
}

func (q blockQuery) DoesThisUserHaveContribOnThisBlock(userID uint64, contrib datatypes.JSON) bool {
	var contribs []UserBlockContributor
	_ = json.Unmarshal(contrib, &contribs)
	for _, v := range contribs {
		fmt.Println(userID)
		fmt.Println(v.Id)
		if v.Id == userID {
			return true
		}
	}
	return false
}

func (q blockQuery) GetBlock(query map[string]interface{}) (*models.Block, error) {
	var block models.Block
	err := postgres.GetDB().Model(&models.Block{}).Where(query).First(&block).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &block, nil
}

func (q blockQuery) UpdateBlock(userID uint64, branchId uint64, data models.PostBlocks, contrib datatypes.JSON) (uint64, error) {
	var query map[string]interface{}
	if data.ID != 0 {
		query = map[string]interface{}{"id": data.ID}
	} else {
		query = map[string]interface{}{"uuid": data.UUID, "canvas_branch_id": branchId}
	}
	blockInstance, err := q.GetBlock(query)
	if err != nil {
		return 0, err
	}
	blockInstance.Version = data.Version
	blockInstance.Type = data.Type
	blockInstance.Rank = data.Rank //(0)
	blockInstance.Children = data.Children

	blockInstance.UpdatedByID = userID
	blockInstance.Attributes = data.Attributes
	updates := map[string]interface{}{
		"version":       data.Version,
		"type":          data.Type,
		"rank":          data.Rank,
		"children":      data.Children,
		"updated_by_id": userID,
		"attributes":    data.Attributes,
	}
	haveContrib := q.DoesThisUserHaveContribOnThisBlock(userID, blockInstance.Contributors)
	if !haveContrib {
		q.Manager.UserPushBlockContributors(blockInstance.ID, contrib.String())
	}

	errUpdating := postgres.GetDB().Model(&models.Block{}).Where("id = ?", blockInstance.ID).Updates(&updates)
	if errUpdating != nil {
		return 0, err
	}

	return blockInstance.ID, nil
}

// Update a Block Instance
func (r blockQuery) UpdateClonedBlock(userID uint64, branchId uint64, data models.PostBlocks) (uint64, error) {
	// Getting block from UUID and Brach because we don't have the right id.
	blockUUID := data.UUID
	blockInstance, err := r.GetBlock(map[string]interface{}{"uuid": blockUUID, "canvas_branch_id": branchId})
	fmt.Println(blockInstance)
	if err != nil {
		return 0, err
	}
	blockInstance.Version = data.Version
	blockInstance.Type = data.Type
	blockInstance.Rank = data.Rank
	blockInstance.Children = data.Children
	blockInstance.UpdatedByID = userID

	errUpdating := postgres.GetDB().Model(&models.Block{}).Where("id = ?", blockInstance.ID).Updates(&blockInstance)
	if errUpdating != nil {
		return 0, err
	}

	return blockInstance.ID, nil
}

func (r blockQuery) CreateBlock2(user models.User, branchId uint64, data models.PostBlocks, firstContrib datatypes.JSON) (uint64, error) {
	//branchInstance, err := App.Repo.Get(map[string]interface{}{"id": branchId})
	branchInstance, err := App.BranchQuery.GetBranchWithRepoAndStudio(branchId)

	if err != nil {
		return 0, err
	}

	var block models.Block

	instance, _ := block.NewBlock(
		data.UUID,
		user.ID,
		branchInstance.CanvasRepositoryID,
		branchInstance.ID,
		nil,
		data.Rank,
		2, //data.Version,
		data.Type,
		data.Children,
		data.Attributes,
		firstContrib,
	)
	r.UpdateDiscordImageURL(instance)
	results := postgres.GetDB().Create(&instance)
	return instance.ID, results.Error
}

// This function will return tru is the cotrib exists for this user on a block

func (r blockQuery) DeleteBlock(blockID uint64, userId uint64) error {
	err := r.Manager.HardDeleteByID(models.BLOCK, blockID)
	if err != nil {
		return err
	}
	return nil
}

func (r blockQuery) DeleteClonedBlock(uuid uuid.UUID, userId uint64, branchID uint64) error {
	blockInstance, errGet := r.GetBlock(map[string]interface{}{"uuid": uuid, "canvas_branch_id": branchID})
	if errGet != nil {
		return errGet
	}
	err := r.Manager.HardDeleteByID(models.BLOCK, blockInstance.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r blockQuery) GetBlocksByBranchID(branchID uint64) (*[]models.Block, error) {
	var blocks []models.Block
	var blocksData []map[string]interface{}
	//err := r.db.Model(&models.Block{}).Where("canvas_branch_id = ?", branchID).Preload("CreatedByUser").Preload("UpdatedByUser").Order("rank ASC").Find(&blocks).Error
	_ = postgres.GetDB().Raw("SELECT bl.*, cb.id as cu_user_id, cb.uuid as cu_user_uuid, cb.email as cu_user_email, cb.username as cu_user_username, cb.avatar_url as cu_user_avatar_url, cb.is_setup_done as cu_is_setup_done, cb.full_name as cu_full_name, ub.id as ub_user_id, ub.uuid as ub_user_uuid, ub.email as ub_user_email, ub.username as ub_user_username, ub.avatar_url as ub_user_avatar_url, ub.is_setup_done as ub_is_setup_done, ub.full_name as ub_full_name   FROM blocks bl Join users cb on cb.id = bl.created_by_id Join users ub on ub.id = bl.updated_by_id Where bl.canvas_branch_id = ? Order by bl.rank", branchID).Scan(&blocksData)

	//.Where("canvas_branch_id = ?", branchID).Preload("CreatedByUser").Preload("UpdatedByUser").Order("rank ASC").Find(&blocks).Error
	//if err != nil {
	//	return nil, err
	//}
	for _, block := range blocksData {
		var newBlock models.Block
		var createdByUser, updatedByUser models.User
		newBlock.CreatedByUser = &createdByUser
		newBlock.UpdatedByUser = &updatedByUser

		newBlock.ID = uint64(block["id"].(int64))
		newBlock.ArchivedByID = uint64(block["archived_by_id"].(int64))
		newBlock.UpdatedByID = uint64(block["updated_by_id"].(int64))
		newBlock.ClonedFromBlockID = uint64(block["cloned_from_block_id"].(int64))
		newBlock.CreatedByID = uint64(block["created_by_id"].(int64))
		newBlock.CommentCount = uint(block["comment_count"].(int64))
		newBlock.ReelCount = uint(block["reel_count"].(int64))
		newBlock.Version = uint(block["version"].(int64)) //block["version"].(uint)
		newBlockCanvasBranchID := block["canvas_branch_id"].(int64)
		newBlockCanvasBranchIDUInt := uint64(newBlockCanvasBranchID)
		newBlock.CanvasBranchID = &newBlockCanvasBranchIDUInt
		newBlock.CanvasRepositoryID = uint64(block["canvas_repository_id"].(int64))

		if block["parent_block"] == nil {
			newBlock.ParentBlock = nil
		} else {
			newBlockParentBlock := block["parent_block"].(int64)
			newBlockParentBlockUInt := uint64(newBlockParentBlock)
			newBlock.ParentBlock = &newBlockParentBlockUInt
		}

		if block["attributes"] == nil {
			newBlock.Attributes = datatypes.JSON(`{}`)
		} else {
			newBlockAttributes := datatypes.JSON(block["attributes"].(string))
			newBlock.Attributes = newBlockAttributes
		}

		if block["children"] == nil {
			newBlock.Children = datatypes.JSON(`{}`)
		} else {
			newBlockChildren := datatypes.JSON(block["children"].(string))
			newBlock.Children = newBlockChildren
		}

		if block["contributors"] == nil {
			newBlock.Contributors = datatypes.JSON(`{}`)
		} else {
			newBlockContributors := datatypes.JSON(block["contributors"].(string))
			newBlock.Contributors = newBlockContributors
		}

		if block["mentions"] == nil {
			x := datatypes.JSON(`{}`)
			newBlock.Mentions = &x
		} else {
			newBlockMentions := datatypes.JSON(block["mentions"].(string))
			newBlock.Mentions = &newBlockMentions
		}

		if block["reactions"] == nil {
			newBlock.Reactions = datatypes.JSON(`{}`)
		} else {
			newBlockReactions := datatypes.JSON(block["reactions"].(string))
			newBlock.Reactions = newBlockReactions
		}

		//fmt.Println(block["archived_at"])
		//xType := reflect.TypeOf(block["archived_at"])
		//fmt.Println(xType) // "[]int [1 2 3]"

		newBlock.ArchivedAt = block["archived_at"].(time.Time)
		newBlock.CreatedAt = block["created_at"].(time.Time)
		newBlock.UpdatedAt = block["updated_at"].(time.Time)
		newBlock.IsArchived = block["is_archived"].(bool)
		newBlock.Rank = block["rank"].(int32)
		newBlock.Type = block["type"].(string)

		//newBlockUUID := block["uuid"]
		newBlock.UUID, _ = uuid.Parse(block["uuid"].(string))

		//fmt.Println("block[\"ub_full_name\"]   --> ", block["ub_full_name"])
		if block["ub_full_name"] == nil {
			newBlock.UpdatedByUser.FullName = ""
		} else {
			newBlock.UpdatedByUser.FullName = block["ub_full_name"].(string)
		}

		newBlock.UpdatedByUser.IsSetupDone = block["ub_is_setup_done"].(bool)
		newBlock.UpdatedByUser.AvatarUrl = block["ub_user_avatar_url"].(string)
		if block["ub_user_email"] == nil {
			newBlock.UpdatedByUser.Email = sql.NullString{
				Valid:  false,
				String: "",
			}
		} else {
			newBlock.UpdatedByUser.Email = sql.NullString{
				Valid:  true,
				String: block["ub_user_email"].(string),
			}
		}

		newBlock.UpdatedByUser.ID = uint64(block["ub_user_id"].(int64))
		newBlock.UpdatedByUser.Username = block["ub_user_username"].(string)

		newBlock.UpdatedByUser.UUID, _ = uuid.Parse(block["ub_user_uuid"].(string))

		if block["ub_full_name"] == nil {
			newBlock.CreatedByUser.FullName = ""
		} else {
			newBlock.CreatedByUser.FullName = block["cu_full_name"].(string)
		}

		newBlock.CreatedByUser.IsSetupDone = block["cu_is_setup_done"].(bool)
		newBlock.CreatedByUser.AvatarUrl = block["cu_user_avatar_url"].(string)
		//newBlock.CreatedByUser.Email = block["cu_user_email"].(sql.NullString)
		if block["cu_user_email"] == nil {
			newBlock.UpdatedByUser.Email = sql.NullString{
				Valid:  false,
				String: "",
			}
		} else {
			newBlock.UpdatedByUser.Email = sql.NullString{
				Valid:  true,
				String: block["cu_user_email"].(string),
			}
		}
		newBlock.CreatedByUser.ID = uint64(block["cu_user_id"].(int64)) //block["cu_user_id"].(uint64)
		newBlock.CreatedByUser.Username = block["cu_user_username"].(string)
		newBlock.CreatedByUser.UUID, _ = uuid.Parse(block["cu_user_uuid"].(string))
		//newBlock.CreatedByUser.UUID = block["cu_user_uuid"].(uuid.UUID)

		//fmt.Printf("%+v\n", newBlock)

		blocks = append(blocks, newBlock)
	}

	return &blocks, nil
}

func (r blockQuery) GetSourceBlocksByBranchID(branchID uint64) ([]*models.Block, error) {
	var blocks []*models.Block
	err := postgres.GetDB().Model(&models.Block{}).Where("canvas_branch_id = ?", branchID).Preload("CreatedByUser").Preload("UpdatedByUser").Order("rank ASC").Find(&blocks).Error
	if err != nil {
		return nil, err
	}
	return blocks, nil
}

// Cloning All the Blocks from Main Branch to Rough Branch
func (r blockQuery) BlocksCloner(fromBranchId uint64, toBranchId uint64, userID uint64) error {
	// Get the blocks
	ogblocks, err := r.GetBlocksByBranchID(fromBranchId)
	if err != nil {
		return err
	}
	// If Size of ogblocks is 0 we return.
	if len(*ogblocks) == 0 {
		return nil
	}
	var newBlocks []models.Block
	// Loop ogBlocks and create newBlocks (These are just objects) DB call is below.
	for _, block := range *ogblocks {
		newBlock := r.CloneBlockInstance(block, toBranchId, userID)
		newBlocks = append(newBlocks, newBlock)
	}
	// This is bulk create (To efficiently insert large number of records, pass a slice to the Create method.
	// GORM will generate a single SQL statement to insert all the data and backfill primary key values)
	// https://gorm.io/docs/create.html
	errBulkCreating := postgres.GetDB().Create(&newBlocks).Error
	if errBulkCreating != nil {
		return errBulkCreating
	}
	// All the Blocks are created now, We neeed to move things
	newblocks, errGettingClonedBlock := r.GetBlocksByBranchID(toBranchId)
	if errGettingClonedBlock != nil {
		return errGettingClonedBlock
	}
	// Get newly created blocks and create things on it.
	for _, newBlock := range *newblocks {
		if newBlock.ClonedFromBlockID != 0 {
			r.CopyBlockCommentsReelsReactionsToNewBlock(newBlock.ClonedFromBlockID, newBlock.ID, toBranchId)
		}
	}
	return nil
}

// This function take a Block Instace and Returns Copy of the Block without the ID but same UUID
func (r blockQuery) CloneBlockInstance(ogBlockInstance models.Block, branchID uint64, userID uint64) models.Block {
	var block models.Block
	block = ogBlockInstance
	block.ClonedFromBlockID = ogBlockInstance.ID
	block.ID = 0 // Reset the PK
	block.CanvasBranchID = &branchID
	/* Refactored : Delete

	block.Attributes = ogBlockInstance.Attributes
	block.UUID = ogBlockInstance.UUID
	block.CanvasRepositoryID = ogBlockInstance.CanvasRepositoryID
	block.CanvasBranchID = &branchID
	block.Version = ogBlockInstance.Version
	block.Type = ogBlockInstance.Type
	block.Rank = ogBlockInstance.Rank
	block.Children = ogBlockInstance.Children
	block.LastAttribution = ogBlockInstance.LastAttribution
	block.CreatedByID = ogBlockInstance.CreatedByID
	block.UpdatedByID = ogBlockInstance.UpdatedByID
	block.ClonedFromBlockID = ogBlockInstance.ID
	// Added this for new block!
	block.CommentCount = ogBlockInstance.CommentCount
	block.ReelCount = ogBlockInstance.ReelCount
	block.Reactions = ogBlockInstance.Reactions
	block.Contributors = ogBlockInstance.Contributors */
	return block
}

func (r blockQuery) UpdateBlockLastAttribution(branchID uint64, blockUUID string, attribution string) error {
	err := postgres.GetDB().Model(&models.Block{}).Where("canvas_branch_id = ? AND uuid = ?", branchID, blockUUID).Update("last_attribution", attribution).Error
	if err != nil {
		return err
	}
	return nil
}

func (r blockQuery) DeleteAllBlocksOnBranch(branchID uint64) error {
	var blocks []models.Block
	err := postgres.GetDB().Model(&models.Block{}).Where("canvas_branch_id = ?", branchID).Delete(blocks).Error
	if err != nil {
		return err
	}
	return nil
}

func (r blockQuery) CopyBlockCommentsReelsReactionsToNewBlock(sourceBlockID uint64, destinationBlockID uint64, destinationBranchID uint64) {
	r.CloneBlockThreadAndComments(sourceBlockID, destinationBlockID, destinationBranchID)
	r.CloneReels(sourceBlockID, destinationBlockID, destinationBranchID)
	r.CloneReactions(sourceBlockID, destinationBlockID, destinationBranchID)
}

func (r blockQuery) CloneBlockThreadAndComments(fromBlockID uint64, toBlockID uint64, branchID uint64) {
	var blockThreads *[]models.BlockThread
	_ = postgres.GetDB().Model(&models.BlockThread{}).Where("start_block_id = ? and is_archived = false", fromBlockID).Find(&blockThreads).Error
	// loop all the block threads and create a copy.
	for _, blockThread := range *blockThreads {
		r.CreateClonedBlockThread(fromBlockID, &blockThread, toBlockID, branchID)
	}
}

func (r blockQuery) CloneReels(fromBlockID uint64, toBlockID uint64, branchID uint64) {
	var reels *[]models.Reel
	_ = postgres.GetDB().Model(&models.Reel{}).Where("start_block_id = ? and is_archived = false", fromBlockID).Find(&reels).Error
	for _, reel := range *reels {
		r.CreateClonedReel(&reel, toBlockID, branchID)
	}
}

// Reactions Copy - BlockReaction /  BlockThreadReaction / BlockCommentReaction
// used to copy the reactions from MAIN branch to ROUGH Branch
// [x] BlockReaction (this function only does this)
// [x] BlockThreadReaction -> UpdateBlockThreadReactions
// [ ] BlockCommentReaction -> Pending
func (r blockQuery) CloneReactions(fromBlockID uint64, toBlockID uint64, branchID uint64) {
	// step 1 : Block Reaction
	var blockReactions *[]models.BlockReaction
	// Get all reactions on Parent Block
	_ = postgres.GetDB().Model(&models.BlockReaction{}).Where("block_id = ?", fromBlockID).Find(&blockReactions).Error
	// loop and clone each instance
	// for speed @todo: more to separate function
	for _, brInstance := range *blockReactions {
		fmt.Println("Cloning Block Reactions")
		ogBrID := brInstance.ID
		instance := brInstance
		instance.ClonedBlockReaction = ogBrID
		instance.UUID = uuid.New()
		instance.BlockID = toBlockID
		instance.CanvasBranchID = &branchID
		instance.ID = 0 // reset the PK
		_ = postgres.GetDB().Create(&instance)
	}
}

func (r blockQuery) CreateClonedBlockThread(fromBlockID uint64, og *models.BlockThread, newBlockID uint64, branchID uint64) {
	// Copy OG to Instance with new block id and branch
	fmt.Println("Cloning Block Thread")
	fmt.Println(og.ID)
	blockThreadClonedID := og.ID
	blockThreadUUID := og.UUID

	instance := og
	instance.StartBlockID = newBlockID
	instance.CanvasBranchID = branchID
	// Saving the reference to the original BlockThread.
	instance.ClonedFromThread = og.ID
	instance.UUID = blockThreadUUID
	instance.ID = 0 // Presumably we are setting ID to NIL
	_ = postgres.GetDB().Create(&instance)
	// Now we have a created a new BlockThread we also need to create Comments copy
	var comments *[]models.BlockComment
	// Get all comments from the parent thread id
	_ = postgres.GetDB().Model(&models.BlockComment{}).Where("thread_id = ? and is_archived = false", blockThreadClonedID).Find(&comments).Error
	fmt.Println(comments)
	r.UpdateBlockThreadReactions(fromBlockID, newBlockID, instance.ID, branchID)
	for _, blockThreadComment := range *comments {
		blockCommentInstance := r.CreateBlockThreadComment(&blockThreadComment, instance.ID, branchID)
		r.UpdateBlockThreadCommentReactions(fromBlockID, newBlockID, instance.ID, branchID, blockCommentInstance.ID)
	}
}

func (r blockQuery) CreateClonedReel(reel *models.Reel, newBlockID uint64, branchID uint64) {
	fmt.Println("Cloning Reel")
	ReelIDOG := reel.ID
	instance := reel
	instance.StartBlockID = newBlockID
	instance.CanvasBranchID = branchID
	// Saving the reference to the original BlockThread.
	instance.ClonedFromReel = ReelIDOG
	// using reel uuid instead of new to keep the reference for notifications
	instance.UUID = reel.UUID
	instance.ID = 0 // Presumably we are setting ID to NIL
	_ = postgres.GetDB().Create(&instance)
	// instance is newly created reel
	// Loop and clone comments
	var comments *[]models.ReelComment
	// Get all comments from the parent reel id
	_ = postgres.GetDB().Model(&models.ReelComment{}).Where("reel_id = ? and is_archived = false", ReelIDOG).Find(&comments).Error
	for _, reelComment := range *comments {
		r.CloneReelComment(&reelComment, instance.ID, branchID)
	}
}

// This is Problamatic as We need to Update the blockThreadReaction when Updating the Block Threads.
//Since The Block Thread Id are also changing
func (r blockQuery) UpdateBlockThreadReactions(parentBlockID, newBlockID, threadID, branchID uint64) {
	var blockThreadReaction *[]models.BlockThreadReaction
	// Get all reactions on Parent Block
	_ = postgres.GetDB().Model(&models.BlockThreadReaction{}).Where("block_id = ?", parentBlockID).Find(&blockThreadReaction).Error
	// loop and clone each instance
	// for speed @todo: more to separate function
	for _, btrInstance := range *blockThreadReaction {
		fmt.Println("Copy Block Thread Reactions to New ID")
		ogBrID := btrInstance.ID
		instance := btrInstance
		instance.ClonedBlockThreadReaction = ogBrID
		instance.UUID = uuid.New()
		instance.BlockID = newBlockID
		instance.CanvasBranchID = &branchID
		instance.BlockThreadID = threadID
		instance.ID = 0 // reset the PK
		_ = postgres.GetDB().Create(&instance)
	}
}

func (r blockQuery) UpdateBlockThreadCommentReactions(parentBlockID, newBlockID, threadID, branchID, commentID uint64) {
	var blockCommentReaction *[]models.BlockCommentReaction
	// Get all reactions on Parent Block
	_ = postgres.GetDB().Model(&models.BlockCommentReaction{}).Where("block_thread_id = ?", threadID).Find(&blockCommentReaction).Error
	// loop and clone each instance
	// for speed @todo: more to separate function
	for _, btrInstance := range *blockCommentReaction {
		fmt.Println("Copy Block Thread comment Reactions to New ID")
		ogBrID := btrInstance.ID
		instance := btrInstance
		instance.ClonedBlockCommentReaction = ogBrID
		instance.UUID = uuid.New()
		instance.BlockID = newBlockID
		instance.CanvasBranchID = &branchID
		instance.BlockThreadID = threadID
		instance.BlockCommentID = commentID
		instance.ID = 0 // reset the PK
		_ = postgres.GetDB().Create(&instance)
	}
}

func (r blockQuery) CloneReelComment(og *models.ReelComment, newReelID uint64, branchID uint64) {
	instance := og
	instance.ID = 0
	instance.ReelID = newReelID
	instance.ClonedFromReelComment = og.ID
	instance.UUID = uuid.New()
	_ = postgres.GetDB().Create(&instance)
}

func (r blockQuery) CreateBlockThreadComment(og *models.BlockComment, newBlockThreadID uint64, branchID uint64) *models.BlockComment {
	//var instance *models.BlockComment
	//instance.
	instance := og
	instance.ID = 0
	instance.ThreadID = newBlockThreadID
	instance.ClonedFromThreadComment = og.ID
	instance.UUID = uuid.New()
	_ = postgres.GetDB().Create(&instance)
	return instance
}
