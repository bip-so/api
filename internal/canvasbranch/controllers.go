package canvasbranch

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"

	"gitlab.com/phonepost/bip-be-platform/internal/studio"
	"gorm.io/gorm"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

var AllowedScope = []string{"", "create", "update", "delete"}

type Controller interface {
	Get(id uint64)
	Create(user models.User, studioId uint64, collectionID uint64, name string, icon string) (*models.CanvasBranch, error)
	Update(id uint64)
	Delete(id uint64)
	BlocksManager(body CanvasBlockPost)
}

func (c canvasBranchController) Create(user models.User, body newCanvasBranchPost, studioID uint64) (*models.CanvasBranch, error) {

	err := App.Git.CommitBranchToGit(&user, body.FromCanvasBranchID, "Creating branch")
	if err != nil {
		return nil, err
	}

	key := utils.NewNanoid()

	commitID, err := App.Git.CreateBranch(&user, body.FromCanvasBranchID, key)
	if err != nil {
		return nil, err
	}

	// Create a Canvas Branch
	cb := App.Service.EmptyCanvasBranchInstance()
	cb.CanvasRepositoryID = body.CanvasRepoID
	cb.Name = key
	cb.FromBranchID = &body.FromCanvasBranchID
	cb.CreatedByID = user.ID
	cb.UpdatedByID = user.ID
	cb.Key = key
	cb.CreatedFromCommitID = commitID

	cbInstance, errCreatingBranch := App.Repo.Create(cb)
	if errCreatingBranch != nil {
		logger.Info(fmt.Sprintf("CR Instance Error:%v\n", errCreatingBranch.Error()))
		return nil, errCreatingBranch
	}

	// Creating canvas branch permission
	err = queries.App.PermsQuery.CreateDefaultCanvasBranchPermission(body.CollectionID, user.ID, studioID, body.CanvasRepoID, cb.ID, &body.ParentCanvasRepoID)
	if err != nil {
		return nil, err
	}

	// Creating a blank block
	//var block models.Block
	//firstContrib := App.Service.BlockContributorFirst(user, body.FromCanvasBranchID)
	//newBlock, errCreatingBlockInst := block.NewBlock(
	//	uuid.New(),
	//	user.ID,
	//	cb.CanvasRepositoryID,
	//	cb.ID,
	//	nil,
	//	1,
	//	2,
	//	models.BlockTypeText,
	//
	//	models.MyFirstBlockJson(),
	//	models.MyFirstEmptyBlockJson(),
	//	models.MyFirstEmptyBlockJson(),
	//	firstContrib)
	//if errCreatingBlockInst != nil {
	//	return nil, errCreatingBlockInst
	//}
	//errNewBlock := blocks.App.Service.Create(newBlock)

	//if errNewBlock != nil {
	//	return nil, errNewBlock
	//}

	errCloneBlocks := App.Service.CloneBlocks(body.FromCanvasBranchID, cb.ID, user.ID)
	// We then send the roughBranchID or 0 with error
	if errCloneBlocks != nil {
		return nil, errCloneBlocks
	}

	return cbInstance, nil
}

func (c canvasBranchController) placeHolder() {
	//_ = App.Service.Create()
	fmt.Println("Here")

	//App.Repo.Create()
}

// Returns Rought branch from this branch
func (c canvasBranchController) RoughBranchBuilder(user *models.User, branchID uint64, studioID uint64, requestBody NewDraftBranchPost) (*models.CanvasBranch, error) {
	err := App.Git.CommitBranchToGit(user, branchID, "Creating branch")
	if err != nil {
		return nil, err
	}
	key := utils.NewNanoid()
	//parentBranchInstance, err := App.Repo.Get(map[string]interface{}{"id": branchID})
	parentBranchInstance, err := queries.App.BranchQuery.GetBranchWithRepoAndStudio(branchID)

	if err != nil {
		return nil, err
	}
	// branchName := parentBranchInstance.CanvasRepository.Name + key
	branchName := key
	commitID, err := App.Git.CreateBranch(user, branchID, branchName)
	if err != nil {
		return nil, err
	}
	// We need to create a new branch which is copy for this branch (branchID)
	roughBranch, err := App.Service.CreateRoughBranch(parentBranchInstance, user.ID, branchName, key, commitID)
	if err != nil {
		return nil, err
	}
	// We need to then copy all the blocks in branchID to roughBranchID (copy should keep the UUID of block same)
	errCloneBlocks := App.Service.CloneBlocks(branchID, roughBranch.ID, user.ID)
	// We then send the roughBranchID or 0 with error
	if errCloneBlocks != nil {
		return nil, errCloneBlocks
	}
	// Creating canvas branch permission for RB
	member, err := queries.App.MemberQuery.GetMember(map[string]interface{}{"user_id": user.ID, "studio_id": studioID})
	if err == gorm.ErrRecordNotFound || member.HasLeft {
		joinStudioBulk := studio.JoinStudioBulkPost{UsersAdded: []uint64{user.ID}}
		_, err = studio.App.Controller.JoinStudioInBulkController(joinStudioBulk, studioID, user.ID)
		member, err = queries.App.MemberQuery.GetMember(map[string]interface{}{"user_id": user.ID, "studio_id": studioID})
	}
	err = permissions.App.Service.CreateCustomCanvasBranchPermission(
		requestBody.CollectionID, user.ID, studioID, requestBody.CanvasRepoID, roughBranch.ID, &requestBody.ParentCanvasRepoID, "pg_canvas_branch_moderate", member.ID)
	if err != nil {
		return nil, err
	}

	return roughBranch, nil
}

// This function will get data in format of CanvasBlockPost
// We need to loop over the data.Blocks and process each block as fast as possible.
// Order of the Block -> Processing
// Slate -> (Deadlock)
// UUID - Duplication
func (c canvasBranchController) BlocksManager(loggedInUserInstance models.User, branchId uint64, data CanvasBlockPost, isAssociationAPI bool, studioInstance *models.Studio) (BlockBulkResponseMetadata, error) {
	var meta BlockBulkResponseMetadata
	// Validation
	validationError := App.Service.ValidateBlocksData(branchId, data, isAssociationAPI)
	if validationError != nil {
		return meta, validationError
	}
	var createdBlockIDs = []uint64{}
	var updatedBlockIDs = []uint64{}
	var failedToUpdatedBlockIDs = []uint64{}
	var failedToDeleteBlockIDs = []uint64{}
	deletedBlocksCount := 0
	failedBlocksCount := 0

	createdBlocks := 0
	updatedBlocks := 0
	deletedBlocks := 0
	// Processing
blockProcessingLoop:
	for _, blockItem := range data.Blocks {
		switch {
		case blockItem.Scope == "create":
			id, err := App.Service.ProcessCreateBlock(loggedInUserInstance, branchId, blockItem)
			if err != nil {
				failedBlocksCount++
			}
			createdBlockIDs = append(createdBlockIDs, id)
			createdBlocks = createdBlocks + 1
			break
		case blockItem.Scope == "update":
			id, err := App.Service.ProcessUpdateBlock(loggedInUserInstance, branchId, blockItem)
			if err != nil {
				failedToUpdatedBlockIDs = append(failedToUpdatedBlockIDs, id)
			} else {
				updatedBlockIDs = append(updatedBlockIDs, id)
			}
			updatedBlocks = updatedBlocks + 1
			break

		case blockItem.Scope == "delete":
			err := App.Service.ProcessDeleteBlock(loggedInUserInstance.ID, branchId, blockItem)
			if err != nil {
				failedToDeleteBlockIDs = append(failedToDeleteBlockIDs, blockItem.ID)
			}
			deletedBlocksCount++
			deletedBlocks = deletedBlocks + 1
			break

		default:
			continue blockProcessingLoop
		}
	}
	meta = BlockBulkResponseMetadata{
		CreatedBlockIDs:         createdBlockIDs,
		UpdatedBlockIDs:         updatedBlockIDs,
		FailedToUpdatedBlockIDs: failedToUpdatedBlockIDs,
		FailedToDeleteBlockIDs:  failedToDeleteBlockIDs,
		DeletedBlocksCount:      deletedBlocksCount,
		FailedBlocksCount:       failedBlocksCount,
	}
	// Studio Logger
	if studioInstance.FeatureFlagHasXP {
		// Now Update the data
		SetStudioLogDataRedis(loggedInUserInstance.ID, studioInstance.ID, branchId, createdBlocks, updatedBlocks, deletedBlocks)
	}

	return meta, nil
}

// This function will get data in format of CanvasBlockPost
// We need to loop over the data.Blocks and process each block as fast as possible.
// Order of the Block -> Processing
// Slate -> (Deadlock)
// UUID - Duplication

/*
	Read Me: This Function will receive an Existing Branch
	Create if any blocks need to be created
	In this case if scope is "" we will copy the block as is to new branch.
	We will still honour Edit and Delete but all on fresh blocks on new branch
*/
// Commenting this function to check is this is used anywhere?
//
//func (c canvasBranchController) RoughBranchBlocksManager1(userInstance models.User, branchId uint64, roughBranchID uint64, data CanvasBlockPost) (BlockBulkResponseMetadata, error) {
//	var meta BlockBulkResponseMetadata
//	// Validation
//
//	validationError := App.Service.ValidateBlocksData(branchId, data, false)
//	if validationError != nil {
//		return meta, validationError
//	}
//
//	// Clone All blocks from Parent Branch to New Branch
//	errCloningblocks := App.Service.CloneBlocks(branchId, roughBranchID, userInstance.ID)
//	if errCloningblocks != nil {
//		return meta, errCloningblocks
//	}
//
//	// We have cloned all the blocks to new branch
//	var createdBlockIDs = []uint64{}
//	var updatedBlockIDs = []uint64{}
//	var failedToUpdatedBlockIDs = []uint64{}
//	var failedToDeleteBlockIDs = []uint64{}
//	deletedBlocksCount := 0
//	failedBlocksCount := 0
//	// Processing rest of Data for only Create and Delete and Update but on Freshly Cloned Blocks
//blockProcessingLoop:
//	for _, blockItem := range data.Blocks {
//		switch {
//		case blockItem.Scope == "create":
//			// Create a block on new branch
//			id, err := App.Service.ProcessCreateBlock(userInstance, roughBranchID, blockItem)
//			if err != nil {
//				failedBlocksCount++
//			}
//			createdBlockIDs = append(createdBlockIDs, id)
//			break
//		case blockItem.Scope == "update":
//			// Need to Change this as the ID is changed.
//			// We have to find the Block with UUID and BranchID
//			id, err := App.Service.ProcessUpdateClonedBlock(userInstance.ID, roughBranchID, blockItem)
//			if err != nil {
//				failedToUpdatedBlockIDs = append(failedToUpdatedBlockIDs, id)
//			} else {
//				updatedBlockIDs = append(updatedBlockIDs, id)
//			}
//			break
//		case blockItem.Scope == "delete":
//			// Need to Change this as the ID is changed.
//			// We have to find the Block with UUID and BranchID
//			// and Delete
//			err := App.Service.ProcessDeleteClonedBlock(userInstance.ID, roughBranchID, blockItem)
//			if err != nil {
//				failedToDeleteBlockIDs = append(failedToDeleteBlockIDs, blockItem.ID)
//			}
//			deletedBlocksCount++
//			break
//		default:
//			continue blockProcessingLoop
//		}
//	}
//
//	meta = BlockBulkResponseMetadata{
//		CreatedBlockIDs:         createdBlockIDs,
//		UpdatedBlockIDs:         updatedBlockIDs,
//		FailedToUpdatedBlockIDs: failedToUpdatedBlockIDs,
//		FailedToDeleteBlockIDs:  failedToDeleteBlockIDs,
//		DeletedBlocksCount:      deletedBlocksCount,
//		FailedBlocksCount:       failedBlocksCount,
//	}
//
//	return meta, nil
//}

//func (c canvasBranchController) AnonymousBranchesSearchController(search string) (*[]GetCanvasBranchSerializer, error) {
//
//	branches, err := App.Repo.GetBranchesAnonymous(body.CanvasID, []string{"view", "edit", "comment"})
//	if err != nil {
//		return nil, err
//	}
//	data := MultiSerializeCanvasBranch(*branches)
//	return data, nil
//}

func (c canvasBranchController) AnonymousBranchesController(body GetCanvasBranches) (*[]GetCanvasBranchSerializer, error) {
	branches, err := App.Repo.GetBranchesAnonymous(body.CanvasID, []string{"view", "edit", "comment"})
	if err != nil {
		return nil, err
	}
	data := MultiSerializeCanvasBranch(*branches)
	return data, nil
}

func (c canvasBranchController) AuthGetCanvasBranchesController(body GetCanvasBranches, user *models.User, studioID uint64) (*[]GetCanvasBranchSerializer, error) {
	var permissionsList map[uint64]map[uint64]string
	accessCanvasBranch := &[]GetCanvasBranchSerializer{}

	branches, err := App.Repo.GetBranches(map[string]interface{}{"canvas_repository_id": body.CanvasID, "is_archived": false, "is_merged": false})
	if err != nil {
		return nil, err
	}

	if body.ParentCanvasID == 0 {
		permissionsList, err = permissions.App.Service.CalculateCanvasRepoPermissions(user.ID, studioID, body.CollectionID)
	} else {
		permissionsList, err = permissions.App.Service.CalculateSubCanvasRepoPermissions(user.ID, studioID, body.CollectionID, body.ParentCanvasID)
	}
	if err != nil {
		return nil, err
	}

	branchPermissions := permissionsList[body.CanvasID]
	for _, branch := range *branches {
		perm := branchPermissions[branch.ID]
		if utils.Contains(permissiongroup.UserAccessCanvasPermissionsList, perm) {
			canvasBranchSerializedData := SerializeCanvasBranch(&branch)
			canvasBranchSerializedData.Permission = perm
			*accessCanvasBranch = append(*accessCanvasBranch, *canvasBranchSerializedData)
		} else {
			if branch.PublicAccess != "private" {
				canvasBranchSerializedData := SerializeCanvasBranch(&branch)
				canvasBranchSerializedData.Permission = permissiongroup.PGCanvasNone().SystemName
				*accessCanvasBranch = append(*accessCanvasBranch, *canvasBranchSerializedData)
			}
		}
	}
	return accessCanvasBranch, nil
}
