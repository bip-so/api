package canvasrepo

import (
	"encoding/json"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/internal/workflows"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/apiClient"
	"net/http"
	"strconv"
	"time"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/internal/studio"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/middlewares"
	"gitlab.com/phonepost/bip-be-platform/pkg/api/response"
)

// @Summary 	Init Canvas Repo with a Fresh Branch
// @Tags		CanvasRepo
// @Security 	bearerAuth
// @Router 		/v1/canvas-repo/init  [post]
// @Param 		body body 		InitCanvasRepoPost true "Init Canvas Repo"
func (r *canvasRepoRoutes) Init(c *gin.Context) {
	var body InitCanvasRepoPost

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Required should be in header
	studioID, _ := r.RouteHelper.GetStudioId(c)
	userID := r.GetLoggedInUserId(c)
	user, errGettingUser := r.GetLoggedInUser(c)
	if errGettingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errGettingUser.Error(),
		})
		return
	}

	if body.ParentCanvasRepositoryID != 0 {
		canvasRepo, _ := App.Repo.Get(map[string]interface{}{"id": body.ParentCanvasRepositoryID})
		if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(user.ID, *canvasRepo.DefaultBranchID, permissiongroup.CANVAS_BRANCH_VIEW_METADATA); err != nil || !hasPermission {
			response.RenderPermissionError(c)
			return
		}
	} else {
		if hasPermission, err := permissions.App.Service.CanUserDoThisOnCollection(user.ID, studioID, body.CollectionID, permissiongroup.COLLECTION_VIEW_METADATA); err != nil || !hasPermission {
			response.RenderPermissionError(c)
			return
		}
	}

	// For debug remove
	fmt.Println(studioID)
	fmt.Println(userID)

	// Adding Studio Header Chcekr
	_, status := middlewares.StudioHeaderRequiredCheck(c)
	if !status {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": "Studio Header Missing",
		})
		return
	}

	repo, err := workflows.WorkflowHelperInitCanvasRepo(workflows.InitCanvasRepoPost{
		CollectionID:             body.CollectionID,
		Name:                     body.Name,
		Icon:                     body.Icon,
		Position:                 body.Position,
		ParentCanvasRepositoryID: body.ParentCanvasRepositoryID,
	}, userID, studioID, *user)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	response.RenderResponse(c, SerializeDefaultCanvasRepo(repo))
}

// GetAllCanvas Get Canvas Repos
// @Summary 	Get all Canvas Related to collection
// @Description
// @Tags		CanvasRepo
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		body 		body 		GetAllCanvasValidator true "Get all Canvas"
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/canvas-repo/get [post]
func (r *canvasRepoRoutes) GetAllCanvas(c *gin.Context) {
	var body GetAllCanvasValidator
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, _ := App.RouteHandler.GetLoggedInUser(c)
	studioId, _ := App.RouteHandler.GetStudioId(c)
	public := c.Query("public")
	var canvasRepoViews *[]CanvasRepoDefaultSerializer
	var err error

	// Checking if the request is made by loggedIn user or not
	// If user == nil we trigger the Anonymous flow to get the canvas or vice-versa
	if user == nil || public == "true" {
		canvasRepoViews, err = App.Controller.AnonymousGetAllCanvasController(&body)
	} else {
		_, err = permissions.App.Repo.GetMember(map[string]interface{}{"studio_id": studioId, "user_id": user.ID})
		if err == gorm.ErrRecordNotFound {
			canvasRepoViews, err = App.Controller.AnonymousGetAllCanvasController(&body)
		} else {
			canvasRepoViews, err = App.Controller.AuthUserGetAllCanvasController(&body, user, studioId)
		}
	}
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	response.RenderResponse(c, canvasRepoViews)
}

// Move Canvas Position
// @Summary 	Move Canvas Position
// @Tags		CanvasRepo
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		body 		body 		MoveCanvasRepoPost true "Move Canvas"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/canvas-repo/move [post]
func (r *canvasRepoRoutes) MoveCanvas(c *gin.Context) {
	/*
		Some cases of move API:

		Mandatory Fields:
			- canvas_repo_id
			- future_position

		Canvas moving from collection1 canvas1-pos2(parent as canvas) -> collection2 canvas2-pos4(Parent as canvas) (or)
		Canvas moving from collection1 canvas-pos2(parent as collection) -> collection2 canvas2-pos4(Parent as canvas) from_parent_canvas_repository_id will be 0 here
			- to_collection_id
			- to_parent_canvas_repository_id

		canvas moving from collection1 canvas(parent as collection) -> collection1 canvas(parent as collection)
			-

		canvas moving from collection1 canvas(parent as collection) -> collection2 canvas(parent as collection)
			- to_collection_id

		canvas moving from collection1 canvas1-pos2(parent as canvas) -> collection1 canvas1-pos4(parent as canvas)
		canvas moving from collection1 canvas1-pos2(parent as canvas) -> collection1 canvas2-pos4(parent as canvas)
			- to_parent_canvas_repository_id

		canvas moving from collection1 canvas1-pos2(parent as canvas) -> collection1 canvas-pos4(parent as collection)
			-

		canvas moving from collection1 canvas(parent as collection) -> collection1 canvas(parent as canvas)
			- to_parent_canvas_repository_id
	*/
	var body MoveCanvasRepoPost
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID := r.GetLoggedInUserId(c)
	studioID, _ := r.GetStudioId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnStudio(userID, studioID, permissiongroup.STUDIO_CHANGE_CANVAS_COLLECTION_POSITION); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}
	canvasRepoInstance, _ := queries.App.RepoQuery.GetRepo(map[string]interface{}{"id": body.CanvasRepoID})
	preCollectionID := canvasRepoInstance.CollectionID

	var err error
	if body.ToCollectionID != 0 && body.ToParentCanvasRepositoryID != 0 {
		permissions.App.Service.UpdateParentHasCanvasRepoOnRemove(*canvasRepoInstance.DefaultBranchID, *canvasRepoInstance.DefaultBranchID)
		err = App.Controller.moveCanvasBetweenCanvasAndCollection(&body)
		if err != nil {
			response.RenderCustomErrorResponse(c, err)
			return
		}
		App.Service.SendCollectionTreeToDiscord(preCollectionID)
	} else if body.ToCollectionID != 0 {
		permissions.App.Service.UpdateParentHasCanvasRepoOnRemove(*canvasRepoInstance.DefaultBranchID, *canvasRepoInstance.DefaultBranchID)
		err = App.Controller.moveCanvasBetweenCollections(&body)
		App.Service.SendCollectionTreeToDiscord(preCollectionID)
	} else if body.ToParentCanvasRepositoryID != 0 {
		permissions.App.Service.UpdateParentHasCanvasRepoOnRemove(*canvasRepoInstance.DefaultBranchID, *canvasRepoInstance.DefaultBranchID)
		err = App.Controller.moveCanvasBetweenCanvas(&body)
	} else {
		err = App.Controller.moveCanvasBetweenSameCollectionsAndCanvas(&body)
	}

	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	go func() {
		fmt.Println("collection id step 2", canvasRepoInstance.Name, canvasRepoInstance.CollectionID, canvasRepoInstance.DefaultBranch.PublicAccess)
		if canvasRepoInstance.DefaultBranch.PublicAccess == models.PRIVATE {
			permissions.App.Service.UpdateParentHasCanvasRepoOnPrivate(*canvasRepoInstance.DefaultBranchID)
		} else {
			permissions.App.Service.UpdateParentHasCanvasRepoOnPublic(*canvasRepoInstance.DefaultBranchID)
		}
		if body.ToCollectionID != 0 {
			payload, _ := json.Marshal(map[string]uint64{"collectionId": body.ToCollectionID})
			apiClient.AddToQueue(apiClient.UpdateDiscordTreeMessage, payload, apiClient.DEFAULT, apiClient.CommonRetry)
		} else {
			payload, _ := json.Marshal(map[string]uint64{"collectionId": canvasRepoInstance.CollectionID})
			apiClient.AddToQueue(apiClient.UpdateDiscordTreeMessage, payload, apiClient.DEFAULT, apiClient.CommonRetry)
		}
	}()
	response.RenderResponse(c, "CanvasRepo Moved Succesfully")
}

// @Summary 	Update Canvas - Note this only updates Name or Icon for a Repo use Move API for positiion
// @Tags		CanvasRepo
// @Router 		/v1/canvas-repo/{canvasRepoID}  [patch]
// @Param 		body body 		UpdateCanvasRepoPost true "UpdateCanvas Repo"
func (r *canvasRepoRoutes) Update(c *gin.Context) {
	var body UpdateCanvasRepoPost
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Required should be in header
	studioID, _ := r.GetStudioId(c)
	userID := r.GetLoggedInUserId(c)

	canvasRepoIDStr := c.Param("canvasRepoID")
	canvasRepoID, _ := strconv.ParseUint(canvasRepoIDStr, 10, 64)

	canvasRepo, _ := App.Repo.Get(map[string]interface{}{"id": canvasRepoID})
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(userID, *canvasRepo.DefaultBranchID, permissiongroup.CANVAS_BRANCH_EDIT_NAME); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	// Todo: For debug remove
	fmt.Println(studioID)
	fmt.Println(userID)
	// Adding Studio Header Chcekr
	_, status := middlewares.StudioHeaderRequiredCheck(c)
	if !status {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": "Studio Header Missing",
		})
		return
	}

	flag := App.Controller.UpdateCanvasRepo(canvasRepoID, body, userID)
	go func() {
		payload, _ := json.Marshal(map[string]uint64{"collectionId": canvasRepo.CollectionID})
		apiClient.AddToQueue(apiClient.UpdateDiscordTreeMessage, payload, apiClient.DEFAULT, apiClient.CommonRetry)
	}()
	if flag {
		c.JSON(http.StatusOK, gin.H{
			"message": "Done",
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Something went wrong.",
		})
		return
	}
}

// Create Language Canvas Repo
// @Summary 	Create Language Canvas Repo
// @Description
// @Tags		CanvasRepo
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		body 		body 		CreateLanguageValidator true "create language data"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/canvas-repo/create-language [post]
func (r *canvasRepoRoutes) CreateLanguage(c *gin.Context) {

	user, _ := r.GetLoggedInUser(c)

	var body CreateLanguageValidator

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID := r.GetLoggedInUserId(c)
	canvasRepo, _ := App.Repo.Get(map[string]interface{}{"id": body.CanvasRepositoryID})
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(userID, *canvasRepo.DefaultBranchID, permissiongroup.CANVAS_BRANCH_VIEW); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	canvasRepos, err, duplicateLanguageCodes := App.Controller.CreateLanguageCanvasRepo(&body, user)
	if err != nil {
		response.RenderErrorResponse(c, response.ServerError(c.Request.Context(), err))
		return
	}
	fmt.Println(canvasRepos, err, duplicateLanguageCodes)
	var canvasRepoSerialized *[]CanvasRepoDefaultSerializer
	var canvasReposData []CanvasRepoDefaultSerializer
	if canvasRepos != nil {
		canvasRepoSerialized = MultiSerializeDefaultCanvasRepo(canvasRepos)
		for _, lRepo := range *canvasRepoSerialized {
			lRepo.DefaultBranch.Permission = models.PGCanvasModerateSysName
			canvasReposData = append(canvasReposData, lRepo)
		}
	}
	resp := map[string]interface{}{
		"data":                   canvasReposData,
		"duplicateLanguageCodes": duplicateLanguageCodes,
	}
	response.RenderCustomResponse(c, resp)
}

// GetCanvas Get Canvas Repos By Key
// @Summary 	Get One Canvas using Key
// @Description
// @Tags		CanvasRepo
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/canvas-repo/ [get]
func (r *canvasRepoRoutes) GetCanvas(c *gin.Context) {
	key := c.Query("key")
	inviteCode := c.Query("inviteCode")
	if key == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No key found ",
		})
		return
	}
	canvasRepoInstance, _ := App.Repo.GetCanvasRepoByKey(key)
	if canvasRepoInstance == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No Repo Found.",
		})
		return
	}

	user, _ := App.RouteHandler.GetLoggedInUser(c)
	// User not found send a graceful 403 and publicAccess of the Default branch is private
	if user == nil && inviteCode != "" {
		canvasBranchAccessToken, err := permissions.App.Service.CheckBranchAccessToken(inviteCode, key)
		if err != nil || canvasBranchAccessToken == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Anonymous user does not have permissions to view this.",
			})
			return
		}
		canvasRepoView, err := App.Controller.AnonymousGetOneCanvasByKeyController(key)
		if err != nil {
			response.RenderCustomErrorResponse(c, err)
			return
		}
		canvasRepoView.DefaultBranch.Permission = canvasBranchAccessToken.PermissionGroup
		response.RenderResponse(c, canvasRepoView)
		return
	} else if user == nil && canvasRepoInstance.DefaultBranch.PublicAccess == models.CANVAS_BRANCH_PUBLIC_ACCESS_PRIVATE {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Anonymous user does not have permissions to view this.",
		})
		return
	}

	userID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(userID, *canvasRepoInstance.DefaultBranchID, permissiongroup.CANVAS_BRANCH_VIEW_METADATA); err != nil || !hasPermission {
		// We are adding access_requested is User already have Access Request
		exists := App.Repo.AcceesRequestExistsSimple(*canvasRepoInstance.DefaultBranchID, user.ID)
		c.JSON(http.StatusForbidden, gin.H{
			"error":            "User does not have permissions to view branch",
			"access_requested": exists,
		})
		return
	}

	var canvasRepoView *CanvasRepoDefaultSerializer
	var err error
	// Checking if the request is made by loggedIn user or not
	// If user == nil we trigger the Anonymous flow to get the canvas or vice-versa
	if user == nil {
		canvasRepoView, err = App.Controller.AnonymousGetOneCanvasByKeyController(key)
	} else {
		canvasRepoView, err = App.Controller.AuthGetOneCanvasByKeyController(key, user.ID)
	}
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	// Adds to the Repo Set

	response.RenderResponse(c, canvasRepoView)
}

// @Summary 	Delete Canvas Repo
// @Tags		CanvasRepo
// @Router 		/v1/canvas-repo/{canvasRepoID}  [delete]
func (r *canvasRepoRoutes) DeleteCanvasRepo(c *gin.Context) {
	canvasRepoIDStr := c.Param("canvasRepoID")
	canvasRepoID, _ := strconv.ParseUint(canvasRepoIDStr, 10, 64)

	// Permission check
	userID := r.GetLoggedInUserId(c)
	canvasRepo, _ := App.Repo.Get(map[string]interface{}{"id": canvasRepoID})
	if canvasRepo.CreatedByID != userID {
		if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(userID, *canvasRepo.DefaultBranchID, permissiongroup.CANVAS_BRANCH_DELETE); err != nil || !hasPermission {
			response.RenderPermissionError(c)
			return
		}
	}

	// TRodo : Refactor move to service ask NR also do we delete or sioft delete
	canvas, _ := queries.App.RepoQuery.GetRepo(map[string]interface{}{"id": canvasRepoID})
	var repo models.CanvasRepository
	err := App.Repo.Manager.HardDeleteByID(repo.TableName(), canvasRepoID)
	if err == nil {
		go func() {
			App.Repo.RearrangeTheOldCanvasRepo(canvas)
			canvasData, _ := json.Marshal(canvas)
			App.Repo.kafka.Publish(configs.KAFKA_TOPICS_DELETED_CANVAS, strconv.FormatUint(canvasRepoID, 10), canvasData)
		}()
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Something went wrong.",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Deleted",
		})
		return
	}

}

// Create
// @Summary 	Create Canvas Repo with a Default Branch
// @Tags		CanvasRepo
// @Security 	bearerAuth
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Router 		/v1/canvas-repo/create  [post]
// @Param 		body body 		NewCanvasRepoPost true "Create Canvas Repo"
func (r *canvasRepoRoutes) Create(c *gin.Context) {
	var body NewCanvasRepoPost
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Required should be in header
	studioID, _ := r.RouteHelper.GetStudioId(c)
	userID := r.GetLoggedInUserId(c)
	user, errGettingUser := r.GetLoggedInUser(c)

	if errGettingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errGettingUser.Error(),
		})
		return
	}

	var parentPublicAccess string
	// permission check
	if body.ParentCanvasRepositoryID != 0 {
		canvasRepo, _ := queries.App.RepoQuery.GetRepo(map[string]interface{}{"id": body.ParentCanvasRepositoryID})
		if canvasRepo.DefaultBranch.PublicAccess == models.EDIT {
			_, err := queries.App.MemberQuery.GetMember(map[string]interface{}{"user_id": userID, "studio_id": studioID})
			if err == gorm.ErrRecordNotFound {
				studio.App.Controller.JoinStudioController(user, studioID)
			}
		} else if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(user.ID, *canvasRepo.DefaultBranchID, permissiongroup.CANVAS_BRANCH_VIEW_METADATA); err != nil || !hasPermission {
			response.RenderPermissionError(c)
			return
		}
		parentPublicAccess = canvasRepo.DefaultBranch.PublicAccess
	} else {
		collection, _ := App.Repo.GetCollection(map[string]interface{}{"id": body.CollectionID})
		parentPublicAccess = collection.PublicAccess
		if parentPublicAccess == models.EDIT {
			_, err := queries.App.MemberQuery.GetMember(map[string]interface{}{"user_id": userID, "studio_id": studioID})
			if err == gorm.ErrRecordNotFound {
				studio.App.Controller.JoinStudioController(user, studioID)
			}
		} else if hasPermission, err := permissions.App.Service.CanUserDoThisOnCollection(user.ID, studioID, body.CollectionID, permissiongroup.COLLECTION_VIEW_METADATA); err != nil || !hasPermission {
			response.RenderPermissionError(c)
			return
		}
	}
	if parentPublicAccess == "" {
		parentPublicAccess = models.PRIVATE
	}

	// For debug remove
	fmt.Println(studioID)
	fmt.Println(userID)
	// Adding Studio Header Chcekr
	_, status := middlewares.StudioHeaderRequiredCheck(c)
	if !status {
		response.RenderCustomResponse(c, map[string]interface{}{
			"error": "Studio Header Missing",
		})
		return
	}

	/*
		repo, err := shared.WorkflowHelperInitCanvasRepo(shared.InitCanvasRepoPost{
				CollectionID:             body.CollectionID,
				Name:                     body.Name,
				Icon:                     body.Icon,
				Position:                 body.Position,
				ParentCanvasRepositoryID: body.ParentCanvasRepositoryID,
			}, userID, studioID, *user)
	*/

	repo, err := App.Controller.CreateCanvasRepo(body, userID, studioID, *user, parentPublicAccess)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}
	serialized := SerializeDefaultCanvasRepo(repo)
	serialized.DefaultBranch.Permission = models.PGCanvasModerateSysName
	serialized.DefaultBranch.CanPublish = true

	fmt.Println("Create Time Taken")
	defer utils.TimeTrack(time.Now())
	response.RenderResponse(c, serialized)
}

// GetMemberCanvas Get Member Canvas Repos for permissions
// @Summary 	Get member Canvas Related to collection
// @Description
// @Tags		CanvasRepo
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		body 	body 		GetAllCanvasValidator true "Get all Canvas"
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Param       userId  path    string  true  "User Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/canvas-repo/user/{userId} [post]
func (r *canvasRepoRoutes) GetMemberCanvas(c *gin.Context) {
	var body GetAllCanvasValidator
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userIDStr := c.Param("userId")
	userID := utils.Uint64(userIDStr)
	studioId, _ := App.RouteHandler.GetStudioId(c)

	// Permission check
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnStudio(authUserID, studioId, permissiongroup.STUDIO_ADD_REMOVE_USER_TO_ROLE); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}
	var canvasRepoViews *[]CanvasRepoDefaultSerializer
	var err error

	canvasRepoViews, err = App.Controller.MemberCanvasController(body, userID, studioId)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	response.RenderResponse(c, canvasRepoViews)
}

// GetRoleCanvas Get Role Canvas Repos for permissions
// @Summary 	Get role Canvas Related to collection
// @Description
// @Tags		CanvasRepo
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		body 	body 		GetAllCanvasValidator true "Get all Canvas"
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Param       roleId  path    string  true  "Role Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/canvas-repo/role/{roleId} [post]
func (r *canvasRepoRoutes) GetRoleCanvas(c *gin.Context) {
	var body GetAllCanvasValidator
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	roleIDStr := c.Param("roleId")
	roleID := utils.Uint64(roleIDStr)
	studioID, _ := r.GetStudioId(c)

	// Permission check
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnStudio(authUserID, studioID, permissiongroup.STUDIO_ADD_REMOVE_USER_TO_ROLE); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	var canvasRepoViews *[]CanvasRepoDefaultSerializer
	var err error

	canvasRepoViews, err = App.Controller.RoleCanvasController(body, roleID)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	response.RenderResponse(c, canvasRepoViews)
}

// GetNextPrevCanvas Get next and prev canvas of that canvas
// @Summary 	Get next and prev canvas of that canvas
// @Description
// @Tags		CanvasRepo
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Param       canvasRepoId  path    string  true  "canvas Repo Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/canvas-repo/next-prev/{canvasRepoId} [get]
func (r *canvasRepoRoutes) GetNextPrevCanvas(c *gin.Context) {
	userID := r.GetLoggedInUserId(c)
	canvasIDStr := c.Param("canvasRepoId")
	canvasID := utils.Uint64(canvasIDStr)

	nextCanvas, PrevCanvas := App.Service.GetCanvasPrevAndNext(userID, canvasID)

	resp := map[string]interface{}{
		"next": SerializeDefaultCanvasRepoMini(nextCanvas),
		"prev": SerializeDefaultCanvasRepoMini(PrevCanvas),
	}

	response.RenderCustomResponse(c, resp)
}

// UserSearchCanvases Search Canvas Repos
// @Summary 	Search Canvas Repos
// @Description
// @Tags		CanvasRepo
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Param       search   query    string  true  "search"
// @Param       userId  path    string  true  "User Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/canvas-repo/search/user/{userId} [get]
func (r *canvasRepoRoutes) UserSearchCanvases(c *gin.Context) {
	search := c.Query("search")
	studioID, _ := r.GetStudioId(c)
	userIDStr := c.Param("userId")
	userID := utils.Uint64(userIDStr)

	// Permission check
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnStudio(authUserID, studioID, permissiongroup.STUDIO_ADD_REMOVE_USER_TO_ROLE); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	canvasRepoViews, err := App.Controller.UserCanvasSearchController(studioID, userID, search)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	response.RenderResponse(c, canvasRepoViews)
}

// RoleSearchCanvases Search Canvas Repos
// @Summary 	Search Canvas Repos
// @Description
// @Tags		CanvasRepo
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Param       search   query    string  true  "search"
// @Param       roleId  path    string  true  "Role Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/canvas-repo/search/role/{roleId} [get]
func (r *canvasRepoRoutes) RoleSearchCanvases(c *gin.Context) {
	search := c.Query("search")
	studioID, _ := r.GetStudioId(c)

	roleIDStr := c.Param("roleId")
	roleID := utils.Uint64(roleIDStr)

	// Permission check
	authUserID := r.GetLoggedInUserId(c)
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnStudio(authUserID, studioID, permissiongroup.STUDIO_ADD_REMOVE_USER_TO_ROLE); err != nil || !hasPermission {
		response.RenderPermissionError(c)
		return
	}

	canvasRepoViews, err := App.Controller.RoleCanvasSearchController(studioID, roleID, search)
	if err != nil {
		response.RenderCustomErrorResponse(c, err)
		return
	}

	response.RenderResponse(c, canvasRepoViews)
}

// GetDistinctRepoLanguages Get distinct languages of all repos
// @Summary 	Get distinct languages of all repos
// @Description
// @Tags		CanvasRepo
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/canvas-repo/distinct-languages [get]
func (r *canvasRepoRoutes) GetDistinctRepoLanguages(c *gin.Context) {
	studioID, _ := r.GetStudioId(c)
	repos, _ := queries.App.RepoQuery.GetReposLanguages(studioID)
	languages := []string{}
	for _, repo := range repos {
		languages = append(languages, *repo.Language)
	}
	resp := map[string]interface{}{
		"languages": languages,
	}

	response.RenderCustomResponse(c, resp)
}

// GetLangNextPrevCanvas Get Lang next and prev canvas of that canvas
// @Summary 	Get Lang next and prev canvas of that canvas
// @Description
// @Tags		CanvasRepo
// @Security 	bearerAuth
// @Accept 		json
// @Produce 	json
// @Param       bip-studio-id  header    string  true  "Studio Id"
// @Param       canvasRepoId  path    string  true  "canvas Repo Id"
// @Success 	200 		{object} 	response.ApiResponse
// @Failure 	401 		{object} 	response.ApiResponse
// @Router 		/v1/canvas-repo/lang-next-prev/{canvasRepoId} [get]
func (r *canvasRepoRoutes) GetLangNextPrevCanvas(c *gin.Context) {
	userID := r.GetLoggedInUserId(c)
	canvasIDStr := c.Param("canvasRepoId")
	canvasID := utils.Uint64(canvasIDStr)
	language := c.Query("language")

	nextCanvas, PrevCanvas := App.Service.GetLanguageCanvasPrevAndNext(userID, canvasID, language)

	resp := map[string]interface{}{
		"next": SerializeDefaultCanvasRepoMini(nextCanvas),
		"prev": SerializeDefaultCanvasRepoMini(PrevCanvas),
	}

	response.RenderCustomResponse(c, resp)
}
