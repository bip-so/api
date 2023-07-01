package workflows

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
)

type InitCanvasRepoPost struct {
	CollectionID             uint64 `json:"collectionID" binding:"required"`
	Name                     string `json:"name" binding:"required"`
	Icon                     string `json:"icon"`
	Position                 uint   `json:"position" binding:"required"`
	ParentCanvasRepositoryID uint64 `json:"parentCanvasRepositoryID"`
}

func WorkflowCreateCanvasRepoInsideCollection(collectionId uint64, collectionCreatedBy uint64, collectionName string, icon string, position uint, studioID uint64, userInstance models.User) (*models.CanvasRepository, error) {
	canvasView := InitCanvasRepoPost{
		CollectionID: collectionId,
		Name:         collectionName,
		Icon:         icon,
		Position:     position,
	}
	repo, errCreatingCanvasBranch := WorkflowHelperInitCanvasRepo(canvasView, collectionCreatedBy, studioID, userInstance)
	return repo, errCreatingCanvasBranch
}

func WorkflowHelperInitCanvasRepo(body InitCanvasRepoPost, userID uint64, studioID uint64, user models.User) (*models.CanvasRepository, error) {
	// Create repo
	repoInstance, errRepoCreate := queries.App.RepoQuery.CreateRepo(body.ParentCanvasRepositoryID, userID, studioID, body.CollectionID, body.Name, body.Icon, body.Position)
	if errRepoCreate != nil {
		return nil, errRepoCreate
	}
	// create branch
	branchInstance, errBranchCreate := queries.App.BranchQuery.CreateBranch(userID, repoInstance.ID, models.CANVAS_BRANCH_NAME_MAIN, true,
		models.PRIVATE,
	)
	if errBranchCreate != nil {
		return nil, errBranchCreate
	}
	// attached the newly created branch as default branch to fresh repo
	_, branchRepoAttachError := queries.App.RepoQuery.UpdateRepo(repoInstance.ID, map[string]interface{}{
		"default_branch_id": branchInstance.ID,
	})
	if branchRepoAttachError != nil {
		return nil, branchRepoAttachError
	}
	// get updated repo
	UpdatedRepoInstance, _ := queries.App.RepoQuery.GetRepo(map[string]interface{}{"id": repoInstance.ID})

	// Creating a First block
	erCreatingFirstBlankBlock := queries.App.BlockQuery.CreateFirstBlock(user, branchInstance, UpdatedRepoInstance)
	if branchRepoAttachError != nil {
		return nil, erCreatingFirstBlankBlock
	}

	err := queries.App.PermsQuery.CreateDefaultCanvasBranchPermission(UpdatedRepoInstance.CollectionID, userID, studioID, UpdatedRepoInstance.ID, branchInstance.ID, UpdatedRepoInstance.ParentCanvasRepositoryID)
	if err != nil {
		return nil, err
	}

	return UpdatedRepoInstance, nil
}
