package post

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/internal/shared"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"time"
)

type PostCommentSerializer struct {
	ID                  uint64                          `json:"id"`
	UUID                string                          `json:"uuid"`
	PostID              uint64                          `json:"postID"`
	IsEdited            bool                            `json:"isEdited"`
	Comment             string                          `json:"comment"`
	CreatedAt           time.Time                       `json:"createdAt"`
	CreatedByID         uint64                          `json:"createdById"`
	UpdatedByID         uint64                          `json:"updatedById"`
	CreatedByUser       shared.CommonUserMiniSerializer `json:"createdByUser"`
	UpdatedByUser       shared.CommonUserMiniSerializer `json:"updatedByUser"`
	ParentPostCommentID *uint64                         `json:"parentPostCommentID"`
	CommentCount        uint                            `json:"commentCount"`
	ReactionCounter     []shared.ReactedCount           `json:"reactions"`
	ReactionCopy        string                          `json:"reactionCopy"`
}

func SinglePostCommentSerializerData(model *models.PostComment, currentUserID uint64) *PostCommentSerializer {
	//creatorByUser, _ := user2.App.Repo.GetUser(map[string]interface{}{"id": model.CreatedByID})
	//updatedByUser, _ := user2.App.Repo.GetUser(map[string]interface{}{"id": model.CreatedByID})

	creatorByUser, _ := queries.App.UserQueries.GetUserByID(model.CreatedByID)
	updatedByUser, _ := queries.App.UserQueries.GetUserByID(model.UpdatedByID)

	view := PostCommentSerializer{
		ID:                  model.ID,
		UUID:                model.UUID.String(),
		CreatedAt:           model.CreatedAt,
		CreatedByID:         model.CreatedByID,
		UpdatedByID:         model.UpdatedByID,
		PostID:              model.PostID,
		Comment:             model.Comment,
		IsEdited:            model.IsEdited,
		ParentPostCommentID: model.ParentPostCommentID,
		CommentCount:        model.CommentCount,
		CreatedByUser: shared.CommonUserMiniSerializer{
			Id:        creatorByUser.ID,
			UUID:      creatorByUser.UUID.String(),
			FullName:  creatorByUser.FullName,
			Username:  creatorByUser.Username,
			AvatarUrl: creatorByUser.AvatarUrl,
		},
		UpdatedByUser: shared.CommonUserMiniSerializer{
			Id:        updatedByUser.ID,
			UUID:      updatedByUser.UUID.String(),
			FullName:  updatedByUser.FullName,
			Username:  updatedByUser.Username,
			AvatarUrl: updatedByUser.AvatarUrl,
		},
		ReactionCounter: TransposeReactionComment(model, currentUserID),
	}
	return &view
}

func ManyPostCommentsSerializerData(modelInstances *[]models.PostComment, currentUserID uint64) *[]PostCommentSerializer {
	comments := &[]PostCommentSerializer{}
	if len(*modelInstances) == 0 {
		return comments
	}
	for _, model := range *modelInstances {
		*comments = append(*comments, *SinglePostCommentSerializerData(&model, currentUserID))
	}
	return comments
}

func TransposeReactionComment(model *models.PostComment, loggedInUserID uint64) []shared.ReactedCount {
	type ReactionData struct {
		Emoji string `json:"emoji"`
		Count int    `json:"count"`
	}
	var Reactions []ReactionData
	postgres.GetDB().Raw("SELECT emoji, COUNT(emoji) as Count FROM post_comment_reactions where  post_comment_id = ? GROUP BY emoji", model.ID).Scan(&Reactions)
	var updated []shared.ReactedCount
	for _, v := range Reactions {
		// We need to query DB for Existence of the User Reaction
		var exists bool
		_ = postgres.GetDB().Raw("SELECT EXISTS (SELECT 1 from post_comment_reactions WHERE created_by_id = ? and post_id = ? and post_comment_id = ? and emoji = ?)", loggedInUserID, model.PostID, model.ID, v.Emoji).Scan(&exists)
		updated = append(updated, shared.ReactedCount{
			Emoji:   v.Emoji,
			Count:   v.Count,
			Reacted: exists,
		})
	}
	return updated
}
