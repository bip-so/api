package studio_integration

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/apiutil"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
)

// @Summary 	Get Integrations Settings
// @description Get discord, slack, dm integration page settings
// @Tags		StudioIntegrations
// @Accept       json
// @Produce      json
// @Router 		/v1/integrations/settings [get]
func (r *studioIntegrationsRoutes) GetSettings(c *gin.Context) {

	_, status := middlewares.StudioHeaderRequiredCheck(c)
	if !status {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": "Studio Header Missing",
		})
		return
	}
	//user, _ := r.GetLoggedInUser(c)
	studioID, _ := r.GetStudioId(c)

	result, err := App.Controller.GetSettings(studioID)
	if err != nil {
		response.RenderErrorResponse(c, err)
		return
	}

	response.RenderResponse(c, result)
}

// @Summary 	Delete Integration
// @description Delete integration from studio
// @Tags		StudioIntegrations
// @Accept       json
// @Produce      json
//@Param 		type 		query  string  true "integration type"
// @Router 		/v1/integrations [delete]
func (r *studioIntegrationsRoutes) DeleteIntegration(c *gin.Context) {
	integrationType := c.Query("type")
	studioId, _ := r.GetStudioId(c)

	_, err := App.Repo.DeleteStudioIntegration(studioId, integrationType)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	response.RenderSuccessResponse(c, "Successfully deleted integration")
	return
}

// @Summary 	Update discord dm notifications
// @description Update discord dm notifications status
// @Tags		StudioIntegrations
// @Accept       json
// @Produce      json
//@Param 		body 		body 		UpdateDiscordNotification true "Update discord dm notifications"
// @Router 		/v1/integrations/discord [put]
func (r *studioIntegrationsRoutes) UpdateDiscordDmNotifications(c *gin.Context) {
	studioId, _ := r.GetStudioId(c)
	var body *UpdateDiscordNotification
	if err := apiutil.Bind(c, &body); err != nil {
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}

	err := App.Repo.UpdateDiscordDmNotification(studioId, body.Status)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	response.RenderSuccessResponse(c, "Successfully updated integration")
	return
}

// @Summary 	Update slack dm notifications
// @description Update slack dm notifications status
// @Tags		StudioIntegrations
// @Accept       json
// @Produce      json
//@Param 		body 		body 		UpdateSlackNotification true "Update slack dm notifications"
// @Router 		/v1/integrations/slack [put]
func (r *studioIntegrationsRoutes) UpdateSlackDmNotifications(c *gin.Context) {
	studioId, _ := r.GetStudioId(c)
	var body *UpdateDiscordNotification
	if err := apiutil.Bind(c, &body); err != nil {
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}

	err := App.Repo.UpdateSlackDmNotification(studioId, body.Status)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	response.RenderSuccessResponse(c, "Successfully updated integration")
	return
}

// @Summary 	Get Integrations Settings
// @description Check if studio needs to the integration done again.
// @Tags		StudioIntegrations
// @Accept       json
// @Produce      json
// @Router 		/v1/integrations/discord/update [get]
func (r *studioIntegrationsRoutes) CheckUpdateIntegration(c *gin.Context) {

	_, status := middlewares.StudioHeaderRequiredCheck(c)
	if !status {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": "Studio Header Missing",
		})
		return
	}
	studioID, _ := r.GetStudioId(c)

	result, err := App.Repo.GetDiscordStudioIntegration(studioID)
	if err != nil {
		response.RenderErrorResponse(c, err)
		return
	}

	resp := map[string]interface{}{}
	if result.MessagesData == nil {
		resp["updateIntegration"] = true
	} else {
		resp["updateIntegration"] = false
	}

	response.RenderCustomResponse(c, resp)
}
