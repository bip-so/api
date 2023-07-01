package blockthread

import (
	"encoding/json"
	"gitlab.com/phonepost/bip-be-platform/internal/shared"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

type DefaultSerializer struct {
	ID                 uint64                          `json:"id"`
	UUID               string                          `json:"uuid"`
	CanvasRepositoryID uint64                          `json:"canvasRepositoryId"`
	CanvasBranchID     uint64                          `json:"canvasBranchId"`
	StartBlockID       uint64                          `json:"startBlockId"`
	StartBlockUUID     uuid.UUID                       `json:"startBlockUUID"`
	Position           uint                            `json:"position"`
	TextRangeStart     uint                            `json:"textRangeStart"`
	TextRangeEnd       uint                            `json:"textRangeEnd"`
	Text               string                          `json:"text"`
	HighlightedText    string                          `json:"highlightedText"`
	ReactionCounter    []ReactedCount                  `json:"reactions"`
	CommentCount       uint                            `json:"commentCount"`
	IsResolved         bool                            `json:"isResolved"`
	ResolvedByID       *uint64                         `json:"resolvedById"`
	ResolvedAt         time.Time                       `json:"resolvedAt"`
	CreatedByID        uint64                          `json:"createdById"`
	UpdatedByID        uint64                          `json:"updatedById"`
	ArchivedByID       uint64                          `json:"archivedById"`
	IsArchived         bool                            `json:"isArchived"`
	ArchivedAt         time.Time                       `json:"archivedAt"`
	CreatedAt          time.Time                       `json:"createdAt"`
	UpdatedAt          time.Time                       `json:"updatedAt"`
	Mentions           *datatypes.JSON                 `json:"mentions"`
	User               shared.CommonUserMiniSerializer `json:"user"`
}

func SerializeDefaultBlockThread(bt *models.BlockThread) *DefaultSerializer {
	view := &DefaultSerializer{
		ID:                 bt.ID,
		UUID:               bt.UUID.String(),
		CanvasRepositoryID: bt.CanvasRepositoryID,
		CanvasBranchID:     bt.CanvasBranchID,
		StartBlockID:       bt.StartBlockID,
		StartBlockUUID:     bt.StartBlockUUID,
		Position:           bt.Position,
		TextRangeStart:     bt.TextRangeStart,
		TextRangeEnd:       bt.TextRangeEnd,
		Text:               bt.Text,
		HighlightedText:    bt.HighlightedText,
		CommentCount:       bt.CommentCount,
		IsResolved:         bt.IsResolved,
		ResolvedByID:       bt.ResolvedByID,
		ResolvedAt:         bt.ResolvedAt,
		CreatedByID:        bt.CreatedByID,
		UpdatedByID:        bt.UpdatedByID,
		ArchivedByID:       bt.ArchivedByID,
		IsArchived:         bt.IsArchived,
		ArchivedAt:         bt.ArchivedAt,
		CreatedAt:          bt.CreatedAt,
		UpdatedAt:          bt.UpdatedAt,
		Mentions:           bt.Mentions,
		//User:               user2.UserMiniSerializerData(bt.CreatedByUser),
		User: shared.CommonUserMiniSerializer{
			Id:        bt.CreatedByUser.ID,
			UUID:      bt.CreatedByUser.UUID.String(),
			FullName:  bt.CreatedByUser.FullName,
			Username:  bt.CreatedByUser.Username,
			AvatarUrl: bt.CreatedByUser.AvatarUrl,
		},
	}
	return view
}

// CC: Template this
func SerializeDefaultManyBlockThread(models *[]models.BlockThread) *[]DefaultSerializer {
	records := &[]DefaultSerializer{}
	for _, record := range *models {
		*records = append(*records, DefaultSerializer{
			ID:                 record.ID,
			UUID:               record.UUID.String(),
			CanvasRepositoryID: record.CanvasRepositoryID,
			CanvasBranchID:     record.CanvasBranchID,
			StartBlockID:       record.StartBlockID,
			Position:           record.Position,
			TextRangeStart:     record.TextRangeStart,
			TextRangeEnd:       record.TextRangeEnd,
			Text:               record.Text,
			HighlightedText:    record.HighlightedText,
			CommentCount:       record.CommentCount,
			IsResolved:         record.IsResolved,
			ResolvedByID:       record.ResolvedByID,
			ResolvedAt:         record.ResolvedAt,
			CreatedByID:        record.CreatedByID,
			UpdatedByID:        record.UpdatedByID,
			ArchivedByID:       record.ArchivedByID,
			IsArchived:         record.IsArchived,
			ArchivedAt:         record.ArchivedAt,
			CreatedAt:          record.CreatedAt,
			UpdatedAt:          record.UpdatedAt,
			Mentions:           record.Mentions,
		})
	}
	return records
}

func SerializeDefaultManyBlockThreadWithReaction(blockthreads *[]models.BlockThread, reactions []models.BlockThreadReaction, user *models.User) *[]DefaultSerializer {
	records := &[]DefaultSerializer{}
	for _, record := range *blockthreads {
		*records = append(*records, *SerializeBlockThreadWithReactionForUser(&record, reactions, user))
	}
	return records
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

func SerializeBlockThreadWithReactionForUser(bt *models.BlockThread, reactions []models.BlockThreadReaction, user *models.User) *DefaultSerializer {
	serialized := SerializeDefaultBlockThread(bt)
	reactionCountList, err := mapReactionCountWithReacted(string(bt.Reactions), bt.ID, reactions, user)
	if err != nil {
	}
	serialized.ReactionCounter = reactionCountList
	return serialized
}

//mapReactionCountWithReacted
func mapReactionCountWithReacted(ReactionCountList string, blockThreadID uint64, reactions []models.BlockThreadReaction, user *models.User) (reactedCountList []ReactedCount, err error) {
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
			if react.BlockThreadID == blockThreadID && reactCount.Emoji == react.Emoji {
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
