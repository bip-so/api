package collection

import (
	"encoding/json"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/apiClient"
	"net/http"
	"strconv"

	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"gorm.io/gorm"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/apiutil"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/context"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"

	"github.com/gin-gonic/gin"
)

func (r collectionRoutes) test(c *gin.Context) {
	fmt.Println("Test")
	var col models.Collection
	//err := App.Repo.Manger.HardDeleteByID(col.TableName(), 29)
	//err := App.Repo.Manger.SoftDeleteByID(col.TableName(), 125, 18)
	result := map[string]interface{}{"name": "I am the UPDATEDDD"}
	err := App.Repo.Manger.UpdateEntityByID(col.TableName(), 148, result)
	fmt.Println(err)
}

// Create New Collection
// @Summary 	Create Collection
// @Description
// @Tags		Collection
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Param 		body 		body 		CollectionCreateValidator true "Create Collection"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/collection/create [post]
func (r collectionRoutes) createCollectionRoute(c *gin.Context) {
	var body *CollectionCreateValidator
	if err := apiutil.Bind(c, &body); err != nil {
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}

	user, err := context.GetAuthUser(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	studioId, err := context.GetAuthStudio(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	// Permission check
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnStudio(authUserID, studioId, permissiongroup.STUDIO_CREATE_COLLECTION); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	collection, err := App.Controller.CreateCollectionController(body, user.ID, studioId)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	serialized := CollectionSerializerData(collection)
	serialized.Permission = models.PGCollectionModerateSysName
	response.RenderResponse(c, serialized)
}

// Update Collection
// @Summary 	Update Collection
// @Description
// @Tags		Collection
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		body 		body 		CollectionUpdateValidator true "Update Collection"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/collection/update [put]
func (r collectionRoutes) updateCollectionRoute(c *gin.Context) {
	var body *CollectionUpdateValidator
	if err := apiutil.Bind(c, &body); err != nil {
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}

	// Permission check
	studioId, _ := r.GetStudioId(c)
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnCollection(authUserID, studioId, body.ID, permissiongroup.COLLECTION_EDIT_NAME); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	collection, err := App.Controller.UpdateCollectionController(body)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}

	go func() {
		payload, _ := json.Marshal(map[string]uint64{"collectionId": collection.ID})
		apiClient.AddToQueue(apiClient.UpdateDiscordTreeMessage, payload, apiClient.DEFAULT, apiClient.CommonRetry)
	}()
	response.RenderResponse(c, CollectionSerializerData(collection))
	return
}

// Collection Visibility
// @Summary 	Update Collection Visibility - "private", "view", "comment", "edit"
// @Tags		Collection
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		collectionId 	path 		string		 		true "Collection Id"
// @Param 		inherit 	query 		string		 		false "inherit visibilty"
// @Param 		body 		body 		VisibilityUpdateValidator true "Update Collection"
// @Success 	200 		{object} 	response.ApiResponse
// @Router 		/v1/collection/{collectionId}/visibility [post]
func (r collectionRoutes) manageVisibilityCollectionRoute(c *gin.Context) {
	var body VisibilityUpdateValidator
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	verr := body.Validate()
	if verr != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), verr))
		return
	}
	inherit := c.Query("inherit")
	collectionId, _ := c.Params.Get("collectionId")
	parsedcollectionId, err := strconv.ParseUint(collectionId, 10, 64)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	// Permission check
	studioId, _ := r.GetStudioId(c)
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnCollection(authUserID, studioId, parsedcollectionId, permissiongroup.COLLECTION_MANAGE_PERMS); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	user, err := context.GetAuthUser(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	updates := map[string]interface{}{}
	if inherit == "true" && body.PublicAccess == models.PRIVATE {
		updates = map[string]interface{}{"public_access": body.PublicAccess, "updated_by_id": user.ID, "has_public_canvas": false}
	} else {
		updates = map[string]interface{}{"public_access": body.PublicAccess, "updated_by_id": user.ID}
	}
	errUpdating := App.Controller.UpdateCollectionVisibility(parsedcollectionId, updates)
	if errUpdating != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), errUpdating))
		return
	}

	go func() {
		if inherit == "true" {
			canvasRepos, _ := App.Repo.GetCanvasRepos(map[string]interface{}{"collection_id": parsedcollectionId})
			for _, repo := range *canvasRepos {
				App.Repo.UpdateCanvasBranch(*repo.DefaultBranchID, map[string]interface{}{"public_access": body.PublicAccess, "updated_by_id": authUserID})
			}
		}
		payload, _ := json.Marshal(map[string]uint64{"collectionId": parsedcollectionId})
		apiClient.AddToQueue(apiClient.UpdateDiscordTreeMessage, payload, apiClient.DEFAULT, apiClient.CommonRetry)
	}()

	response.RenderSuccessResponse(c, "Collection Visibility Updated.")
	return
}

// Delete Collection
// @Summary 	Delete Collection
// @Description
// @Tags		Collection
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		collectionId 	path 		string		 		true "Collection Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/collection/delete/{collectionId} [delete]
func (r collectionRoutes) deleteCollectionRoute(c *gin.Context) {
	collectionId, _ := c.Params.Get("collectionId")
	parsedcollectionId, err := strconv.ParseUint(collectionId, 10, 64)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}

	studioId, err := context.GetAuthStudio(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	// Permission check
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnCollection(authUserID, studioId, parsedcollectionId, permissiongroup.COLLECTION_DELETE); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	err = App.Controller.DeleteCollectionController(parsedcollectionId, authUserID)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}

	go func() {
		collection, _ := App.Repo.GetCollection(map[string]interface{}{"id": collectionId})
		err = App.Repo.ResetCollectionPositionOnDelete(collection)
		if err != nil {
			response.RenderCustomErrorResponse(c, err)
			return
		}
		payload, _ := json.Marshal(map[string]uint64{"collectionId": collection.ID})
		apiClient.AddToQueue(apiClient.UpdateDiscordTreeMessage, payload, apiClient.DEFAULT, apiClient.CommonRetry)
	}()

	response.RenderSuccessResponse(c, "Collection Deleted Succesfully")
	return
}

// Get Studio Collections
// @Summary 	Get all Collections related to studioID
// @Description
// @Tags		Collection
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/collection/get [get]
func (r collectionRoutes) getCollectionRoute(c *gin.Context) {
	studioId, err := App.RouteHandler.GetStudioId(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	user, _ := App.RouteHandler.GetLoggedInUser(c)
	public := c.Query("public")
	// Checking if the request is made by loggedIn user or not
	// If user == nil we trigger the Anonymous flow to get the collection or vice-versa
	if user == nil || public == "true" {
		var collections *[]models.Collection
		collections, err = App.Controller.AnonymousCollectionsController(studioId)
		if err != nil {
			response.RenderCustomErrorResponse(c, err)
			return
		}
		response.RenderResponse(c, MultiCollectionSerializerData(collections))
		return
	} else {
		_, err := permissions.App.Repo.GetMember(map[string]interface{}{"studio_id": studioId, "user_id": user.ID})
		if err == gorm.ErrRecordNotFound {
			var collections *[]models.Collection
			collections, err = App.Controller.AnonymousCollectionsController(studioId)
			if err != nil {
				response.RenderCustomErrorResponse(c, err)
				return
			}
			response.RenderResponse(c, MultiCollectionSerializerData(collections))
			return
		} else {
			var collections *[]CollectionSerializer
			collections, err = App.Controller.AuthUserCollectionsController(studioId, user)
			if err != nil {
				response.RenderCustomErrorResponse(c, err)
				return
			}
			response.RenderResponse(c, collections)
			return
		}
	}
}

// Move Collection Position
// @Summary 	Move Collection Position
// @Description
// @Tags		Collection
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		body 		body 		CollectionMoveValidator true "Move Collection"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/collection/move [post]
func (r collectionRoutes) moveCollectionRoute(c *gin.Context) {
	var body *CollectionMoveValidator
	if err := apiutil.Bind(c, &body); err != nil {
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}

	// Permission check
	studioId, _ := r.GetStudioId(c)
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnStudio(authUserID, studioId, permissiongroup.STUDIO_CHANGE_CANVAS_COLLECTION_POSITION); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	_, err := App.Controller.MoveCollectionController(body)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	go func() {
		payload, _ := json.Marshal(map[string]uint64{"collectionId": body.CollectionId})
		apiClient.AddToQueue(apiClient.UpdateDiscordTreeMessage, payload, apiClient.DEFAULT, apiClient.CommonRetry)
	}()
	response.RenderResponse(c, "Collection Moved Succesfully")
}

// Get Studio Collections of user
// @Summary 	Get all Collections to display in member permissions
// @Description
// @Tags		Collection
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Param       userId  path    string  true  "User Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/collection/user/{userId} [get]
func (r collectionRoutes) getStudioMemberCollections(c *gin.Context) {
	studioID, _ := r.GetStudioId(c)
	userIDStr, _ := c.Params.Get("userId")
	userID := utils.Uint64(userIDStr)

	// Permission check
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnStudio(authUserID, studioID, permissiongroup.STUDIO_ADD_REMOVE_USER_TO_ROLE); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	var collections *[]CollectionSerializer
	collections, err := App.Controller.StudioMemberCollectionsController(studioID, userID)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	response.RenderResponse(c, collections)
	return
}

// Get Studio Collections of role
// @Summary 	Get all Collections to display in role permissions
// @Description
// @Tags		Collection
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Param       roleId  path    string  true  "Role Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/collection/role/{roleId} [get]
func (r collectionRoutes) getStudioRoleCollections(c *gin.Context) {
	studioID, _ := r.GetStudioId(c)
	roleIDStr, _ := c.Params.Get("roleId")
	roleID := utils.Uint64(roleIDStr)

	// Permission check
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnStudio(authUserID, studioID, permissiongroup.STUDIO_ADD_REMOVE_USER_TO_ROLE); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	var collections *[]CollectionSerializer
	collections, err := App.Controller.StudioRoleCollectionsController(studioID, roleID)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	response.RenderResponse(c, collections)
	return
}

// GetNextPrevCanvas Get next and prev collection
// @Summary 	Get next and prev collection
// @Description
// @Tags		Collection
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Param       collectionId  path    string  true  "collectionId"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/collection/next-prev/{collectionId} [post]
func (r *collectionRoutes) getNextPrevCollection(c *gin.Context) {
	userID := r.GetLoggedInUserId(c)
	collectionIDStr := c.Param("collectionId")
	collectionId := utils.Uint64(collectionIDStr)

	nextCollection, PrevCollection := App.Service.GetCollectionPrevAndNext(userID, collectionId)

	resp := map[string]interface{}{
		"next": CollectionSerializerDataMini(nextCollection),
		"prev": CollectionSerializerDataMini(PrevCollection),
	}

	response.RenderCustomResponse(c, resp)
}
