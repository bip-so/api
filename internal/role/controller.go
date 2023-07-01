package role

import (
	"gitlab.com/phonepost/bip-be-platform/internal/member"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
)

var (
	RoleBasicController = roleController{}
)

// Create an studio Role Instance for a given studio id
func (rc roleController) CreateDefaultStudioRole(studioID uint64, memberId uint64) (uint64, error) {
	// Empty (Special ) Role Instance
	getEmptyDefaultAdminRole := App.Service.DefaultStudioAdminRole()

	// studio set
	getEmptyDefaultAdminRole.StudioID = studioID

	var members []models.Member
	var memberInstance *models.Member
	memberInstance = member.App.Controller.GetMemberInstance(memberId)
	members = append(members, *memberInstance)
	getEmptyDefaultAdminRole.Members = members

	sr := NewRoleRepo()
	createdRole, err := RoleRepo.CreateRole(sr, getEmptyDefaultAdminRole)
	if err != nil {
		return 0, err
	}
	//getEmptyDefaultAdminRole.Members = *[]models.Member{}
	//getEmptyDefaultAdminRole.Members
	return createdRole.ID, nil
}

func (rc roleController) CreateDefaultStudioMemberRole(studioID uint64, members []models.Member) (uint64, error) {
	memberRole := &models.Role{
		Name:     models.SYSTEM_ROLE_MEMBER,
		StudioID: studioID,
		Color:    "#ffffff",
		IsSystem: true,
		Icon:     "",
		Members:  members,
	}
	sr := NewRoleRepo()
	createdRole, err := RoleRepo.CreateRole(sr, memberRole)
	if err != nil {
		return 0, err
	}
	return createdRole.ID, nil
}

func (rc roleController) CreateBillingRole(studioID uint64, members []models.Member) (uint64, error) {
	memberRole := &models.Role{
		Name:     models.SYSTEM_BILLING_ROLE,
		StudioID: studioID,
		Color:    "#ffffff",
		IsSystem: true,
		Icon:     "",
		Members:  members,
	}
	sr := NewRoleRepo()
	createdRole, err := RoleRepo.CreateRole(sr, memberRole)
	if err != nil {
		return 0, err
	}
	return createdRole.ID, nil
}

func (rc roleController) GetMemberRoles(studioID uint64, memberID uint64) ([]RoleMembersSerializer, error) {
	memberRoles, err := App.Repo.GetMemberRolesByID(studioID, memberID)
	if err != nil {
		return nil, err
	}
	return SerializeRoleMembers(memberRoles), nil
}
