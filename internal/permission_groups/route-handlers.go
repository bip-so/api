package permissiongroup

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	PermissionGroupRouteHandler permissionGroupRoutes
)

type permissionGroupRoutes struct{}

// Get User details
// @Summary 	Get Studio Permissions Schema
// @Tags		PermissionGroup
// @Produce 	json
// @Success 	200 		{object} 	PermissionsSchemaResponse
// @Router 		/v1/permissions-schema/studio/schema [get]
func (r permissionGroupRoutes) StudioPermsSchema(c *gin.Context) {
	pr := PermissionsSchemaResponse{}
	pr.Version = PGSchemaVersion
	pr.Group = PGTYPESTUDIO

	var prgnone PermissionsTemplate
	var prgadmin PermissionsTemplate

	prgnone = PGStudioNone()
	prgadmin = PGStudioAdmin()

	pr.Permissions = []PermissionsTemplate{prgnone, prgadmin}
	c.JSON(http.StatusOK, pr)
	return

}

// Get User details
// @Summary 	Get Collections Permissions Schema
// @Tags		PermissionGroup
// @Produce 	json
// @Success 	200 		{object} 	PermissionsSchemaResponse
// @Router 		/v1/permissions-schema/collection/schema [get]
func (r permissionGroupRoutes) CollectionPermsSchema(c *gin.Context) {
	pr := PermissionsSchemaResponse{}
	pr.Version = PGSchemaVersion
	pr.Group = PGTYPECOLLECTIION
	var prNone, prView, prComment, prEdit, prModerate, prViewMetadata PermissionsTemplate

	prNone = PGCollectionNone()
	prView = PGCollectionView()
	prComment = PGCollectionComment()
	prEdit = PGCollectionEdit()
	prModerate = PGCollectionModerate()
	prViewMetadata = PGCollectionViewMetadata()

	pr.Permissions = []PermissionsTemplate{prNone, prView, prComment, prEdit, prModerate, prViewMetadata}
	c.JSON(http.StatusOK, pr)
	return

}

// Get User details
// @Summary 	Get Canvas Permissions Schema
// @Tags		PermissionGroup
// @Produce 	json
// @Success 	200 		{object} 	PermissionsSchemaResponse
// @Router 		/v1/permissions-schema/canvasBranch/schema [get]
func (r permissionGroupRoutes) CanvasBranchPermsSchema(c *gin.Context) {
	pr := PermissionsSchemaResponse{}
	pr.Version = PGSchemaVersion
	pr.Group = PGTYPECANVAS
	var prNone, prView, prComment, prEdit, prModerate, prViewMetadata PermissionsTemplate

	prNone = PGCanvasNone()
	prViewMetadata = PGCanvasViewMetaData()
	prView = PGCanvasView()
	prComment = PGCanvasComment()
	prEdit = PGCanvasEdit()
	prModerate = PGCanvasModerate()

	pr.Permissions = []PermissionsTemplate{prNone, prView, prComment, prEdit, prModerate, prViewMetadata}
	c.JSON(http.StatusOK, pr)
	return
}
