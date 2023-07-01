package canvasrepo

import (
	"encoding/json"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/internal/studio"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/apiClient"
	"gorm.io/gorm"

	"github.com/google/uuid"
	"gitlab.com/phonepost/bip-be-platform/internal/blocks"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

func (c canvasRepoController) UpdateCanvasRepo(repoId uint64, data UpdateCanvasRepoPost, userId uint64) bool {
	return App.Service.UpdateCanvasRepo(data.Name, data.Icon, repoId, userId, data.CoverUrl)
}

// DeleteThis
//func (c canvasRepoController) InitCanvasRepo(
//	body InitCanvasRepoPost, userID uint64, studioID uint64, user models.User) (*models.CanvasRepository, error) {
//	// create repo
//	repoInstance, errRepoCreate := App.Service.Create(
//		body.ParentCanvasRepositoryID,
//		userID,
//		studioID,
//		body.CollectionID,
//		body.Name,
//		body.Icon,
//		body.Position,
//	)
//	if errRepoCreate != nil {
//		return nil, errRepoCreate
//	}
//	fmt.Println(repoInstance)
//	// create branch
//	branchInstance, errBranchCreate := App.Service.CreateCanvasBranch(
//		userID,
//		repoInstance.ID,
//		models.CANVAS_BRANCH_NAME_MAIN,
//		true,
//		models.PRIVATE,
//	)
//	if errBranchCreate != nil {
//		return nil, errBranchCreate
//	}
//	fmt.Println(branchInstance)
//
//	_, branchRepoAttachError := App.Service.UpdateCanvasRepoDefaultBranch(repoInstance.ID, branchInstance.ID)
//	// Now we need to set repoInstance.DefaultBranchID = branchInstance.ID save to Database.
//	if branchRepoAttachError != nil {
//		return nil, branchRepoAttachError
//	}
//	repo, _ := App.Repo.GetRepo(map[string]interface{}{"id": repoInstance.ID})
//
//	// Creating a blank block
//	var block models.Block
//	firstContrib := App.Service.BlockContributorFirst(user, branchInstance.ID)
//	newBlock, errCreatingBlockInst := block.NewBlock(
//		uuid.New(),
//		userID,
//		repoInstance.ID,
//		branchInstance.ID,
//		nil,
//		1,
//		2,
//		models.BlockTypeText,
//		models.MyFirstBlockJson(),
//		models.MyFirstEmptyBlockJson(),
//		firstContrib,
//	)
//	if errCreatingBlockInst != nil {
//		return nil, errCreatingBlockInst
//	}
//	errNewBlock := blocks.App.Service.Create(newBlock)
//	if errNewBlock != nil {
//		return nil, errNewBlock
//	}
//
//	err := permissions.App.Service.CreateDefaultCanvasBranchPermission(repo.CollectionID, userID, studioID, repo.ID, branchInstance.ID, repo.ParentCanvasRepositoryID)
//	if err != nil {
//		return nil, err
//	}
//
//	return repo, nil
//}

/*
	Get all the canvas data based on parentCollectionId or parentCanvasRepoId for non-logged-in users.
	Args:
		body *GetAllCanvasValidator
	Returns:
		*[]CanvasRepoDefaultSerializer
		error
*/
func (c canvasRepoController) AnonymousGetAllCanvasController(body *GetAllCanvasValidator) (*[]CanvasRepoDefaultSerializer, error) {
	var canvasRepos *[]models.CanvasRepository
	var err error
	publicAccess := []string{"view", "edit", "comment"}

	if body.ParentCanvasRepositoryID == 0 {
		canvasRepos, err = App.Repo.GetAnonymousCanvasRepos(body.ParentCollectionID, publicAccess)
	} else {
		canvasRepos, err = App.Repo.GetAnonymousSubCanvasRepos(body.ParentCanvasRepositoryID, publicAccess)
	}

	if err != nil {
		return nil, err
	}
	canvasRepoViews := MultiSerializeDefaultCanvasRepo(canvasRepos)
	return canvasRepoViews, nil
}

// Get Anno One
func (c canvasRepoController) AnonymousGetOneCanvasByKeyController(key string) (*CanvasRepoDefaultSerializer, error) {
	var canvasRepo *models.CanvasRepository
	var err error
	canvasRepo, err = App.Repo.GetCanvasRepoByKey(key)
	if err != nil {
		return nil, err
	}

	canvasRepoViews := SerializeDefaultCanvasRepo(canvasRepo)
	canvasRepoViews.DefaultBranch.Permission = models.PGCanvasNoneSysName
	pr, err := App.Repo.GetAcceptedPublishRequestByBranch(canvasRepo.ID)
	if err == nil && pr != nil {
		canvasRepoViews.DefaultBranch.IsPublishedBy = pr.ReviewedByUserID
	} else {
		canvasRepoViews.DefaultBranch.IsPublishedBy = canvasRepo.CreatedByID
	}
	return canvasRepoViews, nil
}

func (c canvasRepoController) AuthGetOneCanvasByKeyController(key string, userID uint64) (*CanvasRepoDefaultSerializer, error) {
	var canvasRepo *models.CanvasRepository
	var err error
	canvasRepo, err = App.Repo.GetCanvasRepoByKey(key)
	if err != nil {
		return nil, err
	}
	canvasRepoViews := SerializeDefaultCanvasRepo(canvasRepo)
	var permissionsList map[uint64]map[uint64]string
	var hasPermission bool
	if canvasRepo.ParentCanvasRepositoryID == nil {
		permissionsList, err = permissions.App.Service.CalculateCanvasRepoPermissions(userID, canvasRepo.StudioID, canvasRepo.CollectionID)
		hasPermission, _ = permissions.App.Service.CanUserDoThisOnCollection(userID, canvasRepo.StudioID, canvasRepo.CollectionID, permissiongroup.COLLECTION_MANAGE_PUBLISH_REQUEST)
	} else {
		permissionsList, err = permissions.App.Service.CalculateSubCanvasRepoPermissions(userID, canvasRepo.StudioID, canvasRepo.CollectionID, *canvasRepo.ParentCanvasRepositoryID)
		hasPermission, _ = permissions.App.Service.CanUserDoThisOnBranch(userID, *canvasRepo.ParentCanvasRepository.DefaultBranchID, permissiongroup.CANVAS_BRANCH_MANAGE_PUBLISH_REQUESTS)
	}
	canvasRepoViews.DefaultBranch.Permission = permissionsList[canvasRepo.ID][*canvasRepo.DefaultBranchID]
	pr, err := App.Repo.GetAcceptedPublishRequestByBranch(canvasRepo.ID)
	if err == nil && pr != nil {
		canvasRepoViews.DefaultBranch.IsPublishedBy = pr.ReviewedByUserID
	} else {
		canvasRepoViews.DefaultBranch.IsPublishedBy = canvasRepo.CreatedByID
	}
	canvasRepoViews.DefaultBranch.CanPublish = hasPermission
	return canvasRepoViews, nil
}

/*
	Get all the canvas data based on parentCollectionId or parentCanvasRepoId with user permissions.
	Args:
		body *GetAllCanvasValidator
	Returns:
		*[]CanvasRepoDefaultSerializer
		error
*/
func (c canvasRepoController) AuthUserGetAllCanvasController(body *GetAllCanvasValidator, user *models.User, studioId uint64) (*[]CanvasRepoDefaultSerializer, error) {
	var canvasRepos *[]models.CanvasRepository
	var permissionsList map[uint64]map[uint64]string
	canvasRepoViews := &[]CanvasRepoDefaultSerializer{}

	var err error

	if body.ParentCanvasRepositoryID == 0 {
		canvasRepos, err = App.Repo.GetCanvasRepos(map[string]interface{}{"collection_id": body.ParentCollectionID, "parent_canvas_repository_id": nil, "is_archived": false, "is_processing": false})
		if err != nil {
			return nil, err
		}
		permissionsList, err = permissions.App.Service.CalculateCanvasRepoPermissions(user.ID, studioId, body.ParentCollectionID)
		if err != nil {
			return nil, err
		}

	} else {
		canvasRepos, err = App.Repo.GetCanvasRepos(map[string]interface{}{"parent_canvas_repository_id": body.ParentCanvasRepositoryID, "is_archived": false, "is_processing": false})
		if err != nil {
			return nil, err
		}
		permissionsList, err = permissions.App.Service.CalculateSubCanvasRepoPermissions(user.ID, studioId, body.ParentCollectionID, body.ParentCanvasRepositoryID)
		if err != nil {
			return nil, err
		}
	}
	for _, repo := range *canvasRepos {
		repoPermissions := permissionsList[repo.ID]
		permissionValues := utils.Values(repoPermissions)
		hasPermission := false
		for _, perm := range permissionValues {
			if utils.Contains(permissiongroup.UserAccessCanvasPermissionsList, perm) {
				repoView := SerializeDefaultCanvasRepo(&repo)
				branchPerm := repoPermissions[repo.DefaultBranch.ID]
				repoView.DefaultBranch.Permission = branchPerm
				*canvasRepoViews = append(*canvasRepoViews, *repoView)
				hasPermission = true
				break
			}
		}
		if !hasPermission {
			if repo.HasPublicCanvas && repo.IsPublished {
				repoView := SerializeDefaultCanvasRepo(&repo)
				branchPerm := permissiongroup.PGCanvasViewMetaData().SystemName
				repoView.DefaultBranch.Permission = branchPerm
				*canvasRepoViews = append(*canvasRepoViews, *repoView)
			} else if repo.DefaultBranch.PublicAccess == models.EDIT && repo.IsPublished {
				repoView := SerializeDefaultCanvasRepo(&repo)
				repoView.DefaultBranch.Permission = models.PGCanvasEditSysName
				*canvasRepoViews = append(*canvasRepoViews, *repoView)
			} else if repo.DefaultBranch.PublicAccess == models.COMMENT && repo.IsPublished {
				repoView := SerializeDefaultCanvasRepo(&repo)
				repoView.DefaultBranch.Permission = models.PGCanvasCommentSysName
				*canvasRepoViews = append(*canvasRepoViews, *repoView)
			} else if repo.DefaultBranch.PublicAccess == models.VIEW && repo.IsPublished {
				repoView := SerializeDefaultCanvasRepo(&repo)
				repoView.DefaultBranch.Permission = models.PGCanvasViewSysName
				*canvasRepoViews = append(*canvasRepoViews, *repoView)
			}
		}
	}

	if err != nil {
		return nil, err
	}
	//canvasRepoViews := MultiSerializeDefaultCanvasRepo(&accessCanvasRepo)
	return canvasRepoViews, nil
}

func (c canvasRepoController) moveCanvasBetweenCollections(body *MoveCanvasRepoPost) error {
	canvasRepo, err := queries.App.RepoQuery.GetRepo(map[string]interface{}{"id": body.CanvasRepoID})
	if err != nil {
		return err
	}

	err = App.Repo.MoveCanvasRepoWithDifferentParentCollection(canvasRepo, body.FuturePosition, body.ToCollectionID)

	// update canvas branch permission collectionID
	branchPerms, err := App.Repo.UpdateCanvasRepoBranchPermissions(
		map[string]interface{}{"canvas_repository_id": canvasRepo.ID},
		map[string]interface{}{"collection_id": body.ToCollectionID, "cbp_parent_canvas_repository_id": nil})
	if err != nil {
		return err
	}

	// Invalidate the collection cache
	go func() {
		branchPerms, _ = App.Repo.GetCanvasBranchPermissions(map[string]interface{}{"collection_id": body.ToCollectionID, "cbp_parent_canvas_repository_id": nil})
		for _, perm := range branchPerms {
			if perm.Member != nil {
				fmt.Println(perm.Member.UserID, perm.StudioID, perm.CollectionId)
				permissions.App.Service.InvalidateCanvasPermissionCache(perm.Member.UserID, perm.StudioID, perm.CollectionId)
				permissions.App.Service.AddMemberToCollectionIfNotPresent(0, *perm.MemberId, perm.CollectionId, perm.StudioID)
				if perm.CbpParentCanvasRepositoryID != nil {
					permissions.App.Service.AddMemberToCanvasIfNotPresent(0, *perm.MemberId, *perm.CbpParentCanvasRepositoryID, perm.StudioID)
				}
			} else if perm.Role != nil {
				fmt.Println(*perm.RoleId, perm.StudioID, perm.CollectionId)
				permissions.App.Service.InvalidateCanvasPermissionCacheByRole(*perm.RoleId, perm.StudioID, perm.CollectionId)
				permissions.App.Service.AddRoleToCollectionIfNotPresent(0, *perm.RoleId, perm.CollectionId, perm.StudioID)
				if perm.CbpParentCanvasRepositoryID != nil {
					permissions.App.Service.AddMemberToCanvasIfNotPresent(0, *perm.RoleId, *perm.CbpParentCanvasRepositoryID, perm.StudioID)
				}
			}
		}
	}()

	if err != nil {
		return err
	}
	return nil
}

func (c canvasRepoController) moveCanvasBetweenCanvas(body *MoveCanvasRepoPost) error {
	canvasRepo, err := queries.App.RepoQuery.GetRepo(map[string]interface{}{"id": body.CanvasRepoID})
	if err != nil {
		return err
	}

	err = App.Repo.MoveCanvasRepoWithDifferentParentAsCanvas(canvasRepo, body.FuturePosition, body.ToParentCanvasRepositoryID)
	if err != nil {
		return err
	}

	// update canvas branch permission collectionID
	branchPerms, err := App.Repo.UpdateCanvasRepoBranchPermissions(
		map[string]interface{}{"canvas_repository_id": canvasRepo.ID},
		map[string]interface{}{"cbp_parent_canvas_repository_id": body.ToParentCanvasRepositoryID},
	)
	if err != nil {
		return err
	}

	// Invalidate the collection cache
	go func() {
		for _, perm := range branchPerms {
			if perm.Member != nil {
				permissions.App.Service.InvalidateCollectionMatchingPermissionCache(perm.Member.UserID, perm.StudioID, perm.CollectionId)
				permissions.App.Service.AddMemberToCollectionIfNotPresent(0, *perm.MemberId, perm.CollectionId, perm.StudioID)
				if perm.CbpParentCanvasRepositoryID != nil {
					permissions.App.Service.AddMemberToCanvasIfNotPresent(0, *perm.MemberId, *perm.CbpParentCanvasRepositoryID, perm.StudioID)
				}
			} else if perm.Role != nil {
				permissions.App.Service.InvalidateCollectionMatchingPermissionCacheByRole(*perm.RoleId, perm.StudioID, perm.CollectionId)
				permissions.App.Service.AddRoleToCollectionIfNotPresent(0, *perm.RoleId, perm.CollectionId, perm.StudioID)
				if perm.CbpParentCanvasRepositoryID != nil {
					permissions.App.Service.AddMemberToCanvasIfNotPresent(0, *perm.RoleId, *perm.CbpParentCanvasRepositoryID, perm.StudioID)
				}
			}
		}
	}()

	return nil
}

func (c canvasRepoController) moveCanvasBetweenSameCollectionsAndCanvas(body *MoveCanvasRepoPost) error {
	canvasRepo, err := queries.App.RepoQuery.GetRepo(map[string]interface{}{"id": body.CanvasRepoID})
	if err != nil {
		return err
	}

	err = App.Repo.MoveCanvasRepoWithParentAsCollection(canvasRepo, body.FuturePosition)
	if err != nil {
		return err
	}

	// update canvas branch permission collectionID
	branchPerms, err := App.Repo.UpdateCanvasRepoBranchPermissions(
		map[string]interface{}{"canvas_repository_id": canvasRepo.ID},
		map[string]interface{}{"cbp_parent_canvas_repository_id": nil},
	)
	if err != nil {
		return err
	}

	// Invalidate the collection cache
	go func() {
		for _, perm := range branchPerms {
			if perm.Member != nil {
				permissions.App.Service.InvalidateCanvasPermissionCache(perm.Member.UserID, perm.StudioID, perm.CollectionId)
				permissions.App.Service.AddMemberToCollectionIfNotPresent(0, *perm.MemberId, perm.CollectionId, perm.StudioID)
				if perm.CbpParentCanvasRepositoryID != nil {
					permissions.App.Service.AddMemberToCanvasIfNotPresent(0, *perm.MemberId, *perm.CbpParentCanvasRepositoryID, perm.StudioID)
				}
			} else if perm.Role != nil {
				permissions.App.Service.InvalidateCanvasPermissionCacheByRole(*perm.RoleId, perm.StudioID, perm.CollectionId)
				permissions.App.Service.AddRoleToCollectionIfNotPresent(0, *perm.RoleId, perm.CollectionId, perm.StudioID)
				if perm.CbpParentCanvasRepositoryID != nil {
					permissions.App.Service.AddMemberToCanvasIfNotPresent(0, *perm.RoleId, *perm.CbpParentCanvasRepositoryID, perm.StudioID)
				}
			}
		}
	}()

	return nil
}

func (c canvasRepoController) moveCanvasBetweenCanvasAndCollection(body *MoveCanvasRepoPost) error {
	canvasRepo, err := queries.App.RepoQuery.GetRepo(map[string]interface{}{"id": body.CanvasRepoID})
	if err != nil {
		return err
	}

	err = App.Repo.MoveCanvasRepoWithDifferentParentCanvasAndCollection(canvasRepo, body.FuturePosition, body.ToCollectionID, body.ToParentCanvasRepositoryID)
	if err != nil {
		return err
	}

	// update canvas branch permission collectionID
	branchPerms, err := App.Repo.UpdateCanvasRepoBranchPermissions(
		map[string]interface{}{"canvas_repository_id": canvasRepo.ID},
		map[string]interface{}{"collection_id": body.ToCollectionID, "cbp_parent_canvas_repository_id": body.ToParentCanvasRepositoryID},
	)
	if err != nil {
		return err
	}

	// Invalidate the collection cache
	go func() {
		for _, perm := range branchPerms {
			if perm.Member != nil {
				permissions.App.Service.InvalidateCanvasPermissionCache(perm.Member.UserID, perm.StudioID, perm.CollectionId)
				permissions.App.Service.AddMemberToCollectionIfNotPresent(0, *perm.MemberId, perm.CollectionId, perm.StudioID)
				if perm.CbpParentCanvasRepositoryID != nil {
					permissions.App.Service.AddMemberToCanvasIfNotPresent(0, *perm.MemberId, *perm.CbpParentCanvasRepositoryID, perm.StudioID)
				}
			} else if perm.Role != nil {
				permissions.App.Service.InvalidateCanvasPermissionCacheByRole(*perm.RoleId, perm.StudioID, perm.CollectionId)
				permissions.App.Service.AddRoleToCollectionIfNotPresent(0, *perm.RoleId, perm.CollectionId, perm.StudioID)
				if perm.CbpParentCanvasRepositoryID != nil {
					permissions.App.Service.AddMemberToCanvasIfNotPresent(0, *perm.RoleId, *perm.CbpParentCanvasRepositoryID, perm.StudioID)
				}
			}
		}
	}()
	return nil
}

func (c canvasRepoController) CreateLanguageCanvasRepo(body *CreateLanguageValidator, user *models.User) (*[]models.CanvasRepository, error, []string) {
	languageCanvasRepos := []models.CanvasRepository{}
	duplicateLanguageCodes := []string{}
	canvasRepo, err := queries.App.RepoQuery.GetRepo(map[string]interface{}{"id": body.CanvasRepositoryID})
	if err != nil {
		return nil, err, duplicateLanguageCodes
	}
	for _, languageCode := range body.Languages {
		canvasCheck, _ := queries.App.RepoQuery.GetRepo(map[string]interface{}{"default_language_canvas_repo_id": body.CanvasRepositoryID, "language": languageCode})
		if canvasCheck != nil {
			duplicateLanguageCodes = append(duplicateLanguageCodes, languageCode)
			continue
		}
		canvasRepository := App.Service.NewLanguageCanvasRepoInstance(canvasRepo, languageCode, body.AutoTranslate, user.ID)
		languageCanvasRepos = append(languageCanvasRepos, *canvasRepository)
	}
	if len(languageCanvasRepos) == 0 {
		return nil, nil, duplicateLanguageCodes
	}

	err = App.Repo.CreateMultiple(&languageCanvasRepos)
	if err != nil {
		return nil, err, duplicateLanguageCodes
	}

	canvasRepoMainBranches := []models.CanvasBranch{}
	for _, repo := range languageCanvasRepos {
		canvasRepoMainBranches = append(canvasRepoMainBranches, *models.NewCanvasBranch(models.CANVAS_BRANCH_NAME_MAIN, repo.ID, user.ID, canvasRepo.DefaultBranch.PublicAccess))
	}

	err = App.Repo.CreateCanvasBranchMultiple(&canvasRepoMainBranches)
	if err != nil {
		return nil, err, duplicateLanguageCodes
	}

	ids := []uint64{}
	for _, branch := range canvasRepoMainBranches {
		_, branchRepoAttachError := queries.App.RepoQuery.UpdateRepo(branch.CanvasRepositoryID, map[string]interface{}{
			"default_branch_id": branch.ID,
		})

		if branchRepoAttachError != nil {
			return nil, branchRepoAttachError, duplicateLanguageCodes
		}
		ids = append(ids, branch.CanvasRepositoryID)
	}

	lCanvasRepos, err := App.Repo.GetMultipleReposByIDS(ids)
	for _, repo := range *lCanvasRepos {
		_, err = permissions.App.Repo.GetMember(map[string]interface{}{"user_id": user.ID, "studio_id": repo.StudioID})
		if err == gorm.ErrRecordNotFound {
			joinStudioBulk := studio.JoinStudioBulkPost{UsersAdded: []uint64{user.ID}}
			_, err = studio.App.Controller.JoinStudioInBulkController(joinStudioBulk, repo.StudioID, user.ID)
		}
		err = queries.App.PermsQuery.CreateDefaultCanvasBranchPermission(repo.CollectionID, user.ID, repo.StudioID, repo.ID, *repo.DefaultBranchID, repo.ParentCanvasRepositoryID)
		if err != nil {
			fmt.Println("Error in creating default canvas permission", err)
			return nil, err, nil
		}
		//permissions.App.Service.InheritLanguageRepoParentPerms(&repo)
	}

	if body.AutoTranslate {
		for _, repo := range *lCanvasRepos {
			App.Repo.Update(map[string]interface{}{"id": repo.ID}, map[string]interface{}{"is_processing": true})
		}
		lRepoStr, _ := json.Marshal(languageCanvasRepos)
		apiClient.AddToQueue(apiClient.TranslateCanvasRepositories, lRepoStr, apiClient.DEFAULT, apiClient.CommonRetry)
	} else {
		// Creating an empty block
		for _, repo := range *lCanvasRepos {
			var block models.Block
			firstContrib := App.Service.BlockContributorFirst(*user, *repo.DefaultBranchID)
			newBlock, errCreatingBlockInst := block.NewBlock(
				uuid.New(),
				user.ID,
				repo.ID,
				*repo.DefaultBranchID,
				nil,
				1,
				2,
				models.BlockTypeText,
				models.MyFirstBlockJson(),
				models.MyFirstEmptyBlockJson(),
				firstContrib,
			)
			if errCreatingBlockInst != nil {
				return nil, errCreatingBlockInst, nil
			}
			errNewBlock := blocks.App.Service.Create(newBlock)
			if errNewBlock != nil {
				return nil, errNewBlock, nil
			}
		}
	}
	return lCanvasRepos, err, duplicateLanguageCodes
}

// CreateCanvasRepo Create New Canvas Repo
//
//func (c canvasRepoController) CreateCanvasRepo(
//	body NewCanvasRepoPost, userID uint64, studioID uint64, user models.User, parentPublicAccess string) (*models.CanvasRepository, error) {
//	// create repo
//	repoInstance, errRepoCreate := App.Service.Create(
//		body.ParentCanvasRepositoryID,
//		userID,
//		studioID,
//		body.CollectionID,
//		body.Name,
//		body.Icon,
//		body.Position,
//	)
//	if errRepoCreate != nil {
//		return nil, errRepoCreate
//	}
//	// update position of all other canvas
//	if body.Position == 1 {
//		App.Repo.UpdateCanvasPositionOnAddingNewCanvas(repoInstance)
//	}
//	// create branch
//	branchInstance, errBranchCreate := App.Service.CreateCanvasBranch(
//		userID,
//		repoInstance.ID,
//		models.CANVAS_BRANCH_NAME_MAIN,
//		true,
//		parentPublicAccess,
//	)
//	if errBranchCreate != nil {
//		return nil, errBranchCreate
//	}
//
//	_, branchRepoAttachError := App.Service.UpdateCanvasRepoDefaultBranch(repoInstance.ID, branchInstance.ID)
//	// Now we need to set repoInstance.DefaultBranchID = branchInstance.ID save to Database.
//	if branchRepoAttachError != nil {
//		return nil, branchRepoAttachError
//	}
//	repo, _ := App.Repo.GetRepo(map[string]interface{}{"id": repoInstance.ID})
//
//	err := permissions.App.Service.CreateDefaultCanvasBranchPermission(repo.CollectionID, userID, studioID, repo.ID, branchInstance.ID, repo.ParentCanvasRepositoryID)
//	if err != nil {
//		return nil, err
//	}
//
//	// Creating a blank block
//	var block models.Block
//
//	firstContrib := App.Service.BlockContributorFirst(user, branchInstance.ID)
//	newBlock, errCreatingBlockInst := block.NewBlock(
//		uuid.New(),
//		userID,
//		repoInstance.ID,
//		branchInstance.ID,
//		nil,
//		1,
//		2,
//		models.BlockTypeText,
//		models.MyFirstBlockJson(),
//		models.MyFirstEmptyBlockJson(),
//		firstContrib,
//	)
//	if errCreatingBlockInst != nil {
//		return nil, errCreatingBlockInst
//	}
//	errNewBlock := blocks.App.Service.Create(newBlock)
//	if errNewBlock != nil {
//		return nil, errNewBlock
//	}
//
//	return repo, nil
//}

func (c canvasRepoController) MemberCanvasController(body GetAllCanvasValidator, userID uint64, studioId uint64) (*[]CanvasRepoDefaultSerializer, error) {
	var canvasRepos *[]models.CanvasRepository
	var permissionsList map[uint64]map[uint64]string
	var RoleActualPermsObject []RoleBranchActualPermissionsObject
	memberObject, _ := App.Repo.GetMemberByUserID(userID, studioId)

	canvasRepoViews := &[]CanvasRepoDefaultSerializer{}

	var err error

	if body.ParentCanvasRepositoryID == 0 {
		canvasRepos, err = App.Repo.GetCanvasRepos(map[string]interface{}{"collection_id": body.ParentCollectionID, "parent_canvas_repository_id": nil, "is_archived": false})
		if err != nil {
			return nil, err
		}
		permissionsList, err = permissions.App.Service.CalculateCanvasRepoPermissions(userID, studioId, body.ParentCollectionID)
		if err != nil {
			return nil, err
		}
	} else {
		canvasRepos, err = App.Repo.GetCanvasRepos(map[string]interface{}{"parent_canvas_repository_id": body.ParentCanvasRepositoryID, "is_archived": false})
		if err != nil {
			return nil, err
		}
		permissionsList, err = permissions.App.Service.CalculateSubCanvasRepoPermissions(userID, studioId, body.ParentCollectionID, body.ParentCanvasRepositoryID)
		if err != nil {
			return nil, err
		}

	}
	fmt.Println(canvasRepos)
	fmt.Println("permissionsList of canvasrepo -->", permissionsList)
	for _, repo := range *canvasRepos {
		repoPermissions := permissionsList[repo.ID]
		permissionValues := utils.Values(repoPermissions)
		hasPermission := false
		for _, perm := range permissionValues {
			if utils.Contains(permissiongroup.UserAccessCanvasPermissionsList, perm) {
				repoView := SerializeDefaultCanvasRepo(&repo)
				branchPerm := repoPermissions[repo.DefaultBranch.ID]
				fmt.Println("Default Branch ID ->", repo.DefaultBranch.ID)
				fmt.Println("Calculated Perms ->", branchPerm)
				repoView.DefaultBranch.Permission = branchPerm
				repoView.DefaultBranch.MemberPermsObject = MemberBranchPermissionActualCalculator(repo.CollectionID, repo.ID, repo.DefaultBranch.ID, memberObject, studioId)
				actualPerms := RoleBranchPermissionActualCalculator(repo.CollectionID, repo.ID, repo.DefaultBranch.ID, memberObject, studioId)
				for _, ap := range actualPerms {
					RoleActualPermsObject = append(RoleActualPermsObject, ap)
				}
				repoView.DefaultBranch.RolePermsObject = actualPerms
				*canvasRepoViews = append(*canvasRepoViews, *repoView)
				hasPermission = true
				break
			}
		}
		if !hasPermission {
			repoView := SerializeDefaultCanvasRepo(&repo)
			branchPerm := permissiongroup.PGCanvasNone().SystemName
			fmt.Println("Has No perms ")
			fmt.Println("Default Branch ID ->", repo.DefaultBranch.ID)
			fmt.Println("Calculated Perms ->", branchPerm)
			repoView.DefaultBranch.Permission = branchPerm
			// set empty
			// Commenting as added for Debug
			//ActualPermsObject = append(ActualPermsObject, BranchActualPermissionsObject{})
			repoView.DefaultBranch.RolePermsObject = []RoleBranchActualPermissionsObject{}
			repoView.DefaultBranch.MemberPermsObject = MemberBranchActualPermissionsObject{}

			*canvasRepoViews = append(*canvasRepoViews, *repoView)
		}
	}

	if err != nil {
		return nil, err
	}
	//canvasRepoViews := MultiSerializeDefaultCanvasRepo(&accessCanvasRepo)
	//fmt.Println("ActualPermsObject")
	//fmt.Printf("%+v\n", ActualPermsObject)

	return canvasRepoViews, nil
}

func (c canvasRepoController) RoleCanvasController(body GetAllCanvasValidator, roleID uint64) (*[]CanvasRepoDefaultSerializer, error) {
	var canvasRepos *[]models.CanvasRepository
	var permissionsList map[uint64]map[uint64]string
	canvasRepoViews := &[]CanvasRepoDefaultSerializer{}

	var err error

	if body.ParentCanvasRepositoryID == 0 {
		canvasRepos, err = App.Repo.GetCanvasRepos(map[string]interface{}{"collection_id": body.ParentCollectionID, "parent_canvas_repository_id": nil, "is_archived": false})
		if err != nil {
			return nil, err
		}
		permissionsList = permissions.App.Service.CalculateCanvasRolePermissions(roleID, []uint64{body.ParentCollectionID})
	} else {
		canvasRepos, err = App.Repo.GetCanvasRepos(map[string]interface{}{"parent_canvas_repository_id": body.ParentCanvasRepositoryID, "is_archived": false})
		if err != nil {
			return nil, err
		}
		permissionsList = permissions.App.Service.CalculateSubCanvasRolePermissions(roleID, body.ParentCanvasRepositoryID)
	}

	for _, repo := range *canvasRepos {
		repoPermissions := permissionsList[repo.ID]
		permissionValues := utils.Values(repoPermissions)
		hasPermission := false
		for _, perm := range permissionValues {
			if utils.Contains(permissiongroup.UserAccessCanvasPermissionsList, perm) {
				repoView := SerializeDefaultCanvasRepo(&repo)
				branchPerm := repoPermissions[repo.DefaultBranch.ID]
				repoView.DefaultBranch.Permission = branchPerm
				*canvasRepoViews = append(*canvasRepoViews, *repoView)
				hasPermission = true
				break
			}
		}
		if !hasPermission {
			repoView := SerializeDefaultCanvasRepo(&repo)
			branchPerm := permissiongroup.PGCanvasNone().SystemName
			repoView.DefaultBranch.Permission = branchPerm
			*canvasRepoViews = append(*canvasRepoViews, *repoView)
		}
	}

	if err != nil {
		return nil, err
	}
	//canvasRepoViews := MultiSerializeDefaultCanvasRepo(&accessCanvasRepo)
	return canvasRepoViews, nil
}

/*
	Get all the canvases count based on studioId with user permissions.
	Args:
		body *AuthUserGetAllStudioCanvasControllerCount
	Returns:
		int
		error
*/
func (c canvasRepoController) AuthUserGetAllStudioCanvasControllerCount(user *models.User, studioId uint64) (int, error) {
	var canvasRepos *[]models.CanvasRepository
	var permissionsList map[uint64]map[uint64]string
	var err error
	canvasRepos, err = App.Repo.GetCanvasReposWithCollections(map[string]interface{}{"studio_id": studioId, "is_archived": false, "is_processing": false})
	if err != nil {
		return 0, err
	}
	canvasCount := 0
	for _, repo := range *canvasRepos {
		if repo.Collection.IsArchived {
			continue
		}
		if repo.IsPublished == false && repo.CreatedByID != user.ID {
			continue
		}
		if repo.ParentCanvasRepositoryID == nil {
			permissionsList, err = permissions.App.Service.CalculateCanvasRepoPermissions(user.ID, studioId, repo.CollectionID)
			if err != nil {
				continue
			}
		} else {
			permissionsList, err = permissions.App.Service.CalculateSubCanvasRepoPermissions(user.ID, studioId, repo.CollectionID, *repo.ParentCanvasRepositoryID)
			if err != nil {
				continue
			}
		}
		repoPermissions := permissionsList[repo.ID]
		permissionValues := utils.Values(repoPermissions)
		hasPermission := false
		for _, perm := range permissionValues {
			if utils.Contains(permissiongroup.UserAccessCanvasPermissionsList, perm) {
				canvasCount++
				hasPermission = true
				break
			}
		}
		if !hasPermission {
			if repo.HasPublicCanvas && repo.IsPublished {
				canvasCount++
			} else if repo.DefaultBranch.PublicAccess == models.EDIT && repo.IsPublished {
				canvasCount++
			} else if repo.DefaultBranch.PublicAccess == models.COMMENT && repo.IsPublished {
				canvasCount++
			} else if repo.DefaultBranch.PublicAccess == models.VIEW && repo.IsPublished {
				canvasCount++
			}
		}
	}

	if err != nil {
		return 0, err
	}
	//canvasRepoViews := MultiSerializeDefaultCanvasRepo(&accessCanvasRepo)
	return canvasCount, nil
}

func (c canvasRepoController) CreateCanvasRepo(
	body NewCanvasRepoPost, userID uint64, studioID uint64, user models.User, parentPublicAccess string) (*models.CanvasRepository, error) {
	// create repo
	repoInstance, errRepoCreate := queries.App.RepoQuery.CreateRepo(
		body.ParentCanvasRepositoryID,
		userID,
		studioID,
		body.CollectionID,
		body.Name,
		body.Icon,
		body.Position,
	)
	if errRepoCreate != nil {
		return nil, errRepoCreate
	}
	// update position of all other canvas
	if body.Position == 1 {
		App.Repo.UpdateCanvasPositionOnAddingNewCanvas(repoInstance)
	}
	// create branch
	branchInstance, errBranchCreate := queries.App.BranchQuery.CreateBranch(
		userID,
		repoInstance.ID,
		models.CANVAS_BRANCH_NAME_MAIN,
		true,
		parentPublicAccess,
	)
	if errBranchCreate != nil {
		return nil, errBranchCreate
	}

	_, branchRepoAttachError := queries.App.RepoQuery.UpdateRepo(repoInstance.ID, map[string]interface{}{
		"default_branch_id": branchInstance.ID,
	})

	// Now we need to set repoInstance.DefaultBranchID = branchInstance.ID save to Database.
	if branchRepoAttachError != nil {
		return nil, branchRepoAttachError
	}
	repo, _ := queries.App.RepoQuery.GetRepo(map[string]interface{}{"id": repoInstance.ID})

	err := queries.App.PermsQuery.CreateDefaultCanvasBranchPermission(repo.CollectionID, userID, studioID, repo.ID, branchInstance.ID, repo.ParentCanvasRepositoryID)
	if err != nil {
		return nil, err
	}

	// Creating a blank block
	var block models.Block

	firstContrib := App.Service.BlockContributorFirst(user, branchInstance.ID)
	newBlock, errCreatingBlockInst := block.NewBlock(
		uuid.New(),
		userID,
		repoInstance.ID,
		branchInstance.ID,
		nil,
		1,
		2,
		models.BlockTypeText,
		models.MyFirstBlockJson(),
		models.MyFirstEmptyBlockJson(),
		firstContrib,
	)
	if errCreatingBlockInst != nil {
		return nil, errCreatingBlockInst
	}
	errNewBlock := blocks.App.Service.Create(newBlock)
	if errNewBlock != nil {
		return nil, errNewBlock
	}

	return repo, nil
}

func (c canvasRepoController) RoleCanvasSearchController(studioID, roleID uint64, search string) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	repoAndCollectionRows := queries.App.CanvasRepoQuery.QueryDB(search, studioID)
	collectionIDs, repoIDs, supportRepoIDs := App.Service.ProcessSearchDump(repoAndCollectionRows)
	// get collections
	collections, err := App.Service.GetRoleCollectionsByID(collectionIDs, roleID)
	if err != nil {
		return nil, err
	}
	// get canvases
	canvasRepos, err := App.Service.GetRoleCanvasesByIDs(repoIDs, roleID)
	if err != nil {
		return nil, err
	}

	supportCanvasRepos, err := App.Service.GetRoleCanvasesByIDs(supportRepoIDs, roleID)
	sCanvasRepos := *supportCanvasRepos
	if err != nil {
		return nil, err
	}
	for i, _ := range sCanvasRepos {
		sCanvasRepos[i].SearchMatch = false
	}
	*canvasRepos = append(*canvasRepos, sCanvasRepos...)

	result["collections"] = collections
	result["canvases"] = canvasRepos
	return result, nil
}

func (c canvasRepoController) UserCanvasSearchController(studioID, userID uint64, search string) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	repoAndCollectionRows := queries.App.CanvasRepoQuery.QueryDB(search, studioID)
	collectionIDs, repoIDs, supportRepoIDs := App.Service.ProcessSearchDump(repoAndCollectionRows)
	// get collections
	collections, err := App.Service.GetStudioMemberCollectionsByIDs(collectionIDs, studioID, userID)
	if err != nil {
		return nil, err
	}
	// get canvases
	canvasRepos, err := App.Service.GetMemberCanvasesByIDs(studioID, userID, repoIDs)
	if err != nil {
		return nil, err
	}

	supportCanvasRepos, err := App.Service.GetMemberCanvasesByIDs(studioID, userID, supportRepoIDs)
	sCanvasRepos := *supportCanvasRepos
	if err != nil {
		return nil, err
	}
	for i, _ := range sCanvasRepos {
		sCanvasRepos[i].SearchMatch = false
	}
	*canvasRepos = append(*canvasRepos, sCanvasRepos...)

	result["collections"] = collections
	result["canvases"] = canvasRepos
	return result, nil
}
