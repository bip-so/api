package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/reel"
)

func (s taskService) AddReelToAlgolia(ctx context.Context, task *asynq.Task) {
	var reelInstance *models.Reel
	json.Unmarshal(task.Payload(), &reelInstance)
	err := reel.App.Service.AddReelToAlgolia(reelInstance.ID)
	if err != nil {
		fmt.Println("Error adding reel to algolia", err.Error())
	}
}

func (s taskService) DeleteReelFromAlgolia(ctx context.Context, task *asynq.Task) {
	var reelInstance *models.Reel
	json.Unmarshal(task.Payload(), &reelInstance)
	err := reel.App.Service.DeleteReelFromAlgolia(reelInstance.ID)
	if err != nil {
		fmt.Println("Error deleting reel to algolia", err.Error())
	}
}
