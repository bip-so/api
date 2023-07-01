package ar

import (
	"errors"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/internal/role"
)

func (s arService) CreateAccessRequest(data queries.CreateAccessRequestPost, userID uint64) (bool, error) {
	// Check for Duplicate
	exists := queries.App.AccessRequestQuery.AccessRequestExists(data, userID)
	if exists {
		return exists, errors.New("User has already requested access for this Canvas")
	}
	// Create Object and Validate
	var ar models.AccessRequest
	// Create ObjectRemember : ACCESS_REQUEST_PENDING is set by default
	newRequestObject := ar.NewAccessRequest(
		data.StudioID,
		data.CollectionID,
		data.CanvasRepositoryID,
		data.CanvasBranchID,
		userID,
	)
	_, errCreatingAr := queries.App.AccessRequestQuery.CreateAccessRequest(newRequestObject)
	if errCreatingAr != nil {
		return exists, errCreatingAr
	}

	// Todo Add Notifications
	go func() {
		extraData := notifications.NotificationExtraData{
			CanvasRepoID:   data.CanvasRepositoryID,
			CanvasBranchID: data.CanvasBranchID,
			Status:         models.ACCESS_REQUEST_PENDING,
		}
		contentObject := models.ACCESSREQUEST
		notifications.App.Service.PublishNewNotification(notifications.AccessRequested, userID, nil, &data.StudioID,
			nil, extraData, &newRequestObject.ID, &contentObject)
	}()
	return exists, nil
}

func (s arService) GetAllAccessRequest(studioID uint64) (*[]models.AccessRequest, error) {
	return queries.App.AccessRequestQuery.GetAllAccessRequests(map[string]interface{}{"studio_id": studioID})
}

// Update Access Object
// Check Status
func (s arService) ManageAccessRequest(accessRequestID uint64, data ManageAccessRequestPost, authUserID uint64) error {
	arInstance, err := queries.App.AccessRequestQuery.AccessRequestInstance(accessRequestID)
	if err != nil {
		return err
	}

	if data.Status == models.ACCESS_REQUEST_ACCEPTED {
		// Check if this user is Member for this Studio
		member, memberErr := queries.App.StudioQueries.SafeAddUserToStudio(arInstance.CreatedByID, arInstance.StudioID)
		if memberErr != nil {
			return memberErr
		}
		// Add User to StudioMemebr Role
		errAddingMemebrToRole := role.App.Repo.AddMembersToMemberRoleForStudio(arInstance.StudioID, member.ID)
		if errAddingMemebrToRole != nil {
			return errAddingMemebrToRole
		}
		// Add Perms
		errCreatingPerms := permissions.App.Service.CreateCanvasBranchPermission(
			arInstance.CollectionID,
			arInstance.CreatedByID,
			arInstance.StudioID,
			arInstance.CanvasRepositoryID,
			arInstance.CanvasBranchID,
			arInstance.CanvasRepository.ParentCanvasRepositoryID,
			data.CanvasBranchPermissionGroup,
		)
		if errCreatingPerms != nil {
			return errCreatingPerms
		}
	} else if data.Status == models.ACCESS_REQUEST_REJECTED {
		// Need to see if anything needs to be done here.
	} else {
		return errors.New("Status is wrong! Rejecting the process.")
	}
	// Check if the user is Member of the studio?
	// Add if not
	// Add Permissions for this User as Member on Canvas
	errUpdatingAccessRequest := queries.App.AccessRequestQuery.UpdateAccessRequests(accessRequestID, data.CanvasBranchPermissionGroup, data.Status)
	if err != errUpdatingAccessRequest {
		return errUpdatingAccessRequest
	}

	go func() {
		extraData := notifications.NotificationExtraData{
			CanvasRepoID:    arInstance.CanvasRepositoryID,
			Status:          data.Status,
			PermissionGroup: data.CanvasBranchPermissionGroup,
			CanvasBranchID:  arInstance.CanvasBranchID,
		}
		if arInstance.CanvasBranchPermissionGroup != nil {
			extraData.PermissionGroup = *arInstance.CanvasBranchPermissionGroup
		}
		if arInstance.Message != nil {
			extraData.Message = *arInstance.Message
		}
		contentObject := models.ACCESSREQUEST
		notifications.App.Service.PublishNewNotification(notifications.AccessRequestedUpdate, authUserID, nil, &arInstance.StudioID,
			nil, extraData, &arInstance.ID, &contentObject)
	}()

	return nil
}
