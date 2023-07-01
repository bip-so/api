package reel

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/phonepost/bip-be-platform/internal/feed"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/internal/reactions"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/search"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

func (s reelService) GetAll(studioID uint64, canvasBranchID uint64, uuid1 uuid.UUID, authUser *models.User) (*[]models.Reel, error) {
	//var uint64Pointer *uint64
	//var uuidPointer *uuid.UUID
	var query map[string]interface{}
	var reels *[]models.Reel
	var err error
	if canvasBranchID != 0 && authUser != nil {
		branch, _ := queries.App.BranchQuery.GetBranchByID(canvasBranchID)
		if branch.PublicAccess == models.PRIVATE {
			hasPerm, err := permissions.App.Service.CanUserDoThisOnBranch(authUser.ID, canvasBranchID, permissiongroup.CANVAS_BRANCH_VIEW)
			fmt.Println("has perm check", hasPerm, err)
			if err != nil || hasPerm == false {
				return reels, errors.New("User doesn't have permission")
			}
		}
		query = map[string]interface{}{"studio_id": studioID, "canvas_branch_id": canvasBranchID, "is_archived": false}
	} else if canvasBranchID != 0 && authUser == nil {
		branch, _ := queries.App.BranchQuery.GetBranchByID(canvasBranchID)
		if branch.PublicAccess == models.PRIVATE {
			return reels, errors.New("User doesn't have permission")
		}
		query = map[string]interface{}{"studio_id": studioID, "canvas_branch_id": canvasBranchID, "is_archived": false}
	} else if uuid1 != uuid.Nil {
		query = map[string]interface{}{"studio_id": studioID, "start_block_uuid": uuid1, "is_archived": false}
	} else {
		query = map[string]interface{}{"studio_id": studioID, "is_archived": false}
	}

	reels, err = App.Repo.GetReels(query)
	if err != nil {
		return nil, nil
	}
	return reels, nil
}

func (s reelService) GetOne(studioID uint64, reelID uint64) (*models.Reel, error) {
	var query map[string]interface{}
	query = map[string]interface{}{"id": reelID, "is_archived": false}
	reel, err := App.Repo.GetReel(query)
	if err != nil {
		return nil, nil
	}
	return reel, nil
}

func (s reelService) GetAnonymousPopular(skip, limit int) (*[]ReelsSerialData, error) {
	reels, err := App.Repo.GetPopular(skip, limit)
	if err != nil {
		return nil, nil
	}
	var reelIDs []uint64
	for _, reel := range *reels {
		reelIDs = append(reelIDs, reel.ID)
	}
	reelReactions, _ := reactions.App.Repo.GetReelReactionByIDs(reelIDs)
	reelsData := SerializeDefaultManyReelsWithReactionsForUser(reels, reelReactions, nil, nil, nil)
	return reelsData, nil
}

func (s reelService) GetAnonymousStudioPopular(skip, limit int, studioID uint64) (*[]ReelsSerialData, error) {
	reels, err := App.Repo.GetStudioPopular(skip, limit, studioID)
	if err != nil {
		return nil, nil
	}
	var reelIDs []uint64
	for _, reel := range *reels {
		reelIDs = append(reelIDs, reel.ID)
	}
	reelReactions, _ := reactions.App.Repo.GetReelReactionByIDs(reelIDs)
	reelsData := SerializeDefaultManyReelsWithReactionsForUser(reels, reelReactions, nil, nil, nil)
	return reelsData, nil
}

//func (s reelService) Create(studioID uint64) (*[]models.Reel, error) {
//	reels, err := App.Repo.GetAll(studioID, false, 0)
//	if err != nil {
//		return nil, nil
//	}
//	return reels, nil
//}

func (s reelService) GetAllComments(studioID uint64, reelID uint64) (*[]models.ReelComment, error) {
	reelcomments, err := App.Repo.GetAllReelComments(studioID, reelID)
	if err != nil {
		return nil, nil
	}
	return reelcomments, nil
}

func (s reelService) GetChildComments(parentReelCommentID uint64) (*[]models.ReelComment, error) {
	reelcomments, err := App.Repo.GetChildReelComments(parentReelCommentID)
	if err != nil {
		return nil, nil
	}
	return reelcomments, nil
}

func (s reelService) GetLoggedInPopular(authUser *models.User, skip int, limit int) (*[]ReelsSerialData, error) {
	reels, err := App.Repo.GetPopular(skip, limit)
	if err != nil {
		return nil, nil
	}
	reelsData, err := s.GetReelsWithConfigData(reels, authUser)
	if err != nil {
		return nil, nil
	}
	return reelsData, nil
}

func (s reelService) GetReelsWithConfigData(reels *[]models.Reel, authUser *models.User) (*[]ReelsSerialData, error) {
	var reelIDs []uint64
	for _, reel := range *reels {
		reelIDs = append(reelIDs, reel.ID)
	}
	reelReactions, _ := reactions.App.Repo.GetUserReelReactionByIDs(reelIDs, authUser.ID)
	members, _ := App.Repo.GetMembersByUserID(authUser.ID)
	userFollowings, _ := App.Repo.GetUserFollowings(authUser.ID)
	reelsData := SerializeDefaultManyReelsWithReactionsForUser(reels, reelReactions, authUser, members, userFollowings)
	return reelsData, nil
}

func (s reelService) AddReelToAlgolia(reelID uint64) error {
	reel, err := App.Repo.GetReel(map[string]interface{}{"id": reelID})
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if reel.IsArchived {
		return nil
	}
	contextData := ReelContextData{}
	selectedBlocks := ReelSelectedBlocks{}
	json.Unmarshal(reel.ContextData, &contextData)
	json.Unmarshal(reel.SelectedBlocks, &selectedBlocks)
	reelDoc := ReelDocument{
		ID:          reel.ID,
		UUID:        reel.UUID.String(),
		Text:        contextData.Text,
		BlockUUIDs:  selectedBlocks.BlockUUIDs,
		CreatedByID: reel.CreatedByID,
		CreatedAt:   reel.CreatedAt,
	}
	err = search.GetIndex(search.ReelDocumentIndexName).SaveRecord(reelDoc)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

func (s reelService) DeleteReelFromAlgolia(reelID uint64) error {
	err := search.GetIndex(search.ReelDocumentIndexName).DeleteRecordByID(reelID)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return err
}

func (s reelService) StudioReels(user *models.User, studioID uint64, offset int, limit int, reels []models.Reel, skipFeed int) ([]models.Reel, error) {
	reelsActivityFeed, err := feed.App.Service.GetReelActivities(user.ID, skipFeed, 100)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	var reelIDs []uint64
	for _, reelActivity := range reelsActivityFeed.Results {
		if reelActivity.ForeignID != "" {
			reelIDs = append(reelIDs, utils.Uint64(reelActivity.ForeignID))
		}
	}
	newReels, err := App.Repo.GetReelsByIDsForStudio(reelIDs, studioID, offset, limit)
	if err != nil {
		return nil, err
	}
	reels = append(reels, newReels...)
	if len(reels) < 15 && len(reelsActivityFeed.Results) > 0 {
		reels, _ = s.StudioReels(user, studioID, offset, limit, reels, skipFeed+limit)
	}
	return reels, nil
}

func (s reelService) StudioUserReelsFromDb(user *models.User, studioID uint64, offset int, limit int) ([]models.Reel, error) {
	userAccessReels := []models.Reel{}
	reels, err := App.Repo.GetReelsPopulatedData(map[string]interface{}{"studio_id": studioID, "is_archived": false})
	if err != nil {
		return nil, err
	}
	for _, reel := range reels {
		if len(userAccessReels) >= offset+limit {
			break
		}
		if reel.CanvasBranch != nil && reel.CanvasBranch.PublicAccess != models.PRIVATE {
			userAccessReels = append(userAccessReels, reel)
		} else {
			hasPerm, err := permissions.App.Service.CanUserDoThisOnBranch(user.ID, reel.CanvasBranchID, permissiongroup.CANVAS_BRANCH_VIEW)
			if err == nil && hasPerm {
				userAccessReels = append(userAccessReels, reel)
			}
		}
	}
	fmt.Println("length of reels", len(userAccessReels))
	return userAccessReels[offset:], nil
}
