package queries

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/logger"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/redis"
	"gitlab.com/phonepost/bip-be-platform/pkg/utils"
	"strconv"
	"time"
)

func (q *branchQuery) EmptyBranchObject() *models.CanvasBranch {
	return &models.CanvasBranch{}
}
func (q *branchQuery) CreateBranch(userID uint64, canvasRepoID uint64, name string, isTrue bool, parentPublicAccess string) (*models.CanvasBranch, error) {
	instance := q.EmptyBranchObject()
	instance.CreatedByID = userID
	instance.UpdatedByID = userID
	instance.Name = name
	instance.CanvasRepositoryID = canvasRepoID
	instance.IsDefault = isTrue
	instance.PublicAccess = parentPublicAccess
	instance.Key = utils.NewNanoid()
	results := postgres.GetDB().Create(&instance)
	return instance, results.Error
}

func (q *branchQuery) GetBranch(query map[string]interface{}, preloadWithRepo bool) (*models.CanvasBranch, error) {
	var branch models.CanvasBranch
	var err error
	// This is stupid will think later
	if preloadWithRepo {
		err = postgres.GetDB().Model(&models.CanvasBranch{}).Where(query).Preload("CanvasRepository").First(&branch).Error
	} else {
		err = postgres.GetDB().Model(&models.CanvasBranch{}).Where(query).Preload("CanvasRepository").First(&branch).Error
	}
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &branch, nil
}

const BranchPreLoader = "branches-preload:"

func (r *branchQuery) GetBranchWithRepoAndStudio(id uint64) (*models.CanvasBranch, error) {
	var branch models.CanvasBranch
	branchIDStr := strconv.FormatUint(id, 10)
	key := BranchPreLoader + branchIDStr
	val, _ := redis.RedisClient().Get(context.Background(), key).Result()
	if val == "" {
		fmt.Println("Sending DB Value: " + key)
		err := postgres.GetDB().Model(&models.CanvasBranch{}).Where("id = ?", id).Preload("CanvasRepository").Preload("CanvasRepository.Studio").First(&branch).Error
		if err != nil {
			logger.Debug(err.Error())
			return nil, err
		}
		branchObjectJson, _ := json.Marshal(branch)
		_ = redis.RedisClient().Set(context.Background(), key, branchObjectJson, 180*time.Second).Err()
	} else {
		fmt.Println("Sending Cached Value: " + key)
		val2, _ := redis.RedisClient().Get(context.Background(), key).Result()
		_ = json.Unmarshal([]byte(val2), &branch)
	}
	defer utils.TimeTrack(time.Now())
	return &branch, nil

}

const BranchFullPreLoader = "branches-rb-fb-repo-preload:"

func (r *branchQuery) GetBranchByIDWithRBFBStudioParentRepo(branchID uint64) (*models.CanvasBranch, error) {

	var branch models.CanvasBranch
	branchIDStr := strconv.FormatUint(branchID, 10)
	key := BranchFullPreLoader + branchIDStr

	val, _ := redis.RedisClient().Get(context.Background(), key).Result()
	if val == "" {
		fmt.Println("GetBranchByIDWithRBFBStudioParentRepo / DB Value: " + key)
		//err := postgres.GetDB().Model(&models.CanvasBranch{}).Where("id = ?", id).Preload("CanvasRepository").Preload("CanvasRepository.Studio").First(&branch).Error
		err := postgres.GetDB().Model(&models.CanvasBranch{}).Where("id = ?", branchID).
			Preload("RoughFromBranch").Where("id = ?", branchID).
			Preload("FromBranch").
			Preload("CanvasRepository.Studio").
			Preload("CanvasRepository.ParentCanvasRepository").
			First(&branch).
			Error
		if err != nil {
			logger.Debug(err.Error())
			return nil, err
		}
		branchObjectJson, _ := json.Marshal(branch)
		_ = redis.RedisClient().Set(context.Background(), key, branchObjectJson, 120*time.Second).Err()
	} else {
		fmt.Println("GetBranchByIDWithRBFBStudioParentRepo / Cached Value: " + key)
		val2, _ := redis.RedisClient().Get(context.Background(), key).Result()
		_ = json.Unmarshal([]byte(val2), &branch)
	}

	defer utils.TimeTrack(time.Now())
	return &branch, nil
}

// Should not be used.
// Use: BranchFullPreLoader instead
func (r branchQuery) GetBranchByID(branchID uint64) (*models.CanvasBranch, error) {
	//return queries.App.BranchQuery.GetBranchByID.GetBranchByIDWithRBFBStudioParentRepo(branchID)
	//
	var branch models.CanvasBranch
	err := postgres.GetDB().Model(&models.CanvasBranch{}).Where("id = ?", branchID).
		Preload("RoughFromBranch").Where("id = ?", branchID).
		Preload("FromBranch").
		Preload("CanvasRepository.Studio").
		Preload("CanvasRepository.ParentCanvasRepository").
		First(&branch).
		Error

	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return &branch, nil
}

// GetPendingEmails : Returns list of unclaimed emails for a branch
func (r branchQuery) GetPendingEmails(id uint64) []string {
	var emails []string
	var invited *[]models.BranchInviteViaEmail
	_ = postgres.GetDB().Model(&models.BranchInviteViaEmail{}).Where("branch_id = ?", id).Find(&invited).Error
	for _, val := range *invited {
		emails = append(emails, val.Email)
	}
	return emails
}

func (r branchQuery) UpdateBranchInstance(branchID uint64, query map[string]interface{}) error {
	err := postgres.GetDB().Model(&models.CanvasBranch{}).Where("id = ?", branchID).Updates(query).Error
	if err != nil {
		return err
	}
	return nil
}
