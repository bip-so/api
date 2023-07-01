package studio

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/apiutil"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"net/http"
)

// @Summary 	Get Studio Vendor
// @Tags		Studio
// @Accept 		json
// @Produce 	json
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studio/external-integration/ping [get]
func (r *studioExternalRoutes) Ping(c *gin.Context) {
	vendorName, _ := c.Get("vendorName")
	vn := vendorName.(string)
	data := map[string]string{"message": "Pong", "vendor_name": vn}
	c.JSON(http.StatusOK, data)
}

// @Summary 	Get Studio Vendor
// @Tags		Studio
// @Accept 		json
// @Produce 	json
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studio/external-integration/ [get]
func (r *studioExternalRoutes) Get(c *gin.Context) {
	vendorName, _ := c.Get("vendorName")
	vn := vendorName.(string)
	guildId := c.Query("guildId")
	vendorObject, err := queries.App.StudioVendorQuery.GetStudioVendor(guildId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	data := map[string]string{
		"integrationId": vendorObject.UUID.String(),
		"guildId":       vendorObject.GuildId,
		"guildName":     vendorObject.GuildName,
		"vendorName":    vn,
		"status":        vendorObject.IntegrationStatus,
		"handle":        vendorObject.Handle,
	}

	c.JSON(http.StatusOK, data)
}

// @Summary 	Register New Integration
// @Tags		Studio
// @Accept 		json
// @Produce 	json
// @Param 		body 		body 		RegisterNewStudioIntegrationValidator true "Register New Studio Integration Validator"
// @Param       Bip-Partner-Key  header    string  true  "Vendor Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studio/external-integration/ [post]
func (r *studioExternalRoutes) Create(c *gin.Context) {
	vendorName, _ := c.Get("vendorName")
	vn := vendorName.(string)
	var body RegisterNewStudioIntegrationValidator
	if err := apiutil.Bind(c, &body); err != nil {
		fmt.Println(err, body)
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}
	// Register the record
	svInstance, _ := queries.App.StudioVendorQuery.CreateStudioVendor(body.GuildId, body.GuildName, vn)

	CURRENT_ENV := configs.GetConfigString("ENV")
	//link := models.MailerRouterPaths[CURRENT_ENV]["BASE_URL"] + "discord-integration?guildId=" + body.GuildId + "&partnerIntegrationId=" + strconv.FormatUint(svInstance.ID, 10)
	link := models.MailerRouterPaths[CURRENT_ENV]["BASE_URL"] + "discord-integration?guildId=" + body.GuildId + "&partnerIntegrationId=" + svInstance.UUID.String()

	// Return discord-integration
	data := map[string]string{
		"url":           link,
		"integrationId": svInstance.UUID.String(),
		"guildId":       body.GuildId,
		"guildName":     body.GuildName,
		"vendorName":    vn,
	}

	c.JSON(http.StatusOK, data)
}

// @Summary 	Get Studio User Points
// @Tags		Studio
// @Accept 		json
// @Produce 	json
// @Param       Bip-Partner-Key  header    string  true  "Vendor Id"
// @Param 		guildId 	query 	string	true "guildId"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/studio/external-integration/user-points [get]
func (r *studioExternalRoutes) GetUserPoints(c *gin.Context) {
	var MainStudioPointsNameSpace = "studio-user-points:"
	vendorName, _ := c.Get("vendorName")
	vn := vendorName.(string)
	guildID := c.Query("guildId")
	data := App.StudioService.cache.HGet(context.Background(), MainStudioPointsNameSpace+vn, guildID).Val()
	var studioUserPoints map[string]map[string]int
	json.Unmarshal([]byte(data), &studioUserPoints)
	response.RenderCustomResponse(c, studioUserPoints)
	return
}
