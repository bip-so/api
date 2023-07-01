package workflows

import (
	"gitlab.com/phonepost/bip-be-platform/internal/collection"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

//WorkflowCreateCollectionAndPerms: This created a Collections and sets proper permissions
func WorkflowCreateCollectionAndPerms(collectionName string, position uint, collectionAccess string, studioCreatedByID uint64, studioID uint64) (*models.Collection, error) {
	collectionView := &collection.CollectionCreateValidator{
		Name:         collectionName,
		Position:     position,
		PublicAccess: collectionAccess,
	}
	collectionInstance, err := collection.App.Controller.CreateCollectionController(collectionView, studioCreatedByID, studioID)
	return collectionInstance, err
}
