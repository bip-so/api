package role

import (
	"fmt"

	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg/stores/postgres"
	"gorm.io/gorm/clause"
)

type RoleRepo interface {
	CreateRole(role *models.Role) (*models.Role, error)
	DeleteRole(roleId uint64) error
	GetRole(roleId uint64) (*models.Role, error)
	UpdateRole(ID uint64, updates map[string]interface{}) error
	GetRolesByStudioID(studioId uint64) ([]models.Role, error)
	FindMembersInRole(role *models.Role) (members []models.Member, err error)
	AddMembersInRole(members []models.Member, role *models.Role) error
	RemoveMembersInRole(members []models.Member, role *models.Role) error
}

func NewRoleRepo() RoleRepo {
	return &roleRepo{}
}

func (sr roleRepo) CreateRole(role *models.Role) (*models.Role, error) {
	result := postgres.GetDB().Create(&role)
	if result.Error != nil {
		return nil, result.Error
	}
	return role, nil
}

func (sr roleRepo) GetRoleByID(roleId uint64) (*models.Role, error) {
	var role *models.Role
	postgres.GetDB().Model(&models.Role{}).Where("id = ?", roleId).Preload("Studio").First(&role)
	return role, nil
}

func (sr roleRepo) GetRole(roleId uint64) (*models.Role, error) {
	var role *models.Role
	postgres.GetDB().Model(&models.Role{}).Where("id = ?", roleId).Preload("Members").Preload("Studio").Preload("Members.User").First(&role)
	return role, nil
}

func (sr roleRepo) GetRolesByIDs(roleIDs []uint64) ([]models.Role, error) {
	var roles []models.Role
	err := postgres.GetDB().Model(&models.Role{}).Where("id in ?", roleIDs).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (sr roleRepo) GetRolesByStudioID(studioID uint64) ([]models.Role, error) {
	var roles []models.Role
	result := postgres.GetDB().Model(&roles).Where("studio_id = ?", studioID).Preload("Members").Preload("Studio").Preload("Members.User").Order("id asc").Find(&roles)

	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}

func (sr roleRepo) DeleteRole(roleId uint64) error {
	roleObject := models.Role{BaseModel: models.BaseModel{ID: roleId}}
	//err := postgres.GetDB().Delete(roleObject).Error
	err := postgres.GetDB().Unscoped().Select(clause.Associations).Delete(roleObject).Error
	return err
}

func (sr roleRepo) UpdateRole(ID uint64, updates map[string]interface{}) error {
	err := postgres.GetDB().Model(&models.Role{}).Where("id = ?", ID).Updates(updates).Error
	return err
}

func (sr roleRepo) FindMembersInRole(role *models.Role) (members []models.Member, err error) {
	err = postgres.GetDB().Model(&role).Association("Members").Find(&members)
	return members, err
}

func (sr roleRepo) AddMembersInRole(addMembers []models.Member, role *models.Role) error {
	//err := postgres.GetDB().Model(&role).Association("Members").Append(addMembers)
	query := `insert into role_members (member_id, role_id) values`
	for i, member := range addMembers {
		query = fmt.Sprintf("%s (%d,%d)", query, member.ID, role.ID)
		if i < (len(addMembers) - 1) {
			query += ","
		}
	}
	println("Query====>  ", query)
	err := postgres.GetDB().Exec(query).Error
	if err != nil {
		fmt.Println("Error in creating role_members: ", err)
	}
	return err
}

func (sr roleRepo) RemoveMembersInRole(removeMembers []models.Member, role *models.Role) error {
	err := postgres.GetDB().Model(&role).Association("Members").Delete(removeMembers)
	return err
}

func (sr roleRepo) GetMemberRolesByID(studioID uint64, memberID uint64) ([]RoleMembersSerializer, error) {
	var roles []RoleMembersSerializer
	result := postgres.GetDB().Raw(`
		select * from roles 
		left join role_members on role_members.role_id = roles.id
		where studio_id = ? and role_members.member_id = ?
	`, studioID, memberID).Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}

// This adds a given member to a Studio for Memeber Role
func (sr roleRepo) AddMembersToMemberRoleForStudio(studioID uint64, memberID uint64) error {
	// Get SYSTEM_ROLE_MEMBER Role
	var role *models.Role
	fmt.Println(studioID, memberID, role)
	roleInstanceError := sr.db.Model(&models.Role{}).Where("studio_id = ? and name = ?", studioID, models.SYSTEM_ROLE_MEMBER).Preload("Members").First(&role).Error
	if roleInstanceError != nil {
		return roleInstanceError
	}
	// Get Member
	var member *models.Member
	errGettingMember := sr.db.Model(&models.Member{}).Where("id = ?", memberID).First(&member).Error
	if errGettingMember != nil {
		return errGettingMember
	}
	// Add Member to Role
	err := sr.db.Model(&role).Association("Members").Append(&member)
	if err != nil {
		return err
	}
	return nil
}

func (sr roleRepo) GetMembersByIDs(memberIds []uint64) ([]models.Member, error) {
	var members []models.Member
	result := postgres.GetDB().Model(&models.Member{}).Where("id in ?", memberIds).Preload("User").Find(&members)

	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}
