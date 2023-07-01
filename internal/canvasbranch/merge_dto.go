package canvasbranch

import (
	"fmt"
	"strconv"
	"strings"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
)

// DeleteBlocks: Deletes a Multiple Block Ids passed as []uint64
func (r canvasBranchRepo) DeleteBlocks(blockIDs []uint64) error {
	err := r.db.Delete(&models.Block{}, "id IN ?", blockIDs).Error
	if err != nil {
		return err
	}
	return nil
}

func (r canvasBranchRepo) DeleteBlocksInBranchID(canvasBranchID uint64) error {
	err := r.db.Delete(&models.Block{}, "canvas_branch_id = ?", canvasBranchID).Error
	if err != nil {
		return err
	}
	return nil
}

func (r canvasBranchRepo) SyncBlockRanks(ranksMap map[string]int32) error {
	ranks := []string{}
	for uuid, rank := range ranksMap {
		ranks = append(ranks, "('"+uuid+"', "+strconv.Itoa(int(rank))+")")
	}
	valuesStr := strings.Join(ranks, ", ")
	query := `update public.blocks as b 
	set rank = c.rank
	from (values` + valuesStr + `) as c(uuid, rank) 
	where c.uuid = b.uuid::text;`
	return r.db.Exec(query).Error
}

func (cr canvasBranchRepo) UpdateBlockMerge(blockID uint64, updates map[string]interface{}) error {
	err := cr.db.Model(&models.Block{}).Where("id = ?", blockID).Updates(updates).Error
	return err
}

func (cr canvasBranchRepo) UpdateBlockSimple(block models.Block) error {
	err := cr.db.Save(block).Error
	return err
}

// Delete (Hard) MR
func (cr canvasBranchRepo) DeleteMergeRequest(mergeRequestID uint64) error {
	err := cr.Manager.HardDeleteByID(models.MERGEREQUEST, mergeRequestID)
	if err != nil {
		return err
	}
	return nil
}

func (r canvasBranchRepo) GetMergeRequest(query map[string]interface{}) (*models.MergeRequest, error) {
	fmt.Println("calls")
	fmt.Println(query)
	var mr models.MergeRequest
	err := r.db.Debug().Table(models.MERGEREQUEST).Preload("CanvasRepository.Studio").Preload("SourceBranch").Preload("DestinationBranch").Where(query).First(&mr).Error
	if err != nil {
		return nil, err
	}
	return &mr, nil
}
func (r canvasBranchRepo) GetMergeRequestWithPreloads(query map[string]interface{}) (*models.MergeRequest, error) {
	fmt.Println("calls")
	fmt.Println(query)
	var mr models.MergeRequest
	err := r.db.Debug().Table(models.MERGEREQUEST).Preload("CanvasRepository.Studio").Preload("SourceBranch").Preload("DestinationBranch").Where(query).First(&mr).Error
	if err != nil {
		return nil, err
	}
	return &mr, nil
}

func (r canvasBranchRepo) UpdateMergeRequest(mrID uint64, query map[string]interface{}) error {
	err := r.db.Model(&models.MergeRequest{}).Where("id = ?", mrID).Updates(query).Error
	if err != nil {
		return err
	}
	return nil
}

func (r canvasBranchRepo) CreateMergeRequest(instance *models.MergeRequest) (*models.MergeRequest, error) {
	results := r.db.Create(&instance)
	return instance, results.Error
}

func (r canvasBranchRepo) GetAllMergeRequests(query map[string]interface{}) (*[]models.MergeRequest, error) {
	var instances []models.MergeRequest
	err := r.db.Model(&models.MergeRequest{}).Where(query).Find(&instances).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &instances, nil
}

// Check for any open request are there while creating a new one.
func (r canvasBranchRepo) MergeRequestExists(branchID uint64, userID uint64) bool {
	var count int64
	_ = r.db.Model(&models.MergeRequest{}).Where("source_branch_id = ? and created_by_id = ? and status = ?", branchID, userID, models.MERGE_REQUEST_OPEN).Count(&count).Error
	if count == 0 {
		return false
	}
	return true
}

func (r canvasBranchRepo) MergeRequestCount(canvasRepoID uint64) int64 {
	var count int64
	_ = r.db.Model(&models.MergeRequest{}).Where("canvas_repository_id = ? and status = ?", canvasRepoID, models.MERGE_REQUEST_OPEN).Count(&count).Error
	if count == 0 {
		return 0
	}
	return count
}
