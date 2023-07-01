package canvasbranch

import (
	"github.com/gin-gonic/gin"
	ar "gitlab.com/phonepost/bip-be-platform/internal/accessrequest"
	"gitlab.com/phonepost/bip-be-platform/internal/bat"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
)

func (a CanvasBranchApp) RegisterRoutes(r *gin.RouterGroup) {
	App.Routes = *r.Group("canvas-branch")
	{
		// get branch metadata
		// CANVAS_BRANCH_VIEW
		App.Routes.GET("/:canvasBranchID", App.RouteHandler.Get)
		App.Routes.GET("/repo/:canvasBranchID", App.RouteHandler.GetRepoBranch)
		// List invited users (pending)
		App.Routes.GET("/:canvasBranchID/invited", App.RouteHandler.Invited)
		// Create branch metadata: Refactor
		// CANVAS_BRANCH_EDIT should be have on the frombranch
		App.Routes.POST("/create", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.Create)
		// Deletes a Branch
		// CANVAS_BRANCH_MANAGE_CONTENT
		App.Routes.DELETE("/:canvasBranchID", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.Delete)
		// CANVAS_BRANCH_VIEW
		App.Routes.GET("/attributions/:canvasBranchID", App.RouteHandler.Attributions)
		// CANVAS_BRANCH_CREATE_MERGE_REQUEST
		App.Routes.GET("/:canvasBranchID/diffblocks", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.DiffBlocks)
		App.Routes.GET("/:canvasBranchID/last-updated", App.RouteHandler.GetBranchLastUpdated)

		// ------------ Branch Ops ------------------------
		branchOpsRoutes := App.Routes.Group("/branch-ops")
		// get all API to get All Branches which are Unpublished or Rough Branch API (Per studio/ User ID)
		// logged in check & response should have the CANVAS_BRANCH_VIEW permission.
		// Returns only the drafts of a specific user. So we don't need to add permission in response.
		branchOpsRoutes.GET("/drafts", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.RoughAndUnpublishedByStudio)
		// Update visibility
		// CANVAS_BRANCH_MANAGE_PERMS
		branchOpsRoutes.POST("/:canvasBranchID/visibility", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.UpdateVisibility)
		// Navigation visibility
		// no logged in check
		branchOpsRoutes.POST("/nav/get-branches", App.RouteHandler.GetCanvasBranches)
		// CANVAS_BRANCH_VIEW
		// Returns the nav based on user permission. so no need to check for permission upfront.
		branchOpsRoutes.GET("/nav/:canvasBranchID/root", App.RouteHandler.GetCanvasBranchesRoots)
		branchOpsRoutes.GET("/blocks/:canvasBranchID/node", App.RouteHandler.GetCanvasBranchesNodes)

		// response should have the CANVAS_BRANCH_VIEW permission
		// Question: Returning all the canvas repos and collections by adding user permission to it
		branchOpsRoutes.POST("/nav/search", App.RouteHandler.SearchCanvasBranches)
		// We need to invite a user with Email.
		// CANVAS_BRANCH_MANAGE_PERMS
		branchOpsRoutes.POST("/:canvasBranchID/create-access-token", bat.App.RouteHandler.CreateAccessToken)
		// no permission check here
		branchOpsRoutes.GET("/get-access-token-detail/:code", bat.App.RouteHandler.GetAccessTokenDetail)
		// CANVAS_BRANCH_MANAGE_PERMS
		branchOpsRoutes.DELETE("/delete-token/:code", bat.App.RouteHandler.DeleteBranchAccessToken)
		// no permission check
		branchOpsRoutes.POST("/:canvasBranchID/join/:code", bat.App.RouteHandler.JoinCurrentUserToStudioAndBranchWithToken)
		// Create a rough branch from a Parent Branch for Edits.
		// CANVAS_BRANCH_EDIT
		branchOpsRoutes.POST("/:canvasBranchID/rough-branch", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.BuildRoughBranch)
		// Get git history for a Branch
		// CANVAS_BRANCH_VIEW
		branchOpsRoutes.GET("/:canvasBranchID/history", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.BranchHistory)
		// We will get list of emails to be added to a branch.

		branchOpsRoutes.POST("/:canvasBranchID/invite-via-emails", bat.App.RouteHandler.InviteViaEmail)

		// ------------ Blocks API's ------------------------
		// Create or Update or Delete Blocks on a Branch
		blockRoutes := App.Routes.Group("/:canvasBranchID/blocks")
		// CANVAS_BRANCH_ADD_COMMENT
		blockRoutes.POST("/associations", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.AssociatesBlock)
		// CANVAS_BRANCH_EDIT
		// check for unpublished or rough.
		blockRoutes.POST("/", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.UpdateBlocks) // This will create/update blocks on a Branch
		// CANVAS_BRANCH_VIEW
		blockRoutes.GET("/", App.RouteHandler.GetBlocks)
		//blockRoutes.GET("/", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.BlockHistoryByCommitID)
		// CANVAS_BRANCH_VIEW
		blockRoutes.GET("/:commitID/blocks-history", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.BlockHistoryByCommitID)

		// ------------ Merge Request Actions  ------------------------

		mergeRoutes := App.Routes.Group("/:canvasBranchID/merge-request")
		/*
			Merge Request Flow
			- Create a Request
			- List All Requests
			- Accept (Partial Request)
			- Reject Merge Request
			- Delete/Cancel Request
			- Get Merge Request with BLOCKS and Branch from GIT
		*/
		// Allows : merge=true // 19 Seconds
		// CANVAS_BRANCH_CREATE_MERGE_REQUEST
		mergeRoutes.POST("/create", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.CreateMergeRequest) // Done
		// CANVAS_BRANCH_MANAGE_MERGE_REQUESTS else
		// creator of merge request
		mergeRoutes.GET("/list", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.ListMergeRequest) // Done
		// CANVAS_BRANCH_MANAGE_MERGE_REQUESTS
		mergeRoutes.POST("/:mergeRequestID/merge-accept", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.MergeRequestAcceptPartial) // Done
		mergeRoutes.POST("/:mergeRequestID/reject", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.MergeRequestRejected)            // DOne
		// only creator of merge request
		mergeRoutes.POST("/:mergeRequestID/delete", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.DeleteMergeRequest) // Done
		// This is a special route as requested by the FE as we don't need the Branch ID here.
		// creator of merge request or CANVAS_BRANCH_MANAGE_MERGE_REQUESTS
		App.Routes.GET("/merge-request/:mergeRequestID/response", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.MergeResponseRequest)
		//mergeRoutes.GET("/:mergeRequestID/response", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.MergeResponseRequest)

		// ------------ Publish Request Actions  ------------------------

		prRoutes := App.Routes.Group("/:canvasBranchID/publish-request")
		// @todo later This should be in the response user perm CANVAS_BRANCH_MANAGE_PUBLISH_REQUESTS or COLLECTION_MANAGE_PUBLISH_REQUEST
		prRoutes.GET("/list", App.RouteHandler.ListPublishRequests)
		// creator of canvas repo
		prRoutes.POST("/init", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.InitPr)
		// Immediate parent having this CANVAS_BRANCH_MANAGE_PUBLISH_REQUESTS or COLLECTION_MANAGE_PUBLISH_REQUEST
		prRoutes.POST("/:publishRequestID/manage", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.ManagePR) // Check perms

		prRoutes.DELETE("/:publishRequestID/delete", middlewares.TokenAuthorizationMiddleware(), App.RouteHandler.DeletePR) // Check perms

		// ------------ Access Request Actions  ------------------------
		arRoutes := App.Routes.Group("/:canvasBranchID/access-request")
		// Anyone loggedIn user who doesn't have access
		arRoutes.POST("/create", middlewares.TokenAuthorizationMiddleware(), ar.App.RouteHandler.CreateAccessRequest)
		// CANVAS_BRANCH_MANAGE_PERMS
		arRoutes.GET("/list", middlewares.TokenAuthorizationMiddleware(), ar.App.RouteHandler.ListAccessRequest)
		// CANVAS_BRANCH_MANAGE_PERMS
		arRoutes.POST("/:accessRequestID/manage", middlewares.TokenAuthorizationMiddleware(), ar.App.RouteHandler.ManageAccessRequest)

	}
}
