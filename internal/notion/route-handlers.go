package notion

import (
	"github.com/gin-gonic/gin"
)

func (r *notionRoutes) ImportNotion(c *gin.Context) {
	// PostImportNotion we get from the request body.
	/*
		1. Expose an API to upload the zip file of notion.
		2. Create a table `Imports` and store the studioID & file metadata for refrence. (Not mandatory I think, but in old code they are doing this.)
		3. Start the extracting zip file and processing it. In old code nodejs is used for extracting the zip file and
			creating the blocks structure. https://gitlab.com/-/ide/project/phonepost/bip-conversion-lambda/tree/master/-/index.js/#L303

			* Need to do some R&D on how to extract zip files in Go lang.
		4. Create the pages and add blocks to it.
		5. Pages should be published here.
		6. We need to decide inside which collection we need to add this imported canvases?
	*/
}
