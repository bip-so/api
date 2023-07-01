package permissions

import (
	"context"

	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

// Get User Studio Permissions
// @Summary 	Gets User Studio Permissions
// @Description
// @Tags		Permissions
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/permission/studio [get]
func (r *permissionRoutes) getStudioPermission(c *gin.Context) {
	user, err := r.GetLoggedInUser(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	permissionList, err := App.Service.CalculateStudioPermissions(user.ID)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	response.RenderResponse(c, permissionList)
	return
}

// Get User Studio Collection Permissions
// @Summary 	Gets User Studio Collection Permissions
// @Description
// @Tags		Permissions
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/permission/collection [get]
func (r *permissionRoutes) getCollectionPermission(c *gin.Context) {
	user, err := r.GetLoggedInUser(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	studioId, err := r.GetStudioId(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	permissionList, err := App.Service.CalculateCollectionPermissions(user.ID, studioId)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	response.RenderResponse(c, permissionList)
	return
}

// Redis Invalidate cache API
// @Summary 	Redis Invalidate cache API
// @Description
// @Tags		RedisAPI
// @Param 		hash 	query 		string		 		true "Hash key or Direct redis key"
// @Param 		key 	query 		string		 		false "Key"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/permission/invalidate-cache [get]
func (r *permissionRoutes) invalidateCache(c *gin.Context) {
	hash := c.Query("hash")
	key := c.Query("key")
	App.Service.InvalidatePermissionCache(c.Request.Context(), hash, key)
	response.RenderSuccessResponse(c, "Invalidated successfully")
	return
}

// Get User canvas Permissions
// @Summary 	Gets User canvas Permissions
// @Description
// @Tags		Permissions
// @Security 	bearerAuth
// @Param       bip-studio-id  header    string  true  "Studio ID"
// @Param       collectionId  path    string  true  "collection Id"
// @Router 		/v1/permission/canvas/{collectionId} [get]
func (r *permissionRoutes) getCanvasPermission(c *gin.Context) {
	user, err := r.GetLoggedInUser(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	studioId, err := r.GetStudioId(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	collectionId := c.Param("collectionId")

	permissionList, err := App.Service.CalculateCanvasRepoPermissions(user.ID, studioId, utils.Uint64(collectionId))
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	response.RenderResponse(c, permissionList)
	return
}

// Get User Sub canvas Permissions
// @Summary 	Gets User canvas Permissions
// @Description
// @Tags		Permissions
// @Security 	bearerAuth
// @Param       bip-studio-id  header    string  true  "Studio ID"
// @Param       collectionId  path    string  true  "collectionId Id"
// @Param       parentCanvasId  path    string  true  "collectionId Id"
// @Router 		/v1/permission/canvas/{collectionId}/{parentCanvasId} [get]
func (r *permissionRoutes) getSubCanvasPermission(c *gin.Context) {
	user, err := r.GetLoggedInUser(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	studioId, err := r.GetStudioId(c)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	collectionId := c.Param("collectionId")
	parentCanvasId := c.Param("parentCanvasId")

	permissionList, err := App.Service.CalculateSubCanvasRepoPermissions(user.ID, studioId, utils.Uint64(collectionId), utils.Uint64(parentCanvasId))
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	response.RenderResponse(c, permissionList)
	return
}

// Redis Flush API
// @Summary 	Redis Flush API
// @Description
// @Tags		RedisAPI
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/permission/flush-cache [get]
func (r *permissionRoutes) flushCache(c *gin.Context) {
	redis.RedisClient().FlushAll(context.Background())
	response.RenderSuccessResponse(c, "Flushed redis cache successfully")
	return
}
