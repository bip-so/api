package bat

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/canvasbranchpermissions"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
)

// @Summary 	InviteViaEmail: Invite People to Canvas who may or may not have accounts
// @Tags		CanvasBranch
// @Accept      json
// @Produce     json
// @Param 		body 	body 		CreateEmailInvite true "Data"
// @Router 		/v1/canvas-branch/branch-ops/:canvasBranchID/invite-via-emails [POST]
func (r batRoutes) InviteViaEmail(c *gin.Context) {
	var body CreateEmailInvite
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	studioID, errCantGetStudio := r.GetStudioId(c)
	if errCantGetStudio != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Studio ID is required",
		})
		return
	}

	// Get the canvas branch
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)

	// permission check (Commented for now)
	/*
		authUserID := r.GetLoggedInUserId(c)
		if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, canvasBranchID, permissiongroup.CANVAS_BRANCH_MANAGE_PERMS); err != nil || !hasPermission {
			response.RenderPermissionError(c)
			return
		}
	*/

	user, errGettingUser := r.GetLoggedInUser(c)
	if errGettingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errGettingUser.Error(),
		})
		return
	}

	errCreatingInvite := App.Controller.InviteViaEmailController(body, user, canvasBranchID, studioID)
	if errCreatingInvite != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errGettingUser.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "We have successfully invited.",
	})
	return
}

// @Summary 	Create a Branch Token with Access
// @Tags		CanvasBranch
// @Accept      json
// @Produce     json
// @Param 		body 	body 		CreateAccessTokenPost true "Data"
// @Router 		/v1/canvas-branch/branch-ops/:canvasBranchID/create-access-token [POST]
func (r batRoutes) CreateAccessToken(c *gin.Context) {
	var body CreateAccessTokenPost
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get the canvas branch
	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)

	// permission check
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, canvasBranchID, permissiongroup.CANVAS_BRANCH_MANAGE_PERMS); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	user, errGettingUser := r.GetLoggedInUser(c)
	if errGettingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errGettingUser.Error(),
		})
		return
	}
	branchAccessToken, err := App.Service.CreateBranchAccessToken(user, canvasBranchID, body.PermissionGroup)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Token is created successfully",
		"inviteCode":      branchAccessToken.InviteCode,
		"permissionGroup": branchAccessToken.PermissionGroup,
	})
	return
}

// @Summary 	Join user with a given branch Token with Access
// @Tags		CanvasBranch
// @Accept      json
// @Produce     json
// @Param 		body 	body 		PlaceHolder true "Data"
// @Router 		/v1/canvas-branch/branch-ops/{canvasBranchID}/join/{code} [POST]
func (r batRoutes) JoinCurrentUserToStudioAndBranchWithToken(c *gin.Context) {
	var body PlaceHolder
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	studioID, errGettingStudioID := r.GetStudioId(c)
	if errGettingStudioID != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errGettingStudioID.Error(),
		})
		return
	}

	canvasBranchIDStr := c.Param("canvasBranchID")
	canvasBranchID, _ := strconv.ParseUint(canvasBranchIDStr, 10, 64)
	code := c.Param("code")
	user, errGettingUser := r.GetLoggedInUser(c)
	if errGettingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errGettingUser.Error(),
		})
		return
	}
	branchAccessToken, errBat := App.Service.GetBranchAccessToken(code)
	if errBat != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errBat.Error(),
		})
		return
	}

	errorPerformingAction := App.Service.JoinUserStudioBranch(studioID, canvasBranchID, user, branchAccessToken)
	if errorPerformingAction != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errorPerformingAction.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User has been added to this studio and branch.",
	})
	return
}

// Get Details of a short code
// @Summary 	Get Details of a short code, by invite code.
// @Tags		CanvasBranch
// @Accept       json
// @Produce      json
// @Router 		/v1/canvas-branch/branch-ops/get-access-token-detail/:code [get]
func (r batRoutes) GetAccessTokenDetail(c *gin.Context) {
	code := c.Param("code")
	branchAccessToken, err := App.Service.GetBranchAccessToken(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//branch, _ := App.Repo.Get(map[string]interface{}{"id": branchAccessToken.BranchID})
	branch, _ := queries.App.BranchQuery.GetBranchWithRepoAndStudio(branchAccessToken.BranchID)

	// Automagically add the Logged-in user to the Branch and Studio.
	user, _ := r.GetLoggedInUser(c)
	if user != nil {
		// This is conditional Flow
		studioID, _ := r.GetStudioId(c)
		// updating permission of user and adding user to the studio if not present.
		member, _ := queries.App.MemberQuery.GetMember(map[string]interface{}{"user_id": user.ID, "studio_id": studioID})
		if member == nil {
			_ = App.Service.JoinUserStudioBranch(studioID, branchAccessToken.BranchID, user, branchAccessToken)
			member, _ = queries.App.MemberQuery.GetMember(map[string]interface{}{"user_id": user.ID, "studio_id": studioID})
		}
		permData := canvasbranchpermissions.NewCanvasBranchPermissionCreatePost{
			CollectionId:       branch.CanvasRepository.CollectionID,
			CanvasBranchId:     branch.ID,
			CanvasRepositoryID: branch.CanvasRepositoryID,
			PermGroup:          branchAccessToken.PermissionGroup,
			MemberID:           member.ID,
		}
		if branch.CanvasRepository.ParentCanvasRepositoryID != nil {
			permData.CbpParentCanvasRepositoryID = *branch.CanvasRepository.ParentCanvasRepositoryID
		}
		_, err = canvasbranchpermissions.App.Controller.CreateCanvasBranchPermissionController(permData, studioID, user.ID, "false")
		if err != nil {
			fmt.Println("error while updating canvas branch perm", err)
		}
		c.JSON(http.StatusOK, gin.H{
			"inviteCode":      branchAccessToken.InviteCode,
			"permissionGroup": branchAccessToken.PermissionGroup,
			"createdById":     branchAccessToken.CreatedByID,
			"isActive":        branchAccessToken.IsActive,
			"branchId":        branchAccessToken.BranchID,
			"repoId":          branch.CanvasRepository.ID,
			"repoKey":         branch.CanvasRepository.Key,
			"user_added":      true,
			"studioHandle":    branch.CanvasRepository.Studio.Handle,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"inviteCode":      branchAccessToken.InviteCode,
		"permissionGroup": branchAccessToken.PermissionGroup,
		"createdById":     branchAccessToken.CreatedByID,
		"isActive":        branchAccessToken.IsActive,
		"branchId":        branchAccessToken.BranchID,
		"repoId":          branch.CanvasRepository.ID,
		"repoKey":         branch.CanvasRepository.Key,
		"user_added":      false,
		"studioHandle":    branch.CanvasRepository.Studio.Handle,
	})
	return
}

// @Summary 	Delete a Branch Acceess Token By Key
// @Tags		CanvasBranch
// @Accept       json
// @Produce      json
// @Param       code  path    string  true  "code"
// @Router 		/v1/canvas-branch/branch-ops/delete-token/{code} [delete]
func (r batRoutes) DeleteBranchAccessToken(c *gin.Context) {
	code := c.Param("code")
	branchAccessToken, err := App.Service.GetBranchAccessToken(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// permission check
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(authUserID, branchAccessToken.BranchID, permissiongroup.CANVAS_BRANCH_MANAGE_PERMS); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	errDelete := App.Service.DeleteBAT(branchAccessToken.ID)
	if errDelete != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errDelete.Error(),
		})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{
		"message": "Resource deleted successfully",
	})
	return
}
