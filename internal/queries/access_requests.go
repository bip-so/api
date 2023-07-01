package queries

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

type CreateAccessRequestPost struct {
	StudioID           uint64 `json:"studioID"`
	CollectionID       uint64 `json:"collectionID"`
	CanvasRepositoryID uint64 `json:"canvasRepositoryID"`
	CanvasBranchID     uint64 `json:"canvasBranchID"`
}

func (r accessRequestQuery) AccessRequestInstance(arID uint64) (*models.AccessRequest, error) {
	var instances *models.AccessRequest
	err := postgres.GetDB().Model(&models.AccessRequest{}).Where("id = ?", arID).Preload("CanvasRepository").Find(&instances).Error
	if err != nil {
		return nil, err
	}
	return instances, nil
}

func (r accessRequestQuery) AccessRequestExists(data CreateAccessRequestPost, userID uint64) bool {
	var count int64
	_ = postgres.GetDB().Model(&models.AccessRequest{}).Where("canvas_repository_id = ? and canvas_branch_id = ? and created_by_id = ? and status = ?", data.CanvasRepositoryID, data.CanvasBranchID, userID, models.ACCESS_REQUEST_PENDING).Count(&count).Error
	if count == 0 {
		return false
	}
	return true
}

func (r accessRequestQuery) AccessRequestExistsSimple(branchID uint64, userID uint64) bool {
	var count int64
	_ = postgres.GetDB().Model(&models.AccessRequest{}).Where("canvas_branch_id = ? and created_by_id = ? , status = ?", branchID, userID, models.ACCESS_REQUEST_PENDING).Count(&count).Error
	if count == 0 {
		return false
	}
	return true
}

func (r accessRequestQuery) CreateAccessRequest(instance *models.AccessRequest) (*models.AccessRequest, error) {
	results := postgres.GetDB().Create(&instance)
	return instance, results.Error
}

func (r accessRequestQuery) GetAllAccessRequests(query map[string]interface{}) (*[]models.AccessRequest, error) {
	var instances []models.AccessRequest
	err := postgres.GetDB().Model(&models.AccessRequest{}).Where(query).Find(&instances).Error
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &instances, nil
}

func (r accessRequestQuery) UpdateAccessRequests(id uint64, pg string, status string) error {
	//var instances []models.AccessRequest
	err := postgres.GetDB().Model(&models.AccessRequest{}).Where("id = ?", id).Updates(map[string]interface{}{
		"canvas_branch_permission_group": pg,
		"status":                         status,
	}).Error
	if err != nil {
		logger.Debug(err.Error())
		return err
	}
	return nil
}
