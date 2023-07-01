package feed

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

func (fr feedRepo) GetCanvasBranchPermission(query map[string]interface{}) ([]models.CanvasBranchPermission, error) {
	var canvasBranchPerms []models.CanvasBranchPermission
	postgres.GetDB().
		Model(&canvasBranchPerms).
		Where(query).
		Preload("Studio").
		Preload("Role").
		Preload("Role.Members").
		Preload("Member").
		Find(&canvasBranchPerms)
	return canvasBranchPerms, nil
}
