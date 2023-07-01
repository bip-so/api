package canvasbranch

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	stores "gitlab.com/phonepost/bip-be-platform/pkg/stores/git"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

type CanvasBranchServiceInterface interface {
	EmptyCanvasBranchInstance() *models.CanvasBranch
}

type BlockOperations interface {
	Clone(id uint64)
	Delete(id uint64)
	ChangeBranch(from uint64, to uint64)
}

// Returns Empty Instance
func (s canvasBranchService) EmptyCanvasBranchInstance() *models.CanvasBranch {
	return &models.CanvasBranch{}
}

func (s canvasBranchService) GetCanvasBranchInstance(canvasBranchID uint64) (*models.CanvasBranch, error) {
	//branch, err := queries.App.BranchQuery.GetBranchByID.GetBranchWithRepoAndStudio(canvasBranchID)
	branch, err := App.Repo.Get(map[string]interface{}{"id": canvasBranchID})
	if err != nil {
		return nil, err
	}
	return branch, nil
}

func (s canvasBranchService) Create(userID uint64, canvasRepoID uint64, name string, isTrue bool) (*models.CanvasBranch, error) {
	instance := s.EmptyCanvasBranchInstance()
	instance.CreatedByID = userID
	instance.UpdatedByID = userID
	instance.Name = name
	instance.CanvasRepositoryID = canvasRepoID
	instance.IsDefault = isTrue
	instance.Key = utils.NewNanoid()
	created, err := App.Repo.Create(instance)
	return created, err
}

// CreateRoughBranch Get a Rough Branch Instance and Error
func (s canvasBranchService) CreateRoughBranch(parentBranch *models.CanvasBranch, userID uint64, branchName, key, commitID string) (*models.CanvasBranch, error) {
	instance, err := App.Repo.CreateRoughBranch(parentBranch, userID, branchName, key, commitID)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (s canvasBranchService) UpdateCanvasBranchVisibility(branchID uint64, userID uint64, visibility string) error {
	err := App.Repo.UpdateBranchInstance(branchID, map[string]interface{}{"updated_by_id": userID, "public_access": visibility})
	if err != nil {
		return err
	}
	return nil
}

func (s canvasBranchService) UpdateCanvasLanguageBranchesVisibility(branchID uint64, userID uint64, visibility string) error {
	//canvasBranch, err := App.Repo.Get(map[string]interface{}{"id": branchID})
	canvasBranch, err := queries.App.BranchQuery.GetBranchWithRepoAndStudio(branchID)

	if err != nil {
		return err
	}
	canvasRepos, err := App.Repo.GetCanvasRepos(map[string]interface{}{"default_language_canvas_repo_id": canvasBranch.CanvasRepositoryID})
	if err != nil {
		return err
	}
	for _, repo := range *canvasRepos {
		App.Repo.UpdateBranchInstance(*repo.DefaultBranchID, map[string]interface{}{"updated_by_id": userID, "public_access": visibility})
	}
	return nil
}

func (s canvasBranchService) PublishCanvasBranch(branchID uint64, user *models.User, status bool) error {
	branch, err := queries.App.BranchQuery.GetBranchByID(branchID)
	if err != nil {
		return err
	}
	if status {
		//branchCreatedUser, _ := user2.App.Repo.GetUser(map[string]interface{}{"id": branch.CreatedByID})
		branchCreatedUser, _ := queries.App.UserQueries.GetUserByID(branch.CreatedByID)

		err = App.Git.CommitBranchToGit(branchCreatedUser, branchID, "Initial Snapshot")
		if err != nil {
			fmt.Println("failed git commit")
			return err
		}
		go App.Git.FetchAndUpdateLatestBlockAttributionsForBranch(user, *branch.CanvasRepository.DefaultBranchID)
	}
	err = App.Repo.UpdateBranchInstance(branchID, map[string]interface{}{"updated_by_id": user.ID, "is_draft": status})
	if err != nil {
		return err
	}
	err = App.Repo.PublishCanvasRepository(branch.CanvasRepositoryID)
	if err != nil {
		return err
	}

	// Add permissions of parent collection/canvasBranch to this canvasBranch.
	// Parent as canvas
	// Get all the canvasBranch permissions of the parent canvasRepo Default branch and add it to this.
	// Parent as collection
	// Get all the collection permissions and match it with the canvasBranch permissions and create canvasBranch permissions
	// canvasRepoID, parentCanvasRepoID, parentCollectionID
	go permissions.App.Service.InheritParentPermissions(branch.ID)

	return nil
}

// checkIfBranchIsInsideRootRepo:  we need to check if the given Branch mostly Default branch is inside a repo
// Which is ROOT or have a Parent Repo
func (s canvasBranchService) doesBelongToRootRepo(branch models.CanvasBranch) bool {
	// Branch -> CanvasRepository -> ParentCanvasRepositoryID == 0 for Root.
	if branch.CanvasRepository.ParentCanvasRepositoryID == nil {
		return true
	}
	return false
}

// Given a Branch, We need to Return the Default branch on the Parent Repo not this Repo
func (s canvasBranchService) defaultBranchOnParentCanvasRepo(branch models.CanvasBranch) uint64 {
	parentRepoID := branch.CanvasRepository.ParentCanvasRepositoryID
	parentRepoInstance, err := App.Repo.GetCanvasRepoInstance(map[string]interface{}{"id": parentRepoID})
	if err != nil {
		return 0
	}
	return *parentRepoInstance.DefaultBranchID
}

func (s canvasBranchService) draftBranches(studioID uint64, userID uint64) (*[]models.CanvasBranch, error) {
	var branches *[]models.CanvasBranch
	var err error
	branches, err = App.Repo.GetRoughUnpublishedByStudioAndUser(studioID, userID)
	if err != nil {
		return nil, err
	}
	return branches, nil
}

func (s canvasBranchService) BuildRootNavObject(branchID uint64, userID uint64, public string) (*[]CollectionRootNavSerializer, error) {
	branch, err := App.Repo.GetBranchRepoCollection(branchID)

	if err != nil {
		return nil, err
	}
	//fmt.Println(branch)
	return RootSerialized(branch, userID, public), nil

}

func (s canvasBranchService) BuildRootNodeObject(branchID uint64, userID uint64) (*[]CanvasRepoDefaultSerializer, error) {

	branch, err := App.Repo.GetBranchRepoCollection(branchID)

	if err != nil {
		return nil, err
	}
	//fmt.Println(branch)
	return NodeSerialized(branch, userID), nil

}

func (s canvasBranchService) BranchHistory(user *models.User, branchId uint64, startCommitID string) ([]*stores.GitLog, *[]models.User, string, error) {
	return App.Git.FetchBranchHistoryFromGit(user, branchId, startCommitID)
}

func (s canvasBranchService) GetBranchAttributions(user *models.User, branchId uint64) (*[]models.Attribution, error) {
	return App.Git.FetchAllAttributionsForBranch(user, branchId)
}

func (s canvasBranchService) GetCanvasBranchData(canvasBranchID uint64, authUser *models.User, inviteCode string) (*BranchMeta, map[string]interface{}) {
	var authUserID uint64
	if authUser != nil {
		authUserID = authUser.ID
	}
	canUserViewBranch, errGettingPermissions := permissions.App.Service.CanUserDoThisOnBranch(authUserID, canvasBranchID, permissiongroup.CANVAS_BRANCH_VIEW)
	if errGettingPermissions != nil {
		return nil, map[string]interface{}{
			"error": errGettingPermissions.Error(),
		}
	}

	branch, err := App.Repo.GetBranchWithRepo(map[string]interface{}{"id": canvasBranchID})
	if err != nil {
		return nil, map[string]interface{}{
			"error": err.Error(),
		}
	}

	if authUser == nil && inviteCode != "" {
		canvasBranchAccessToken, err := permissions.App.Service.CheckBranchAccessToken(inviteCode, branch.CanvasRepository.Key)
		if err != nil || canvasBranchAccessToken == nil {
			return nil, map[string]interface{}{
				"error": "Anonymous user does not have permissions to view this.",
			}
		}
	} else if branch.PublicAccess == "private" && !canUserViewBranch {
		// We are adding a requested is User already have Access Request
		exists := queries.App.AccessRequestQuery.AccessRequestExistsSimple(canvasBranchID, authUserID)
		return nil, map[string]interface{}{
			"error":            "User does not have permissions to view branch",
			"access_requested": exists,
		}
	}

	attributions, _ := App.Service.GetBranchAttributions(authUser, canvasBranchID)
	branchData := BranchMetaWithMRSerializer(branch, canvasBranchID, authUserID, attributions)
	return &branchData, nil
}
