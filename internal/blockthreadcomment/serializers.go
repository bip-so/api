package blockThreadCommentcomment

import (
	"encoding/json"
	"gitlab.com/phonepost/bip-be-platform/internal/shared"
	"strconv"
	"time"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gorm.io/datatypes"
)

type DefaultSerializer struct {
	ID              uint64                          `json:"id"`
	UUID            string                          `json:"uuid"`
	ThreadID        uint64                          `json:"threadId"`
	ParentID        *uint64                         `json:"parentId"`
	Position        uint                            `json:"position"`
	Data            datatypes.JSON                  `json:"data"`
	ReactionCounter []ReactedCount                  `json:"reactions"`
	IsEdited        bool                            `json:"isEdited"`
	IsReply         bool                            `json:"isReply"`
	CreatedByID     uint64                          `json:"createdByID"`
	UpdatedByID     uint64                          `json:"updatedByID"`
	ArchivedByID    uint64                          `json:"archivedByID"`
	IsArchived      bool                            `json:"isArchived"`
	ArchivedAt      time.Time                       `json:"archivedAt"`
	Mentions        *datatypes.JSON                 `json:"mentions"`
	User            shared.CommonUserMiniSerializer `json:"user"`
}

func SerializeDefaultBlockThreadComment(btc *models.BlockComment) *DefaultSerializer {
	view := &DefaultSerializer{
		ID:           btc.ID,
		UUID:         btc.UUID.String(),
		ThreadID:     btc.ThreadID,
		ParentID:     btc.ParentID,
		Position:     btc.Position,
		Data:         btc.Data,
		IsEdited:     btc.IsEdited,
		IsReply:      btc.IsReply,
		CreatedByID:  btc.CreatedByID,
		UpdatedByID:  btc.UpdatedByID,
		ArchivedByID: btc.ArchivedByID,
		IsArchived:   btc.IsArchived,
		ArchivedAt:   btc.ArchivedAt,
		Mentions:     btc.Mentions,
		//User:         user.UserMiniSerializerData(btc.CreatedByUser),
		User: shared.CommonUserMiniSerializer{
			Id:        btc.CreatedByUser.ID,
			UUID:      btc.CreatedByUser.UUID.String(),
			FullName:  btc.CreatedByUser.FullName,
			Username:  btc.CreatedByUser.Username,
			AvatarUrl: btc.CreatedByUser.AvatarUrl,
		},
	}
	return view
}

func SerializeDefaultManyBlockThreadComment(btc *[]models.BlockComment) *[]DefaultSerializer {
	blockCommentsData := &[]DefaultSerializer{}
	for _, btc := range *btc {
		*blockCommentsData = append(*blockCommentsData, DefaultSerializer{
			ID:           btc.ID,
			UUID:         btc.UUID.String(),
			ThreadID:     btc.ThreadID,
			ParentID:     btc.ParentID,
			Position:     btc.Position,
			Data:         btc.Data,
			IsEdited:     btc.IsEdited,
			IsReply:      btc.IsReply,
			CreatedByID:  btc.CreatedByID,
			UpdatedByID:  btc.UpdatedByID,
			ArchivedByID: btc.ArchivedByID,
			IsArchived:   btc.IsArchived,
			ArchivedAt:   btc.ArchivedAt,
			Mentions:     btc.Mentions,
			//User:         user.UserMiniSerializerData(btc.CreatedByUser),
			User: shared.CommonUserMiniSerializer{
				Id:        btc.CreatedByUser.ID,
				UUID:      btc.CreatedByUser.UUID.String(),
				FullName:  btc.CreatedByUser.FullName,
				Username:  btc.CreatedByUser.Username,
				AvatarUrl: btc.CreatedByUser.AvatarUrl,
			},
		})
	}
	return blockCommentsData
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

func SerializeBlockThreadCommentWithReactionForUser(btc *models.BlockComment, reactions []models.BlockCommentReaction, user *models.User) *DefaultSerializer {
	serialized := SerializeDefaultBlockThreadComment(btc)
	reactionCountList, err := mapReactionCountWithReacted(string(btc.Reactions), btc.ID, reactions, user)
	if err != nil {
	}
	serialized.ReactionCounter = reactionCountList
	return serialized
}

//mapReactionCountWithReacted
func mapReactionCountWithReacted(ReactionCountList string, blockCommentID uint64, reactions []models.BlockCommentReaction, user *models.User) (reactedCountList []ReactedCount, err error) {
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
			if react.BlockCommentID == blockCommentID && reactCount.Emoji == react.Emoji {
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

func SerializeDefaultManyBlockThreadCommentWithReaction(blockcomments *[]models.BlockComment, reactions []models.BlockCommentReaction, user *models.User) *[]DefaultSerializer {
	records := &[]DefaultSerializer{}
	for _, record := range *blockcomments {
		*records = append(*records, *SerializeBlockThreadCommentWithReactionForUser(&record, reactions, user))
	}
	return records
}
