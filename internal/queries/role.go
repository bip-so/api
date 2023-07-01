package queries

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
)

func (q *roleQuery) EmptyRoleObject() *models.Role {
	return &models.Role{}
}

func (m *roleQuery) GetRole(roleId uint64) (*models.Role, error) {
	var role *models.Role
	postgres.GetDB().Model(&models.Role{}).Where("id = ?", roleId).First(&role)
	return role, nil
}

func (q *roleQuery) GetStudioMemberRole(studioID uint64) *models.Role {
	var memberRole *models.Role
	postgres.GetDB().Model(&models.Role{}).Where("name = ? and studio_id = ?", "Member", studioID).First(&memberRole)
	return memberRole
}

func (m *roleQuery) GetStudioAdminRole(studioID uint64) (role *models.Role, err error) {
	err = postgres.GetDB().Model(models.Role{}).
		Where("studio_id = ? and name = ?", studioID, models.SYSTEM_ADMIN_ROLE).
		Preload("Members").
		Preload("Members.User").
		First(&role).Error
	return
}

func (q *roleQuery) GetRoleMembers(roleID uint64, skip, limit int) ([]models.RoleMember, error) {
	roleMembers := []models.RoleMember{}
	postgres.GetDB().Table("role_members").Where("role_id = ?", roleID).Offset(skip).Limit(limit).Find(&roleMembers)
	//postgres.GetDB().Model(&role).Where(query).Preload("Members").Preload("Members.User").First(&role)
	return roleMembers, nil
}
