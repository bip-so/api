package auth

import (
	"encoding/json"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/canvasbranch"
	"gitlab.com/phonepost/bip-be-platform/internal/workflows"
	"gorm.io/datatypes"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"gitlab.com/phonepost/bip-be-platform/internal/canvasrepo"
	"gitlab.com/phonepost/bip-be-platform/internal/collection"
	"gitlab.com/phonepost/bip-be-platform/internal/feed"
	"gitlab.com/phonepost/bip-be-platform/internal/member"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/internal/role"
	"gitlab.com/phonepost/bip-be-platform/internal/studio"
	"gitlab.com/phonepost/bip-be-platform/internal/studiopermissions"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"gorm.io/gorm"
)

const MaxStudioLength = 20

// PostStudioSetup : After a studio is created we create a collection, canvas repo, canvas branch, 1 Paragraph Block.
func (s *authCreateStudioService) PostStudioSetup(std *models.Studio, user *models.User) error {
	// Create Default Documents in Personal studio
	err := s.CreateDefaultDocsInPersonalStudio(std, user)
	if err != nil {
		fmt.Println("Error in creating default docs", err)
		return err
	}
	// Create a collection
	//collectionInstance, err := shared.WorkflowCreateCollectionAndPerms(
	//	"My new collection",
	//	2,
	//	"private",
	//	std.CreatedByID,
	//	std.ID)
	//if err != nil {
	//	return err
	//}
	//// Create a Repo
	//_, errCreatingBranch := shared.WorkflowCreateCanvasRepoInsideCollection(
	//	collectionInstance.ID,
	//	collectionInstance.CreatedByID,
	//	"Default Canvas",
	//	"📋",
	//	1,
	//	std.ID,
	//	*user,
	//)
	//if errCreatingBranch != nil {
	//	return errCreatingBranch
	//}
	//
	//collectionInstance, err = queries.App.CollectionQuery.UpdateCollection(collectionInstance.ID, map[string]interface{}{
	//	"computed_root_canvas_count": 1,
	//	"computed_all_canvas_count":  1,
	//})

	feed.App.Service.JoinStudio(std.ID, std.CreatedByID)
	return nil
}

func (s *authCreateStudioService) CreateMember(member *models.Member) uint64 {
	//result := postgres.GetDB().Create(member)
	//return member.ID, result.Error
	if err := postgres.GetDB().Create(&member).Error; err != nil {
		log.Fatalln(err)
		return 0
	}

	return member.ID
}

func (s *authCreateStudioService) NewMember(userId uint64, studioId uint64) *models.Member {
	return &models.Member{
		UserID:      userId,
		StudioID:    studioId,
		CreatedByID: userId,
		UpdatedByID: userId,
	}
}

func (s authCreateStudioService) GetStudioByHandle(handle string) (*models.Studio, error) {
	//var studio models.Studio
	//err := s.db.Model(models.Studio{}).Where("handle = ? and is_archived = ?", handle, false).Preload("Topics").First(&studio).Error
	studio, err := queries.App.StudioQueries.GetStudioQuery(map[string]interface{}{
		"handle":      handle,
		"is_archived": false,
	})
	return studio, err
}
func (s authCreateStudioService) checkHandleAvailability(handle string) bool {
	//var user models.User
	//err := s.db.Model(&models.User{}).Where("username = ?", handle).First(&user).Error
	_, err := queries.App.UserQueries.GetUser(map[string]interface{}{"username": handle})
	if err == nil {
		return false
	}
	if err != gorm.ErrRecordNotFound {
		return false
	}
	_, err = s.GetStudioByHandle(handle)
	if err == nil {
		return false
	}
	if err != gorm.ErrRecordNotFound {
		return false
	}
	return true
}

func (s *authCreateStudioService) CreateDefaultStudio(user *models.User) {
	var studioHandle string
	if len(user.Username) > MaxStudioLength {
		// username is greater than allowed len
		studioHandle = user.Username + "-" + utils.HandleExtender(3)
		exists := s.checkHandleAvailability(studioHandle)
		if !exists {
			studioHandle = user.Username + "-" + utils.HandleExtender(4)
		}
	} else {
		studioHandle = user.Username + "-" + "studio"
		exists := s.checkHandleAvailability(studioHandle)
		if !exists {
			studioHandle = studioHandle + "-" + utils.HandleExtender(4)
		}
	}
	// remove space from the studioHandle
	studioHandleCleaned := strings.ReplaceAll(studioHandle, " ", "-")

	//studioHandle := user.Username // can be changed later
	// ADD THIS TO A LOOP FOR A DUP CHECKER
	//handlesClean := reg.ReplaceAllString(studioHandle, "")

	// We need to create a studio for this user
	studioName := user.Username + "'s Studio"
	studioDesc := "My personal workspace" // can be changed later

	newStudio := queries.App.StudioQueries.NewStudioInstance(studioName, studioHandleCleaned, studioDesc, "", "", user.ID)
	errCreatingStudio := studio.App.StudioRepo.CreateStudio(newStudio)
	if errCreatingStudio != nil {
		logger.Info("StudioRepo.CreateStudio / errCreatingStudio" + errCreatingStudio.Error())
	}

	// Action 1: Create a member instance on studio creation
	var memberId uint64
	memberId = member.App.Controller.CreateStudioMemberController(user.ID, newStudio.ID)
	memberInstance := member.App.Controller.GetMemberInstance(memberId)
	// logger.Info("New member ID : " + fmt.Sprintf("%d", memberId))
	fmt.Println("New member ID : ", memberId)

	// Action 2: Create and attach Admin Role to StudioPerms (2 steps)
	// Admin Role
	roleId, _ := role.RoleBasicController.CreateDefaultStudioRole(newStudio.ID, memberInstance.ID)
	fmt.Println("Role just created: ", roleId)
	// StudioPermission (PG:"Admin")
	// Create a StudioPermssions Object
	// Move pg_studio_admin = Contact
	_, err2 := studiopermissions.StudioPermissionService.NewStudioPermission(newStudio.ID, "pg_studio_admin", &roleId, nil, false)
	fmt.Println(err2)

	// Action3: Create StudioMember Role, Permission
	memberRoleId, _ := role.RoleBasicController.CreateDefaultStudioMemberRole(newStudio.ID, []models.Member{*memberInstance})
	fmt.Println("Role just created: ", roleId)
	_, err2 = studiopermissions.StudioPermissionService.NewStudioPermission(newStudio.ID, "pg_studio_none", &memberRoleId, nil, false)
	fmt.Println(err2)
	permissions.App.Service.InvalidateStudioPermissionCache(user.ID)
	err23 := s.PostStudioSetup(newStudio, user)
	if err23 != nil {
		fmt.Println(err23.Error())
	}
	// Final step
	//App.Repo.updateUserDefaultStudio(user.ID, newStudio.ID)
	queries.App.UserQueries.UpdateUser(user.ID, map[string]interface{}{"default_studio_id": newStudio.ID}, true)
	queries.App.StudioQueries.DeleteUserAssociatedStudioDataByUserID(user.ID)
}

func (s *authCreateStudioService) CreateDefaultDocsInPersonalStudio(std *models.Studio, user *models.User) error {
	welcomeToBipCanvasBranchID := uint64(7361)
	playgroundCanvasBranchID := uint64(7360)
	howBipHelpCommunitiesCanvasBranchID := uint64(7359)
	communityWorkspaceCanvasBranchID := uint64(7358)
	if configs.GetConfigString("APP_MODE") == "production" {
		welcomeToBipCanvasBranchID = uint64(34294)
		playgroundCanvasBranchID = uint64(36407)
		howBipHelpCommunitiesCanvasBranchID = uint64(34296)
		communityWorkspaceCanvasBranchID = uint64(34299)
	}

	collectionView := &collection.CollectionCreateValidator{
		Name:         "INTRODUCTION",
		Position:     1,
		PublicAccess: "private",
	}
	collectionInstance, err := collection.App.Controller.CreateCollectionController(collectionView, std.CreatedByID, std.ID)
	if err != nil {
		fmt.Println("Error in creating collection", err)
		return err
	}

	welcomeToBipCanvasView := canvasrepo.InitCanvasRepoPost{
		CollectionID: collectionInstance.ID,
		Name:         "Welcome to bip",
		Icon:         "",
		Position:     1,
	}
	welcomeToBipCanvas, errCreatingCanvasBranch := workflows.WorkflowHelperInitCanvasRepo(workflows.InitCanvasRepoPost{
		CollectionID:             welcomeToBipCanvasView.CollectionID,
		Name:                     welcomeToBipCanvasView.Name,
		Icon:                     welcomeToBipCanvasView.Icon,
		Position:                 welcomeToBipCanvasView.Position,
		ParentCanvasRepositoryID: welcomeToBipCanvasView.ParentCanvasRepositoryID,
	}, collectionInstance.CreatedByID, collectionInstance.StudioID, *user)
	if errCreatingCanvasBranch != nil {
		return errCreatingCanvasBranch
	}
	err = BlocksCloner(welcomeToBipCanvasBranchID, *welcomeToBipCanvas.DefaultBranchID, user, welcomeToBipCanvas.ID)
	if err != nil {
		fmt.Println("Error in cloning blocks of welcome to bip", err)
	}

	playgroundCanvasView := canvasrepo.InitCanvasRepoPost{
		CollectionID: collectionInstance.ID,
		Name:         "Playground",
		Icon:         "",
		Position:     2,
	}
	playgrounCanvas, errCreatingCanvasBranch := workflows.WorkflowHelperInitCanvasRepo(workflows.InitCanvasRepoPost{
		CollectionID:             playgroundCanvasView.CollectionID,
		Name:                     playgroundCanvasView.Name,
		Icon:                     playgroundCanvasView.Icon,
		Position:                 playgroundCanvasView.Position,
		ParentCanvasRepositoryID: playgroundCanvasView.ParentCanvasRepositoryID,
	}, collectionInstance.CreatedByID, collectionInstance.StudioID, *user)
	if errCreatingCanvasBranch != nil {
		return errCreatingCanvasBranch
	}
	err = BlocksCloner(playgroundCanvasBranchID, *playgrounCanvas.DefaultBranchID, user, playgrounCanvas.ID)
	if err != nil {
		fmt.Println("Error in cloning blocks of playground", err)
	}

	howBipHelpCommunitiesCanvasView := canvasrepo.InitCanvasRepoPost{
		CollectionID: collectionInstance.ID,
		Name:         "How bip helps communities",
		Icon:         "",
		Position:     3,
	}
	howBipHelpCommunitiesCanvas, errCreatingCanvasBranch := workflows.WorkflowHelperInitCanvasRepo(workflows.InitCanvasRepoPost{
		CollectionID:             howBipHelpCommunitiesCanvasView.CollectionID,
		Name:                     howBipHelpCommunitiesCanvasView.Name,
		Icon:                     howBipHelpCommunitiesCanvasView.Icon,
		Position:                 howBipHelpCommunitiesCanvasView.Position,
		ParentCanvasRepositoryID: howBipHelpCommunitiesCanvasView.ParentCanvasRepositoryID,
	}, collectionInstance.CreatedByID, collectionInstance.StudioID, *user)
	if errCreatingCanvasBranch != nil {
		return errCreatingCanvasBranch
	}
	err = BlocksCloner(howBipHelpCommunitiesCanvasBranchID, *howBipHelpCommunitiesCanvas.DefaultBranchID, user, howBipHelpCommunitiesCanvas.ID)
	if err != nil {
		fmt.Println("Error in cloning blocks of howBipHelpCommunitiesCanvas", err)
	}

	communityWorkspaceCanvasView := canvasrepo.InitCanvasRepoPost{
		CollectionID: collectionInstance.ID,
		Name:         "Create your community Workspace",
		Icon:         "",
		Position:     4,
	}
	communityWorkspaceCanvas, errCreatingCanvasBranch := workflows.WorkflowHelperInitCanvasRepo(workflows.InitCanvasRepoPost{
		CollectionID:             communityWorkspaceCanvasView.CollectionID,
		Name:                     communityWorkspaceCanvasView.Name,
		Icon:                     communityWorkspaceCanvasView.Icon,
		Position:                 communityWorkspaceCanvasView.Position,
		ParentCanvasRepositoryID: communityWorkspaceCanvasView.ParentCanvasRepositoryID,
	}, collectionInstance.CreatedByID, collectionInstance.StudioID, *user)
	if errCreatingCanvasBranch != nil {
		return errCreatingCanvasBranch
	}
	err = BlocksCloner(communityWorkspaceCanvasBranchID, *communityWorkspaceCanvas.DefaultBranchID, user, communityWorkspaceCanvas.ID)
	if err != nil {
		fmt.Println("Error in cloning blocks of welcome to bip", err)
	}

	canvasBranchMap := map[uint64]uint64{
		// Prod mapping
		34294: *welcomeToBipCanvas.DefaultBranchID,
		36407: *playgrounCanvas.DefaultBranchID,
		34296: *howBipHelpCommunitiesCanvas.DefaultBranchID,
		34299: *communityWorkspaceCanvas.DefaultBranchID,

		// Stage mapping
		7361: *welcomeToBipCanvas.DefaultBranchID,
		7360: *playgrounCanvas.DefaultBranchID,
		7359: *howBipHelpCommunitiesCanvas.DefaultBranchID,
		7358: *communityWorkspaceCanvas.DefaultBranchID,
	}
	UpdateMentionsInBlocks(*welcomeToBipCanvas.DefaultBranchID, canvasBranchMap, user, std)
	UpdateMentionsInBlocks(*playgrounCanvas.DefaultBranchID, canvasBranchMap, user, std)
	UpdateMentionsInBlocks(*howBipHelpCommunitiesCanvas.DefaultBranchID, canvasBranchMap, user, std)

	collectionInstance, err = queries.App.CollectionQuery.UpdateCollection(
		collectionInstance.ID, map[string]interface{}{
			"computed_root_canvas_count": 4,
			"computed_all_canvas_count":  4,
		})

	// Publishing the branches
	canvasbranch.App.Service.PublishCanvasBranch(*welcomeToBipCanvas.DefaultBranchID, user, true)
	canvasbranch.App.Service.PublishCanvasBranch(*playgrounCanvas.DefaultBranchID, user, true)
	canvasbranch.App.Service.PublishCanvasBranch(*howBipHelpCommunitiesCanvas.DefaultBranchID, user, true)
	canvasbranch.App.Service.PublishCanvasBranch(*communityWorkspaceCanvas.DefaultBranchID, user, true)

	return nil
}

func BlocksCloner(fromBranchId uint64, toBranchId uint64, user *models.User, newCanvasRepoID uint64) error {
	// Get the blocks
	postgres.GetDB().Where("canvas_branch_id = ?", toBranchId).Delete(models.Block{})
	ogblocks, err := queries.App.BlockQuery.GetBlocksByBranchID(fromBranchId)
	if err != nil {
		return err
	}
	// If Size of og blocks is 0 we return.
	if len(*ogblocks) == 0 {
		return nil
	}
	var newBlocks []models.Block
	for _, block := range *ogblocks {
		newBlock := CloneBlockInstance(block, toBranchId, user, newCanvasRepoID)
		newBlocks = append(newBlocks, newBlock)
	}
	errBulkCreating := postgres.GetDB().Create(&newBlocks).Error
	if errBulkCreating != nil {
		return errBulkCreating
	}
	// All the Blocks are created now, We need to move things
	_, errGettingClonedBlock := queries.App.BlockQuery.GetBlocksByBranchID(toBranchId)
	if errGettingClonedBlock != nil {
		return errGettingClonedBlock
	}
	return nil
}

func CloneBlockInstance(ogBlockInstance models.Block, branchID uint64, user *models.User, newCanvasRepoID uint64) models.Block {
	var block models.Block
	block = ogBlockInstance
	block.ClonedFromBlockID = ogBlockInstance.ID
	block.ID = 0 // Reset the PK
	block.UUID = uuid.New()
	block.CanvasBranchID = &branchID
	block.CanvasRepositoryID = newCanvasRepoID
	block.CreatedByUser = nil
	block.UpdatedByUser = nil
	block.CreatedByID = 0
	block.UpdatedByID = 0
	block.ArchivedByID = 0
	block.CreatedByID = user.ID
	block.UpdatedByID = user.ID
	block.CreatedAt = time.Now()
	block.UpdatedAt = time.Now()
	contributors := []map[string]interface{}{
		{
			"id":        user.ID,
			"uuid":      user.UUID.String(),
			"repoID":    newCanvasRepoID,
			"branchID":  branchID,
			"fullName":  user.FullName,
			"username":  user.Username,
			"avatarUrl": user.AvatarUrl,
			"timestamp": time.Now(),
		},
	}
	contributorsStr, _ := json.Marshal(contributors)
	block.Contributors = contributorsStr
	block.CommentCount = 0
	return block
}

func UpdateMentionsInBlocks(branchID uint64, canvasBranchMap map[uint64]uint64, user *models.User, studio *models.Studio) {
	newblocks, err := GetBlocksByBranchID(branchID)
	if err != nil {
		fmt.Println("Error in getting blocks of branch", err)
		return
	}
	for blockIndex, block := range newblocks {
		if block.Type == models.BlockSimpleTableV1 {
			newblocks[blockIndex].Children = DefaultTableBlockChildren(block.UUID.String())
			continue
		}
		var blockChildren []map[string]interface{}
		json.Unmarshal(block.Children, &blockChildren)
		pageMentionPresent := false
		for childIndex, children := range blockChildren {
			fmt.Println(children["type"])
			if children["type"] != nil && children["type"].(string) == "pageMention" {
				oldMention := children["mention"].(map[string]interface{})
				newBranchID := canvasBranchMap[uint64(oldMention["id"].(float64))]
				newBranch, _ := queries.App.BranchQuery.GetBranchByID(newBranchID)
				newMention := map[string]interface{}{
					"id":                    newBranch.ID,
					"key":                   newBranch.Key,
					"name":                  newBranch.Name,
					"type":                  "branch",
					"uuid":                  newBranch.UUID.String(),
					"repoID":                newBranch.CanvasRepositoryID,
					"repoKey":               newBranch.CanvasRepository.Key,
					"repoName":              strings.TrimSuffix(newBranch.CanvasRepository.Name, "\n"),
					"repoUUID":              newBranch.CanvasRepository.UUID.String(),
					"studioID":              newBranch.CanvasRepository.StudioID,
					"createdByUserID":       user.ID,
					"createdByUserFullName": user.FullName,
					"createdByUserUsername": user.Username,
				}
				blockChildren[childIndex]["mention"] = newMention
				blockChildren[childIndex]["uuid"] = uuid.New().String()
				pageMentionPresent = true
			} else if children["type"] != nil && children["type"].(string) == "userMention" {
				fmt.Println("came inside userMention", user.FullName)
				newMention := map[string]interface{}{
					"id":                    user.ID,
					"type":                  "user",
					"uuid":                  user.UUID.String(),
					"fullName":              user.FullName,
					"studioID":              studio.ID,
					"username":              user.Username,
					"avatarUrl":             user.AvatarUrl,
					"createdByUserID":       user.ID,
					"createdByUserFullName": user.FullName,
					"createdByUserUsername": user.FullName,
				}
				blockChildren[childIndex]["mention"] = newMention
				blockChildren[childIndex]["uuid"] = uuid.New().String()
				blockChildren[childIndex]["children"] = []map[string]string{
					{
						"text": fmt.Sprintf("<@%s>", user.FullName),
					},
				}
				pageMentionPresent = true
			}
		}
		if pageMentionPresent {
			blockChildrenStr, _ := json.Marshal(blockChildren)
			newblocks[blockIndex].Children = blockChildrenStr
		}
	}
	postgres.GetDB().Save(newblocks)
}

func GetBlocksByBranchID(branchID uint64) ([]models.Block, error) {
	var blocks []models.Block
	err := postgres.GetDB().Model(&models.Block{}).Where("canvas_branch_id = ?", branchID).Preload("CreatedByUser").Preload("UpdatedByUser").Order("rank ASC").Find(&blocks).Error
	if err != nil {
		return nil, err
	}
	return blocks, nil
}

func DefaultTableBlockChildren(tableUUID string) datatypes.JSON {
	children := fmt.Sprintf(`[
	  {
		"type": "table-row",
		"uuid": "6db85d00-022b-47d2-a106-a9a8372d0892",
		"children": [
		  {
			"type": "table-cell",
			"uuid": "98d9b17f-ea72-43bd-b3a9-8748d9c6e022",
			"rowUUID": "6db85d00-022b-47d2-a106-a9a8372d0892",
			"children": [
			  {
				"rank": 28625,
				"type": "ulist",
				"uuid": "84460445-80d5-4ba5-876a-87af94b07c1d",
				"rowUUID": "6db85d00-022b-47d2-a106-a9a8372d0892",
				"cellUUID": "98d9b17f-ea72-43bd-b3a9-8748d9c6e022",
				"children": [
				  {
					"text": "Bullet points"
				  }
				],
				"tableUUID": "%s",
				"attributes": {}
			  },
			  {
				"type": "ulist",
				"uuid": "191002a8-bf25-4767-adfe-54d2a5024983",
				"rowUUID": "6db85d00-022b-47d2-a106-a9a8372d0892",
				"cellUUID": "98d9b17f-ea72-43bd-b3a9-8748d9c6e022",
				"children": [
				  {
					"text": "They work on table too"
				  }
				],
				"tableUUID": "%s",
				"attributes": {},
				"contributors": []
			  }
			],
			"tableUUID": "%s"
		  },
		  {
			"type": "table-cell",
			"uuid": "9e27599c-34c1-40be-a9e5-ebbbb69adb2b",
			"rowUUID": "6db85d00-022b-47d2-a106-a9a8372d0892",
			"children": [
			  {
				"type": "text",
				"uuid": "99d41154-d100-4db6-87e9-ffe0c5020efb",
				"rowUUID": "6db85d00-022b-47d2-a106-a9a8372d0892",
				"cellUUID": "9e27599c-34c1-40be-a9e5-ebbbb69adb2b",
				"children": [
				  {
					"text": ""
				  }
				],
				"tableUUID": "%s",
				"attributes": {}
			  }
			],
			"tableUUID": "%s"
		  }
		],
		"tableUUID": "%s"
	  },
	  {
		"type": "table-row",
		"uuid": "64a9d954-4fae-4d84-84ab-bed8e291d02b",
		"children": [
		  {
			"type": "table-cell",
			"uuid": "8b7f0690-331e-4355-b272-7d151d3c2261",
			"rowUUID": "64a9d954-4fae-4d84-84ab-bed8e291d02b",
			"children": [
			  {
				"type": "checklist",
				"uuid": "949ae1ef-648d-43be-981d-e8f1906afbaa",
				"rowUUID": "64a9d954-4fae-4d84-84ab-bed8e291d02b",
				"cellUUID": "15de6e60-9ef3-4b75-a799-b4e64b9529a4",
				"children": [
				  {
					"text": "Yes"
				  }
				],
				"tableUUID": "%s",
				"attributes": {
				  "level": 1,
				  "checked": false
				}
			  }
			],
			"tableUUID": "%s"
		  },
		  {
			"type": "table-cell",
			"uuid": "f3305953-2e58-4f77-8c66-93b9eb403046",
			"rowUUID": "64a9d954-4fae-4d84-84ab-bed8e291d02b",
			"children": [
			  {
				"rank": 28937,
				"type": "text",
				"uuid": "d753e136-4f1b-4708-9b49-5bd7370343f3",
				"rowUUID": "edd3b064-45e0-4ba9-bf7b-6d53f188a101",
				"cellUUID": "4ac07596-aab2-46e4-aac2-bb2877fc5054",
				"children": [
				  {
					"text": ""
				  }
				],
				"tableUUID": "%s",
				"updatedAt": "2022-11-09T13:51:11.601Z",
				"attributes": {
				  "level": 1
				},
				"updatedById": 254611,
				"contributors": []
			  }
			],
			"tableUUID": "%s"
		  }
		],
		"tableUUID": "%s"
	  }
	]`, tableUUID, tableUUID, tableUUID, tableUUID, tableUUID, tableUUID, tableUUID, tableUUID, tableUUID, tableUUID, tableUUID)
	return []byte(children)
}
