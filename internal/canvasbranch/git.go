package canvasbranch

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"strings"
	"time"

	"github.com/google/uuid"
	"gitlab.com/phonepost/bip-be-platform/internal/message"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gorm.io/datatypes"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/user"
	stores "gitlab.com/phonepost/bip-be-platform/pkg/stores/git"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

func (g gitService) CommitBranchToGit(usr *models.User, branchID uint64, message string) error {
	fmt.Println(usr, branchID, message)
	branch, err := queries.App.BranchQuery.GetBranchByID(branchID)
	if err != nil {
		return err
	}

	if branch.Committed {
		return nil
	}

	blocks, err := queries.App.BlockQuery.GetBlocksByBranchID(branch.ID)
	if err != nil {
		return err
	}

	blockViews := []interface{}{}
	blockData := BulkBlocksGetSerializerData(blocks)
	for _, blk := range blockData {
		blockViews = append(blockViews, blk)
	}
	fromBranchName := ""
	if !branch.IsDefault {
		if branch.RoughFromBranch != nil {
			fromBranchName = branch.RoughFromBranch.Name
		} else if branch.RoughFromBranch == nil && branch.FromBranch == nil {
			return errors.New("invalid branch")
		}
	}
	//fmt.Println("-----Megamega---")
	//fmt.Println(user.UUID.String())
	//fmt.Println(user.Username)
	//fmt.Println(branch.CanvasRepository.Studio.UUID.String())
	//fmt.Println(user.UUID.String(), user.Username, branch.CanvasRepository.Studio.UUID.String(), branch.CanvasRepository.UUID.String())
	// @todo: Discuss PW changes user.FullName to user.Username
	gitUserRepo, err := stores.CreateGitUserRepo(usr.UUID.String(), usr.Username, branch.CanvasRepository.Studio.UUID.String(), branch.CanvasRepository.UUID.String())
	if err != nil {
		fmt.Println("Broke @ CreateGitUserRepo")
		return err
	}

	err = gitUserRepo.CreateSnapshot(&blockViews, fromBranchName, branch.Name, message)
	if err != nil {
		return err
	}

	err = App.Repo.SetBranchCommited(branch.ID, true)
	if err != nil {
		return err
	}
	return nil
}

func (g gitService) RejectMergeRequest(authUser *models.User, canvasBranchID, parentBranchID uint64, mergeRequestID uint64) (bool, error) {
	// Merge Request
	mergeReq, err := App.Repo.GetMergeRequest(map[string]interface{}{"id": mergeRequestID})
	if err != nil {
		return false, err
	}

	if mergeReq.Status != models.MERGE_REQUEST_OPEN {
		return false, errors.New("Merge request already merged")
	}

	query := map[string]interface{}{
		"status":            models.MERGE_REQUEST_REJECTED,
		"updated_by_id":     authUser.ID,
		"closed_by_user_id": authUser.ID,
		"closed_at":         time.Now(),
	}
	errUpdateMR := App.Repo.UpdateMergeRequest(mergeReq.ID, query)
	if errUpdateMR != nil {
		return false, errUpdateMR
	}
	go App.Repo.SetBranchCommited(mergeReq.SourceBranchID, false)

	go func() {
		extraData := notifications.NotificationExtraData{
			CollectionID:   mergeReq.CanvasRepository.CollectionID,
			CanvasRepoID:   mergeReq.CanvasRepositoryID,
			CanvasBranchID: mergeReq.DestinationBranchID,
			Status:         models.MERGE_REQUEST_REJECTED,
			Message:        mergeReq.CommitMessage,
		}
		contentObject := models.MERGEREQUEST
		notifications.App.Service.PublishNewNotification(notifications.MergeRequestedUpdate, authUser.ID, nil,
			&mergeReq.CanvasRepository.Studio.ID, nil, extraData, &mergeReq.ID, &contentObject)
	}()

	return true, nil
}

func (g gitService) MergeMergeRequest(authUser *models.User, canvasBranchID, parentBranchID uint64, body MergeRequestAcceptPartialPost, mr *models.MergeRequest, isAutoMerge bool) (bool, error) {
	if body.MergeStatus == models.MERGE_REQUEST_REJECTED {
		return false, errors.New("Please use Reject Merge Request API.")
	}
	if body.MergeStatus == models.MERGE_REQUEST_OPEN {
		return false, errors.New("MergeRequest is already Can't send  `OPEN` Again. ")
	}
	if !utils.SliceContainsItem([]string{models.MERGE_REQUEST_ACCEPTED, models.MERGE_REQUEST_REJECTED, models.MERGE_REQUEST_PARTIALLY_ACCEPTED}, body.MergeStatus) {
		return false, errors.New("Please use correct merge status : Allowed `ACCEPTED`,`REJECTED`,`PARTIALLY_ACCEPTED`")
	}
	if mr.Status != models.MERGE_REQUEST_OPEN {
		return false, errors.New("Merge request already merged")
	}

	gitUserRepo, err := stores.CreateGitUserRepo(authUser.UUID.String(), authUser.Username, mr.CanvasRepository.Studio.UUID.String(), mr.CanvasRepository.UUID.String())
	if err != nil {
		return false, err
	}

	// ??? Creating a Record on GIY?
	commitID, srcCommitID, destCommitID, err := gitUserRepo.MergeMergeRequest(mr.DestinationBranch.Name, mr.SourceBranch.Name, body.MergeStatus, mr.SourceBranch.CreatedFromCommitID, body.ChangesAccepted)
	if err != nil {
		return false, err
	}

	err = App.Service.StartMerge(canvasBranchID, parentBranchID, authUser, body.MergeStatus, body.ChangesAccepted)
	if err != nil {
		return false, err
	}
	//go App.Service.StartMerge(canvasBranchID, parentBranchID, authUser, body.MergeStatus, body.ChangesAccepted)

	var query map[string]interface{}
	if body.MergeStatus == models.MERGE_REQUEST_PARTIALLY_ACCEPTED {
		query = map[string]interface{}{
			"status":                body.MergeStatus,
			"updated_by_id":         authUser.ID,
			"changes_accepted":      body.ChangesAccepted,
			"commit_id":             commitID,
			"source_commit_id":      srcCommitID,
			"destination_commit_id": destCommitID,
			"closed_by_user_id":     authUser.ID,
			"closed_at":             time.Now(),
		}
	} else {
		query = map[string]interface{}{
			"status":                body.MergeStatus,
			"updated_by_id":         authUser.ID,
			"commit_id":             commitID,
			"source_commit_id":      srcCommitID,
			"destination_commit_id": destCommitID,
			"closed_by_user_id":     authUser.ID,
			"closed_at":             time.Now(),
		}
	}
	err = App.Repo.UpdateMergeRequest(mr.ID, query)
	if err != nil {
		return false, err
	}

	// Todo: Double Check
	go func() {

		// THIS IS DECISION PENDING...PARTIAL_MERGE and FULL MERGE.
		// if body.MergeStatus == models.MERGE_REQUEST_ACCEPTED {
		// 	App.Git.DeleteBranch(authUser, mr.SourceBranchID)
		// }

		// if mergeStatus == models.MERGE_REQUEST_ACCEPTED {
		// 	models.SendBipMarkAuthorAndMentionNotificationOnMerge(ctx, mergeReq)
		// }
		fmt.Println("merge request", mr)
		// create notification
		mrString, _ := json.Marshal(mr)
		extraData := notifications.NotificationExtraData{
			CollectionID:   mr.CanvasRepository.CollectionID,
			CanvasRepoID:   mr.CanvasRepositoryID,
			CanvasBranchID: mr.DestinationBranchID,
			Status:         body.MergeStatus,
			Data:           string(mrString),
			Message:        mr.CommitMessage,
		}
		contentObject := models.MERGEREQUEST
		if isAutoMerge {
			notifications.App.Service.PublishNewNotification(notifications.CanvasMerged, authUser.ID, nil,
				&mr.CanvasRepository.Studio.ID, nil, extraData, &mr.ID, &contentObject)
		} else {
			notifications.App.Service.PublishNewNotification(notifications.MergeRequestedUpdate, authUser.ID, nil,
				&mr.CanvasRepository.Studio.ID, nil, extraData, &mr.ID, &contentObject)
		}
		notifications.App.Service.PublishNewNotification(notifications.BipMarkMessageAdded, authUser.ID, nil,
			&mr.CanvasRepository.StudioID, nil, extraData, &mr.ID, &contentObject)

		// Delete All Rough Branch Blocks doing in go routine
		_ = queries.App.BlockQuery.DeleteAllBlocksOnBranch(mr.SourceBranchID)
		// Delete the Rough Branch on the Success of Above Execution.
		fmt.Println("Deleting branch with merge reqeust", mr.SourceBranchID, mr)
		_ = App.Repo.DeleteBranch(mr.SourceBranchID)
	}()

	// cache related
	// update branch last edited
	// models.UpdateBranchesLastEdited(ctx, []string{mergeReq.ToBranchID, mergeReq.FromBranchID})

	// calling invalidation signal
	// keyData := []string{mergeReq.PageID}

	// cache.InvalidateCache(ctx, "canvas", keyData, true)
	// Post Success Flow:
	App.Service.InvalidateBranchBlocks(mr.DestinationBranchID)
	App.Service.InvalidateBranchBlocks(mr.SourceBranchID)
	App.Repo.UpdateBranchLastEdited(mr.DestinationBranchID)

	defer utils.TimeTrack(time.Now())
	return true, nil
}

func (g gitService) CreateBranch(user *models.User, branchID uint64, newBranchName string) (string, error) {

	branch, err := queries.App.BranchQuery.GetBranchByID(branchID)
	if err != nil {
		return "", err
	}

	gitUserRepo, err := stores.CreateGitUserRepo(user.UUID.String(), user.Username, branch.CanvasRepository.Studio.UUID.String(), branch.CanvasRepository.UUID.String())
	if err != nil {
		return "", err
	}

	commitID, err := gitUserRepo.CreateBranch(branch.Name, newBranchName)
	if err != nil {
		return "", err
	}

	return commitID, nil
}

func (g gitService) DeleteBranch(user *models.User, branchID uint64) error {

	branch, err := queries.App.BranchQuery.GetBranchByID(branchID)
	if err != nil {
		return err
	}

	gitUserRepo, err := stores.CreateGitUserRepo(user.UUID.String(), user.FullName, branch.CanvasRepository.Studio.UUID.String(), branch.CanvasRepository.UUID.String())
	if err != nil {
		return err
	}

	err = gitUserRepo.DeleteBranch(branch.Name)
	if err != nil {
		return err
	}

	return nil
}

func (g gitService) FetchBranchHistoryFromGit(usr *models.User, branchID uint64, startCommitID string) ([]*stores.GitLog, *[]models.User, string, error) {
	branch, err := queries.App.BranchQuery.GetBranchByID(branchID)
	if err != nil {
		return nil, nil, "", err
	}

	gitUserRepo, err := stores.CreateGitUserRepo(usr.UUID.String(), usr.Username, branch.CanvasRepository.Studio.UUID.String(), branch.CanvasRepository.UUID.String())
	if err != nil {
		return nil, nil, "", err
	}
	gitLogs, next, err := gitUserRepo.FetchBranchHistoryFromGit(branch.Name, startCommitID)
	if err != nil {
		return nil, nil, "", err
	}

	userUUIDs := []string{}
	for i, log := range gitLogs {
		userID := strings.Split(log.AuthorEmail, "@")[0]
		gitLogs[i].UserID = userID
		if !utils.SliceContainsItem(userUUIDs, userID) {
			userUUIDs = append(userUUIDs, userID)
		}
	}

	users, err := user.App.Repo.GetUsersByUUIDs(userUUIDs)
	if err != nil {
		return nil, nil, "", err
	}

	return gitLogs, users, next, nil
}

func (g gitService) FetchAndUpdateLatestBlockAttributionsForBranch(usr *models.User, branchID uint64) error {

	branch, err := queries.App.BranchQuery.GetBranchByID(branchID)
	if err != nil {
		return err
	}

	gitUserRepo, err := stores.CreateGitUserRepo(usr.UUID.String(), usr.FullName, branch.CanvasRepository.Studio.UUID.String(), branch.CanvasRepository.UUID.String())
	if err != nil {
		return err
	}
	gitAttributions, err := gitUserRepo.FetchLatestBlockAttributionsForBranch(branch.Name)
	if err != nil {
		return err
	}

	userIDs := []string{}
	for _, attr := range gitAttributions {
		userID := strings.Split(attr.AuthorEmail, "@")[0]
		attr.UserID = userID
		userIDs = append(userIDs, userID)
	}

	users, err := user.App.Repo.GetUsersByUUIDs(userIDs)
	if err != nil {
		return err
	}

	userMap := map[string]*models.User{}
	for _, user := range *users {
		userMap[user.UUID.String()] = &user
	}

	for _, attr := range gitAttributions {
		blockAttr := models.BlockAttribution{
			UserUUID:  attr.UserID,
			Username:  userMap[attr.UserID].Username,
			FullName:  userMap[attr.UserID].FullName,
			TimeStamp: attr.UpdatedAt,
			AvatarUrl: userMap[attr.UserID].AvatarUrl,
		}
		data, _ := json.Marshal(blockAttr)
		fmt.Println(data)
		// Go we need this?
		// App.Repo.UpdateBlockLastAttribution(branchID, attr.BlockID, string(data))
	}

	return nil
}

//  Todo: We need to check if the attributions is empty and build the attributions and then send again
/*
		When get ATTRIB FOR A BRANCH  // NON-ROUGH BRANCH
	    CHECK EMPTY
	    IF EMPTY
	    THEN WE CHECK IF REPO IS PUBLISHED
	        IF THEN
	            WE NEED TO CHECK IF REPO IS_COMMITTED
	                WE COMMIT THIS BRANCH IN THE NAME OF THE CREATED BY
		COMMIT BRANCH TO GIT
*/
func (g gitService) FetchAllAttributionsForBranch(usr *models.User, branchID uint64) (*[]models.Attribution, error) {
	// get branch data
	branch, err := queries.App.BranchQuery.GetBranchByID(branchID)
	if err != nil {
		return nil, err
	}

	// Attributions from DB
	attributions, err := queries.App.AttributionQuery.GetAllAttributionsForBranch(branchID)
	if err != nil {
		return nil, err
	}

	// Setting Defaults to Admin and Change only if User is found.
	userID := "admin"
	userName := "admin"
	// Logged user in user
	if usr != nil {
		userID = usr.UUID.String()
		userName = usr.Username
	}

	// Creates a GitUserRepo
	gitUserRepo, err := stores.CreateGitUserRepo(userID, userName, branch.CanvasRepository.Studio.UUID.String(), branch.CanvasRepository.UUID.String())
	if err != nil {
		return nil, err
	}
	gitAttributions, startCommitID, err := gitUserRepo.FetchAllAttributionsForBranch(branch.Name, branch.LastSyncedAllAttributionsCommitID)
	if err != nil {
		return nil, err
	}

	//if the commid id's won't mtch we wiull rebuild
	if startCommitID == branch.LastSyncedAllAttributionsCommitID {
		return &attributions, nil
	}

	// Setting as empty DB Attrib
	if attributions == nil {
		attributions = []models.Attribution{}
	}

	// Make small change here.
	attributionMap := map[string]int{}
	for i, attr := range attributions {
		attributionMap[attr.User.UUID.String()] = i
	}

	//type GitAttribution struct {
	//	UserID      string `json:"-"`
	//	AuthorEmail string `json:"authorEmail"`
	//	Edits       int    `json:"edits"`
	//}

	// Existing Data
	newAttributionUserIDs := []string{}
	for i, log := range gitAttributions {
		userID := strings.Split(log.AuthorEmail, "@")[0]
		gitAttributions[i].UserID = userID
		if _, exists := attributionMap[userID]; !exists {
			newAttributionUserIDs = append(newAttributionUserIDs, userID)
		}
	}

	users, err := user.App.Repo.GetUsersByUUIDs(newAttributionUserIDs)
	if err != nil {
		return nil, err
	}

	userMap := map[string]*models.User{}
	for i := range *users {
		usr := (*users)[i]
		userMap[usr.UUID.String()] = &usr
	}

	for i, _ := range gitAttributions {
		userID := gitAttributions[i].UserID
		if idx, exists := attributionMap[userID]; exists {
			attributions[idx].Edits = attributions[idx].Edits + gitAttributions[i].Edits
		} else {
			attributions = append(attributions, *models.NewAttribution(branch.CanvasRepositoryID, branchID, userMap[userID].ID, gitAttributions[i].Edits))
		}
	}

	err = queries.App.AttributionQuery.ReplaceOrCreateNewAttributions(attributions, branchID)
	if err == nil {
		App.Repo.SetBranchLastSyncedAllAttributionsCommitID(branchID, startCommitID)
	}

	for i, att := range attributions {
		var foundUser *models.User = nil
		for _, user := range *users {
			if att.UserID == user.ID {
				foundUser = &user
				break
			}
		}
		if foundUser != nil {
			attributions[i].User = *foundUser
		}
	}

	return &attributions, nil
}

func (g gitService) FetchCommitFromGit(user *models.User, commitID string, canvasRepoID uint64) ([]*models.Block, error) {

	page, err := queries.App.RepoQuery.GetRepo(map[string]interface{}{"id": canvasRepoID})
	if err != nil {
		return nil, err
	}

	gitUserRepo, err := stores.CreateGitUserRepo(user.UUID.String(), user.Username, page.Studio.UUID.String(), page.UUID.String())
	if err != nil {
		return nil, err
	}

	gitBlocks, err := gitUserRepo.FetchCommitFromGit(commitID)
	if err != nil {
		return nil, err
	}

	return convertGitBlocksToBlocks(gitBlocks, page.ID, canvasRepoID)
}

func (g gitService) FetchBranchFromGit(user *models.User, branchID uint64) ([]*models.Block, error) {
	branch, err := queries.App.BranchQuery.GetBranchByID(branchID)
	if err != nil {
		return nil, err
	}

	if !branch.Committed {
		return nil, errors.New("branch not committed")
	}

	gitUserRepo, err := stores.CreateGitUserRepo(user.UUID.String(), user.Username, branch.CanvasRepository.Studio.UUID.String(), branch.CanvasRepository.UUID.String())
	if err != nil {
		return nil, err
	}

	gitBlocks, err := gitUserRepo.FetchBranchFromGit(branch.Name)
	if err != nil {
		return nil, err
	}

	return convertGitBlocksToBlocks(gitBlocks, branch.CanvasRepositoryID, branch.ID)
}

func convertGitBlocksToBlocks(gitBlocks []*stores.GitBlockV2, canvasRepoID uint64, canvasBranchID uint64) ([]*models.Block, error) {
	blocks := []*models.Block{}

	for _, gitBlock := range gitBlocks {
		var children datatypes.JSON
		data, _ := json.Marshal(gitBlock.Children)
		children.Scan(string(data))

		var attributes datatypes.JSON
		data, _ = json.Marshal(gitBlock.Attributes)
		attributes.Scan(string(data))

		block := models.Block{
			UUID:               uuid.MustParse(gitBlock.UUID),
			Version:            gitBlock.Version,
			Type:               gitBlock.Type,
			Children:           children,
			CanvasRepositoryID: canvasRepoID,
			CanvasBranchID:     &canvasBranchID,
			Rank:               gitBlock.Rank,
			Attributes:         attributes,
			CreatedAt:          gitBlock.CreatedAt,
			UpdatedAt:          gitBlock.UpdatedAt,
		}
		blocks = append(blocks, &block)
	}
	return blocks, nil
}

func (g gitService) CommitMessageBlockToGit(blockID uint64, messageID string) error {

	block, err := queries.App.BlockQuery.GetBlock(map[string]interface{}{"id": blockID})
	if err != nil {
		return err
	}

	msg, err := message.GetMessageByUUID(messageID)
	if err != nil {
		return err
	}

	branch, err := queries.App.BranchQuery.GetBranchByID(*block.CanvasBranchID)
	if err != nil {
		return err
	}

	if branch.Committed {
		return nil
	}

	blocks, err := queries.App.BlockQuery.GetBlocksByBranchID(branch.ID)
	if err != nil {
		return err
	}

	blockViews := []interface{}{}
	blockData := BulkBlocksGetSerializerData(blocks)
	for _, blk := range blockData {
		blockViews = append(blockViews, blk)
	}

	gitUserRepo, err := stores.CreateGitUserRepo(msg.Author.UUID.String(), msg.Author.Username, branch.CanvasRepository.Studio.UUID.String(), branch.CanvasRepository.UUID.String())
	if err != nil {
		return err
	}

	err = gitUserRepo.CreateSnapshotForMessageBlock(&blockViews, branch.Name, block.UUID.String())
	if err != nil {
		return err
	}

	message.MarkMessageAsUsed(msg.ID, true)
	return nil
}
