package pr

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
)

// @Summary 	Get All the Publish Requests on a Studio
// @Tags		PublishRequest
// @Accept       json
// @Produce      json
// @Router 		/v1/publish-requests/ [get]
func (r *prRoutes) ListPublishRequests(c *gin.Context) {
	studioID, _ := r.GetStudioId(c)
	userID := r.GetLoggedInUserId(c)
	instances, err := App.Service.GetPublishRequestsByStudio(studioID, userID)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}

	response.RenderResponse(c, PublishRequestSerializerMany(instances))
}
