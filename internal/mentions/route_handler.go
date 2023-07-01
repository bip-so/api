package mentions

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
	"net/http"
)

// @Summary 	Create a Mention
// @Description Requires a scope: scope can be "block","block_comment","reel","reel_comment", "block_thread"
// @Accept       json
// @Produce      json
// @Tags		Mentions
// @Security 	bearerAuth
// @Param 		body  		body 		MentionPost true "Create Mention"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/mentions/ [post]
func (r *mentionsRoutes) AddMention(c *gin.Context) {
	var body MentionPost
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	studioID, errGetStudioId := r.GetStudioId(c)
	if errGetStudioId != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errGetStudioId.Error(),
		})
		return
	}

	validate := body.Validate()
	if validate != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), validate))
		return
	}
	user, errUserGet := r.GetLoggedInUser(c)
	if errUserGet != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), errUserGet))
		return
	}
	mappeddata, err := App.Service.AddMentionManager(body, user, studioID)
	if err != nil && err.Error() == response.NoPermissionError {
		response.RenderPermissionError(c)
		return
	}
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	c.JSON(http.StatusOK, mappeddata)
}
