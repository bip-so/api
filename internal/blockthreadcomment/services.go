package blockThreadCommentcomment

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gorm.io/datatypes"
)

func (s blockThreadCommentService) InitBlockThreadInstance(
	userID uint64, threadID uint64, parentID *uint64, position uint, data datatypes.JSON, isReply bool, isEdited bool,
) *models.BlockComment {
	var ThisParentID *uint64
	if *parentID == 0 {
		ThisParentID = nil
	} else {
		ThisParentID = parentID
	}

	CommentBase := models.CommentBase{
		position,
		data,
		isEdited,
		isReply,
	}
	return &models.BlockComment{
		CreatedByID: userID,
		UpdatedByID: userID,
		ThreadID:    threadID,
		ParentID:    ThisParentID,
		CommentBase: CommentBase,
	}
}

func (s blockThreadCommentService) Create(body *PostBlockThreadComment, userID uint64) (*models.BlockComment, error) {
	instance := s.InitBlockThreadInstance(userID, body.ThreadID, body.ParentID, body.Position, body.Data, body.IsReply, body.IsEdited)
	created, err := App.Repo.Create(instance)
	return created, err
}

func (s blockThreadCommentService) Update(body *PatchBlockThreadComment, userID uint64) error {

	updates := map[string]interface{}{
		"updated_by_id": userID,
		"thread_id":     body.ThreadID,
		"parent_id":     body.ParentID,
		"data":          body.Data,
		"position":      body.Position,
	}
	err := App.Repo.Update(body.ID, updates)
	return err
}
