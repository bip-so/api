package global

import (
	"fmt"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

func GetRolesByStudioID(studioID uint64, name string) ([]models.Role, error) {
	var roles []models.Role
	result := postgres.GetDB().Model(&roles).Where("studio_id = ? and name ILIKE ? ", studioID, name+"%").Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	fmt.Println(roles)
	return roles, nil
}

func StudioMembersCount(studioID uint64) int64 {
	var count int64
	postgres.GetDB().Model(models.Member{}).Where("studio_id = ?", studioID).Count(&count)
	return count
}

func GetStudioTags(studioID uint64) *models.Studio {
	var studio *models.Studio
	result := postgres.GetDB().Model(&studio).Where("studio_id = ?", studioID).Preload("Topics").Find(&studio)
	if result.Error != nil {
		fmt.Println("Error on fetching studio tags", result.Error)
		return nil
	}
	return studio
}

func FirstCanvasRepo(studioID uint64) (*models.CanvasRepository, error) {
	var repo *models.CanvasRepository
	err := postgres.GetDB().Model(models.CanvasRepository{}).Where("studio_id = ?", studioID).Order("created_at ASC").First(&repo).Error
	return repo, err
}
