package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"gitlab.com/phonepost/bip-be-platform/internal/parser2"
)

func (s taskService) NotionImportHandler(ctx context.Context, task *asynq.Task) {
	fmt.Println("Found the correct task")
	var body parser2.ImportTask
	json.Unmarshal(task.Payload(), &body)
	file, err := body.File.Open()
	if err != nil {
		return
	}
	fmt.Println("Starting the task")
	parser2.App.Service.NotionImportZipHandler(file, body.File.Size, body.User, body.StudioID)
}
