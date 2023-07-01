package queries

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

func (q collectionQuery) UpdateCollection(collectionID uint64, updates map[string]interface{}) (*models.Collection, error) {
	var collection *models.Collection
	err := postgres.GetDB().Model(models.Collection{}).Where(
		models.Collection{BaseModel: models.BaseModel{ID: collectionID}}).Updates(updates).Find(&collection).Error
	return collection, err
}

func (q collectionQuery) GetCollections(query map[string]interface{}) ([]models.Collection, error) {
	var collections []models.Collection
	err := postgres.GetDB().Model(models.Collection{}).Where(query).Order("position ASC").Find(&collections).Error
	return collections, err
}

func (q collectionQuery) GetCollection(query map[string]interface{}) (*models.Collection, error) {
	var collections *models.Collection
	err := postgres.GetDB().Model(models.Collection{}).Where(query).Order("position ASC").First(&collections).Error
	return collections, err
}

func (q collectionQuery) SearchCollections(studioID uint64, query string) ([]models.Collection, error) {
	var collections []models.Collection
	err := postgres.GetDB().Model(models.Collection{}).Where("studio_id = ? and name ILIKE ?  and is_archived = false", studioID, "%"+query+"%").Order("position ASC").Find(&collections).Error
	return collections, err
}
