package canvasbranch

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
)

type CanvasBranchRepoInterface interface {
	Get(query map[string]interface{}) (*models.CanvasBranch, error)
	Create(id uint64)
	Update(id uint64)
	Delete(id uint64)
}

// Should not be used.
// Replace with
//branch, err := queries.App.BranchQuery.GetBranchByID.GetBranchWithRepoAndStudio(id)
func (r canvasBranchRepo) Get(query map[string]interface{}) (*models.CanvasBranch, error) {
	var branch models.CanvasBranch
	err := r.db.Model(&models.CanvasBranch{}).Where(query).Preload("CanvasRepository").Preload("CanvasRepository.Studio").First(&branch).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &branch, nil
}

// GetBranchRepoCollection: Branch with CanvasRepo and Collection
func (r canvasBranchRepo) GetBranchRepoCollection(branchID uint64) (*models.CanvasBranch, error) {
	var branch models.CanvasBranch
	err := r.db.Model(&models.CanvasBranch{}).Where("id = ?", branchID).Preload("CanvasRepository").Preload("CanvasRepository.Collection").First(&branch).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &branch, nil
}

func (r canvasBranchRepo) SetBranchCommited(branchID uint64, commited bool) error {
	err := r.db.Model(&models.CanvasBranch{}).Where("id = ?", branchID).Update("committed", commited).Error
	if err != nil {
		logger.Debug(err.Error())
		return err
	}
	return nil
}

func (r canvasBranchRepo) GetBranchesAnonymous(canvasID uint64, publicAccess []string) (*[]models.CanvasBranch, error) {
	var branches []models.CanvasBranch
	err := r.db.Model(&models.CanvasBranch{}).Where("canvas_repository_id = ? AND public_access IN ? AND is_archived = false AND is_merged = false", canvasID, publicAccess).Find(&branches).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &branches, nil
}

func (r canvasBranchRepo) GetBranchesNoPreload(reposID []uint64) (*[]models.CanvasBranch, error) {
	var branches []models.CanvasBranch
	err := r.db.Model(&models.CanvasBranch{}).Where("canvas_repository_id IN ?", reposID).Find(&branches).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &branches, nil
}

func (r canvasBranchRepo) GetBranchesSimple(query map[string]interface{}) (*[]models.CanvasBranch, error) {
	var branches []models.CanvasBranch
	err := r.db.Model(&models.CanvasBranch{}).Where(query).Find(&branches).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &branches, nil
}
func (r canvasBranchRepo) GetBranches(query map[string]interface{}) (*[]models.CanvasBranch, error) {
	var branches []models.CanvasBranch
	err := r.db.Model(&models.CanvasBranch{}).Where(query).Preload("CanvasRepository").Find(&branches).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &branches, nil
}

func (r canvasBranchRepo) GetCanvasRepos(query map[string]interface{}) (*[]models.CanvasRepository, error) {
	var canvasRepos []models.CanvasRepository
	err := r.db.Model(&models.CanvasRepository{}).Where(query).Find(&canvasRepos).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &canvasRepos, nil
}

func (r canvasBranchRepo) Create(instance *models.CanvasBranch) (*models.CanvasBranch, error) {
	results := r.db.Create(&instance)
	return instance, results.Error
}

//type BlockCBRepoInterface interface {
//}

// Create a new branch with
// Note to SR to remove (IsRoughBranch = true) from NavList
// IsRoughBranch = true
// RoughFromBranchID = mainBranchId
// RoughBranchCreatorID = userID
func (r canvasBranchRepo) CreateRoughBranch(parentBranch *models.CanvasBranch, userID uint64, branchName, key, commitID string) (*models.CanvasBranch, error) {
	var roughBranch models.CanvasBranch
	// New Instance Filling
	roughBranch.IsRoughBranch = true
	roughBranch.RoughFromBranchID = &parentBranch.ID
	roughBranch.RoughBranchCreatorID = &userID
	roughBranch.CanvasRepositoryID = parentBranch.CanvasRepositoryID
	roughBranch.Name = branchName
	roughBranch.IsDraft = true
	roughBranch.IsMerged = false
	roughBranch.IsDefault = false
	roughBranch.CreatedByID = userID
	roughBranch.UpdatedByID = userID
	roughBranch.Key = key
	roughBranch.CreatedFromCommitID = commitID

	// Create Rough Branch
	branch, errCreating := r.Create(&roughBranch)
	if errCreating != nil {
		return nil, errCreating
	}

	return branch, nil
}

func (r canvasBranchRepo) DeleteBranch(branchID uint64) error {
	err := r.Manager.HardDeleteByID(models.CANVAS_BRANCH, branchID)
	if err != nil {
		return err
	}
	return nil
}

func (r canvasBranchRepo) UpdateBranchInstance(branchID uint64, query map[string]interface{}) error {
	err := r.db.Model(&models.CanvasBranch{}).Where("id = ?", branchID).Updates(query).Error
	if err != nil {
		return err
	}
	return nil
}

func (r canvasBranchRepo) PublishCanvasRepository(canvasRepoId uint64) error {
	err := r.db.Model(&models.CanvasRepository{}).Where("id = ?", canvasRepoId).Update("is_published", true).Error
	if err != nil {
		return err
	}
	return nil
}

func (r canvasBranchRepo) GetBranchWithRepo(query map[string]interface{}) (*models.CanvasBranch, error) {
	var branch models.CanvasBranch
	err := r.db.Model(&models.CanvasBranch{}).Where(query).Preload("CanvasRepository").First(&branch).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &branch, nil
}

// Canvas Repo DTO
func (r canvasBranchRepo) GetCanvasRepoInstance(query map[string]interface{}) (*models.CanvasRepository, error) {
	var repo models.CanvasRepository
	err := r.db.Model(&models.CanvasRepository{}).Where(query).First(&repo).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &repo, nil
}

func (r canvasBranchRepo) GetBranchRepoInstance(query map[string]interface{}) (*models.CanvasBranch, error) {
	var branch models.CanvasBranch
	err := r.db.Model(&models.CanvasBranch{}).Where(query).Preload("CanvasRepository").Preload("CanvasRepository.ParentCanvasRepository").First(&branch).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &branch, nil
}

//func (r canvasBranchRepo) GetAllRepos(studioID uint64) (*[]models.CanvasRepository, error) {
//
//}

//func (r canvasBranchRepo) GetRoughUnpublishedbyStudioandUser(studio uint64) (*[]models.CanvasBranch, error) {
//	var branches []models.CanvasBranch
//	err := r.db.Model(&models.CanvasBranch{}).Where("canvas_repository_id = ? AND public_access IN ? AND is_archived = false", canvasID, publicAccess).Find(&branches).Error
//	if err != nil {
//		logger.Debug(err.Error())
//		return nil, err
//	}
//	return &branches, nil
//}

func (r canvasBranchRepo) CanvasRepo(id uint64) models.CanvasRepository {
	var repo models.CanvasRepository
	r.db.Table(models.CANVAS_REPO).Where("id = ?", id).First(&repo)
	fmt.Println("RepoKey", repo.Key)
	return repo
}

func (r canvasBranchRepo) GetRoughUnpublishedByStudioAndUser(studioID uint64, userID uint64) (*[]models.CanvasBranch, error) {
	var branches []models.CanvasBranch
	queryErr := r.db.Raw("select canvas_branches.* from canvas_repositories left join canvas_branches ON canvas_branches.canvas_repository_id = canvas_repositories.id where studio_id = @studioid and canvas_branches.created_by_id = @userid  and (canvas_branches.is_draft is true or canvas_branches.is_rough_branch is true) and canvas_branches.is_merged is false and canvas_branches.is_archived is false and canvas_repositories.is_published = false",
		map[string]interface{}{"studioid": studioID, "userid": userID}).Find(&branches).Error
	if queryErr != nil {
		logger.Debug(queryErr.Error())
		return nil, queryErr
	}
	return &branches, nil
}

func (cr canvasBranchRepo) GetCollections(query map[string]interface{}) (*[]models.Collection, error) {
	var collections []models.Collection
	err := cr.db.Model(&models.Collection{}).Where(query).Order("position ASC").Find(&collections).Error
	if err != nil {
		return nil, err
	}
	return &collections, nil
}

func (cr canvasBranchRepo) SetBranchLastSyncedAllAttributionsCommitID(branchID uint64, lastSyncedAllAttributionsCommitID string) error {
	return cr.db.Model(&models.CanvasBranch{}).Where(models.CanvasBranch{BaseModel: models.BaseModel{ID: branchID}}).Update("last_synced_all_attributions_commit_id", lastSyncedAllAttributionsCommitID).Error
}

func (r canvasBranchRepo) GetListOfClonedFromThreadIDs(blockID uint64) []uint64 {
	var ids []uint64
	var cleaned []uint64
	r.db.Model(&models.BlockThread{}).Where("start_block_id = ?", blockID).Pluck("cloned_from_thread", &ids)
	for _, id := range ids {
		if id != 0 {
			cleaned = append(cleaned, id)
		}
	}
	return cleaned
}

func (r canvasBranchRepo) UpdateBranchLastEdited(branchID uint64) error {
	err := r.db.Model(models.CanvasBranch{}).Where("id = ?", branchID).Updates(map[string]interface{}{}).Error
	return err
}
