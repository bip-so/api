package global

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gitlab.com/phonepost/bip-be-platform/internal/follow"
	"gitlab.com/phonepost/bip-be-platform/internal/member"
	"gitlab.com/phonepost/bip-be-platform/internal/message"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/internal/reactions"
	"gitlab.com/phonepost/bip-be-platform/internal/reel"
	"gitlab.com/phonepost/bip-be-platform/internal/role"
	"gitlab.com/phonepost/bip-be-platform/internal/studio"
	"gitlab.com/phonepost/bip-be-platform/internal/studiopermissions"
	"gitlab.com/phonepost/bip-be-platform/internal/user"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/s3"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/search"
	"gorm.io/gorm"
	"io"
	"path/filepath"
)

var (
	GlobalController = globalController{}
)

type globalController struct{}

func (gc globalController) CreateStudioController(validator *studio.CreateStudioValidator, user *models.User) (*models.Studio, error) {
	logger.Info("Creating a studio")
	std, err := studio.App.Controller.CreateStudioController(validator, user)
	if err != nil {
		logger.Info("Error in creating a studio")
		logger.Debug(err.Error())
		return nil, err
	}

	var memberId uint64
	// Action 1: Create a member instance on studio creation
	memberId = member.App.Controller.CreateStudioMemberController(user.ID, std.ID)
	memberInstance := member.App.Controller.GetMemberInstance(memberId)
	if memberId == 0 {
		return nil, err
	}
	logger.Info("New member ID : " + fmt.Sprintf("%d", memberId))

	// Action 2: Create and attach Admin Role to StudioPerms (2 steps)
	// Admin Role
	roleId, err := role.RoleBasicController.CreateDefaultStudioRole(std.ID, memberInstance.ID)
	fmt.Println("Role just created: ", roleId)
	// StudioPermission (PG:"Admin")
	// Create a StudioPermssions Object
	// Move pg_studio_admin = Contact
	_, err2 := studiopermissions.StudioPermissionService.NewStudioPermission(std.ID, "pg_studio_admin", &roleId, nil, false)
	if err2 != nil {
		return nil, err2
	}

	// Action3: Create StudioMember Role, Permission
	memberRoleId, err := role.RoleBasicController.CreateDefaultStudioMemberRole(std.ID, []models.Member{*memberInstance})
	billingRoleID, err := role.RoleBasicController.CreateBillingRole(std.ID, []models.Member{*memberInstance})
	// only for debug
	fmt.Println(billingRoleID)

	fmt.Println("Role just created: ", roleId)
	_, err2 = studiopermissions.StudioPermissionService.NewStudioPermission(std.ID, "pg_studio_none", &memberRoleId, nil, false)
	if err2 != nil {
		return nil, err2
	}

	// Invalidating the redis cache
	permissions.App.Service.InvalidateStudioPermissionCache(user.ID)

	// Post studio creation setup
	err = PostStudioSetup(std, user)
	if err != nil {
		return nil, err
	}

	return std, err
}

func (gc globalController) CheckHandleAvailable(handle string) (bool, error) {

	_, err := user.App.Repo.GetUserByUsername(handle)
	if err == nil {
		return false, nil
	}
	if err != gorm.ErrRecordNotFound {
		return false, err
	}
	_, err = studio.App.StudioRepo.GetStudioByHandle(handle)
	if err == nil {
		return false, nil
	}
	if err != gorm.ErrRecordNotFound {
		return false, err
	}
	return true, nil
}

func (gc globalController) PopularUsersController(authUser *models.User, skipInt int) (*[]models.User, *[]models.FollowUser, error) {

	users, err := user.App.Controller.PopularUsersController(skipInt)
	if err != nil {
		return nil, nil, err
	}

	if authUser == nil {
		return users, nil, nil
	}

	userIDs := []uint64{}
	for _, usr := range *users {
		userIDs = append(userIDs, usr.ID)
	}

	followUsers, err := follow.App.Repo.GetIsUserFollowingUser(authUser.ID, userIDs)
	if err != nil {
		return nil, nil, err
	}

	return users, followUsers, nil
}

func (gc globalController) SearchController(query, objectType string, authUser *models.User, skipInt int, studioID uint64) ([]studio.StudioDocument, []user.UserDocument, []role.RoleDocument, []reel.ReelsSerialData, error) {
	var studioDocs []studio.StudioDocument
	var userDocs []user.UserDocument
	var rolesDocs []role.RoleDocument
	var reelDocs []reel.ReelDocument
	reelsData := &[]reel.ReelsSerialData{}
	roles, _ := GetRolesByStudioID(studioID, query)
	if roles != nil {
		rolesDocs = role.GetRoleSearch(&roles)
	}

	if objectType == search.StudioDocumentIndexName || objectType == "" {
		queryResults, err := search.GetIndex(search.StudioDocumentIndexName).Search(query, skipInt)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		err = queryResults.UnmarshalHits(&studioDocs)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		studioIDs := []uint64{}
		for _, std := range studioDocs {
			studioIDs = append(studioIDs, std.ID)
		}
		if len(studioIDs) > 0 && authUser != nil {
			members, _ := queries.App.MemberQuery.GetMembersOfUserInMultipleStudios(studioIDs, authUser.ID)
			if members != nil {
				memMap := map[uint64]*models.Member{}
				for _, member := range *members {
					memMap[member.StudioID] = &member
				}
				varFalse := false
				for i, stdio := range studioDocs {
					if mem, exists := memMap[stdio.ID]; exists {
						isJoined := mem.UserID == authUser.ID
						studioDocs[i].IsJoined = &isJoined
					} else {
						studioDocs[i].IsJoined = &varFalse
					}
					studioDocs[i].MembersCount = StudioMembersCount(stdio.ID)
					studioDocs[i].IsRequested = studio.App.StudioService.CheckIsRequested(authUser.ID, stdio.ID)
				}
			}
		} else {
			for i, stdio := range studioDocs {
				studioDocs[i].MembersCount = StudioMembersCount(stdio.ID)
			}
		}
		if objectType != "" {
			return studioDocs, nil, rolesDocs, nil, nil
		}
	}
	if objectType == search.UserDocumentIndexName || objectType == "" {
		queryResults, err := search.GetIndex(search.UserDocumentIndexName).Search(query, skipInt)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		err = queryResults.UnmarshalHits(&userDocs)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		userIDs := []uint64{}
		for _, usr := range userDocs {
			userIDs = append(userIDs, usr.ID)
		}
		if len(userIDs) > 0 && authUser != nil {
			followUsers, _ := follow.App.Repo.GetIsUserFollowingUser(authUser.ID, userIDs)
			if followUsers != nil {
				followMap := map[uint64]*models.FollowUser{}
				for _, fUser := range *followUsers {
					followMap[fUser.UserId] = &fUser
				}
				varFalse := false
				for i, usrDoc := range userDocs {
					resp, _ := follow.App.Controller.GetUserFollowFollowCountHandler(usrDoc.ID)
					userDocs[i].Followers = resp.Followers
					userDocs[i].Following = resp.Following
					if fUser, exists := followMap[usrDoc.ID]; exists {
						isFollowing := fUser.FollowerId == authUser.ID
						userDocs[i].IsFollowing = &isFollowing
					} else {
						userDocs[i].IsFollowing = &varFalse
					}
				}
			}
		}
		if objectType != "" {
			return nil, userDocs, rolesDocs, nil, nil
		}
	}
	if objectType == search.ReelDocumentIndexName || objectType == "" {
		queryResults, err := search.GetIndex(search.ReelDocumentIndexName).Search(query, skipInt)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		err = queryResults.UnmarshalHits(&reelDocs)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		reelIDs := []uint64{}
		for _, reelD := range reelDocs {
			reelIDs = append(reelIDs, reelD.ID)
		}
		if len(reelIDs) > 0 && authUser != nil {
			reels, _ := reel.App.Repo.GetReelsByIDs(reelIDs)
			reelsData, _ = reel.App.Service.GetReelsWithConfigData(reels, authUser)
		} else if len(reelIDs) > 0 && authUser == nil {
			reels, _ := reel.App.Repo.GetReelsByIDs(reelIDs)
			reelReactions, _ := reactions.App.Repo.GetReelReactionByIDs(reelIDs)
			reelsData = reel.SerializeDefaultManyReelsWithReactionsForUser(reels, reelReactions, nil, nil, nil)
		}
		if objectType != "" {
			return nil, nil, rolesDocs, *reelsData, nil
		}
	}
	return studioDocs, userDocs, rolesDocs, *reelsData, nil
}

func (gc *globalController) UpdateImage(file io.Reader, model string, uuid2 string, fileName string, studioID uint64, repoID uint64) (*string, error) {

	var studioImagePath string
	// Get Extention and Name
	extension := filepath.Ext(fileName)
	name := fileName[0 : len(fileName)-len(extension)]
	// CleanFile Name
	updateFileName := slug.Make(name)

	newFileName := updateFileName + extension

	if model == "blocks" {
		// Check is studio id is bot 0
		if studioID == 0 {
			studioImagePath = fmt.Sprintf("nostudio/%s/%s/%s", model, uuid2, newFileName)
		} else {
			studioImagePath = fmt.Sprintf("%d/%d/%s/%s/%s", studioID, repoID, model, uuid2, newFileName)
		}

	} else if model == "canvasrepocover" {
		id := uuid.New()
		randomFileName := id.String() + extension
		studioImagePath = fmt.Sprintf("cover-images/%s", randomFileName)
	} else {
		studioImagePath = fmt.Sprintf("%s/%s/%s", model, uuid2, newFileName)
	}

	fmt.Println(studioImagePath)
	response, err := s3.UploadObjectToBucket(studioImagePath, file, true)
	if err != nil {
		logger.Error(fmt.Sprintf("Error on updating studio image url %s", err.Error()))
		return nil, err
	}
	return &response, err
}

func (gc *globalController) GetMessages(userID uint64, skip int) (*[]SerializedMessage, error) {
	messages, err := message.GetMessages(userID, skip)
	if err != nil {
		return nil, err
	}

	return SerializeMessages(messages), nil
}

func (gc *globalController) DeleteMessage(userID uint64, messageID uint64) error {
	err := message.DeleteMessageById(userID, messageID)
	if err != nil {
		return err
	}

	return nil
}
