package collection

import (
	"errors"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
)

func (c collectionController) CreateCollectionController(requestBody *CollectionCreateValidator, userID uint64, studioId uint64) (*models.Collection, error) {

	permissionList, err := permissions.App.Service.CalculateStudioPermissions(userID)
	if err != nil {
		return nil, err
	}
	if permissionList[studioId] != permissiongroup.PG_STUIDO_ADMIN {
		return nil, errors.New("user doesn't have access to create collection")
	}

	collection := App.Service.CollectionInstance(requestBody, userID, studioId)
	err = App.Repo.CreateCollection(&collection)
	if err != nil {
		return nil, err
	}

	err = permissions.App.Service.CreateDefaultCollectionPermission(collection.ID, collection.CreatedByID, collection.StudioID)
	if err != nil {
		return nil, err
	}

	return &collection, nil
}

func (c collectionController) UpdateCollectionController(requestBody *CollectionUpdateValidator) (*models.Collection, error) {

	updates := App.Service.Update(requestBody)
	collection, err := queries.App.CollectionQuery.UpdateCollection(requestBody.ID, updates)
	if err != nil {
		return nil, err
	}
	return collection, nil
}

func (c collectionController) DeleteCollectionController(collectionId uint64, userID uint64) error {
	err := App.Repo.Manger.SoftDeleteByID(models.COLLECTION, collectionId, userID)
	if err != nil {
		return err
	}
	return nil
}

func (c collectionController) UpdateCollectionVisibility(collectionId uint64, updates map[string]interface{}) error {
	_, err := queries.App.CollectionQuery.UpdateCollection(collectionId, updates)
	if err != nil {
		return err
	}
	return nil
}

/*
	Returns the list of the collections that has publicAccess as view.
	Args:
		studioId uint64
	Returns:
		*[]models.Collection
		error
*/
func (c collectionController) AnonymousCollectionsController(studioId uint64) (*[]models.Collection, error) {
	collections, err := App.Repo.AnonymousGetCollections(studioId)
	if err != nil {
		return nil, err
	}
	return collections, nil
}

/*
	Returns the list of the collections that user has access to.
	Args:
		studioId uint64
	Returns:
		*[]models.Collection
		error
*/
func (c collectionController) AuthUserCollectionsController(studioId uint64, user *models.User) (*[]CollectionSerializer, error) {
	collections, err := App.Repo.GetCollections(map[string]interface{}{"studio_id": studioId, "is_archived": false})
	if err != nil {
		return nil, err
	}

	permissionList, err := permissions.App.Service.CalculateCollectionPermissions(user.ID, studioId)
	if err != nil {
		return nil, err
	}

	accessCollections := &[]CollectionSerializer{}
	for _, collection := range *collections {
		collectionPermission := permissionList[collection.ID]
		if collectionPermission != "" && collectionPermission != permissiongroup.PGCollectionNone().SystemName {
			collectionData := CollectionSerializerData(&collection)
			collectionData.Permission = collectionPermission
			*accessCollections = append(*accessCollections, collectionData)
		} else {
			if collection.HasPublicCanvas == true {
				collectionData := CollectionSerializerData(&collection)
				collectionData.Permission = models.PGCollectionNoneSysName
				*accessCollections = append(*accessCollections, collectionData)
			} else if collection.PublicAccess == models.EDIT {
				collectionData := CollectionSerializerData(&collection)
				collectionData.Permission = models.PGCollectionEditSysName
				*accessCollections = append(*accessCollections, collectionData)
			} else if collection.PublicAccess == models.COMMENT {
				collectionData := CollectionSerializerData(&collection)
				collectionData.Permission = models.PGCollectionCommentSysName
				*accessCollections = append(*accessCollections, collectionData)
			} else if collection.PublicAccess == models.VIEW {
				collectionData := CollectionSerializerData(&collection)
				collectionData.Permission = models.PGCollectionViewSysName
				*accessCollections = append(*accessCollections, collectionData)
			}
		}
	}
	return accessCollections, nil
}

/*
	Move collection position to forward or backward.
	Args:
		requestBody *CollectionMoveValidator
	Returns:
		*models.Collection
		error
*/
func (c collectionController) MoveCollectionController(requestBody *CollectionMoveValidator) (collection *models.Collection, err error) {

	result, err := App.Repo.GetCollection(map[string]interface{}{"id": requestBody.CollectionId})
	if err != nil {
		return nil, err
	}

	err = App.Repo.MoveCollection(result, requestBody.Position)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c collectionController) StudioRoleCollectionsController(studioId uint64, roleID uint64) (*[]CollectionSerializer, error) {
	collections, err := App.Repo.GetCollections(map[string]interface{}{"studio_id": studioId, "is_archived": false})
	if err != nil {
		return nil, err
	}

	permissionList := permissions.App.Service.CalculateCollectionRolePermissions(roleID)
	accessCollections := &[]CollectionSerializer{}
	for _, collection := range *collections {
		collectionPermission := permissionList[collection.ID]
		if collectionPermission != "" && collectionPermission != permissiongroup.PGCollectionNone().SystemName {
			collectionData := CollectionSerializerData(&collection)
			collectionData.Permission = collectionPermission
			*accessCollections = append(*accessCollections, collectionData)
		} else {
			collectionData := CollectionSerializerData(&collection)
			collectionData.Permission = permissiongroup.PGCollectionNone().SystemName
			*accessCollections = append(*accessCollections, collectionData)
		}
	}
	return accessCollections, nil
}

/*
	Returns the list of the collections for the studio member
	Args:
		studioId uint64
	Returns:
		*[]models.Collection
		error
*/
func (c collectionController) StudioMemberCollectionsController(studioId uint64, userID uint64) (*[]CollectionSerializer, error) {
	var allTheCollectionIDs []uint64
	var ActualPermsObject []CollectionActualPermissionsObject
	memberObject, _ := App.Repo.GetMemberByUserID(userID, studioId)
	collections, err := App.Repo.GetCollections(map[string]interface{}{"studio_id": studioId, "is_archived": false})
	if err != nil {
		return nil, err
	}

	// Making list of all ht ecollections
	for _, collectionInstance := range *collections {
		allTheCollectionIDs = append(allTheCollectionIDs, collectionInstance.ID)
	}

	permissionList, err := permissions.App.Service.CalculateCollectionPermissions(userID, studioId)
	if err != nil {
		return nil, err
	}

	// We are trying to get actual perms for logged in user. via vi the collection
	// Building the new array
	for _, collectionID := range allTheCollectionIDs {
		if permissionList[collectionID] == "" {
			ActualPermsObject = append(ActualPermsObject, CollectionActualPermissionsObject{})
		} else {
			// We have calculated value of the permissions
			actualPerms := CollectionPermissionActual(collectionID, memberObject.ID, studioId)
			for _, ap := range actualPerms {
				ActualPermsObject = append(ActualPermsObject, ap)
			}
		}
	}
	fmt.Printf("%+v\n", ActualPermsObject)

	accessCollections := &[]CollectionSerializer{}
	for _, collection := range *collections {
		collectionPermission := permissionList[collection.ID]
		if collectionPermission != "" && collectionPermission != permissiongroup.PGCollectionNone().SystemName {
			collectionData := CollectionSerializerData(&collection)
			collectionData.Permission = collectionPermission
			collectionData.ActualPermsObject = PluckTheObject(ActualPermsObject, collection.ID)
			collectionData.MemberPermsObject = MemberCollectionPermissionActualCalculator(collection.ID, memberObject, studioId)
			collectionData.RolePermsObject = RoleCollectionPermissionActualCalculator(collection.ID, memberObject, studioId)
			*accessCollections = append(*accessCollections, collectionData)
		} else {
			collectionData := CollectionSerializerData(&collection)
			collectionData.Permission = permissiongroup.PGCollectionNone().SystemName
			collectionData.ActualPermsObject = PluckTheObject(ActualPermsObject, collection.ID)
			*accessCollections = append(*accessCollections, collectionData)
		}
	}
	return accessCollections, nil
}
