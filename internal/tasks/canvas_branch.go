package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/canvasrepo"
	"gitlab.com/phonepost/bip-be-platform/internal/permissions"
	cache "gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"time"

	"github.com/hibiken/asynq"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
)

func (s taskService) DeleteMergeRequestNotifications(ctx context.Context, task *asynq.Task) {
	notifications.App.Service.MergeRequestDeleteHandler(utils.Uint64((string(task.Payload()))))
}

func (s taskService) DeletePublishRequestNotifications(ctx context.Context, task *asynq.Task) {
	var prInstance *models.PublishRequest
	json.Unmarshal(task.Payload(), &prInstance)
	notifications.App.Service.PublishRequestDeleteHandler(prInstance)
}

func (s taskService) DeleteModsOnCanvas(ctx context.Context, task *asynq.Task) {
	var prInstance *models.PublishRequest
	json.Unmarshal(task.Payload(), &prInstance)
	userIDs, roleIDs := notifications.App.Service.DeleteModUsersOnCanvas(prInstance)
	var canvasRepo models.CanvasRepository
	err := App.Repo.db.Model(models.CanvasRepository{}).Where("id = ?", prInstance.CanvasRepositoryID).Preload("Studio").First(&canvasRepo).Error
	if err != nil {
		fmt.Println("Error in getting canvas repo", err)
		return
	}
	for _, userId := range userIDs {
		permissions.App.Service.InvalidateCanvasPermissionCache(userId, canvasRepo.StudioID, canvasRepo.CollectionID)
	}
	for _, roleId := range roleIDs {
		permissions.App.Service.InvalidateCanvasPermissionCacheByRole(roleId, canvasRepo.StudioID, canvasRepo.CollectionID)
	}
}

func (s taskService) CanvasBranchAccessCron(ctx context.Context, task *asynq.Task) {
	yesterdayDateAsKey := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	cacheInstance := cache.NewCache()
	iter := cacheInstance.GetAllMatchingKeys(ctx, fmt.Sprintf("plans-auto-publish:%s:*", yesterdayDateAsKey))
	for iter.Next(ctx) {
		value := cacheInstance.Get(ctx, iter.Val())
		canvasrepo.App.PlansAutoPublish.ProcessCanvasBranchAccess(value.(string))
	}
}
