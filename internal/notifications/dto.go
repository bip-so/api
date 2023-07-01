package notifications

import (
	"fmt"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

func (r notificationRepo) GetUserSettings(userID uint64) (*[]models.UserSettings, error) {
	var userSettings *[]models.UserSettings
	err := r.db.Model(&models.UserSettings{}).Where(&models.UserSettings{UserID: userID}).Preload("User").Find(&userSettings).Error
	return userSettings, err
}

func (r notificationRepo) CreateNewUserSettings(userSetting *models.UserSettings) bool {
	result := r.db.Create(userSetting)
	return result.RowsAffected == 1
}

func (r notificationRepo) GetUser(userID uint64) (models.User, error) {
	var user models.User
	err := r.db.Model(&models.User{}).Where("id = ?", userID).First(&user).Error
	return user, err
}

func (r notificationRepo) GetUserSocialAuth(userID uint64, providerName string) (*models.UserSocialAuth, error) {
	var userSocialAuth *models.UserSocialAuth
	err := r.db.Model(&models.UserSocialAuth{}).Where("user_id = ? and provider_name = ?", userID, providerName).First(&userSocialAuth).Error
	return userSocialAuth, err
}

func (r notificationRepo) GetUserByUUID(uuid string) (models.User, error) {
	var user models.User
	err := r.db.Model(&models.User{}).Where("uuid = ?", uuid).Preload("User").First(&user).Error
	return user, err
}

func (r notificationRepo) Create(notification *models.Notification) error {
	err := r.db.Create(&notification).Error
	return err
}

func (r notificationRepo) Get(query map[string]interface{}) (*[]models.Notification, error) {
	var data *[]models.Notification
	err := r.db.Model(models.Notification{}).Where(query).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r notificationRepo) GetNotification(query map[string]interface{}) (*models.Notification, error) {
	var data *models.Notification
	err := r.db.Model(models.Notification{}).Where(query).First(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r notificationRepo) GetMessage(query map[string]interface{}) (*models.Message, error) {
	var data *models.Message
	err := r.db.Model(models.Message{}).Where(query).First(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r notificationRepo) UpdateNotifications(query, updates map[string]interface{}) error {
	err := r.db.Model(models.Notification{}).Where(query).Updates(updates).Error
	if err != nil {
		return err
	}
	return nil
}

func (r notificationRepo) Update(query, updates map[string]interface{}) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.Model(models.Notification{}).Where(query).Updates(updates).Find(&notifications).Error
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r notificationRepo) Delete(query map[string]interface{}) error {
	var notifications models.Notification
	err := r.db.Where(query).Delete(&notifications).Error
	if err != nil {
		return err
	}
	return nil
}

func (r notificationRepo) CreateNotificationCount(notificationCount *models.NotificationCount) error {
	err := r.db.Create(&notificationCount).Error
	return err
}

func (r notificationRepo) GetNotificationCount(query map[string]interface{}) (*models.NotificationCount, error) {
	var data models.NotificationCount
	err := r.db.Model(models.NotificationCount{}).Where(query).First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r notificationRepo) GetNotificationCountForUser(userID uint64) map[string]interface{} {
	var counts struct {
		All      int64
		Personal int64
	}
	countList := map[string]interface{}{}
	studioCountList := []studioCount{} //map[string]interface{}{}
	err := r.db.Raw(`SELECT ( SELECT count(*) FROM
									notifications
								WHERE
									notifier_id = ?
									AND is_archived = FALSE
									AND seen = FALSE
									AND is_website_notification = TRUE ) AS all,
								( SELECT count(*) FROM
									notifications
								WHERE
									notifier_id = ?
									AND is_archived = FALSE
									AND seen = FALSE
									AND is_personal = TRUE
									AND is_website_notification = TRUE ) AS personal`, userID, userID).Find(&counts).Error
	if err != nil {
		fmt.Println(" error while getting counts for user ", err)
		return countList
	}
	userStudios, err := r.GetUserStudiosByID(userID)
	if err != nil {
		fmt.Println(" error while getting studio for user ", err)
		return countList
	}
	r.db.Model(&models.Notification{}).
		Debug().
		Select("studio_id, count(studio_id) as count").
		Where(map[string]interface{}{"notifier_id": userID, "is_archived": false, "seen": false, "is_personal": false}).
		Group("studio_id").
		Find(&studioCountList)
	countList["all"] = counts.All
	countList["personal"] = counts.Personal
	studioCountListMap := App.Service.mapStudioWithCount(&userStudios, studioCountList)
	countList["studio"] = studioCountListMap
	return countList
}

func (r notificationRepo) SaveNotificationCount(notificationCount *models.NotificationCount) error {
	updates := map[string]interface{}{
		"all":      notificationCount.All,
		"personal": notificationCount.Personal,
		"studio":   notificationCount.Studio,
	}
	err := r.db.Model(models.NotificationCount{}).Where("id = ?", notificationCount.ID).Updates(updates).Error
	return err
}

func (r notificationRepo) GetStudioByID(studioID uint64) (*models.Studio, error) {
	var studio models.Studio
	err := r.db.Model(models.Studio{}).Where("id = ? and is_archived = ?", studioID, false).First(&studio).Error
	return &studio, err
}

func (r notificationRepo) GetUserStudiosByID(userID uint64) ([]models.Studio, error) {
	var members []models.Member
	result := postgres.GetDB().Model(&models.Member{}).Where("user_id = ? AND has_left = false AND is_removed = false", userID).Preload("Studio").Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	var studios []models.Studio
	for _, member := range members {
		studios = append(studios, *member.Studio)
	}
	return studios, nil
}

func (r notificationRepo) AddNotificationToDbAndStream(notification *models.Notification) error {
	err := r.Create(notification)
	if err != nil {
		return err
	}
	err = App.Service.AddNotificationToStream(notification)
	if err != nil {
		return err
	}
	return nil
}

func (r notificationRepo) GetMembersByStudioID(studioID uint64) ([]models.Member, error) {
	var members []models.Member
	result := postgres.GetDB().Select("user_id").Model(&models.Member{}).
		Where("studio_id = ? AND has_left = false AND is_removed = false", studioID).
		Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

func (r notificationRepo) GetNotificationCountByUserIDs(userIDs []uint64) (*[]models.NotificationCount, error) {
	var data *[]models.NotificationCount
	err := r.db.Model(models.NotificationCount{}).Where("user_id IN ?", userIDs).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r notificationRepo) GetRolesByID(roleIDs *[]uint64) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Model(&models.Role{}).Where("id IN ?", *roleIDs).Preload("Members").Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r notificationRepo) GetBlockReactionByID(blockReactionId uint64) (*models.BlockReaction, error) {
	var blockReaction models.BlockReaction
	err := r.db.Model(&models.BlockReaction{}).Where("id = ?", blockReactionId).
		Preload("CanvasBranch").
		Preload("CanvasBranch.CanvasRepository").
		First(&blockReaction).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &blockReaction, nil
}

func (r notificationRepo) GetReelReactionByID(reelReactionId uint64) (*models.ReelReaction, error) {
	var reelReaction models.ReelReaction
	err := r.db.Model(&models.ReelReaction{}).Where("id = ?", reelReactionId).
		Preload("Reel").
		Preload("Reel.Block").
		First(&reelReaction).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &reelReaction, nil
}

func (r notificationRepo) GetBlockThreadCommentByID(blockThreadCommentId uint64) (*models.BlockComment, error) {
	var repo models.BlockComment
	err := r.db.Model(&models.BlockComment{}).Where("id = ?", blockThreadCommentId).
		Preload("Thread").
		Preload("Thread.CanvasRepository").
		Preload("Thread.Block").
		First(&repo).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &repo, nil
}

func (r notificationRepo) GetBlockCommentsByThreadID(threadID uint64) ([]models.BlockComment, error) {
	var repo []models.BlockComment
	err := r.db.Model(&models.BlockComment{}).Where("thread_id = ?", threadID).
		Preload("Thread").
		Preload("Thread.CanvasRepository").
		First(&repo).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return repo, nil
}

func (r notificationRepo) GetReelCommentReactionByID(reelCommentReactionID uint64) (models.ReelCommentReaction, error) {
	var reelCommentReaction models.ReelCommentReaction
	err := r.db.Model(models.ReelCommentReaction{}).Where("id = ?", reelCommentReactionID).Preload("Reel").Preload("Reel.Block").Preload("ReelComment").First(&reelCommentReaction).Error
	return reelCommentReaction, err
}

func (r notificationRepo) GetReelByID(reelID uint64) (models.Reel, error) {
	var reel models.Reel
	err := r.db.Model(models.Reel{}).Where("id = ?", reelID).Preload("CanvasRepository").Preload("Block").First(&reel).Error
	return reel, err
}

func (r notificationRepo) GetReelCommentByID(reelCommentID uint64) (models.ReelComment, error) {
	var reelComment models.ReelComment
	err := r.db.Model(models.ReelComment{}).Where("id = ?", reelCommentID).Preload("Reel").Preload("Reel.Block").First(&reelComment).Error
	return reelComment, err
}

func (r notificationRepo) GetReelComments(query map[string]interface{}) ([]models.ReelComment, error) {
	var reelComments []models.ReelComment
	err := r.db.Model(models.ReelComment{}).Where(query).Preload("Reel").First(&reelComments).Error
	return reelComments, err
}

func (r notificationRepo) GetBlockThreadByID(blockThreadId uint64) (models.BlockThread, error) {
	var blockThread models.BlockThread
	err := r.db.Model(models.BlockThread{}).Where("id = ?", blockThreadId).Preload("CanvasRepository").Preload("Block").First(&blockThread).Error
	return blockThread, err
}

func (r notificationRepo) GetBlockThreadReactionByID(blockThreadReactionId uint64) (models.BlockThreadReaction, error) {
	var blockThreadReaction models.BlockThreadReaction
	err := r.db.Model(models.BlockThreadReaction{}).Where("id = ?", blockThreadReactionId).Preload("BlockThread").Preload("BlockThread.Block").Preload("BlockThread.CanvasRepository").First(&blockThreadReaction).Error
	return blockThreadReaction, err
}

func (r notificationRepo) GetBlockThreadCommentReactionByID(blockThreadCommentReactionId uint64) (models.BlockCommentReaction, error) {
	var blockThreadCommentReaction models.BlockCommentReaction
	err := r.db.Model(models.BlockCommentReaction{}).Where("id = ?", blockThreadCommentReactionId).Preload("BlockComment").Preload("BlockThread").Preload("BlockThread.CanvasRepository").Preload("BlockThread.Block").First(&blockThreadCommentReaction).Error
	return blockThreadCommentReaction, err
}

func (r notificationRepo) GetBlockByID(blockId uint64) (models.Block, error) {
	var block models.Block
	err := r.db.Model(models.Block{}).Where("id = ?", blockId).Preload("CanvasBranch").Preload("CanvasBranch.CanvasRepository").First(&block).Error
	return block, err
}

func (r notificationRepo) GetBlocks(query map[string]interface{}) ([]models.Block, error) {
	var blocks []models.Block
	err := r.db.Model(models.Block{}).Where(query).Find(&blocks).Error
	return blocks, err
}

func (r notificationRepo) GetMergeRequestByID(mergeRequestId uint64) (models.MergeRequest, error) {
	var mergeRequest models.MergeRequest
	err := r.db.Model(models.MergeRequest{}).Where("id = ?", mergeRequestId).First(&mergeRequest).Error
	return mergeRequest, err
}

func (r notificationRepo) GetCollectionByID(collectionId uint64) (models.Collection, error) {
	var collection models.Collection
	err := r.db.Model(models.Collection{}).Where("id = ?", collectionId).First(&collection).Error
	return collection, err
}

func (r notificationRepo) GetCanvasRepoByID(canvasId uint64) (models.CanvasRepository, error) {
	var canvasRepo models.CanvasRepository
	err := r.db.Model(models.CanvasRepository{}).Where("id = ?", canvasId).Preload("Studio").First(&canvasRepo).Error
	return canvasRepo, err
}

func (r notificationRepo) GetAccessRequestByID(accessRequestID uint64) (models.AccessRequest, error) {
	var accessRequest models.AccessRequest
	err := r.db.Model(models.AccessRequest{}).Where("id = ?", accessRequestID).Preload("Studio").First(&accessRequest).Error
	return accessRequest, err
}

func (r notificationRepo) GetCanvasBranchByID(branchID uint64) (models.CanvasBranch, error) {
	var canvasBranch models.CanvasBranch
	err := r.db.Model(models.CanvasBranch{}).Where("id = ?", branchID).First(&canvasBranch).Error
	return canvasBranch, err
}

func (r notificationRepo) GetCanvasBranchByIDPreload(branchID uint64) (models.CanvasBranch, error) {
	var canvasBranch models.CanvasBranch
	err := r.db.Model(models.CanvasBranch{}).Where("id = ?", branchID).Preload("CanvasRepository").First(&canvasBranch).Error
	return canvasBranch, err
}

func (r notificationRepo) GetCollectionModeratorUserIDs(collectionID uint64) ([]uint64, error) {
	var userIDs []uint64
	var collectionPerms []models.CollectionPermission
	err := r.db.Model(models.CollectionPermission{}).
		Where("collection_id = ? and permission_group = ?", collectionID, models.PGCollectionModerateSysName).
		Preload("Member").
		Preload("Role").
		Preload("Role.Members").
		Find(&collectionPerms).Error
	if err != nil {
		return nil, err
	}
	for _, perm := range collectionPerms {
		if perm.MemberId != nil {
			userIDs = append(userIDs, perm.Member.UserID)
		} else if perm.RoleId != nil {
			for _, member := range perm.Role.Members {
				userIDs = append(userIDs, member.UserID)
			}
		}
	}
	return userIDs, nil
}

func (r notificationRepo) GetUserEmailsByIDs(userIDs []uint64) (users []models.User, err error) {
	err = r.db.Model(models.User{}).Select("email").Where("id in ?", userIDs).Find(&users).Error
	return users, err
}

func (r notificationRepo) IsUserDiscordNotificationSentAnytime(userID uint64) bool {
	var count int64
	_ = r.db.Model(&models.Notification{}).Where("notifier_id = ? and discord_dm_id is not null", userID).Count(&count).Error
	return count == 0
}

func (r notificationRepo) IsUserSlackNotificationSentAnytime(userID uint64) bool {
	var count int64
	_ = r.db.Model(&models.Notification{}).Where("notifier_id = ? and slack_dm_id is not null", userID).Count(&count).Error
	return count == 0
}

func (r notificationRepo) GetMemberByUserAndStudioID(userID, studioID uint64) (*models.Member, error) {
	var member *models.Member
	err := r.db.Model(models.Member{}).Where("user_id = ? and studio_id = ?", userID, studioID).First(&member).Error
	return member, err
}

func (r notificationRepo) GetNotifications(query map[string]interface{}, skip, limit int) ([]models.Notification, error) {
	var data []models.Notification
	err := r.db.Model(models.Notification{}).Where(query).Offset(skip).Limit(limit).Order("created_at DESC").Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r notificationRepo) GetStudioIntegrationByID(id uint64) (*models.StudioIntegration, error) {
	var studioIntegration models.StudioIntegration
	err := postgres.GetDB().Where(map[string]interface{}{
		"id": id,
	}).Preload("Studio").First(&studioIntegration).Error
	return &studioIntegration, err
}

func (r notificationRepo) StudioIntegrationUpdate(integrationID uint64, updates map[string]interface{}) error {
	err := postgres.GetDB().Model(&models.StudioIntegration{}).Where("id = ?", integrationID).Updates(updates).Error
	if err != nil {
		return err
	}
	return nil
}
