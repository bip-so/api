package reel

import (
	"encoding/json"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/studio"
	"strconv"
	"time"

	"gitlab.com/phonepost/bip-be-platform/internal/shared"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"

	"gitlab.com/phonepost/bip-be-platform/internal/user"

	"github.com/google/uuid"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gorm.io/datatypes"
)

type CanvasRepoMini struct {
	ID              uint64  `json:"id"`
	UUID            string  `json:"uuid"`
	CollectionID    uint64  `json:"collectionID"`
	Name            string  `json:"name"`
	DefaultBranchID *uint64 `json:"defaultBranchID"`
	Key             string  `json:"key"`
}

type ReelsSerialData struct {
	ID                 uint64                            `json:"id"`
	UUID               uuid.UUID                         `json:"uuid"`
	StudioID           uint64                            `json:"studioID"`
	CanvasRepositoryID uint64                            `json:"canvasRepositoryID"`
	CanvasRepository   CanvasRepoMini                    `json:"canvasRepository"`
	CanvasBranchID     uint64                            `json:"canvasBranchID"`
	CommentCount       uint                              `json:"commentCount"`
	StartBlockID       uint64                            `json:"startBlockID"`
	StartBlockUUID     uuid.UUID                         `json:"startBlockUUID"`
	SelectedBlocks     datatypes.JSON                    `json:"selectedBlocks"`
	TextRangeStart     uint                              `json:"textRangeStart"`
	TextRangeEnd       uint                              `json:"textRangeEnd"`
	RangeStart         datatypes.JSON                    `json:"rangeStart"`
	RangeEnd           datatypes.JSON                    `json:"rangeEnd"`
	HighlightedText    datatypes.JSON                    `json:"highlightedText"`
	ContextData        datatypes.JSON                    `json:"contextData"`
	ReactionCounter    []ReactedCount                    `json:"reactions"`
	CreatedByID        uint64                            `json:"createdByID"`
	UpdatedByID        uint64                            `json:"updatedByID"`
	ArchivedByID       *uint64                           `json:"archivedByID"`
	IsArchived         bool                              `json:"isArchived"`
	ArchivedAt         time.Time                         `json:"archivedAt"`
	CreatedAt          time.Time                         `json:"createdAt"`
	UpdatedAt          time.Time                         `json:"updatedAt"`
	Mentions           *datatypes.JSON                   `json:"mentions"`
	User               user.UserMiniSerializer           `json:"user"`
	Studio             shared.CommonStudioMiniSerializer `json:"studio"`
	IsUserFollower     bool                              `json:"isUserFollower"`
	IsStudioMember     bool                              `json:"isStudioMember"`
	IsUserStudioAdmin  bool                              `json:"isUserStudioAdmin"`
	ReactionCopy       string                            `json:"reactionCopy"`
}

func SerializeDefaultReel(model *models.Reel, loggedInUser uint64) *ReelsSerialData {
	canvasRepo, _ := App.Repo.GetRepo(map[string]interface{}{"id": model.CanvasRepositoryID})
	view := ReelsSerialData{
		ID:                 model.ID,
		UUID:               model.UUID,
		StudioID:           model.StudioID,
		CanvasRepositoryID: model.CanvasRepositoryID,
		CanvasBranchID:     model.CanvasBranchID,
		CanvasRepository: CanvasRepoMini{
			ID:              canvasRepo.ID,
			UUID:            canvasRepo.UUID.String(),
			CollectionID:    canvasRepo.CollectionID,
			Name:            canvasRepo.Name,
			DefaultBranchID: canvasRepo.DefaultBranchID,
			Key:             canvasRepo.Key,
		},
		CommentCount:    model.CommentCount,
		StartBlockID:    model.StartBlockID,
		StartBlockUUID:  model.StartBlockUUID,
		SelectedBlocks:  model.SelectedBlocks,
		TextRangeStart:  model.TextRangeStart,
		TextRangeEnd:    model.TextRangeEnd,
		RangeStart:      model.RangeStart,
		RangeEnd:        model.RangeEnd,
		HighlightedText: model.HighlightedText,
		ContextData:     model.ContextData,
		CreatedByID:     model.CreatedByID,
		UpdatedByID:     model.UpdatedByID,
		ArchivedByID:    model.ArchivedByID,
		IsArchived:      model.IsArchived,
		ArchivedAt:      model.ArchivedAt,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
		Mentions:        model.Mentions,
		User:            user.UserMiniSerializerData(model.CreatedByUser),
		ReactionCopy:    BuildReelsReactionString(model, loggedInUser),
	}
	if model.Studio != nil {
		view.Studio = StudioMiniSerializerData(model.Studio)
		view.Studio.IsRequested = studio.App.StudioService.CheckIsRequested(loggedInUser, view.StudioID)
	}

	return &view
}

func SerializeReelWithReactionForUser(reel *models.Reel, reactions []models.ReelReaction, user *models.User) *ReelsSerialData {
	var userID uint64
	if user != nil {
		userID = user.ID
	}
	serialized := SerializeDefaultReel(reel, userID)
	reactionCountList, err := mapReelReactionCountWithReacted(string(reel.Reactions), reel.ID, reactions, user)
	if err != nil {
	}
	serialized.ReactionCounter = reactionCountList
	return serialized
}

//func SerializeDefaultManyReels(modelInstances *[]models.Reel) *[]ReelsSerialData {
//	reelData := &[]ReelsSerialData{}
//
//	for _, model := range *modelInstances {
//		*reelData = append(*reelData, *SerializeDefaultReel(&model))
//	}
//	return reelData
//}

func SerializeDefaultManyReelsWithReactionsForUser(modelInstances *[]models.Reel, reactions []models.ReelReaction, user *models.User, members []models.Member, userFollowings *[]models.FollowUser) *[]ReelsSerialData {
	reelData := &[]ReelsSerialData{}

	for _, model := range *modelInstances {
		reel := *SerializeReelWithReactionForUser(&model, reactions, user)
		if user != nil {
			reel.IsStudioMember = checkIsStudioMember(reel.StudioID, user.ID, members)
			if reel.CreatedByID == user.ID {
				reel.IsUserFollower = true
			} else {
				reel.IsUserFollower = checkIsUserFollowing(reel.CreatedByID, userFollowings)
			}
			reel.IsUserStudioAdmin = shared.IsUserStudioAdmin(user.ID, reel.StudioID)
		}
		*reelData = append(*reelData, reel)
	}
	return reelData
}

func SerializeDefaultSingleReelsWithReactionsForUser(modelInstance *models.Reel, reactions []models.ReelReaction, user *models.User, members []models.Member, userFollowings *[]models.FollowUser) *ReelsSerialData {
	reelData := &ReelsSerialData{}
	fmt.Println(modelInstance)
	reel := *SerializeReelWithReactionForUser(modelInstance, reactions, user)
	if user != nil {
		reel.IsStudioMember = checkIsStudioMember(reel.StudioID, user.ID, members)
		if reel.CreatedByID == user.ID {
			reel.IsUserFollower = true
		} else {
			reel.IsUserFollower = checkIsUserFollowing(reel.CreatedByID, userFollowings)
		}
	}
	*reelData = reel
	return reelData
}

func SerializePopularReels(modelInstances *[]models.Reel) *[]ReelsSerialData {
	reelData := &[]ReelsSerialData{}
	for _, model := range *modelInstances {
		ReelData := &ReelsSerialData{
			ID:                 model.ID,
			UUID:               model.UUID,
			StudioID:           model.StudioID,
			CanvasRepositoryID: model.CanvasRepositoryID,
			CanvasBranchID:     model.CanvasBranchID,
			CommentCount:       model.CommentCount,
			StartBlockID:       model.StartBlockID,
			StartBlockUUID:     model.StartBlockUUID,
			SelectedBlocks:     model.SelectedBlocks,
			TextRangeStart:     model.TextRangeStart,
			TextRangeEnd:       model.TextRangeEnd,
			RangeStart:         model.RangeStart,
			RangeEnd:           model.RangeEnd,
			HighlightedText:    model.HighlightedText,
			ContextData:        model.ContextData,
			CreatedByID:        model.CreatedByID,
			UpdatedByID:        model.UpdatedByID,
			ArchivedByID:       model.ArchivedByID,
			IsArchived:         model.IsArchived,
			ArchivedAt:         model.ArchivedAt,
			CreatedAt:          model.CreatedAt,
			UpdatedAt:          model.UpdatedAt,
			Mentions:           model.Mentions,
		}
		if model.CreatedByUser != nil {
			ReelData.User = user.UserMiniSerializerData(model.CreatedByUser)
		}
		if model.Studio != nil {
			ReelData.Studio = StudioMiniSerializerData(model.Studio)
		}
		*reelData = append(*reelData, *ReelData)
	}
	return reelData
}

//////////////// REELS COMMENTS ////////////////////////////////////
type ReelsCommentsSerialData struct {
	ID              uint64                  `json:"id"`
	UUID            uuid.UUID               `json:"uuid"`
	ReelID          uint64                  `json:"reelID"`
	ParentID        *uint64                 `json:"parentID"`
	Position        uint                    `json:"position"`
	Data            datatypes.JSON          `json:"rangeStart"`
	IsEdited        bool                    `json:"isEdited"`
	IsReply         bool                    `json:"isReply"`
	ReactionCounter []ReactedCount          `json:"reactionCounter"`
	UpdatedByID     uint64                  `json:"updatedByID"`
	ArchivedByID    *uint64                 `json:"archivedByID"`
	IsArchived      bool                    `json:"isArchived"`
	ArchivedAt      time.Time               `json:"archivedAt"`
	CreatedAt       time.Time               `json:"createdAt"`
	UpdatedAt       time.Time               `json:"updatedAt"`
	Mentions        *datatypes.JSON         `json:"mentions"`
	User            user.UserMiniSerializer `json:"user"`
	CommentCount    uint                    `json:"commentCount"`
}

func SerializeDefaultReelComment(model *models.ReelComment) *ReelsCommentsSerialData {
	view := ReelsCommentsSerialData{
		ID:           model.ID,
		UUID:         model.UUID,
		ReelID:       model.ReelID,
		ParentID:     model.ParentID,
		Data:         model.Data,
		IsReply:      model.IsReply,
		IsEdited:     model.IsEdited,
		UpdatedByID:  model.UpdatedByID,
		ArchivedByID: &model.ArchivedByID,
		IsArchived:   model.IsArchived,
		ArchivedAt:   model.ArchivedAt,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
		Mentions:     model.Mentions,
		User:         user.UserMiniSerializerData(model.CreatedByUser),
		CommentCount: model.CommentCount,
	}

	return &view
}

func SerializeDefaultManyReelComments(modelInstances *[]models.ReelComment) *[]ReelsCommentsSerialData {
	reelCommentData := &[]ReelsCommentsSerialData{}
	for _, model := range *modelInstances {
		*reelCommentData = append(*reelCommentData, *SerializeDefaultReelComment(&model))
	}
	return reelCommentData
}

//mapReactionCountWithReacted
func mapReelReactionCountWithReacted(ReactionCountList string, reelID uint64, reactions []models.ReelReaction, user *models.User) (reactedCountList []ReactedCount, err error) {
	if reactions == nil {
		return nil, err
	}
	reactionCountList := []models.ReactionCounter{}
	err = json.Unmarshal([]byte(ReactionCountList), &reactionCountList)
	if err != nil {
		return nil, err
	}
	for _, reactCount := range reactionCountList {
		found := false
		for _, react := range reactions {
			if react.ReelID == reelID && reactCount.Emoji == react.Emoji {
				found = true
				break
			}
		}
		cnt, _ := strconv.Atoi(reactCount.Count)
		if user == nil {
			reactedCountList = append(reactedCountList, ReactedCount{Emoji: reactCount.Emoji, Count: cnt, Reacted: nil})
		} else {
			reactedCountList = append(reactedCountList, ReactedCount{Emoji: reactCount.Emoji, Count: cnt, Reacted: &found})
		}
	}
	return reactedCountList, err
}

type ReactedCount struct {
	Emoji   string `json:"emoji"`
	Count   int    `json:"count"`
	Reacted *bool  `json:"reacted"`
}

type ReactCountListView struct {
	Emoji string `json:"emoji"`
	Count int    `json:"count"`
}

func SerializeReelComentWithReactionForUser(reel *models.ReelComment, reactions []models.ReelCommentReaction, user *models.User) *ReelsCommentsSerialData {
	serialized := SerializeDefaultReelComment(reel)

	reactionCountList, err := mapReelCommentReactionCountWithReacted(string(reel.Reactions), reel.ID, reactions, user)
	if err != nil {
	}
	serialized.ReactionCounter = reactionCountList
	return serialized
}

//mapReactionCountWithReacted
func mapReelCommentReactionCountWithReacted(ReactionCountList string, reelCommentID uint64, reactions []models.ReelCommentReaction, user *models.User) (reactedCountList []ReactedCount, err error) {
	if reactions == nil {
		return nil, err
	}
	reactionCountList := []models.ReactionCounter{}
	err = json.Unmarshal([]byte(ReactionCountList), &reactionCountList)
	if err != nil {
		return nil, err
	}
	for _, reactCount := range reactionCountList {

		var found bool
		for _, react := range reactions {
			if react.ReelCommentID == reelCommentID && reactCount.Emoji == react.Emoji {
				found = true
				break
			}
		}
		cnt, _ := strconv.Atoi(reactCount.Count)

		if user == nil {
			reactedCountList = append(reactedCountList, ReactedCount{Emoji: reactCount.Emoji, Count: cnt, Reacted: nil})
		} else {
			reactedCountList = append(reactedCountList, ReactedCount{Emoji: reactCount.Emoji, Count: cnt, Reacted: &found})
		}
	}
	return reactedCountList, err
}

func SerializeDefaultManyReelCommentsWithReactionsForUser(modelInstances *[]models.ReelComment, reactions []models.ReelCommentReaction, user *models.User) *[]ReelsCommentsSerialData {
	reelData := &[]ReelsCommentsSerialData{}

	if len(*modelInstances) == 0 {
		return reelData
	}
	for _, model := range *modelInstances {
		*reelData = append(*reelData, *SerializeReelComentWithReactionForUser(&model, reactions, user))
	}
	return reelData
}

type StudioMiniSerializer struct {
	ID          uint64 `json:"id"`
	UUID        string `json:"uuid"`
	DisplayName string `json:"displayName"`
	Handle      string `json:"handle"`
	ImageURL    string `json:"imageUrl"`
	CreatedByID uint64 `json:"createdById"`
}

func StudioMiniSerializerData(studio *models.Studio) shared.CommonStudioMiniSerializer {
	//fmt.Println(studio)
	return shared.CommonStudioMiniSerializer{
		ID:                    studio.ID,
		UUID:                  studio.UUID.String(),
		DisplayName:           studio.DisplayName,
		Handle:                studio.Handle,
		ImageURL:              studio.ImageURL,
		CreatedByID:           studio.CreatedByID,
		AllowPublicMembership: studio.AllowPublicMembership,
	}
}

func checkIsStudioMember(studioID uint64, userID uint64, members []models.Member) bool {
	for _, member := range members {
		if member.UserID == userID && member.StudioID == studioID {
			return true
		}
	}
	return false
}

func checkIsUserFollowing(userID uint64, userFollowings *[]models.FollowUser) bool {
	if userFollowings == nil {
		return false
	}
	for _, userFollowing := range *userFollowings {
		if userFollowing.UserId == userID {
			return true
		}
	}
	return false
}

type ReelContextData struct {
	Text string `json:"text"`
}

type ReelSelectedBlocks struct {
	BlockUUIDs []string    `json:"blockUUIDs"`
	BlocksData interface{} `json:"blocksData"`
}

type ReelDocument struct {
	ID          uint64    `json:"id"`
	UUID        string    `json:"uuid"`
	CreatedByID uint64    `json:"createdById"`
	Text        string    `json:"text"`
	BlockUUIDs  []string  `json:"blockUUIDs"`
	CreatedAt   time.Time `json:"createdAt"`
}

func BuildReelsReactionString(model *models.Reel, loggedInUserId uint64) string {
	count := shared.CountResults(map[string]interface{}{"reel_id": model.ID}, &models.ReelReaction{})
	if count == 0 {
		return ""
	}
	pr := models.ReelReaction{}
	if count == 1 {
		_ = postgres.GetDB().Model(&models.ReelReaction{}).
			Where("reel_id = ?", model.ID).
			Preload("CreatedByUser").
			Limit(1).
			First(&pr)
		return pr.CreatedByUser.FullName
	}

	pr2 := models.ReelReaction{}
	// This query returns Random
	_ = postgres.GetDB().Model(&models.ReelReaction{}).
		Where("reel_id = ? and created_by_id != ?", model.ID, loggedInUserId).
		Order("created_at desc").
		Preload("CreatedByUser").
		Limit(1).
		First(&pr2)
	var name string
	if pr2.ID == 0 {
		name = ""
	} else {
		name = pr2.CreatedByUser.FullName
	}
	//countForString := count - 1
	//return fmt.Sprintf("%s and %d others reacted to this post", name, countForString)
	return fmt.Sprintf("%s", name)
	// and 62 others reacted to this post.
}
