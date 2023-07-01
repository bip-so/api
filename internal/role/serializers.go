package role

import (
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/phonepost/bip-be-platform/internal/member"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

type RoleSerializer struct {
	ID         uint64                    `json:"id"`
	UUID       string                    `json:"uuid"`
	StudioID   uint64                    `json:"studioID"`
	Name       string                    `json:"name"`
	Color      string                    `json:"color"`
	IsSystem   bool                      `json:"isSystem"`
	Icon       string                    `json:"icon"`
	Members    []member.MemberSerializer `json:"members"`
	IsNonPerms bool                      `json:"isNonPerms"`
	CreatedAt  time.Time                 `json:"createdAt"`
	UpdatedAt  time.Time                 `json:"updatedAt"`
}

func SerializeRole(role *models.Role) *RoleSerializer {
	view := RoleSerializer{
		ID:         role.ID,
		UUID:       role.UUID.String(),
		StudioID:   role.StudioID,
		Name:       role.Name,
		Color:      role.Color,
		IsSystem:   role.IsSystem,
		Icon:       role.Icon,
		IsNonPerms: role.IsNonPerms,
	}

	if role.Members != nil {
		for _, mem := range role.Members {
			view.Members = append(view.Members, *member.SerializeMember(&mem))
		}
	}
	return &view
}

type RoleGenericSerializer struct {
	c *gin.Context
}

func (self *RoleGenericSerializer) GetStudioSerializer(role *models.Role) *RoleSerializer {
	view := RoleSerializer{
		ID:         role.ID,
		UUID:       role.UUID.String(),
		StudioID:   role.StudioID,
		Name:       role.Name,
		Color:      role.Color,
		IsSystem:   role.IsSystem,
		Icon:       role.Icon,
		IsNonPerms: role.IsNonPerms,
		CreatedAt:  role.CreatedAt,
		UpdatedAt:  role.UpdatedAt,
	}
	return &view
}

type RoleMembersSerializer struct {
	MemberID   uint64 `json:"memberId"`
	Id         uint64 `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	IsSystem   bool   `json:"isSystem"`
	IsNonPerms bool   `json:"isNonPerms"`
}

func SerializeRoleMembers(roleMembers []RoleMembersSerializer) []RoleMembersSerializer {
	var roles []RoleMembersSerializer
	for _, roleMember := range roleMembers {
		roles = append(roles, roleMember)
	}
	return roles
}
