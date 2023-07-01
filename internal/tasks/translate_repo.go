package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"gitlab.com/phonepost/bip-be-platform/internal/blockthread"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/aws"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"strings"
)

func (s taskService) TranslateCanvasRepositories(ctx context.Context, task *asynq.Task) {

	var canvasRepos []models.CanvasRepository
	json.Unmarshal(task.Payload(), &canvasRepos)

	repoIDs := []uint64{}
	for _, repo := range canvasRepos {
		repoIDs = append(repoIDs, repo.ID)
	}
	var languagePageRepos []models.CanvasRepository
	err := postgres.GetDB().Model(models.CanvasRepository{}).Where("id in ?", repoIDs).Find(&languagePageRepos).Error
	if err != nil {
		fmt.Println("error in getting language pages", err)
		return
	}

	for _, lrepo := range languagePageRepos {
		mainCanvasRepo, err := queries.App.RepoQuery.GetRepo(map[string]interface{}{"id": lrepo.DefaultLanguageCanvasRepoID})
		fmt.Println("mainCanavsRepo", mainCanvasRepo.ID, mainCanvasRepo.DefaultLanguageCanvasRepoID)
		blocks, err := queries.App.BlockQuery.GetBlocksByBranchID(*mainCanvasRepo.DefaultBranchID)
		if err != nil {
			fmt.Println("Error in getting blocks", err)
		}
		srcLanguage := "en"
		if mainCanvasRepo.Language != nil {
			srcLanguage = *mainCanvasRepo.Language
		}
		if srcLanguage == "" {
			srcLanguage = "en"
		}
		translatedCanvasName, err := aws.Translate(lrepo.Name, srcLanguage, *lrepo.Language)
		lrepo.Name = translatedCanvasName
		postgres.GetDB().Save(&lrepo)
		for _, block := range *blocks {
			var childrenData []map[string]interface{}
			json.Unmarshal(block.Children, &childrenData)
			for i, children := range childrenData {
				if block.Type == "code" {
					continue
				}
				if children["text"] != nil && children["text"].(string) != "" {
					text := children["text"].(string)
					translatedText, err := aws.Translate(text, srcLanguage, *lrepo.Language)
					if err != nil {
						fmt.Println("Error in translating", err, text, translatedText, "Src language", srcLanguage, "Target Language:", *lrepo.Language)
					}
					childrenData[i]["text"] = translatedText
				}
			}
			parentBlockID := block.ID
			blockChildren, _ := json.Marshal(childrenData)
			newBlock := block
			newBlock.ID = 0
			newBlock.UUID = uuid.New()
			newBlock.CanvasBranchID = lrepo.DefaultBranchID
			newBlock.Children = blockChildren
			newBlock.CanvasRepositoryID = lrepo.ID
			newBlock.ParentBlock = &parentBlockID
			err = postgres.GetDB().Create(&newBlock).Error
			if err != nil {
				fmt.Println("Error in creating newBlock", block.ID, newBlock.ID)
				continue
			}
			queries.App.BlockQuery.CopyBlockCommentsReelsReactionsToNewBlock(block.ID, newBlock.ID, *lrepo.DefaultBranchID)
		}
		blockThreads, err := blockthread.App.Repo.GetAllThread(map[string]interface{}{"canvas_branch_id": *lrepo.DefaultBranchID})
		if err != nil {
			fmt.Println("Error in getting block threads", err)
		}
		for _, thread := range *blockThreads {
			block, _ := queries.App.BlockQuery.GetBlock(map[string]interface{}{"id": thread.StartBlockID})
			var childrenData []map[string]interface{}
			json.Unmarshal(block.Children, &childrenData)
			for i, children := range childrenData {
				for key, val := range children {
					if strings.Contains(key, "commentThread") {
						newCommentKey := "commentThread_" + thread.UUID.String()
						childrenData[i][newCommentKey] = val
					}
				}
			}
			updatedChildren, _ := json.Marshal(childrenData)
			block.Children = updatedChildren
			postgres.GetDB().Save(&block)
			thread.StartBlockUUID = block.UUID
			postgres.GetDB().Save(&thread)
		}
		reels, err := blockthread.App.Repo.GetAllReels(map[string]interface{}{"canvas_branch_id": *lrepo.DefaultBranchID})
		if err != nil {
			fmt.Println("Error in getting reels", err)
		}
		for _, reel := range reels {
			block, _ := queries.App.BlockQuery.GetBlock(map[string]interface{}{"id": reel.StartBlockID})
			reel.StartBlockUUID = block.UUID
			postgres.GetDB().Save(&reel)
		}
		lrepo.IsProcessing = false
		postgres.GetDB().Save(&lrepo)
		notifications.App.Service.PublishNewNotification(notifications.TranslateCanvas, 0, []uint64{lrepo.CreatedByID}, &lrepo.StudioID,
			nil, notifications.NotificationExtraData{
				CanvasBranchID: *lrepo.DefaultBranchID,
				CanvasRepoID:   lrepo.ID,
				CollectionID:   lrepo.CollectionID,
			}, nil, nil)
	}
}
