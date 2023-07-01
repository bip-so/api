package collectionpermissions

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"

	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
)

var (
	CollectionPermissionsRouteHandler collectionPermissionsRoutes
)

// Get Collection Permissions
// @Summary 	Gets collection permissions by collection id
// @Description
// @Tags		Collection Permission
// @Accept 		json
// @Produce 	json
// @Param 		collectionid 	path 		string		 		true "Collection Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/collectionpermission/{collectionid} [get]
func (r *collectionPermissionsRoutes) getCollectionPermissionsRouteHandler(c *gin.Context) {
	collectionIDStr := c.Param("collectionId")
	collectionID, err := strconv.ParseUint(collectionIDStr, 10, 64)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}

	// Permission check
	studioID, _ := r.GetStudioId(c)
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnCollection(authUserID, studioID, collectionID, permissiongroup.COLLECTION_VIEW_METADATA); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	collectionPerms, err := CollectionPermissionService.GetCollectionPermissions(map[string]interface{}{"collection_id": collectionID})

	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	var result []CollectionPermissionsSerializer
	for _, perm := range collectionPerms {
		result = append(result, *SerializeCollectionPermission(&perm))
	}
	response.RenderResponse(c, result)

}

// Update Collection Permissions
// @Summary 	Create/Update collection permissions
// @Description
// @Tags		Collection Permission
// @Accept 		json
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Param 		body  	body 		CollectionPermissionValidator true "collection permission data"
// @Produce 	json
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/collectionpermission/update [post]
func (r *collectionPermissionsRoutes) createCollectionPermissionsRouteHandler(c *gin.Context) {

	rawStudioId, _ := c.Get("currentStudio")
	studioIDInt := rawStudioId.(uint64)
	user, _ := r.GetLoggedInUser(c)
	inheritPerms := c.Query("inherit")

	var err error
	var collPerm *models.CollectionPermission
	var body CollectionPermissionValidator
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Permission check
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnCollection(authUserID, studioIDInt, body.CollectionId, permissiongroup.COLLECTION_MANAGE_PERMS); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	if body.RoleID != 0 && body.MemberID != 0 {
		err = errors.New("provide either roleId or memberId")
	}
	if body.RoleID != 0 {
		collPerm, err = CollectionPermissionService.UpdateCollectionPermissions(map[string]interface{}{"collection_id": body.CollectionId, "role_id": body.RoleID}, body, studioIDInt, user.ID)

	} else if body.MemberID != 0 {
		collPerm, err = CollectionPermissionService.UpdateCollectionPermissions(map[string]interface{}{"collection_id": body.CollectionId, "member_id": body.MemberID}, body, studioIDInt, user.ID)

	} else {
		err = errors.New("roleId or memberId not provided")
	}

	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}

	// clearing redis cache
	// @todo Later move this to kafka
	go func() {
		if inheritPerms == "true" {
			App.Service.InheritUserPermsToCanvas(body, studioIDInt, authUserID)
		}
		inv := &permissions.InvalidatePermissions{
			MemberID:       &body.MemberID,
			RoleID:         &body.RoleID,
			InvalidationOn: "collection",
		}
		err = permissions.App.Service.InvalidatePermissions(inv)
		if err != nil {
			fmt.Println(err)
		}
		if body.RoleID != 0 {
			permissions.App.Service.InvalidateRolePermissionCache(body.RoleID, studioIDInt)
		}
	}()

	response.RenderResponse(c, *SerializeCollectionPermission(collPerm))

}

// Delete Collection Permissions
// @Summary 	Delete collection permissions by collection id
// @Description
// @Tags		Collection Permission
// @Accept 		json
// @Produce 	json
// @Param       collectionPermissionId  path    string  true  "collectionId Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/collectionpermission/{collectionPermissionId} [delete]
func (r *collectionPermissionsRoutes) deleteCollectionPermissionsRouteHandler(c *gin.Context) {

	collectionIDStr := c.Param("collectionPermissionId")
	collectionPermID, err := strconv.ParseUint(collectionIDStr, 10, 64)
	inheritPerms := c.Query("inherit")

	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}

	var col models.CollectionPermission
	collectionPermission, _ := App.Repo.Get(map[string]interface{}{"id": collectionPermID})

	// Permission check
	studioID, _ := r.GetStudioId(c)
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnCollection(authUserID, studioID, collectionPermission.CollectionId, permissiongroup.COLLECTION_MANAGE_PERMS); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	err = App.Repo.Manger.HardDeleteByID(col.TableName(), collectionPermID)

	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}

	// clearing redis cache
	// @todo Later move this to kafka
	go func() {
		if inheritPerms == "true" {
			App.Service.DeleteInheritCollectionPermission(collectionPermission)
		} else {
			// we can add view metadata permission for the removed permission on collection
			App.Service.AddViewMetaDataPermOnCollection(collectionPermission, authUserID)
		}
		inv := &permissions.InvalidatePermissions{
			MemberID:       collectionPermission.MemberId,
			RoleID:         collectionPermission.RoleId,
			InvalidationOn: "collection",
		}
		err = permissions.App.Service.InvalidatePermissions(inv)
		if err != nil {
			fmt.Println(err)
		}
	}()

	response.RenderSuccessResponse(c, "deleted permission successfully.")

}
