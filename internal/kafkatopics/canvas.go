package kafkatopics

import (
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"gitlab.com/phonepost/bip-be-platform/internal/canvasrepo"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/queries"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
)

// CreateNewCanvas Triggers from kafka topics When a new studio is created.
// Commented by CC
//func CreateNewCanvas(msg *kafka.Message, collectionInstance *models.Collection) {
//
//	canvasView := canvasrepo.InitCanvasRepoPost{
//		CollectionID: collectionInstance.ID,
//		Name:         "Default Canvas",
//		Icon:         "📋",
//		Position:     1,
//	}
//	repo, err := canvasrepo.App.Controller.InitCanvasRepo(canvasView, collectionInstance.CreatedByID, collectionInstance.StudioID, )
//	if err != nil {
//		logger.Error(err.Error())
//		KafkaConsumerError(msg, err)
//		return
//	}
//
//	logger.Debug(fmt.Sprintf(
//		"Collection ID %d And canvas ID %d are created for studio ID %d",
//		collectionInstance.ID, repo.ID, collectionInstance.StudioID))
//}

func UpdateCollectionCanvasCount(msg *kafka.Message, created bool) {
	var canvasRepo models.CanvasRepository
	err := json.Unmarshal(msg.Value, &canvasRepo)
	if err != nil {
		logger.Error(err.Error())
		KafkaConsumerError(msg, err)
		return
	}

	var canvasRootRepoCount int64
	var canvasAllReposCount int64
	canvasRootRepoCount, _ = canvasrepo.App.Repo.GetCanvasReposCount(map[string]interface{}{"collection_id": canvasRepo.CollectionID, "parent_canvas_repository_id": nil})
	canvasAllReposCount, _ = canvasrepo.App.Repo.GetCanvasReposCount(map[string]interface{}{"collection_id": canvasRepo.CollectionID})

	if canvasRepo.ParentCanvasRepositoryID != nil {
		subCanvasRootRepoCount, _ := canvasrepo.App.Repo.GetCanvasReposCount(map[string]interface{}{"collection_id": canvasRepo.CollectionID, "parent_canvas_repository_id": *canvasRepo.ParentCanvasRepositoryID})
		err = canvasrepo.App.Repo.Manager.UpdateEntityByID(models.CANVAS_REPO, *canvasRepo.ParentCanvasRepositoryID, map[string]interface{}{"sub_canvas_count": int(subCanvasRootRepoCount)})
		if err != nil {
			logger.Error(err.Error())
			KafkaConsumerError(msg, err)
			return
		}
	}
	_, err = queries.App.CollectionQuery.UpdateCollection(
		canvasRepo.CollectionID, map[string]interface{}{
			"computed_root_canvas_count": int(canvasRootRepoCount),
			"computed_all_canvas_count":  int(canvasAllReposCount),
		})
	if err != nil {
		logger.Error(err.Error())
		KafkaConsumerError(msg, err)
		return
	}
	// Send updated collection to discord
	canvasrepo.App.Service.SendCollectionTreeToDiscord(canvasRepo.CollectionID)
}
