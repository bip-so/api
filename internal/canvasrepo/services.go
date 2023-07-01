package canvasrepo

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	permissiongroup "gitlab.com/phonepost/bip-be-platform/internal/permission_groups"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/pkg/integrations"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"gorm.io/datatypes"
	"net/url"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/phonepost/bip-be-platform/internal/collection"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

const DISCORD_MAX_CANAVS_COUNT = 40

type ServiceInterface interface {
	EmptyCanvasRepoInstance()
	Create(userID uint64, studioID uint64, collectionID uint64, name string, icon string) (*models.CanvasRepository, error)
	EmptyCanvasBranchInstance() *models.CanvasBranch
	CreateCanvasBranch(userID uint64, canvasRepoID uint64, name string, isTrue bool) (*models.CanvasBranch, error)
}

// Returns Empty Instance
func (s canvasRepoService) EmptyCanvasRepoInstance() *models.CanvasRepository {
	return &models.CanvasRepository{}
}

func (s canvasRepoService) UpdateCanvasRepo(name string, icon string, repoId uint64, userId uint64, cover string) bool {
	err := App.Repo.UpdateRepoNameIconOnly(name, icon, repoId, userId, cover)
	if err != nil {
		return false
	}

	return true
}

// Returns new language canvas Instance
func (s canvasRepoService) NewLanguageCanvasRepoInstance(defaultLanguageCanvasRepo *models.CanvasRepository, languageCode string, autoTranslate bool, userID uint64) *models.CanvasRepository {
	return &models.CanvasRepository{
		Name:                        defaultLanguageCanvasRepo.Name,
		Icon:                        defaultLanguageCanvasRepo.Icon,
		CollectionID:                defaultLanguageCanvasRepo.CollectionID,
		StudioID:                    defaultLanguageCanvasRepo.StudioID,
		IsPublished:                 false,
		DefaultBranchID:             nil,
		Position:                    defaultLanguageCanvasRepo.Position,
		DefaultLanguageCanvasRepoID: &defaultLanguageCanvasRepo.ID,
		Language:                    &languageCode,
		IsLanguageCanvas:            true,
		AutoTranslated:              autoTranslate,
		CreatedByID:                 userID,
		UpdatedByID:                 userID,
		Key:                         utils.NewNanoid(),
	}
}

func (s canvasRepoService) CreateStudioFileTree(studio *models.Studio, user *models.User) *discordgo.MessageEmbed {
	// loop the collections
	// initiate a string and start with the collection name with no space at starting.
	// Next create another method with recursion which will recursive start adding the canvasRepos to the string
	//var collections []models.Collection
	//err := postgres.GetDB().Model(models.Collection{}).Where("studio_id = ?", studioID).Find(&collections).Error
	//if err != nil {
	//	fmt.Println("Error in fetching collections", err)
	//}
	tree := ""
	indent := ""
	indentLevel := 1
	canvasCount := 0
	//fmt.Println(len(collections))
	collections, _ := collection.App.Controller.AuthUserCollectionsController(studio.ID, user)

	for _, collectionInstance := range *collections {
		tree += fmt.Sprintf("**%s**\n", collectionInstance.Name)
		var canvasTree string
		canvasTree, canvasCount = s.CreateCanvasRepoFileTree(collectionInstance.Id, 0, "", indent, studio, user, canvasCount, indentLevel)
		tree += canvasTree + "\n"
		fmt.Println(collectionInstance.Name, canvasCount)
		if canvasCount > DISCORD_MAX_CANAVS_COUNT {
			allStudioCanvasesCount, _ := App.Controller.AuthUserGetAllStudioCanvasControllerCount(user, studio.ID)
			diff := allStudioCanvasesCount - DISCORD_MAX_CANAVS_COUNT
			if diff <= 0 {
				diff = 1
			}
			tree += "\n" + fmt.Sprintf("[+%d more](%s)", diff, fmt.Sprintf("%s/%s", configs.GetAppInfoConfig().FrontendHost, studio.Handle))
			break
		}
	}
	fmt.Println(tree)
	embed := &discordgo.MessageEmbed{
		Title:       "Canvases List",
		Description: tree,
		Color:       0x44B244,
		Type:        "rich",
	}
	return embed
}

func (s canvasRepoService) CreateCanvasRepoFileTree(collectionID uint64, parentCanvasRepoID uint64, canvasTree string, indent string, studio *models.Studio, user *models.User, canvasCount int, indentLevel int) (string, int) {
	var subCanvas *[]CanvasRepoDefaultSerializer
	subCanvas, _ = App.Controller.AuthUserGetAllCanvasController(&GetAllCanvasValidator{
		ParentCollectionID:       collectionID,
		ParentCanvasRepositoryID: parentCanvasRepoID,
	}, user, studio.ID)
	if indent == "" {
		indent += "|—"
	} else {
		indent = "|　" + indent
	}
	charLimit := 37 - (indentLevel * 3)
	indentLevel += 1
	for _, canvas := range *subCanvas {
		if canvas.IsLanguageCanvas {
			continue
		}
		canvasCount++
		if canvasCount > DISCORD_MAX_CANAVS_COUNT {
			return canvasTree, canvasCount
		}
		canvasUrlTitle := notifications.App.Service.GenerateCanvasUrlTitle(canvas.Name, *canvas.DefaultBranchID)
		canvasName := canvas.Name
		if len(canvasName) > charLimit {
			canvasName = canvasName[:charLimit] + "..."
		}
		canvasTree += indent + fmt.Sprintf("[%s](%s)", canvasName, fmt.Sprintf("%s/%s/%s", configs.GetAppInfoConfig().FrontendHost, studio.Handle, url.QueryEscape(canvasUrlTitle)))
		//fmt.Println(tree)
		if canvas.DefaultBranch.PublicAccess == "private" {
			canvasTree += ":lock:"
		}
		canvasTree += "\n"
		var languageCanvases []models.CanvasRepository
		postgres.GetDB().Model(models.CanvasRepository{}).Where("default_language_canvas_repo_id = ? and is_language_canvas = true", canvas.ID).Find(&languageCanvases)
		for _, lcanvas := range languageCanvases {
			canvasCount++
			if canvasCount > DISCORD_MAX_CANAVS_COUNT {
				break
			}
			lcanvasName := lcanvas.Name
			if len(lcanvasName) > charLimit {
				lcanvasName = lcanvasName[:charLimit] + "..."
			}
			lcanvasUrlTitle := notifications.App.Service.GenerateCanvasUrlTitle(lcanvas.Name, *lcanvas.DefaultBranchID)
			canvasTree += indent + fmt.Sprintf("[|अ %s](%s)", lcanvasName, fmt.Sprintf("%s/%s/%s", configs.GetAppInfoConfig().FrontendHost, studio.Handle, url.QueryEscape(lcanvasUrlTitle)))
			canvasTree += "\n"
		}
		var subCanvasTree string
		subCanvasTree, canvasCount = s.CreateCanvasRepoFileTree(collectionID, canvas.ID, "", indent, studio, user, canvasCount, indentLevel)
		canvasTree += subCanvasTree
	}
	return canvasTree, canvasCount
}

func (s canvasRepoService) GetCanvasPrevAndNext(userID, canvasID uint64) (*models.CanvasRepository, *models.CanvasRepository) {
	var nextCanvas *models.CanvasRepository
	var prevCanvas *models.CanvasRepository
	canvasRepository, err := App.Repo.GetRepoWithCollection(map[string]interface{}{"id": canvasID})
	if err != nil {
		fmt.Println("Error in getting", canvasRepository.ID)
	}
	// Get the canvas repos position > canvasRepository.position order by position asc
	// if canvasRepos found take the first for the next canvas.
	// else
	// Get the collections position > canvasRepository collection position order by position asc
	// if collections found take the first collection and get the canvas of that collection by position order by position asc
	// if canvases found take the first canvas
	nextCanvas = s.GetNextSubCanvasPublicRepo(canvasRepository)
	if nextCanvas == nil {
		nextCanvases, _ := App.Repo.GetNextCanvases(canvasRepository, canvasRepository.Position)
		if len(nextCanvases) == 0 && canvasRepository.ParentCanvasRepositoryID != nil {
			nextCanvases, _ = App.Repo.GetNextCanvases(canvasRepository.ParentCanvasRepository, canvasRepository.ParentCanvasRepository.Position)
		}
		if len(nextCanvases) > 0 {
			for _, repo := range nextCanvases {
				fmt.Println(repo.ID, repo.Name, repo.DefaultBranch.PublicAccess)
				//if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(userID, *repo.DefaultBranchID, permissiongroup.CANVAS_BRANCH_VIEW); err == nil || hasPermission {
				//	nextCanvas = &repo
				//	break
				//}
				if repo.DefaultBranch.PublicAccess != models.PRIVATE {
					nextCanvas = &repo
					break
				}
				if repo.HasPublicCanvas {
					nextCanvas = s.GetNextSubCanvasPublicRepo(&repo)
				}
			}
		} else {
			nextCollections, _ := App.Repo.GetNextCollections(canvasRepository)
			if len(nextCollections) > 0 {
				nextCollection := nextCollections[0]
				nextCanvases, _ = App.Repo.GetCanvasesOrderByPositionAsc(map[string]interface{}{"collection_id": nextCollection.ID, "parent_canvas_repository_id": nil, "is_archived": false, "is_published": true, "is_language_canvas": false})
				for _, repo := range nextCanvases {
					//if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(userID, *repo.DefaultBranchID, permissiongroup.CANVAS_BRANCH_VIEW); err == nil || hasPermission {
					//	nextCanvas = &repo
					//	break
					//}
					if repo.DefaultBranch.PublicAccess != models.PRIVATE {
						nextCanvas = &repo
						break
					}
				}
			}
		}
	}

	// Get the canvas repos position < canvasRepository.position order by position desc
	// if canvasRepos found take the first for the next canvas.
	// else
	// Get the collections position < canvasRepository collection position order by position desc
	// if collections found take the first collection and get the canvas of that collection by position order by position desc
	// if canvases found take the first canvas
	sameLevel := true
	prevCanvases, _ := App.Repo.GetPrevCanvases(canvasRepository, canvasRepository.Position)
	if len(prevCanvases) == 0 && canvasRepository.ParentCanvasRepositoryID != nil {
		sameLevel = false
		prevCanvases, _ = App.Repo.GetPrevCanvases(canvasRepository.ParentCanvasRepository, canvasRepository.ParentCanvasRepository.Position+1)
	}
	if len(prevCanvases) > 0 {
		prevCanvas = s.GetPrevSubCanvasPublicRepo(prevCanvases, sameLevel)
	} else {
		prevCollections, _ := App.Repo.GetPrevCollections(canvasRepository)
		if len(prevCollections) > 0 {
			prevCollection := prevCollections[0]
			prevCanvases, _ = App.Repo.GetCanvasesOrderByPositionDesc(map[string]interface{}{"collection_id": prevCollection.ID, "parent_canvas_repository_id": nil, "is_archived": false, "is_published": true, "is_language_canvas": false})
			for _, repo := range prevCanvases {
				//if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(userID, *repo.DefaultBranchID, permissiongroup.CANVAS_BRANCH_VIEW); err == nil || hasPermission {
				//	prevCanvas = &repo
				//	break
				//}
				if repo.DefaultBranch.PublicAccess != models.PRIVATE || repo.HasPublicCanvas {
					prevCanvas = &repo
					break
				}
			}
			if prevCanvas != nil {
				prevCanvas = s.GetCollectionLastCanvas(prevCanvas)
			}
		}
	}
	return nextCanvas, prevCanvas
}

func (s canvasRepoService) GetNextSubCanvasPublicRepo(repo *models.CanvasRepository) *models.CanvasRepository {
	nextCanvases, _ := App.Repo.GetNextSubCanvases(repo, 0)
	var nextSubCanvasRepo *models.CanvasRepository
	for _, canvasRepo := range nextCanvases {
		if canvasRepo.DefaultBranch.PublicAccess != "private" {
			nextSubCanvasRepo = &canvasRepo
			break
		}
		if canvasRepo.HasPublicCanvas {
			nextSubCanvasRepo = s.GetNextSubCanvasPublicRepo(&canvasRepo)
			if nextSubCanvasRepo != nil {
				return nextSubCanvasRepo
			}
		}
	}
	return nextSubCanvasRepo
}

func (s canvasRepoService) GetLangNextSubCanvasPublicRepo(repo *models.CanvasRepository, language string) *models.CanvasRepository {
	nextCanvases, _ := App.Repo.GetNextSubCanvases(repo, 0)
	var nextSubCanvasRepo *models.CanvasRepository
	for _, canvasRepo := range nextCanvases {
		if canvasRepo.DefaultBranch.PublicAccess != "private" {
			nextLangCanvases, _ := App.Repo.GetLangCanvases(canvasRepo.ID, language)
			if len(nextLangCanvases) > 0 {
				nextSubCanvasRepo = &canvasRepo
				break
			}
		}
		if canvasRepo.HasPublicCanvas {
			nextSubCanvasRepo = s.GetLangNextSubCanvasPublicRepo(&canvasRepo, language)
			if nextSubCanvasRepo != nil {
				return nextSubCanvasRepo
			}
		}
	}
	return nextSubCanvasRepo
}

func (s canvasRepoService) GetPrevSubCanvasPublicRepo(prevCanvases []models.CanvasRepository, sameLevel bool) *models.CanvasRepository {
	var prevCanvas *models.CanvasRepository
	for _, pRepo := range prevCanvases {
		fmt.Println(pRepo.Name)
		if pRepo.DefaultBranch.PublicAccess != models.PRIVATE {
			prevCanvas = &pRepo
			if sameLevel {
				prevCanvas = s.GetCollectionLastCanvas(prevCanvas)
				break
			}
			break
		}
	}
	if prevCanvas != nil {
		return prevCanvas
	}
	prevNonPermCanvas := prevCanvases[0]
	if prevNonPermCanvas.ParentCanvasRepositoryID != nil {
		prevCanvases, _ = App.Repo.GetPrevCanvases(prevNonPermCanvas.ParentCanvasRepository, prevNonPermCanvas.ParentCanvasRepository.Position)
		if len(prevCanvases) == 0 {
			return prevCanvas
		}
		if prevCanvas == nil {
			prevCanvas = s.GetPrevSubCanvasPublicRepo(prevCanvases, false)
		}
	}
	return prevCanvas
}

func (s canvasRepoService) GetLangPrevSubCanvasPublicRepo(prevCanvases []models.CanvasRepository, sameLevel bool, language string) *models.CanvasRepository {
	var prevCanvas *models.CanvasRepository
	var prevLangCanvas *models.CanvasRepository
	for _, pRepo := range prevCanvases {
		fmt.Println(pRepo.Name)
		if pRepo.DefaultBranch.PublicAccess != models.PRIVATE {
			prevCanvas = &pRepo
			if sameLevel {
				prevCanvas = s.GetCollectionLastCanvas(prevCanvas)
				if prevCanvas != nil {
					prevLangCanvases, _ := App.Repo.GetLangCanvases(prevCanvas.ID, language)
					if len(prevLangCanvases) > 0 {
						prevLangCanvas = prevCanvas
						break
					}
				}
			}
			if prevCanvas != nil {
				prevLangCanvases, _ := App.Repo.GetLangCanvases(prevCanvas.ID, language)
				if len(prevLangCanvases) > 0 {
					prevLangCanvas = prevCanvas
					break
				}
			}
		}
	}
	if prevCanvas != nil && prevLangCanvas != nil {
		return prevCanvas
	}
	prevNonPermCanvas := prevCanvases[0]
	if prevNonPermCanvas.ParentCanvasRepositoryID != nil {
		prevCanvases, _ = App.Repo.GetPrevCanvases(prevNonPermCanvas.ParentCanvasRepository, prevNonPermCanvas.ParentCanvasRepository.Position)
		if len(prevCanvases) == 0 {
			return prevCanvas
		}
		if prevCanvas == nil {
			prevCanvas = s.GetLangPrevSubCanvasPublicRepo(prevCanvases, false, language)
		}
	}
	return prevCanvas
}

func (s canvasRepoService) GetCollectionLastCanvas(lastCanvas *models.CanvasRepository) *models.CanvasRepository {
	nextCanvases, _ := App.Repo.GetNextSubCanvases(lastCanvas, 0)
	if len(nextCanvases) == 0 {
		return lastCanvas
	}
	var newLastCanvas *models.CanvasRepository
	for i, _ := range nextCanvases {
		canvas := nextCanvases[len(nextCanvases)-i-1]
		if canvas.DefaultBranch.PublicAccess != models.PRIVATE {
			newLastCanvas = &canvas
			break
		}
		if canvas.HasPublicCanvas {
			subLastCanvas := s.GetCollectionLastCanvas(&canvas)
			if subLastCanvas != nil {
				return subLastCanvas
			}
		}
	}
	if newLastCanvas == nil {
		return lastCanvas
	}
	lastCanvas = s.GetCollectionLastCanvas(newLastCanvas)
	return lastCanvas
}

/*
	SendCollectionTreeToDiscord method is used to send or update tree message in discord.
	Args:
		collectionID uint64

	Algorithm:
		- If collection is private we delete the collection from discord messages list.
		- If collection is updated to public or moved we then updates the tree structure in discord.
		- If only canvas is added or updated we only updates that particular collection message.
*/
func (s canvasRepoService) SendCollectionTreeToDiscord(collectionID uint64) {
	collectionInstance, err := App.Repo.GetCollectionPreloadStudio(map[string]interface{}{"id": collectionID})
	if err != nil {
		fmt.Println("collection Instance not found by ID", err)
		return
	}
	studioIntegration, err := App.Repo.GetStudioIntegration(collectionInstance.StudioID, models.DISCORD_INTEGRATION_TYPE)
	if err != nil {
		fmt.Println("studio integration not found", err)
		return
	}
	discordMessagesData := DiscordMessagesData{}
	json.Unmarshal(*studioIntegration.MessagesData, &discordMessagesData)
	messagesData := discordMessagesData.CollectionsMap
	collectionMessage := map[string]interface{}{}
	collectionIDStr := utils.String(collectionID)
	bipCanvasesChannelID := discordMessagesData.CanvasesChannelID
	discordMessageIDs := discordMessagesData.MessageIDs
	if (collectionInstance.PublicAccess == models.PRIVATE && !collectionInstance.HasPublicCanvas) || collectionInstance.IsArchived {
		fmt.Println("collection has public Access private", collectionInstance.PublicAccess, collectionInstance.HasPublicCanvas)
		if messagesData != nil {
			// Delete the message from discord if present
			// Delete the messageId
			collectionMessage = messagesData[collectionIDStr]
			if collectionMessage != nil {
				s.DeleteDiscordCollectionFromTree(studioIntegration, bipCanvasesChannelID, discordMessagesData, collectionMessage, collectionIDStr)
			}
		}
		return
	}
	if messagesData != nil {
		collectionMessage = messagesData[collectionIDStr]
		if collectionMessage == nil {
			// here write the logic
			// if collection is newly added, and it doesn't have any public canvases then we are returning.
			var publicCanvasCount int64
			postgres.GetDB().Table("canvas_repositories").
				Joins("left join canvas_branches on canvas_repositories.id = canvas_branches.canvas_repository_id").
				Where("canvas_repositories.is_archived = false and canvas_repositories.is_published = true and canvas_repositories.collection_id = ? and canvas_branches.public_access <> 'private'", collectionID).
				Count(&publicCanvasCount)
			if publicCanvasCount == 0 {
				return
			} else {
				s.RearrangeDiscordCanvasOnNewCollection(collectionInstance, bipCanvasesChannelID, studioIntegration, discordMessageIDs, discordMessagesData)
			}
		}
		valuePositionFloat, isOk := collectionMessage["position"].(float64)
		if isOk {
			collectionPositionInt := uint(valuePositionFloat)
			if collectionPositionInt != collectionInstance.Position {
				s.RearrangeDiscordCanvasTree(discordMessagesData, collectionInstance, bipCanvasesChannelID, studioIntegration, discordMessageIDs)
				return
			} else {
				var publicCanvasCount int64
				postgres.GetDB().Table("canvas_repositories").
					Joins("left join canvas_branches on canvas_repositories.id = canvas_branches.canvas_repository_id").
					Where("canvas_repositories.is_archived = false and canvas_repositories.is_published = true and canvas_repositories.collection_id = ? and canvas_branches.public_access <> 'private'", collectionID).
					Count(&publicCanvasCount)
				fmt.Println("publicCanvasCount", publicCanvasCount)
				if publicCanvasCount == 0 {
					s.DeleteDiscordCollectionFromTree(studioIntegration, bipCanvasesChannelID, discordMessagesData, collectionMessage, collectionIDStr)
					fmt.Println("came here after deleting the collection")
					return
				}
				fmt.Println("collection message Index", collectionMessage["index"])
				embed, _ := s.BuildDiscordCollectionTree(collectionID)
				if collectionMessage["index"] != nil && collectionMessage["index"].(float64) == 0 {
					embed.Title = "Public Canvases List"
				}
				// Edit the message
				_, err = integrations.EditDiscordEmbedComplexToChannel(bipCanvasesChannelID, collectionMessage["discordMessageId"].(string), embed)
				if err != nil {
					fmt.Println("Error in sending discord messages: ", err)
					return
				}
			}
			return
		}
	} else {
		// have to send the new message
		// and get the message id and index is 0
		// save to the db
		var publicCanvasCount int64
		postgres.GetDB().Table("canvas_repositories").
			Joins("left join canvas_branches on canvas_repositories.id = canvas_branches.canvas_repository_id").
			Where("canvas_repositories.is_archived = false and canvas_repositories.is_published = true and canvas_repositories.collection_id = ? and canvas_branches.public_access <> 'private'", collectionID).
			Count(&publicCanvasCount)
		if publicCanvasCount == 0 {
			return
		}
		embed, _ := s.BuildDiscordCollectionTree(collectionID)
		embed.Title = "Public Canvases List"
		message, err := integrations.SendDiscordEmbedToChannel(bipCanvasesChannelID, embed)
		if err != nil {
			fmt.Println("Error in sending discord messages: ", err)
			return
		}
		messagesData[collectionIDStr] = map[string]interface{}{"discordMessageId": message.ID, "index": 0, "position": collectionInstance.Position}
		discordMessageIDs = append(discordMessageIDs, message.ID)
		discordMessagesData.MessageIDs = discordMessageIDs
		discordMessagesData.CollectionsMap = messagesData
		discordMessagesDataStr, _ := json.Marshal(discordMessagesData)
		App.Repo.StudioIntegrationUpdate(studioIntegration.ID, map[string]interface{}{"messages_data": datatypes.JSON(discordMessagesDataStr)})
	}
}

func (s canvasRepoService) DeleteDiscordCollectionFromTree(studioIntegration *models.StudioIntegration, bipCanvasesChannelID string, discordMessagesData DiscordMessagesData, collectionMessage map[string]interface{}, collectionIDStr string) {
	err := integrations.DeleteDiscordEmbedToChannel(bipCanvasesChannelID, collectionMessage["discordMessageId"].(string))
	if err != nil {
		fmt.Println("Error in deleting discord embed channel", err)
		return
	}
	messagesData := discordMessagesData.CollectionsMap
	delete(messagesData, collectionIDStr)
	collectionMessageIndexFloat := collectionMessage["index"].(float64)
	collectionMessageIndexInt := int(collectionMessageIndexFloat)
	for key, val := range messagesData {
		if val["index"] != nil {
			indexFloat := val["index"].(float64)
			fmt.Println("index float", indexFloat)
			indexInt := int(indexFloat)
			if indexInt > collectionMessageIndexInt {
				newValue := val
				newValue["index"] = indexInt - 1
				messagesData[key] = newValue
			}
		}
	}
	discordMessageIDs := discordMessagesData.MessageIDs
	fmt.Println(collectionMessage["discordMessageId"].(string))
	fmt.Println("before", discordMessageIDs)
	updatedDiscordMessageIDs := utils.Remove(discordMessageIDs, collectionMessage["discordMessageId"].(string))
	fmt.Println("after", updatedDiscordMessageIDs)
	discordMessagesData.MessageIDs = updatedDiscordMessageIDs
	discordMessagesData.CollectionsMap = messagesData
	discordMessagesDataStr, _ := json.Marshal(discordMessagesData)
	App.Repo.StudioIntegrationUpdate(studioIntegration.ID, map[string]interface{}{"messages_data": datatypes.JSON(discordMessagesDataStr)})
}

/*
	BuildDiscordCollectionTree method returns the discord embed message
	Args:
		collectionId uint64
	Based on collection we get all the repos and construct a tree structure.
	We are considering only the public canvases here.
	If canvases are more than 40 then we just +x canvases more at the last.
*/
func (s canvasRepoService) BuildDiscordCollectionTree(collectionID uint64) (*discordgo.MessageEmbed, error) {
	collectionInstance, err := App.Repo.GetCollectionPreloadStudio(map[string]interface{}{"id": collectionID})
	if err != nil {
		fmt.Println("collection Instance not found by ID", err)
		return nil, err
	}
	title := "　"
	tree := ""
	indent := ""
	indentLevel := 1
	tree += fmt.Sprintf("__**%s**__\n", collectionInstance.Name)
	canvasTree, canvasesCount := s.CreatePublicCanvasRepoFileTree(collectionInstance.ID, 0, "", indent, 0, indentLevel)
	var canvasRepos []models.CanvasRepository
	var collectionRootCount int
	postgres.GetDB().Model(models.CanvasRepository{}).Where("collection_id = ? and is_archived = false and is_published = true", collectionInstance.ID).Preload("DefaultBranch").Find(&canvasRepos)
	for _, repo := range canvasRepos {
		if repo.IsLanguageCanvas || (repo.DefaultBranch.PublicAccess == models.PRIVATE && !repo.HasPublicCanvas) {
			continue
		}
		var languageCanvases []models.CanvasRepository
		postgres.GetDB().Model(models.CanvasRepository{}).Where("default_language_canvas_repo_id = ? and is_language_canvas = true and is_archived = false", repo.ID).Preload("DefaultBranch").Find(&languageCanvases)
		for _, lcanvas := range languageCanvases {
			if lcanvas.DefaultBranch.PublicAccess == models.PRIVATE || !lcanvas.IsPublished {
				continue
			}
			collectionRootCount++
		}
		collectionRootCount++
	}
	fmt.Println("canvas count", collectionRootCount, canvasesCount)
	diffCount := collectionRootCount - canvasesCount
	if diffCount > 0 {
		canvasTree += "\n" + fmt.Sprintf("[+%d more](%s)", diffCount, fmt.Sprintf("%s/%s", configs.GetAppInfoConfig().FrontendHost, collectionInstance.Studio.Handle))
	}
	tree += canvasTree + "\n"
	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: tree,
		Color:       0x44B244,
		Type:        "rich",
	}
	return embed, nil
}

/*
	RearrangeDiscordCanvasOnNewCollection method re-arranges the messages in bip-canvases discord channel.
	When a new collection is added or any old collection is made public.

	Args:
		collectionInstance *models.Collection
		bipCanvasesChannelID string
		studioIntegration *models.StudioIntegration
		discordMessageIDs []string, discordMessagesData DiscordMessagesData
*/
func (s canvasRepoService) RearrangeDiscordCanvasOnNewCollection(collectionInstance *models.Collection, bipCanvasesChannelID string, studioIntegration *models.StudioIntegration, discordMessageIDs []string, discordMessagesData DiscordMessagesData) {
	newMessageData := map[string]map[string]interface{}{}
	var collections []models.Collection
	err := postgres.GetDB().Model(models.Collection{}).Where("studio_id = ? and is_archived = false", collectionInstance.StudioID).Preload("Studio").Order("position ASC").Find(&collections).Error
	if err != nil {
		fmt.Println("Error in fetching collections", err)
	}
	var publicCollections []models.Collection
	for _, col := range collections {
		if col.PublicAccess == models.PRIVATE && !col.HasPublicCanvas {
			continue
		}
		var publicCanvasCount int64
		postgres.GetDB().Table("canvas_repositories").
			Joins("left join canvas_branches on canvas_repositories.id = canvas_branches.canvas_repository_id").
			Where("canvas_repositories.is_archived = false and canvas_repositories.is_published = true and canvas_repositories.collection_id = ? and canvas_branches.public_access <> 'private'", col.ID).
			Count(&publicCanvasCount)
		if publicCanvasCount == 0 {
			continue
		}
		publicCollections = append(publicCollections, col)
	}
	fmt.Println("length of collections", len(collections))
	for i, col := range publicCollections {
		messageID := discordMessageIDs[i]
		embed, _ := s.BuildDiscordCollectionTree(col.ID)
		if i == 0 {
			embed.Title = "Public Canvases List"
		}
		fmt.Println("message id printing here", messageID)
		message, err := integrations.EditDiscordEmbedComplexToChannel(bipCanvasesChannelID, messageID, embed)
		if err != nil {
			fmt.Println("Error in sending discord messages: ", err)
			continue
		}
		data := map[string]interface{}{"discordMessageId": message.ID, "index": i, "position": col.Position}
		newMessageData[utils.String(col.ID)] = data
	}
	discordComponents := []interface{}{
		notifications.ActionRowsComponent{
			Type: 1,
			Components: []interface{}{
				notifications.MessageBtnComponent{
					Type:     2,
					Label:    "Check your access",
					Style:    2,
					CustomID: "checkCanvasAccess",
				},
			},
		},
	}
	msg, err := integrations.SendDiscordDMMessageToChannel(bipCanvasesChannelID, []string{"Above is a list of all public canvases. Use the button below to check out list of private canvases you have access to!\n\nAlternatively, you can use `/bip-search` here to find a specific canvas or create a new one using `/bip-new`!"}, discordComponents)
	fmt.Println("msg", msg.GuildID, err)
	fmt.Println("message Data", newMessageData)
	data := map[string]interface{}{"discordMessageId": msg.ID, "index": nil, "position": nil}
	newMessageData["menu"] = data
	discordMessageIDs = append(discordMessageIDs, msg.ID)
	discordMessagesData.MessageIDs = discordMessageIDs
	discordMessagesData.CollectionsMap = newMessageData
	// save to db
	discordMessagesDataStr, _ := json.Marshal(discordMessagesData)
	App.Repo.StudioIntegrationUpdate(studioIntegration.ID, map[string]interface{}{"messages_data": datatypes.JSON(discordMessagesDataStr)})
}

/*
	RearrangeDiscordCanvasTree method re-arranges when a collection is moved from its position.
	Args:
		discordMessagesData
		collectionInstance *models.Collection
		bipCanvasesChannelID string
		studioIntegration *models.StudioIntegration
		discordMessageIDs []string
*/
func (s canvasRepoService) RearrangeDiscordCanvasTree(discordMessagesData DiscordMessagesData, collectionInstance *models.Collection, bipCanvasesChannelID string, studioIntegration *models.StudioIntegration, discordMessageIDs []string) {
	newMessageData := map[string]map[string]interface{}{}
	var collections []models.Collection
	err := postgres.GetDB().Model(models.Collection{}).Where("studio_id = ? and is_archived = false", collectionInstance.StudioID).Preload("Studio").Order("position ASC").Find(&collections).Error
	if err != nil {
		fmt.Println("Error in fetching collections", err)
	}
	var publicCollections []models.Collection
	for _, col := range collections {
		if col.PublicAccess == models.PRIVATE && !col.HasPublicCanvas {
			continue
		}
		var publicCanvasCount int64
		postgres.GetDB().Table("canvas_repositories").
			Joins("left join canvas_branches on canvas_repositories.id = canvas_branches.canvas_repository_id").
			Where("canvas_repositories.is_archived = false and canvas_repositories.is_published = true and canvas_repositories.collection_id = ? and canvas_branches.public_access <> 'private'", col.ID).
			Count(&publicCanvasCount)
		if publicCanvasCount == 0 {
			continue
		}
		publicCollections = append(publicCollections, col)
	}
	for i, col := range publicCollections {
		embed, _ := s.BuildDiscordCollectionTree(col.ID)
		if i == 0 {
			embed.Title = "Public Canvases List"
		}
		messageID := discordMessageIDs[i]
		message, err := integrations.EditDiscordEmbedComplexToChannel(bipCanvasesChannelID, messageID, embed)
		if err != nil {
			fmt.Println("Error in sending discord messages: ", err)
			continue
		}
		data := map[string]interface{}{"discordMessageId": message.ID, "index": i, "position": col.Position}
		newMessageData[utils.String(col.ID)] = data
	}
	// save to db
	messagesData := discordMessagesData.CollectionsMap
	newMessageData["menu"] = messagesData["menu"]
	discordMessagesData.CollectionsMap = newMessageData
	discordMessagesDataStr, _ := json.Marshal(discordMessagesData)
	App.Repo.StudioIntegrationUpdate(studioIntegration.ID, map[string]interface{}{"messages_data": datatypes.JSON(discordMessagesDataStr)})
}

/*
	CreatePublicCanvasRepoFileTree method is triggered recursively to create the tree structure.
	Args:
		collectionID uint64
		parentCanvasRepoID uint64
		canvasTree string
		indent int
		canvasesCount int
*/
func (s canvasRepoService) CreatePublicCanvasRepoFileTree(collectionID uint64, parentCanvasRepoID uint64, canvasTree string, indent string, canvasesCount int, indentLevel int) (string, int) {
	var subCanvas []models.CanvasRepository
	if parentCanvasRepoID == 0 {
		postgres.GetDB().Model(models.CanvasRepository{}).Where("collection_id = ? and parent_canvas_repository_id is null and is_archived = false and is_published = true", collectionID).Preload("Studio").Preload("DefaultBranch").Order("position ASC").Find(&subCanvas)
	} else {
		postgres.GetDB().Model(models.CanvasRepository{}).Where("collection_id = ? and parent_canvas_repository_id = ? and is_archived = false and is_published = true", collectionID, parentCanvasRepoID).Preload("Studio").Preload("DefaultBranch").Order("position ASC").Find(&subCanvas)
	}
	if indent == "" {
		indent += "|—"
	} else {
		indent = "|　" + indent
	}
	charLimit := 38 - (indentLevel * 3)
	indentLevel += 1
	for _, canvas := range subCanvas {
		if canvas.IsLanguageCanvas || (canvas.DefaultBranch.PublicAccess == models.PRIVATE && !canvas.HasPublicCanvas) {
			continue
		}
		canvasesCount++
		if canvasesCount > DISCORD_MAX_CANAVS_COUNT {
			return canvasTree, canvasesCount
		}
		//fmt.Println(tree)
		canvasName := canvas.Name
		if len(canvasName) > charLimit {
			canvasName = canvasName[:charLimit] + "..."
		}
		if canvas.DefaultBranch.PublicAccess == models.PRIVATE {
			canvasTree += indent + fmt.Sprintf("%s", canvasName)
		} else {
			canvasUrlTitle := notifications.App.Service.GenerateCanvasUrlTitle(canvas.Name, *canvas.DefaultBranchID)
			canvasTree += indent + fmt.Sprintf("[%s](%s)", canvasName, fmt.Sprintf("%s/%s/%s", configs.GetAppInfoConfig().FrontendHost, canvas.Studio.Handle, url.QueryEscape(canvasUrlTitle)))
		}
		canvasTree += "\n"
		var languageCanvases []models.CanvasRepository
		postgres.GetDB().Model(models.CanvasRepository{}).Where("default_language_canvas_repo_id = ? and is_language_canvas = true and is_archived = false", canvas.ID).Preload("DefaultBranch").Order("position ASC").Find(&languageCanvases)
		for _, lcanvas := range languageCanvases {
			if lcanvas.DefaultBranch.PublicAccess == models.PRIVATE || !lcanvas.IsPublished {
				continue
			}
			canvasesCount++
			if canvasesCount > DISCORD_MAX_CANAVS_COUNT {
				break
			}
			lcanvasUrlTitle := notifications.App.Service.GenerateCanvasUrlTitle(lcanvas.Name, *lcanvas.DefaultBranchID)
			lcanvasName := lcanvas.Name
			if len(lcanvasName) > charLimit {
				lcanvasName = lcanvasName[:charLimit] + "..."
			}
			canvasTree += indent + fmt.Sprintf("[|अ %s](%s)", lcanvasName, fmt.Sprintf("%s/%s/%s", configs.GetAppInfoConfig().FrontendHost, canvas.Studio.Handle, url.QueryEscape(lcanvasUrlTitle)))
			canvasTree += "\n"
		}
		var subCanvasTree string
		subCanvasTree, canvasesCount = s.CreatePublicCanvasRepoFileTree(collectionID, canvas.ID, "", indent, canvasesCount, indentLevel)
		canvasTree += subCanvasTree
	}
	return canvasTree, canvasesCount
}

func (s canvasRepoService) GetCanvasRepoWithKey(key string, user *models.User) (*CanvasRepoDefaultSerializer, map[string]interface{}) {
	canvasRepoInstance, err := App.Repo.GetCanvasRepoByKey(key)
	if canvasRepoInstance == nil {
		return nil, map[string]interface{}{
			"error": "No key found ",
		}
	}

	if user == nil && canvasRepoInstance.DefaultBranch.PublicAccess == models.CANVAS_BRANCH_PUBLIC_ACCESS_PRIVATE {
		return nil, map[string]interface{}{
			"error": "Anonymous user does not have permissions to view this.",
		}
	}

	var userID uint64
	if user != nil {
		userID = user.ID
	}
	if hasPermission, err := permissions.App.Service.CanUserDoThisOnBranch(userID, *canvasRepoInstance.DefaultBranchID, permissiongroup.CANVAS_BRANCH_VIEW_METADATA); err != nil || !hasPermission {
		// We are adding access_requested is User already have Access Request
		exists := App.Repo.AcceesRequestExistsSimple(*canvasRepoInstance.DefaultBranchID, user.ID)
		return nil, map[string]interface{}{
			"error":            "User does not have permissions to view branch",
			"access_requested": exists,
		}
	}

	var canvasRepoView *CanvasRepoDefaultSerializer
	// Checking if the request is made by loggedIn user or not
	// If user == nil we trigger the Anonymous flow to get the canvas or vice-versa
	if user == nil {
		canvasRepoView, err = App.Controller.AnonymousGetOneCanvasByKeyController(key)
	} else {
		canvasRepoView, err = App.Controller.AuthGetOneCanvasByKeyController(key, user.ID)
	}
	if err != nil {
		return nil, map[string]interface{}{
			"error": err.Error(),
		}
	}
	return canvasRepoView, nil
}

func (s canvasRepoService) GetFirstPublicCanvasOfStudio(studioID uint64) (*models.CanvasRepository, error) {
	var canvasRepository *models.CanvasRepository
	collections, _ := App.Repo.AnonymousGetCollections(studioID)
	if collections == nil || len(*collections) == 0 {
		return nil, errors.New("public collection not found")
	}
	for _, col := range *collections {
		canvasRepository, _ = s.FirstAnonymousCanvas(col.ID, studioID, 0)
		if canvasRepository != nil {
			return canvasRepository, nil
		}
	}
	return canvasRepository, nil
}

func (s canvasRepoService) GetFirstUserCanvasOfStudio(studioID uint64, authUser *models.User) (*models.CanvasRepository, error) {
	var canvasRepository *models.CanvasRepository
	collections, _ := collection.App.Controller.AuthUserCollectionsController(studioID, authUser)
	if collections == nil || len(*collections) == 0 {
		return nil, errors.New("public collection not found")
	}
	for _, col := range *collections {
		canvasRepository, _ = s.FirstUserAccessCanvasFromCollectionID(studioID, col.Id, 0, authUser)
		if canvasRepository != nil {
			return canvasRepository, nil
		}
	}
	return canvasRepository, nil
}

func (s canvasRepoService) FirstAnonymousCanvas(collectionID, studioID, parentCanvasID uint64) (*models.CanvasRepository, error) {
	var canvasRepository *models.CanvasRepository
	var canvasRepos *[]models.CanvasRepository

	if parentCanvasID == 0 {
		canvasRepos, _ = App.Repo.GetAnonymousCanvasRepos(collectionID, []string{models.VIEW, models.COMMENT, models.EDIT})
	} else {
		canvasRepos, _ = App.Repo.GetAnonymousSubCanvasRepos(parentCanvasID, []string{models.VIEW, models.COMMENT, models.EDIT})
	}
	for _, repo := range *canvasRepos {
		if repo.IsLanguageCanvas {
			continue
		}
		if repo.DefaultBranch.PublicAccess != models.PRIVATE {
			return &repo, nil
		}
		canvasRepository, _ = s.FirstAnonymousCanvas(collectionID, studioID, repo.ID)
		if canvasRepository != nil {
			return canvasRepository, nil
		}
	}
	return canvasRepository, nil
}

func (s canvasRepoService) FirstUserAccessCanvasFromCollectionID(studioID, collectionID, parentCanvasRepoID uint64, authUser *models.User) (*models.CanvasRepository, error) {
	var canvasRepository *models.CanvasRepository
	var permissionsList map[uint64]map[uint64]string
	var canvasRepos *[]models.CanvasRepository
	var err error

	if parentCanvasRepoID == 0 {
		canvasRepos, err = App.Repo.GetCanvasRepos(map[string]interface{}{"collection_id": collectionID, "parent_canvas_repository_id": nil, "is_archived": false, "is_processing": false, "is_language_canvas": false})
		if err != nil {
			return nil, err
		}
		permissionsList, err = permissions.App.Service.CalculateCanvasRepoPermissions(authUser.ID, studioID, collectionID)
		if err != nil {
			return nil, err
		}

	} else {
		canvasRepos, err = App.Repo.GetCanvasRepos(map[string]interface{}{"parent_canvas_repository_id": parentCanvasRepoID, "is_archived": false, "is_processing": false, "is_language_canvas": false})
		if err != nil {
			return nil, err
		}
		permissionsList, err = permissions.App.Service.CalculateSubCanvasRepoPermissions(authUser.ID, studioID, collectionID, parentCanvasRepoID)
		if err != nil {
			return nil, err
		}
	}

	for _, repo := range *canvasRepos {
		if utils.Contains([]string{models.VIEW, models.EDIT, models.COMMENT}, repo.DefaultBranch.PublicAccess) {
			canvasRepository = &repo
			return canvasRepository, nil
		} else {
			repoPermissions := permissionsList[repo.ID]
			permissionValues := utils.Values(repoPermissions)
			for _, perm := range permissionValues {
				if utils.Contains(permissiongroup.UserAccessViewCanvasPermissionsList, perm) {
					canvasRepository = &repo
					fmt.Println(canvasRepository.Name, canvasRepository.Key, canvasRepository.ID)
					return canvasRepository, nil
				}
			}
		}
		canvasRepository, _ = s.FirstUserAccessCanvasFromCollectionID(studioID, collectionID, repo.ID, authUser)
		if canvasRepository != nil {
			return canvasRepository, nil
		}
	}
	return canvasRepository, nil
}

func (s canvasRepoService) ProcessRoleCollectionsPermissions(collections []models.Collection, roleID uint64) (*[]collection.CollectionSerializer, error) {
	permissionList := permissions.App.Service.CalculateCollectionRolePermissions(roleID)
	accessCollections := &[]collection.CollectionSerializer{}
	for _, col := range collections {
		collectionPermission := permissionList[col.ID]
		if collectionPermission != "" && collectionPermission != permissiongroup.PGCollectionNone().SystemName {
			collectionData := collection.CollectionSerializerData(&col)
			collectionData.Permission = collectionPermission
			*accessCollections = append(*accessCollections, collectionData)
		} else {
			collectionData := collection.CollectionSerializerData(&col)
			collectionData.Permission = permissiongroup.PGCollectionNone().SystemName
			*accessCollections = append(*accessCollections, collectionData)
		}
	}
	return accessCollections, nil
}

func (s canvasRepoService) ProcessUserCollectionsPermissions(collections []models.Collection, studioId uint64, userID uint64) (*[]collection.CollectionSerializer, error) {
	var allTheCollectionIDs []uint64
	var ActualPermsObject []collection.CollectionActualPermissionsObject
	memberObject, _ := App.Repo.GetMemberByUserID(userID, studioId)
	// Making list of all ht ecollections
	for _, collectionInstance := range collections {
		allTheCollectionIDs = append(allTheCollectionIDs, collectionInstance.ID)
	}
	permissionList, err := permissions.App.Service.CalculateCollectionPermissions(userID, studioId)
	if err != nil {
		return nil, err
	}
	// We are trying to get actual perms for logged in user. via vi the collection
	// Building the new array
	for _, collectionID := range allTheCollectionIDs {
		if permissionList[collectionID] == "" {
			ActualPermsObject = append(ActualPermsObject, collection.CollectionActualPermissionsObject{})
		} else {
			// We have calculated value of the permissions
			actualPerms := collection.CollectionPermissionActual(collectionID, memberObject.ID, studioId)
			for _, ap := range actualPerms {
				ActualPermsObject = append(ActualPermsObject, ap)
			}
		}
	}
	fmt.Printf("%+v\n", ActualPermsObject)
	accessCollections := &[]collection.CollectionSerializer{}
	for _, col := range collections {
		collectionPermission := permissionList[col.ID]
		if collectionPermission != "" && collectionPermission != permissiongroup.PGCollectionNone().SystemName {
			collectionData := collection.CollectionSerializerData(&col)
			collectionData.Permission = collectionPermission
			collectionData.ActualPermsObject = collection.PluckTheObject(ActualPermsObject, col.ID)
			collectionData.MemberPermsObject = collection.MemberCollectionPermissionActualCalculator(col.ID, memberObject, studioId)
			collectionData.RolePermsObject = collection.RoleCollectionPermissionActualCalculator(col.ID, memberObject, studioId)
			*accessCollections = append(*accessCollections, collectionData)
		} else {
			collectionData := collection.CollectionSerializerData(&col)
			collectionData.Permission = permissiongroup.PGCollectionNone().SystemName
			collectionData.ActualPermsObject = collection.PluckTheObject(ActualPermsObject, col.ID)
			*accessCollections = append(*accessCollections, collectionData)
		}
	}
	return accessCollections, nil
}

func (s canvasRepoService) ProcessRoleCanvasesPermissions(canvasRepos []models.CanvasRepository, roleID uint64) (*[]CanvasRepoDefaultSerializer, error) {
	var permissionsList map[uint64]map[uint64]string
	canvasRepoViews := &[]CanvasRepoDefaultSerializer{}
	var err error
	collectionIDs := []uint64{}
	for _, repo := range canvasRepos {
		collectionIDs = append(collectionIDs, repo.CollectionID)
	}
	permissionsList = permissions.App.Service.CalculateCanvasRolePermissions(roleID, collectionIDs)
	for _, repo := range canvasRepos {
		repoPermissions := permissionsList[repo.ID]
		permissionValues := utils.Values(repoPermissions)
		hasPermission := false
		for _, perm := range permissionValues {
			if utils.Contains(permissiongroup.UserAccessCanvasPermissionsList, perm) {
				repoView := SerializeDefaultCanvasRepo(&repo)
				branchPerm := repoPermissions[repo.DefaultBranch.ID]
				repoView.DefaultBranch.Permission = branchPerm
				repoView.SearchMatch = true
				*canvasRepoViews = append(*canvasRepoViews, *repoView)
				hasPermission = true
				break
			}
		}
		if !hasPermission {
			repoView := SerializeDefaultCanvasRepo(&repo)
			branchPerm := permissiongroup.PGCanvasNone().SystemName
			repoView.DefaultBranch.Permission = branchPerm
			repoView.SearchMatch = true
			*canvasRepoViews = append(*canvasRepoViews, *repoView)
		}
	}
	if err != nil {
		return nil, err
	}
	return canvasRepoViews, nil
}

func (s canvasRepoService) ProcessUserCanvasesPermissions(canvasRepos []models.CanvasRepository, studioId uint64, userID uint64) (*[]CanvasRepoDefaultSerializer, error) {
	var permissionsList map[uint64]map[uint64]string
	var RoleActualPermsObject []RoleBranchActualPermissionsObject
	memberObject, _ := App.Repo.GetMemberByUserID(userID, studioId)
	canvasRepoViews := &[]CanvasRepoDefaultSerializer{}
	var err error
	for _, repo := range canvasRepos {
		if repo.ParentCanvasRepositoryID == nil {
			canvasPermissionsList, err := permissions.App.Service.CalculateCanvasRepoPermissions(userID, studioId, repo.CollectionID)
			if err != nil {
				return nil, err
			}
			permissionsList = utils.MergeMaps(permissionsList, canvasPermissionsList)
		} else {
			canvasPermissionsList, err := permissions.App.Service.CalculateSubCanvasRepoPermissions(userID, studioId, repo.CollectionID, *repo.ParentCanvasRepositoryID)
			if err != nil {
				return nil, err
			}
			permissionsList = utils.MergeMaps(permissionsList, canvasPermissionsList)
		}
	}
	for _, repo := range canvasRepos {
		repoPermissions := permissionsList[repo.ID]
		permissionValues := utils.Values(repoPermissions)
		hasPermission := false
		for _, perm := range permissionValues {
			if utils.Contains(permissiongroup.UserAccessCanvasPermissionsList, perm) {
				repoView := SerializeDefaultCanvasRepo(&repo)
				branchPerm := repoPermissions[repo.DefaultBranch.ID]
				repoView.DefaultBranch.Permission = branchPerm
				repoView.DefaultBranch.MemberPermsObject = MemberBranchPermissionActualCalculator(repo.CollectionID, repo.ID, repo.DefaultBranch.ID, memberObject, studioId)
				actualPerms := RoleBranchPermissionActualCalculator(repo.CollectionID, repo.ID, repo.DefaultBranch.ID, memberObject, studioId)
				for _, ap := range actualPerms {
					RoleActualPermsObject = append(RoleActualPermsObject, ap)
				}
				repoView.DefaultBranch.RolePermsObject = actualPerms
				*canvasRepoViews = append(*canvasRepoViews, *repoView)
				hasPermission = true
				break
			}
		}
		if !hasPermission {
			repoView := SerializeDefaultCanvasRepo(&repo)
			branchPerm := permissiongroup.PGCanvasNone().SystemName
			repoView.DefaultBranch.Permission = branchPerm
			repoView.DefaultBranch.RolePermsObject = []RoleBranchActualPermissionsObject{}
			repoView.DefaultBranch.MemberPermsObject = MemberBranchActualPermissionsObject{}

			*canvasRepoViews = append(*canvasRepoViews, *repoView)
		}
	}
	if err != nil {
		return nil, err
	}
	return canvasRepoViews, nil
}

func (s canvasRepoService) GetRoleCollectionsByID(collectionIDs []uint64, roleID uint64) ([]collection.CollectionSerializer, error) {
	collections, err := queries.App.CollectionQuery.GetCollections(map[string]interface{}{"id": collectionIDs, "is_archived": false})
	if err != nil {
		return nil, err
	}
	accessCollections, err := s.ProcessRoleCollectionsPermissions(collections, roleID)
	if err != nil {
		return nil, err
	}
	return *accessCollections, nil
}

func (s canvasRepoService) GetRoleCanvasesByIDs(canvasRepoIDs []uint64, roleID uint64) (*[]CanvasRepoDefaultSerializer, error) {
	var canvasRepos []models.CanvasRepository
	canvasRepos, err := queries.App.CanvasRepoQuery.GetCanvasRepos(map[string]interface{}{"id": canvasRepoIDs})
	if err != nil {
		return nil, err
	}
	canvasRepoViews, err := s.ProcessRoleCanvasesPermissions(canvasRepos, roleID)
	if err != nil {
		return nil, err
	}
	return canvasRepoViews, nil
}

func (s canvasRepoService) GetStudioMemberCollectionsByIDs(collectionIDs []uint64, studioId uint64, userID uint64) ([]collection.CollectionSerializer, error) {
	collections, err := queries.App.CollectionQuery.GetCollections(map[string]interface{}{"id": collectionIDs, "is_archived": false})
	if err != nil {
		return nil, err
	}
	accessCollections, err := s.ProcessUserCollectionsPermissions(collections, studioId, userID)
	if err != nil {
		return nil, err
	}
	return *accessCollections, nil
}

func (s canvasRepoService) GetMemberCanvasesByIDs(studioId uint64, userID uint64, repoIDs []uint64) (*[]CanvasRepoDefaultSerializer, error) {
	var canvasRepos []models.CanvasRepository
	canvasRepos, err := queries.App.CanvasRepoQuery.GetCanvasRepos(map[string]interface{}{"id": repoIDs})
	if err != nil {
		return nil, err
	}
	canvasRepoViews, err := s.ProcessUserCanvasesPermissions(canvasRepos, studioId, userID)
	if err != nil {
		return nil, err
	}
	return canvasRepoViews, nil
}

func (s canvasRepoService) ProcessSearchDump(records *[]models.CanvasRepoFullRow) ([]uint64, []uint64, []uint64) {
	collectionIds := []uint64{}
	repoIDs := []uint64{}
	supportRepoIDs := []uint64{}
	repoIDMap := make(map[uint64]bool)

	for _, v := range *records {
		collectionIds = append(collectionIds, v.ID2)
		repoIDs = append(repoIDs, v.ID)
		repoIDMap[v.ID] = true
	}

	for _, v := range *records {
		if v.ParentCanvasRepositoryID != nil && !repoIDMap[*v.ParentCanvasRepositoryID] {
			parentRecordsIDs := []uint64{}
			parentRecordsIDs, repoIDMap = s.GetParentRecordIDs(*v.ParentCanvasRepositoryID, repoIDMap, []uint64{})
			supportRepoIDs = append(supportRepoIDs, parentRecordsIDs...)
		}
	}
	return collectionIds, repoIDs, supportRepoIDs
}

func (s canvasRepoService) GetParentRecordIDs(canvasID uint64, repoIDMap map[uint64]bool, canvasParentRecordIDs []uint64) ([]uint64, map[uint64]bool) {
	canvasRepo, _ := queries.App.CanvasRepoQuery.GetCanvasRepoInstance(map[string]interface{}{"id": canvasID})
	if canvasRepo.ParentCanvasRepositoryID != nil && !repoIDMap[canvasRepo.ID] {
		canvasParentRecordIDs, repoIDMap = s.GetParentRecordIDs(*canvasRepo.ParentCanvasRepositoryID, repoIDMap, canvasParentRecordIDs)
	}
	if repoIDMap[canvasRepo.ID] {
		return canvasParentRecordIDs, repoIDMap
	}
	repoIDMap[canvasRepo.ID] = true
	canvasParentRecordIDs = append(canvasParentRecordIDs, canvasRepo.ID)
	return canvasParentRecordIDs, repoIDMap
}

func (s canvasRepoService) GetLanguageCanvasPrevAndNext(userID, canvasID uint64, language string) (*models.CanvasRepository, *models.CanvasRepository) {
	var nextCanvas *models.CanvasRepository
	var prevCanvas *models.CanvasRepository
	var langNextCanvas *models.CanvasRepository
	var langPrevCanvas *models.CanvasRepository
	canvasRepository, err := App.Repo.GetRepoWithCollection(map[string]interface{}{"id": canvasID})
	if err != nil {
		fmt.Println("Error in getting", canvasRepository.ID)
	}
	// Get the canvas repos position > canvasRepository.position order by position asc
	// if canvasRepos found take the first for the next canvas.
	// else
	// Get the collections position > canvasRepository collection position order by position asc
	// if collections found take the first collection and get the canvas of that collection by position order by position asc
	// if canvases found take the first canvas
	nextCanvas = s.GetLangNextSubCanvasPublicRepo(canvasRepository, language)
	if nextCanvas == nil {
		nextCanvases, _ := App.Repo.GetNextCanvases(canvasRepository, canvasRepository.Position)
		if len(nextCanvases) == 0 && canvasRepository.ParentCanvasRepositoryID != nil {
			nextCanvases, _ = App.Repo.GetNextCanvases(canvasRepository.ParentCanvasRepository, canvasRepository.ParentCanvasRepository.Position)
		}
		if len(nextCanvases) > 0 {
			for _, repo := range nextCanvases {
				fmt.Println(repo.ID, repo.Name, repo.DefaultBranch.PublicAccess)
				if repo.DefaultBranch.PublicAccess != models.PRIVATE {
					nextCanvas = &repo
					if nextCanvas != nil {
						langCanvases, _ := App.Repo.GetLangCanvases(nextCanvas.ID, language)
						if len(langCanvases) > 0 {
							langNextCanvas = &langCanvases[0]
							break
						}
					}
				}
				if repo.HasPublicCanvas {
					nextCanvas = s.GetLangNextSubCanvasPublicRepo(&repo, language)
					if nextCanvas != nil {
						langCanvases, _ := App.Repo.GetLangCanvases(nextCanvas.ID, language)
						if len(langCanvases) > 0 {
							langNextCanvas = &langCanvases[0]
						}
						break
					}
				}
			}
		}
		if (len(nextCanvases) > 0 && langNextCanvas == nil) || len(nextCanvases) <= 0 {
			nextCollections, _ := App.Repo.GetNextCollections(canvasRepository)
			if len(nextCollections) > 0 {
				nextCollection := nextCollections[0]
				nextCanvases, _ = App.Repo.GetCanvasesOrderByPositionAsc(map[string]interface{}{"collection_id": nextCollection.ID, "parent_canvas_repository_id": nil, "is_archived": false, "is_published": true, "is_language_canvas": false})
				for _, repo := range nextCanvases {
					if repo.DefaultBranch.PublicAccess != models.PRIVATE {
						langCanvases, _ := App.Repo.GetLangCanvases(repo.ID, language)
						if len(langCanvases) > 0 {
							nextCanvas = &repo
							langNextCanvas = &langCanvases[0]
							break
						}
					}
				}
			}
		}
	} else {
		if nextCanvas != nil {
			langCanvases, _ := App.Repo.GetLangCanvases(nextCanvas.ID, language)
			if len(langCanvases) > 0 {
				langNextCanvas = &langCanvases[0]
			}
		}
	}

	// Get the canvas repos position < canvasRepository.position order by position desc
	// if canvasRepos found take the first for the next canvas.
	// else
	// Get the collections position < canvasRepository collection position order by position desc
	// if collections found take the first collection and get the canvas of that collection by position order by position desc
	// if canvases found take the first canvas
	sameLevel := true
	prevCanvases, _ := App.Repo.GetPrevCanvases(canvasRepository, canvasRepository.Position)
	if len(prevCanvases) == 0 && canvasRepository.ParentCanvasRepositoryID != nil {
		sameLevel = false
		prevCanvases, _ = App.Repo.GetPrevCanvases(canvasRepository.ParentCanvasRepository, canvasRepository.ParentCanvasRepository.Position+1)
	}
	if len(prevCanvases) > 0 {
		prevCanvas = s.GetLangPrevSubCanvasPublicRepo(prevCanvases, sameLevel, language)
		if prevCanvas != nil {
			langCanvases, _ := App.Repo.GetLangCanvases(prevCanvas.ID, language)
			if len(langCanvases) > 0 {
				langPrevCanvas = &langCanvases[0]
			}
		}
	} else {
		prevCollections, _ := App.Repo.GetPrevCollections(canvasRepository)
		if len(prevCollections) > 0 {
			prevCollection := prevCollections[0]
			prevCanvases, _ = App.Repo.GetCanvasesOrderByPositionDesc(map[string]interface{}{"collection_id": prevCollection.ID, "parent_canvas_repository_id": nil, "is_archived": false, "is_published": true, "is_language_canvas": false})
			for _, repo := range prevCanvases {
				if repo.DefaultBranch.PublicAccess != models.PRIVATE || repo.HasPublicCanvas {
					prevCanvas = &repo
					if prevCanvas != nil {
						langCanvases, _ := App.Repo.GetLangCanvases(prevCanvas.ID, language)
						if len(langCanvases) > 0 {
							langPrevCanvas = &langCanvases[0]
							break
						}
					}
				}
			}
			if prevCanvas != nil {
				prevCanvas = s.GetCollectionLastCanvas(prevCanvas)
				if prevCanvas != nil {
					langCanvases, _ := App.Repo.GetLangCanvases(prevCanvas.ID, language)
					if len(langCanvases) > 0 {
						langPrevCanvas = &langCanvases[0]
					}
				}
			}
		}
	}
	return langNextCanvas, langPrevCanvas
}
