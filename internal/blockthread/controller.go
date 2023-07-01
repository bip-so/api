package blockthread

import (
	"errors"

	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/reactions"
)

func (c blockThreadController) Get(blockThreadID uint64, user *models.User) (*DefaultSerializer, error) {

	blockThread, err := App.Repo.Get(map[string]interface{}{"id": blockThreadID})
	if err != nil {
		return nil, err
	}

	var userID uint64
	if user == nil {
		userID = 0
	}
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(userID, blockThread.CanvasBranchID, permissiongroup.CANVAS_BRANCH_VIEW); err != nil || !hasPermission {
		return nil, errors.New(response.NoPermissionError)
	}

	btReactions := []models.BlockThreadReaction{}
	if user != nil {
		btReactions, _ = reactions.App.Repo.GetBlockThreadReaction(map[string]interface{}{"block_thread_id": blockThreadID, "created_by_id": user.ID})
	}
	blockThreadSerializerData := SerializeBlockThreadWithReactionForUser(blockThread, btReactions, user)

	return blockThreadSerializerData, nil
}
func (c blockThreadController) GetAllByBranch(canvasBranchID uint64, user *models.User, showResolved string) (*[]DefaultSerializer, error) {
	var blockThreads *[]models.BlockThread
	var errorGettingBlockThreads error
	if showResolved == "true" {
		blockThreads, errorGettingBlockThreads = App.Repo.GetAllThread(map[string]interface{}{"canvas_branch_id": canvasBranchID, "is_archived": false})
	} else {
		blockThreads, errorGettingBlockThreads = App.Repo.GetAllThread(map[string]interface{}{"canvas_branch_id": canvasBranchID, "is_resolved": false, "is_archived": false})
	}

	if errorGettingBlockThreads != nil {
		return nil, errorGettingBlockThreads
	}

	btReactions := []models.BlockThreadReaction{}
	if user != nil {
		btReactions, _ = reactions.App.Repo.GetBlockThreadReaction(map[string]interface{}{"canvas_branch_id": canvasBranchID, "created_by_id": user.ID})
	}
	blockThreadSerializerData := SerializeDefaultManyBlockThreadWithReaction(blockThreads, btReactions, user)

	return blockThreadSerializerData, nil
}

func (c blockThreadController) Create(body PostBlockThread, user *models.User) (*DefaultSerializer, error) {

	blockThread, err := App.Service.Create(&body, user.ID)
	if err != nil {
		return nil, err
	}
	blockThread.CreatedByUser = user
	blockThreadSerializerData := SerializeDefaultBlockThread(blockThread)
	return blockThreadSerializerData, nil
}

func (c blockThreadController) Update(body PatchBlockThread, user *models.User) error {

	err := App.Service.Update(&body, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (c blockThreadController) Delete(blockThreadID uint64, user *models.User) error {
	err := App.Repo.Delete(blockThreadID, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (c blockThreadController) DeleteClonedCommentsOnRoughBranch(blockThread *models.BlockThread, branch *models.CanvasBranch, user *models.User) error {
	// get ids for all rough branches.
	getRoughBranchIds := App.Repo.AllRoughBranchesForGivenBranch(branch.ID)
	// get ids for all rough branches.
	for _, roughBranchID := range getRoughBranchIds {
		// delete the comment with this branch and blockthread iD
		App.Repo.DeleteClonedCommentOnRoughBranch(roughBranchID, blockThread.ID)
	}
	return nil
}

func (c blockThreadController) Resolve(threadID uint64, userID uint64) error {
	err := App.Repo.Resolve(threadID, userID)
	if err != nil {
		return err
	}
	return nil
}
