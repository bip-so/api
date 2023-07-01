package blockThreadCommentcomment

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/reactions"
)

func (c blockThreadCommentController) Get(blockThreadID uint64, user *models.User) (*[]DefaultSerializer, error) {

	blockThreadComments, err := App.Repo.GetAllComments(map[string]interface{}{"thread_id": blockThreadID})
	if err != nil {
		return nil, err
	}
	btcReactions := []models.BlockCommentReaction{}
	if user != nil {
		btcReactions, _ = reactions.App.Repo.GetBlockThreadCommentReaction(map[string]interface{}{"block_thread_id": blockThreadID, "created_by_id": user.ID})
	}
	fmt.Println(btcReactions)

	blockThreadCommentSerializerData := SerializeDefaultManyBlockThreadCommentWithReaction(blockThreadComments, btcReactions, user)
	return blockThreadCommentSerializerData, nil
}

func (c blockThreadCommentController) GetReply(blockThreadID uint64, parentCommentID uint64) (*[]DefaultSerializer, error) {

	blockThreadComments, err := App.Repo.GetAllComments(map[string]interface{}{"thread_id": blockThreadID, "parent_id": parentCommentID})
	if err != nil {
		return nil, err
	}

	blockThreadCommentSerializerData := SerializeDefaultManyBlockThreadComment(blockThreadComments)
	return blockThreadCommentSerializerData, nil
}

func (c blockThreadCommentController) Create(body PostBlockThreadComment, user *models.User, studioID uint64) (*DefaultSerializer, error) {

	blockThreadComment, err := App.Service.Create(&body, user.ID)
	if err != nil {
		return nil, err
	}
	blockThreadComment.CreatedByUser = user

	blockThreadCommentSerializerData := SerializeDefaultBlockThreadComment(blockThreadComment)

	go func() {
		blockThread, _ := App.Repo.GetBlockThread(map[string]interface{}{"id": blockThreadComment.ThreadID})
		extraData := notifications.NotificationExtraData{
			CanvasRepoID:   blockThread.CanvasRepositoryID,
			CanvasBranchID: blockThread.CanvasBranchID,
		}
		contentObject := models.BLOCK_THREAD_COMMENT
		notifications.App.Service.PublishNewNotification(notifications.CommentReply,
			user.ID, nil, &studioID, nil, extraData, &blockThreadComment.ID, &contentObject)
	}()
	return blockThreadCommentSerializerData, nil
}

func (c blockThreadCommentController) Update(body PatchBlockThreadComment, user *models.User) error {

	err := App.Service.Update(&body, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (c blockThreadCommentController) Delete(blockThreadCommentID uint64, user *models.User) error {

	err := App.Repo.Delete(blockThreadCommentID, user.ID)
	if err != nil {
		return err
	}

	return nil
}
