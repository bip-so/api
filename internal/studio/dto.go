package studio

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"gorm.io/gorm"
)

// *****************
//    Studio
// *****************

type StudioRepoInterface interface {
	CreateStudio(studio *models.Studio) error
	GetStudioByID(studioID uint64) (*models.Studio, error)
	GetStudioByName(handle string) (*models.Studio, error)
	UpdateStudioByID(studioID uint64, updates map[string]interface{}) error
	DeleteStudio(studioID uint64, deletedById uint64) error
	GetPopularStudios() (*[]models.Studio, error)
	GetStudiosByIDs(studioIDs []uint64) (*[]models.Studio, error)
}

func (sr studioRepo) CreateStudio(studio *models.Studio) error {
	err := sr.db.Create(studio).Error
	fmt.Println(err)
	return err
}

func (sr studioRepo) GetStudioByID(studioID uint64) (*models.Studio, error) {
	var studio models.Studio
	err := sr.db.Model(models.Studio{}).Where("id = ? and is_archived = ?", studioID, false).Preload("Topics").First(&studio).Error
	return &studio, err
}

func (sr studioRepo) GetStudioByHandle(handle string) (*models.Studio, error) {
	var studio models.Studio
	err := sr.db.Model(models.Studio{}).Where("LOWER(handle) = ? and is_archived = ?", strings.ToLower(handle), false).Preload("Topics").First(&studio).Error
	return &studio, err
}

func (sr studioRepo) UpdateStudioByID(studioID uint64, updates map[string]interface{}) (*models.Studio, error) {
	var studioInstance *models.Studio
	err := sr.db.Model(models.Studio{}).Where(models.Studio{BaseModel: models.BaseModel{ID: studioID}}).Updates(updates).Preload("Topics").Find(&studioInstance).Error
	if err == nil {
		go func() {
			std, _ := sr.GetStudioByID(studioID)
			stdData, _ := json.Marshal(std)
			sr.kafka.Publish(configs.KAFKA_TOPICS_UPDATE_STUDIO, strconv.FormatUint(std.ID, 10), stdData)
		}()
	}
	return studioInstance, err
}

func (sr studioRepo) GetPopularStudios() ([]models.Studio, error) {
	var studios []models.Studio
	//err := sr.db.Model(models.Studio{}).Where("is_archived = ?", false).Order("computed_follower_count desc").Preload("Topics").Limit(configs.PAGINATION_LIMIT).Find(&studios).Error
	err := sr.db.Debug().Preload("Topics").Table("studios").
		Select("studios.*, COUNT(DISTINCT members.id) mem, COUNT(DISTINCT canvas_repositories.id) canvasCount").
		Joins("LEFT OUTER JOIN canvas_repositories on studios.id = canvas_repositories.studio_id").
		Joins("LEFT OUTER JOIN canvas_branches on canvas_repositories.id = canvas_branches.canvas_repository_id").
		Joins("LEFT OUTER JOIN members ON studios.id = members.studio_id").
		Joins("LEFT OUTER JOIN collections ON canvas_repositories.collection_id = collections.id").
		Where("members.has_left = false AND members.is_removed = false and canvas_branches.public_access <> 'private' and canvas_repositories.is_archived = false and canvas_repositories.is_published = true and studios.is_archived = false and collections.is_archived = false").
		Group("studios.id").
		Order("canvasCount desc, mem desc").
		Find(&studios).Limit(100).Error
	return studios, err
}

func (sr studioRepo) GetStudioRepoCount(id uint64) int64 {
	var count int64
	sr.db.Model(models.CanvasRepository{}).Where("studio_id = ? AND is_archived = false", id).Count(&count)
	return count
}

func (sr studioRepo) GetAllStudios() (*[]models.Studio, error) {
	var studios []models.Studio
	err := sr.db.Model(models.Studio{}).Where("is_archived = ?", false).Find(&studios).Error
	return &studios, err
}

func (sr studioRepo) GetStudiosByIDs(studioIDs []uint64) (*[]models.Studio, error) {
	var studios []models.Studio
	err := sr.db.Model(models.Studio{}).Where("id IN ? and is_archived = ?", studioIDs, false).Find(&studios).Error
	return &studios, err
}

func (sr studioRepo) DeleteStudio(studioID uint64, deletedById uint64) error {
	var studio *models.Studio
	updates := map[string]interface{}{
		"is_archived":    true,
		"archived_at":    time.Now(),
		"archived_by_id": deletedById,
	}
	err := sr.db.Model(models.Studio{}).Where(models.Studio{BaseModel: models.BaseModel{ID: studioID}}).Updates(updates).First(&studio).Error
	if err == nil {
		go func() {
			stdData, _ := json.Marshal(studio)
			sr.kafka.Publish(configs.KAFKA_TOPICS_DELETED_STUDIO, strconv.FormatUint(studio.ID, 10), stdData)
		}()
	}
	return err
}

func (sr studioRepo) checkHandleAvailablity(handle string) (bool, error) {
	var user models.User
	err := sr.db.Model(&models.User{}).Where("username = ?", handle).First(&user).Error
	if err == nil {
		return false, nil
	}
	if err != gorm.ErrRecordNotFound {
		return false, err
	}
	_, err = sr.GetStudioByHandle(handle)
	if err == nil {
		return false, nil
	}
	if err != gorm.ErrRecordNotFound {
		return false, err
	}
	return true, nil
}

func (sr studioRepo) studioMembersCount(studioID uint64) int64 {
	var count int64
	sr.db.Model(models.Member{}).Where("studio_id = ? AND has_left = false AND is_removed = false", studioID).Count(&count)
	return count
}

// *****************
//    Topic
// *****************

type TopicRepoInterface interface {
	CreateTopics(topics *[]models.Topic) error
	FindTopics(topicNames []string) (*[]models.Topic, error)
}

func (tr studioTopicRepo) CreateTopics(topics *[]models.Topic) error {
	err := tr.db.Create(topics).Error
	return err
}

func (tr studioTopicRepo) FindTopics(topicNames []string) (*[]models.Topic, error) {
	var topics []models.Topic
	err := tr.db.Where("name IN ?", topicNames).Find(&topics).Error
	return &topics, err
}

func (tr studioTopicRepo) UpdateStudioTopics(studio *models.Studio, topics []models.Topic) error {
	err := tr.db.Model(studio).Association("Topics").Replace(topics)
	if err == nil {
		go func() {
			std, _ := App.StudioRepo.GetStudioByID(studio.ID)
			stdData, _ := json.Marshal(std)
			tr.kafka.Publish(configs.KAFKA_TOPICS_UPDATE_STUDIO, strconv.FormatUint(std.ID, 10), stdData)
		}()
	}
	return err
}

// ********
//  User Associated Studios
// ********

type Interface interface {
	CreateUserAssociatedStudio(userStudios *models.UserAssociatedStudio) error
	GetUserAssociatedStudioDataByUserID(userID uint64) (*models.UserAssociatedStudio, error)
}

func (ur *userAssociatedStudioRepo) CreateUserAssociatedStudio(userStudios *models.UserAssociatedStudio) error {
	err := ur.db.Create(userStudios).Error
	return err
}

func (ur *userAssociatedStudioRepo) GetUserAssociatedStudioDataByUserID(userID uint64) (*models.UserAssociatedStudio, error) {
	var userStudios models.UserAssociatedStudio
	err := ur.db.Where(models.UserAssociatedStudio{UserID: userID}).First(&userStudios).Error
	return &userStudios, err
}

func (sr studioRepo) GetUserFollows(userID, followerID uint64) (*models.FollowUser, error) {
	var userFollow *models.FollowUser
	err := sr.db.Model(&models.FollowUser{}).Where("user_id = ? AND follower_id = ?", userID, followerID).Find(&userFollow).Error
	return userFollow, err
}

func (sr studioRepo) UpdateMembershipRequestByID(membershipRequestID uint64, updates map[string]interface{}) (*models.StudioMembersRequest, error) {
	var instance *models.StudioMembersRequest
	err := sr.db.Model(models.StudioMembersRequest{}).Where(models.StudioMembersRequest{BaseModel: models.BaseModel{ID: membershipRequestID}}).Updates(updates).First(&instance).Error
	return instance, err
}

func (sr studioRepo) GetUserStudioMemberRequest(userID, studioID uint64) (*models.StudioMembersRequest, error) {
	var instance *models.StudioMembersRequest
	err := sr.db.Model(models.StudioMembersRequest{}).Where("user_id = ? and studio_id = ? and action = ?", userID, studioID, "Pending").First(&instance).Error
	return instance, err
}
