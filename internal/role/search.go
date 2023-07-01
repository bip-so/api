package role

import (
	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

type RoleDocument struct {
	ID         uint64 `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	IsSystem   bool   `json:"isSystem"`
	Icon       string `json:"icon"`
	IsNonPerms bool   `json:"isNonPerms"`
}

func RoleModelToUserDocument(role *models.Role) *RoleDocument {
	return &RoleDocument{
		ID:         role.ID,
		Name:       role.Name,
		Color:      role.Color,
		IsSystem:   role.IsSystem,
		Icon:       role.Icon,
		IsNonPerms: role.IsNonPerms,
	}
}
func GetRoleSearch(roles *[]models.Role) []RoleDocument {
	var rd []RoleDocument
	for _, role := range *roles {
		rd2 := RoleModelToUserDocument(&role)
		rd = append(rd, *rd2)
	}
	return rd
}
