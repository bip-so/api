package studio

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"io"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gosimple/slug"

	"gitlab.com/phonepost/bip-be-platform/internal/feed"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/supabase"

	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/internal/member"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/kafka"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/s3"
)

func (sc *studioController) CreateStudioController(validator *CreateStudioValidator, user *models.User) (*models.Studio, error) {
	logger.Info("Studio: CreateStudioController")
	studio, err := App.StudioService.NewStudio(validator.Name, validator.Handle, validator.Description, validator.Website, "", user.ID)
	if err != nil {
		logger.Info("StudioService.NewStudio " + err.Error())
		return nil, err
	}

	topicNames := []string{}
	topicNames = append(topicNames, validator.Topics...)
	if len(topicNames) > 0 {
		existingTopics, _ := App.TopicRepo.FindTopics(topicNames)
		allTopics := *existingTopics
		if len(*existingTopics) != len(topicNames) {
			var topics []models.Topic
			for _, name := range topicNames {
				found := false
				for _, eTopic := range *existingTopics {
					if eTopic.Name == name {
						found = true
						break
					}
				}
				if found {
					continue
				}
				topic, err := App.TopicService.NewTopic(name)
				if err != nil {
					logger.Info("TopicService.NewTopic " + err.Error())
					continue
				}
				topics = append(topics, *topic)
			}
			if len(topics) != 0 {
				err = App.TopicRepo.CreateTopics(&topics)
				if err != nil {
					logger.Info("TopicRepo.Create " + err.Error())
					return nil, err
				}
			}
			allTopics = append(allTopics, topics...)
		}
		studio.Topics = allTopics
	}

	err = App.StudioRepo.CreateStudio(studio)
	if err != nil {
		logger.Info("StudioRepo.CreateStudio " + err.Error())
		return nil, err
	}
	go func() {
		stdData, _ := json.Marshal(studio)
		kafkaClient := kafka.GetKafkaClient()
		kafkaClient.Publish(configs.KAFKA_TOPICS_NEW_STUDIO, strconv.FormatUint(studio.ID, 10), stdData)
	}()

	return studio, nil
}

func (sc *studioController) GetStudioController(studioID uint64, authUser *models.User) (*models.Studio, *models.Member, error) {
	studio, err := App.StudioRepo.GetStudioByID(studioID)
	if err != nil {
		return nil, nil, err
	}
	member, err := queries.App.MemberQuery.GetMember(map[string]interface{}{"studio_id": studioID, "user_id": authUser.ID})
	if err != nil {
		return studio, nil, nil
	}
	return studio, member, err
}

func (sc *studioController) UpdateStudioImage(studioID uint64, file io.Reader, fileName string) (*models.Studio, error) {
	studio, err := App.StudioRepo.GetStudioByID(studioID)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	// Get Extention and Name
	extension := filepath.Ext(fileName)
	name := fileName[0 : len(fileName)-len(extension)]
	// CleanFile Name
	updateFileName := slug.Make(name)
	newFileName := updateFileName + extension

	studioImagePath := fmt.Sprintf("studio/%s/%s", studio.UUID, newFileName)
	response, err := s3.UploadImageToBucket(studioImagePath, file, true, true)
	if err != nil {
		logger.Error(err.Error())
	}

	studio, err = App.StudioRepo.UpdateStudioByID(studioID, map[string]interface{}{"image_url": response.URL})
	if err != nil {
		logger.Error(fmt.Sprintf("Error on updating studio image url %s", err.Error()))
	}
	return studio, err
}

func (sc *studioController) ToggleStudioMembershipController(studioID uint64, loggedInUserID uint64) {
	var newState bool
	data := map[string]interface{}{}
	existingStudioInstance, _ := App.StudioRepo.GetStudioByID(studioID)
	currentState := existingStudioInstance.AllowPublicMembership
	if currentState {
		newState = false
	} else {
		newState = true
	}
	data["allow_public_membership"] = newState
	data["updated_by_id"] = loggedInUserID
	App.StudioRepo.UpdateStudioByID(studioID, data)
	// Ok so we are now converting all the requests to Member now.
	if newState {
		// We need to now loop the pending request and make them members
		pendingRequestForStudio, _ := App.StudioService.GetRequestStudioList(studioID)
		if len(*pendingRequestForStudio) != 0 {
			var UserIds []uint64
			for _, req := range *pendingRequestForStudio {
				UserIds = append(UserIds, req.UserID)
				// Delete the request
				_ = postgres.GetDB().Delete(req)
			}
			// Convert  this request to Member now
			_, _ = queries.App.MemberQuery.AddMembersInBulkWithUserID(UserIds, studioID)
		}
	}

}

// todo : Fix spelling
func (sc *studioController) editStudioControler(studioID uint64, validator *UpdateStudioValidator) (*models.Studio, error) {
	existingStudioInstance, err := App.StudioRepo.GetStudioByID(studioID)
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{}
	if validator.Name != existingStudioInstance.DisplayName {
		if len(validator.Name) > 100 {
			return nil, StudioErrorsNameLenMax
		}
		data["display_name"] = validator.Name
	}
	if validator.Handle != existingStudioInstance.Handle {
		if len(validator.Handle) < 2 {
			return nil, StudioErrorsHandleInvalid
		}
		if len(validator.Handle) > 36 {
			return nil, StudioErrorsHandleLenMax
		}
		isValid := models.CheckHandleValidity(validator.Handle)
		if !isValid {
			return nil, StudioErrorsHandleInvalid
		}
		available, err := App.StudioRepo.checkHandleAvailablity(validator.Handle)
		if err != nil {
			return nil, err
		}
		if !available {
			return nil, ErrHandleUnavailable
		}
		data["handle"] = validator.Handle
	}
	if validator.Description != existingStudioInstance.Description {
		if len(validator.Description) > 300 {
			return nil, StudioErrorsDescriptionLenMax
		}
		data["description"] = validator.Description
	}
	if validator.Website != existingStudioInstance.Website {
		data["website"] = validator.Website
	}
	addedTopicNames, removedTopicNames := existingStudioInstance.TopicsDiff(validator.Topics)
	if len(addedTopicNames)+len(removedTopicNames) != 0 {
		existingTopics, err := App.TopicRepo.FindTopics(validator.Topics)
		if err != nil {
			return nil, err
		}
		notExistingTopics := []models.Topic{}
		existingTopicNamesMap := map[string]*models.Topic{}
		for _, topic := range *existingTopics {
			existingTopicNamesMap[topic.Name] = &topic
		}
		for _, topicName := range validator.Topics {
			if _, exists := existingTopicNamesMap[topicName]; !exists {
				topic, err := App.TopicService.NewTopic(topicName)
				if err != nil {
					return nil, err
				}
				notExistingTopics = append(notExistingTopics, *topic)
			}
		}
		if len(notExistingTopics) != 0 {
			err = App.TopicRepo.CreateTopics(&notExistingTopics)
			if err != nil {
				logger.Info("TopicRepo.Create " + err.Error())
				return nil, err
			}
		}
		allTopics := append(*existingTopics, notExistingTopics...)
		err = App.TopicRepo.UpdateStudioTopics(existingStudioInstance, allTopics)
		if err != nil {
			return nil, err
		}
	}

	studioInstance, err := App.StudioRepo.UpdateStudioByID(studioID, data)
	if err != nil {
		return nil, err
	}

	return studioInstance, err
}

// Todo: Delete Studio on Stripe too
func (sc *studioController) deleteStudioController(studioID uint64, deletedById uint64) error {
	err := App.StudioRepo.DeleteStudio(studioID, deletedById)
	return err
}

func (sc *studioController) integrations(c *gin.Context) {
	return
}

func (sc *studioController) imageURL(c *gin.Context) {
	return
}

// studios list to be shown on explore page
func (sc *studioController) GetPopularStudioController(authUser *models.User, skip, limit int) ([]models.Studio, *[]models.Member, error) {
	cache := redis.NewCache()
	cacheKey := "popular:studios"
	var members *[]models.Member
	if value := cache.Get(context.Background(), cacheKey); value != nil {
		var studios []models.Studio
		json.Unmarshal([]byte(value.(string)), &studios)
		if len(studios) > limit+skip {
			studios = studios[skip : limit+skip]
		} else {
			studios = studios[skip:]
		}
		if authUser != nil {
			var err error
			studioIDs := []uint64{}
			for _, std := range studios {
				studioIDs = append(studioIDs, std.ID)
			}
			members, err = queries.App.MemberQuery.GetMembersOfUserInMultipleStudios(studioIDs, authUser.ID)
			if err != nil {
				return nil, nil, err
			}
		}
		return studios, members, nil
	}
	studios, err := App.StudioRepo.GetPopularStudios()
	if err != nil {
		return nil, nil, err
	}
	if authUser != nil {
		studioIDs := []uint64{}
		for _, std := range studios {
			studioIDs = append(studioIDs, std.ID)
		}
		members, err = queries.App.MemberQuery.GetMembersOfUserInMultipleStudios(studioIDs, authUser.ID)
		if err != nil {
			return nil, nil, err
		}
	}
	cacheData := studios
	go func() {
		studiosData, _ := json.Marshal(cacheData)
		cache.Set(context.Background(), cacheKey, studiosData, &redis.Options{
			Expiration: 24 * time.Hour,
		})
	}()
	studios = studios[skip : limit+skip]
	return studios, members, nil
}

// studio list icons to be shown in extreme left rail based on user
func (sc *studioController) studioList(c *gin.Context) {
	return
}

func (sc *studioController) GetStudioByHandleController(handle string, authUser *models.User) (*models.Studio, *[]models.Member, error) {
	studio, err := App.StudioRepo.GetStudioByHandle(handle)
	if err != nil {
		return nil, nil, err
	}
	if authUser == nil {
		return studio, nil, err
	}
	members, err := queries.App.MemberQuery.GetMembersOfUserInMultipleStudios([]uint64{studio.ID}, authUser.ID)
	if err != nil {
		return studio, nil, err
	}
	return studio, members, err
}

func (sc *studioController) JoinStudioController(user *models.User, studioId uint64) error {
	// This is two tasks
	// Add this user to Members Table
	// Add this user to The "Member" Role on the Studio

	var err error
	// Get the Member on This Studio
	mem, err := queries.App.MemberQuery.GetMember(map[string]interface{}{"user_id": user.ID, "studio_id": studioId})
	// If member found
	if err == nil && mem != nil {
		if mem.IsRemoved {
			return errors.New("can't join as you were banned")
		} else if mem.HasLeft {
			return queries.App.MemberQuery.JoinStudio([]uint64{user.ID}, studioId)
		}

	} else {
		// Member not found

		member := queries.App.MemberQuery.AddUserIDToStudio(user.ID, studioId)
		if member == nil {
			return errors.New("Error in adding user as member to studio")
		}
		memberObj, err := queries.App.MemberQuery.GetMember(map[string]interface{}{"id": member.ID})
		if err != nil {
			return err
		}
		err = queries.App.MemberQuery.AddMembersToStudioInMemberRole(studioId, []models.Member{*memberObj})
		if err != nil {
			return err
		}
	}
	queries.App.StudioQueries.DeleteUserAssociatedStudioDataByUserID(user.ID)
	supabase.UpdateUserSupabase(user.ID, true)
	feed.App.Service.JoinStudio(studioId, user.ID)

	return nil
}

func (sc *studioController) GetStudioStatsController(studioId uint64) map[string]interface{} {
	return App.StudioService.StudioStats(studioId)
}

func (sc *studioController) MemberCountController(studioId uint64) (uint64, error) {
	count, err := queries.App.MemberQuery.GetMemberCountForStudio(studioId)
	return uint64(count), err
}

func (sc *studioController) JoinStudioInBulkController(body JoinStudioBulkPost, studioId uint64, authUserId uint64) ([]models.Member, error) {
	var err error
	members, err := queries.App.MemberQuery.GetMembersByUserIDs(body.UsersAdded, studioId)
	if err != nil {
		return nil, err
	}
	var addUserIds []uint64
	var joinbackUserIds []uint64
	for _, userId := range body.UsersAdded {

		flag := false
		for _, memb := range members {

			if userId == memb.UserID {

				flag = true
				if memb.IsRemoved {
					continue
				} else if memb.HasLeft {
					joinbackUserIds = append(joinbackUserIds, userId)
				}
				break
			}
		}
		if !flag {
			addUserIds = append(addUserIds, userId)
		}

	}

	// add userIds to studio

	_, err = queries.App.MemberQuery.AddMembersInBulkWithUserID(addUserIds, studioId)

	// update hasLeft
	err = queries.App.MemberQuery.JoinStudio(joinbackUserIds, studioId)

	var allUserIDs []uint64
	allUserIDs = append(allUserIDs, addUserIds...)
	allUserIDs = append(allUserIDs, joinbackUserIds...)
	members, err = queries.App.MemberQuery.GetMembersByStudioIDandUserIDs(studioId, allUserIDs, 0)
	// Finally
	// Ok this is an edge case we are implementing here.
	// We are checking if this user has a request pending.
	for _, userIdx := range allUserIDs {
		var quickcheck models.StudioMembersRequest
		_ = postgres.GetDB().Model(&models.StudioMembersRequest{}).Where(map[string]interface{}{
			"studio_id": studioId,
			"user_id":   userIdx,
		}).First(&quickcheck).Error
		if quickcheck.ID != 0 {
			// Delete this Instance
			_ = postgres.GetDB().Delete(quickcheck)
		}
	}
	// End of Edge Case
	// Sending an event to supabase so frontend can update the list
	// @todo later move this to kafka based on member joined or updated.
	go func() {
		for _, userID := range allUserIDs {
			queries.App.StudioQueries.DeleteUserAssociatedStudioDataByUserID(userID)
			supabase.UpdateUserSupabase(userID, true)

			notifications.App.Service.PublishNewNotification(notifications.StudioInviteByName, authUserId, []uint64{userID}, &studioId,
				nil, notifications.NotificationExtraData{}, nil, nil)
		}
		// follow studio feed
		if len(allUserIDs) > 0 {
			feed.App.Service.BulkJoinStudio(studioId, allUserIDs)
		}
	}()
	return members, err
}

func (sc *studioController) StudioAdminMembers(studioId uint64) ([]member.MemberSerializer, error) {
	studioAdminRole, err := queries.App.RoleQuery.GetStudioAdminRole(studioId)
	if err != nil {
		return nil, err
	}
	members := member.BulkSerializeMembers(studioAdminRole.Members)
	return members, nil
}

func (sc *studioController) GetRequestToJoinStudioListController(studioId uint64) (*[]models.StudioMembersRequest, error) {
	return App.StudioService.GetRequestStudioList(studioId)
}

func (sc *studioController) CreateRequestToJoinStudioController(studioID uint64, loggedInUserID uint64) error {
	joinStudioRequestObject := App.StudioService.MockStudioMembershipRequestObject(studioID, loggedInUserID, loggedInUserID)
	return App.StudioService.CreateStudioRequestInstance(joinStudioRequestObject)
}

func (sc *studioController) RejectRequestToJoinStudioController(membershipRequestID uint64, loggedInUserID uint64) (*models.StudioMembersRequest, error) {
	membershipRequest, err := App.StudioRepo.UpdateMembershipRequestByID(membershipRequestID, map[string]interface{}{"action": "Rejected", "action_by_id": loggedInUserID})
	return membershipRequest, err
}

func (sc *studioController) AcceptRequestToJoinStudioController(membershipRequestID uint64, loggedInUserID uint64) (*models.StudioMembersRequest, error) {
	membershipRequest, err := App.StudioRepo.UpdateMembershipRequestByID(membershipRequestID, map[string]interface{}{"action": "Accepted", "action_by_id": loggedInUserID})
	return membershipRequest, err
}

func (sc *studioController) StudioStats(studioID uint64) map[string]interface{} {
	return App.StudioService.BuildStudioStats(studioID)
}
