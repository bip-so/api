package canvasbranch

import (
	"errors"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/blocks"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"time"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

func (s canvasBranchService) MergeRequestRejectValidation(roughBranchInstance *models.CanvasBranch, userID uint64) (*models.CanvasBranch, error) {
	// Check Rough
	//if !roughBranchInstance.IsRoughBranch {
	//	return nil, errors.New("Branch needs to be rough or draft branch.")
	//}
	//// Edge Case : Very rare UI Should prevent this case.
	//if roughBranchInstance.Committed {
	//	return nil, errors.New("Merge Request already exists for this branch. ")
	//}
	// Check main branch or not
	if roughBranchInstance.IsDefault {
		return nil, errors.New("cannot merge the default branch")
	}
	return roughBranchInstance, nil
}
func (s canvasBranchService) MergeRequestAcceptValidation(roughBranchInstance *models.CanvasBranch, userID uint64) (*models.CanvasBranch, error) {
	// Check main branch or not
	if roughBranchInstance.IsDefault {
		return nil, errors.New("cannot merge the default branch")
	}
	// Check Rough
	//if !roughBranchInstance.IsRoughBranch {
	//	return nil, errors.New("Branch needs to be rough or draft branch.")
	//}
	//// Edge Case : Very rare UI Should prevent this case.
	//if roughBranchInstance.Committed {
	//	return nil, errors.New("Merge Request already exists for this branch. ")
	//}
	return roughBranchInstance, nil
}

func (s *canvasBranchService) MergeRequestCreationValidation(roughBranchInstance *models.CanvasBranch, userID uint64) (*models.CanvasBranch, error) {
	// Todo: CC: We need to add a double check if the Requester has perms then don't show this error and Jump to Automerge
	if App.Repo.MergeRequestExists(roughBranchInstance.ID, userID) {
		return nil, errors.New("Merge Request Already Exists for this Branch and This User.! ")
	}
	//roughBranchInstance, errGettingBranch := App.Repo.Get(map[string]interface{}{"id": roughBranchID})
	//roughBranchInstance, errGettingBranch := queries.App.BranchQuery.GetBranchByID.GetBranchWithRepoAndStudio(roughBranchID)

	//if errGettingBranch != nil {
	//	return nil, errGettingBranch
	//}
	// Check Rough
	//if !roughBranchInstance.IsRoughBranch {
	//	return nil, errors.New("Branch needs to be rough or draft branch.")
	//}
	// Edge Case : Very rare UI Should prevent this case.
	//if roughBranchInstance.Committed {
	//	return nil, errors.New("Merge Request already exists for this branch. ")
	//}
	// Check main branch or not
	if roughBranchInstance.IsDefault {
		return nil, errors.New("cannot merge the default branch")
	}

	defer utils.TimeTrack(time.Now())

	return roughBranchInstance, nil
}

// ValidateMergeRequest :Merge Request Validation
// This will return (ParentBranchID, CanRequestMerge, canManageMerge, err)
// Is this a Rough Branch?
// Check for "CANVAS_BRANCH_CREATE_MERGE_REQUEST" Permission for this (Branch/User) on RoughBranch -> 400
// Check for "CANVAS_BRANCH_MANAGE_MERGE_REQUESTS" Permissions for this (Branch/User) on ParentBranch
// Flow1 : (Create Merge Request Instance and Set Committed = True) -> Merge request is created
func (s canvasBranchService) ValidateAndPrepareMergeRequest(roughBranchID uint64, userID uint64) (uint64, bool, bool, error) {
	fmt.Println("userIDuserIDuserIDuserIDuserID")
	fmt.Println(userID)

	if App.Repo.MergeRequestExists(roughBranchID, userID) {
		return 0, false, false, errors.New("Merge Request Already Exists for this Branch and This User.! ")
	}
	// Rough Branch Permissions
	//roughBranchInstance, errGettingBranch := App.Repo.Get(map[string]interface{}{"id": roughBranchID})
	roughBranchInstance, errGettingBranch := queries.App.BranchQuery.GetBranchWithRepoAndStudio(roughBranchID)
	if errGettingBranch != nil {
		return 0, false, false, errGettingBranch
	}

	// Check Rough
	if !roughBranchInstance.IsRoughBranch {
		return 0, false, false, errors.New("Branch needs to be rough or draft branch.")
	}

	// Edge Case : Very rare UI Should prevent this case.
	if roughBranchInstance.Committed {
		return 0, false, false, errors.New("Merge Request already exists for this branch. ")
	}

	// Can User CANVAS_BRANCH_CREATE_MERGE_REQUEST on This Branch
	canCreateMergeRequest, errGettingCreateMergeRequestPerms := permissions.App.Service.CanUserDoThisOnBranch(userID, roughBranchID, permissiongroup.CANVAS_BRANCH_CREATE_MERGE_REQUEST)
	if errGettingCreateMergeRequestPerms != nil {
		return 0, false, false, errors.New("Did not get valid permissions for this user")
	}

	canManageMergeRequest, errGettingManageMergeRequestPerms := permissions.App.Service.CanUserDoThisOnBranch(userID, *roughBranchInstance.RoughFromBranchID, permissiongroup.CANVAS_BRANCH_MANAGE_MERGE_REQUESTS)
	if errGettingManageMergeRequestPerms != nil {
		return 0, false, false, errors.New("Did not get valid permissions for this user")
	}
	fmt.Println("canCreateMergeRequest")
	fmt.Println(canCreateMergeRequest)
	fmt.Println("canManageMergeRequest")
	fmt.Println(canManageMergeRequest)
	// parentBranchID uint64, canRequest bool, canManageMergeRequest bool, err error
	return *roughBranchInstance.RoughFromBranchID, canCreateMergeRequest, canManageMergeRequest, nil
}

func (s canvasBranchService) StartMerge(roughBranch uint64, parentBranch uint64, loggedInUser *models.User, mergeStatus string, changesAccepted *map[string]interface{}) error {
	// Build Maps of IDs and UUI's to start with "UUID" : {ID, UUID}
	roughBranchLookupMap := App.Service.UtilBranchLookupMapMaker(roughBranch)
	parentBranchLookupMap := App.Service.UtilBranchLookupMapMaker(parentBranch)
	finalGitBlocks, errGitBlocks := App.Git.FetchBranchFromGit(loggedInUser, parentBranch)
	if errGitBlocks != nil {
		return errGitBlocks
	}
	finalGitBlocksMap := map[string]*models.Block{}
	finalGitBlocksRanksMap := map[string]int32{}
	for _, blk := range finalGitBlocks {
		finalGitBlocksRanksMap[blk.UUID.String()] = blk.Rank
		finalGitBlocksMap[blk.UUID.String()] = blk
	}

	// Now we have three lists
	// blocksUUIDToDelete: uuid's of the Blocks we need to delete on ParentBranch
	// blocksUUIDToCreate : Blocks we need to be present on Parent Branch
	// blocksUUIDToUpdate : Blocks we need to be Update on Parent Branch from the Rough Branch
	blocksUUIDToDelete := App.Service.UtilBlocksToDelete(roughBranchLookupMap, parentBranchLookupMap, finalGitBlocksMap)
	fmt.Println("IDs to delete", blocksUUIDToDelete)
	blocksUUIDToCreate := App.Service.UtilBlocksToCreate(roughBranchLookupMap, parentBranchLookupMap, finalGitBlocksMap)
	fmt.Println("IDs to create", blocksUUIDToCreate)
	blocksUUIDToUpdate := App.Service.UtilBlocksToUpdate(roughBranchLookupMap, parentBranchLookupMap)
	fmt.Println("IDs to update", blocksUUIDToUpdate)

	// Action: Delete these uuid[] from ParentBranch
	errDeletingBlock := s.ProcessBlocksToBeDeleted(blocksUUIDToDelete, parentBranchLookupMap, changesAccepted)

	if errDeletingBlock != nil {
		return errDeletingBlock
	}
	// Action: New Blocks We need to Create on Parent Branch
	// This is not always creation it's just adjustments
	errCreatingNewBlocks := s.ProcessBlocksToBeCreated(blocksUUIDToCreate, parentBranchLookupMap, roughBranchLookupMap, roughBranch, parentBranch, loggedInUser.ID, changesAccepted)
	if errCreatingNewBlocks != nil {
		return errCreatingNewBlocks
	}
	// Action: Update the Blocks on Master (Parent Branch)
	errUpdatingBlocks := s.ProcessBlocksUpdates(blocksUUIDToUpdate, parentBranchLookupMap, roughBranchLookupMap, roughBranch, parentBranch, loggedInUser.ID, finalGitBlocksMap, changesAccepted)
	if errUpdatingBlocks != nil {
		return errUpdatingBlocks
	}

	//go App.Repo.SyncBlockRanks(finalGitBlocksRanksMap)
	err := App.Repo.SyncBlockRanks(finalGitBlocksRanksMap)
	if err != nil {
		return err
	}

	// Clean UP
	//errPostMergeCleanup := s.PostMergeCleanUP(parentBranchLookupMap, roughBranchLookupMap, roughBranch, parentBranch, loggedInUser.ID)
	//if errPostMergeCleanup != nil {
	//	return errPostMergeCleanup
	//}
	return nil
}

func (s canvasBranchService) ProcessBlocksToBeCreated(
	blocksUUIDToCreate []string,
	parentBranchLookupMap map[string]BranchIDUUID,
	roughBranchLookupMap map[string]BranchIDUUID,
	roughBranchID uint64,
	parentBranchID uint64,
	userID uint64,
	changesAccepted *map[string]interface{},
) error {
	// @todo: Refactor this to do single update insted of 2 calls on each update
	for _, v := range blocksUUIDToCreate {
		// We are only processing the Blocks which are "Accepted"
		if changesAccepted != nil {
			if value, isBool := (*changesAccepted)[v].(bool); isBool && !value {
				continue
			}
		}
		// Get Block Instance By UUID and RoughBranchID
		blockDBInstance, err := queries.App.BlockQuery.GetBlock(map[string]interface{}{"uuid": v, "canvas_branch_id": roughBranchID})
		if err != nil {
			s.logg.Debug("Problem fetching block db instance")
			fmt.Println(err.Error())
		}
		//blocks = append(blocks, *blockDBInstance)
		errBlockUpdate := App.Repo.UpdateBlockMerge(blockDBInstance.ID, map[string]interface{}{"canvas_branch_id": &parentBranchID}) // commented out "updated_by_id": userID
		if errBlockUpdate != nil {
			s.logg.Debug("Problem fetching block db instance")
			fmt.Println(errBlockUpdate.Error())
		}
		App.Service.MoveThreadCommentsReelsReactionsMainBranchNewBlock(blockDBInstance.ID, &parentBranchID)
	}
	return nil
}

// ProcessBlocksToBeDeleted: Process block UUID []
// Create an array list of ID's and call for deletion by ID (PK)
func (s canvasBranchService) ProcessBlocksToBeDeleted(blocksUUIDToDelete []string, parentBranchLookupMap map[string]BranchIDUUID, changesAccepted *map[string]interface{}) error {
	var ids []uint64
	for _, v := range blocksUUIDToDelete {
		// We are only processing the Blocks which are "Accepted"

		if changesAccepted != nil {
			if value, isBool := (*changesAccepted)[v].(bool); isBool && !value {
				continue
			}
		}
		ids = append(ids, parentBranchLookupMap[v].id)
	}
	err := App.Repo.DeleteBlocks(ids)
	if err != nil {
		return err
	}

	return nil
}

type UpdateBlockData struct {
	parentBlockId uint64
	roughBlockId  uint64
}
type UpdateBlock struct {
	uuid string
	data UpdateBlockData
}

// Update Blocks from Rough Branch to Parent Branch
func (s canvasBranchService) ProcessBlocksUpdates(blocksUUIDToUpdate []string,
	parentBranchLookupMap map[string]BranchIDUUID,
	roughBranchLookupMap map[string]BranchIDUUID,
	roughBranchID uint64,
	parentBranchID uint64,
	loggedInUserID uint64,
	finalGitBlocksMap map[string]*models.Block,
	changesAccepted *map[string]interface{}) error {
	var thisBlock UpdateBlock
	for _, v := range blocksUUIDToUpdate {
		// This is skipping the blocks which are rejected in the Merge
		// We are only processing the Blocks which are "Accepted"
		if changesAccepted != nil {
			if value, isBool := (*changesAccepted)[v].(bool); isBool && !value {
				continue
			}
		}
		thisBlock.uuid = v

		thisBlock.data.parentBlockId = parentBranchLookupMap[v].id
		thisBlock.data.roughBlockId = roughBranchLookupMap[v].id
		// Debug
		// Get BLOCKS instances
		roughBranchBlockDBInstance, _ := queries.App.BlockQuery.GetBlock(map[string]interface{}{"id": thisBlock.data.roughBlockId})
		parentBranchBlockDBInstance, _ := queries.App.BlockQuery.GetBlock(map[string]interface{}{"id": thisBlock.data.parentBlockId})

		// We Replace following items
		// Updated by
		// Removed as per EM commentrs
		//parentBranchBlockDBInstance.UpdatedByID = userID
		parentBranchBlockDBInstance.Type = roughBranchBlockDBInstance.Type
		// Data
		parentBranchBlockDBInstance.Children = finalGitBlocksMap[v].Children
		parentBranchBlockDBInstance.Attributes = finalGitBlocksMap[v].Attributes
		// We need to merge the Contributors from both Blocks
		parentBranchBlockDBInstance.Rank = roughBranchBlockDBInstance.Rank
		// Recalculate Reel Count + CommentCount

		// Considering the difference of comment count only if parent is less in commentCount. Otherwise, duplicate comment count is added on every merge
		diffPresentCommentCount := roughBranchBlockDBInstance.CommentCount - parentBranchBlockDBInstance.CommentCount
		if roughBranchBlockDBInstance.CommentCount < parentBranchBlockDBInstance.CommentCount {
			diffPresentCommentCount = 0
		}
		diffPresentReelCount := roughBranchBlockDBInstance.ReelCount - parentBranchBlockDBInstance.ReelCount
		if roughBranchBlockDBInstance.ReelCount < parentBranchBlockDBInstance.ReelCount {
			diffPresentReelCount = 0
		}
		parentBranchBlockDBInstance.CommentCount = parentBranchBlockDBInstance.CommentCount + diffPresentCommentCount
		parentBranchBlockDBInstance.ReelCount = parentBranchBlockDBInstance.ReelCount + diffPresentReelCount
		mergedContributors := blocks.App.Service.BlockContributorMerge(parentBranchBlockDBInstance.Contributors, roughBranchBlockDBInstance.Contributors)
		parentBranchBlockDBInstance.Contributors = mergedContributors
		mergedReactions := blocks.App.Service.BlockReactionsMerge(parentBranchBlockDBInstance.Reactions, roughBranchBlockDBInstance.Reactions)
		parentBranchBlockDBInstance.Reactions = mergedReactions
		parentBranchBlockDBInstance.UpdatedAt = roughBranchBlockDBInstance.UpdatedAt
		// Save the Instance with merged data
		_ = App.Repo.UpdateBlockSimple(*parentBranchBlockDBInstance)
		// Move Things from Child Block to Parent block
		// Converting this in to Go Routine Let's see.
		App.Service.MoveThreadCommentsReelsReactionsMainBranchExistingBlock(thisBlock.data.roughBlockId, thisBlock.data.parentBlockId, parentBranchID)
	}

	return nil
}

// PostMergeCleanUP: Later with CreateMergeRequest
func (s canvasBranchService) PostMergeCleanUP(
	parentBranchLookupMap map[string]BranchIDUUID,
	roughBranchLookupMap map[string]BranchIDUUID,
	roughBranchID uint64,
	parentBranchID uint64,
	userID uint64) error {

	// delete the orphan blocks from rough branch after merge
	App.Repo.DeleteBlocksInBranchID(roughBranchID)
	// B1 / RB1

	// Move X from RB Block to X to Parent
	//- reaction
	//- comments
	//- mentions  json json deep
	//- contributers
	//- updateby ??? -> Don't Do or Do. (Accept)
	// Last edit by Check karo

	return nil
}

// CreateMergeRequest This is simple, we can do later
func (s canvasBranchService) CreateMergeRequest(roughBranchID uint64, parentBranchID uint64, user *models.User, message string) (*models.MergeRequest, error) {
	branchInstance, err := App.Service.GetCanvasBranchInstance(roughBranchID)
	if err != nil {
		return nil, err
	}

	if !branchInstance.Committed {
		if message == "" {
			message = "Creating merge request"
		}
		err := App.Git.CommitBranchToGit(user, roughBranchID, message)
		if err != nil {
			return nil, err
		}
	}

	var mr models.MergeRequest
	mergeRequestInstance := mr.NewMergeRequest(branchInstance.CanvasRepositoryID, roughBranchID, parentBranchID, message, "OPEN", user.ID, user.ID)
	mergeRequestInstaceDB, err := App.Repo.CreateMergeRequest(mergeRequestInstance)
	if err != nil {
		return nil, err
	}

	return mergeRequestInstaceDB, nil
}

// CloseMergeRequest: Git Internal will call
// PW : You can close a MR Instance now
func (s canvasBranchService) CloseMergeRequest(mergeRequestID uint64, status string) error {
	// We are closing
	var query = map[string]interface{}{"status": status}
	errUpdatingMergeRequest := App.Repo.UpdateMergeRequest(mergeRequestID, query)
	if errUpdatingMergeRequest != nil {
		return errUpdatingMergeRequest
	}
	return nil
}

// merge API - (
/// partially-accepted
// Fields Changes accepted MR Instance

//var changesAccepted *map[string]interface{} = nil
//if status == models.MERGE_REQUEST_PARTIALLY_ACCEPTED {
//changesAcceptedStr, isOK := p.Args["changesAccepted"].(string)
//if !isOK {
//return nil, errors.New("changesAccepted is not passed")
//}
//changesAccepted = new(map[string]interface{})
//err = json.Unmarshal([]byte(changesAcceptedStr), &changesAccepted)
//if err != nil {
//return nil, err
//}
//}

func (s canvasBranchService) MergeRequestService(mergeRequestInstance *models.MergeRequest, user *models.User) (*MergeResponseObject, error) {
	var diffResponse MergeResponseObject
	// We have Branch
	branch, _ := App.Service.GetCanvasBranchInstance(mergeRequestInstance.SourceBranchID)
	if branch != nil {
		serializedBranch := SimpleBranchDiffSerializer(branch)
		diffResponse.Branch = serializedBranch
	}
	cr := mergeRequestInstance.CanvasRepository
	diffResponse.CanvasRepo = MiniCanvasRepoSerializer{
		ID:           cr.ID,
		UUID:         cr.UUID.String(),
		Key:          cr.Key,
		CollectionID: cr.CollectionID,
		Name:         cr.Name,
		Position:     cr.Position,
		Icon:         cr.Icon,
		CreatedAt:    cr.CreatedAt,
		UpdatedAt:    cr.UpdatedAt,
		CreatedByID:  cr.CreatedByID,
		UpdatedByID:  cr.UpdatedByID,
	}
	if cr.ParentCanvasRepositoryID != nil {
		permissionList, _ := permissions.App.Service.CalculateSubCanvasRepoPermissions(user.ID, cr.StudioID, cr.CollectionID, *cr.ParentCanvasRepositoryID)
		diffResponse.CanvasRepo.Permission = permissionList[cr.ID][*cr.DefaultBranchID]
	} else {
		permissionList, _ := permissions.App.Service.CalculateCanvasRepoPermissions(user.ID, cr.StudioID, cr.CollectionID)
		diffResponse.CanvasRepo.Permission = permissionList[cr.ID][*cr.DefaultBranchID]
	}
	serialMr := SimpleMergeRequestSerializer(mergeRequestInstance)
	// Get blocks from source
	sourceBlocks, destinationBlocks := s.GetSourceAndDestinationBlocksFromGit(mergeRequestInstance, user)
	// We have Merge Request ID
	diffResponse.MergeRequest = serialMr
	diffResponse.SourceBlocks = sourceBlocks
	diffResponse.DestinationBlocks = destinationBlocks
	return &diffResponse, nil
}

func (s canvasBranchService) BlockBeforeMergeService(branch *models.CanvasBranch, authUser *models.User) (*DiffBeforeMergeResponseObject, error) {
	var diffResponse DiffBeforeMergeResponseObject
	// Get blocks from source
	sourceBlocks, _ := App.Service.GetAllBlockByBranchID(branch.ID)
	serializedSourceBlocks := BulkBlocksGetSerializerData(sourceBlocks)

	destinationBlocks, err := App.Git.FetchCommitFromGit(authUser, branch.CreatedFromCommitID, branch.CanvasRepositoryID)
	if err != nil {
		return nil, err
	}
	serializedDestinationBlocks := BulkGitBlocksSerializerData(destinationBlocks)

	diffResponse.SourceBlocks = &serializedSourceBlocks
	diffResponse.DestinationBlocks = serializedDestinationBlocks
	return &diffResponse, nil
}

func (s canvasBranchService) TransformBlock(gitBlock models.Block, bdBlock models.Block) models.Block {
	return models.Block{}
}

func (s canvasBranchService) GetSourceAndDestinationBlocksFromGit(mergeRequest *models.MergeRequest, user *models.User) (*[]BulkBlocks, *[]BulkBlocks) {
	var sourceBlocks []*models.Block
	var destinationBlocks []*models.Block
	if mergeRequest.Status == models.MERGE_REQUEST_OPEN {
		sourceBlocks, _ = App.Git.FetchBranchFromGit(user, mergeRequest.SourceBranchID)
		if mergeRequest.SourceBranch.CreatedFromCommitID == "" {
			destinationBlocks, _ = App.Git.FetchBranchFromGit(user, mergeRequest.DestinationBranchID)
		} else {
			destinationBlocks, _ = App.Git.FetchCommitFromGit(user, mergeRequest.SourceBranch.CreatedFromCommitID, mergeRequest.CanvasRepositoryID)
		}
	} else if mergeRequest.Status == models.MERGE_REQUEST_REJECTED {
		sourceBlocks, _ = queries.App.BlockQuery.GetSourceBlocksByBranchID(mergeRequest.SourceBranchID)
		destinationBlocks, _ = App.Git.FetchBranchFromGit(user, mergeRequest.DestinationBranchID)
	} else {
		sourceBlocks, _ = App.Git.FetchCommitFromGit(user, mergeRequest.SourceCommitID, mergeRequest.CanvasRepositoryID)
		destinationBlocks, _ = App.Git.FetchCommitFromGit(user, mergeRequest.DestinationCommitID, mergeRequest.CanvasRepositoryID)
	}
	serializedSourceBlocks := BulkGitBlocksSerializerData(sourceBlocks)
	serializedDestinationBlocks := BulkGitBlocksSerializerData(destinationBlocks)
	return serializedSourceBlocks, serializedDestinationBlocks
}

// MERGE_REQUEST_OPEN
func (s canvasBranchService) GetMergeRequestsByBranch(branchID uint64, userID uint64) (*[]models.MergeRequest, error) {

	modUserIDs, err := notifications.App.Repo.GetCanvasBranchModeratorsUserIDs(branchID)
	if err != nil {
		return nil, err
	}

	query := map[string]interface{}{"destination_branch_id": branchID, "status": models.MERGE_REQUEST_OPEN}
	if !utils.SliceContainsInt(modUserIDs, userID) {
		query["created_by_id"] = userID
	}

	instances, err := App.Repo.GetAllMergeRequests(query) //May need to add more filter as PW.
	if err != nil {
		return nil, err
	}
	return instances, nil
}
