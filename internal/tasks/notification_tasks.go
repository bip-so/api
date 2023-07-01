package tasks

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
)

func (s taskService) RoughBranchNotificationsOnMerge(ctx context.Context, task *asynq.Task) {
	var branch *models.CanvasBranch
	json.Unmarshal(task.Payload(), &branch)
	notifications.App.Service.ExecuteAllRoughBranchNotifications(branch.ID)
}
