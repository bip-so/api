package parser2

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/apiutil"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"net/url"
	"os"
)

// GetCanvasBranchMarkdownFile
// Download markdown of the blocks data
// @Summary 	Get markdown file of the canvas Branch
// @Description
// @Tags		Parser
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		canvasBranchID 	path 	string	true "Canvas Branch Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/parser/markdown/{canvasBranchID} [get]
func (r *parser2Routes) GetCanvasBranchMarkdownFile(c *gin.Context) {
	canvasBranchID := c.Param("canvasBranchID")
	App.Service.StartMdBlockConversionForBranch(utils.Uint64(canvasBranchID))
	branch, _ := queries.App.BranchQuery.GetBranchByID(utils.Uint64(canvasBranchID))
	filePath := fmt.Sprintf("%s-%s.md", url.QueryEscape(branch.CanvasRepository.Name), canvasBranchID)
	c.FileAttachment(filePath, filePath)
	go func() {
		e := os.Remove(filePath)
		if e != nil {
			fmt.Println("Error in removing file", e)
		}
	}()
	return
}

// ImportNotion
// @Summary 	Import notion files via zip
// @Description
// @Tags		Parser
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		file		formData file false "Zip file file"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/parser/import-notion [post]
func (r *parser2Routes) ImportNotion(c *gin.Context) {
	var body *ImportNotionValidator
	if err := apiutil.Bind(c, &body); err != nil {
		fmt.Println("error here", err)
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}

	file, err := body.File.Open()
	if err != nil {
		return
	}
	authUser, _ := r.GetLoggedInUser(c)
	studioID, _ := r.GetStudioId(c)
	importTask := ImportTask{}
	importTask.File = body.File
	importTask.StudioID = studioID
	importTask.User = authUser
	go App.Service.NotionImportZipHandler(file, body.File.Size, authUser, studioID)
	//bodyStr, _ := json.Marshal(importTask)
	//apiClient.AddToQueue(apiClient.NotionImportHandler, bodyStr, apiClient.DEFAULT, apiClient.CommonRetry)
	response.RenderSuccessResponse(c, "Successfully started notion import")
	return
}

// ImportFile
// @Summary 	Import single file
// @Description
// @Tags		Parser
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		file		formData file false "file"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/parser/import-file [post]
func (r *parser2Routes) ImportFile(c *gin.Context) {
	var body *ImportNotionValidator
	if err := apiutil.Bind(c, &body); err != nil {
		fmt.Println("error on parsing body", err)
		response.RenderErrorResponse(c, response.BadRequestError(c.Request.Context()))
		return
	}
	authUser, _ := r.GetLoggedInUser(c)
	studioID, _ := r.GetStudioId(c)
	importTask := ImportTask{}
	importTask.File = body.File
	importTask.StudioID = studioID
	importTask.User = authUser
	go App.Service.NotionImportFileHandler(body.File, authUser, studioID)
	response.RenderSuccessResponse(c, "Successfully imported")
	return
}

func Mdplain(g *gin.Context) {
	//649
	// 307
	Parser2{}.PlainText(646)
	g.JSON(200, gin.H{
		"done": "Converted",
	})
}

func (r *parser2Routes) ExportStudio(c *gin.Context) {
	
	studioId, _ := r.GetStudioId(c)
	email := c.Query("email")
	response.RenderSuccessResponse(c, "Exporting started you will get email")
	go func() {
		zipFileName := App.Service.StartExportingStudio(studioId, email)
		e := os.Remove(zipFileName)
		if e != nil {
			fmt.Println("Error in removing file", e)
		}
	}()
	return
}
