package blockthread

import (
	"github.com/google/uuid"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/internal/notifications"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
)

func (s blockThreadService) InitThreadInstance(
	userID uint64, canvasRepositoryID uint64, canvasBranchID uint64, startBlockID uint64,
	position uint, textRangeStart uint, textRangeEnd uint, text string, highlightedText string,
	startBlockUUID uuid.UUID,
) *models.BlockThread {
	return &models.BlockThread{
		CreatedByID:        userID,
		UpdatedByID:        userID,
		CanvasRepositoryID: canvasRepositoryID,
		CanvasBranchID:     canvasBranchID,
		StartBlockID:       startBlockID,
		StartBlockUUID:     startBlockUUID,
		Position:           position,
		TextRangeStart:     textRangeStart,
		TextRangeEnd:       textRangeEnd,
		Text:               text,
		HighlightedText:    highlightedText,
	}
}
func (r blockThreadRepo) GetBlockByUUIDAndBranchID(query map[string]interface{}) (*models.Block, error) {
	var block models.Block
	err := r.db.Model(&models.Block{}).Where(query).First(&block).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &block, nil
}

func (s blockThreadService) Create(body *PostBlockThread, userID uint64) (*models.BlockThread, error) {
	blockInstance, _ := App.Repo.GetBlockByUUIDAndBranchID(map[string]interface{}{"uuid": body.StartBlockUUID, "canvas_branch_id": body.CanvasBranchID})

	instance := s.InitThreadInstance(
		userID, body.CanvasRepositoryID, body.CanvasBranchID, blockInstance.ID,
		body.Position, body.TextRangeStart, body.TextRangeEnd, body.Text, body.HighlightedText,
		body.StartBlockUUID,
	)
	created, err := App.Repo.Create(instance)
	// When BlockThread is created we have to update Block Instance too
	go func() {
		extraData := notifications.NotificationExtraData{
			CanvasRepoID:   created.CanvasRepositoryID,
			CanvasBranchID: created.CanvasBranchID,
		}
		contentObject := models.BLOCK_THREAD
		notifications.App.Service.PublishNewNotification(notifications.BlockComment,
			userID, nil, nil, nil, extraData, &created.ID, &contentObject)
	}()
	return created, err
}

func (s blockThreadService) Update(body *PatchBlockThread, userID uint64) error {

	updates := map[string]interface{}{
		"updated_by_id":        userID,
		"canvas_repository_id": body.CanvasRepositoryID,
		"canvas_branch_id":     body.CanvasBranchID,
		"start_block_id":       body.StartBlockID,
		"position":             body.Position,
		"text_range_start":     body.TextRangeStart,
		"text_range_end":       body.TextRangeEnd,
		"text":                 body.Text,
		"highlighted_text":     body.HighlightedText,
	}
	err := App.Repo.Update(body.ID, updates)
	return err
}
