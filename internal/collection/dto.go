package collection

import (
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"time"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gorm.io/gorm"
)

func (cr collectionRepo) GetCollection(query map[string]interface{}) (*models.Collection, error) {
	var collection models.Collection
	err := cr.db.Model(&models.Collection{}).Where(query).First(&collection).Error
	if err != nil {
		return nil, err
	}
	return &collection, nil
}

func (cr collectionRepo) CreateCollection(collection *models.Collection) error {
	err := cr.db.Create(collection).Error
	return err
}

func (cr collectionRepo) DeleteCollection(collectionId uint64, userId uint64) error {
	err := cr.db.Model(&models.Collection{}).Where("id = ?", collectionId).Updates(map[string]interface{}{
		"is_archived":    true,
		"archived_at":    time.Now(),
		"archived_by_id": userId,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (cr collectionRepo) GetCollections(query map[string]interface{}) (*[]models.Collection, error) {
	var collections []models.Collection
	err := cr.db.Model(&models.Collection{}).Where(query).Order("position ASC").Find(&collections).Error
	if err != nil {
		return nil, err
	}
	return &collections, nil
}

// Returns all the Permissions for this COllection
func (cr collectionRepo) GetCollectionPermission(query map[string]interface{}) (*models.CollectionPermission, error) {
	var collectionPerm *models.CollectionPermission
	postgres.GetDB().Model(&collectionPerm).Where(query).Preload("Role").Preload("Member").Preload("Member.User").Preload("Role.Members").First(&collectionPerm)
	return collectionPerm, nil
}

func (cr collectionRepo) GetCollectionsPermissionAll(query map[string]interface{}) (*[]models.CollectionPermission, error) {
	var collectionPerm *[]models.CollectionPermission
	postgres.GetDB().Model(&collectionPerm).Where(query).Preload("Role").Preload("Member").Preload("Member.User").Preload("Role.Members").Find(&collectionPerm)
	return collectionPerm, nil
}

// Get Member ID for this user id
func (cr collectionRepo) GetMemberByUserID(userId uint64, studioID uint64) (*models.Member, error) {
	var member *models.Member
	result := postgres.GetDB().Model(&models.Member{}).
		Where("user_id = ? and studio_id = ? AND has_left = false AND is_removed = false", userId, studioID).First(&member)

	if result.Error != nil {
		return nil, result.Error
	}
	return member, nil
}

func (cr collectionRepo) GetCanvasRepos(query map[string]interface{}) (*[]models.CanvasRepository, error) {
	var canvasRepos []models.CanvasRepository
	err := cr.db.Model(&models.CanvasRepository{}).Where(query).Order("position ASC").Find(&canvasRepos).Error
	if err != nil {
		return nil, err
	}
	return &canvasRepos, nil
}

func (cr collectionRepo) UpdateCanvasBranch(canvasBranchId uint64, updates map[string]interface{}) (*models.CanvasBranch, error) {
	var canvasBranch *models.CanvasBranch
	err := cr.db.Model(models.CanvasBranch{}).Where("id = ?", canvasBranchId).Updates(updates).Find(&canvasBranch).Error
	return canvasBranch, err
}

func (cr collectionRepo) AnonymousGetCollections(studioID uint64) (*[]models.Collection, error) {
	var collections []models.Collection
	err := cr.db.Model(&models.Collection{}).Where(
		"studio_id = ? AND (public_access != 'private' or has_public_canvas = true) AND is_archived = false",
		studioID).Order("position ASC").Find(&collections).Error
	if err != nil {
		return nil, err
	}
	return &collections, nil
}

/*
	Dto method of moving collection forward or backward. We added the logic based on gorm db condition
	to update the position of the collection.

	Logic Explanation:

	* Checking the future position with the current position so we can know to move forward or backward.

	    Forward move:
			Example: positions = 1,2,3,4,5
			Here if 2 wants to move to 4
			We are fetching collections based on studio_id, position > collectionCurrentPosition &
			position < futurePosition + 1(equals to position <= futurePosition)

			And we are decreasing by -1 for all these collections

			So positions will be = 1,2,2,3,5
			And finally we change 2 position to 4

		Backward move:
			Example: positions = 1,2,3,4,5
			Here if 4 wants to move to 2
			We are fetching collections based on studio_id, position > futurePosition -1 &
			position < collectionCurrentPosition

			And we are increasing by +1 for all these collections

			So position will be = 1,3,4,4,5
			And finally we change 4 position to 2

	* Finally setting up the futurePosition to the collection that needs to be moved.
		So again positions will be = 1,2,3,4,5

	Args:
		collection *models.Collection
		futurePosition uint
	Returns:
		error
*/
func (cr collectionRepo) MoveCollection(collection *models.Collection, futurePosition uint) error {
	var err error
	// Forward move
	if futurePosition > collection.Position {
		err := cr.db.Model(&models.Collection{}).Where(
			"studio_id = ? AND position > ? AND position < ?", collection.StudioID, collection.Position, futurePosition+1).Update("position", gorm.Expr("position - 1")).Error
		if err != nil {
			return err
		}
		// Backward move
	} else if futurePosition < collection.Position {
		err := cr.db.Model(&models.Collection{}).Where(
			"studio_id = ? AND position > ? AND position < ?", collection.StudioID, futurePosition-1, collection.Position).Update("position", gorm.Expr("position + 1")).Error
		if err != nil {
			return err
		}
	}
	// Setting up the moved collection position.
	err = cr.db.Model(&models.Collection{}).Where("id = ?", collection.ID).Update("position", futurePosition).Error
	if err != nil {
		return err
	}

	return nil
}

func (cr collectionRepo) ResetCollectionPositionOnDelete(collection *models.Collection) error {
	err := cr.db.Model(&models.Collection{}).Where(
		"studio_id = ? AND position > ?", collection.StudioID, collection.Position).Update("position", gorm.Expr("position - 1")).Error
	if err != nil {
		return err
	}
	return nil
}

func (cr collectionRepo) GetNextCollections(col *models.Collection) ([]models.Collection, error) {
	var collections []models.Collection
	err := cr.db.Model(models.Collection{}).Where(
		"position > ? and studio_id = ? and is_archived = false and (public_access <> ? or has_public_canvas = ?)", col.Position, col.StudioID, models.PRIVATE, true,
	).Order("position ASC").Find(&collections).Error
	return collections, err
}

func (cr collectionRepo) GetPrevCollections(col *models.Collection) ([]models.Collection, error) {
	var collections []models.Collection
	err := cr.db.Model(models.Collection{}).Where(
		"position < ? and studio_id = ? and is_archived = false and (public_access <> ? or has_public_canvas = ?)", col.Position, col.StudioID, models.PRIVATE, true,
	).Order("position DESC").Find(&collections).Error
	return collections, err
}
