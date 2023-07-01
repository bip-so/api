package canvasbranch

import (
	"encoding/json"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/internal/shared"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gorm.io/datatypes"
)

type BlockBulkResponseMetadata struct {
	CreatedBlockIDs         []uint64 `json:"createdBlockIDs"`
	UpdatedBlockIDs         []uint64 `json:"updatedBlockIDs"`
	FailedToUpdatedBlockIDs []uint64 `json:"failedToUpdatedBlockIDs"`
	FailedToDeleteBlockIDs  []uint64 `json:"failedToDeleteBlockIDs"`
	DeletedBlocksCount      int      `json:"deletedBlocksCount"`
	FailedBlocksCount       int      `json:"failedBlocksCount"`
}

type BulkBlocks struct {
	ID              uint64                          `json:"id"`
	UUID            uuid.UUID                       `json:"uuid"`
	Version         uint                            `json:"version"`
	Type            string                          `json:"type"`
	Position        uint                            `json:"position"`
	Children        datatypes.JSON                  `json:"children"`
	CreatedAt       time.Time                       `json:"createdAt"`
	UpdatedAt       time.Time                       `json:"updatedAt"`
	CreatedByID     uint64                          `json:"createdByID"`
	UpdatedByID     uint64                          `json:"updatedByID"`
	ArchivedByID    *uint64                         `json:"archivedByID"`
	IsArchived      bool                            `json:"isArchived"`
	ArchivedAt      time.Time                       `json:"archivedAt"`
	Attributes      datatypes.JSON                  `json:"attributes"`
	ReactionCounter []ReactedCount                  `json:"reactions"`
	Mentions        *datatypes.JSON                 `json:"mentions"`
	CommentCount    uint                            `json:"commentCount"`
	ReelCount       uint                            `json:"reelCount"`
	CreatedByUser   shared.CommonUserMiniSerializer `json:"createdByUser"`
	UpdatedByUser   shared.CommonUserMiniSerializer `json:"updatedByUser"`
	Rank            int32                           `json:"rank"`
	Contributors    datatypes.JSON                  `json:"contributors"`
}

type BlockBulkDiffResponse struct {
	Data []BulkBlocks `json:"data"`
}

type BlockBulkResponse struct {
	ResponseMetadata BlockBulkResponseMetadata `json:"responseMetadata,omitempty"`
	Data             []BulkBlocks              `json:"data"`
}
type BlockBulkRoughResponse struct {
	RoughBranchID    uint64                    `json:"roughBranchID"`
	ResponseMetadata BlockBulkResponseMetadata `json:"responseMetadata,omitempty"`
	Data             []BulkBlocks              `json:"data"`
}

type BlockBulkResponseNoMeta struct {
	Data []BulkBlocks `json:"data"`
}

func BulkBlocksSerializerData(blocks *[]models.Block, meta BlockBulkResponseMetadata) BlockBulkResponse {
	var blocksData []BulkBlocks

	for _, block := range *blocks {
		// Commented
		//cu, _ := user2.App.Repo.GetUser(map[string]interface{}{"id": block.CreatedByID})
		//uu, _ := user2.App.Repo.GetUser(map[string]interface{}{"id": block.UpdatedByID})

		cu, _ := queries.App.UserQueries.GetUserByID(block.CreatedByID)
		uu, _ := queries.App.UserQueries.GetUserByID(block.UpdatedByID)

		blocksData = append(blocksData, BulkBlocks{
			ID:           block.ID,
			UUID:         block.UUID,
			Version:      block.Version,
			Type:         block.Type,
			Children:     block.Children,
			CreatedAt:    block.CreatedAt,
			UpdatedAt:    block.UpdatedAt,
			CreatedByID:  block.CreatedByID,
			UpdatedByID:  block.UpdatedByID,
			IsArchived:   block.IsArchived,
			Attributes:   block.Attributes,
			Mentions:     block.Mentions,
			CommentCount: block.CommentCount,
			ReelCount:    block.ReelCount,
			Rank:         block.Rank,
			Contributors: block.Contributors,
			CreatedByUser: shared.CommonUserMiniSerializer{
				Id:        cu.ID,
				UUID:      cu.UUID.String(),
				FullName:  cu.FullName,
				Username:  cu.Username,
				AvatarUrl: cu.AvatarUrl,
			},
			UpdatedByUser: shared.CommonUserMiniSerializer{
				Id:        uu.ID,
				UUID:      uu.UUID.String(),
				FullName:  uu.FullName,
				Username:  uu.Username,
				AvatarUrl: uu.AvatarUrl,
			},
		})
	}
	metadata := BlockBulkResponseMetadata{
		CreatedBlockIDs:         meta.CreatedBlockIDs,
		UpdatedBlockIDs:         meta.UpdatedBlockIDs,
		FailedToUpdatedBlockIDs: meta.FailedToUpdatedBlockIDs,
		FailedToDeleteBlockIDs:  meta.FailedToDeleteBlockIDs,
		DeletedBlocksCount:      meta.DeletedBlocksCount,
		FailedBlocksCount:       meta.FailedBlocksCount,
	}

	Reponse := BlockBulkResponse{}
	Reponse.Data = blocksData
	Reponse.ResponseMetadata = metadata
	return Reponse

}

func BulkBlocksGetSerializerData(blocks *[]models.Block) []BulkBlocks {
	var blocksData []BulkBlocks

	for _, block := range *blocks {
		//cu, _ := user2.App.Repo.GetUser(map[string]interface{}{"id": block.CreatedByID})
		//uu, _ := user2.App.Repo.GetUser(map[string]interface{}{"id": block.UpdatedByID})

		cu, _ := queries.App.UserQueries.GetUserByID(block.CreatedByID)
		uu, _ := queries.App.UserQueries.GetUserByID(block.UpdatedByID)

		blocksData = append(blocksData, BulkBlocks{
			ID:           block.ID,
			UUID:         block.UUID,
			Version:      block.Version,
			Type:         block.Type,
			Children:     block.Children,
			CreatedAt:    block.CreatedAt,
			UpdatedAt:    block.UpdatedAt,
			CreatedByID:  block.CreatedByID,
			UpdatedByID:  block.UpdatedByID,
			Attributes:   block.Attributes,
			Mentions:     block.Mentions,
			CommentCount: block.CommentCount,
			ReelCount:    block.ReelCount,
			Rank:         block.Rank,
			Contributors: block.Contributors,
			CreatedByUser: shared.CommonUserMiniSerializer{
				Id:        cu.ID,
				UUID:      cu.UUID.String(),
				FullName:  cu.FullName,
				Username:  cu.Username,
				AvatarUrl: cu.AvatarUrl,
			},
			UpdatedByUser: shared.CommonUserMiniSerializer{
				Id:        uu.ID,
				UUID:      uu.UUID.String(),
				FullName:  uu.FullName,
				Username:  uu.Username,
				AvatarUrl: uu.AvatarUrl,
			},
		})
	}
	return blocksData

}

func BulkBlocksReactionsSerializerData(blocks *[]models.Block, blockReactions *[]models.BlockReaction, user *models.User) []BulkBlocks {
	var blocksData []BulkBlocks

	for _, block := range *blocks {
		// cu, _ := user2.App.Repo.GetUser(map[string]interface{}{"id": block.CreatedByID})
		// uu, _ := user2.App.Repo.GetUser(map[string]interface{}{"id": block.UpdatedByID})
		blockData := BulkBlocks{
			ID:           block.ID,
			UUID:         block.UUID,
			Version:      block.Version,
			Type:         block.Type,
			Children:     block.Children,
			CreatedAt:    block.CreatedAt,
			UpdatedAt:    block.UpdatedAt,
			CreatedByID:  block.CreatedByID,
			UpdatedByID:  block.UpdatedByID,
			Attributes:   block.Attributes,
			Mentions:     block.Mentions,
			CommentCount: block.CommentCount,
			ReelCount:    block.ReelCount,
			Rank:         block.Rank,
			Contributors: block.Contributors,
			CreatedByUser: shared.CommonUserMiniSerializer{
				Id:        block.CreatedByUser.ID,
				UUID:      block.CreatedByUser.UUID.String(),
				FullName:  block.CreatedByUser.FullName,
				Username:  block.CreatedByUser.Username,
				AvatarUrl: block.CreatedByUser.AvatarUrl,
			},
			UpdatedByUser: shared.CommonUserMiniSerializer{
				Id:        block.UpdatedByUser.ID,
				UUID:      block.UpdatedByUser.UUID.String(),
				FullName:  block.UpdatedByUser.FullName,
				Username:  block.UpdatedByUser.Username,
				AvatarUrl: block.UpdatedByUser.AvatarUrl,
			},
		}
		if blockReactions != nil {
			reactionCountList, _ := mapReactionCountWithReacted(string(block.Reactions), block.ID, *blockReactions, user)
			blockData.ReactionCounter = reactionCountList
		}
		blocksData = append(blocksData, blockData)
	}
	return blocksData

}

func mapReactionCountWithReacted(ReactionCountList string, blockID uint64, reactions []models.BlockReaction, user *models.User) (reactedCountList []ReactedCount, err error) {
	reactionCountList := []models.ReactionCounter{}
	err = json.Unmarshal([]byte(ReactionCountList), &reactionCountList)
	if err != nil {
		return nil, err
	}
	for _, reactCount := range reactionCountList {
		found := false
		for _, react := range reactions {
			if react.BlockID == blockID && reactCount.Emoji == react.Emoji {
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

//
//func BulkRoughBlocksSerializerData(blocks *[]models.Block, meta BlockBulkResponseMetadata, roughBranchID uint64) BlockBulkRoughResponse {
//	var blocksData []BulkBlocks
//
//	for _, block := range *blocks {
//		cu, _ := user2.App.Repo.GetUser(map[string]interface{}{"id": block.CreatedByID})
//		uu, _ := user2.App.Repo.GetUser(map[string]interface{}{"id": block.UpdatedByID})
//		blocksData = append(blocksData, BulkBlocks{
//			ID:           block.ID,
//			UUID:         block.UUID,
//			Version:      block.Version,
//			Type:         block.Type,
//			Children:     block.Children,
//			CreatedAt:    block.CreatedAt,
//			UpdatedAt:    block.UpdatedAt,
//			CreatedByID:  block.CreatedByID,
//			UpdatedByID:  block.UpdatedByID,
//			IsArchived:   block.IsArchived,
//			Attributes:   block.Attributes,
//			Mentions:     block.Mentions,
//			CommentCount: block.CommentCount,
//			ReelCount:    block.ReelCount,
//			Rank:         block.Rank,
//			Contributors: block.Contributors,
//			CreatedByUser: shared.CommonUserMiniSerializer{
//				Id:        cu.ID,
//				UUID:      cu.UUID.String(),
//				FullName:  cu.FullName,
//				Username:  cu.Username,
//				AvatarUrl: cu.AvatarUrl,
//			},
//			UpdatedByUser: shared.CommonUserMiniSerializer{
//				Id:        uu.ID,
//				UUID:      uu.UUID.String(),
//				FullName:  uu.FullName,
//				Username:  uu.Username,
//				AvatarUrl: uu.AvatarUrl,
//			},
//		})
//	}
//	metadata := BlockBulkResponseMetadata{
//		CreatedBlockIDs:         meta.CreatedBlockIDs,
//		UpdatedBlockIDs:         meta.UpdatedBlockIDs,
//		FailedToUpdatedBlockIDs: meta.FailedToUpdatedBlockIDs,
//		FailedToDeleteBlockIDs:  meta.FailedToDeleteBlockIDs,
//		DeletedBlocksCount:      meta.DeletedBlocksCount,
//		FailedBlocksCount:       meta.FailedBlocksCount,
//	}
//
//	Response := BlockBulkRoughResponse{}
//	Response.Data = blocksData
//	Response.ResponseMetadata = metadata
//	Response.RoughBranchID = roughBranchID
//	return Response
//
//}

func BulkGitBlocksSerializerData(blocks []*models.Block) *[]BulkBlocks {
	var blocksData []BulkBlocks
	for _, block := range blocks {
		blocksData = append(blocksData, BulkBlocks{
			UUID:         block.UUID,
			Version:      block.Version,
			Type:         block.Type,
			Children:     block.Children,
			Attributes:   block.Attributes,
			CreatedAt:    block.CreatedAt,
			UpdatedAt:    block.UpdatedAt,
			Rank:         block.Rank,
			Contributors: block.Contributors,
		})
	}
	return &blocksData
}
