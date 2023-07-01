package discord

import (
	"database/sql"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"log"
	"regexp"

	"gitlab.com/phonepost/bip-be-platform/internal/canvasrepo"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	studiointegration "gitlab.com/phonepost/bip-be-platform/internal/studio_integration"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"gorm.io/datatypes"
)

const (
	DISCORD_USER_TYPE = "discord_user"
	DISCORD_PROVIDER  = "discord"
)

func FindUsersByDiscordIDs(userIds []string) (users []models.UserSocialAuth, err error) {
	err = postgres.GetDB().Model(&models.UserSocialAuth{}).Where("provider_id IN ?", userIds).Preload("User").Find(&users).Error

	return
}
func GetDiscordStudioIntegration(studioId uint64) (integration []models.StudioIntegration, err error) {
	condition := models.StudioIntegration{StudioID: studioId}

	condition.Type = studiointegration.DISCORD_INTEGRATION_TYPE
	err = postgres.GetDB().Model(&models.StudioIntegration{}).Where(condition).Find(&integration).Error
	return
}
func FindUsersByDiscordID(userId string) (user *models.UserSocialAuth, err error) {
	err = postgres.GetDB().Model(&models.UserSocialAuth{}).Where("provider_id = ?", userId).Preload("User").Find(&user).Error

	return
}

func GetProductIntegrationByDiscordTeamId(teamID string) (integration []models.StudioIntegration, err error) {
	condition := models.StudioIntegration{TeamID: teamID}

	condition.Type = studiointegration.DISCORD_INTEGRATION_TYPE
	err = postgres.GetDB().Model(&models.StudioIntegration{}).Where(condition).Preload("Studio").Find(&integration).Error

	return
}

func NewDiscordUser(userId uint64, providerId string, metadata datatypes.JSON) *models.UserSocialAuth {
	return &models.UserSocialAuth{
		UserID:       userId,
		ProviderName: DISCORD_PROVIDER,
		ProviderID:   providerId,
		Metadata:     metadata,
	}
}

func CreateNewUser(email, password, username, avatarURL string) *models.User {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	processedUsername := reg.ReplaceAllString(username, "")
	user := &models.User{
		Email: sql.NullString{
			String: email,
			Valid:  true,
		},
		Password: password,
		Username: processedUsername,
	}

	_ = queries.App.UserQueries.CreateUser(user)
	return user
}

// CreateCanvasRepo Create New Canvas Repo
func CreateCanvasRepo(name, icon string, userID uint64, studioID uint64, collectionID uint64, position uint, parentCanvasRepositoryID uint64) (*models.CanvasRepository, error) {
	// create repo
	repoInstance := &models.CanvasRepository{}
	repoInstance.CreatedByID = userID
	repoInstance.UpdatedByID = userID
	repoInstance.Name = name
	repoInstance.Icon = icon
	repoInstance.CollectionID = collectionID
	repoInstance.StudioID = studioID
	repoInstance.Position = position
	if parentCanvasRepositoryID != 0 {
		repoInstance.ParentCanvasRepositoryID = &parentCanvasRepositoryID
	}

	repoInstance.IsPublished = false
	repoInstance.DefaultBranchID = nil // Generated later
	repoInstance.Key = utils.NewNanoid()

	created := postgres.GetDB().Create(&repoInstance)
	err := created.Error
	/* todo kafka publish
	go func() {
		canvasRepo, _ := json.Marshal(repoInstance)

		r.kafka.Publish(configs.KAFKA_TOPICS_NEW_CANVAS, strconv.FormatUint(repoInstance.ID, 10), canvasRepo)

	}()
	*/

	if err != nil {
		return nil, err
	}
	fmt.Println(repoInstance)
	// create branch
	branchInstance := models.CanvasBranch{}
	branchInstance.CreatedByID = userID
	branchInstance.UpdatedByID = userID
	branchInstance.Name = models.CANVAS_BRANCH_NAME_MAIN
	branchInstance.CanvasRepositoryID = repoInstance.ID
	branchInstance.IsDefault = true
	branchInstance.PublicAccess = "private"
	branchInstance.Key = utils.NewNanoid()

	results := postgres.GetDB().Create(&branchInstance)
	err = results.Error

	if err != nil {
		return nil, err
	}
	fmt.Println(branchInstance)

	// repoInstance.ID, branchInstance.ID
	var repo *models.CanvasRepository
	err = postgres.GetDB().Model(&models.CanvasRepository{}).Where("id = ?", repoInstance.ID).Update("default_branch_id", branchInstance.ID).First(&repo).Error
	//return &repo, err

	if err != nil {
		return nil, err
	}

	err = postgres.GetDB().Model(&models.CanvasRepository{}).Where("id = ?", repoInstance.ID).Preload("DefaultBranch").First(&repo).Error
	if err != nil {
		return nil, err
	}

	err = queries.App.PermsQuery.CreateDefaultCanvasBranchPermission(repo.CollectionID, userID, studioID, repo.ID, branchInstance.ID, repo.ParentCanvasRepositoryID)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func AnonymousGetAllCanvasController(parentCollectionID uint64, parentCanvasRepositoryID uint64) (*[]models.CanvasRepository, error) {
	var canvasRepos *[]models.CanvasRepository
	var err error
	publicAccess := []string{"view", "edit", "comment"}

	if parentCollectionID != 0 {
		canvasRepos, err = canvasrepo.App.Repo.GetAnonymousCanvasRepos(parentCollectionID, publicAccess)
		fmt.Println(canvasRepos)
	} else {
		canvasRepos, err = canvasrepo.App.Repo.GetAnonymousSubCanvasRepos(parentCanvasRepositoryID, publicAccess)
	}

	if err != nil {
		return nil, err
	}

	return canvasRepos, nil
}

func AuthUserGetAllCanvasController(parentCollectionID uint64, parentCanvasRepositoryID uint64, user *models.User, studioId uint64) (*[]models.CanvasRepository, error) {
	var canvasRepos *[]models.CanvasRepository
	var accessCanvasRepo []models.CanvasRepository
	var permissionsList map[uint64]map[uint64]string

	var err error

	if parentCollectionID != 0 {
		canvasRepos, err = canvasrepo.App.Repo.GetCanvasRepos(map[string]interface{}{"collection_id": parentCollectionID, "parent_canvas_repository_id": nil, "is_archived": false})
		if err != nil {
			return nil, err
		}
		permissionsList, err = permissions.App.Service.CalculateCanvasRepoPermissions(user.ID, studioId, parentCollectionID)
		if err != nil {
			return nil, err
		}

	} else {
		canvasRepos, err = canvasrepo.App.Repo.GetCanvasRepos(map[string]interface{}{"parent_canvas_repository_id": parentCanvasRepositoryID, "is_archived": false})
		if err != nil {
			return nil, err
		}
		permissionsList, err = permissions.App.Service.CalculateSubCanvasRepoPermissions(user.ID, studioId, parentCollectionID, parentCanvasRepositoryID)
		if err != nil {
			return nil, err
		}
	}

	for _, repo := range *canvasRepos {
		repoPermissions := permissionsList[repo.ID]
		permissionValues := utils.Values(repoPermissions)
		for _, perm := range permissionValues {
			if utils.Contains(permissiongroup.UserAccessCanvasPermissionsList, perm) {
				accessCanvasRepo = append(accessCanvasRepo, repo)
				break
			}
		}
	}

	if err != nil {
		return nil, err
	}

	return canvasRepos, nil
}
