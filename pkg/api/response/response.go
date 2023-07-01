package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseView is the struct of response
type ResponseView struct {
	Data    interface{} `json:"data,omitempty"`
	Next    string      `json:"next,omitempty"`
	Error   error       `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// RenderResponse respond a special data, e.g.: Topic, Category etc.
func RenderResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, ResponseView{Data: data})
}

// RenderPaginatedResponse responds data with next url
func RenderPaginatedResponse(c *gin.Context, data interface{}, next string) {
	c.JSON(http.StatusOK, ResponseView{Data: data, Next: next})
}

// RenderErrorResponse respond an error response
func RenderErrorResponse(c *gin.Context, err error) {
	sessionError, ok := err.(Error)
	if !ok {
		sessionError = ServerError(c.Request.Context(), err)
	}
	// Todo: why?
	if sessionError.Code == 10001 {
		sessionError.Code = 500
	}
	c.JSON(sessionError.Status, ResponseView{Error: sessionError})
}

func RenderEntityNotUnprocessable(c *gin.Context, err error) {
	sessionError, ok := err.(Error)
	if !ok {
		sessionError = ServerError(c.Request.Context(), err)
	}
	c.JSON(http.StatusUnprocessableEntity, ResponseView{Error: sessionError})
}

// RenderBlankResponse respond a blank response
func RenderBlankResponse(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{})
}

// RenderCustomResponse response custom response
func RenderCustomResponse(c *gin.Context, response interface{}) {
	c.JSON(http.StatusOK, response)
}

// RenderSuccessResponse with a text
func RenderSuccessResponse(c *gin.Context, message string) {
	c.JSON(http.StatusOK, ResponseView{Message: message})
}

// RenderSuccessResponse with a json
func RenderOkWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, ResponseView{Data: data})
}

func RenderNotFoundResponse(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ResponseView{Message: message})
}

func RenderCustomErrorResponse(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, map[string]interface{}{
		"error": err.Error(),
	})
}

func RenderCustomErrorDataResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusBadRequest, data)
}

func RenderPermissionError(c *gin.Context) {
	c.JSON(http.StatusForbidden, map[string]interface{}{
		"error": "User doesn't have permission",
	})
}

func RenderPermissionErrorDataResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusForbidden, data)
}

func RenderInternalServerErrorResponse(c *gin.Context, err interface{}) {
	c.JSON(http.StatusInternalServerError, map[string]interface{}{
		"error": err,
	})
}

type ApiResponse struct{}
