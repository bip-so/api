package queries

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

func (q *studioPermissionQuery) GetStudioPermission(query map[string]interface{}) (models.StudioPermission, error) {
	var studioPerms models.StudioPermission
	err := postgres.GetDB().Model(&studioPerms).Where(query).Find(&studioPerms).Error
	return studioPerms, err
}
