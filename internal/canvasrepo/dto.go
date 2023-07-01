package canvasrepo

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gorm.io/gorm"
)

type RepoInterface interface {
	Create(id uint64)
	CreateCanvasBranch(instance *models.CanvasBranch) (*models.CanvasBranch, error)
}

func (r canvasRepoRepo) Get(query map[string]interface{}) (*models.CanvasRepository, error) {
	var repo models.CanvasRepository
	err := r.db.Model(&models.CanvasRepository{}).Where(query).First(&repo).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &repo, nil
}

func (r canvasRepoRepo) Update(query, updates map[string]interface{}) (*models.CanvasRepository, error) {
	var repo models.CanvasRepository
	err := r.db.Model(&models.CanvasRepository{}).Where(query).Updates(updates).First(&repo).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &repo, nil
}

// get Repo

// get Repo
func (r canvasRepoRepo) GetMultipleReposByIDS(IDs []uint64) (*[]models.CanvasRepository, error) {
	var repos []models.CanvasRepository
	err := r.db.Model(&models.CanvasRepository{}).Where("id IN ?", IDs).Preload("DefaultBranch").Find(&repos).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &repos, nil
}

// Canvas Repo
func (r canvasRepoRepo) Create(instance *models.CanvasRepository) (*models.CanvasRepository, error) {
	results := r.db.Create(&instance)
	go func() {
		canvasRepo, _ := json.Marshal(instance)
		r.kafka.Publish(configs.KAFKA_TOPICS_NEW_CANVAS, strconv.FormatUint(instance.ID, 10), canvasRepo)
	}()
	return instance, results.Error
}

// Canvas Repo
func (r canvasRepoRepo) CreateMultiple(instances *[]models.CanvasRepository) error {
	results := r.db.Create(instances)
	go func() {
		for _, instance := range *instances {
			canvasRepo, _ := json.Marshal(instance)
			r.kafka.Publish(configs.KAFKA_TOPICS_NEW_CANVAS, strconv.FormatUint(instance.ID, 10), canvasRepo)
		}
	}()
	return results.Error
}

// Canvas Branch Create
func (r canvasRepoRepo) CreateCanvasBranch(instance *models.CanvasBranch) (*models.CanvasBranch, error) {
	results := r.db.Create(&instance)
	return instance, results.Error
}

// Canvas Branch Create
func (r canvasRepoRepo) CreateCanvasBranchMultiple(instances *[]models.CanvasBranch) error {
	results := r.db.Create(instances)
	return results.Error
}

// Canvas Branch Create
func (r canvasRepoRepo) UpdateCanvasRepoDefaultBranch(repoID uint64, branchId uint64) (*models.CanvasRepository, error) {
	var repo models.CanvasRepository
	err := r.db.Model(&models.CanvasRepository{}).Where("id = ?", repoID).Update("default_branch_id", branchId).First(&repo).Error
	return &repo, err
}

// get Repo
func (r canvasRepoRepo) UpdateRepoNameIconOnly(name string, icon string, repoId uint64, userId uint64, cover string) error {
	var err error
	// If cover url is empty don't update

	//if cover == "" {
	//	err = r.db.Model(&models.CanvasRepository{}).Where("id = ?", repoId).Updates(map[string]interface{}{
	//		"name":          name,
	//		"icon":          icon,
	//		"updated_by_id": userId,
	//	}).Error
	//
	//} else {
	//	err = r.db.Model(&models.CanvasRepository{}).Where("id = ?", repoId).Updates(map[string]interface{}{
	//		"name":          name,
	//		"icon":          icon,
	//		"updated_by_id": userId,
	//		"cover_url":     cover,
	//	}).Error
	//}

	err = r.db.Model(&models.CanvasRepository{}).Where("id = ?", repoId).Updates(map[string]interface{}{
		"name":          name,
		"icon":          icon,
		"updated_by_id": userId,
		"cover_url":     cover,
	}).Error

	if err != nil {
		return err
	}
	return nil
}

// GetAnonymousCanvasRepos based on different filters and by order by position ASC.
func (r canvasRepoRepo) GetAnonymousCanvasRepos(collectionId uint64, publicAccess []string) (*[]models.CanvasRepository, error) {
	var repos []models.CanvasRepository
	err := r.db.Table("canvas_repositories").
		Joins("LEFT JOIN canvas_branches ON canvas_branches.canvas_repository_id = canvas_repositories.id").
		Where("collection_id = ? AND canvas_repositories.is_archived = false AND parent_canvas_repository_id is null AND (canvas_branches.public_access in ? or canvas_repositories.has_public_canvas = true) AND is_published = true AND is_processing = false", collectionId, publicAccess).
		Order("position ASC").
		Group("canvas_repositories.id").
		Preload("DefaultBranch").
		Find(&repos).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &repos, nil
}

// GetAnonymousSubCanvasRepos based on different filters and by order by position ASC.
func (r canvasRepoRepo) GetAnonymousSubCanvasRepos(parentCanvasRepoId uint64, publicAccess []string) (*[]models.CanvasRepository, error) {
	var repos []models.CanvasRepository
	err := r.db.Table("canvas_repositories").
		Joins("LEFT JOIN canvas_branches ON canvas_branches.canvas_repository_id = canvas_repositories.id").
		Where("parent_canvas_repository_id = ? AND canvas_repositories.is_archived = false AND (canvas_branches.public_access in ? or canvas_repositories.has_public_canvas = true) AND is_published = true AND is_processing = false", parentCanvasRepoId, publicAccess).
		Order("position ASC").
		Group("canvas_repositories.id").
		Preload("DefaultBranch").
		Find(&repos).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &repos, nil
}

func (r canvasRepoRepo) GetCanvasRepos(query map[string]interface{}) (*[]models.CanvasRepository, error) {
	var repo []models.CanvasRepository
	err := r.db.Model(&models.CanvasRepository{}).Where(query).Preload("DefaultBranch").Order("position ASC").Find(&repo).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &repo, nil
}

func (r canvasRepoRepo) GetCanvasRepoByKey(key string) (*models.CanvasRepository, error) {
	var repo models.CanvasRepository
	err := r.db.Model(&models.CanvasRepository{}).Where("LOWER(key) = ?", strings.ToLower(key)).Preload("DefaultBranch").Preload("ParentCanvasRepository").First(&repo).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &repo, nil
}

/*
	Dto method of moving canvas forward or backward. We added the logic based on gorm db condition
	to update the position of the canvas.

	Moves the canvas which are on same collection and are having parent as collection

	Logic Explanation:
	This logic is same for all the below methods.
		- MoveCanvasRepoWithSameParentCollection
		- MoveCanvasRepoWithDifferentParentCollection
		- MoveCanvasRepoWithSameParentCanvas
		- MoveCanvasRepoWithDifferentParentAsCanvas
		- MoveCanvasRepoWithDifferentParentCanvasAndCollection


	* Checking the future position with the current position so we can know to move forward or backward.

	    Forward move:
			Example: positions = 1,2,3,4,5
			Here if 2 wants to move to 4
			We are fetching collections based on collectionId or parentCanvasId on different menthods &
			position > collectionCurrentPosition &
			position < futurePosition + 1(equals to position <= futurePosition)

			And we are decreasing by -1 for all these collections

			So positions will be = 1,2,2,3,5
			And finally we change 2 position to 4

		Backward move:
			Example: positions = 1,2,3,4,5
			Here if 4 wants to move to 2
			We are fetching collections based on collectionId or parentCanvasId on different menthods &
			position > futurePosition -1 &
			position < collectionCurrentPosition

			And we are increasing by +1 for all these collections

			So position will be = 1,3,4,4,5
			And finally we change 4 position to 2

	* Finally setting up the futurePosition to the collection that needs to be moved.
		So again positions will be = 1,2,3,4,5

	Args:
		canvas *models.CanvasRepository
		futurePosition uint
	Returns:
		error
*/
func (r canvasRepoRepo) MoveCanvasRepoWithParentAsCollection(canvas *models.CanvasRepository, futurePosition uint) error {
	var err error
	// Forward move
	if canvas.ParentCanvasRepositoryID != nil {
		err = r.db.Model(&models.CanvasRepository{}).Where(
			"collection_id = ? AND position > ? AND parent_canvas_repository_id is null",
			canvas.CollectionID, futurePosition-1).Update("position", gorm.Expr("position + 1")).Error
		r.RearrangeTheOldCanvasRepo(canvas)
	} else if futurePosition > canvas.Position {
		err = r.db.Model(&models.CanvasRepository{}).Where(
			"collection_id = ? AND position > ? AND position < ? AND parent_canvas_repository_id is null",
			canvas.CollectionID, canvas.Position, futurePosition+1).Update("position", gorm.Expr("position - 1")).Error
		if err != nil {
			return err
		}
		// Backward move
	} else if futurePosition < canvas.Position {
		err = r.db.Model(&models.CanvasRepository{}).Where(
			"collection_id = ? AND position > ? AND position < ? AND parent_canvas_repository_id is null",
			canvas.CollectionID, futurePosition-1, canvas.Position).Update("position", gorm.Expr("position + 1")).Error
		if err != nil {
			return err
		}
	}
	// Setting up the moved canvas position.
	err = r.db.Model(&models.CanvasRepository{}).Where("id = ?", canvas.ID).Updates(map[string]interface{}{
		"position":                    futurePosition,
		"parent_canvas_repository_id": nil,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

/*
	Moves the canvas which are on different collection and are having parent as collection

	* After setting up the position for the required canvas we are re-arranging the old canvas positions.

	Args:
		canvas *models.CanvasRepository
		futurePosition uint
		futureParentCollectionID uint64
	Returns:
		error
*/
func (r canvasRepoRepo) MoveCanvasRepoWithDifferentParentCollection(canvas *models.CanvasRepository, futurePosition uint, futureParentCollectionID uint64) error {
	var err error
	// Forward move
	//if futurePosition > canvas.Position {
	//	err = r.db.Model(&models.CanvasRepository{}).Where(
	//		"collection_id = ? AND position > ? AND position < ? AND parent_canvas_repository_id is null",
	//		futureParentCollectionID, canvas.Position, futurePosition+1).Update("position", gorm.Expr("position - 1")).Error
	//	// Backward move
	//} else if futurePosition < canvas.Position {
	//	err = r.db.Model(&models.CanvasRepository{}).Where(
	//		"collection_id = ? AND position > ? AND position < ? AND parent_canvas_repository_id is null",
	//		futureParentCollectionID, futurePosition-1, canvas.Position).Update("position", gorm.Expr("position + 1")).Error
	//}
	err = r.db.Model(&models.CanvasRepository{}).Where(
		"collection_id = ? AND position > ? AND parent_canvas_repository_id is null",
		futureParentCollectionID, futurePosition-1).Update("position", gorm.Expr("position + 1")).Error
	if err != nil {
		return err
	}

	r.RearrangeTheOldCanvasRepo(canvas)

	r.UpdateCollectionCountOnRemoveCanvas(canvas.CollectionID)

	if canvas.ParentCanvasRepositoryID != nil {
		r.UpdateParentCanvasCountOnRemoveCanvas(*canvas.ParentCanvasRepositoryID)
	}

	// Setting up the moved canvas position.
	err = r.db.Model(&models.CanvasRepository{}).Where("id = ?", canvas.ID).Updates(
		map[string]interface{}{"position": futurePosition, "collection_id": futureParentCollectionID, "parent_canvas_repository_id": nil}).Error
	if err != nil {
		return err
	}

	r.UpdateCollectionCountOnAddingCanvas(futureParentCollectionID)
	return nil
}

/*
	@todo Remove this method later
	Moves the canvas which are on same canvas and are having parent as canvas.

	Args:
		canvas *models.CanvasRepository
		futurePosition uint
	Returns:
		error
*/
// func (r canvasRepoRepo) MoveCanvasRepoWithSameParentCanvas(canvas *models.CanvasRepository, futurePosition uint) error {
// 	var err error
// 	// Forward move
// 	if futurePosition > canvas.Position {
// 		err = r.db.Model(&models.CanvasRepository{}).Where(
// 			"collection_id = ? AND parent_canvas_repository_id is null AND position > ? AND position < ?",
// 			canvas.CollectionID, canvas.Position, futurePosition+1).Update("position", gorm.Expr("position - 1")).Error
// 		// Backward move
// 	} else if futurePosition < canvas.Position {
// 		err = r.db.Model(&models.CanvasRepository{}).Where(
// 			"collection_id = ? AND parent_canvas_repository_id is null AND position > ? AND position < ?",
// 			canvas.CollectionID, futurePosition-1, canvas.Position).Update("position", gorm.Expr("position + 1")).Error
// 	}
// 	if err != nil {
// 		return err
// 	}

// 	// Setting up the moved canvas position.
// 	err = r.db.Model(&models.CanvasRepository{}).Where("id = ?", canvas.ID).Update("position", futurePosition).Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

/*
	Moves the canvas which are on different canvas but same collection and are having parent as canvas.

	* After setting up the position for the required canvas we are re-arranging the old canvas positions.

	Args:
		canvas *models.CanvasRepository
		futurePosition uint
		futureParentCanvasRepositoryID uint64
	Returns:
		error
*/
func (r canvasRepoRepo) MoveCanvasRepoWithDifferentParentAsCanvas(canvas *models.CanvasRepository, futurePosition uint, futureParentCanvasRepositoryID uint64) error {
	var err error
	// Forward move
	if canvas.ParentCanvasRepositoryID != nil && *canvas.ParentCanvasRepositoryID == futureParentCanvasRepositoryID {
		if futurePosition > canvas.Position {
			err = r.db.Model(&models.CanvasRepository{}).Where(
				"collection_id = ? AND parent_canvas_repository_id = ? AND position > ? AND position < ?",
				canvas.CollectionID, futureParentCanvasRepositoryID, canvas.Position, futurePosition+1).Update("position", gorm.Expr("position - 1")).Error
			// Backward move
		} else if futurePosition < canvas.Position {
			err = r.db.Model(&models.CanvasRepository{}).Where(
				"collection_id = ? AND parent_canvas_repository_id = ? AND position > ? AND position < ?",
				canvas.CollectionID, futureParentCanvasRepositoryID, futurePosition-1, canvas.Position).Update("position", gorm.Expr("position + 1")).Error
		}
	} else {
		err = r.db.Model(&models.CanvasRepository{}).Where(
			"collection_id = ? AND parent_canvas_repository_id = ? AND position > ?",
			canvas.CollectionID, futureParentCanvasRepositoryID, futurePosition-1).Update("position", gorm.Expr("position + 1")).Error
		r.RearrangeTheOldCanvasRepo(canvas)
		if canvas.ParentCanvasRepositoryID != nil {
			r.UpdateParentCanvasCountOnRemoveCanvas(*canvas.ParentCanvasRepositoryID)
		}
	}
	if err != nil {
		return err
	}

	// Setting up the moved canvas position.
	err = r.db.Model(&models.CanvasRepository{}).Where(
		"id = ?", canvas.ID).Updates(map[string]interface{}{"position": futurePosition, "parent_canvas_repository_id": futureParentCanvasRepositoryID}).Error
	if err != nil {
		return err
	}
	r.UpdateParentCanvasCountOnAddingCanvas(futureParentCanvasRepositoryID)
	return nil
}

/*
	Moves the canvas which are on different canvas and different collection and are having parent as canvas.

	* After setting up the position for the required canvas we are re-arranging the old canvas positions.

	Args:
		canvas *models.CanvasRepository
		futurePosition uint
		futureCollectionID uint64
		futureParentCanvasRepositoryID uint64
	Returns:
		error
*/
func (r canvasRepoRepo) MoveCanvasRepoWithDifferentParentCanvasAndCollection(canvas *models.CanvasRepository, futurePosition uint, futureCollectionID uint64, futureParentCanvasRepositoryID uint64) error {
	var err error
	// Forward move
	//if futurePosition > canvas.Position {
	//	err = r.db.Model(&models.CanvasRepository{}).Where(
	//		"collection_id = ? AND position > ? AND position < ? AND parent_canvas_repository_id = ?",
	//		futureCollectionID, canvas.Position, futurePosition+1, futureParentCanvasRepositoryID).Update("position", gorm.Expr("position - 1")).Error
	//	// Backward move
	//} else if futurePosition < canvas.Position {
	//	err = r.db.Model(&models.CanvasRepository{}).Where(
	//		"collection_id = ? AND position > ? AND position < ? AND parent_canvas_repository_id = ?",
	//		futureCollectionID, futurePosition-1, canvas.Position, futureParentCanvasRepositoryID).Update("position", gorm.Expr("position + 1")).Error
	//}
	err = r.db.Model(&models.CanvasRepository{}).Where(
		"collection_id = ? AND position > ? AND parent_canvas_repository_id = ?",
		futureCollectionID, futurePosition-1, futureParentCanvasRepositoryID).Update("position", gorm.Expr("position + 1")).Error
	if err != nil {
		return err
	}

	r.RearrangeTheOldCanvasRepo(canvas)

	r.UpdateCollectionCountOnRemoveCanvas(canvas.CollectionID)

	if canvas.ParentCanvasRepositoryID != nil {
		r.UpdateParentCanvasCountOnRemoveCanvas(*canvas.ParentCanvasRepositoryID)
	}

	// Setting up the moved canvas position.
	err = r.db.Model(&models.CanvasRepository{}).Where(
		"id = ?", canvas.ID).Updates(map[string]interface{}{
		"position":                    futurePosition,
		"parent_canvas_repository_id": futureParentCanvasRepositoryID,
		"collection_id":               futureCollectionID,
	}).Error
	r.UpdateCollectionCountOnAddingCanvas(futureCollectionID)
	r.UpdateParentCanvasCountOnAddingCanvas(futureParentCanvasRepositoryID)
	if err != nil {
		return err
	}
	return nil
}

/*
	Re-arranges the canvas repos & collections when a canvas is moved to different collection or
	to different parent canvas ID.

	Args:
		canvas *models.CanvasRepository
	Returns:
		error
*/
func (r canvasRepoRepo) RearrangeTheOldCanvasRepo(canvas *models.CanvasRepository) error {
	var err error
	if canvas.ParentCanvasRepositoryID == nil {
		err = r.db.Model(&models.CanvasRepository{}).Where(
			"collection_id = ? AND position > ? AND parent_canvas_repository_id is null",
			canvas.CollectionID, canvas.Position).Update("position", gorm.Expr("position - 1")).Error
	} else {
		err = r.db.Model(&models.CanvasRepository{}).Where(
			"collection_id = ? AND position > ? AND parent_canvas_repository_id = ?",
			canvas.CollectionID, canvas.Position, canvas.ParentCanvasRepositoryID).Update("position", gorm.Expr("position - 1")).Error
	}
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

func (r canvasRepoRepo) checkHasPublishRequestOnDefaultBranch(instance *models.CanvasBranch) bool {
	var count int64
	_ = r.db.Model(models.PublishRequest{}).Where("canvas_branch_id = ?", instance.ID).Count(&count)
	if count > 0 {
		return true
	}
	return false
}

// Canvas Branch Create
func (r canvasRepoRepo) UpdateCanvasRepoBranchPermissions(query, updates map[string]interface{}) ([]models.CanvasBranchPermission, error) {
	var repo []models.CanvasBranchPermission
	err := r.db.Model(&models.CanvasBranchPermission{}).Where(query).Updates(updates).Preload("Member").First(&repo).Error
	return repo, err
}

func (r canvasRepoRepo) UpdateCollectionCountOnRemoveCanvas(collectionID uint64) error {
	err := r.db.Model(&models.Collection{}).
		Where("id = ?", collectionID).
		Update("computed_root_canvas_count", gorm.Expr("computed_root_canvas_count - 1")).
		Update("computed_all_canvas_count", gorm.Expr("computed_all_canvas_count - 1")).
		Error
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (r canvasRepoRepo) UpdateCollectionCountOnAddingCanvas(collectionID uint64) error {
	err := r.db.Model(&models.Collection{}).
		Where("id = ?", collectionID).
		Update("computed_root_canvas_count", gorm.Expr("computed_root_canvas_count + 1")).
		Update("computed_all_canvas_count", gorm.Expr("computed_all_canvas_count + 1")).
		Error
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (r canvasRepoRepo) UpdateParentCanvasCountOnRemoveCanvas(canvasID uint64) error {
	err := r.db.Model(&models.CanvasRepository{}).
		Where("id = ?", canvasID).
		Update("sub_canvas_count", gorm.Expr("sub_canvas_count - 1")).
		Error
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (r canvasRepoRepo) UpdateParentCanvasCountOnAddingCanvas(canvasID uint64) error {
	err := r.db.Model(&models.CanvasRepository{}).
		Where("id = ?", canvasID).
		Update("sub_canvas_count", gorm.Expr("sub_canvas_count + 1")).
		Error
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (r canvasRepoRepo) GetCanvasReposCount(query map[string]interface{}) (int64, error) {
	var count int64
	err := r.db.Model(&models.CanvasRepository{}).Where(query).Count(&count).Error
	if err != nil {
		logger.Debug(err.Error())
		return 0, err
	}
	return count, nil
}

func (r canvasRepoRepo) GetCanvasBranchPermissions(query map[string]interface{}) ([]models.CanvasBranchPermission, error) {
	var canvasBranchPerms []models.CanvasBranchPermission
	postgres.GetDB().Model(&canvasBranchPerms).Where(query).Preload("Member").Find(&canvasBranchPerms)
	return canvasBranchPerms, nil
}

func (r canvasRepoRepo) AcceesRequestExistsSimple(branchID uint64, userID uint64) bool {
	var count int64
	_ = r.db.Model(&models.AccessRequest{}).Where("canvas_branch_id = ? and created_by_id = ? and status = ?", branchID, userID, models.ACCESS_REQUEST_PENDING).Count(&count).Error
	if count == 0 {
		return false
	}
	return true
}

func (r canvasRepoRepo) GetCollection(query map[string]interface{}) (*models.Collection, error) {
	var collection *models.Collection
	err := r.db.Model(&models.Collection{}).Where(query).First(&collection).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return collection, nil
}

func (r canvasRepoRepo) GetCollectionPreloadStudio(query map[string]interface{}) (*models.Collection, error) {
	var collection *models.Collection
	err := r.db.Model(&models.Collection{}).Where(query).Preload("Studio").First(&collection).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return collection, nil
}

func (r canvasRepoRepo) MergeRequestCount(canvasRepoID uint64) int64 {
	var count int64
	_ = r.db.Model(&models.MergeRequest{}).Where("canvas_repository_id = ? and status = ?", canvasRepoID, models.MERGE_REQUEST_OPEN).Count(&count).Error
	if count == 0 {
		return 0
	}
	return count
}

func (r canvasRepoRepo) GetAcceptedPublishRequestByBranch(canvasRepoID uint64) (*models.PublishRequest, error) {
	var pr *models.PublishRequest
	err := r.db.Model(models.PublishRequest{}).Where(map[string]interface{}{"canvas_repository_id": canvasRepoID, "status": models.PUBLISH_REQUEST_ACCEPTED}).First(&pr).Error
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func (r canvasRepoRepo) GetMemberByUserID(userId uint64, studioID uint64) (*models.Member, error) {
	var member *models.Member
	result := postgres.GetDB().Model(&models.Member{}).
		Where("user_id = ? and studio_id = ? AND has_left = false AND is_removed = false", userId, studioID).First(&member)

	if result.Error != nil {
		return nil, result.Error
	}
	return member, nil
}

func (r canvasRepoRepo) GetBranchPermission(query map[string]interface{}) (*models.CanvasBranchPermission, error) {
	var branchPerm *models.CanvasBranchPermission
	postgres.GetDB().Model(&branchPerm).Where(query).Preload("Role").Preload("Member").Preload("Member.User").Preload("Role.Members").Find(&branchPerm)
	return branchPerm, nil
}

func (r canvasRepoRepo) GetBranchesPermissionAll(query map[string]interface{}) (*[]models.CanvasBranchPermission, error) {
	var branchPerms *[]models.CanvasBranchPermission
	postgres.GetDB().Model(&branchPerms).Where(query).Preload("Role").Preload("Member").Preload("Member.User").Preload("Role.Members").Find(&branchPerms)
	return branchPerms, nil
}

func (r canvasRepoRepo) UpdateCanvasPositionOnAddingNewCanvas(canvas *models.CanvasRepository) error {
	var err error
	if canvas.ParentCanvasRepositoryID == nil {
		err = r.db.Model(&models.CanvasRepository{}).Where(
			"collection_id = ? AND parent_canvas_repository_id is null AND id <> ?",
			canvas.CollectionID, canvas.ID).Update("position", gorm.Expr("position + 1")).Error
	} else {
		err = r.db.Model(&models.CanvasRepository{}).Where(
			"collection_id = ? AND parent_canvas_repository_id = ? AND id <> ?",
			canvas.CollectionID, canvas.ParentCanvasRepositoryID, canvas.ID).Update("position", gorm.Expr("position + 1")).Error
	}
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

func (r canvasRepoRepo) GetRepoWithCollection(query map[string]interface{}) (*models.CanvasRepository, error) {
	var repo models.CanvasRepository
	err := r.db.Model(&models.CanvasRepository{}).Where(query).Preload("Collection").Preload("ParentCanvasRepository").First(&repo).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &repo, nil
}

func (r canvasRepoRepo) GetCanvasesOrderByPositionAsc(query map[string]interface{}) ([]models.CanvasRepository, error) {
	var repos []models.CanvasRepository
	err := r.db.Model(models.CanvasRepository{}).Where(query).Order("position ASC").Preload("DefaultBranch").Find(&repos).Error
	if err != nil {
		fmt.Println("Error while getting canvas by position ASC", err)
		return nil, err
	}
	return repos, nil
}

func (r canvasRepoRepo) GetNextCanvases(canvasRepository *models.CanvasRepository, position uint) ([]models.CanvasRepository, error) {
	var repos []models.CanvasRepository
	var err error
	if canvasRepository.ParentCanvasRepositoryID == nil {
		err = r.db.Model(models.CanvasRepository{}).Where("position > ? and collection_id = ? and parent_canvas_repository_id is null and is_archived = false and is_published = true and is_language_canvas = false", position, canvasRepository.CollectionID).Preload("DefaultBranch").Order("position ASC").Find(&repos).Error
	} else {
		err = r.db.Model(models.CanvasRepository{}).Where("position > ? and collection_id = ? and parent_canvas_repository_id = ? and is_archived = false and is_published = true and is_language_canvas = false", position, canvasRepository.CollectionID, *canvasRepository.ParentCanvasRepositoryID).Preload("DefaultBranch").Order("position ASC").Find(&repos).Error
	}
	if err != nil {
		fmt.Println("Error while getting canvas by position ASC", err)
		return nil, err
	}
	return repos, nil
}

func (r canvasRepoRepo) GetNextSubCanvases(canvasRepository *models.CanvasRepository, position uint) ([]models.CanvasRepository, error) {
	var repos []models.CanvasRepository
	var err error
	err = r.db.Model(models.CanvasRepository{}).Where("position > ? and collection_id = ? and parent_canvas_repository_id = ? and is_archived = false and is_published = true and is_language_canvas = false", position, canvasRepository.CollectionID, canvasRepository.ID).Preload("DefaultBranch").Order("position ASC").Find(&repos).Error
	if err != nil {
		fmt.Println("Error while getting canvas by position ASC", err)
		return nil, err
	}
	return repos, nil
}

func (r canvasRepoRepo) GetNextCollections(canvasRepository *models.CanvasRepository) ([]models.Collection, error) {
	var collections []models.Collection
	err := r.db.Model(models.Collection{}).Where("position > ? and studio_id = ? and is_archived = false", canvasRepository.Collection.Position, canvasRepository.StudioID).Order("position ASC").Find(&collections).Error
	if err != nil {
		fmt.Println("Error while getting collection by position ASC", err)
		return nil, err
	}
	return collections, nil
}

func (r canvasRepoRepo) GetCanvasesOrderByPositionDesc(query map[string]interface{}) ([]models.CanvasRepository, error) {
	var repos []models.CanvasRepository
	err := r.db.Model(models.CanvasRepository{}).Where(query).Order("position Desc").Preload("DefaultBranch").Find(&repos).Error
	if err != nil {
		fmt.Println("Error while getting canvas by position ASC", err)
		return nil, err
	}
	return repos, nil
}

func (r canvasRepoRepo) GetPrevCanvases(canvasRepository *models.CanvasRepository, position uint) ([]models.CanvasRepository, error) {
	var repos []models.CanvasRepository
	var err error
	if canvasRepository.ParentCanvasRepositoryID == nil {
		err = r.db.Model(models.CanvasRepository{}).Where("position < ? and collection_id = ? and parent_canvas_repository_id is null and is_archived = false and is_published = true and is_language_canvas = false", position, canvasRepository.CollectionID).Preload("ParentCanvasRepository").Preload("DefaultBranch").Order("position DESC").Find(&repos).Error
	} else {
		err = r.db.Model(models.CanvasRepository{}).Where("position < ? and collection_id = ? and parent_canvas_repository_id = ? and is_archived = false and is_published = true and is_language_canvas = false", position, canvasRepository.CollectionID, *canvasRepository.ParentCanvasRepositoryID).Preload("ParentCanvasRepository").Preload("DefaultBranch").Order("position DESC").Find(&repos).Error
	}
	if err != nil {
		fmt.Println("Error while getting canvas by position ASC", err)
		return nil, err
	}
	return repos, nil
}

func (r canvasRepoRepo) GetPrevCollections(canvasRepository *models.CanvasRepository) ([]models.Collection, error) {
	var collections []models.Collection
	err := r.db.Model(models.Collection{}).Where("position < ? and studio_id = ? and is_archived = false", canvasRepository.Collection.Position, canvasRepository.StudioID).Order("position DESC").Find(&collections).Error
	if err != nil {
		fmt.Println("Error while getting collection by position ASC", err)
		return nil, err
	}
	return collections, nil
}

func (r canvasRepoRepo) GetCanvasReposWithCollections(query map[string]interface{}) (*[]models.CanvasRepository, error) {
	var repo []models.CanvasRepository
	err := r.db.Model(&models.CanvasRepository{}).Where(query).Preload("DefaultBranch").Preload("Collection").Order("position ASC").Find(&repo).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &repo, nil
}

func (r canvasRepoRepo) StudioIntegrationUpdate(integrationID uint64, updates map[string]interface{}) error {
	err := r.db.Model(&models.StudioIntegration{}).Where("id = ?", integrationID).Updates(updates).Error
	if err != nil {
		return err
	}
	return nil
}

func (r canvasRepoRepo) GetStudioIntegration(studioId uint64, integrationType string) (integration *models.StudioIntegration, err error) {
	err = postgres.GetDB().Model(&models.StudioIntegration{}).Where("studio_id = ? and type = ?", studioId, integrationType).First(&integration).Error
	return
}

func (r canvasRepoRepo) AnonymousGetCollections(studioID uint64) (*[]models.Collection, error) {
	var collections []models.Collection
	err := r.db.Model(&models.Collection{}).Where(
		"studio_id = ? AND (public_access != 'private' or has_public_canvas = true) AND is_archived = false",
		studioID).Order("position ASC").Find(&collections).Error
	if err != nil {
		return nil, err
	}
	return &collections, nil
}

func (r canvasRepoRepo) GetLangCanvases(canvasRepoID uint64, language string) ([]models.CanvasRepository, error) {
	var langCanvases []models.CanvasRepository
	err := r.db.Model(models.CanvasRepository{}).Where("default_language_canvas_repo_id = ? and language = ? and is_archived = false", canvasRepoID, language).Find(&langCanvases).Error
	return langCanvases, err
}
