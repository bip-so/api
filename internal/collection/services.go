package collection

import (
	"strings"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

func (s collectionService) CollectionInstance(collection *CollectionCreateValidator, userID uint64, studioID uint64) models.Collection {
	newCollection := models.Collection{
		Name:         collection.Name,
		Position:     collection.Position,
		Icon:         collection.Icon,
		StudioID:     studioID,
		PublicAccess: strings.ToLower(collection.PublicAccess),
		CreatedByID:  userID,
		UpdatedByID:  userID,
	}

	if collection.ParentCollectionID != 0 {
		newCollection.ParentCollectionID = &collection.ParentCollectionID
	}

	return newCollection
}

func (s collectionService) Update(requestBody *CollectionUpdateValidator) map[string]interface{} {
	updates := map[string]interface{}{
		"name":          requestBody.Name,
		"icon":          requestBody.Icon,
		"public_access": requestBody.PublicAccess,
	}
	return updates
}

func (s collectionService) GetCollectionPrevAndNext(userID, collectionID uint64) (*models.Collection, *models.Collection) {
	var prevCollection models.Collection
	var nextCollection models.Collection
	collectionInstance, _ := App.Repo.GetCollection(map[string]interface{}{"id": collectionID})
	nextCollections, _ := App.Repo.GetNextCollections(collectionInstance)
	if len(nextCollections) > 0 {
		nextCollection = nextCollections[0]
	}
	prevCollections, _ := App.Repo.GetPrevCollections(collectionInstance)
	if len(prevCollections) > 0 {
		prevCollection = prevCollections[0]
	}
	return &nextCollection, &prevCollection
}
