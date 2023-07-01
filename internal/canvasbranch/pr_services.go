package canvasbranch

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/apiClient"
)

// We need to determine is user can do a PR self or not?
func (s canvasBranchService) InitPrRequest(branchID uint64, userID uint64) (*models.CanvasBranch, bool, error) {
	// Add validation branchID & userID
	if queries.App.PublishRequestQuery.PublishRequestExists(branchID, userID) {
		return nil, false, errors.New("Publish Request Already Exists for this Branch and This User.! ")
	}
	branchInstance, errGettingBranch := App.Repo.GetBranchRepoInstance(map[string]interface{}{"id": branchID})
	if errGettingBranch != nil {
		return nil, false, errGettingBranch
	}
	// Validation : Default Branch for this Repo. Check is the Branch ID is same as the Default Branch ID on the Canvas Repository
	if branchInstance.ID != *branchInstance.CanvasRepository.DefaultBranchID {
		return nil, false, errors.New("Publish request can only be applied to Default Branches")
	}
	// Check :  Does this Branch belong to a Root Repo or ChildRepo
	// Col _ R1 DEFDranh R2 ->

	isBranchInsideRootRepo := App.Service.doesBelongToRootRepo(*branchInstance)
	if isBranchInsideRootRepo {
		// If isBranchInsideRootRepo == True => Workflow 1
		// Branch is in a Root Repo and a Default Branch, now we check the permission on the Collection this belogs to
		// Get PG with `CollectionBranchPermissionGroupByUserID`
		canManagePR, errGettingPerms := permissions.App.Service.CanUserDoThisOnCollection(userID, branchInstance.CanvasRepository.StudioID, branchInstance.CanvasRepository.CollectionID, permissiongroup.COLLECTION_MANAGE_PUBLISH_REQUEST)
		if errGettingPerms != nil {
			return nil, false, errors.New("Did not get valid permissions for this user")
		}
		//// Check permission on COLLECTION, Can this person do PR or Not.
		//canManagePR := models.CollectionPermissionsMap[*pg][permissiongroup.COLLECTION_MANAGE_PUBLISH_REQUEST]
		if canManagePR {
			// Person can publish as they have perms on
			//Collection COLLECTION_MANAGE_PUBLISH_REQUEST == 1
			return branchInstance, true, nil
			//return errors.New(permissiongroup.CANVAS_BRANCH_CREATE_MERGE_REQUEST + " user does not have permission to perform this action.")
		} else {
			// Person can ONLY create a Publish Request as they have lack perms on
			// Collection COLLECTION_MANAGE_PUBLISH_REQUEST == 0
			return branchInstance, false, nil
		}
	}
	// This Flow is for Repo being a Child
	// Here we have to check the permission on the Parent Repo -> Default Branch
	defaultBranchIDofParentRepo := App.Service.defaultBranchOnParentCanvasRepo(*branchInstance)
	if defaultBranchIDofParentRepo == 0 {
		return nil, false, errors.New("We didn't get any permission on Canvas Repo -> Default ")
	}
	defaultBranchPG, errGettingPg := permissions.App.Service.CanvasBranchPermissionGroupByUserID(userID, defaultBranchIDofParentRepo)
	if errGettingPg != nil {
		return nil, false, errors.New("Did not get valid permissions for this user")
	}
	canManagePR := models.CanvasPermissionsMap[*defaultBranchPG][permissiongroup.CANVAS_BRANCH_MANAGE_PUBLISH_REQUESTS]
	if canManagePR == 1 {
		// Person can publish as they have perms on
		// Default Branch on the Parent Repo  CANVAS_BRANCH_MANAGE_PUBLISH_REQUESTS == 1
		return branchInstance, true, nil
		//return errors.New(permissiongroup.CANVAS_BRANCH_CREATE_MERGE_REQUEST + " user does not have permission to perform this action.")
	} else {
		// Person can ONLY create a Publish Request as they have lack perms on
		// Canvas Default Branch of the parent Repo   CANVAS_BRANCH_MANAGE_PUBLISH_REQUESTS == 0
		return branchInstance, false, nil
	}

}

// Manage PR
func (s canvasBranchService) ManagePublishRequest(branchID uint64, prID uint64, user *models.User, data ManagePublishRequest) error {
	var status string
	if data.Accept {
		status = models.PUBLISH_REQUEST_ACCEPTED
	} else {
		status = models.PUBLISH_REQUEST_REJECTED
	}
	var query = map[string]interface{}{
		"status":              status,
		"reviewed_by_user_id": user.ID,
	}
	err := queries.App.PublishRequestQuery.UpdatePublishRequest(prID, query)
	// We update the Publish Request Instance and Then if true, we'll also update `PublishCanvasBranch`
	if data.Accept {
		_ = App.Service.PublishCanvasBranch(branchID, user, true)
	}

	if err != nil {
		return err
	}

	go func() {
		contentObject := models.PUBLISHREQUEST
		branch, _ := queries.App.BranchQuery.GetBranchByID(branchID)
		extraData := notifications.NotificationExtraData{
			CollectionID:   branch.CanvasRepository.CollectionID,
			CanvasRepoID:   branch.CanvasRepositoryID,
			CanvasBranchID: branchID,
			Status:         status,
		}
		notifications.App.Service.PublishNewNotification(notifications.PublishRequestedUpdate, user.ID, nil,
			&branch.CanvasRepository.StudioID, nil, extraData, &prID, &contentObject)

		if data.Accept {
			notifications.App.Service.PublishNewNotification(notifications.BipMarkMessageAdded, user.ID, nil,
				&branch.CanvasRepository.StudioID, nil, extraData, &prID, &contentObject)
			go func() {
				payload, _ := json.Marshal(map[string]uint64{"collectionId": branch.CanvasRepository.CollectionID})
				apiClient.AddToQueue(apiClient.UpdateDiscordTreeMessage, payload, apiClient.DEFAULT, apiClient.CommonRetry)
			}()
		} else {
			prInstance, _ := queries.App.PublishRequestQuery.PublishRequestGetter(branchID, user.ID, prID)
			prInstanceStr, _ := json.Marshal(prInstance)
			apiClient.AddToQueue(apiClient.DeleteModsOnCanvas, []byte(prInstanceStr), apiClient.DEFAULT, apiClient.CommonRetry)
		}
	}()
	return nil

}

func (s canvasBranchService) DirectPublishRequest(instance models.CanvasBranch, message string, user models.User) error {
	err := App.Service.PublishCanvasBranch(instance.ID, &user, true)
	if err != nil {
		return err
	}

	go func() {
		extraData := notifications.NotificationExtraData{
			CollectionID:   instance.CanvasRepository.CollectionID,
			CanvasRepoID:   instance.CanvasRepositoryID,
			CanvasBranchID: instance.ID,
			Status:         models.PUBLISH_REQUEST_ACCEPTED,
			Message:        message,
		}
		notifications.App.Service.PublishNewNotification(notifications.BipMarkMessageAdded, user.ID, nil,
			&instance.CanvasRepository.StudioID, nil, extraData, nil, nil)
	}()
	return nil

}

func (s canvasBranchService) CreatePublishRequest(instance models.CanvasBranch, message string, userID uint64) (*models.PublishRequest, error) {
	var requestModel models.PublishRequest
	pr := requestModel.NewPublishRequest(
		instance.CanvasRepository.StudioID,
		instance.CanvasRepositoryID,
		instance.ID,
		models.PUBLISH_REQUEST_PENDING,
		message,
		userID,
		userID,
	)
	fmt.Println(pr)
	pr, err := queries.App.PublishRequestQuery.CreatePublishRequest(pr)
	fmt.Println("Creating publish requst id", pr.ID)
	if err != nil {
		return pr, err
	}

	// Add the mods of the immediate parent to this canvas when canvas publish is requested.
	go permissions.App.Service.SetCanvasPermissionsOnPublishRequest(instance)

	go func() {
		extraData := notifications.NotificationExtraData{
			CollectionID:   instance.CanvasRepository.CollectionID,
			CanvasRepoID:   instance.CanvasRepositoryID,
			CanvasBranchID: instance.ID,
			Status:         models.PUBLISH_REQUEST_PENDING,
			Message:        pr.Message,
		}
		contentObject := models.PUBLISHREQUEST
		notifications.App.Service.PublishNewNotification(notifications.PublishRequested, userID, nil,
			&instance.CanvasRepository.StudioID, nil, extraData, &pr.ID, &contentObject)
	}()
	return pr, nil
}

func (s canvasBranchService) GetPublishRequestsByBranch(studioID uint64, branchID uint64, userID uint64) (*[]models.PublishRequest, error) {
	instances, err := queries.App.PublishRequestQuery.GetAllPublishRequests(map[string]interface{}{"studio_id": studioID, "canvas_branch_id": branchID, "status": models.PUBLISH_REQUEST_PENDING})
	if err != nil {
		return nil, err
	}
	return instances, nil
}
